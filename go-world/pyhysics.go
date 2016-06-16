package go_world

import (
	"github.com/go-gl/mathgl/mgl32"
    "fmt"
)

/*
Extends a plain object with a physical representation
*/
type RigidBody struct {
	object          *Object
	velocity        mgl32.Vec3
	angularVelocity mgl32.Vec3
	centerOfMass    mgl32.Vec3
}

func NewRigidBody() *RigidBody {
	rigidBody := new(RigidBody)
	rigidBody.angularVelocity = mgl32.Vec3{0, 0, 0}
	return rigidBody
}

func (rigidBody *RigidBody) SetVelocity(x, y, z float32) *RigidBody {
	rigidBody.velocity = mgl32.Vec3{x, y, z}
	return rigidBody
}

func (rigidBody *RigidBody) SetAngularVelocity(x, y, z float32) *RigidBody {
	rigidBody.angularVelocity = mgl32.Vec3{x, y, z}
	return rigidBody
}

type Physics struct {
	rigidBodies []*RigidBody
	forceFields []ForceField
}

func NewPhysics() *Physics {
	physics := new(Physics)
	return physics
}

func (physics *Physics) RegisterObject(object *Object) *RigidBody {
	rigidBody := new(RigidBody)
	rigidBody.object = object
	rigidBody.velocity = mgl32.Vec3{0, 0, 0}

	physics.rigidBodies = append(physics.rigidBodies, rigidBody)
	return rigidBody
}

func (physics *Physics) Update(timeDelta float32) {
	physics.animate(timeDelta)
	physics.applyForces(timeDelta)
}

func (physics *Physics) AddForceField(forceField ForceField) {
	physics.forceFields = append(physics.forceFields, forceField)
}

func (physics *Physics) animate(time_delta float32) {
	for _, rigidBody := range physics.rigidBodies {
		rigidBody.object.position[0] += (float32)(rigidBody.velocity[0] * time_delta)
		rigidBody.object.position[1] += (float32)(rigidBody.velocity[1] * time_delta)
		rigidBody.object.position[2] += (float32)(rigidBody.velocity[2] * time_delta)
		rigidBody.object.RotateX(rigidBody.angularVelocity[0])
		rigidBody.object.RotateY(rigidBody.angularVelocity[1])
		rigidBody.object.RotateZ(rigidBody.angularVelocity[2])
	}
}

func (physics *Physics) applyForces(timeDelta float32) {
	for _, rigidBody := range physics.rigidBodies {
		for _, forceField := range physics.forceFields {
            fmt.Printf("", rigidBody, forceField);
            //forceField.Apply(rigidBody, timeDelta)
		}
	}
}
