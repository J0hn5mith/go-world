package go_world_physics

import (
	mgl "github.com/go-gl/mathgl/mgl64"
	"math"
)

type MassParticle struct {
	position    mgl.Vec3 // Position relative to the center of mass
	velocity    mgl.Vec3
	radius      float64
	collided    bool
	springs     []*Spring
	springForce mgl.Vec3
}

/*
Creates a new mass particles.
*/
func CreateMassParticle(x, y, z, radius float64) *MassParticle {
	massParticle := new(MassParticle)
	massParticle.position = mgl.Vec3{x, y, z}
	massParticle.velocity = mgl.Vec3{0, 0, 0}
	massParticle.radius = radius
	massParticle.collided = false

	return massParticle
}

/*
Returns the position of the mass particle.
*/
func (massParticle *MassParticle) Position() mgl.Vec3 {
	return massParticle.position
}

/*
Sets the position of the mass particle.
*/
func (massParticle *MassParticle) SetPosition(position mgl.Vec3) *MassParticle {
	massParticle.position = position
	return massParticle
}

/*
Shifts the position of the particle by provided values.
*/
func (massParticle *MassParticle) ShiftPosition(shift mgl.Vec3) *MassParticle {
	massParticle.SetPosition(massParticle.Position().Add(shift))
	return massParticle
}

/*
Returns the radius of the partilce.
*/
func (massParticle *MassParticle) Radius() float64 {
	return massParticle.radius
}

/*
Returns the position of the particle.
*/
func (massParticle *MassParticle) Velocity() mgl.Vec3 {
	return massParticle.velocity
}

func (massParticle *MassParticle) Collided() bool {
	return massParticle.collided
}

func (massParticle *MassParticle) SetCollided(collided bool) *MassParticle {
	massParticle.collided = collided
	return massParticle
}

func (massParticle *MassParticle) SetVelocity(velocity mgl.Vec3) *MassParticle {
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

func (massParticle *MassParticle) GetSpringForce() mgl.Vec3 {
	return massParticle.springForce
}

func (massParticle *MassParticle) SetSpringForce(x, y, z float64) *MassParticle {
	massParticle.springForce = mgl.Vec3{x, y, z}
	return massParticle
}

func (massParticle *MassParticle) ApplySpringForce(force mgl.Vec3) *MassParticle {
	massParticle.springForce = massParticle.springForce.Add(force)
	return massParticle
}

func (massParticle *MassParticle) ApplyForce(force mgl.Vec3) *MassParticle {
	v := massParticle.Velocity().Add(force)
	massParticle.SetVelocity(v)
	return massParticle
}

type Spring struct {
	source         *MassParticle
	target         *MassParticle
	length         float64
	springConstant float64
	damperConstant float64
}

func SpringFromMassParticles(particle1, particle2 *MassParticle, springConstant, damperConstant float64) *Spring {
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
func AddMassParticle(body PhysicalBody, dimension mgl.Vec3, diameter float64) {
	x, y, z := dimension.Elem()
	radius := diameter / 2
	lenX := int(x / diameter)
	lenY := int(y / diameter)
	lenZ := int(z / diameter)
	offsetX := x / 2
	offsetY := y / 2
	offsetZ := z / 2

	// TODO: Consider rest of division
	op := body.Position()
	for z := 0; z < lenZ; z++ {
		for x := 0; x < lenX; x++ {
			for y := 0; y < lenY; y++ {
				//if (x == 0 || x == lenX-1) && (y == 0 || y == lenY-1) && (z == 0 || z == lenZ-1) {
					//xPos := -offsetX + radius + float64(x)*diameter + op.X()
					//yPos := -offsetY + radius + float64(y)*diameter + op.Y()
					//zPos := -offsetZ + radius + float64(z)*diameter + op.Z()
					//radiusHalf := radius / 2
					//body.AddMassParticle(CreateMassParticle(xPos-radiusHalf, yPos-radiusHalf, zPos, radiusHalf))
					//body.AddMassParticle(CreateMassParticle(xPos+radiusHalf, yPos-radiusHalf, zPos, radiusHalf))
					//body.AddMassParticle(CreateMassParticle(xPos+radiusHalf, yPos+radiusHalf, zPos, radiusHalf))
					//body.AddMassParticle(CreateMassParticle(xPos-radiusHalf, yPos+radiusHalf, zPos, radiusHalf))

				//} else {
					xPos := -offsetX + radius + float64(x)*diameter + op.X()
					yPos := -offsetY + radius + float64(y)*diameter + op.Y()
					zPos := -offsetZ + radius + float64(z)*diameter + op.Z()
					body.AddMassParticle(CreateMassParticle(xPos, yPos, zPos, radius))
					//if x == 0 || x == lenX-1 {
						//if y != lenY-1 {
							//if z != lenZ-1 {
								//xPos := -offsetX + radius + float64(x)*diameter + op.X()
								//yPos := -offsetY + 2*radius + float64(y)*diameter + op.Y()
								//body.AddMassParticle(
									//CreateMassParticle(xPos, yPos, zPos, radius),
								//)
							//}
						//}
					//}
					//if y == 0 || y == lenY-1 {
						//if x != lenX-1 {
							//if z != lenZ-1 {
								//xPos := -offsetX + 2*radius + float64(x)*diameter + op.X()
								//yPos := -offsetY + radius + float64(y)*diameter + op.Y()
								//zPos := -offsetZ + radius + float64(z)*diameter + op.Z()
								//body.AddMassParticle(
									//CreateMassParticle(xPos, yPos, zPos, radius),
								//)
							//}
						//}
					//}
					//if z == 0 || z == lenZ-1 {
						//if x != lenX-1 {
							//if y != lenZ-1 {
								//xPos := -offsetX + radius + float64(x)*diameter + op.X()
								//yPos := -offsetY + radius + float64(y)*diameter + op.Y()
								//zPos := -offsetZ + 2*radius + float64(z)*diameter + op.Z()
								//body.AddMassParticle(
									//CreateMassParticle(xPos, yPos, zPos, radius),
								//)
							//}
						//}
					//}
				//}
			}
		}
	}
	bsRadius := math.Sqrt(
		math.Pow((x/2), 2) + math.Pow(y/2, 2) + math.Pow(z/2, 2),
	)
	body.AddBoundingSphere(&Sphere{mgl.Vec3{0, 0, 0}, bsRadius})
}

func (tree *SphereTree) CreateSphereTree(rectangle Rectangle, sphereRadius float64) *SphereTree {
	tree.root = tree.decomposeRectangle(rectangle, sphereRadius)
	return tree
}

func (tree *SphereTree) decomposeRectangle(rectangle Rectangle, sphereRadius float64) *SphereTreeNode {
	radius := math.Sqrt(
		math.Pow((rectangle.Dimension.X()/2), 2) +
			math.Pow(rectangle.Dimension.Y()/2, 2),
	)

	if radius < sphereRadius {
		node := SphereTreeNode{
			&Sphere{rectangle.Position, radius},
			[2]*SphereTreeNode{nil, nil},
		}
		tree.leafs = append(tree.leafs, &node)
		return &node
	}

	var shift mgl.Vec3
	if rectangle.Dimension.X()/rectangle.Dimension.Y() < 0.5 { // Horizontal Splitting
		shift = mgl.Vec3{0, rectangle.Dimension.Y() / 2, 0}
	} else {
		shift = mgl.Vec3{rectangle.Dimension.X() / 2, 0, 0}
	}
	newSize := rectangle.Dimension.Sub(shift.Mul(2))
	positionA := rectangle.Position.Sub(shift)
	positionB := rectangle.Position.Add(shift)
	nodeA := tree.decomposeRectangle(Rectangle{positionA, newSize}, sphereRadius)
	nodeB := tree.decomposeRectangle(Rectangle{positionB, newSize}, sphereRadius)

	return &SphereTreeNode{
		&Sphere{rectangle.Position, radius},
		[2]*SphereTreeNode{nodeA, nodeB},
	}
}
