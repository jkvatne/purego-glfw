package glfw

import (
	"fmt"
	"image"
	"image/draw"
	"syscall"
	"time"
	"unsafe"

	"golang.design/x/clipboard"
)

// Window related hints/attributes.
const (
	Focused                Hint = 0x00020001 // Specifies whether the window will be given input focus when created. This hint is ignored for full screen and initially hidden windows.
	Iconified              Hint = 0x00020002 // Specifies whether the window will be minimized.
	Maximized              Hint = 0x00020008 // Specifies whether the window is maximized.
	Visible                Hint = 0x00020004 // Specifies whether the window will be initially visible.
	Hovered                Hint = 0x0002000B // Specifies whether the cursor is currently directly over the content area of the window, with no other windows between. See Cursor enter/leave events for details.
	Resizable              Hint = 0x00020003 // Specifies whether the window will be resizable by the user.
	Decorated              Hint = 0x00020005 // Specifies whether the window will have window decorations such as a border, a close widget, et
	Floating               Hint = 0x00020007 // Specifies whether the window will be always-on-top.
	AutoIconify            Hint = 0x00020006 // Specifies whether fullscreen windows automatically iconify (and restore the previous video mode) on focus loss.
	CenterCursor           Hint = 0x00020009 // Specifies whether the cursor should be centered over newly created full screen windows. This hint is ignored for windowed mode windows.
	TransparentFramebuffer Hint = 0x0002000A // Specifies whether the framebuffer should be transparent.
	FocusOnShow            Hint = 0x0002000C // Specifies whether the window will be given input focus when glfwShowWindow is called.
	ScaleToMonitor         Hint = 0x0002200C // Specified whether the window content area should be resized based on the monitor content scale of any monitor it is placed on. This includes the initial placement when the window is created.
	MousePassthrough       Hint = 0x0002000D
	PositionX              Hint = 0x0002000E
	PositionY              Hint = 0x0002000F
	ScaleFramebuffer       Hint = 0x0002200D
	CocoaRetinaFramebuffer Hint = 0x00023001
)

// Context related hints.
const (
	ClientAPI               Hint = 0x00022001 // Specifies which client API to create the context for. Hard constraint.
	ContextVersionMajor     Hint = 0x00022002 // Specifies the client API version that the created context must be compatible with.
	ContextVersionMinor     Hint = 0x00022003 // Specifies the client API version that the created context must be compatible with.
	ContextRobustness       Hint = 0x00022005 // Specifies the robustness strategy to be used by the context.
	ContextReleaseBehavior  Hint = 0x00022009 // Specifies the release behavior to be used by the context.
	OpenGLForwardCompatible Hint = 0x00022006 // Specifies whether the OpenGL context should be forward-compatible. Hard constraint.
	OpenGLDebugContext      Hint = 0x00022007 // Specifies whether to create a debug OpenGL context, which may have additional error and performance issue reporting functionality. If OpenGL ES is requested, this hint is ignored.
	OpenGLProfile           Hint = 0x00022008 // Specifies which OpenGL profile to create the context for. Hard constraint.
	ContextCreationAPI      Hint = 0x0002200B // Specifies which context creation API to use to create the context.
	ContextNoError          Hint = 0x0002200A
)

// Framebuffer related hints.
const (
	ContextRevision        Hint = 0x00022004
	RedBits                Hint = 0x00021001 // Specifies the desired bit depth of the default framebuffer.
	GreenBits              Hint = 0x00021002 // Specifies the desired bit depth of the default framebuffer.
	BlueBits               Hint = 0x00021003 // Specifies the desired bit depth of the default framebuffer.
	AlphaBits              Hint = 0x00021004 // Specifies the desired bit depth of the default framebuffer.
	DepthBits              Hint = 0x00021005 // Specifies the desired bit depth of the default framebuffer.
	StencilBits            Hint = 0x00021006 // Specifies the desired bit depth of the default framebuffer.
	AccumRedBits           Hint = 0x00021007 // Specifies the desired bit depth of the accumulation buffer.
	AccumGreenBits         Hint = 0x00021008 // Specifies the desired bit depth of the accumulation buffer.
	AccumBlueBits          Hint = 0x00021009 // Specifies the desired bit depth of the accumulation buffer.
	AccumAlphaBits         Hint = 0x0002100A // Specifies the desired bit depth of the accumulation buffer.
	AuxBuffers             Hint = 0x0002100B // Specifies the desired number of auxiliary buffers.
	Stereo                 Hint = 0x0002100C // Specifies whether to use stereoscopic rendering. Hard constraint.
	Samples                Hint = 0x0002100D // Specifies the desired number of samples to use for multisampling. Zero disables multisampling.
	SRGBCapable            Hint = 0x0002100E // Specifies whether the framebuffer should be sRGB capable.
	RefreshRate            Hint = 0x0002100F // Specifies the desired refresh rate for full screen windows. If set to zero, the highest available refresh rate will be used. This hint is ignored for windowed mode windows.
	DoubleBuffer           Hint = 0x00021010 // Specifies whether the framebuffer should be double buffered. You nearly always want to use double buffering. This is a hard constraint.
	CocoaGraphicsSwitching Hint = 0x00023003 // Specifies whether to in Automatic Graphics Switching, i.e. to allow the system to choose the integrated GPU for the OpenGL context and move it between GPUs if necessary or whether to force it to always run on the discrete GPU.
)

// Naming related hints. (Use with glfw.WindowHintString)
const (
	CocoaFrameName  Hint = 0x00023002 // Specifies the UTF-8 encoded name to use for autosaving the window frame, or if empty disables frame autosaving for the window.
	X11ClassName    Hint = 0x00024001 // Specifies the desired ASCII encoded class parts of the ICCCM WM_CLASS window property.nd instance parts of the ICCCM WM_CLASS window property.
	X11InstanceName Hint = 0x00024002 // Specifies the desired ASCII encoded instance parts of the ICCCM WM_CLASS window property.nd instance parts of the ICCCM WM_CLASS window property.
	WaylandAppId    Hint = 0x00026001
)

const (
	OpenGlAnyProfile    = 0
	OpenGLCoreProfile   = 0x00032001
	OpenGLCompatProfile = 0x00032002
	AnyPosition         = int32(-0x80000000)
)

// Values for the ClientAPI hint.
const (
	OpenGLAPI   int32 = 0x00030001
	OpenGLESAPI int32 = 0x00030002
	NoAPI       int32 = 0
)

// Values for ContextCreationAPI hint.
const (
	NativeContextAPI int32 = 0x00036001
	EGLContextAPI    int32 = 0x00036002
	OSMesaContextAPI int32 = 0x00036003
)

// Values for the ContextRobustness hint.
const (
	NoRobustness        int32 = 0
	NoResetNotification int32 = 0x00031001
	LoseContextOnReset  int32 = 0x00031002
)

// Values for ContextReleaseBehavior hint.
const (
	AnyReleaseBehavior   int32 = 0
	ReleaseBehaviorFlush int32 = 0x00035001
	ReleaseBehaviorNone  int32 = 0x00035002
)

// Other values.
const (
	True     = 1 // GL_TRUE
	False    = 0 // GL_FALSE
	DontCare = -1
)

// InputMode corresponds to an input mode.
type InputMode int

const (
	CursorMode            = 0x00033001
	StickyKeys            = 0x00033002
	StickyMouseButtons    = 0x00033003
	LockKeyMods           = 0x00033004
	RawMouseMotion        = 0x00033005
	UnlimitedMouseButtons = 0x00033006
)

// Cursor modes
const (
	CursorNormal   = 0x00034001
	CursorHidden   = 0x00034002
	CursorDisabled = 0x00034003
	CursorCaptured = 0x00034004
)

// MouseButton definitions
type MouseButton int

const (
	MouseButtonFirst  MouseButton = 0
	MouseButtonLeft   MouseButton = 0
	MouseButtonRight  MouseButton = 1
	MouseButtonMiddle MouseButton = 2
	MouseButtonLast   MouseButton = 2
)

// Exported cursor types
const (
	ArrowCursor      = 0x00036001
	IBeamCursor      = 0x00036002
	CrosshairCursor  = 0x00036003
	HandCursor       = 0x00036004
	HResizeCursor    = 0x00036005
	VResizeCursor    = 0x00036006
	ResizeNwseCursor = 0x00036007
	ResizeNeswCursor = 0x00036008
	ResizeAllCursor  = 0x00036009
	NotAllowedCursor = 0x0003600A
)

type Action int

type StandardCursor uint16

type Hint uint32

// Window represents a Window.
type Window = _GLFWwindow

type Cursor struct {
	next   *Cursor
	handle syscall.Handle
}

// PollEvents processes only those events that have already been received and
// then returns immediately. Processing events will cause the Window and input
// callbacks associated with those events to be called.
func PollEvents() {
	glfwPollEvents()
}

func WaitEvents() {
	WaitMessage()
	glfwPollEvents()
}

// WaitEventsTimeout waits a number of seconds or until an event is detected
func WaitEventsTimeout(timeout float64) {
	if timeout < 0.0 {
		panic("Wait time must be positive")
	}
	MsgWaitForMultipleObjects(0, nil, 0, uint32(timeout*1e3), _QS_ALLINPUT)
	glfwPollEvents()
}

func WindowHint(hint Hint, v int) error {
	value := int32(v)
	switch hint {
	case RedBits:
		_glfw.hints.framebuffer.redBits = value
	case GreenBits:
		_glfw.hints.framebuffer.greenBits = value
	case BlueBits:
		_glfw.hints.framebuffer.blueBits = value
	case AlphaBits:
		_glfw.hints.framebuffer.alphaBits = value
	case DepthBits:
		_glfw.hints.framebuffer.depthBits = value
	case StencilBits:
	case AccumRedBits:
		_glfw.hints.framebuffer.accumRedBits = value
	case AccumGreenBits:
		_glfw.hints.framebuffer.accumGreenBits = value
	case AccumBlueBits:
		_glfw.hints.framebuffer.accumBlueBits = value
	case AccumAlphaBits:
		_glfw.hints.framebuffer.accumAlphaBits = value
	case AuxBuffers:
		_glfw.hints.framebuffer.auxBuffers = value
	case DoubleBuffer:
		_glfw.hints.framebuffer.doublebuffer = value != 0
	case TransparentFramebuffer:
		_glfw.hints.framebuffer.transparent = value != 0
	case Samples:
		_glfw.hints.framebuffer.samples = value
	case SRGBCapable:
		_glfw.hints.framebuffer.sRGB = value != 0
	case Resizable:
		_glfw.hints.window.resizable = value != 0
	case Decorated:
		_glfw.hints.window.decorated = value != 0
	case Focused:
		_glfw.hints.window.focused = value != 0
	case AutoIconify:
		_glfw.hints.window.autoIconify = value != 0
	case Floating:
		_glfw.hints.window.floating = value != 0
	case Maximized:
		_glfw.hints.window.maximized = value != 0
	case Visible:
		_glfw.hints.window.visible = value != 0
	case PositionX:
		_glfw.hints.window.xpos = int32(value)
	case PositionY:
		_glfw.hints.window.ypos = int32(value)
	case ScaleToMonitor:
		_glfw.hints.window.scaleToMonitor = value != 0
	case ScaleFramebuffer, CocoaRetinaFramebuffer:
		_glfw.hints.window.scaleFramebuffer = value != 0
	case CenterCursor:
		_glfw.hints.window.centerCursor = value != 0
	case FocusOnShow:
		_glfw.hints.window.focusOnShow = value != 0
	case MousePassthrough:
		_glfw.hints.window.mousePassthrough = value != 0
	case ClientAPI:
		_glfw.hints.context.client = int32(value)
	case ContextCreationAPI:
		_glfw.hints.context.source = int32(value)
	case ContextVersionMajor:
		_glfw.hints.context.major = int32(value)
	case ContextVersionMinor:
		_glfw.hints.context.minor = int32(value)
	case ContextRobustness:
		_glfw.hints.context.robustness = int32(value)
	case OpenGLForwardCompatible:
		_glfw.hints.context.forward = value != 0
	case OpenGLDebugContext:
		_glfw.hints.context.debug = value != 0
	case ContextNoError:
		_glfw.hints.context.noerror = value != 0
	case OpenGLProfile:
		_glfw.hints.context.profile = value
	case ContextReleaseBehavior:
		_glfw.hints.context.release = value
	case RefreshRate:
		_glfw.hints.refreshRate = value
	default:
		return fmt.Errorf("Invalid window hint %d with value %d", hint, value)
	}
	return nil
}

func WindowHintString(hint Hint, value string) {
	switch hint {
	case CocoaFrameName:
		_glfw.hints.window.ns.frameName = value
	case X11ClassName:
		// _glfw.hints.window.x11.classNam = value
	case X11InstanceName:
		// _glfw.hints.window.x11.instanceName = value
	case WaylandAppId:
		// _glfw.hints.window.wl.appId =  value
	default:
		panic(fmt.Sprintf("Invalid Window hint %d with value %s", hint, value))
	}
}

// GetClipboardString returns the contents of the system clipboard
// if it contains or is convertible to a UTF-8 encoded string.
// This function may only be called from the main thread.
func GetClipboardString() string {
	return glfwGetClipboardString()
}

// SetClipboardString sets the system clipboard to the specified UTF-8 encoded string.
// This function may only be called from the main thread.
func SetClipboardString(str string) {
	glfwSetClipboardString(str)
}

func CreateCursor(image image.Image, xhot int, yhot int) *Cursor {
	if image == nil || image.Bounds().Dx() <= 0 || image.Bounds().Dy() <= 0 {
		panic("CreateCursor: image is nil or invalid")
	}
	var cursor Cursor
	cursor.next = _glfw.cursorListHead
	_glfw.cursorListHead = &cursor
	im := imageToGLFW(image)
	cursor.handle = createIcon(&im, int32(xhot), int32(yhot), false)
	return &cursor
}

// CreateStandardCursor returns a cursor with a standard shape,
// that can be set for a Window with SetCursor.
func CreateStandardCursor(shape int) *Cursor {
	var cursor = Cursor{}
	cursor.next = _glfw.cursorListHead
	_glfw.cursorListHead = &cursor
	var id uint16
	switch shape {
	case ArrowCursor:
		id = IDC_ARROW
	case IBeamCursor:
		id = IDC_IBEAM
	case CrosshairCursor:
		id = IDC_CROSS
	case HResizeCursor:
		id = IDC_SIZEWE
	case VResizeCursor:
		id = IDC_SIZENS
	case HandCursor:
		id = IDC_HAND
	case ResizeAllCursor:
		id = IDC_SIZEALL
	case ResizeNeswCursor:
		id = IDC_SIZENESW
	case ResizeNwseCursor:
		id = IDC_SIZENWSE
	case NotAllowedCursor:
		id = IDC_NO
	default:
		panic("Win32: Unknown or unsupported standard cursor")
	}
	cursor.handle = LoadCursor(id)
	if cursor.handle == 0 {
		panic("Failed to create standard cursor")
	}
	return &cursor
}

func CreateWindow(width, height int, title string, monitor *Monitor, share *Window) (*Window, error) {
	wnd, err := glfwCreateWindow(int32(width), int32(height), title, monitor, share)
	if err != nil {
		return nil, fmt.Errorf("CreateWindow failed: %v", err)
	}
	return wnd, nil
}

// Destroy destroys the specified window and its context.
func (w *Window) Destroy() {
	glfwDestroyWindow(w)
}

// ShouldClose reports the close flag value for the specified Window.
func (w *Window) ShouldClose() bool {
	return w.shouldClose
}

// SetShouldClose sets the value of the close flag of the window. This can be
// used to override the user's attempt to close the window, or to signal that it
// should be closed.
func (w *Window) SetShouldClose(value bool) {
	w.shouldClose = value
}

// SetTitle sets the window title, encoded as UTF-8, of the window.
//
// This function may only be called from the main thread.
func (w *Window) SetTitle(title string) {
	glfwSetTitle(w, title)

}

// SetIcon sets the icon of the specified window. If passed an array of candidate images,
// those of or closest to the sizes desired by the system are selected. If no images are
// specified, the window reverts to its default icon.
// The image is ideally provided in the form of *image.NRGBA.
// The pixels are 32-bit, little-endian, non-premultiplied RGBA, i.e. eight
// bits per channel with the red channel first. They are arranged canonically
// as packed sequential rows, starting from the top-left corner. If the image
// type is not *image.NRGBA, it will be converted to it.
// The desired image sizes varies depending on platform and system settings. The selected
// images will be rescaled as needed. Good sizes include 16x16, 32x32 and 48x48.
func (w *Window) SetIcon(images []image.Image) {
	count := len(images)
	cImages := make([]*GLFWimage, count)
	for i, img := range images {
		im := imageToGLFW(img)
		cImages[i] = &im
	}
	glfwSetWindowIcon(w, count, cImages)
}

// GetPos returns the position, in screen coordinates, of the upper-left
// corner of the client area of the window.
func (w *Window) GetPos() (x, y int) {
	xx, yy := glfwGetPos(w)
	return int(xx), int(yy)
}

// SetPos sets the position, in screen coordinates, of the Window's upper-left corner
func (w *Window) SetPos(xPos, yPos int) {
	glfwSetPos(w, int32(xPos), int32(yPos))
}

// GetSize returns the size, in screen coordinates, of the client area of the
// specified Window.
func (w *Window) GetSize() (width int, height int) {
	var wi, h int32
	glfwGetWindowSize(w, &wi, &h)
	return int(wi), int(h)
}

// SetSize sets the size, in screen coordinates, of the client area of the Window.
func (w *Window) SetSize(width, height int) {
	if w.monitor != nil {
		if w.monitor.window == w {
			acquireMonitor(w)
			fitToMonitor(w)
		}
	} else {
		rect := RECT{0, 0, int32(width), int32(height)}
		adjustWindowRect(&rect, getWindowStyle(w), 0, getWindowExStyle(w), GetDpiForWindow(w.Win32.handle), "glfwSetWindowSize")
		SetWindowPos(w.Win32.handle, 0, 0, 0, int32(width), int32(height), SWP_NOACTIVATE|SWP_NOOWNERZORDER|SWP_NOMOVE|SWP_NOZORDER)
	}
}

func (w *Window) SetSizeLimits(minw, minh, maxw, maxh int) {
	if (minw == DontCare || minh == DontCare) && (maxw == DontCare || maxh == DontCare) {
		return
	}
	area := GetWindowRect(w.Win32.handle)
	MoveWindow(w.Win32.handle, area.Left, area.Top, area.Right-area.Left, area.Bottom-area.Top, true)
}

// SetAspectRatio sets the required aspect ratio of the client area of the specified window.
func (w *Window) SetAspectRatio(numer, denom int) {
	glfwSetWindowAspectRatio(w, int32(numer), int32(denom))
}

func (w *Window) GetFramebufferSize() (int, int) {
	return glfwGetFramebufferSize(w)
}

// GetFrameSize retrieves the size, in screen coordinates, of each edge of the frame
// This size includes the title bar if the Window has one.
func (w *Window) GetFrameSize() (left, top, right, bottom int) {
	var l, t, r, b int32
	glfwGetWindowFrameSize(w, &l, &t, &r, &b)
	return int(l), int(t), int(r), int(b)
}

// GetContentScale function retrieves the content scale for the specified
// Window. The content scale is the ratio between the current DPI and the
// platform's default DPI.
func (w *Window) GetContentScale() (float32, float32) {
	return glfwGetContentScale(w)
}

// GetOpacity function returns the opacity of the window
func (w *Window) GetOpacity() float32 {
	return glfwGetWindowOpacity(w)
}

// SetOpacity function sets the opacity of the window (0 to 1.0)
func (w *Window) SetOpacity(opacity float64) {
	if opacity < 0.0 || (opacity > 1.0) {
		panic("SetOpacity: opacity must be between 0.0 and 1.0")
	}
	glfwSetWindowOpacity(w, opacity)
}

// RequestWindowAttention funciton requests user attention to the specified window.
func (w *Window) RequestAttention() {
	glfwRequestWindowAttention(w)
}

// Focus brings the specified Window to front and sets input focus.
func (w *Window) Focus() {
	glfwFocusWindow(w)
}

// Iconify iconifies/minimizes the window, if it was previously restored.
func (w *Window) Iconify() {
	w.Win32.maximized = false
	w.Win32.iconified = true
	glfwShowWindow(w)
}

// Maximize maximizes the specified window if it was previously not maximized.
func (w *Window) Maximize() {
	w.Win32.iconified = false
	w.Win32.maximized = true
	glfwShowWindow(w)
}

// Restore restores the window, if it was previously iconified/minimized.
func (w *Window) Restore() {
	glfwRestoreWindow(w)
}

// Show makes the Window visible if it was previously hidden.
func (w *Window) Show() {
	if w.monitor != nil {
		return
	}
	glfwShowWindow(w)
	if w.focusOnShow {
		glfwFocusWindow(w)
	}
}

// Show makes the Window visible if it was previously hidden.
func (w *Window) Hide() {
	glfwHideWindow(w)
}

// GetMonitor returns the handle of the monitor that the window is in fullscreen on.
// Returns nil if the window is in windowed mode.
func (w *Window) GetMonitor() *Monitor {
	return glfwGetWindowMonitor(w)
}

// SetMonitor sets the monitor that the window uses for full screen mode or,
// if the monitor is NULL, makes it windowed mode.
func (w *Window) SetMonitor(monitor *Monitor, xpos, ypos, width, height, refreshRate int) {
	glfwSetWindowMonitor(w, monitor, int32(xpos), int32(ypos), int32(width), int32(height), int32(refreshRate))
}

// GetAttrib returns an attribute of the window.
func (w *Window) GetAttrib(attrib Hint) int {
	return int(glfwGetWindowAttrib(w, attrib))
}

// SetAttrib function sets the value of an attribute of the specified window.
func (w *Window) SetAttrib(attrib Hint, value int) {
	glfwSetWindowAttrib(w, attrib, int32(value))
}

// SwapBuffers swaps the front and back buffers of the Window.
func (w *Window) SwapBuffers() {
	glfwSwapBuffers(w)
}

// imageToGLFW converts img to be compatible with C.GLFWimage.
func imageToGLFW(img image.Image) (r GLFWimage) {
	b := img.Bounds()
	r.Width = int32(b.Dx())
	r.Height = int32(b.Dy())
	var pixels *[16384]uint8
	if m, ok := img.(*image.NRGBA); ok && m.Stride == b.Dx()*4 {
		p := m.Pix[:m.PixOffset(m.Rect.Min.X, m.Rect.Max.Y)]
		pixels = (*[16384]uint8)(unsafe.Pointer(&p[0]))
	} else {
		m := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(m, m.Bounds(), img, b.Min, draw.Src)
		pixels = (*[16384]uint8)(unsafe.Pointer(&m.Pix[0]))
	}
	r.Pixels = &pixels[0]
	return r
}

type GLFWimage struct {
	Width  int32  // The width, in pixels, of this image.
	Height int32  // The height, in pixels, of this image.
	Pixels *uint8 // The pixel data of this image, arranged left-to-right, top-to-bottom.
}

// SetCursor sets the cursor image to be used when the cursor is over the client area
func (w *Window) SetCursor(c *Cursor) {
	glfwSetCursor(w, c)
}

// GetCursorPos returns the last reported position of the cursor.
func (w *Window) GetCursorPos() (x float64, y float64) {
	var xPos, yPos int32
	glfwGetCursorPos(w, &xPos, &yPos)
	return float64(xPos), float64(yPos)
}

func (w *Window) MakeContextCurrent() {
	_ = w.context.makeCurrent(w)
}

// DetachCurrentContext detaches the current context.
func DetachCurrentContext() {
	glfwDetachCurrentContext()
}

// GetCurrentContext returns the window whose context is current.
func GetCurrentContext() *Window {
	return glfwGetCurrentContext()
}

// Terminate destroys all remaining Windows, frees any allocated resources and
// sets the library to an uninitialized state.
func Terminate() {
	glfwTerminate()
}

func getTime() float64 {
	return float64(time.Now().UnixNano()) / 1.0e9
}

var startTime = float64(time.Now().UnixNano()) / 1.0e9

func GetTime() float64 {
	return getTime() - startTime
}

func SetTime(t float64) {
	startTime = getTime() - t
}

// Init is glfwInit(void)
func Init() error {
	// Repeated calls do nothing
	if _glfw.initialized {
		return nil
	}
	_glfw.initialized = true
	_glfw.hints.init = _GLFWinitconfig{}
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
	if err = glfwPlatformInit(); err != nil {
		return err
	}
	err = glfwPlatformCreateTls(&_glfw.errorSlot)
	if err != nil {
		return err
	}
	err = glfwPlatformCreateTls(&_glfw.contextSlot)
	if err != nil {
		return err
	}
	return nil
}

func (window *Window) GetInputMode(mode InputMode) int {
	switch mode {
	case CursorMode:
		return window.cursorMode
	case StickyKeys:
		return toInt(window.stickyKeys)
	case StickyMouseButtons:
		return toInt(window.stickyMouseButtons)
	case LockKeyMods:
		return toInt(window.lockKeyMods)
	case RawMouseMotion:
		return int(window.rawMouseMotion)
	default:
		panic("Unknown InputMode")
	}
	return 0
}

func (window *Window) Focused() bool {
	return window.Win32.handle == GetActiveWindow()
}

func (window *Window) SetCursorMode(mode int) {
	if window.Focused() {
		if mode == CursorDisabled {
			_glfw.win32.restoreCursorPosX, _glfw.win32.restoreCursorPosY = window.GetCursorPos()
			// Center Cursor In Content Area
			x, y := window.GetCursorPos()
			window.SetCursorPos(x/2, y/2)
			if window.rawMouseMotion != 0 {
				enableRawMouseMotion(window)
			}
		} else if _glfw.win32.disabledCursorWindow == window {
			if window.rawMouseMotion != 0 {
				disableRawMouseMotion(window)
			}
		}
		if mode == CursorDisabled || mode == CursorCaptured {
			captureCursor(window)
		} else {
			releaseCursor()
		}
		if mode == CursorDisabled {
			_glfw.win32.disabledCursorWindow = window
		} else if _glfw.win32.disabledCursorWindow == window {
			_glfw.win32.disabledCursorWindow = nil
			window.SetCursorPos(_glfw.win32.restoreCursorPosX, _glfw.win32.restoreCursorPosY)
		}
	}
	if cursorInContentArea(window) {
		updateCursorImage(window)
	}
}

func (window *Window) SetInputMode(mode int, value int) {
	switch mode {
	case CursorMode:
		if value != CursorNormal &&
			value != CursorHidden &&
			value != CursorDisabled &&
			value != CursorCaptured {
			fmt.Printf("Invalid cursor mode 0x%08X", value)
		}
		if window.cursorMode == value {
			return
		}
		window.cursorMode = value
		window.virtualCursorPosX, window.virtualCursorPosY = window.GetCursorPos()
		window.SetCursorMode(value)
	case StickyKeys:
		value = min(1, max(0, value))
		if window.stickyKeys == (value != 0) {
			return
		}
		if value == 0 {
			// Release all sticky keys
			for i := 0; i <= KeyLast; i++ {
				if window.keys[i] == Stick {
					window.keys[i] = Release
				}
			}
			window.stickyKeys = value != 0
		}
	case StickyMouseButtons:
		value = min(1, max(0, value))
		if window.stickyMouseButtons == (value != 0) {
			return
		}
		if value == 0 {
			// Release all sticky mouse buttons
			for i := MouseButton(0); i <= MouseButtonLast; i++ {
				if window.mouseButtons[i] == Stick {
					window.mouseButtons[i] = Release
				}
			}
			window.stickyMouseButtons = value != 0
		}
	case LockKeyMods:
		value = min(1, max(0, value))
		window.lockKeyMods = value != 0
	case RawMouseMotion:
		value = min(1, max(0, value))
		if window.rawMouseMotion == value {
			return
		}
		window.rawMouseMotion = value
		setRawMouseMotion(window, value != 0)
	case UnlimitedMouseButtons:
		value = min(1, max(0, value))
		window.disableMouseButtonLimit = value != 0
	default:
		panic(fmt.Sprintf("Invalid input mode 0x%08X", mode))
	}
}

func (window *Window) SetCursorPos(x, y float64) {
	pos := POINT{int32(x), int32(y)}
	// Store the new position so it can be recognized later
	window.lastCursorPosX = float64(pos.X)
	window.lastCursorPosY = float64(pos.Y)
	pos = ClientToScreen(window.Win32.handle, pos)
	SetCursorPos(pos.X, pos.Y)
}

func setRawMouseMotion(window *Window, enabled bool) {
	if _glfw.win32.disabledCursorWindow != window {
		return
	}
	if enabled {
		enableRawMouseMotion(window)
	} else {
		disableRawMouseMotion(window)
	}
}

func RawMouseMotionSupported() bool {
	return true
}

func DestroyCursor(cursor *Cursor) {
	if cursor == nil {
		return
	}
	// Make sure the cursor is not being used by any window
	for window := _glfw.windowListHead; window != nil; window = window.next {
		if window.cursor == cursor {
			window.SetCursor(nil)
		}
		destroyCursor(cursor)
		// Unlink cursor from global linked list
		prev := &_glfw.cursorListHead
		for *prev != cursor {
			prev = &((*prev).next)
			*prev = cursor.next
		}
	}
}

func DefaultWindowHints() {
	_glfw.hints.context.client = OpenGLAPI
	_glfw.hints.context.source = NativeContextAPI
	_glfw.hints.context.major = 1
	_glfw.hints.context.minor = 0
	// The default is a focused, visible, resizable window with decorations
	_glfw.hints.window.resizable = true
	_glfw.hints.window.visible = true
	_glfw.hints.window.decorated = true
	_glfw.hints.window.focused = true
	_glfw.hints.window.autoIconify = false
	_glfw.hints.window.centerCursor = true
	_glfw.hints.window.focusOnShow = true
	_glfw.hints.window.xpos = AnyPosition
	_glfw.hints.window.ypos = AnyPosition
	_glfw.hints.window.scaleFramebuffer = true
	// The default is 24 bits of color, 24 bits of depth and 8 bits of stencil, double buffered
	_glfw.hints.framebuffer.redBits = 8
	_glfw.hints.framebuffer.greenBits = 8
	_glfw.hints.framebuffer.blueBits = 8
	_glfw.hints.framebuffer.alphaBits = 8
	_glfw.hints.framebuffer.depthBits = 24
	_glfw.hints.framebuffer.stencilBits = 8
	_glfw.hints.framebuffer.doublebuffer = true
	// The default is to select the highest available refresh rate
	_glfw.hints.refreshRate = DontCare
	// The default is to use full Retina resolution framebuffers
	_glfw.hints.window.ns.retina = true
}

func SwapInterval(interval int) {
	window := glfwGetCurrentContext()
	if window == nil {
		panic("glfwSwapInterval: window == nil")
	}
	window.context.swapInterval(interval)
}

func (w *Window) PostEmptyEvent() {
	glfwPostEmptyEvent(w)
}
