package go_world_physics

import (
	mgl "github.com/go-gl/mathgl/mgl64"
    "go-world/go-world"
)

type Rectangle struct {
	Position  mgl.Vec3
	Dimension mgl.Vec3
}

type Sphere struct {
	Position mgl.Vec3
	Radius   float64
}

type Collision struct {
	Direction mgl.Vec3
	Magnitude float64
}

type SphereTreeNode struct {
    sphere *Sphere
    subNodes [2]*SphereTreeNode
}

type SphereTree struct {
    root *SphereTreeNode
    position mgl.Vec3
    leafs []*SphereTreeNode
}

/*
Checks two circles for a collision based on their center and radius
Legacy
*/
func CircleCollision(p1, p2 mgl.Vec3, r1, r2 float64) Collision {
	distance := p1.Sub(p2).Len()
	magnitude := -(distance - (r1 + r2))
	if magnitude > 0 {
		normal := p1.Sub(p2).Normalize()
		return Collision{normal, go_world.Round(magnitude)}
	}
	return Collision{mgl.Vec3{0, 0, 0}, 0}
}

func SphereCollision(sphereA, sphereB *Sphere) Collision {
	delta := sphereA.Position.Sub(sphereB.Position)
	distance := delta.Len()
	magnitude := -(distance - (sphereA.Radius + sphereB.Radius))
	if magnitude > 0 {
		normal := delta.Normalize()
		return Collision{normal, go_world.Round(magnitude)}
	}
	return Collision{mgl.Vec3{0, 0, 0}, 0}
}

/*
Checks two particle for a collision
*/
func DetectInterParticleCollision(particleA, particleB *MassParticle) Collision {
	return SphereCollision(
		&Sphere{particleA.Position(), particleA.Radius()},
		&Sphere{particleB.Position(), particleB.Radius()},
	)
}

func InterSphereCollisions(spheresA, spheresB []*Sphere) []Collision {
	var collisions []Collision
	for _, sphereA := range spheresA {
		for _, sphereB := range spheresB {
			collision := SphereCollision(sphereA, sphereB)
			if collision.Magnitude > 0 {
				collisions = append(collisions, collision)
			}
		}
	}
	return collisions
}
