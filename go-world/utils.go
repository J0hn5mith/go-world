package go_world

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Collision struct {
	Direction mgl32.Vec3
	Magnitude float32
}

func TestCicrcleCollision(p1, p2 mgl32.Vec3, r1, r2 float32) Collision {
	magnitude := -(p1.Sub(p2).Len() - (r1 + r2))
	if magnitude > 0 {
		spring := p1.Sub(p2).Normalize()
		return Collision{spring, magnitude}
	}
	return Collision{mgl32.Vec3{0, 0, 0}, 0}
}
