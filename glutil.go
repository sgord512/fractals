package main

import (
	"fmt"
	"github.com/go-gl/gl"
	"io/ioutil"
)

func compileShader(filename string, shaderType gl.GLenum) gl.Shader {
	shader := gl.CreateShader(shaderType)
	shaderSource, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error while reading shader from: %s\n", filename)
		panic(err)
	}
	shader.Source(string(shaderSource))
	shader.Compile()
	if shader.Get(gl.COMPILE_STATUS) == gl.FALSE {
		fmt.Printf("Failed to compile shader: %s with following error\n",
			filename)
		fmt.Println(shader.GetInfoLog())
		panic(fmt.Errorf("Fatal shader compilation error"))
	}
	return shader
}
