package go_world_physics

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Sphere struct {
	Position mgl32.Vec3
	Radius   float32
}
type Collision struct {
	Direction mgl32.Vec3
	Magnitude float32
}

/*
Checks two circles for a collision based on their center and radius
Legacy
*/
func CircleCollision(p1, p2 mgl32.Vec3, r1, r2 float32) Collision {
	distance := p1.Sub(p2).Len()
	magnitude := -(distance - (r1 + r2))
	if magnitude > 0 {
		normal := p1.Sub(p2).Normalize()
		return Collision{normal, magnitude}
	}
	return Collision{mgl32.Vec3{0, 0, 0}, 0}
}

func SphereCollision(sphereA, sphereB *Sphere) Collision {
    delta := sphereA.Position.Sub(sphereB.Position)
	distance := delta.Len()
	magnitude := -(distance - (sphereA.Radius + sphereB.Radius))
	if magnitude > 0 {
		normal := delta.Normalize()
		return Collision{normal, magnitude}
	}
	return Collision{mgl32.Vec3{0, 0, 0}, 0}
}

/*
Checks two particle for a collision
*/
func DetectInterParticleCollision(particleA, particleB *MassParticle) Collision {
	return CircleCollision(
		particleA.Position(),
		particleB.Position(),
		particleA.Radius(),
		particleB.Radius(),
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
