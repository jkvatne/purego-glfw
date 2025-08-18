// Multisample anti-aliasing test
// This test renders two high contrast, slowly rotating quads, one aliased
// and one (hopefully) anti-aliased, thus allowing for visual verification
// of whether MSAA is indeed enabled
package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"unsafe"

	"github.com/jkvatne/jkvgui/gl"
	glfw "github.com/jkvatne/purego-glfw"
)

// static const vec2 vertices[4] = {{ -0.6f, -0.6f },{  0.6f, -0.6f },{  0.6f,  0.6f },{ -0.6f,  0.6f }};
var vertices2 = [8]float32{-0.6, 0.6, 0.6, -0.6, -0.6, -0.6, 0.6, 0.6}

// {-0.6, -0.6, 0.6, -0.6, 0.6, 0.6, -0.6, 0.6}

var vertex_shader_text2 = `
#version 110
uniform mat4 MVP;
attribute vec2 vPos;
void main()
{
    gl_Position = MVP * vec4(vPos, 0.0, 1.0);	
}
` + "\x00"

var fragment_shader_text2 = `
#version 110
void main()
{
    gl_FragColor = vec4(0.0, 1.0, 1.0, 1.0);
}
` + "\x00"

func key_callback2(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Press {
		return
	}
	switch key {
	case glfw.KeySpace:
		glfw.SetTime(0.0)
	case glfw.KeyEscape:
		window.SetShouldClose(true)
	}
}

func checkError(sts uint32, program uint32, source string) {
	var status int32
	gl.GetShaderiv(program, sts, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(program, gl.INFO_LOG_LENGTH, &logLength)
		if logLength > 0 {
			txt := strings.Repeat("\x00", int(logLength+1))
			gl.GetShaderInfoLog(program, logLength, nil, gl.Str(txt))
			fmt.Printf("Shader error, source=%s, error=%s\n", source, txt)
		}
	}
}

func msaa() {
	var samples int32 = 4
	var window *glfw.Window
	var vertex_buffer2 uint32
	fmt.Printf("Multisample anti-aliasing test.\n")
	fmt.Printf("Notice bottom side aliasing at the left and not at the right side\n")

	// TODO glfwSetErrorCallback(error_callback);
	runtime.LockOSThread()

	err := glfw.Init()
	if err != nil {
		fmt.Println("init glfw error: " + err.Error())
		os.Exit(1)
	}
	fmt.Printf("Requesting MSAA with %d samples\n", samples)

	_ = glfw.WindowHint(glfw.Samples, samples)
	_ = glfw.WindowHint(glfw.ContextVersionMajor, 2)
	_ = glfw.WindowHint(glfw.ContextVersionMinor, 0)

	window, err = glfw.CreateWindow(800, 400, "Aliasing Detector", nil, nil)
	if err != nil {
		glfw.Terminate()
		fmt.Printf("CreateWindow error, " + err.Error())
	}
	if window == nil {
		glfw.Terminate()
		fmt.Printf("CreateWindow error, " + err.Error())
		os.Exit(2)
	}
	_ = window.SetKeyCallback(key_callback2)
	window.MakeContextCurrent()
	err = gl.Init()
	if err != nil {
		glfw.Terminate()
		fmt.Printf("gl Init error, " + err.Error())
	}
	glfw.SwapInterval(1)
	gl.Enable(gl.MULTISAMPLE)
	gl.GetIntegerv(gl.SAMPLES, &samples)
	if samples != 0 {
		fmt.Printf("Context reports MSAA is available with %d samples\n", samples)
	} else {
		fmt.Printf("Context reports MSAA is unavailable\n")
		os.Exit(0)
	}
	gl.GenBuffers(1, &vertex_buffer2)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertex_buffer2)
	gl.BufferData(gl.ARRAY_BUFFER, 32, unsafe.Pointer(&vertices2[0]), gl.STATIC_DRAW)

	vertex_shader2 := gl.CreateShader(gl.VERTEX_SHADER)
	csources, free := gl.Strs(vertex_shader_text2)
	gl.ShaderSource(vertex_shader2, 1, csources, nil)
	free()
	gl.CompileShader(vertex_shader2)
	CheckError(gl.COMPILE_STATUS, vertex_shader2, vertex_shader_text2)

	fragment_shader2 := gl.CreateShader(gl.FRAGMENT_SHADER)
	csources, free = gl.Strs(fragment_shader_text2)
	gl.ShaderSource(fragment_shader2, 1, csources, nil)
	free()
	gl.CompileShader(fragment_shader2)
	CheckError(gl.COMPILE_STATUS, fragment_shader2, fragment_shader_text2)

	program := gl.CreateProgram()
	gl.AttachShader(program, vertex_shader2)
	gl.AttachShader(program, fragment_shader2)
	gl.LinkProgram(program)
	checkError(gl.LINK_STATUS, program, "Linker")

	mvp_location := gl.GetUniformLocation(program, gl.Str("MVP\x00"))
	vpos_location := gl.GetAttribLocation(program, gl.Str("vPos\x00"))

	gl.EnableVertexAttribArray(uint32(vpos_location))
	gl.VertexAttribPointer(uint32(vpos_location), 2, gl.FLOAT, false, 4, nil)
	angle := 0.0
	glfw.SetTime(0)
	for !window.ShouldClose() && glfw.GetTime() < 5.0 {
		var m, p, mvp Mat4
		angle = angle + 0.0001
		if angle > 0.02 {
			angle = 0
		}

		width, height := window.GetFramebufferSize()
		ratio := float32(width) / float32(height)
		gl.Viewport(0, 0, width, height)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.UseProgram(program)
		p = mat4x4_ortho(-ratio, ratio, -1, 1, 0, 1)
		m = mat4x4_translate(-1, 0, 0)
		r := mat4x4_rotate_Z(m, float32(angle))
		mvp = mat4x4_mul(p, r)
		gl.UniformMatrix4fv(mvp_location, 1, false, &mvp[0])
		gl.Disable(gl.MULTISAMPLE)
		gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)

		m = mat4x4_translate(1, 0, 0)
		m = mat4x4_rotate_Z(m, float32(angle))
		mvp = mat4x4_mul(p, m)
		gl.UniformMatrix4fv(mvp_location, 1, false, &mvp[0])
		gl.Enable(gl.MULTISAMPLE)
		gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)
		window.SwapBuffers()
		glfw.PollEvents()
	}
	window.Destroy()
	glfw.Terminate()
	fmt.Printf("Multisample anti-aliasing test finished\n")
}
