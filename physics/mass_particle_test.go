package go_world_physics
import (
    mgl "github.com/go-gl/mathgl/mgl64"
	"testing"
)

func TestCreateSphereTree(t * testing.T){
    var tree SphereTree
    rectangle := Rectangle{mgl.Vec3{0,0,0}, mgl.Vec3{1,1,1}}
    if tree.CreateSphereTree(rectangle, 0.5).root == nil{
        t.Error("Test")
    }
}
