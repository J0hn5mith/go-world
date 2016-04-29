package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
    "github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
    viewMatrix mgl32.Mat4
    program    uint32
}

func NewCamera(program uint32) *Camera {
    camera := new(Camera)
	
	//projection matrix
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 10.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])
	
	// camera position
	viewMatrix := mgl32.LookAtV(mgl32.Vec3{0, 0, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &viewMatrix[0])

	camera.viewMatrix = viewMatrix
    camera.program = program
    return camera
}


