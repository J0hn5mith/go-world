package go_world_utils

import (
	"go-world/go-world"
    "fmt"
)

type GravityForceField struct{}

func (forceField *GravityForceField) Apply(body go_world.PhysicalBody, timeDelta float32) {
	G := float32(9.81)
	body.SetVelocity(
		body.GetVelocity()[0],
		body.GetVelocity()[1]-G*timeDelta,
		body.GetVelocity()[2],
	)
}

func CreateGravityForceField() *GravityForceField {
	return new(GravityForceField)
}

type BasicPhysicsCollisionHandler struct{}

func CreateBasicPhysicsCollisionHandler() *BasicPhysicsCollisionHandler {
	return new(BasicPhysicsCollisionHandler)
}

func (collisionHandler *BasicPhysicsCollisionHandler) Apply(dynamicBodies, staticBodies []*go_world.RigidBody) {
    fmt.Println("Go!")
}
