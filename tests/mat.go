package main

import "math"

type Mat4 [16]float32

func mat4x4_ortho(left, right, bottom, top, near, far float32) Mat4 {
	var m Mat4
	// Implementation details would involve calculating the matrix elements
	// based on the orthographic projection formula.
	// This is a simplified representation.
	m[0] = 2 / (right - left)
	m[5] = 2 / (top - bottom)
	m[10] = -2 / (far - near)
	m[12] = -(right + left) / (right - left)
	m[13] = -(top + bottom) / (top - bottom)
	m[14] = -(far + near) / (far - near)
	m[15] = 1
	return m
}
func mat4x4_mul(a Mat4, b Mat4) (temp Mat4) {
	for c := 0; c < 4; c++ {
		for r := 0; r < 4; r++ {
			temp[c+4*r] = 0
			for k := 0; k < 4; k++ {
				temp[c+4*r] += a[k+4*r] * b[c+4*k]
			}
		}
	}
	return temp
}

func mat4x4_translate(x, y, z float32) Mat4 {
	return Mat4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		x, y, z, 1}
}

func mat4x4_rotate_Z(M Mat4, angle float32) Mat4 {
	s := float32(math.Sin(float64(angle)))
	c := float32(math.Cos(float64(angle)))
	R := Mat4{c, s, 0, 0,
		-s, c, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1}
	return mat4x4_mul(M, R)
}
