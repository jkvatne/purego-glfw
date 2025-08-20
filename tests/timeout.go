// Event wait timeout test
// This test is intended to verify that waiting for events with timeout works
package main

import (
	"fmt"
	"runtime"

	"github.com/go-gl/gl/all-core/gl"
	glfw "github.com/jkvatne/purego-glfw"
	"golang.org/x/exp/rand"
)

func error_callback(error int, description string) {
	fmt.Printf("Error: %s\n", description)
}

func key_callback3(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}
}

func timeout() {
	fmt.Printf("Test timeouts at 1 second. Window should change color each second\n")
	runtime.LockOSThread()
	glfw.SetErrorCallback(error_callback)
	err := glfw.Init()
	if err != nil {
		panic("glfw.Init error: " + err.Error())
	}
	window, err := glfw.CreateWindow(640, 480, "Event Wait Timeout Test", nil, nil)
	if err != nil {
		glfw.Terminate()
		fmt.Printf("Could not create window: %v\n", err)
	}
	window.MakeContextCurrent()
	gl.Init()
	window.SetKeyCallback(key_callback)
	glfw.SetTime(0)
	for !window.ShouldClose() && glfw.GetTime() < 5.0 {
		width, height := window.GetFramebufferSize()
		gl.Viewport(0, 0, int32(width), int32(height))
		gl.ClearColor(rand.Float32(), rand.Float32(), rand.Float32(), 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		window.SwapBuffers()
		glfw.WaitEventsTimeout(1.0)
	}
	window.Destroy()
	glfw.Terminate()
	fmt.Printf("Timeout test finished\n")
}
