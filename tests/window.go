package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/jkvatne/jkvgui/gl"
	glfw "github.com/jkvatne/purego-glfw"
)

func key_callback_window(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Press {
		return
	}
	switch key {
	case glfw.KeyI:
		printInfo(window)
	}
}

func printInfo(window *glfw.Window) {
	w, h = window.GetFramebufferSize()
	fmt.Printf("Window Framebuffer Size:%dx%d\n", w, h)
	sx, sy := window.GetContentScale()
	fmt.Printf("ContentScale: %0.3fx%0.3f\n", sx, sy)
	l, t, r, b := window.GetFrameSize()
	fmt.Printf("Frame size: l=%d t=%d r=%d b=%d\n", l, t, r, b)
	o := window.GetOpacity()
	fmt.Printf("Opacity: %v\n", o)
	decorated := window.GetAttrib(glfw.Decorated)
	fmt.Printf("Decorated: %v\n", decorated)
	resizeable := window.GetAttrib(glfw.Resizable)
	fmt.Printf("Resizable: %v\n", resizeable)
	floating := window.GetAttrib(glfw.Floating)
	fmt.Printf("Floating: %v\n", floating)
	mousePassthrough := window.GetAttrib(glfw.MousePassthrough)
	fmt.Printf("MousePassthrough: %v\n", mousePassthrough)
	autoIconify := window.GetAttrib(glfw.AutoIconify)
	fmt.Printf("AutoIconify: %v\n", autoIconify)
	focused := window.GetAttrib(glfw.Focused)
	fmt.Printf("Focused: %v\n", focused)
	hovered := window.GetAttrib(glfw.Hovered)
	fmt.Printf("Hovered: %v\n", hovered)
	visible := window.GetAttrib(glfw.Visible)
	fmt.Printf("Visible: %v\n", visible)
	iconified := window.GetAttrib(glfw.Iconified)
	fmt.Printf("Iconified: %v\n", iconified)
	maximized := window.GetAttrib(glfw.Maximized)
	fmt.Printf("Maximized: %v\n", maximized)
	transparent := window.GetAttrib(glfw.TransparentFramebuffer)
	fmt.Printf("Transparent: %v\n", transparent)
}

func window() {
	fmt.Printf("\nPrint some window info\n")
	runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	window, err := glfw.CreateWindow(800, 400, "Window info", nil, nil)
	if err != nil {
		glfw.Terminate()
		os.Exit(1)
	}
	window.SetKeyCallback(key_callback_window)
	window.MakeContextCurrent()
	gl.Init()
	glfw.SwapInterval(1)
	glfw.SetTime(0)
	for !window.ShouldClose() {
		gl.ClearColor(100, 100, 120, 256)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		window.SwapBuffers()
		glfw.PollEvents()
		if glfw.GetTime() > 2.0 {
			printInfo(window)
			break
		}
	}
	window.Destroy()
	glfw.Terminate()
}
