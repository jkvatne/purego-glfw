package glfw

type SizeCallback func(w *Window, width int, height int)
type CursorPosCallback func(w *Window, xpos float64, ypos float64)
type KeyCallback func(w *Window, key Key, scancode int, action Action, mods ModifierKey)
type DropCallback func(w *Window, names []string)
type IconifyCallback func(w *Window, iconified bool)
type CharCallback func(w *Window, char rune)
type ContentScaleCallback func(w *Window, x float32, y float32)
type CursorEnterCallback func(w *Window, entered bool)
type RefreshCallback func(w *Window)
type FocusCallback func(w *Window, focused bool)
type MaximizeCallback func(w *Window, maximized bool)
type ScrollCallback func(w *Window, xoff float64, yoff float64)
type MouseButtonCallback func(w *Window, button MouseButton, action Action, mods ModifierKey)
type CloseCallback func(w *Window)
type ErrorCallbackFunc func(e int, description string)

// SetCursorPosCallback sets the cursor position callback which is called
// when the cursor is moved. The callback is provided with the position relative
// to the upper-left corner of the client area of the Window.
func (w *Window) SetCursorPosCallback(cbfun CursorPosCallback) (previous CursorPosCallback) {
	w.cursorPosCallback, previous = cbfun, w.cursorPosCallback
	return previous
}

// SetKeyCallback sets the key callback which is called when a key is pressed, repeated or released.
func (w *Window) SetKeyCallback(cbfun KeyCallback) (previous KeyCallback) {
	w.keyCallback, previous = cbfun, w.keyCallback
	return previous
}

// SetCharCallback sets the character callback which is called when a Unicode character is input.
func (w *Window) SetCharCallback(cbfun CharCallback) (previous CharCallback) {
	w.charCallback, previous = cbfun, w.charCallback
	return previous
}

// SetDropCallback sets the drop callback
func (w *Window) SetDropCallback(cbfun DropCallback) (previous DropCallback) {
	w.dropCallback, previous = cbfun, w.dropCallback
	return previous
}

// SetContentScaleCallback function sets the Window content scale callback of
// the specified Window, which is called when the content scale of the specified Window changes.
func (w *Window) SetContentScaleCallback(cbfun ContentScaleCallback) (previous ContentScaleCallback) {
	w.contentScaleCallback, previous = cbfun, w.contentScaleCallback
	return previous
}

// SetRefreshCallback sets the refresh callback of the Window, which
// is called when the client area of the Window needs to be redrawn,
func (w *Window) SetRefreshCallback(cbfun RefreshCallback) (previous RefreshCallback) {
	w.refreshCallback, previous = cbfun, w.refreshCallback
	return previous
}

// SetFocusCallback sets the focus callback of the Window, which is called when
// the Window gains or loses focus.
//
// After the focus callback is called for a Window that lost focus, synthetic key
// and mouse button release events will be generated for all such that had been
// pressed. For more information, see SetKeyCallback and SetMouseButtonCallback.
func (w *Window) SetFocusCallback(cbfun FocusCallback) (previous FocusCallback) {
	w.focusCallback, previous = cbfun, w.focusCallback
	return previous
}

// SetSizeCallback sets the size callback of the Window, which is called when
// the Window is resized. The callback is provided with the size, in screen
// coordinates, of the client area of the Window.
func (w *Window) SetSizeCallback(cbfun SizeCallback) (previous SizeCallback) {
	w.sizeCallback, previous = cbfun, w.sizeCallback
	return previous
}

// SetFramebufferSizeCallback sets the size callback of the Window, which is called when
// the Window is resized. The callback is provided with the size, in screen
// coordinates, of the client area of the Window.
func (w *Window) SetFramebufferSizeCallback(cbfun SizeCallback) (previous SizeCallback) {
	w.framebufferSizeCallback, previous = cbfun, w.framebufferSizeCallback
	return previous
}

func (w *Window) SetIconifyCallback(cbfun IconifyCallback) (previous IconifyCallback) {
	w.iconifyCallback, previous = cbfun, w.iconifyCallback
	return previous
}

// SetMouseButtonCallback sets the mouse button callback which is called when a
// mouse button is pressed or released.
//
// When a window loses focus, it will generate synthetic mouse button release
// events for all pressed mouse buttons. You can tell these events from
// user-generated events by the fact that the synthetic ones are generated after
// the window has lost focus, i.e. Focused will be false and the focus
// callback will have already been called.
func (w *Window) SetMouseButtonCallback(cbfun MouseButtonCallback) (previous MouseButtonCallback) {
	w.mouseButtonCallback, previous = cbfun, w.mouseButtonCallback
	return previous
}

// SetScrollCallback sets the scroll callback which is called when a scrolling
// device is used, such as a mouse wheel or scrolling area of a touchpad.
func (w *Window) SetScrollCallback(cbfun ScrollCallback) (previous ScrollCallback) {
	w.scrollCallback, previous = cbfun, w.scrollCallback
	return previous
}

// SetWindowCloseCallback will set set the close callback of the specified window, which is
// called when the user attempts to close the window, for example by clicking
// the close widget in the title bar.
func (w *Window) SetCloseCallback(cbfun CloseCallback) (previous CloseCallback) {
	w.windowCloseCallback, previous = cbfun, w.windowCloseCallback
	return previous
}

// SetErrorCallback is not used
func SetErrorCallback(cbfun ErrorCallbackFunc) (previous ErrorCallbackFunc) {
	_glfw.errorCallback, previous = cbfun, _glfw.errorCallback
	return previous
}
