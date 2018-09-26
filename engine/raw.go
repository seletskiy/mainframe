package engine

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

// #include <X11/Xlib.h>
// #include "../../../go-gl/glfw/v3.2/glfw/glfw/include/GLFW/glfw3.h"
// #cgo linux LDFLAGS: -lX11
import "C"

// Some X11 magic here.
//
// We need to set internal override redirect flag as window attribute so WM
// will not manage our window, so we it can be created without focus.
//
// Useful for making notification-like windows.
func overrideRedirect(window *glfw.Window) {
	var attrs C.XSetWindowAttributes
	attrs.override_redirect = 1

	C.XChangeWindowAttributes(
		(*C.Display)(glfw.GetX11Display()),
		(C.Window)(window.GetX11Window()),
		C.CWOverrideRedirect,
		&attrs,
	)
}
