package main

import (
	"flag"
)

var testName string

func main() {
	flag.StringVar(&testName, "t", "", "Test can be: tear_title, msaa, reopen, timeout, monitor")
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
