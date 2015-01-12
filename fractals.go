package main

import (
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"math"
	"runtime"
)

const (
	ul                                        = 256
	epsilon                                   = 0.0001
	fragmentShaderFile                        = "basic.frag"
	vertexShaderFile                          = "basic.vert"
	viewStepSize                              = 1.0 / 24.0
	zoomBase                                  = 1.03
	zoomLimit                                 = 10000
	startingWindowWidth, startingWindowHeight = 800, 800
)

type Point struct {
	x, y float64
}

var (
	poly           Polynomial
	vbo            gl.Buffer
	vao            gl.VertexArray
	vertices       []float32
	program        gl.Program
	positionAttrib gl.AttribLocation
	viewOrigin     *Point
	viewZoom       float64
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
	var delta float64
	if math.Abs(yoff) < 1 {
		delta = yoff / 10
	} else {
		delta = math.Abs(yoff) * yoff / 100
	}
	viewZoom += delta
	if viewZoom > zoomLimit {
		viewZoom = zoomLimit
	} else if viewZoom < -zoomLimit {
		viewZoom = -zoomLimit
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

func hsvToRgb(hsv [3]float32) [3]float32 {
	hsv_x, hsv_y, hsv_z := hsv[0], hsv[1], hsv[2]
	c := hsv_z * hsv_y
	h := hsv_x * 6.0
	x := c * (1 - float32(math.Abs(math.Mod(float64(h), 2)-1)))
	m := hsv_z - c
	var rgb [3]float32
	if h < 1 {
		rgb = [3]float32{c, x, 0}
	} else if h < 2 {
		rgb = [3]float32{x, c, 0}
	} else if h < 3 {
		rgb = [3]float32{0, c, x}
	} else if h < 4 {
		rgb = [3]float32{0, x, c}
	} else if h < 5 {
		rgb = [3]float32{x, 0, c}
	} else {
		rgb = [3]float32{c, 0, x}
	}
	rgb[0] += m
	rgb[1] += m
	rgb[2] += m
	return rgb
}

func makeColormap() []float32 {
	img := make([]float32, ul*3)
	for i := 0; i < ul; i++ {
		rgb := hsvToRgb([3]float32{float32(i) / (ul - 1), 1.0, 1.0})
		img[3*i] = rgb[0]
		img[3*i+1] = rgb[1]
		img[3*i+2] = rgb[2]
	}
	return img
}

func init() {
	poly = Polynomial{[]Root{
		PlainRoot(complex(-0.5, -0.86603)),
		PlainRoot(complex(-0.5, 0.86603)),
		PlainRoot(complex(1, 0)),
	}}

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

	gl.Enable(gl.TEXTURE_1D)
	colormap_texture := gl.GenTexture()
	colormap_texture.Bind(gl.TEXTURE_1D)

	gl.TexParameteri(gl.TEXTURE_1D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_1D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.TexParameteri(gl.TEXTURE_1D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_1D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	pixels := makeColormap()
	img := make([]uint8, len(pixels))
	for ix, pix := range pixels {
		img[ix] = uint8(math.Floor(float64(pix * (ul - 1))))
	}
	var imgArray [ul * 3]uint8
	copy(imgArray[:], img[0:len(img)])
	gl.TexImage1D(gl.TEXTURE_1D, 0, gl.RGB,
		ul, 0, gl.RGB, gl.UNSIGNED_BYTE,
		&imgArray)

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

	gl.ActiveTexture(gl.TEXTURE0 + 0)
	fmt.Printf("gl.TEXTURE0 is %d\n", gl.TEXTURE0)
	colormap_texture.Bind(gl.TEXTURE_1D)
	colormapLoc := program.GetUniformLocation("colormap")
	colormapLoc.Uniform1i(0)

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
