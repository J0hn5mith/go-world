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

/*
Creates a new mass particles.
*/
func CreateMassParticle(x, y, z, radius float32) *MassParticle {
	massParticle := new(MassParticle)
	massParticle.position = mgl32.Vec3{x, y, z}
	massParticle.velocity = mgl32.Vec3{0, 0, 0}
	massParticle.radius = radius
	massParticle.collided = false

	return massParticle
}

/*
Returns the position of the mass particle.
*/
func (massParticle *MassParticle) Position() mgl32.Vec3 {
	return massParticle.position
}

/*
Sets the position of the mass particle.
*/
func (massParticle *MassParticle) SetPosition(position mgl32.Vec3) *MassParticle {
	massParticle.position = position
	return massParticle
}

/*
Shifts the position of the particle by provided values.
*/
func (massParticle *MassParticle) ShiftPosition(shift mgl32.Vec3) *MassParticle {
	massParticle.SetPosition(massParticle.Position().Add(shift))
	return massParticle
}

/*
Returns the radius of the partilce.
*/
func (massParticle *MassParticle) Radius() float32 {
	return massParticle.radius
}

/*
Returns the position of the particle.
*/
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

func (massParticle *MassParticle) SetVelocity(velocity mgl32.Vec3) *MassParticle {
	massParticle.velocity = velocity
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

func (massParticle *MassParticle) ApplyForce(force mgl32.Vec3) *MassParticle {
	v := massParticle.Velocity().Add(force)
	massParticle.SetVelocity(v)
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
            if (x == 0 || x == lenX-1) && (y == 0 || y == lenY-1) {
				xPos := -offsetX + radius + float32(x)*diameter + op.X()
				yPos := -offsetY + radius + float32(y)*diameter + op.Y()
				radiusHalf := radius / 2
				body.AddMassParticle(CreateMassParticle(xPos-radiusHalf, yPos-radiusHalf, 0, radiusHalf))
				body.AddMassParticle(CreateMassParticle(xPos+radiusHalf, yPos-radiusHalf, 0, radiusHalf))
				body.AddMassParticle(CreateMassParticle(xPos+radiusHalf, yPos+radiusHalf, 0, radiusHalf))
				body.AddMassParticle(CreateMassParticle(xPos-radiusHalf, yPos+radiusHalf, 0, radiusHalf))

			} else {
				xPos := -offsetX + radius + float32(x)*diameter + op.X()
				yPos := -offsetY + radius + float32(y)*diameter + op.Y()
				body.AddMassParticle(CreateMassParticle(xPos, yPos, 0, radius))
				if x == 0 || x == lenX-1 {
					if y != lenY-1 {
						xPos := -offsetX + radius + float32(x)*diameter + op.X()
						yPos := -offsetY + 2*radius + float32(y)*diameter + op.Y()
						body.AddMassParticle(
							CreateMassParticle(xPos, yPos, 0, radius),
						)
					}
				}
				if y == 0 || y == lenY-1 {
					if x != lenX-1 {
						xPos := -offsetX + 2*radius + float32(x)*diameter + op.X()
						yPos := -offsetY + radius + float32(y)*diameter + op.Y()
						body.AddMassParticle(
							CreateMassParticle(xPos, yPos, 0, radius),
						)
					}
				}
			}
		}
	}
}
