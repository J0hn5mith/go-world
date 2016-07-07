package go_world

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
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

	Object() *Object

	Mass() float32
	SetMass(mass float32) PhysicalBody
}

type RigidBody struct {
	object          *Object
	velocity        mgl32.Vec3
	angularVelocity mgl32.Vec3
	centerOfMass    mgl32.Vec3
	massParticles   []*MassParticle
	mass            float32
}

func NewRigidBody() *RigidBody {
	rigidBody := new(RigidBody)
	rigidBody.angularVelocity = mgl32.Vec3{0, 0, 0}
	return rigidBody
}

func CreateRigidBody(object *Object) *RigidBody {
	rigidBody := new(RigidBody)
	rigidBody.object = object
	rigidBody.velocity = mgl32.Vec3{0, 0, 0}
	rigidBody.mass = 1.0

	return rigidBody
}

func (rigidBody *RigidBody) SetVelocity(x, y, z float32) PhysicalBody {
	rigidBody.velocity = mgl32.Vec3{x, y, z}
	return rigidBody
}

func (rigidBody *RigidBody) GetVelocity() mgl32.Vec3 {
	return rigidBody.velocity
}

func (rigidBody *RigidBody) Velocity() mgl32.Vec3 {
	return rigidBody.velocity
}

func (rigidBody *RigidBody) SetAngularVelocity(x, y, z float32) PhysicalBody {
	rigidBody.angularVelocity = mgl32.Vec3{x, y, z}
	return rigidBody
}

func (rigidBody *RigidBody) AngularVelocity() mgl32.Vec3 {
	return rigidBody.angularVelocity
}

func (rigidBody *RigidBody) AddMassParticle(massParticle *MassParticle) PhysicalBody {
	rigidBody.massParticles = append(rigidBody.massParticles, massParticle)
	return rigidBody
}

func (rigidBody *RigidBody) GetMassParticles() []*MassParticle {
	return rigidBody.massParticles
}

func (rigidBody *RigidBody) ApplyForce(x, y, z float32) PhysicalBody {
	v := rigidBody.GetVelocity()
	rigidBody.SetVelocity(
		v.X()+x,
		v.Y()+y,
		v.Z()+z,
	)
	return rigidBody
}

func (rigidBody *RigidBody) Object() *Object {
	return rigidBody.object
}

func (rigidBody *RigidBody) Mass() float32 {
	return rigidBody.mass
}
func (rigidBody *RigidBody) SetMass(mass float32) PhysicalBody {
	rigidBody.mass = mass
	return rigidBody
}

func (rigidBody *RigidBody) SetPosition(x, y, z float32) PhysicalBody {
	rigidBody.object.SetPosition(x, y, z)
	return rigidBody
}

func (rigidBody *RigidBody) Position() mgl32.Vec3 {
	return rigidBody.object.Position()
}

func (rigidBody *RigidBody) ShiftPosition(x, y, z float32) PhysicalBody {
	p := rigidBody.Position()
	rigidBody.object.SetPosition(
		p[0]+x,
		p[1]+y,
		p[2]+z,
	)
	return rigidBody
}

type Physics struct {
	dynamicBodies    []*RigidBody
	staticBodies     []*RigidBody
	softBodies       []*SoftBody
	forceFields      []ForceField
	collisionHandler PhysicsCollisionHandler
}

func NewPhysics() *Physics {
	physics := new(Physics)
	return physics
}

func (physics *Physics) RegisterDynamicObject(object *Object) *RigidBody {
	rigidBody := CreateRigidBody(object)
	physics.dynamicBodies = append(physics.dynamicBodies, rigidBody)
	return rigidBody
}

func (physics *Physics) SoftBodies() []*SoftBody {
    return physics.softBodies
}

func (physics *Physics) RegisterSoftBody(body *SoftBody) PhysicalBody {
	physics.softBodies = append(physics.softBodies, body)
	return body
}

func (physics *Physics) RegisterStaticObject(object *Object) *RigidBody {
	rigidBody := CreateRigidBody(object)
	physics.staticBodies = append(physics.staticBodies, rigidBody)
	return rigidBody
}

func (physics *Physics) Update(timeDelta float32) {
    physics.applySpringForces(timeDelta)
	physics.applyForceFields(timeDelta)
	physics.animate(timeDelta)
	if physics.collisionHandler != nil {
		physics.collisionHandler.Apply(
			physics.dynamicBodies,
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

func (physics *Physics) animate(time_delta float32) {
	for _, rigidBody := range physics.dynamicBodies {
		rigidBody.object.position[0] += (float32)(rigidBody.velocity[0] * time_delta)
		rigidBody.object.position[1] += (float32)(rigidBody.velocity[1] * time_delta)
		rigidBody.object.position[2] += (float32)(rigidBody.velocity[2] * time_delta)
		rigidBody.object.RotateX(rigidBody.angularVelocity[0])
		rigidBody.object.RotateY(rigidBody.angularVelocity[1])
		rigidBody.object.RotateZ(rigidBody.angularVelocity[2])
	}

	for _, softBody := range physics.softBodies {
		for _, particle := range softBody.GetMassParticles() {
			newPosition := particle.Velocity().Mul(time_delta).Add(particle.Position())
			particle.SetPosition(
				newPosition.X(),
				newPosition.Y(),
				newPosition.Z(),
			)
		}
        softBody.Object().Geometry().SetDrawMethod(gl.TRIANGLE_STRIP)
		//softBody.Object().Geometry().SetDrawMethod(gl.LINES)
		softBody.Object().Geometry().UpdateVertices(softBody.GetVertices())
	}
}

func (physics *Physics) applyForceFields(timeDelta float32) {
	for _, forceField := range physics.forceFields {
		for _, rigidBody := range physics.dynamicBodies {
			forceField.Apply(rigidBody, timeDelta)
		}
		for _, softBody := range physics.softBodies {
			forceField.ApplySoft(softBody, timeDelta)
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
	Apply(dynamicBodies, staticBodies []*RigidBody, softBodies []*SoftBody)
}


