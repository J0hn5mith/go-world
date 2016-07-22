package go_world_physics

import (
	"github.com/go-gl/mathgl/mgl32"
	"go-world/go-world"
)

type RigidBody struct {
	object          *go_world.Object
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

func CreateRigidBody(object *go_world.Object) *RigidBody {
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

//TODO: Deprecated
func (rigidBody *RigidBody) GetMassParticles() []*MassParticle {
	return rigidBody.massParticles
}

func (rigidBody *RigidBody) MassParticles() []*MassParticle {
	return rigidBody.massParticles
}

func (rigidBody *RigidBody) ApplyForce(x, y, z float32) PhysicalBody {
	for _, particle := range rigidBody.MassParticles() {
		particle.ApplyForce(
			x,
			y,
			z,
		)
	}
	return rigidBody
}

func (rigidBody *RigidBody) Object() *go_world.Object {
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
