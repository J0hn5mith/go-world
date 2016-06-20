package go_world

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

type Renderer struct {
	camera *Camera
}

func NewRenderer(camera *Camera) *Renderer {
	renderer := new(Renderer)
	renderer.camera = camera

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	return renderer
}

func (r *Renderer) Render(world *World) {
	gl.ClearDepth(1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	r.render(world.Scene)
    r.debugRenderSoftBodies(world.Physics().softBodies)//TODO: check if present!!

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
		object.geometry.draw_method,
		0,
		int32(len(object.geometry.vertices)),
	)
}

func (r *Renderer) debugRenderSoftBodies(softBodies []*SoftBody) {
	particleGeometry := createCircleGeometry(100, 0.05).Load(r.camera.program)
	for _, softBody := range softBodies {
		positionBody := softBody.Position()
        translationBody := mgl32.Translate3D(
            positionBody.X(),
            positionBody.Y(),
            positionBody.Z(),
        )
		for _, particle := range softBody.GetMassParticles() {
			position := particle.Position()
			trans := mgl32.Translate3D(
				position.X(),
				position.Y(),
				position.Z(),
			)
	        mat := mgl32.Ident4()
			mat = mat.Mul4(translationBody).Mul4(trans)
			modelUniform := gl.GetUniformLocation(
				r.camera.program,
				gl.Str("model\x00"),
			)
			gl.UniformMatrix4fv(modelUniform, 1, false, &mat[0])

			gl.BindVertexArray(particleGeometry.vao)
			gl.DrawArrays(
				particleGeometry.DrawMethod(),
				0,
				int32(len(particleGeometry.Vertices())),
			)
		}
	}
}

/*TODO: Remove*/
func createCircleGeometry(num_vertices int, radius float32) *Geometry {
	vertices := []float32{}

	px, py := angleToCoords(0)

	for i := 1; i < int(num_vertices); i++ {
		var angle float64 = (2 * math.Pi * float64(i) / float64(num_vertices))
		x, y := angleToCoords(angle)
		vertices = append(
			vertices,
			0, 0, 0,
			float32(x)*radius, float32(y)*radius, 0.0,
			float32(px)*radius, float32(py)*radius, 0.0,
		)
		px = x
		py = y
	}

	x, y := angleToCoords(0)
	vertices = append(
		vertices,
		0, 0, 0,
		float32(x)*radius, float32(y)*radius, 0.0,
		float32(px)*radius, float32(py)*radius, 0.0,
	)

	geometry := NewGeometry(vertices)
	geometry.SetDrawMethod(gl.TRIANGLES)
	return geometry
}

func angleToCoords(angle float64) (float64, float64) {
	x := math.Sin(angle)
	y := math.Cos(angle)
	return x, y
}
