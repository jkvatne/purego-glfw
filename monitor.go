package glfw

import (
	"unsafe"
)

const EDS_ROTATEDMODE = 0x00000004
const CDS_TEST = 0x00000002

// Monitor structure
//
type Monitor struct {
	name        [128]byte
	userPointer unsafe.Pointer
	widthMM     int
	heightMM    int
	modes       []GLFWvidmode
	currentMode GLFWvidmode
	// The window whose video mode is current on this monitor
	window *_GLFWwindow
	Win32  _GLFWMonitorWin32
}

type _GLFWMonitorWin32 struct {
	hMonitor          HMONITOR
	hDc               HDC
	Bounds            RECT
	adapterName       [32]uint16
	displayName       [32]uint16
	publicAdapterName string
	publicDisplayName string
	modesPruned       bool
	modeChanged       bool
}

type POINT struct {
	X, Y int32
}

// GetMonitors returns a slice of handles for all currently connected monitors.
func GetMonitors() []*Monitor {
	return _glfw.monitors
}

// GetPrimaryMonitor returns the primary monitor. This is usually the monitor
// where elements like the Windows task bar or the OS X menu bar is located.
func GetPrimaryMonitor() *Monitor {
	if len(_glfw.monitors) == 0 {
		return nil
	}
	return _glfw.monitors[0]
}

// GetPos returns the position, in screen coordinates, of the upper-left
// corner of the monitor.
func (m *Monitor) GetPos() (x, y int) {
	return glfwGetMonitorPos(m)
}

// GetPhysicalSize returns the size, in millimetres, of the display area of the monitor.
func (m *Monitor) GetPhysicalSize() (width, height int) {
	return m.widthMM, m.heightMM
}

// GetWorkarea returns the position, in screen coordinates, of the upper-left
// corner of the work area of the specified monitor along with the work area
// size in screen coordinates.
func (m *Monitor) GetWorkarea() (x, y, width, height int) {
	mi := GetMonitorInfo(m.Win32.hMonitor)
	x = int(mi.RcWork.Left)
	y = int(mi.RcWork.Top)
	width = int(mi.RcWork.Right - mi.RcWork.Left)
	height = int(mi.RcWork.Bottom - mi.RcWork.Top)
	return x, y, width, height
}

// GetContentScale function retrieves the content scale for the specified monitor.
// The content scale is the ratio between the current DPI and the platform's
// default DPI. If you scale all pixel dimensions by this scale then your content
// should appear at an appropriate size. This is especially important for text
// and any UI elements.
//
// This function must only be called from the main thread.
func (m *Monitor) GetContentScale() (float32, float32) {
	return glfwGetMonitorContentScale(m)
}

func (m *Monitor) GetMonitorName() string {
	s := GoStr(&m.name[0])
	return s
}

func (monitor *Monitor) GetVideoMode() GLFWvidmode {
	return glfwGetVideoMode(monitor)
}

func (m *Monitor) GetVideoModes() []GLFWvidmode {
	if !refreshVideoModes(m) {
		return nil
	}
	return m.modes
}
