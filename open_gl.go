package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"log"
	"math/rand"
)

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

	camera := NewCamera(program)

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	position := *new(mgl32.Vec4)
	position[0] = 1

	scene := NewScene(program)
    particleSystem := NewParticleSystem(scene)
    createParticle(particleSystem)
    createParticle(particleSystem)
    createParticle(particleSystem)


	for !window.ShouldClose() {

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

        particleSystem.animate(0.05)
        camera.render(scene)

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
