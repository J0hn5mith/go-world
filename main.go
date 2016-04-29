// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Renders a textured spinning cube using GLFW 3.1 and OpenGL 4.1 core forward-compatible profile.
package main

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"go/build"
	//"image"
	//"image/draw"
	//_ "image/png"
	"log"
	"os"
	"runtime"
	"strings"
	"github.com/go-gl/glfw/v3.1/glfw"
	"math/rand"
	"github.com/go-gl/mathgl/mgl32"
)

const windowWidth = 800
const windowHeight = 800

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
    render()
}

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}


//func newTexture(file string) (uint32, error) {
	//imgFile, err := os.Open(file)
	//if err != nil {
		//return 0, fmt.Errorf("texture %q not found on disk: %v", file, err)
	//}
	//img, _, err := image.Decode(imgFile)
	//if err != nil {
		//return 0, err
	//}

	//rgba := image.NewRGBA(img.Bounds())
	//if rgba.Stride != rgba.Rect.Size().X*4 {
		//return 0, fmt.Errorf("unsupported stride")
	//}
	//draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	//var texture uint32
	//gl.GenTextures(1, &texture)
	//gl.ActiveTexture(gl.TEXTURE0)
	//gl.BindTexture(gl.TEXTURE_2D, texture)
	//gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	//gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	//gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	//gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	//gl.TexImage2D(
		//gl.TEXTURE_2D,
		//0,
		//gl.RGBA,
		//int32(rgba.Rect.Size().X),
		//int32(rgba.Rect.Size().Y),
		//0,
		//gl.RGBA,
		//gl.UNSIGNED_BYTE,
		//gl.Ptr(rgba.Pix))

	//return texture, nil
//}


// Set the working directory to the root of Go package, so that its assets can be accessed.
func init() {
	dir, err := importPathToDir("github.com/go-gl/examples/glfw31-gl41core-cube")
	if err != nil {
		log.Fatalln("Unable to find Go package in your GOPATH, it's needed to load assets:", err)
	}
	err = os.Chdir(dir)
	if err != nil {
		log.Panicln("os.Chdir:", err)
	}
}

// importPathToDir resolves the absolute path from importPath.
// There doesn't need to be a valid Go package inside that import path,
// but the directory must exist.
func importPathToDir(importPath string) (string, error) {
	p, err := build.Import(importPath, "", build.FindOnly)
	if err != nil {
		return "", err
	}
	return p.Dir, nil
}

func render() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate() // What does the defer keyword mean?
	window, err := createWindow()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	// Configure the vertex and fragment shaders
	program, err := newProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}

	gl.UseProgram(program)


	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))


	camera := NewCamera(program)
	scene := NewScene(program)
	renderer := Renderer{camera}
	object := NewObject(createCircleGeometry(60))
	object.setScale(0.1)
	object.configure(scene.program)
	scene.addObject(object)

	object2 := NewObject(createTriangle2DGeometry(2))
	object2.setScale(0.1)
	object2.configure(scene.program)
	scene.addObject(object2)

    //particleSystem := NewParticleSystem(scene)
    //createParticle(particleSystem)
    //createParticle(particleSystem)
    //createParticle(particleSystem)


	for !window.ShouldClose() {

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

        //particleSystem.animate(0.05)
        renderer.render(scene)

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}

}

func createWindow() (*glfw.Window, error) {
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Cube", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	return window, err
}

func createParticle(particleSystem *ParticleSystem) *Particle{
    particle := particleSystem.newParticle()
	particle.object.setScale(.5)
	particle.velocity[0] = rand.Float32()-0.5
	particle.velocity[1] = rand.Float32()-0.5
	particle.velocity[2] = rand.Float32()-0.5
    return particle
}
