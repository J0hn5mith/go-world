package go_world

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	//"log"
	//"runtime"
	"strings"
)

const windowWidth = 800
const windowHeight = 800

type World struct {
	Scene   *Scene
	camera  *Camera
	window  *glfw.Window
	program uint32
}

func StartWorld(window *glfw.Window) (Renderer, *World) {

	// Initialize Glow
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

	camera := NewCamera(program)
	scene := NewScene(program)
	renderer := Renderer{camera}

	//object := NewObject(createCircleGeometry(60))
	//object.setScale(0.1)
	//object.configure(scene.program)
	//scene.addObject(object)

    world := new(World)
    world.Scene = scene
    world.window = window
    world.program = program
    //&renderer.Start(world)

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
