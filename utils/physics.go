package go_world_utils

import (
	mgl "github.com/go-gl/mathgl/mgl64"
	physics "go-world/physics"
    "math"
)

var G float64 = 9.81
var K float64 = 155
var B float64 = 1.141
var FRICTION float64 = 0.6

/*
Default implementations for the physics stuff
*/

type GravityForceField struct{}

func (forceField *GravityForceField) Apply(body physics.PhysicalBody, timeDelta float64) {
    body.ApplyForce(mgl.Vec3{0, -G*timeDelta*2.0, 0.0})
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
            if len(physics.InterSphereCollisions(
                bodyA.BoundingSpheres(), bodyB.BoundingSpheres(),
            )) < 1 {
                continue
            }
			for _, particleA := range bodyA.MassParticles() {
				for _, particleB := range bodyB.MassParticles() {
					col := physics.DetectInterParticleCollision(particleA, particleB)
					if col.Magnitude > 0 {
                        mag := math.Pow(1 + col.Magnitude, 2) - 1
                        if !bodyA.Static() {
                            springForce := col.Direction.Mul( -K * -mag  - B * particleA.Velocity().Dot(col.Direction))
                            particleA.ApplyForce(springForce)
                            friction := particleB.Velocity().Mul(-FRICTION)
                            particleB.ApplyForce(friction)
                        }
                        if !bodyB.Static() {
                            springForce := col.Direction.Mul( -K * mag  - B * particleB.Velocity().Dot(col.Direction))
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
	posA := mgl.TransformCoordinate(
		particleA.Position(),
		bodyA.Object().TransformationMatrix(),
	)
	posB := mgl.TransformCoordinate(
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

func calculateNewVelocity(body1, body2 physics.PhysicalBody, collision physics.Collision) (mgl.Vec3, mgl.Vec3) {
	k := (body1.Position().
		Sub(body2.Position())).
		Mul(1 / (body1.Position().Sub(body2.Position()).Len()))
	a := k.Mul(2).Dot(body1.Velocity().Sub(body2.Velocity())) * (1 / (1/body1.Mass() + 1/body2.Mass()))
	p1_vel := body1.Velocity().Sub(k.Mul(a / body1.Mass()))
	p2_vel := body2.Velocity().Add(k.Mul(a / body2.Mass()))

	return p1_vel, p2_vel
}
