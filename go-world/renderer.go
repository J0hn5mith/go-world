package go_world

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Renderer struct {
	camera         *Camera
	debugRenderers []DebugRenderer
}

func NewRenderer(camera *Camera) *Renderer {
	renderer := new(Renderer)
	renderer.camera = camera

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	return renderer
}

func (r *Renderer) Camera() *Camera{
    return r.camera
}

func (r *Renderer) AddDebugRenderer(debugRenderer DebugRenderer) *Renderer{
    r.debugRenderers = append(r.debugRenderers, debugRenderer)
    return r
}

func (r *Renderer) Render(world *World) {
	gl.ClearDepth(1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	r.render(world.Scene)
	for _, debugRenderer := range r.debugRenderers {
		debugRenderer.Render(world)
	}

	world.window.SwapBuffers()
	glfw.PollEvents()
}

func (renderer *Renderer) render(scene *Scene) {
	gl.UseProgram(renderer.camera.program)
	for _, object := range scene.objects {
		uniColor := gl.GetUniformLocation(
			renderer.camera.program,
			gl.Str("modelColor\x00"),
		)
        r, g, b := Vec3To32(object.Geometry().Color()).Elem()
		gl.Uniform3f(uniColor, r, g, b)
		//renderer.renderObject(object)
	}
}

func (r *Renderer) renderObject(object *Object) {
	mat := mgl32.Ident4()
	trans := mgl32.Translate3D( Vec3To32(object.Position()).Elem())
	scale := mgl32.Scale3D( Vec3To32(object.Scale()).Elem())
	mat = mat.Mul4(scale).Mul4(trans).Mul4(Mat4To32(object.Rotation()))
	modelUniform := gl.GetUniformLocation(
		r.camera.program,
		gl.Str("model\x00"),
	)
	gl.UniformMatrix4fv(modelUniform, 1, false, &mat[0])

	gl.BindVertexArray(object.geometry.vao)
	gl.DrawArrays(
		object.geometry.drawMethod,
		0,
		int32(len(object.geometry.vertices)),
	)
}

type DebugRenderer interface {
	Renderer() *Renderer
	Render(world *World)
}
