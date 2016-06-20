package go_world

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
)

var SPRING_CONSTANT float32 = 20000.0
var DAMPER_CONSTANT float32 = 10.03

type SoftBody struct {
	object          *Object
	velocity        mgl32.Vec3
	angularVelocity mgl32.Vec3
	centerOfMass    mgl32.Vec3
	massParticles   [][]*MassParticle
	mass            float32
	springs         []*Spring
}

func CreateSoftBody(object *Object) *SoftBody {
	softBody := new(SoftBody)
	softBody.object = object
	softBody.velocity = mgl32.Vec3{0, 0, 0}
	softBody.mass = 1.0

	return softBody
}

func (softBody *SoftBody) SetVelocity(x, y, z float32) PhysicalBody {
	softBody.velocity = mgl32.Vec3{x, y, z}
	return softBody
}

func (softBody *SoftBody) GetVelocity() mgl32.Vec3 {
	return softBody.velocity
}

func (softBody *SoftBody) Velocity() mgl32.Vec3 {
	return softBody.velocity
}

func (softBody *SoftBody) SetAngularVelocity(x, y, z float32) PhysicalBody {
	softBody.angularVelocity = mgl32.Vec3{x, y, z}
	return softBody
}

func (softBody *SoftBody) AngularVelocity() mgl32.Vec3 {
	return softBody.angularVelocity
}

func (softBody *SoftBody) AddMassParticle(massParticle *MassParticle) PhysicalBody {
	//softBody.AddMassParticleXYZ(massParticle, 0, 0, 0)
	fmt.Print("AddMassParticle is not implemented")
	return softBody
}

func (softBody *SoftBody) SetMassParticles(massParticles [][]*MassParticle) {
	softBody.massParticles = massParticles
}

func (softBody *SoftBody) GetMassParticles() []*MassParticle {
	var particles []*MassParticle
	for _, row := range softBody.massParticles {
		particles = append(particles, row...)
	}
	return particles
}

func (softBody *SoftBody) ApplyForce(x, y, z float32) PhysicalBody {
	v := softBody.GetVelocity()
	softBody.SetVelocity(
		v.X()+x,
		v.Y()+y,
		v.Z()+z,
	)
	return softBody
}

func (softBody *SoftBody) Object() *Object {
	return softBody.object
}

func (softBody *SoftBody) Mass() float32 {
	return softBody.mass
}
func (softBody *SoftBody) SetMass(mass float32) PhysicalBody {
	softBody.mass = mass
	return softBody
}
func (softBody *SoftBody) Springs() []*Spring {
	return softBody.springs
}

func (softBody *SoftBody) SetPosition(x, y, z float32) PhysicalBody {
	softBody.object.SetPosition(x, y, z)
	return softBody
}

func (softBody *SoftBody) Position() mgl32.Vec3 {
	return softBody.object.Position()
}

func (softBody *SoftBody) ShiftPosition(x, y, z float32) PhysicalBody {
	p := softBody.Position()
	softBody.object.SetPosition(
		p[0]+x,
		p[1]+y,
		p[2]+z,
	)
	return softBody
}

func (softBody *SoftBody) GetVertices() []float32 {
	var array []float32
	var lastRow []*MassParticle

	for _, row := range softBody.massParticles {
		if lastRow != nil {
			for i, particle := range row {
				x, y, z := VectorToFloats(lastRow[i].Position())
				array = append(array, x, y, z)
				x, y, z = VectorToFloats(particle.Position())
				array = append(array, x, y, z)
			}
		}
		lastRow = row
	}
	return array
	//return []float32{-0.5, -0.5, 0, 0.5, -0.5, 0, -0.5, 0.5, 0, 0.5, 0.5, 0}
}

func (softBody *SoftBody) UpdateSpringForces() *SoftBody {
	for _, particle := range softBody.GetMassParticles() {
		particle.SetSpringForce(0.0, 0.0, 0.0)
	}

	for _, spring := range softBody.Springs() {
		spring.Apply()
	}
	return softBody
}
func (softBody *SoftBody) ApplySpringForces(timeDelta float32) *SoftBody {
	for _, particle := range softBody.GetMassParticles() {
		force := particle.GetSpringForce().Mul(timeDelta)
		particle.ApplyForce(
			force.X(),
			force.Y(),
			force.Z(),
		)
	}
	return softBody
}
func AddMassParticleSoft2D(softBody *SoftBody, x, y, diameter float32) {
	radius := diameter / 2
	lenX := int(x / diameter)
	lenY := int(y / diameter)
	offsetX := x / 2
	offsetY := y / 2

	// TODO: Consider rest of division
	particles := [][]*MassParticle{}
	for x := 0; x < lenX; x++ {
		row := []*MassParticle{}
		for y := 0; y < lenY; y++ {
			xPos := -offsetX + radius + float32(x)*diameter
			yPos := -offsetY + radius + float32(y)*diameter
			row = append(row, CreateMassParticle(xPos, yPos, 0, radius))
		}
		particles = append(particles, row)
	}
	softBody.massParticles = particles

	springs := []*Spring{}
	for y, row := range softBody.massParticles {
		for x, particle := range row {
			if x+1 < lenX {
				springs = append(springs, AddSpring(particle, row[x+1]))
			}
			if y+1 < lenY {
                targetRow := particles[y+1]
                springs = append(springs, AddSpring(particle, targetRow[x]))
				if x+1 < lenX {
					springs = append(springs, AddSpring(particle, targetRow[x+1]))
				}
			}
		}
	}
	softBody.springs = springs
}

func AddSpring(source, destination *MassParticle) *Spring {
	spring := SpringFromMassParticles(
		source,
		destination,
		SPRING_CONSTANT,
		DAMPER_CONSTANT,
	)
	source.AddSpring(spring)
	return spring
}
