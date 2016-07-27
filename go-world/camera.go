package go_world

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	mgl "github.com/go-gl/mathgl/mgl64"
)

type Camera struct {
	viewMatrix mgl.Mat4
	program    uint32
	ratio      float64
	position   mgl.Vec3
}

func NewCamera(program uint32, windowWidth, windowHeight float64) *Camera {
	camera := new(Camera)
	camera.program = program
	camera.ratio = windowWidth/windowHeight
	camera.position = mgl.Vec3{0, 0, 0}
	camera.updateViewMatrix()

	return camera
}

func (camera *Camera) SetPosition(x, y, z float32) *Camera {
	camera.position = mgl.Vec3{
        float64(x),
        float64(y),
        float64(z),
    }
	camera.updateViewMatrix()
	return camera
}

func (camera *Camera) Program() uint32 {
    return camera.program
}

func (camera *Camera) updateViewMatrix() {
	projectionMatrix := Mat4To32(mgl.Perspective(
		mgl.DegToRad(45.0),
		camera.ratio,
		0.1,
		300.0,
	))
	projectionUniform := gl.GetUniformLocation(
		camera.program,
		gl.Str("projection\x00"),
	)
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projectionMatrix[0])
    camera.viewMatrix = mgl.LookAtV(
        camera.position,
        mgl.Vec3{0, 3, 0},
        mgl.Vec3{0, 1, 0},
    )

    viewMatrix := Mat4To32(camera.viewMatrix)
        cameraUniform := gl.GetUniformLocation(
        camera.program,
        gl.Str("camera\x00"),
    )
    gl.UniformMatrix4fv(cameraUniform, 1, false, &viewMatrix[0])
}
