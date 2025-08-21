package main

import (
	"flag"
)

var testname string

func main() {
	flag.StringVar(&testname, "t", "", "Test can be: tear_title, msaa, reopen, timeout, monitor")
	flag.Parse()
	switch testname {
	case "tear_title":
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
		opacity()
		title()
		msaa()
		timeout()
		reopen()
		icon()
		threads()
		windowinfo()
		monitor()
		cursor()
	}
}
