package go_world_physics
import (
    mgl "github.com/go-gl/mathgl/mgl64"
	"testing"
)

func TestCircleCollision(t * testing.T){
    // None
    c := CircleCollision(mgl.Vec3{0,0,0}, mgl.Vec3{2,0,0}, 1, 1)
    if c.Magnitude < 0 {
        t.Error("False positive", c.Magnitude, "> 0 ")
    }

    c2 := CircleCollision(mgl.Vec3{0,0,0}, mgl.Vec3{2,0,0}, 1, 1.5)
    if c2.Magnitude !=  0.5 {
        t.Error("Magnitude is incorrect", c2.Magnitude, "!= 0.5")
    }

    c3 := CircleCollision(mgl.Vec3{0,0,0}, mgl.Vec3{2,0,0}, 1.9, 1.9)
    if c3.Magnitude !=  1.8 {
        t.Error("Magnitude is incorrect", c3.Magnitude, "!= 1.8")
    }

    c4 := CircleCollision(mgl.Vec3{0,0,0}, mgl.Vec3{0.5,0,0}, 1, 1)
    if c4.Magnitude !=  1.5 {
        t.Error("Magnitude is incorrect", c4.Magnitude, "!= 1.5")
    }


}
