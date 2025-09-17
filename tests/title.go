package main

import (
	"fmt"
	"os"
	"runtime"

	glfw "github.com/jkvatne/purego-glfw"
	"github.com/neclepsio/gl/all-core/gl"
)

// UTF-8 window tear_title test

func TitleMain() {
	fmt.Printf("Windows title should show non-latin utf characters\n")
	runtime.LockOSThread()
	// Initialize the glfw library
	err := glfw.Init()
	if err != nil {
		panic(err.Error())
	}
	defer glfw.Terminate()
	glfw.SetErrorCallback(error_callback)

	window, err := glfw.CreateWindow(800, 400, "Temp", nil, nil)
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

	// Test clipboard
	glfw.SetClipboardString("Test")
	if err != nil {
		fmt.Printf("Could not set clipboard: %v\n", err)
	} else {
		fmt.Printf("Clipboard set to Test\n")
	}
	s := glfw.GetClipboardString()
	if err != nil {
		fmt.Printf("Could not gset clipboard: %v\n", err)
	} else {
		fmt.Printf("Clipboard read correctly the value %s\n", s)
	}
	if s != "Test" {
		panic("Invalid clipboard string")
	}

	glfw.SwapInterval(1)
	glfw.SetTime(0)
	_ = window.SetTitle("English 日本語 русский язык 官話")
	for !window.ShouldClose() && glfw.GetTime() < 2.0 {
		gl.ClearColor(100, 100, 0, 256)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		window.SwapBuffers()
		glfw.PollEvents()
	}
	window.Destroy()
	glfw.Terminate()
	fmt.Printf("Window tear_title test finished\n")
}
