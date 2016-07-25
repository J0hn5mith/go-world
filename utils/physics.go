package go_world_utils

import (
	"github.com/go-gl/mathgl/mgl32"
	physics "go-world/physics"
)

var G float32 = 9.81
var K float32 = 500
var B float32 = 1.641
var FRICTION float32 = 0.00

/*
Default implementations for the physics stuff
*/

type GravityForceField struct{}

func (forceField *GravityForceField) Apply(body physics.PhysicalBody, time_delta float32) {
    body.ApplyForce( mgl32.Vec3{0, -G*time_delta, 0,})
}

func CreateGravityForceField() *GravityForceField {
	return new(GravityForceField)
}


type BasicPhysicsCollisionHandler struct{}

func CreateBasicPhysicsCollisionHandler() *BasicPhysicsCollisionHandler {
	return new(BasicPhysicsCollisionHandler)
}

/*
    Handles the collision among bodies. Spring forces are used to handle the
    effect of inter body collision.
*/
func (collisionHandler *BasicPhysicsCollisionHandler) Apply(bodies []*physics.RigidBody) {
	for _, bodyA := range bodies {
		for _, bodyB := range bodies {
            if bodyA == bodyB { continue }
			for _, particleA := range bodyA.MassParticles() {
				for _, particleB := range bodyB.MassParticles() {
					col := physics.DetectInterParticleCollision(particleA, particleB)
					if col.Magnitude > 0 {
                        if !bodyA.Static() {
                            springForce := col.Direction.Mul( -K * -col.Magnitude  - B * particleA.Velocity().Dot(col.Direction))
                            particleA.ApplyForce(springForce)
                            friction := particleB.Velocity().Mul(-FRICTION)
                            particleB.ApplyForce(friction)
                        }
                        if !bodyB.Static() {
                            springForce := col.Direction.Mul( -K * col.Magnitude  - B * particleB.Velocity().Dot(col.Direction))
                            particleB.ApplyForce(springForce)
                            friction := particleB.Velocity().Mul(-FRICTION)
                            particleB.ApplyForce(friction)
                        }
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
	return physics.CircleCollision(
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
