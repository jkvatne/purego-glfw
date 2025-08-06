// Monitor information tool
// This test prints monitor and video mode information or verifies video modes
package main

import (
	"fmt"
	"github.com/jkvatne/jkvgui/gl"
	glfw "github.com/jkvatne/purego-glfw"
	"math/rand/v2"
	"os"
	"runtime"
)

func euclid(a, b int32) int32 {
	if b != 0 {
		return euclid(b, a%b)
	}
	return a
}

func format_mode(mode *glfw.GLFWvidmode) string {
	gcd := euclid(mode.Width, mode.Height)
	s := fmt.Sprintf("%d x %d x %d (%d:%d) (%d %d %d) %d Hz",
		mode.Width, mode.Height,
		mode.RedBits+mode.GreenBits+mode.BlueBits,
		mode.Width/gcd, mode.Height/gcd,
		mode.RedBits, mode.GreenBits, mode.BlueBits,
		mode.RefreshRate)

	return s
}

func framebuffer_size_callback(window *glfw.Window, width int32, height int32) {
	fmt.Printf("Framebuffer resized to %ix%i\n", width, height)
	gl.Viewport(0, 0, width, height)
}

func key_callback4(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape {
		window.SetWindowShouldClose(true)
	}
}

func list_modes(monitor *glfw.Monitor) {
	mode := glfw.GetVideoMode(monitor)
	modes := glfw.GetVideoModes(monitor)
	x, y := monitor.GetPos()
	width_mm, height_mm := monitor.GetPhysicalSize()
	xscale, yscale := monitor.GetContentScale()
	workarea_x, workarea_y, workarea_width, workarea_height := monitor.GetWorkarea()
	ps := "secondary"
	if monitor == glfw.GetPrimaryMonitor() {
		ps = "primary"
	}
	fmt.Printf("Name: %s (%s)\n", monitor.GetMonitorName(), ps)
	fmt.Printf("Current mode: %s\n", format_mode(&mode))
	fmt.Printf("Virtual position: %d, %d\n", x, y)
	fmt.Printf("Content scale: %f x %f\n", xscale, yscale)
	fmt.Printf("Physical size: %d x %d mm (%0.2f dpi at %d x %d)\n",
		width_mm, height_mm, float64(mode.Width)*25.4/float64(width_mm), mode.Width, mode.Height)
	fmt.Printf("Monitor work area: %d x %d starting at %d, %d\n\n",
		workarea_width, workarea_height, workarea_x, workarea_y)
	fmt.Printf("Modes:\n")
	for i := 0; i < len(modes); i++ {
		fmt.Printf("%3d: %s", i, format_mode(&modes[i]))
		if mode == modes[i] {
			fmt.Printf(" (current mode)\n")
		}
		fmt.Printf("\n")
	}
}

func test_modes(monitor *glfw.Monitor) {
	modes := glfw.GetVideoModes(monitor)
	for i := 0; i < len(modes); i++ {
		mode := modes[i]
		// TODO _ = glfw.WindowHint(glfw.Samples, mode.Samples)
		_ = glfw.WindowHint(glfw.RedBits, mode.RedBits)
		_ = glfw.WindowHint(glfw.GreenBits, mode.GreenBits)
		_ = glfw.WindowHint(glfw.BlueBits, mode.BlueBits)
		_ = glfw.WindowHint(glfw.RefreshRate, mode.RefreshRate)

		fmt.Printf("Testing mode %d on monitor %s: %s\n", i, monitor.GetMonitorName(), format_mode(&mode))
		window, err := glfw.CreateWindow(mode.Width, mode.Height, "Video Mode Test", monitor, nil)
		if err != nil {
			fmt.Printf("Failed to enter mode %d: %s\n", i, format_mode(&mode))
			continue
		}

		// TODO glfw.SetFramebufferSizeCallback(window, framebuffer_size_callback);
		window.SetKeyCallback(key_callback)
		window.MakeContextCurrent()
		err = gl.Init()
		if err != nil {
			fmt.Printf("Failed to enter mode %d: %s\n", i, format_mode(&mode))
		}
		glfw.SwapInterval(1)
		glfw.SetTime(0.0)
		gl.ClearColor(rand.Float32(), rand.Float32(), rand.Float32(), 1)
		for glfw.GetTime() < 1.0 {
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

		current.Width, current.Height = window.GetSize()
		if current.RedBits != mode.RedBits || current.GreenBits != mode.GreenBits || current.BlueBits != mode.BlueBits {
			fmt.Printf("*** Color bit mismatch: (%d %d %d) instead of (%d %d %id)\n",
				current.RedBits, current.GreenBits, current.BlueBits,
				mode.RedBits, mode.GreenBits, mode.BlueBits)
		}
		if current.Width != mode.Width || current.Height != mode.Height {
			fmt.Printf("*** Size mismatch: %dx%d instead of %dx%d\n",
				current.Width, current.Height,
				mode.Width, mode.Height)
		}
		window.Destroy()
		window = nil
		glfw.PollEvents()
	}
}

func monitor() {
	fmt.Printf("Monitor test started\n")
	runtime.LockOSThread()
	glfw.SetErrorCallback(error_callback)
	err := glfw.Init()
	if err != nil {
		panic("glfw.Init error: " + err.Error())
	}

	mons := glfw.GetMonitors()
	for i := 0; i < len(mons); i++ {
		list_modes(mons[i])
	}
	// TODO test both monitors
	for i := 1; i < len(mons); i++ {
		test_modes(mons[i])
	}
	glfw.Terminate()
	fmt.Printf("Monitor test finished\n")
}
