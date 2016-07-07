package go_world_utils

import (
    "github.com/go-gl/mathgl/mgl32"
    "github.com/go-gl/gl/v4.1-core/gl"
    "go-world/go-world"
)

type MassParticleDebugRenderer struct{
    renderer *go_world.Renderer
}

func CreateMassParticleDebugRenderer(renderer * go_world.Renderer) *MassParticleDebugRenderer{
    debugRenderer :=  new(MassParticleDebugRenderer)
    debugRenderer.renderer = renderer
    return debugRenderer
}

func (debugRenderer * MassParticleDebugRenderer) Renderer() *go_world.Renderer{
    return debugRenderer.renderer
}

func (debugRenderer *MassParticleDebugRenderer ) Render(world *go_world.World) {
    softBodies := world.Physics().SoftBodies()
    renderer := debugRenderer.Renderer()
    particleGeometry := createCircleGeometry(100, 0.05).Load(renderer.Camera().Program())
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
    }
}
