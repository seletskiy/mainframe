package engine

import (
	"runtime"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/reconquest/karma-go"
	"github.com/seletskiy/mainframe/pkg/fonts"
	"github.com/seletskiy/mainframe/pkg/log"
	"github.com/seletskiy/mainframe/pkg/protocol/messages"
)

const (
	DefaultWindowWidth  = 640
	DefaultWindowHeight = 480
)

type Delegate struct {
	barrier  chan struct{}
	callback func()
}

type Engine struct {
	vertices struct {
		points []float32

		buffers struct {
			triangles  uint32
			glyphs     uint32
			attributes uint32
			colors     uint32
		}
	}

	shaders struct {
		program  uint32
		vertex   uint32
		fragment uint32
	}

	contexts map[*Screen]*Context

	queue struct {
		sync.Mutex
		screens map[*Screen]bool
	}

	font struct {
		handle  *fonts.Font
		texture uint32
	}

	delegates chan Delegate

	running bool
}

func New() *Engine {
	engine := &Engine{
		delegates: make(chan Delegate, 0),
	}

	engine.contexts = map[*Screen]*Context{}
	engine.queue.screens = map[*Screen]bool{}

	return engine
}

func (engine *Engine) Init() error {
	// Required for OpenGL to work.
	runtime.LockOSThread()

	err := gl.Init()
	if err != nil {
		return karma.Format(
			err,
			"{gl} unable to init",
		)
	}

	err = glfw.Init()
	if err != nil {
		return karma.Format(
			err,
			"{glfw} unable to init",
		)
	}

	engine.running = true

	return nil
}

func (engine *Engine) CreateWindow(options *messages.Open) (*Context, error) {
	var (
		err     error
		context *Context
	)

	var (
		width  = DefaultWindowWidth
		height = DefaultWindowHeight
	)

	if options.Width != nil {
		width = int(*options.Width)
	}

	if options.Height != nil {
		height = int(*options.Height)
	}

	if options.Columns != nil {
		width = *options.Columns * engine.font.handle.GetWidth()
	}

	if options.Rows != nil {
		height = *options.Rows * engine.font.handle.GetHeight()
	}

	// All GL commands should be evaluated in same system thread, so we
	// need to send them to main engine thread to execute.
	engine.delegate(
		func() {
			// TODO return error
			context = engine.createWindow(width, height, options)
		},
	)
	if err != nil {
		return nil, karma.Format(
			err,
			"{glfw} unable to create window",
		)
	}

	return context, nil
}

func (engine *Engine) Render(screen *Screen) {
	engine.queue.Lock()
	engine.queue.screens[screen] = true
	engine.queue.Unlock()

	glfw.PostEmptyEvent()
}

func (engine *Engine) Running() bool {
	return engine.running
}

func (engine *Engine) Loop() error {
	// When there are nothing to draw, we just wait for new window to create.
	if len(engine.contexts) == 0 {
		delegate := <-engine.delegates
		delegate.callback()
		delegate.barrier <- struct{}{}
	} else {
		engine.queue.Lock()
		for screen, _ := range engine.queue.screens {
			delete(engine.queue.screens, screen)

			engine.render(screen)
		}
		engine.queue.Unlock()

		select {
		case delegate := <-engine.delegates:
			delegate.callback()
			delegate.barrier <- struct{}{}
		default:
			glfw.WaitEvents()
		}
	}

	if len(engine.contexts) == 0 {
		engine.free()
	}

	return nil
}

func (engine *Engine) SetFont(font *fonts.Font) {
	engine.font.handle = font
}

func (engine *Engine) GetFont() *fonts.Font {
	return engine.font.handle
}

func (engine *Engine) Stop() {
	engine.delegate(func() {
		glfw.Terminate()
		engine.running = false
	})
}

func (engine *Engine) createWindow(
	width,
	height int,
	options *messages.Open,
) *Context {
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.DoubleBuffer, glfw.False)

	var position bool

	if options.X != nil && options.Y != nil {
		position = true
	}

	// FIXME: hidden has no use now, because there is no way to show hidden
	// window from mainframe API.
	if options.Hidden || options.Raw || position {
		glfw.WindowHint(glfw.Visible, glfw.False)
	} else {
		glfw.WindowHint(glfw.Visible, glfw.True)
	}

	// FIXME: not really works with i3 for the reason it doesn't check for
	// NET_WM_STATE_ABOVE when setting floating mode:
	//
	// https://github.com/i3/i3/blob/next/src/manage.c#L439
	//if options.Floating {
	//    glfw.WindowHint(glfw.Floating, glfw.True)
	//} else {
	//    glfw.WindowHint(glfw.Floating, glfw.False)
	//}

	if options.Fixed {
		glfw.WindowHint(glfw.Resizable, glfw.False)
	} else {
		glfw.WindowHint(glfw.Resizable, glfw.True)
	}

	if options.Bare {
		glfw.WindowHint(glfw.Decorated, glfw.False)
	} else {
		glfw.WindowHint(glfw.Decorated, glfw.True)
	}

	var parent *glfw.Window
	for _, context := range engine.contexts {
		parent = context.Window
		break
	}

	var window *glfw.Window
	window, err := glfw.CreateWindow(
		width,
		height,
		options.Title,
		nil,
		parent,
	)
	if err != nil {
		return nil
	}

	if position {
		window.SetPos(*options.X, *options.Y)
	}

	if options.Raw {
		overrideRedirect(window)
	}

	context := NewContext()
	context.Window = window
	context.Screen = NewScreen(
		width,
		height,
		engine.font.handle,
		engine.Render,
	)

	window.SetCharModsCallback(
		func(
			_ *glfw.Window,
			char rune,
			mods glfw.ModifierKey,
		) {
			context.Input(char, mods)
		},
	)

	window.SetKeyCallback(
		func(
			_ *glfw.Window,
			key glfw.Key,
			scancode int,
			action glfw.Action,
			mods glfw.ModifierKey,
		) {
			context.Key(action, mods, key, scancode)
		},
	)

	window.SetRefreshCallback(
		func(
			_ *glfw.Window,
		) {
			err := engine.render(context.Screen)
			if err != nil {
				panic(err)
			}
		},
	)

	//window.SetCloseCallback(
	//    func(
	//        _ *glfw.Window,
	//    ) {
	//        fmt.Fprintln(os.Stderr, "XXXXXX engine.go:308  CLOSE")
	//    },
	//)

	window.MakeContextCurrent()

	glfw.SwapInterval(0)

	engine.clear()

	gl.Enable(gl.DEBUG_OUTPUT)
	gl.DebugMessageCallback(engine.debug, nil)
	gl.GenVertexArrays(1, &context.vao)

	engine.contexts[context.Screen] = context

	if !options.Hidden && (options.Raw || position) {
		window.Show()
	}

	return context
}

func (engine *Engine) render(screen *Screen) error {
	context := engine.contexts[screen]

	context.tick = time.Now().UnixNano() / int64(time.Microsecond)

	if context.Window.ShouldClose() {
		context.Window.Destroy()
		delete(engine.contexts, context.Screen)
		return nil
	}

	var (
		windowWidth, windowHeight = context.Window.GetSize()
		screenWidth, screenHeight = context.Screen.GetSize()
	)

	if windowWidth != screenWidth || windowHeight != screenHeight {
		context.Resize(windowWidth, windowHeight)
	}

	screen.Lock()
	defer screen.Unlock()

	context.Window.MakeContextCurrent()

	err := engine.initShaders()
	if err != nil {
		return err
	}

	err = engine.initVertices()
	if err != nil {
		return err
	}

	engine.initTextures()

	gl.BindVertexArray(context.vao)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	engine.clear()

	var (
		glyphWidth  = engine.font.handle.GetWidth()
		glyphHeight = engine.font.handle.GetHeight()
	)

	gl.Viewport(0, 0, int32(windowWidth), int32(windowHeight))

	gl.Uniform2i(0, int32(windowWidth), int32(windowHeight))
	gl.Uniform2i(1, int32(glyphWidth), int32(glyphHeight))

	gl.BindBuffer(gl.ARRAY_BUFFER, engine.vertices.buffers.triangles)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.BindBuffer(gl.ARRAY_BUFFER, engine.vertices.buffers.glyphs)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		4*len(context.Screen.GetCells()),
		gl.Ptr(context.Screen.GetCells()),
		gl.DYNAMIC_DRAW,
	)
	gl.VertexAttribIPointer(1, 2, gl.INT, 2*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribDivisor(1, 1)

	gl.BindBuffer(gl.ARRAY_BUFFER, engine.vertices.buffers.attributes)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		4*len(context.Screen.GetAttrs()),
		gl.Ptr(context.Screen.GetAttrs()),
		gl.DYNAMIC_DRAW,
	)
	gl.VertexAttribIPointer(2, 1, gl.INT, 1*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(2)
	gl.VertexAttribDivisor(2, 1)

	gl.BindBuffer(gl.ARRAY_BUFFER, engine.vertices.buffers.colors)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		4*len(context.Screen.GetColors()),
		gl.Ptr(context.Screen.GetColors()),
		gl.DYNAMIC_DRAW,
	)
	gl.VertexAttribIPointer(3, 2, gl.INT, 2*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(3)
	gl.VertexAttribDivisor(3, 1)

	gl.DrawArraysInstanced(
		gl.TRIANGLE_STRIP,
		0,
		6,
		int32(context.Screen.GetArea()),
	)

	gl.Finish()

	return nil
}

func (engine *Engine) clear() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(0, 0, 0, 1)
}

func (engine *Engine) initTextures() error {
	if engine.font.texture > 0 {
		gl.BindTexture(gl.TEXTURE_2D, engine.font.texture)

		return nil
	}

	gl.GenTextures(1, &engine.font.texture)
	gl.BindTexture(gl.TEXTURE_2D, engine.font.texture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_BORDER)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_BORDER)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(engine.font.handle.Image.Bounds().Size().X),
		int32(engine.font.handle.Image.Bounds().Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(engine.font.handle.Image.Pix),
	)

	return nil
}

func (engine *Engine) initShaders() error {
	if engine.shaders.program > 0 {
		gl.UseProgram(engine.shaders.program)

		return nil
	}

	var err error

	engine.shaders.vertex, err = engine.compileShader(
		gl.VERTEX_SHADER,
		vertexShader,
	)
	if err != nil {
		return karma.Format(
			err,
			"{gl} unable to compile vertex shader",
		)
	}

	engine.shaders.fragment, err = engine.compileShader(
		gl.FRAGMENT_SHADER,
		fragmentShader,
	)
	if err != nil {
		return karma.Format(
			err,
			"{gl} unable to compile fragment shader",
		)
	}

	engine.shaders.program = gl.CreateProgram()

	gl.AttachShader(engine.shaders.program, engine.shaders.vertex)
	gl.AttachShader(engine.shaders.program, engine.shaders.fragment)
	gl.LinkProgram(engine.shaders.program)
	gl.UseProgram(engine.shaders.program)

	return nil
}

func (engine *Engine) initVertices() error {
	if engine.vertices.points != nil {
		return nil
	}

	//         ->
	// (0; 1) X--X (1; 1)
	//      ^ |\ | |
	//      | | \| V
	// (0; 0) X--X (1; 0)
	//         <-
	engine.vertices.points = []float32{
		0, 0,
		0, 1,
		1, 0,

		1, 0,
		0, 1,
		1, 1,
	}

	// First buffer for vertices, which form triangles, which form cells.
	gl.GenBuffers(1, &engine.vertices.buffers.triangles)

	gl.BindBuffer(gl.ARRAY_BUFFER, engine.vertices.buffers.triangles)

	gl.BufferData(
		gl.ARRAY_BUFFER,
		len(engine.vertices.points)*4,
		gl.Ptr(engine.vertices.points),
		gl.STATIC_DRAW,
	)

	// Second buffer for vertices data, in this case it's glyph coordinates
	// in font.
	gl.GenBuffers(1, &engine.vertices.buffers.glyphs)

	// Third buffer for vertices data, in this case it's cell attributes.
	gl.GenBuffers(1, &engine.vertices.buffers.attributes)

	// Fourth buffer for background/foreground colors.
	gl.GenBuffers(1, &engine.vertices.buffers.colors)

	return nil
}

func (engine *Engine) free() {
	gl.DeleteBuffers(1, &engine.vertices.buffers.triangles)
	gl.DeleteBuffers(1, &engine.vertices.buffers.glyphs)
	gl.DeleteBuffers(1, &engine.vertices.buffers.attributes)
	gl.DeleteBuffers(1, &engine.vertices.buffers.colors)
	gl.DeleteTextures(1, &engine.font.texture)
	gl.DeleteShader(engine.shaders.vertex)
	gl.DeleteShader(engine.shaders.fragment)
	gl.DeleteProgram(engine.shaders.program)

	engine.vertices.points = nil
	engine.vertices.buffers.triangles = 0
	engine.vertices.buffers.glyphs = 0
	engine.vertices.buffers.attributes = 0
	engine.font.texture = 0
	engine.shaders.vertex = 0
	engine.shaders.fragment = 0
	engine.shaders.program = 0
}

func (engine *Engine) delegate(callback func()) {
	barrier := make(chan struct{}, 1)

	engine.delegates <- Delegate{
		barrier,
		callback,
	}

	<-barrier
}

func (engine *Engine) compileShader(
	kind uint32,
	source string,
) (uint32, error) {
	handle := gl.CreateShader(kind)

	buffer, free := gl.Strs(source + "\x00")
	defer free()

	gl.ShaderSource(handle, 1, buffer, nil)
	gl.CompileShader(handle)

	var result int32
	gl.GetShaderiv(handle, gl.COMPILE_STATUS, &result)
	if result == gl.FALSE {
		var length int32
		gl.GetShaderiv(handle, gl.INFO_LOG_LENGTH, &length)

		err := strings.Repeat("\x00", int(length+1))
		gl.GetShaderInfoLog(handle, length, nil, gl.Str(err))

		return 0, karma.Describe("source", source).Format(
			err,
			"{shader} compilation error",
		)
	}

	return handle, nil
}

func (engine *Engine) debug(
	source uint32,
	kind uint32,
	id uint32,
	severity uint32,
	length int32,
	message string,
	userParam unsafe.Pointer,
) {
	logger := log.Warningf

	switch kind {
	case gl.DEBUG_TYPE_ERROR:
		logger = log.Errorf
	case gl.DEBUG_TYPE_OTHER:
		logger = log.Tracef
	}

	logger(
		"{gl} %s | type=0x%x severity=0x%x source=0x%x",
		message,
		kind,
		severity,
		source,
	)
}
