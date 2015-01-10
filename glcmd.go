package main { 

import (
	"C"
	"github.com/go-gl/gl"
)

type glvec interface {
	sendVertex()
}

type v2f struct {
	x, y C.float
}

type v2s struct {
	x, y C.short
}

type v2i struct {
	x, y C.int
}

type v2d struct {
	x, y C.double
}

type v3f struct {
	x, y, z C.float
}

type v3s struct {
	x, y, z C.short
}

type v3i struct {
	x, y, z C.int
}

type v3d struct {
	x, y, z C.double
}

type v4f struct {
	x, y, z, w C.float
}

type v4s struct {
	x, y, z, w C.short
}

type v4i struct {
	x, y, z, w C.int
}

type v4d struct {
	x, y, z, w C.double
}


type glShape interface {
	sendVertices([]glvec)
}
