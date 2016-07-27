package go_world_physics

import (
	"github.com/go-gl/mathgl/mgl32"
	"go-world/go-world"
)

/*
Extends a plain object with a physical representation
*/
type PhysicalBody interface {
	SetVelocity(velocity mgl32.Vec3) PhysicalBody
	GetVelocity() mgl32.Vec3 //TODO: Deprecated
	Velocity() mgl32.Vec3

	SetPosition(position mgl32.Vec3) PhysicalBody
	Position() mgl32.Vec3
	ShiftPosition(x, y, z float32) PhysicalBody

	SetAngularVelocity(x, y, z float32) PhysicalBody
	AngularVelocity() mgl32.Vec3

	AddMassParticle(particle *MassParticle) PhysicalBody
	MassParticles() []*MassParticle

	ApplyForce(force mgl32.Vec3) PhysicalBody

	Object() *go_world.Object

	Mass() float32
	SetMass(mass float32) PhysicalBody

    AddBoundingSphere(boundingSphere *Sphere) *RigidBody
    BoundingSpheres() []*Sphere
}

type Physics struct {
	bodies           []*RigidBody
	forceFields      []ForceField
	collisionHandler PhysicsCollisionHandler
	airResistance    float32
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

func (physics *Physics) Update(timeDelta float32) {

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

var ALPHA float32 = 1

//func (physics *Physics) getPreviousPositions(timeDelta float32) {
//var previousPosition []mgl32.Vec3
//for particle := range particles {
//preivousPositions = append(previousPosition, particle.Position())
//}
//}
func (physics *Physics) updateVelocity(timeDelta float32) {
	for _, body := range physics.bodies {
		if !body.Static() {
			cm := getCenterOfMass(body)
			new_cm := getNewCenterOfMass(body, timeDelta)
			rotation := getRotationMatrix(body, cm, new_cm, timeDelta)
			for _, particle := range body.MassParticles() {
				goalPosition := particle.Position().Sub(cm)
				goalPosition = mgl32.TransformCoordinate(goalPosition, rotation.Mat4())
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
func getNewCenterOfMass(body *RigidBody, time_delta float32) mgl32.Vec3 {
	new_center := mgl32.Vec3{0, 0, 0}
	for _, particle := range body.MassParticles() {
		new_position := particle.Position().Add(particle.Velocity().Mul(time_delta))
		new_center = new_center.Add(new_position)
	}
	return new_center.Mul(1.0 / float32(len(body.MassParticles())))
}

func getCenterOfMass(body *RigidBody) mgl32.Vec3 {
	new_center := mgl32.Vec3{0, 0, 0}
	for _, particle := range body.MassParticles() {
		new_center = new_center.Add(particle.Position())
	}
	return new_center.Mul(1.0 / float32(len(body.MassParticles())))
}

//func implicitEuerl(particles, preivousPositions){
//}

/*
   Returns the rotation matrix for the next state of a body.
   The rotation matrix is computed based on the position
   and velocity of the body's particle.
*/
func getRotationMatrix(body *RigidBody, old_center, new_center mgl32.Vec3, time_delta float32) mgl32.Mat3 {
	var old_positions []mgl32.Vec3
	var new_positions []mgl32.Vec3
	for _, particle := range body.MassParticles() {
		new_position := particle.Position().Add(
			particle.Velocity().Mul(time_delta),
		)
		old_positions = append(old_positions, particle.Position().Sub(old_center))
		new_positions = append(new_positions, new_position.Sub(new_center))
	}
	return ExtractRotationFromPositions(old_positions, new_positions)
}

func (physics *Physics) updatePosition(timeDelta float32) {
	for _, body := range physics.bodies {
		if !body.Static() {
			physics.updatePositionBody(timeDelta, body)
		}
	}
}

func (physics *Physics) updatePositionBody(timeDelta float32, body PhysicalBody) {
	for _, particle := range body.MassParticles() {
		delta := particle.Velocity().Mul(timeDelta)
		particle.ShiftPosition(delta)
	}
}

func (physics *Physics) applyForceFields(timeDelta float32) {
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
func (physics *Physics) applyAirResistance(timeDelta float32) {
	//for _, rigidBody := range physics.bodies {
        //for _, particle := range rigidBody.MassParticles(){
            //mag  := (1 - 1/((1 + particle.Velocity().Len())*(1 + particle.Velocity().Len())) 
			//particle.SetVelocity(particle.Velocity().Mul(mag * physics.airResistance))
        //}
	//}
}

func (physics *Physics) applySpringForces(timeDelta float32) {
    // Not used because I removed the soft body class
	//for _, body := range physics.bodies {
		//body.UpdateSpringForces()
	//}
}

type ForceField interface {
	Apply(p PhysicalBody, timeDetla float32)
}

type PhysicsCollisionHandler interface {
	Apply(bodies []*RigidBody)
}
