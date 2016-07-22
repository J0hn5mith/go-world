package go_world_physics

import (
	"github.com/go-gl/mathgl/mgl32"
	"go-world/go-world"
)

/*
Extends a plain object with a physical representation
*/
type PhysicalBody interface {
	SetVelocity(x, y, z float32) PhysicalBody
	GetVelocity() mgl32.Vec3 //TODO: Deprecated
	Velocity() mgl32.Vec3

	SetPosition(x, y, z float32) PhysicalBody
	Position() mgl32.Vec3
	ShiftPosition(x, y, z float32) PhysicalBody

	SetAngularVelocity(x, y, z float32) PhysicalBody
	AngularVelocity() mgl32.Vec3

	AddMassParticle(particle *MassParticle) PhysicalBody
	GetMassParticles() []*MassParticle

	ApplyForce(x, y, z float32) PhysicalBody

	Object() *go_world.Object

	Mass() float32
	SetMass(mass float32) PhysicalBody
}

type Physics struct {
	bodies    []*RigidBody
	staticBodies     []*RigidBody
	softBodies       []*SoftBody
	forceFields      []ForceField
	collisionHandler PhysicsCollisionHandler
	airResistance    float32
}

func NewPhysics() *Physics {
	physics := new(Physics)
	physics.airResistance = 0.95
	return physics
}

func (physics *Physics) RegisterDynamicObject(object *go_world.Object) *RigidBody {
	rigidBody := CreateRigidBody(object)
	physics.bodies = append(physics.bodies, rigidBody)
	return rigidBody
}

func (physics *Physics) RigidBodies() []*RigidBody {
	return append(physics.staticBodies, physics.bodies...)
}

func (physics *Physics) SoftBodies() []*SoftBody {
	return physics.softBodies
}

func (physics *Physics) RegisterSoftBody(body *SoftBody) PhysicalBody {
	physics.softBodies = append(physics.softBodies, body)
	return body
}

func (physics *Physics) RegisterStaticObject(object *go_world.Object) *RigidBody {
	rigidBody := CreateRigidBody(object)
	physics.staticBodies = append(physics.staticBodies, rigidBody)
	return rigidBody
}

func (physics *Physics) Update(timeDelta float32) {
	physics.updateVelocity(timeDelta)

	//physics.applySpringForces(timeDelta)
    physics.applyForceFields(timeDelta)
	//physics.applyAirResistance(timeDelta)

    physics.updatePosition(timeDelta)

    if physics.collisionHandler != nil {
        physics.collisionHandler.Apply(
            physics.bodies,
            physics.staticBodies,
            physics.softBodies,
        )
    }

}

func (physics *Physics) AddForceField(forceField ForceField) {
	physics.forceFields = append(physics.forceFields, forceField)
}

func (physics *Physics) AddCollisionHandler(collisionHandler PhysicsCollisionHandler) {
	physics.collisionHandler = collisionHandler
}

var ALPHA float32 = 1.0

func (physics *Physics) updateVelocity(timeDelta float32) {
	for _, body := range physics.bodies {
        cm := getCenterOfMass(body)
        new_cm := getNewCenterOfMass(body, timeDelta)
        rotation := getRotationMatrix(body, cm, new_cm,  timeDelta)
		for _, particle := range body.MassParticles() {
            goalPosition := particle.Position().Sub(cm)
            goalPosition  =  mgl32.TransformCoordinate(goalPosition, rotation.Mat4())
            goalPosition = goalPosition.Add(new_cm)
            positionDelta := goalPosition.Sub(
                particle.Position().Add(particle.Velocity().Mul(timeDelta)),
            )

            v_delta := positionDelta.Mul(ALPHA/timeDelta)
			v_new := particle.Velocity().Add(v_delta)
			particle.SetVelocity(
				v_new.X(),
				v_new.Y(),
				v_new.Z(),
			)
		}
	}
}

/*
    Returns the new center of mass for a body. The new center is
    based on the bodies particle and their current velocity.
*/
func getNewCenterOfMass(body *RigidBody, time_delta float32) mgl32.Vec3 {
    new_center := mgl32.Vec3{0,0,0}
    for _, particle := range body.MassParticles() {
        new_position := particle.Position().Add(particle.Velocity().Mul(time_delta))
        new_center = new_center.Add(new_position)
    }
    return new_center.Mul(1.0/float32(len(body.MassParticles())))
}

func getCenterOfMass(body *RigidBody) mgl32.Vec3 {
    new_center := mgl32.Vec3{0,0,0}
    for _, particle := range body.MassParticles() {
        new_center = new_center.Add(particle.Position())
    }
    return new_center.Mul(1.0/float32(len(body.MassParticles())))
}

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
        physics.updatePositionBody(timeDelta, body)
	}
}

func (physics *Physics) updatePositionBody(timeDelta float32, body PhysicalBody) {
	for _, particle := range body.GetMassParticles() {
		v := particle.Velocity()
		particle.ShiftPosition(
			v.X()*timeDelta,
			v.Y()*timeDelta,
			v.Z()*timeDelta,
		)
	}
}

func (physics *Physics) applyForceFields(timeDelta float32) {
	for _, forceField := range physics.forceFields {
		for _, rigidBody := range physics.bodies {
			forceField.Apply(rigidBody, timeDelta)
		}
	}
}

//TODO: Could this iteratoin be done using functional programming? Since I use
//it twice
func (physics *Physics) applyAirResistance(timeDelta float32) {
	for _, forceField := range physics.forceFields {
		for _, rigidBody := range physics.bodies {
			rigidBody.velocity = rigidBody.velocity.Mul(physics.airResistance)
		}
		for _, softBody := range physics.softBodies {
			forceField.ApplySoft(softBody, timeDelta)
			softBody.velocity = softBody.velocity.Mul(physics.airResistance)
		}
	}
}

func (physics *Physics) applySpringForces(timeDelta float32) {
	for _, softBody := range physics.softBodies {
		softBody.UpdateSpringForces()
	}
	for _, softBody := range physics.softBodies {
		softBody.ApplySpringForces(timeDelta)
	}
}

type ForceField interface {
	Apply(p PhysicalBody, timeDetla float32)
	ApplySoft(softBody *SoftBody, timeDetla float32)
}

type PhysicsCollisionHandler interface {
	Apply(bodies, staticBodies []*RigidBody, softBodies []*SoftBody)
}
