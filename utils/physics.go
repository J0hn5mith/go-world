package go_world_utils

import (
    "github.com/go-gl/mathgl/mgl32"
    "go-world/go-world"
)
var G float32 = 9.81

type GravityForceField struct{}

func (forceField *GravityForceField) Apply(body go_world.PhysicalBody, timeDelta float32) {
    body.SetVelocity(
        body.GetVelocity()[0],
        body.GetVelocity()[1]-G*timeDelta,
        body.GetVelocity()[2],
    )
}

func (forceField *GravityForceField) ApplySoft(softBody *go_world.SoftBody, timeDelta float32) {
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

func (collisionHandler *BasicPhysicsCollisionHandler) Apply(dynamicBodies, staticBodies []*go_world.RigidBody, softBodies []*go_world.SoftBody) {
    for _, bodyA := range dynamicBodies {
        for _, bodyB := range staticBodies {
            foundCollision := false
            for _, particleA := range bodyA.GetMassParticles() {
                for _, particleB := range bodyB.GetMassParticles() {
                    col := detectCollision(bodyA, bodyB, particleA, particleB)
                    if col.Magnitude > 0 {
                        v := bodyA.Velocity()
                        shift := col.Direction.Mul(col.Magnitude)
                        bodyA.SetVelocity(v.X(), -v.Y(), v.Z())
                        bodyA.ShiftPosition(shift.X(), shift.Y(), shift.Z())

                        av := bodyA.AngularVelocity()
                        newPos := particleA.Position().Sub(shift)
                        _, _, phi := mgl32.CartesianToSpherical(particleA.Position())
                        _, _, phiN := mgl32.CartesianToSpherical(newPos)
                        bodyA.SetAngularVelocity(av.X(), av.Y(), phi-phiN)

                        foundCollision = true
                        break
                    }
                    if foundCollision {
                        break
                    }
                }
                if foundCollision {
                    foundCollision = false
                    break
                }
            }
        }
        for _, bodyB := range dynamicBodies {
            for _, particleA := range bodyA.GetMassParticles() {
                if bodyA != bodyB {
                    for _, particleB := range bodyB.GetMassParticles() {
                        col := detectCollision(bodyA, bodyB, particleA, particleB)
                        if col.Magnitude > 0 {
                            shiftA := col.Direction.Mul(col.Magnitude * 0.5)
                            shiftB := col.Direction.Mul(col.Magnitude * -0.5)
                            velA, velB := calculateNewVelocity(bodyA, bodyB, col)
                            bodyA.SetVelocity(
                                velA.X(),
                                velA.Y(),
                                velA.Z(),
                            )
                            bodyB.SetVelocity(
                                velB.X(),
                                velB.Y(),
                                velB.Z(),
                            )
                            bodyA.ShiftPosition(shiftA.X(), 4*shiftA.Y(), shiftA.Z())
                            bodyB.ShiftPosition(shiftB.X(), 4*shiftB.Y(), shiftB.Z())
                            av := bodyA.AngularVelocity()
                            newPosA := particleA.Position().Add(shiftA)
                            _, _, phi := mgl32.CartesianToSpherical(particleA.Position())
                            _, _, phiN := mgl32.CartesianToSpherical(newPosA)
                            bodyA.SetAngularVelocity(av.X(), av.Y(), phi-phiN)

                            newPosB := particleB.Position().Add(shiftB)
                            _, _, phi = mgl32.CartesianToSpherical(particleA.Position())
                            _, _, phiN = mgl32.CartesianToSpherical(newPosB)
                            bodyB.SetAngularVelocity(av.X(), av.Y(), phi-phiN)
                        }
                    }
                }
            }
        }
    }
    for _, bodyA := range softBodies {
        for _, bodyB := range staticBodies {
            for _, particleA := range bodyA.GetMassParticles() {
                for _, particleB := range bodyB.GetMassParticles() {
                    col := detectCollision(bodyA, bodyB, particleA, particleB)
                    if col.Magnitude > 0 {
                        v := particleA.Velocity()
                        shift := col.Direction.Mul(col.Magnitude)
                        particleA.ShiftPosition(shift.X(), shift.Y(), shift.Z())
                        particleA.SetVelocity(v.X(), -v.Y(), v.Z())
                        particleA.SetCollided(true)
                    }
                }
            }
        }
    }
}

func detectCollision(bodyA, bodyB go_world.PhysicalBody, particleA, particleB *go_world.MassParticle) go_world.Collision {
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

func calculateNewVelocity(body1, body2 go_world.PhysicalBody, collision go_world.Collision) (mgl32.Vec3, mgl32.Vec3) {
    k := (body1.Position().
    Sub(body2.Position())).
    Mul(1 / (body1.Position().Sub(body2.Position()).Len()))
    a := k.Mul(2).Dot(body1.Velocity().Sub(body2.Velocity())) * (1 / (1/body1.Mass() + 1/body2.Mass()))
    p1_vel := body1.Velocity().Sub(k.Mul(a / body1.Mass()))
    p2_vel := body2.Velocity().Add(k.Mul(a / body2.Mass()))
    return p1_vel, p2_vel
}
