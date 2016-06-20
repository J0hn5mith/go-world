package go_world

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"strings"
)

type World struct {
	Scene   *Scene
	camera  *Camera
	physics *Physics
	window  *glfw.Window
	program uint32
}

func StartWorld(window *glfw.Window) (*Renderer, *World) {

	if err := gl.Init(); err != nil {
		panic(err)
	}
	program, err := newProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}

	gl.UseProgram(program)

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	width, height := window.GetSize()
	camera := NewCamera(program, width, height)
	scene := NewScene(program)
	renderer := NewRenderer(camera)

	world := new(World)
	world.Scene = scene
	world.window = window
	world.program = program
	world.camera = camera

	return renderer, world
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

func (w World) Program() uint32 {
	return w.program
}

func (w World) Camera() *Camera {
	return w.camera
}

func (world *World) Physics() *Physics {
	return world.physics
}

func (world *World) SetPhysics(physics *Physics) *World {
	world.physics = physics
    return world
}

func (w World) NewObject(geometry *Geometry) *Object {
	object := NewObject(geometry)
	return object
}

func (world *World) Update(timeDelta float32) *World {
	if world.physics != nil {
		world.physics.Update(timeDelta)
	}
    return world
}
