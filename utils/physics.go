package go_world_utils

import (
	"github.com/go-gl/mathgl/mgl32"
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


func (collisionHandler *BasicPhysicsCollisionHandler) Apply(bodies []*physics.RigidBody) {
	for _, bodyA := range bodies {
		for _, bodyB := range bodies {
            if bodyA == bodyB {
                continue
            }
			for _, particleA := range bodyA.GetMassParticles() {
				for _, particleB := range bodyB.GetMassParticles() {
					col := physics.DetectInterParticleCollision(particleA, particleB)
					if col.Magnitude > 0 {
                        springForce := col.Direction.Mul( -1 * col.Magnitude  - 10 * particleB.Velocity().Dot(col.Direction))
                        particleB.ApplyForce(
                            springForce.X(), springForce.Y(), springForce.Z(),
                        )
					}
				}
			}
		}
	}
}

func detectCollision(bodyA, bodyB physics.PhysicalBody, particleA, particleB *physics.MassParticle) physics.Collision {
	posA := mgl32.TransformCoordinate(
		particleA.Position(),
		bodyA.Object().TransformationMatrix(),
	)
	posB := mgl32.TransformCoordinate(
		particleB.Position(),
		bodyB.Object().TransformationMatrix(),
	)
	return physics.TestCicrcleCollision(
		posA,
		posB,
		particleA.Radius(),
		particleB.Radius(),
	)
}

func calculateNewVelocity(body1, body2 physics.PhysicalBody, collision physics.Collision) (mgl32.Vec3, mgl32.Vec3) {
	k := (body1.Position().
		Sub(body2.Position())).
		Mul(1 / (body1.Position().Sub(body2.Position()).Len()))
	a := k.Mul(2).Dot(body1.Velocity().Sub(body2.Velocity())) * (1 / (1/body1.Mass() + 1/body2.Mass()))
	p1_vel := body1.Velocity().Sub(k.Mul(a / body1.Mass()))
	p2_vel := body2.Velocity().Add(k.Mul(a / body2.Mass()))
	return p1_vel, p2_vel
}
