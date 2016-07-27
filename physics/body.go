package go_world_physics

import (
	mgl "github.com/go-gl/mathgl/mgl64"
	"go-world/go-world"
)

type RigidBody struct {
	object          *go_world.Object
	position        mgl.Vec3
	velocity        mgl.Vec3
	angularVelocity mgl.Vec3
	centerOfMass    mgl.Vec3
	massParticles   []*MassParticle
	mass            float64
	static          bool
	boundingSpheres []*Sphere
}

func NewRigidBody() *RigidBody {
	body := new(RigidBody)
	body.angularVelocity = mgl.Vec3{0, 0, 0}
	return body
}

func CreateDynamicBody() *RigidBody {
	return CreateBody(false)
}

func CreateStaticBody() *RigidBody {
	return CreateBody(true)
}

func CreateBody(static bool) *RigidBody {
	body := new(RigidBody)
	body.object = nil
	body.position = mgl.Vec3{0, 0, 0}
	body.velocity = mgl.Vec3{0, 0, 0}
	body.mass = 1.0
	body.static = static

	return body
}

func (body *RigidBody) SetVelocity(velocity mgl.Vec3) PhysicalBody {
	body.velocity = velocity
	for _, particle := range body.MassParticles() {
		particle.SetVelocity(velocity)
	}

	return body
}

func (body *RigidBody) GetVelocity() mgl.Vec3 {
	return body.velocity
}

func (body *RigidBody) Velocity() mgl.Vec3 {
	return body.velocity
}

func (body *RigidBody) SetAngularVelocity(velocity mgl.Vec3) PhysicalBody {
	body.angularVelocity = velocity
	return body
}

func (body *RigidBody) AngularVelocity() mgl.Vec3 {
	return body.angularVelocity
}

func (body *RigidBody) AddMassParticle(massParticle *MassParticle) PhysicalBody {
	body.massParticles = append(body.massParticles, massParticle)
	return body
}

func (body *RigidBody) MassParticles() []*MassParticle {
	return body.massParticles
}

func (body *RigidBody) ApplyForce(force mgl.Vec3) PhysicalBody {
	for _, particle := range body.MassParticles() {
		particle.ApplyForce(force)
	}
	return body
}

func (body *RigidBody) Object() *go_world.Object {
	return body.object
}

func (body *RigidBody) Mass() float64 {
	return body.mass
}
func (body *RigidBody) SetMass(mass float64) PhysicalBody {
	body.mass = mass
	return body
}

func (body *RigidBody) Static() bool {
	return body.static
}

func (body *RigidBody) SetStatic(static bool) *RigidBody {
	body.static = static
	return body
}

func (body *RigidBody) SetPosition(position mgl.Vec3) PhysicalBody {
	if body.object != nil {
		body.object.SetPosition(position)
	}
	shift := position.Sub(body.position)
	for _, particle := range body.MassParticles() {
		particle.ShiftPosition(shift)
	}
	body.ShiftBoundingSpheres(shift)
	body.position = position
	return body
}

func (body *RigidBody) Position() mgl.Vec3 {
	return body.position
}

func (body *RigidBody) ShiftPosition(shift mgl.Vec3) PhysicalBody {
	p := body.Position().Add(shift)
	body.object.SetPosition(p)
	return body
}

func (body *RigidBody) BoundingSpheres() []*Sphere {
	return body.boundingSpheres
}

func (body *RigidBody) AddBoundingSphere(boundingSphere *Sphere) *RigidBody {
	body.boundingSpheres = append(body.boundingSpheres, boundingSphere)
	return body
}

func (body *RigidBody) ShiftBoundingSpheres(shift mgl.Vec3) *RigidBody {
	for _, sphere := range body.BoundingSpheres() {
		sphere.Position = sphere.Position.Add(shift)
	}
	return body
}
