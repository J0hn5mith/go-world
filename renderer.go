package go_world

import (
    "github.com/go-gl/gl/v4.1-core/gl"
    "github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/glfw/v3.1/glfw"
)

type Renderer struct {
    camera *Camera
}

func (r *Renderer) Render(world *World) {

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
        r.render(world.Scene)

		world.window.SwapBuffers()
		glfw.PollEvents()
}

func (r *Renderer) render(scene *Scene) {

    gl.UseProgram(r.camera.program)
    for _, object := range scene.objects {
        r.renderObject(object)
    }
}

func (r *Renderer) renderObject(object *Object) {

    // TODO: Make naming and logic consistent
    mat := mgl32.Ident4()
    rot := mgl32.Ident4()
    trans := mgl32.Translate3D(
        object.position[0],
        object.position[1],
        object.position[2],
    )
    scale := mgl32.Scale3D(
        object.scale[0],
        object.scale[1],
        //object.scale[2], // TODO
        0.0,
    )

    rot = mgl32.HomogRotate3D(float32(object.angle), mgl32.Vec3{0, 1, 0})
    mat = mat.Mul4(scale).Mul4(trans).Mul4(rot)

    // Render calls
    gl.UniformMatrix4fv(object.modelUniform, 1, false, &mat[0])
    gl.BindVertexArray(object.vao)
    gl.DrawArrays( 
        object.geometry.draw_method, 
        0,
        int32(len(object.geometry.vertices)),
    )

}
