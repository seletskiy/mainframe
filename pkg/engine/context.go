package engine

import (
	"sync"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/seletskiy/mainframe/pkg/protocol/messages"
)

const (
	SubscriptionResize   = 1
	SubscriptionKeyboard = 2
	SubscriptionInput    = 3
)

type Context struct {
	Window *glfw.Window
	Screen *Screen

	vao  uint32
	tick int64

	subscriptions struct {
		sync.Mutex

		clients map[int][]*Client
	}
}

func NewContext() *Context {
	context := &Context{}
	context.subscriptions.clients = map[int][]*Client{}
	return context
}

func (context *Context) Subscribe(client *Client, subscription int) {
	context.subscriptions.Lock()
	context.subscriptions.clients[subscription] = append(
		context.subscriptions.clients[subscription],
		client,
	)
	context.subscriptions.Unlock()
}

func (context *Context) Resize(width, height int) {
	rows, columns := context.Screen.Resize(width, height)

	// TODO move out of render loop
	subscribers := context.subscriptions.clients[SubscriptionResize]
	for _, client := range subscribers {
		client.Send(&messages.EventResize{
			Event: messages.Event{
				Tick: context.tick,
				Kind: "resize",
			},

			Width:   width,
			Height:  height,
			Columns: columns,
			Rows:    rows,
		})
	}
}

func (context *Context) Input(
	char rune,
	mods glfw.ModifierKey,
) {
	subscribers := context.subscriptions.clients[SubscriptionInput]
	for _, client := range subscribers {
		client.Send(&messages.EventInput{
			Event: messages.Event{
				Tick: context.tick,
				Kind: "input",
			},

			Char: char,

			Shift: 0 != mods&glfw.ModShift,
			Ctrl:  0 != mods&glfw.ModControl,
			Alt:   0 != mods&glfw.ModAlt,
			Super: 0 != mods&glfw.ModSuper,
		})
	}
}

func (context *Context) Key(
	action glfw.Action,
	mods glfw.ModifierKey,
	key glfw.Key,
	scancode int,
) {
	subscribers := context.subscriptions.clients[SubscriptionKeyboard]

	for _, client := range subscribers {
		client.Send(&messages.EventKeyboard{
			Event: messages.Event{
				Tick: context.tick,
				Kind: "keyboard",
			},

			Press:   action == glfw.Press,
			Release: action == glfw.Release,
			Repeat:  action == glfw.Repeat,

			Symbol: MapKeyToSymbol(key, mods),
			Code:   scancode,

			Shift: 0 != mods&glfw.ModShift,
			Ctrl:  0 != mods&glfw.ModControl,
			Alt:   0 != mods&glfw.ModAlt,
			Super: 0 != mods&glfw.ModSuper,
		})
	}
}

func (context *Context) Close() {
	if context.Window != nil {
		context.Window.SetShouldClose(true)
		context.Screen.Render()
	}
}
