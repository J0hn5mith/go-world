package go_world

import (
	"github.com/go-gl/mathgl/mgl32"
    "math"
)

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

func AngleToCoordinates32(angle float32) (float32, float32) {
    x, y := AngleToCoordinates(float64(angle))
    return float32(x), float32(y)
}

func AngleToCoordinates(angle float64) (float64, float64) {
    return angleToCoords(angle)
}

func angleToCoords(angle float64) (float64, float64) {
	x := math.Cos(angle)
	y := math.Sin(angle)
	return x, y
}

func reverse(numbers []interface{}) []interface{} {
	for i := 0; i < len(numbers)/2; i++ {
		j := len(numbers) - i - 1
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
	return numbers
}

func Abs32(a float32) float32 {
    return float32(math.Abs(float64(a)))
}

func Max32(a , b float32) float32 {
    return float32(math.Max(float64(a), float64(b)))
}

func Min32(a , b float32) float32 {
    return float32(math.Min(float64(a), float64(b)))
}

func Sqrt32(a float32) float32{
    return float32(math.Sqrt(float64(a)))
}

func Pow32(a, b float32) float32{
    return float32(math.Pow(float64(a), float64(b)))
}

