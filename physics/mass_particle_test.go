package go_world_physics
import (
    "github.com/go-gl/mathgl/mgl32"
	"testing"
)

func TestCreateSphereTree(t * testing.T){
    var tree SphereTree
    rectangle := Rectangle{mgl32.Vec3{0,0,0}, mgl32.Vec3{1,1,1}}
    if tree.CreateSphereTree(rectangle, 0.5).root == nil{
        t.Error("Test")
    }
    t.Error(len(tree.leafs))
}
