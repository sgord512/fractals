package main

import (
	"math"
)

type Root struct {
	base complex64
	exp  complex64
}

func PlainRoot(base complex64) Root {
	return Root{base, complex(1, 0)}
}

func NewRoot(base, exp complex64) Root {
	return Root{base, exp}
}

type rootFunction func(t float32) Root

var rootFunctions []rootFunction = []rootFunction{
	func(t float32) Root {
		whole, frac64 := math.Modf(float64(4.0 * t))
		tPrime := int(whole)
		frac := float32(frac64)
		if tPrime == 0 {
			return NewRoot(complex(1, 0), complex(frac, 0))
		} else if tPrime == 1 {
			return NewRoot(complex(1, 0), complex(1, frac))
		} else if tPrime == 2 {
			return NewRoot(complex(1, 0), complex((2.0-frac)/2.0, 1))
		} else {
			return NewRoot(complex(1, 0), complex(0.5, 1-frac))
		}
	},
	func(t float32) Root {
		return PlainRoot(complex(-0.5, -0.86603))
	},
	func(t float32) Root {
		return PlainRoot(complex(-0.5, 0.86603))
	},
}

func computePolynomial(time float32) Polynomial {
	roots := make([]Root, len(rootFunctions))
	for ix, rootFunc := range rootFunctions {
		roots[ix] = rootFunc(time)
	}
	return Polynomial{roots}
}

type Polynomial struct {
	roots []Root
}
