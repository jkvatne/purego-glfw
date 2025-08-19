package glfw

import (
	"syscall"
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

	// This is defined in the window API's platform.h _GLFW_PLATFORM_MONITOR_STATE;
	hMonitor HMONITOR
	hDc      HDC
	// The window whose video mode is current on this monitor
	window *_GLFWwindow
	Bounds RECT
	// This size matches the static size of DISPLAY_DEVICE.DeviceName
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
	// This is from _glfwPlatformGetMonitorPos
	var dm DEVMODEW
	dm.dmSize = uint16(unsafe.Sizeof(dm))
	EnumDisplaySettingsEx(&m.adapterName[0],
		ENUM_CURRENT_SETTINGS,
		&dm,
		EDS_ROTATEDMODE)

	x = int(dm.dmPosition.X)
	y = int(dm.dmPosition.Y)
	return x, y
	// return int(m.Bounds.Left), int(m.Bounds.Top)
}

func enumMonitorCallback(monitor HMONITOR, hdc HDC, bounds RECT, lParam uintptr) bool {
	m := (*Monitor)(unsafe.Pointer(lParam))
	m.hMonitor = monitor
	m.hDc = hdc
	m.Bounds = bounds
	// monitorInfo := GetMonitorInfo(m.hMonitor)
	// Monitors = append(Monitors, &m)
	return true
}

// NewEnumDisplayMonitorsCallback is used in EnumDisplayMonitors to create the callback.
func NewEnumDisplayMonitorsCallback(callback func(monitor HMONITOR, hdc HDC, bounds RECT, lParam uintptr) bool) uintptr {
	return syscall.NewCallback(
		func(monitor HMONITOR, hdc HDC, bounds *RECT, lParam uintptr) uintptr {
			var r RECT
			if bounds != nil {
				r = *bounds
			}
			if callback(monitor, hdc, r, lParam) {
				return 1
			}
			return 0
		},
	)
}

// GetPhysicalSize returns the size, in millimetres, of the display area of the monitor.
func (m *Monitor) GetPhysicalSize() (width, height int) {
	return m.widthMM, m.heightMM
}

// GetWorkarea returns the position, in screen coordinates, of the upper-left
// corner of the work area of the specified monitor along with the work area
// size in screen coordinates.
func (m *Monitor) GetWorkarea() (x, y, width, height int) {
	mi := GetMonitorInfo(m.hMonitor)
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
	var dpiX, dpiY int
	if IsWindows8Point1OrGreater() {
		dpiX, dpiY = GetDpiForMonitor(m.hMonitor, _MDT_EFFECTIVE_DPI)
	} else {
		dc := getDC(0)
		dpiX = GetDeviceCaps(dc, _LOGPIXELSX)
		dpiX = GetDeviceCaps(dc, _LOGPIXELSY)
		releaseDC(0, dc)
	}
	return float32(dpiX) / _USER_DEFAULT_SCREEN_DPI, float32(dpiY) / _USER_DEFAULT_SCREEN_DPI
}

func (m *Monitor) GetMonitorName() string {
	s := GoStr(&m.name[0])
	return s
}

func GetVideoMode(monitor *Monitor) GLFWvidmode {
	mode := monitor.currentMode
	var dm DEVMODEW
	dm.dmSize = uint16(unsafe.Sizeof(dm))
	EnumDisplaySettingsEx(&monitor.adapterName[0], ENUM_CURRENT_SETTINGS, &dm, 0)
	mode.Width = dm.dmPelsWidth
	mode.Height = dm.dmPelsHeight
	mode.RefreshRate = dm.dmDisplayFrequency
	mode.RedBits, mode.GreenBits, mode.BlueBits = SplitBpp(dm.dmBitsPerPel)
	return mode
}

func SplitBpp(bitsPerPel int32) (int32, int32, int32) {
	bitsPerPel = min(24, bitsPerPel)
	n := bitsPerPel / 3
	blueBits := n
	redBits := n
	greenBits := n
	delta := bitsPerPel - redBits*3
	if delta >= 1 {
		greenBits++
	}
	if delta == 2 {
		redBits++
	}
	return redBits, greenBits, blueBits
}

func GetVideoModes(monitor *Monitor) (result []GLFWvidmode) {
	modeIndex := 0
	count := 0
	for {
		var mode GLFWvidmode
		var dm DEVMODEW
		dm.dmSize = uint16(unsafe.Sizeof(dm))
		n := EnumDisplaySettingsEx(&monitor.adapterName[0], modeIndex, &dm, 0)
		if n == 0 {
			break
		}
		if dm.dmSize == 0 {
			break
		}
		modeIndex++
		// Skip modes with less than 15 BPP
		if dm.dmBitsPerPel < 15 {
			continue
		}
		mode.Width = dm.dmPelsWidth
		mode.Height = dm.dmPelsHeight
		mode.RefreshRate = dm.dmDisplayFrequency
		mode.RedBits, mode.GreenBits, mode.BlueBits = SplitBpp(dm.dmBitsPerPel)
		i := 0
		for ; i < count; i++ {
			if glfwCompareVideoModes(&result[i], &mode) == 0 {
				break
			}
		}
		// Skip duplicate modes
		if i < count {
			continue
		}
		if monitor.modesPruned {
			// Skip modes not supported by the connected displays
			if ChangeDisplaySettingsEx(&monitor.adapterName[0], &dm, 0, CDS_TEST, uintptr(0)) != 0 {
				continue
			}
		}
		count++
		result = append(result, mode)
	}
	if count == 0 {
		// HACK: Report the current mode if no valid modes were found
		result = append(result, GetVideoMode(monitor))
		count = 1
	}
	return result
}

// Lexically compare video modes, used by qsort
//
func glfwCompareVideoModes(fp, sp *GLFWvidmode) int32 {
	fbpp := fp.RedBits + fp.GreenBits + fp.BlueBits
	sbpp := sp.RedBits + sp.GreenBits + sp.BlueBits
	farea := fp.Width * fp.Height
	sarea := sp.Width * sp.Height
	// First sort on color bits per pixel
	if fbpp != sbpp {
		return fbpp - sbpp
	}
	// Then sort on screen area
	if farea != sarea {
		return farea - sarea
	}
	// Then sort on width
	if fp.Width != sp.Width {
		return fp.Width - sp.Width
	}
	// Lastly sort on refresh rate
	return fp.RefreshRate - sp.RefreshRate
}

func glfwGetHMONITORContentScale(handle HMONITOR) (xscale float32, yscale float32) {
	var xdpi, ydpi int
	if IsWindows8Point1OrGreater() {
		xdpi, ydpi = GetDpiForMonitor(handle, _MDT_EFFECTIVE_DPI)
	} else {
		dc := getDC(0)
		xdpi = GetDeviceCaps(dc, _LOGPIXELSX)
		ydpi = GetDeviceCaps(dc, _LOGPIXELSY)
		releaseDC(0, dc)
	}
	return float32(xdpi) / _USER_DEFAULT_SCREEN_DPI, float32(ydpi) / _USER_DEFAULT_SCREEN_DPI
}
