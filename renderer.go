package main

import (
    "github.com/go-gl/gl/v4.1-core/gl"
    "github.com/go-gl/mathgl/mgl32"
)

type Renderer struct {
    camera *Camera
}

func (r *Renderer) render(scene *Scene) {

    gl.UseProgram(r.camera.program)
    for _, object := range scene.objects {
        r.renderObject(object)
    }
}

func (r *Renderer) renderObject(object *Object) {

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

    //rot = mgl32.HomogRotate3D(float32(object.angle), mgl32.Vec3{0, 1, 0})
    object.model = trans.Mul4(rot)
    object.model = object.model.Mul4(scale)

    // Render calls
    gl.UniformMatrix4fv(object.modelUniform, 1, false, &object.model[0])
    gl.BindVertexArray(object.vao)
    gl.DrawArrays( 
        object.geometry.draw_method, 
        //gl.LINES,
        0,
        int32(len(object.geometry.vertices)),
    )

}
