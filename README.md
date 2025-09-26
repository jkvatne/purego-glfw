# GLFW for Windows in pure Go 

## Introduction

This is a translation of the official glfw from C to Go.
It implements the same interface as found in https://github.com/go-gl/glfw/tree/master/v3.3/glfw
The original C code comes from https://github.com/glfw/glfw

Only Windows is currently supported because the other platforms can easily use the
orignal glfw. On Windows there is often problems using the CGO functions.
This code can be used without any C compiler. It is pure Go.

Some of the original tests are also translated, and seems to run fine.

The software is mostly complete. Some functions may be missing.
Please report any errors found.

## Installation

```
go get -u github.com/jkvatne/purego-glfw 
```

## Dependencies

```
golang.org/x/sys
```
This is used to import the types LazyDLL and LazyProc.

```
github.com/neclepsio/gl
```
This is used for testing only, in order to draw graphics in the test windows.
It is a fork of go-gl that does not use CGO. This means you can run the tests
with CGO disabled. (Set the environment variable CGO_ENABLED=0)

## Known limitations

- Only standard OpenGl is supported. No Vulkan or OpenGL ES.
- Only Windows 10 or later is supported (Can perhaps work on Windows 8 and 10).
- Monitor connect/disconnect is not detected while the app is running
- Joystick is not supported
