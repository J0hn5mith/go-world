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

func (r *Renderer) render(scene *Scene) {
	gl.UseProgram(r.camera.program)
	for _, object := range scene.objects {
		uniColor := gl.GetUniformLocation(
			r.camera.program,
			gl.Str("modelColor\x00"),
		)
		gl.Uniform3f(uniColor,
			object.geometry.color[0],
			object.geometry.color[1],
			object.geometry.color[2],
		)
		r.renderObject(object)
	}
}

func (r *Renderer) renderObject(object *Object) {
	mat := mgl32.Ident4()
	trans := mgl32.Translate3D(
		object.position[0],
		object.position[1],
		object.position[2],
	)
	scale := mgl32.Scale3D(
		object.scale[0],
		object.scale[1],
		object.scale[2],
	)
	mat = mat.Mul4(scale).Mul4(trans).Mul4(object.rotation)
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

//func (r *Renderer) debugRenderSoftBodies(softBodies []*SoftBody) {
//particleGeometry := createCircleGeometry(100, 0.05).Load(r.camera.program)
//for _, softBody := range softBodies {
//positionBody := softBody.Position()
//translationBody := mgl32.Translate3D(
//positionBody.X(),
//positionBody.Y(),
//positionBody.Z(),
//)
//for _, particle := range softBody.GetMassParticles() {
//position := particle.Position()
//trans := mgl32.Translate3D(
//position.X(),
//position.Y(),
//position.Z(),
//)
//mat := mgl32.Ident4()
//mat = mat.Mul4(translationBody).Mul4(trans)
//modelUniform := gl.GetUniformLocation(
//r.camera.program,
//gl.Str("model\x00"),
//)
//gl.UniformMatrix4fv(modelUniform, 1, false, &mat[0])

//gl.BindVertexArray(particleGeometry.vao)
//gl.DrawArrays(
//particleGeometry.DrawMethod(),
//0,
//int32(len(particleGeometry.Vertices())),
//)
//}
//}
//}
