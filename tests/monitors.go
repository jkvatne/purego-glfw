// Monitor information tool
// This test prints monitor and video mode information or verifies video modes
package main

import (
	"fmt"
	"os"
	"runtime"

	glfw "github.com/jkvatne/purego-glfw"
	"github.com/neclepsio/gl/all-core/gl"
)

type colDef struct {
	name string
	r    float32
	g    float32
	b    float32
}

var colors = []colDef{
	{name: "Red", r: 1, g: 0, b: 0},
	{name: "Green", r: 0, g: 1, b: 0},
	{name: "Blue", r: 0, g: 0, b: 1},
	{name: "White", r: 1, g: 1, b: 1},
	{name: "Yellow", r: 1, g: 1, b: 0},
	{name: "Cyan", r: 0, g: 1, b: 1},
	{name: "Magenta", r: 1, g: 0, b: 1},
	{name: "Red", r: 1, g: 0, b: 0},
	{name: "Green", r: 0, g: 1, b: 0},
	{name: "Blue", r: 0, g: 0, b: 1},
	{name: "White", r: 1, g: 1, b: 1},
	{name: "Yellow", r: 1, g: 1, b: 0},
	{name: "Cyan", r: 0, g: 1, b: 1},
	{name: "Magenta", r: 1, g: 0, b: 1},
	{name: "Red", r: 1, g: 0, b: 0},
	{name: "Green", r: 0, g: 1, b: 0},
	{name: "Blue", r: 0, g: 0, b: 1},
	{name: "White", r: 1, g: 1, b: 1},
	{name: "Yellow", r: 1, g: 1, b: 0},
	{name: "Cyan", r: 0, g: 1, b: 1},
	{name: "Magenta", r: 1, g: 0, b: 1},
	{name: "Red", r: 1, g: 0, b: 0},
	{name: "Green", r: 0, g: 1, b: 0},
	{name: "Blue", r: 0, g: 0, b: 1},
	{name: "White", r: 1, g: 1, b: 1},
	{name: "Yellow", r: 1, g: 1, b: 0},
	{name: "Cyan", r: 0, g: 1, b: 1},
	{name: "Magenta", r: 1, g: 0, b: 1},
	{name: "Red", r: 1, g: 0, b: 0},
	{name: "Green", r: 0, g: 1, b: 0},
	{name: "Blue", r: 0, g: 0, b: 1},
	{name: "White", r: 1, g: 1, b: 1},
	{name: "Yellow", r: 1, g: 1, b: 0},
	{name: "Cyan", r: 0, g: 1, b: 1},
	{name: "Magenta", r: 1, g: 0, b: 1},
}

func euclid(a, b int32) int32 {
	if b != 0 {
		return euclid(b, a%b)
	}
	return a
}

func formatMode(mode *glfw.GLFWvidmode) string {
	gcd := euclid(mode.Width, mode.Height)
	s := fmt.Sprintf("%d x %d x %d (%d:%d) (%d %d %d) %d Hz",
		mode.Width, mode.Height,
		mode.RedBits+mode.GreenBits+mode.BlueBits,
		mode.Width/gcd, mode.Height/gcd,
		mode.RedBits, mode.GreenBits, mode.BlueBits,
		mode.RefreshRate)

	return s
}

func framebuffer_size_callback(window *glfw.Window, width int, height int) {
	fmt.Printf("Framebuffer resized to %dx%d\n", width, height)
	gl.Viewport(0, 0, int32(width), int32(height))
}

func keyCallback4(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape {
		window.SetShouldClose(true)
		os.Exit(1)
	}
}

func list_modes(monitor *glfw.Monitor) {
	mode := monitor.GetVideoMode()
	modes := monitor.GetVideoModes()
	x, y := monitor.GetPos()
	width_mm, height_mm := monitor.GetPhysicalSize()
	xScale, yScale := monitor.GetContentScale()
	workAreaX, workAreaY, workAreaWidth, workAreaHeight := monitor.GetWorkarea()
	ps := "secondary"
	if monitor == glfw.GetPrimaryMonitor() {
		ps = "primary"
	}
	fmt.Printf("Name: %s (%s)\n", monitor.GetMonitorName(), ps)
	fmt.Printf("Current mode: %s\n", formatMode(&mode))
	fmt.Printf("Virtual position: %d, %d\n", x, y)
	fmt.Printf("Content scale: %f x %f\n", xScale, yScale)
	fmt.Printf("Physical size: %d x %d mm (%0.2f dpi at %d x %d)\n",
		width_mm, height_mm, float64(mode.Width)*25.4/float64(width_mm), mode.Width, mode.Height)
	fmt.Printf("Monitor work area: %d x %d starting at %d, %d\n\n",
		workAreaWidth, workAreaHeight, workAreaX, workAreaY)
	for i := 0; i < len(modes); i++ {
		fmt.Printf("%3d: %s", i, formatMode(&modes[i]))
		if mode == modes[i] {
			fmt.Printf(" (current mode)")
		}
		fmt.Printf("\n")
	}
}

func test_mode(monitor *glfw.Monitor, i int, timeShownSec float64) {
	modes := monitor.GetVideoModes()
	if i >= len(modes) {
		fmt.Printf("Trying to test mode %d, but monitor has only %d modes\n", i, len(modes))
		return
	}
	mode := modes[i]
	_ = glfw.WindowHint(glfw.Samples, 1)
	_ = glfw.WindowHint(glfw.RedBits, int(mode.RedBits))
	_ = glfw.WindowHint(glfw.GreenBits, int(mode.GreenBits))
	_ = glfw.WindowHint(glfw.BlueBits, int(mode.BlueBits))
	_ = glfw.WindowHint(glfw.RefreshRate, int(mode.RefreshRate))
	color := colors[i%len(colors)]
	fmt.Printf("Testing mode %d on monitor %s: %s color=%v\n", i, monitor.GetMonitorName(), formatMode(&mode), color)
	window, err := glfw.CreateWindow(int(mode.Width), int(mode.Height), "Video Mode Test", monitor, nil)
	if err != nil {
		fmt.Printf("Failed to enter mode %d: %s\n", i, formatMode(&mode))
		return
	}
	window.SetFramebufferSizeCallback(framebuffer_size_callback)
	window.SetKeyCallback(keyCallback4)
	window.MakeContextCurrent()
	err = gl.Init()
	if err != nil {
		panic(fmt.Sprintf("Failed to enter mode %d: %s\n", i, formatMode(&mode)))
	}
	glfw.SwapInterval(1)
	glfw.SetTime(0.0)
	gl.ClearColor(color.r, color.g, color.b, 1.0)
	for glfw.GetTime() < timeShownSec && !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)
		window.SwapBuffers()
		glfw.PollEvents()
		if window.ShouldClose() {
			fmt.Printf("User terminated program\n")
			glfw.Terminate()
			os.Exit(0)
		}
	}
	var current glfw.GLFWvidmode

	gl.GetIntegerv(gl.RED_BITS, &current.RedBits)
	gl.GetIntegerv(gl.GREEN_BITS, &current.GreenBits)
	gl.GetIntegerv(gl.BLUE_BITS, &current.BlueBits)
	if current.RedBits != mode.RedBits || current.GreenBits != mode.GreenBits || current.BlueBits != mode.BlueBits {
		fmt.Printf("*** Color bit mismatch: (%d %d %d) instead of (%d %d %d)\n",
			current.RedBits, current.GreenBits, current.BlueBits,
			mode.RedBits, mode.GreenBits, mode.BlueBits)
	}

	w, h := window.GetSize()
	current.Width, current.Height = int32(w), int32(h)
	if current.Width != mode.Width || current.Height != mode.Height {
		fmt.Printf("*** Size mismatch: %dx%d instead of %dx%d\n",
			current.Width, current.Height,
			mode.Width, mode.Height)
	}
	window.Destroy()
	window = nil
	glfw.PollEvents()
}

func MonitorMain() {
	fmt.Printf("\nMonitor test started\n")
	runtime.LockOSThread()
	// Initialize the glfw library
	err := glfw.Init()
	if err != nil {
		panic(err.Error())
	}
	defer glfw.Terminate()
	glfw.SetErrorCallback(error_callback)

	monitors := glfw.GetMonitors()
	for i := 0; i < len(monitors); i++ {
		list_modes(monitors[i])
	}
	fmt.Println("\npress escape to exit")
	// NB: On some monitors it takes up to 5 sec for the mode to be shown.
	// The following tests are testing just 4 modes, 2 on primary and 2 on secondary
	// It depends on the existence of the mode on the actual monitors.
	if len(monitors) > 1 {
		modes := monitors[1].GetVideoModes()
		test_mode(monitors[1], min(25, len(modes)-1), 5.0)
		test_mode(monitors[1], len(modes)-1, 5.0)
	}
	modes := monitors[0].GetVideoModes()
	test_mode(monitors[0], min(19, len(modes)-1), 5.0)
	test_mode(monitors[0], len(modes)-1, 5.0)

	glfw.Terminate()
	fmt.Printf("Monitor test finished\n")
}
