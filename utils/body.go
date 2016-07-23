package go_world_utils

import (
    "github.com/go-gl/mathgl/mgl32"
    //"go-world/go-world"
    "go-world/physics"
)

func CreateBoxBodyTower(physics *go_world_physics.Physics, boxDimension mgl32.Vec3, height int, origin mgl32.Vec3){
    for i := 1; i <= height; i++ {
        body := go_world_physics.CreateDynamicBody()
	    go_world_physics.AddMassParticle2D(body, boxDimension.X(), boxDimension.Y(), 0.08)
        body.SetPosition(origin.Add(mgl32.Vec3{
            0,
            (float32(i) + 0.2) * boxDimension.Y() -0.35,
            0,
        }))
	    physics.RegisterBody(body)
    }
}
