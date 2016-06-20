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
	physics.enforceConstraints()
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
			//rotation := softBody.AngularVelocity().Mul(time_delta)
			//softBody.object.RotateX(rotation.X())
			//softBody.object.RotateY(rotation.Y())
			//softBody.object.RotateZ(rotation.Z())
		}
		softBody.Object().Geometry().SetDrawMethod(gl.TRIANGLE_STRIP)
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

func (physics *Physics) enforceConstraints() {
	//for _, softBody := range physics.softBodies {
	//softBody.UpdateSpringForces()
	//}
}

type MassParticle struct {
	position    mgl32.Vec3 // Position relative to the center of mass
	velocity    mgl32.Vec3
	radius      float32
	collided    bool
	springs     []*Spring
	springForce mgl32.Vec3
}

func CreateMassParticle(x, y, z, radius float32) *MassParticle {
	massParticle := new(MassParticle)
	massParticle.position = mgl32.Vec3{x, y, z}
	massParticle.velocity = mgl32.Vec3{0, 0, 0}
	massParticle.radius = radius
	massParticle.collided = false

	return massParticle
}

func (massParticle *MassParticle) Position() mgl32.Vec3 {
	return massParticle.position
}

func (massParticle *MassParticle) SetPosition(x, y, z float32) *MassParticle {
	massParticle.position = mgl32.Vec3{x, y, z}
	return massParticle
}

func (massParticle *MassParticle) ShiftPosition(x, y, z float32) *MassParticle {
	p := massParticle.Position()
	massParticle.SetPosition(
		p[0]+x,
		p[1]+y,
		p[2]+z,
	)
	return massParticle
}

func (massParticle *MassParticle) Radius() float32 {
	return massParticle.radius
}

func (massParticle *MassParticle) Velocity() mgl32.Vec3 {
	return massParticle.velocity
}

func (massParticle *MassParticle) Collided() bool {
	return massParticle.collided
}

func (massParticle *MassParticle) SetCollided(collided bool) *MassParticle {
	massParticle.collided = collided
	return massParticle
}

func (massParticle *MassParticle) SetVelocity(x, y, z float32) *MassParticle {
	massParticle.velocity = mgl32.Vec3{x, y, z}
	return massParticle
}

func (massParticle *MassParticle) AddSpring(spring *Spring) *MassParticle {
	massParticle.springs = append(massParticle.springs, spring)
	return massParticle
}

func (massParticle *MassParticle) Springs() []*Spring {
	return massParticle.springs
}

func (massParticle *MassParticle) GetSpringForce() mgl32.Vec3 {
	return massParticle.springForce
}

func (massParticle *MassParticle) SetSpringForce(x, y, z float32) *MassParticle {
	massParticle.springForce = mgl32.Vec3{x, y, z}
	return massParticle
}

func (massParticle *MassParticle) ApplySpringForce(force mgl32.Vec3) *MassParticle {
	massParticle.springForce = massParticle.springForce.Add(force)
	return massParticle
}

func (massParticle *MassParticle) ApplyForce(x, y, z float32) *MassParticle {
	v := massParticle.Velocity()
	massParticle.SetVelocity(
		v.X()+x,
		v.Y()+y,
		v.Z()+z,
	)
	return massParticle
}

type Spring struct {
	source         *MassParticle
	target         *MassParticle
	length         float32
	springConstant float32
	damperConstant float32
}

func SpringFromMassParticles(particle1, particle2 *MassParticle, springConstant, damperConstant float32) *Spring {
	spring := new(Spring)
	spring.source = particle1
	spring.target = particle2
	spring.length = particle1.Position().Sub(particle2.Position()).Len()
	spring.springConstant = springConstant
	spring.damperConstant = damperConstant
	return spring
}

func (spring *Spring) Apply() {
	offset := spring.target.Position().Sub(spring.source.Position())
	direction := offset.Normalize()
    damperForce := direction.Mul(
        spring.target.Velocity().Dot(direction) - spring.source.Velocity().Dot(direction),
    ).Mul(spring.damperConstant)
	springForce := direction.Mul(offset.Len() - spring.length).Mul(spring.springConstant)

    totalForce := damperForce.Add(springForce)
    spring.source.ApplySpringForce(totalForce)
    spring.target.ApplySpringForce(totalForce.Mul(-1))

}

type ForceField interface {
	Apply(p PhysicalBody, timeDetla float32)
	ApplySoft(softBody *SoftBody, timeDetla float32)
}

type PhysicsCollisionHandler interface {
	Apply(dynamicBodies, staticBodies []*RigidBody, softBodies []*SoftBody)
}

/*
Creates mass particles for a box object
*/
func AddMassParticle2D(body PhysicalBody, x, y, diameter float32) {
	radius := diameter / 2
	lenX := int(x / diameter)
	lenY := int(y / diameter)
	offsetX := x / 2
	offsetY := y / 2

	// TODO: Consider rest of division
	for x := 0; x < lenX; x++ {
		for y := 0; y < lenY; y++ {
			xPos := -offsetX + radius + float32(x)*diameter
			yPos := -offsetY + radius + float32(y)*diameter
			body.AddMassParticle(
				CreateMassParticle(xPos, yPos, 0, radius),
			)
		}
	}
}
