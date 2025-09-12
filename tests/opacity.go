package main

// Test opacity of window

import (
	"fmt"
	"os"
	"runtime"
	"time"
	"unsafe"

	glfw "github.com/jkvatne/purego-glfw"
	"github.com/neclepsio/gl/all-core/gl"
)

func opacity() {
	fmt.Printf("\nThis windows should gradualy be transparent\n")
	runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 0)
	window, err := glfw.CreateWindow(400, 400, "Transparent", nil, nil)
	if err != nil {
		glfw.Terminate()
		os.Exit(1)
	}
	window.MakeContextCurrent()
	err = gl.Init()
	if err != nil {
		glfw.Terminate()
		fmt.Printf("Could not init gl: %v\n", err)
		os.Exit(2)
	}
	glfw.SwapInterval(1)
	glfw.SetTime(0)

	var vertexShaderSource = `
		#version 110
		uniform mat4 MVP;
		attribute vec2 vPos;
		void main()
		{
			gl_Position = MVP * vec4(vPos, 0.0, 1.0);	
		}
	` + "\x00"

	var fragmentShaderSource = `
		#version 110
		void main()
		{
			gl_FragColor = vec4(1.0);
		}
	` + "\x00"
	var vertex_buffer uint32
	gl.GenBuffers(1, &vertex_buffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertex_buffer)
	vertex_shader := gl.CreateShader(gl.VERTEX_SHADER)
	csources, free := gl.Strs(vertexShaderSource)
	gl.ShaderSource(vertex_shader, 1, csources, nil)
	free()
	gl.CompileShader(vertex_shader)

	fragment_shader := gl.CreateShader(gl.FRAGMENT_SHADER)
	csources, free = gl.Strs(fragmentShaderSource)
	gl.ShaderSource(fragment_shader, 1, csources, nil)
	free()
	gl.CompileShader(fragment_shader)

	program := gl.CreateProgram()
	gl.AttachShader(program, vertex_shader)
	gl.AttachShader(program, fragment_shader)
	gl.LinkProgram(program)

	mvp_location := gl.GetUniformLocation(program, gl.Str("MVP\x00"))
	vpos_location := gl.GetAttribLocation(program, gl.Str("vPos\x00"))

	gl.EnableVertexAttribArray(uint32(vpos_location))
	gl.VertexAttribPointer(uint32(vpos_location), 2, gl.FLOAT, false, 8, nil)
	gl.UseProgram(program)

	for !window.ShouldClose() && glfw.GetTime() < 4.0 {
		gl.ClearColor(0.5, 0.5, 0, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		ww, hh := window.GetFramebufferSize()
		gl.Viewport(0, 0, int32(ww), int32(hh))
		var v = [3][2]float32{{100, 100}, {300, 100}, {300, 300}}
		gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(v)), unsafe.Pointer(&v), gl.STREAM_DRAW)
		mvp := mat4x4_ortho(0.0, 300, 0.0, 300, 0.0, 1.0)
		gl.UniformMatrix4fv(mvp_location, 1, false, &mvp[0])
		gl.DrawArrays(gl.TRIANGLES, 0, 3)
		window.SwapBuffers()
		glfw.PollEvents()
		window.SetOpacity(min(glfw.GetTime(), max(0, min(1.0, 4.0-glfw.GetTime()))))
		time.Sleep(time.Millisecond * 20)
	}
	window.Destroy()
	glfw.Terminate()
}
