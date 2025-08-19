// Window re-opener (open/close stress test)
package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"runtime"
	"strings"
	"unsafe"

	"github.com/jkvatne/jkvgui/gl"
	glfw "github.com/jkvatne/purego-glfw"
)

// It closes and re-opens the GLFW window every five seconds, alternating
// between windowed and full screen mode
// It also times and logs opening and closing actions and attempts to separate
// user initiated window closing from its own

var vertex_shader_text = `
#version 110
uniform mat4 MVP;
attribute vec2 vPos;
void main()
{
    gl_Position = MVP * vec4(vPos, 0.0, 1.0);	
}
` + "\x00"

var fragment_shader_text = `
#version 110
void main()
{
    gl_FragColor = vec4(0.0, 1.0, 1.0, 1.0);
}
` + "\x00"

var vertices = [4][2]float32{{-0.5, -0.5}, {0.5, -0.5}, {0.5, 0.5}, {-0.5, 0.5}}

func window_close_callback(window *glfw.Window) {
	fmt.Printf("Close callback triggered\n")
}

func key_callback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Press {
		return
	}
	switch key {
	case glfw.KeyQ:
	case glfw.KeyEscape:
		window.SetShouldClose(true)
	}
}

func close_window(window *glfw.Window) {
	base := glfw.GetTime()
	window.Destroy()
	fmt.Printf("Closing window took %0.3f seconds\n", glfw.GetTime()-base)
}

func CheckError(sts uint32, program uint32, source string) {
	var status int32
	gl.GetShaderiv(program, sts, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(program, gl.INFO_LOG_LENGTH, &logLength)
		if logLength > 0 {
			txt := strings.Repeat("\x00", int(logLength+1))
			gl.GetShaderInfoLog(program, logLength, nil, gl.Str(txt))
			fmt.Printf("Shader error %s, Source: %s", txt, source)
		}
	}
}

func reopen() {
	runtime.LockOSThread()
	fmt.Printf("Test reopening windows, including full-screen\n")
	count := 1
	var monitor *glfw.Monitor
	err := glfw.Init()
	if err != nil {
		panic("glfwInit err: " + err.Error())
	}
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 0)
	glfw.SetTime(0)
	for glfw.GetTime() < 5.0 {
		monitor = nil
		if count&1 == 0 {
			monitors := glfw.GetMonitors()
			monitor = monitors[rand.Int()%len(monitors)]
		}
		width := 400
		height := 400
		if monitor != nil {
			mode := monitor.GetVideoMode()
			width = int(mode.Width)
			height = int(mode.Height)
			x, y := monitor.GetPos()
			fmt.Printf("Monitor x=%d, y=%d, w=%d, h=%d", x, y, mode.Width, mode.Height)
		}
		base := glfw.GetTime()
		window, err := glfw.CreateWindow(width, height, "Window Re-opener", monitor, nil)
		if err != nil {
			glfw.Terminate()
			fmt.Printf("Could not create window: %v\n", err)
		}
		if monitor != nil {
			fmt.Printf("Opening full screen window on monitor %s took %0.3f seconds\n",
				monitor.GetMonitorName(),
				glfw.GetTime()-base)
		} else {
			fmt.Printf("Opening regular window took %0.3f seconds\n", glfw.GetTime()-base)
		}
		window.MakeContextCurrent()
		gl.Init()
		glfw.SwapInterval(1)

		gl.ClearColor(100, 100, 0, 256)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		// window.SwapBuffers()
		window.SetWindowCloseCallback(window_close_callback)
		window.SetKeyCallback(key_callback)
		window.MakeContextCurrent()
		gl.Init()
		glfw.SwapInterval(1)
		vertex_shader := gl.CreateShader(gl.VERTEX_SHADER)
		csources, free := gl.Strs(vertex_shader_text)
		gl.ShaderSource(vertex_shader, 1, csources, nil)
		free()
		gl.CompileShader(vertex_shader)
		CheckError(gl.COMPILE_STATUS, vertex_shader, vertex_shader_text)
		fragment_shader := gl.CreateShader(gl.FRAGMENT_SHADER)
		csources, free = gl.Strs(fragment_shader_text)
		gl.ShaderSource(fragment_shader, 1, csources, nil)
		free()
		gl.CompileShader(fragment_shader)
		CheckError(gl.COMPILE_STATUS, fragment_shader, fragment_shader_text)
		program := gl.CreateProgram()
		gl.AttachShader(program, vertex_shader)
		gl.AttachShader(program, fragment_shader)
		gl.LinkProgram(program)
		CheckError(gl.LINK_STATUS, program, "Linker")
		mvp_location := gl.GetUniformLocation(program, gl.Str("MVP\x00"))
		vpos_location := gl.GetAttribLocation(program, gl.Str("vPos\x00"))
		var vertex_buffer uint32
		gl.GenBuffers(1, &vertex_buffer)
		gl.BindBuffer(gl.ARRAY_BUFFER, vertex_buffer)
		gl.BufferData(gl.ARRAY_BUFFER, 2*4*4, unsafe.Pointer(&vertices[0][0]), gl.STATIC_DRAW)

		gl.EnableVertexAttribArray(uint32(vpos_location))
		gl.VertexAttribPointer(uint32(vpos_location), 2, gl.FLOAT, false, 2*4, nil)

		t := glfw.GetTime()
		for {
			if glfw.GetTime() > t+1.0 {
				break
			}
			w, h := window.GetFramebufferSize()
			ratio := float32(w) / float32(h)
			gl.Viewport(0, 0, int32(w), int32(h))
			gl.Clear(gl.COLOR_BUFFER_BIT)
			p := mat4x4_ortho(-ratio, ratio, -1, 1, 0, 1)
			m := Mat4{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
			m = mat4x4_rotate_Z(m, float32(glfw.GetTime()))
			mvp := mat4x4_mul(p, m)
			gl.UseProgram(program)
			gl.UniformMatrix4fv(mvp_location, 1, false, &mvp[0])
			gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)
			window.SwapBuffers()
			glfw.PollEvents()
			if window.ShouldClose() {
				close_window(window)
				fmt.Printf("User closed window\n")
				glfw.Terminate()
				os.Exit(0)
			}
		}
		close_window(window)
		count++
	}
	glfw.Terminate()
	fmt.Printf("Re-opening test finished\n")
}
