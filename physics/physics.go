package go_world_physics

import (
	mgl "github.com/go-gl/mathgl/mgl64"
	"go-world/go-world"
)

/*
Extends a plain object with a physical representation
*/
type PhysicalBody interface {
	SetVelocity(velocity mgl.Vec3) PhysicalBody
	GetVelocity() mgl.Vec3 //TODO: Deprecated
	Velocity() mgl.Vec3

	SetPosition(position mgl.Vec3) PhysicalBody
	Position() mgl.Vec3
	ShiftPosition(mgl.Vec3) PhysicalBody

	SetAngularVelocity(velocity mgl.Vec3) PhysicalBody
	AngularVelocity() mgl.Vec3

	AddMassParticle(particle *MassParticle) PhysicalBody
	MassParticles() []*MassParticle

	ApplyForce(force mgl.Vec3) PhysicalBody

	Object() *go_world.Object

	Mass() float64
	SetMass(mass float64) PhysicalBody

    AddBoundingSphere(boundingSphere *Sphere) *RigidBody
    BoundingSpheres() []*Sphere
}

type Physics struct {
	bodies           []*RigidBody
	forceFields      []ForceField
	collisionHandler PhysicsCollisionHandler
	airResistance    float64
}

func NewPhysics() *Physics {
	physics := new(Physics)
	physics.airResistance = 0.9
	return physics
}

func (physics *Physics) RegisterObject(object *go_world.Object) *RigidBody {
	body := CreateDynamicBody()
    body.object = object
	physics.bodies = append(physics.bodies, body)
	return body
}

func (physics *Physics) Bodies() []*RigidBody {
	return physics.bodies
}

func (physics *Physics) RegisterBody(body *RigidBody) *Physics {
	physics.bodies = append(physics.bodies, body)
    return physics
}

func (physics *Physics) Update(timeDelta float64) {

    physics.applySpringForces(timeDelta)
    physics.applyForceFields(timeDelta)
    physics.applyAirResistance(timeDelta)
    physics.updateVelocity(timeDelta)
	physics.updatePosition(timeDelta)

    if physics.collisionHandler != nil {
        physics.collisionHandler.Apply(physics.bodies)
    }

}

func (physics *Physics) AddForceField(forceField ForceField) {
	physics.forceFields = append(physics.forceFields, forceField)
}

func (physics *Physics) AddCollisionHandler(collisionHandler PhysicsCollisionHandler) {
	physics.collisionHandler = collisionHandler
}

var ALPHA float64 = 1

func (physics *Physics) updateVelocity(timeDelta float64) {
	for _, body := range physics.bodies {
		if !body.Static() {
			cm := getCenterOfMass(body)
            new_cm := getNewCenterOfMass(body, timeDelta)
            rotation := getRotationMatrix(body, cm, new_cm, timeDelta)
			for _, particle := range body.MassParticles() {
				goalPosition := particle.Position().Sub(cm)
                goalPosition = mgl.TransformCoordinate(goalPosition, rotation.Mat4())
				goalPosition = goalPosition.Add(new_cm)
				positionDelta := goalPosition.Sub(
					particle.Position().Add(particle.Velocity().Mul(timeDelta)),
				)
				v_delta := positionDelta.Mul(ALPHA / timeDelta)
				v_new := particle.Velocity().Add(v_delta)
				particle.SetVelocity(v_new)
			}

            // TODO: Shift to better place, makes no sense that's here
            shift := new_cm.Sub(cm)
            body.ShiftBoundingSpheres(shift)
		}
	}
}

/*
   Returns the new center of mass for a body. The new center is
   based on the bodies particle and their current velocity.
*/
func getNewCenterOfMass(body *RigidBody, time_delta float64) mgl.Vec3 {
	newCenter := mgl.Vec3{0, 0, 0}
	for _, particle := range body.MassParticles() {
		newPosition := particle.Position().Add(particle.Velocity().Mul(time_delta))
		newCenter = newCenter.Add(newPosition)
	}
	return newCenter.Mul(1.0 / float64(len(body.MassParticles())))
}

func getCenterOfMass(body *RigidBody) mgl.Vec3 {
	newCenter := mgl.Vec3{0, 0, 0}
	for _, particle := range body.MassParticles() {
		newCenter = newCenter.Add(particle.Position())
	}
	return newCenter.Mul(1.0 / float64(len(body.MassParticles())))
}

/*
   Returns the rotation matrix for the next state of a body.
   The rotation matrix is computed based on the position
   and velocity of the body's particle.
*/
func getRotationMatrix(body *RigidBody, oldCenter, newCenter mgl.Vec3, timeDelta float64) mgl.Mat3 {
	var oldPositions []mgl.Vec3
	var newPositions []mgl.Vec3
	for _, particle := range body.MassParticles() {
		newPosition := particle.Position().Add(
			particle.Velocity().Mul(timeDelta),
		)
		oldPositions = append(oldPositions, particle.Position().Sub(oldCenter))
		newPositions = append(newPositions, newPosition.Sub(newCenter))
	}
	return ExtractRotationFromPositions(oldPositions, newPositions)
}

func (physics *Physics) updatePosition(timeDelta float64) {
	for _, body := range physics.bodies {
		if !body.Static() {
			physics.updatePositionBody(timeDelta, body)
		}
	}
}

func (physics *Physics) updatePositionBody(timeDelta float64, body PhysicalBody) {
	for _, particle := range body.MassParticles() {
		delta := particle.Velocity().Mul(timeDelta)
		particle.ShiftPosition(delta)
	}
}

func (physics *Physics) applyForceFields(timeDelta float64) {
	for _, forceField := range physics.forceFields {
		for _, body := range physics.bodies {
			if !body.static {
				forceField.Apply(body, timeDelta)
			}
		}
	}
}

//TODO: Could this iteratoin be done using functional programming? Since I use
//it twice
func (physics *Physics) applyAirResistance(timeDelta float64) {
	//for _, rigidBody := range physics.bodies {
        //for _, particle := range rigidBody.MassParticles(){
            //mag  := (1 - 1/((1 + particle.Velocity().Len())*(1 + particle.Velocity().Len())) 
			//particle.SetVelocity(particle.Velocity().Mul(mag * physics.airResistance))
        //}
	//}
}

func (physics *Physics) applySpringForces(timeDelta float64) {
    // Not used because I removed the soft body class
	//for _, body := range physics.bodies {
		//body.UpdateSpringForces()
	//}
}

type ForceField interface {
	Apply(p PhysicalBody, timeDetla float64)
}

type PhysicsCollisionHandler interface {
	Apply(bodies []*RigidBody)
}
