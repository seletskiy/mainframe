package engine

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Context struct {
	vao uint32

	Window *glfw.Window
	Screen *Screen
}

func (context *Context) Close() {
	context.Window.SetShouldClose(true)
}
