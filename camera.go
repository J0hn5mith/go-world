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
    projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 10.0)
    projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
    gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

    viewMatrix := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
    cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
    gl.UniformMatrix4fv(cameraUniform, 1, false, &viewMatrix[0])

    camera.viewMatrix = viewMatrix
    camera.program = program
    return camera
}

// TODO: Change to scene
func (camera *Camera) render(scene *Scene) {

    gl.UseProgram(camera.program)
    for _, object := range scene.objects {
        camera.renderObject(object)
    }
}

func (camera *Camera) renderObject(object *Object) {

    // TODO: Make naming and logic consistent
    trans := mgl32.Ident4()
    rot := mgl32.Ident4()
    trans = mgl32.Translate3D(
        object.position[0], 
        object.position[1], 
        object.position[2],
    )
    scale := mgl32.Scale3D(
        object.scale[0],
        object.scale[1],
        object.scale[2],
    )

    rot = mgl32.HomogRotate3D(float32(object.angle), mgl32.Vec3{0, 1, 0})
    object.model = trans.Mul4(rot)
    object.model = object.model.Mul4(scale)

	// Render calls
    gl.UniformMatrix4fv(object.modelUniform, 1, false, &object.model[0])
    gl.BindVertexArray(object.vao)
    gl.DrawArrays(gl.TRIANGLES, 0, int32(len(object.geometry.vertices)))

}
