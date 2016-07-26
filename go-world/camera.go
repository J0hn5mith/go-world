package go_world

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
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

func (camera *Camera) Program() uint32 {
    return camera.program
}

func (camera *Camera) updateViewMatrix() {
	projectionMatrix := mgl32.Perspective(
		mgl32.DegToRad(45.0),
		camera.ratio,
		0.1,
		100.0,
	)
	projectionUniform := gl.GetUniformLocation(
		camera.program,
		gl.Str("projection\x00"),
	)
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projectionMatrix[0])
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
