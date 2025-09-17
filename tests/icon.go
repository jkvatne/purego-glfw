// Window icon test program
// This program is used to test the icon feature.
package main

import (
	"fmt"
	"image"
	"runtime"
	"time"

	glfw "github.com/jkvatne/purego-glfw"
	"github.com/neclepsio/gl/all-core/gl"
)

// a simple glfw logo
var logo = [16]string{
	"................",
	"................",
	"...0000..0......",
	"...0.....0......",
	"...0.00..0......",
	"...0..0..0......",
	"...0000..0000...",
	"................",
	"................",
	"...000..0...0...",
	"...0....0...0...",
	"...000..0.0.0...",
	"...0....0.0.0...",
	"...0....00000...",
	"................",
	"................",
}

var icon_colors = [5][4]uint8{
	{0, 0, 0, 255},       // black
	{255, 0, 0, 255},     // red
	{0, 255, 0, 255},     // green
	{0, 0, 255, 255},     // blue
	{255, 255, 255, 255}, // white
}

var cur_icon_color int

func set_icon(window *glfw.Window, icon_color int) {
	var images []image.Image
	img := image.NewNRGBA(image.Rect(0, 0, 16, 16))
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			for i := 0; i < 4; i++ {
				if logo[y][x] == '0' {
					img.Pix[(x+(y*16))*4+i] = icon_colors[icon_color][i]
				} else {
					img.Pix[(x+(y*16))*4+i] = 0
				}
			}
		}
	}
	images = append(images, img)
	window.SetIcon(images)
}

func key_callback8(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Press {
		return
	}
	switch key {
	case glfw.KeyEscape:
		window.SetShouldClose(true)
	case glfw.KeySpace:
		cur_icon_color = (cur_icon_color + 1) % 5
		set_icon(window, cur_icon_color)
	case glfw.KeyX:
		window.SetIcon(nil)
	}
}

func IconMain() {
	fmt.Println("\nWindows icon test. Changes the icon in the upper left corner of the window")
	fmt.Println("and the icon on the task bar. Press space to change color, escape to exit, X to clear icon.")
	runtime.LockOSThread()

	// Initialize the glfw library
	err := glfw.Init()
	if err != nil {
		panic(err.Error())
	}
	defer glfw.Terminate()
	glfw.SetErrorCallback(error_callback)

	window, err := glfw.CreateWindow(500, 500, "Window Icon", nil, nil)
	if err != nil {
		panic("Failed to create window, " + err.Error())
	}
	window.MakeContextCurrent()
	gl.Init()
	if err != nil {
		panic("gl Init error, " + err.Error())
	}
	window.SetKeyCallback(key_callback8)
	set_icon(window, cur_icon_color)

	for !window.ShouldClose() && glfw.GetTime() < 5.0 {
		set_icon(window, cur_icon_color)
		gl.ClearColor(1, 0.5, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		window.SwapBuffers()
		time.Sleep(time.Second)
		glfw.PollEvents()
		cur_icon_color = (cur_icon_color + 1) % 5

	}

	window.Destroy()
	glfw.Terminate()
}
