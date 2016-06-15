package go_world

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
    "fmt"
)

type Camera struct {
	viewMatrix mgl32.Mat4
	program    uint32
	ratio      float32
	position   mgl32.Vec3
}

func NewCamera(program uint32, windowWidth int, windowHeight int) *Camera {
	camera := new(Camera)
	camera.program = program
	camera.ratio = float32(windowWidth) / float32(windowHeight)
	camera.position = mgl32.Vec3{0, 0, 0}
	camera.updateViewMatrix()

	return camera
}

func (camera *Camera) SetPosition(x, y, z float32) *Camera {
	camera.position = mgl32.Vec3{x, y, z}
	camera.updateViewMatrix()
	return camera
}

func (camera *Camera) updateViewMatrix() {
	// projection matrix
	projectionMatrix := mgl32.Perspective(
		mgl32.DegToRad(45.0),
		camera.ratio,
		0.1,
		10.0,
	)

	//TODO what doeas that do
	projectionUniform := gl.GetUniformLocation(
		camera.program,
		gl.Str("projection\x00"),
	)
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projectionMatrix[0])

    fmt.Println(camera.position)
    viewMatrix := mgl32.LookAtV(
        camera.position,
        mgl32.Vec3{0, 0, 0},
        mgl32.Vec3{0, 1, 0},
    )

    cameraUniform := gl.GetUniformLocation(
        camera.program,
        gl.Str("camera\x00"),
    )
    gl.UniformMatrix4fv(cameraUniform, 1, false, &viewMatrix[0])

    camera.viewMatrix = viewMatrix
}
