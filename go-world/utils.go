package go_world

import (
	"github.com/go-gl/mathgl/mgl32"
    "math"
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

func VectorToFloats(vector mgl32.Vec3) (x, y, z float32) {
		return vector[0], vector[1], vector[2]
}

func VectorsToFloats(vertices ...mgl32.Vec3) []float32 {
	var array []float32
	for _, v := range vertices {
		x, y, z := VectorToFloats(v)
		array = append(array, x, y, z)
	}
	return array
}

func angleToCoords(angle float64) (float64, float64) {
	x := math.Sin(angle)
	y := math.Cos(angle)
	return x, y
}
func reverse(numbers []interface{}) []interface{} {
	for i := 0; i < len(numbers)/2; i++ {
		j := len(numbers) - i - 1
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
	return numbers
}
