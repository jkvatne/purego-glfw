// This is the tests for the purego-glfw implementations.
// Type go run ./tests from the main directory.
// To avoid using gcc the "github.com/neclepsio/gl/all-core/gl" is used
// instead of the "github.com/go-gl/gl/all-core/gl" repository.
package main

import (
	"flag"
)

var testName string

func main() {
	flag.StringVar(&testName, "t", "", "Test can be: -tearing, -window, -cursor, -threadds, -msaa, -reopen, -timeout, -monitor, -opacity")
	flag.Parse()
	switch testName {
	case "title":
		title()
	case "reopen":
		reopen()
	case "msaa":
		msaa()
	case "timeout":
		timeout()
	case "monitor":
		monitor()
	case "threads":
		threads()
	case "icon":
		icon()
	case "cursor":
		cursor()
	case "window":
		windowinfo()
	case "opacity":
		opacity()
	case "tearing":
		tearing()
	default:
		title()
		threads()
		opacity()
		msaa()
		timeout()
		reopen()
		icon()
		windowinfo()
		monitor()
		cursor()
	}
}
