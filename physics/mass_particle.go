package go_world_physics

import (
	"github.com/go-gl/mathgl/mgl32"
)

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
	op := body.Position()
	for x := 0; x < lenX; x++ {
		for y := 0; y < lenY; y++ {
			xPos := -offsetX + radius + float32(x)*diameter + op.X()
			yPos := -offsetY + radius + float32(y)*diameter + op.Y()
			body.AddMassParticle(
				CreateMassParticle(xPos, yPos, 0, radius),
			)
		}
	}
}
