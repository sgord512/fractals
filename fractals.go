// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Originally put together by github.com/segfault88, but
// I thought it might be useful to somebody else too.

// It took me quite a lot of frustration and messing around
// to get a basic example of glfw3 with modern OpenGL (3.3)
// with shaders etc. working. Hopefully this will save you
// some trouble. Enjoy!

package main

import (
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"math"
	"runtime"
)

const (
	ul                                        = 255
	epsilon                                   = 0.0001
	fragmentShaderFile                        = "basic.frag"
	vertexShaderFile                          = "basic.vert"
	viewStepSize                              = 1.0 / 24.0
	zoomBase                                  = 1.01
	startingWindowWidth, startingWindowHeight = 800, 800
)

type Point struct {
	x, y float64
}

var (
	poly             Polynomial
	vbo              gl.Buffer
	vao              gl.VertexArray
	vertices         []float32
	program          gl.Program
	positionAttrib   gl.AttribLocation
	rootColors       []float32
	defaultRootColor []float32
	viewOrigin       *Point
	viewZoom         float64
)

func draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.DrawArrays(gl.TRIANGLES, 0, len(vertices)/3)
}

func setRoots() {
	rootCount := len(poly.roots)

	rootCountLoc := program.GetUniformLocation("rootCount")
	rootCountLoc.Uniform1i(rootCount)

	rootBaseArray := make([]float32, rootCount*2)
	for ix, root := range poly.roots {
		rootBaseArray[2*ix] = real(root.base)
		rootBaseArray[2*ix+1] = imag(root.base)
	}

	rootBaseLoc := program.GetUniformLocation("rootBase")
	rootBaseLoc.Uniform2fv(rootCount, rootBaseArray)

	rootColorLoc := program.GetUniformLocation("rootColor")
	rootColorLoc.Uniform3fv(len(rootColors)/3, rootColors)

	defaultRootColorLoc := program.GetUniformLocation("defaultRootColor")
	defaultRootColorLoc.Uniform3fv(1, defaultRootColor)
}

func reshape(width, height int) {
	uniformSizeLoc := program.GetUniformLocation("windowSize")
	widthF, heightF := float32(width), float32(height)
	uniformSizeLoc.Uniform2fv(1, []float32{widthF, heightF})
}

func setupCamera() {
	uniformViewOriginLoc := program.GetUniformLocation("viewOrigin")
	uniformViewOriginLoc.Uniform2fv(1, []float32{
		float32(viewOrigin.x), float32(viewOrigin.y)})

	zoom := math.Pow(zoomBase, viewZoom)
	uniformViewZoomLoc := program.GetUniformLocation("viewZoom")
	uniformViewZoomLoc.Uniform1f(float32(zoom))
}

func scroll(window *glfw.Window, xoff, yoff float64) {
	fmt.Printf("yoff: %G\n", yoff)
	fmt.Printf("viewZoom: %G\n", viewZoom)
	viewZoom += yoff / 10
	if viewZoom > 1000 {
		viewZoom = 1000
	} else if viewZoom < -1000 {
		viewZoom = -1000
	}
}

func key(window *glfw.Window, k glfw.Key, s int, action glfw.Action, mods glfw.ModifierKey) {
	if !(action == glfw.Press || action == glfw.Repeat) {
		return
	}

	zoom := math.Pow(zoomBase, viewZoom)
	viewDelta := viewStepSize / zoom

	switch glfw.Key(k) {
	case glfw.KeyEscape:
		window.SetShouldClose(true)
	case glfw.KeyI:
		fmt.Println("up")
		viewOrigin.y += viewDelta
	case glfw.KeyK:
		fmt.Println("down")
		viewOrigin.y -= viewDelta
	case glfw.KeyJ:
		fmt.Println("left")
		viewOrigin.x -= viewDelta
	case glfw.KeyL:
		fmt.Println("right")
		viewOrigin.x += viewDelta
	default:
		return
	}
}

func init() {
	poly = Polynomial{[]Root{
		PlainRoot(complex(-0.5, -0.86603)),
		PlainRoot(complex(-0.5, 0.86603)),
		PlainRoot(complex(1, 0)),
	}}

	rootColors = []float32{
		1.0, 0.0, 0.0,
		0.0, 1.0, 0.0,
		0.0, 0.0, 1.0,
	}
	defaultRootColor = []float32{
		0.0, 0.0, 0.0,
	}

	viewOrigin = &Point{0.0, 0.0}
	// viewZoom is the exponent used to compute magnification, so
	// zoom = zoomBase ^ viewZoom
	viewZoom = 0

}

func main() {
	// lock glfw/gl calls to a single thread
	runtime.LockOSThread()

	glfw.Init()
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenglForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenglProfile, glfw.OpenglCoreProfile)

	window, err := glfw.CreateWindow(startingWindowWidth, startingWindowHeight, "Example", nil, nil)
	if err != nil {
		panic(err)
	}

	defer window.Destroy()

	window.SetKeyCallback(key)
	window.SetScrollCallback(scroll)

	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	gl.Init()

	vao = gl.GenVertexArray()
	vao.Bind()

	vbo = gl.GenBuffer()
	vbo.Bind(gl.ARRAY_BUFFER)

	vertices = []float32{
		-1, 1, 0,
		-1, -1, 0,
		1, -1, 0,

		1, -1, 0,
		1, 1, 0,
		-1, 1, 0,
	}

	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, vertices, gl.STATIC_DRAW)

	vertex_shader := compileShader(vertexShaderFile, gl.VERTEX_SHADER)
	defer vertex_shader.Delete()

	fragment_shader := compileShader(fragmentShaderFile, gl.FRAGMENT_SHADER)
	defer fragment_shader.Delete()

	program = gl.CreateProgram()
	program.AttachShader(vertex_shader)
	program.AttachShader(fragment_shader)

	program.BindFragDataLocation(0, "outColor")
	program.Link()
	program.Use()
	defer program.Delete()

	ulLoc := program.GetUniformLocation("ul")
	ulLoc.Uniform1i(ul)

	epsilonLoc := program.GetUniformLocation("epsilon")
	epsilonLoc.Uniform1f(epsilon)

	positionAttrib = program.GetAttribLocation("position")
	positionAttrib.AttribPointer(3, gl.FLOAT, false, 0, nil)
	positionAttrib.EnableArray()
	defer positionAttrib.DisableArray()

	gl.ClearColor(0.3, 0.3, 0.3, 1.0)

	for !window.ShouldClose() {
		width, height := window.GetFramebufferSize()
		reshape(width, height)
		setupCamera()
		setRoots()
		draw()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
