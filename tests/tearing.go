// Vsync enabling test
// This test renders a high contrast, horizontally moving bar, allowing for
// visual verification of whether the set swap interval is indeed obeyed

package main

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"unsafe"

	glfw "github.com/jkvatne/purego-glfw"
	gl "github.com/jkvatne/purego-glfw/gl"
)

var (
	swap_tear  bool
	tear_title string
	frame_rate float64
)

func update_window_title(window *glfw.Window) {
	if swap_tear && swap_interval < 0 {
		tear_title = fmt.Sprintf("Tearing detector (interval %d (swap tear), %6.1f Hz)", swap_interval, frame_rate)
	} else {
		tear_title = fmt.Sprintf("Tearing detector (interval %d, %0.1f Hz)", swap_interval, frame_rate)
	}
	_ = window.SetTitle(tear_title)
}

func set_swap_interval(window *glfw.Window, interval int) {
	swap_interval = interval
	glfw.SwapInterval(swap_interval)
	update_window_title(window)
}

func key_callback10(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Press {
		return
	}
	switch key {
	case glfw.KeyUp:
		if swap_interval+1 > swap_interval {
			set_swap_interval(window, swap_interval+1)
		}
	case glfw.KeyDown:
		if swap_tear {
			if swap_interval-1 < swap_interval {
				set_swap_interval(window, swap_interval-1)
			}
		} else {
			if swap_interval-1 >= 0 {
				set_swap_interval(window, swap_interval-1)
			}
		}
	case glfw.KeyEscape:
		window.SetShouldClose(true)

	case glfw.KeyEnter:
		var x, y, width, height int
		if mods != glfw.ModAlt {
			return
		}
		if window.GetMonitor() != nil {
			window.SetMonitor(nil, x, y, width, height, 0)
		} else {
			monitor := glfw.GetPrimaryMonitor()
			mode := monitor.GetVideoMode()
			x, y = window.GetPos()
			w, h = window.GetSize()
			window.SetMonitor(monitor, 0, 0, int(mode.Width), int(mode.Height), int(mode.RefreshRate))
		}
	}
}

func TearingMain() {
	fmt.Println("Tearing test")
	runtime.LockOSThread()

	frame_count := 0
	var vertex_buffer2 uint32

	// Initialize the glfw library
	err := glfw.Init()
	if err != nil {
		panic(err.Error())
	}
	defer glfw.Terminate()
	glfw.SetErrorCallback(error_callback)

	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 0)
	window, err := glfw.CreateWindow(640, 480, "Tearing detector", nil, nil)
	if err != nil {
		glfw.Terminate()
		os.Exit(1)
	}
	window.MakeContextCurrent()
	window.SetKeyCallback(key_callback10)
	err = gl.Init()
	if err != nil {
		glfw.Terminate()
		fmt.Printf("Failed to initialize gl: %v", err)
	}
	set_swap_interval(window, 0)
	swap_tear = glfw.ExtensionSupported("WGL_EXT_swap_control_tear") || glfw.ExtensionSupported("GLX_EXT_swap_control_tear")
	fmt.Printf("Extension for controling tear is %v\n", swap_tear)
	gl.GenBuffers(1, &vertex_buffer2)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertex_buffer2)
	gl.BufferData(gl.ARRAY_BUFFER, 32, unsafe.Pointer(&vertices2[0]), gl.STATIC_DRAW)

	vertex_shader := gl.CreateShader(gl.VERTEX_SHADER)
	cSources, free := gl.Strs(vertex_shader_text2)
	gl.ShaderSource(vertex_shader, 1, cSources, nil)
	free()
	gl.CompileShader(vertex_shader)

	fragment_shader := gl.CreateShader(gl.FRAGMENT_SHADER)
	cSources, free = gl.Strs(fragment_shader_text2)
	gl.ShaderSource(fragment_shader, 1, cSources, nil)
	free()
	gl.CompileShader(fragment_shader)

	program := gl.CreateProgram()
	gl.AttachShader(program, vertex_shader)
	gl.AttachShader(program, fragment_shader)
	gl.LinkProgram(program)

	mvpLocation := gl.GetUniformLocation(program, gl.Str("MVP\x00"))
	vPosLocation := gl.GetAttribLocation(program, gl.Str("vPos\x00"))

	gl.EnableVertexAttribArray(uint32(vPosLocation))
	gl.VertexAttribPointer(uint32(vPosLocation), 2, gl.FLOAT, false, 4, nil)
	glfw.SetTime(0)
	for !window.ShouldClose() {
		position := math.Cos(float64(glfw.GetTime())*4.0) * 0.75
		width, height := window.GetFramebufferSize()
		gl.Viewport(0, 0, int32(width), int32(height))
		gl.Clear(gl.COLOR_BUFFER_BIT)
		p := mat4x4_ortho(-1.0, 1.0, -1.0, 1.0, 0.0, 1.0)
		m := mat4x4_translate(float32(position), 0.0, 0.0)
		mvp := mat4x4_mul(p, m)
		gl.UseProgram(program)
		gl.UniformMatrix4fv(mvpLocation, 1, false, &mvp[0])
		gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)
		window.SwapBuffers()
		glfw.PollEvents()
		frame_count++
		current_time := glfw.GetTime()
		if glfw.GetTime() > 5.0 {
			frame_rate = float64(frame_count) / current_time
			break
		}
		update_window_title(window)
	}

	glfw.Terminate()
}
