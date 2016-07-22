package go_world_utils

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"go-world/go-world"
	physics "go-world/physics"
)

var G float32 = 9.81

/*
Default implementations for the physics stuff
*/

type GravityForceField struct{}

func (forceField *GravityForceField) Apply(body physics.PhysicalBody, time_delta float32) {
	body.ApplyForce(
		0,
		-G*time_delta,
		0,
	)
}

func (forceField *GravityForceField) ApplySoft(softBody *physics.SoftBody, timeDelta float32) {
	for _, particle := range softBody.GetMassParticles() {
		particle.SetVelocity(
			particle.Velocity()[0],
			particle.Velocity()[1]-G*timeDelta,
			particle.Velocity()[2],
		)
	}
}

func CreateGravityForceField() *GravityForceField {
	return new(GravityForceField)
}

type BasicPhysicsCollisionHandler struct{}

func CreateBasicPhysicsCollisionHandler() *BasicPhysicsCollisionHandler {
	return new(BasicPhysicsCollisionHandler)
}

func (collisionHandler *BasicPhysicsCollisionHandler) Apply(dynamicBodies, staticBodies []*physics.RigidBody, softBodies []*physics.SoftBody) {
	for _, bodyA := range dynamicBodies {
		for _, bodyB := range staticBodies {
			for _, particleA := range bodyA.GetMassParticles() {
				for _, particleB := range bodyB.GetMassParticles() {
					col := detectCollision(bodyA, bodyB, particleA, particleB)
					if col.Magnitude > 0 {
						fmt.Println("Collision")
					}
					//v := bodyA.Velocity()
					//shift := col.Direction.Mul(col.Magnitude)
					//bodyA.SetVelocity(v.X(), -v.Y(), v.Z())
					//bodyA.ShiftPosition(shift.X(), shift.Y(), shift.Z())

					//av := bodyA.AngularVelocity().Mul(0.5)
					//newPos := particleA.Position().Sub(shift)
					//_, _, phi := mgl32.CartesianToSpherical(particleA.Position())
					//_, _, phiN := mgl32.CartesianToSpherical(newPos)
					//bodyA.SetAngularVelocity(av.X(), av.Y(), phi-phiN)

					//foundCollision = true
					//break
					//}
					//if foundCollision {
					//break
					//}
					//}
					//if foundCollision {
					//foundCollision = false
					//break
					//}
				}
			}
		}
	}
}

func detectCollision(bodyA, bodyB physics.PhysicalBody, particleA, particleB *physics.MassParticle) go_world.Collision {
	posA := mgl32.TransformCoordinate(
		particleA.Position(),
		bodyA.Object().TransformationMatrix(),
	)
	posB := mgl32.TransformCoordinate(
		particleB.Position(),
		bodyB.Object().TransformationMatrix(),
	)
	return go_world.TestCicrcleCollision(
		posA,
		posB,
		particleA.Radius(),
		particleB.Radius(),
	)
}

func calculateNewVelocity(body1, body2 physics.PhysicalBody, collision go_world.Collision) (mgl32.Vec3, mgl32.Vec3) {
	k := (body1.Position().
		Sub(body2.Position())).
		Mul(1 / (body1.Position().Sub(body2.Position()).Len()))
	a := k.Mul(2).Dot(body1.Velocity().Sub(body2.Velocity())) * (1 / (1/body1.Mass() + 1/body2.Mass()))
	p1_vel := body1.Velocity().Sub(k.Mul(a / body1.Mass()))
	p2_vel := body2.Velocity().Add(k.Mul(a / body2.Mass()))
	return p1_vel, p2_vel
}
