package main

// Test opacity of window

import (
	"fmt"
	"os"
	"runtime"
	"time"
	"unsafe"

	glfw "github.com/jkvatne/purego-glfw"
	gl "github.com/jkvatne/purego-glfw/gl"
)

func OpacityMain() {
	fmt.Printf("\nThis windows should gradualy be transparent\n")
	runtime.LockOSThread()
	// Initialize the glfw library
	err := glfw.Init()
	if err != nil {
		panic(err.Error())
	}
	defer glfw.Terminate()
	glfw.SetErrorCallback(error_callback)

	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 0)
	window, err := glfw.CreateWindow(400, 400, "Transparent", nil, nil)
	if err != nil {
		glfw.Terminate()
		os.Exit(1)
	}

	// Initialize Open-gl on current window
	window.MakeContextCurrent()
	err = gl.Init()
	if err != nil {
		glfw.Terminate()
		fmt.Printf("Could not init gl: %v\n", err)
		os.Exit(2)
	}
	glfw.SwapInterval(1)

	// Setup vertex shader
	var vertexShaderSource = `
		#version 110
		attribute vec2 vPos;
		void main()
		{
			gl_Position = vec4(vPos, 0.0, 1.0);	
		}
	` + "\x00"
	vertex_shader := gl.CreateShader(gl.VERTEX_SHADER)
	source, free := gl.Strs(vertexShaderSource)
	gl.ShaderSource(vertex_shader, 1, source, nil)
	free()
	gl.CompileShader(vertex_shader)

	// Setup fragment shader
	var fragmentShaderSource = `
		#version 110
		void main()
		{
			gl_FragColor = vec4(1.0);
		}
	` + "\x00"
	fragment_shader := gl.CreateShader(gl.FRAGMENT_SHADER)
	source, free = gl.Strs(fragmentShaderSource)
	gl.ShaderSource(fragment_shader, 1, source, nil)
	free()
	gl.CompileShader(fragment_shader)

	// Create program
	program := gl.CreateProgram()
	gl.AttachShader(program, vertex_shader)
	gl.AttachShader(program, fragment_shader)
	gl.LinkProgram(program)

	// Create and bind buffers
	var vertex_buffer uint32
	gl.GenBuffers(1, &vertex_buffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertex_buffer)

	// Create and setup attribute (vPos)
	vposLocation := gl.GetAttribLocation(program, gl.Str("vPos\x00"))
	gl.EnableVertexAttribArray(uint32(vposLocation))
	gl.VertexAttribPointer(uint32(vposLocation), 2, gl.FLOAT, false, 4, nil)

	glfw.SetTime(0)
	for !window.ShouldClose() && glfw.GetTime() < 4.0 {
		t := glfw.GetTime()
		o := min(t/2, max(0, min(1.0, 2.0-t/2)))
		window.SetOpacity(o)

		// Clear screen to a brown color
		gl.ClearColor(0.5, 0.5, 0, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// Draw a triangle
		var v = [3][2]float32{{-0.7, -0.7}, {0.7, 0.7}, {0.7, -0.7}}
		gl.UseProgram(program)
		gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(v)), unsafe.Pointer(&v), gl.STREAM_DRAW)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)
		gl.UseProgram(0)

		// Do end of frame housekeeping
		window.SwapBuffers()
		glfw.PollEvents()
		time.Sleep(time.Millisecond * 20)
	}
	window.Destroy()
	glfw.Terminate()
}
