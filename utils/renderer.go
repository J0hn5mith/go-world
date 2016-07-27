package go_world_utils

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"go-world/go-world"
	physics "go-world/physics"
)

type MassParticleDebugRenderer struct {
	renderer *go_world.Renderer
	physics  *physics.Physics
}

func CreateMassParticleDebugRenderer(renderer *go_world.Renderer, physics *physics.Physics) *MassParticleDebugRenderer {
	debugRenderer := new(MassParticleDebugRenderer)
	debugRenderer.renderer = renderer
	debugRenderer.physics = physics
	return debugRenderer
}

func (debugRenderer *MassParticleDebugRenderer) Renderer() *go_world.Renderer {
	return debugRenderer.renderer
}

func (debugRenderer *MassParticleDebugRenderer) Render(world *go_world.World) {
	bodies := debugRenderer.physics.Bodies()
	renderer := debugRenderer.Renderer()
	for _, body := range bodies {
		renderParticles(renderer, body)
		renderBoundingSpheres(renderer, body)
	}
}

func renderParticles(renderer *go_world.Renderer, body physics.PhysicalBody) {
	particleGeometry := createCircleGeometry(100, body.MassParticles()[0].Radius()).Load(renderer.Camera().Program())
	for _, particle := range body.MassParticles() {
		position := particle.Position()
		trans := mgl32.Translate3D(go_world.Vec3To32(position).Elem())
		mat := mgl32.Ident4()
		mat = mat.Mul4(trans)
		modelUniform := gl.GetUniformLocation(
			renderer.Camera().Program(),
			gl.Str("model\x00"),
		)
		gl.UniformMatrix4fv(modelUniform, 1, false, &mat[0])

		gl.BindVertexArray(particleGeometry.Vao())
		gl.DrawArrays(
			particleGeometry.DrawMethod(),
			0,
			int32(len(particleGeometry.Vertices())-1),
		)
	}
	particleGeometry.Delete()
}

func renderBoundingSpheres(renderer *go_world.Renderer, body physics.PhysicalBody) {
	for _, sphere := range body.BoundingSpheres() {
		position := sphere.Position
		particleGeometry := createCircleLineGeometry(100, sphere.Radius).Load(renderer.Camera().Program())
		trans := mgl32.Translate3D(go_world.Vec3To32(position).Elem())
		mat := mgl32.Ident4()
		mat = mat.Mul4(trans)
		modelUniform := gl.GetUniformLocation(
			renderer.Camera().Program(),
			gl.Str("model\x00"),
		)
		gl.UniformMatrix4fv(modelUniform, 1, false, &mat[0])

		gl.BindVertexArray(particleGeometry.Vao())
		gl.DrawArrays(
			particleGeometry.DrawMethod(),
			0,
			int32(len(particleGeometry.Vertices())-1),
		)
		particleGeometry.Delete()
	}
}
