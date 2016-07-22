package go_world_physics
import (
    "github.com/go-gl/mathgl/mgl32"
	"testing"
    "fmt"
)

func TestOuterProductSum(t *testing.T) {
    points := []mgl32.Vec3{{1, 1, 1},{1,1,1}, {1,1,1}}
    check := outerProductSum(points, points)
    for _, value := range(check){
        if value != 3 {
            t.Error("Value is not 3")
        }
    }
}

func TestExtractRotationMatrix(t *testing.T) {
    // Simplest test
    matrix := mgl32.Mat3{
        1, 1, 1,
        1, 1, 1,
        1, 1, 1,
    }
    check := ExtractRotationFromMatrix(matrix)
    fmt.Println(check)
    for _, value := range(check.Diag()){
        if value != 1 {
            t.Error("1 != ", value)
        }
    }

    // Test for simple case
    points_a := []mgl32.Vec3{{-1,0,0}, {1,0,0}}
    points_b := []mgl32.Vec3{{0,1,0}, {0, -1, 0}}
    matrix_result := mgl32.Mat3{
        0, -1, 0,
        -1, 0, 0,
        0, 0, 1,
    }
    tmp := outerProductSum(points_a, points_b)
    check_2 := ExtractRotationFromMatrix(tmp)
    for i, value := range(check_2){
        if value != matrix_result[i] {
            t.Error(matrix_result[i], " != ", value)
        }
    }

    // Test for slightly instable case
    points_c := []mgl32.Vec3{{-0.19999996, -0.20000005, 0}}
    points_d := []mgl32.Vec3{{-0.20000002, -0.20000005, 0}}
    tmp_3 := outerProductSum(points_c, points_d)
    check_3 := ExtractRotationFromMatrix(tmp_3)
    for _, value := range(check_3.Diag()){
        if value != 1 {
            t.Error("1 != ", value)
        }
    }
}

//flatten3x3(matrix mgl32.Mat3) []float32{
    //for 
//}
