package go_world_utils

import (
    mgl "github.com/go-gl/mathgl/mgl64"
    "go-world/physics"
)

func CreateBoxBody(physics *go_world_physics.Physics, boxDimension mgl.Vec3) *go_world_physics.RigidBody{
        body := go_world_physics.CreateDynamicBody()
	    go_world_physics.AddMassParticle(body, boxDimension, 1)
	    physics.RegisterBody(body)
        return body
}

func CreateBoxBodyTower(physics *go_world_physics.Physics, boxDimension mgl.Vec3, height int, origin mgl.Vec3) []*go_world_physics.RigidBody{
    var bodies []*go_world_physics.RigidBody

    for i := 1; i <= height; i++ {
        body := go_world_physics.CreateDynamicBody()
	    go_world_physics.AddMassParticle(body, boxDimension, 2)
        body.SetPosition(origin.Add(mgl.Vec3{
            0,
            (float64(i) + 0.2) * boxDimension.Y() -0.1,
            0,
        }))
	    physics.RegisterBody(body)
        bodies = append(bodies, body)
    }

    return bodies
}
