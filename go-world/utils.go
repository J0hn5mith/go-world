package go_world

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
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

func Max32(a, b float32) float32 {
	return float32(math.Max(float64(a), float64(b)))
}

func Min32(a, b float32) float32 {
	return float32(math.Min(float64(a), float64(b)))
}

func Sqrt32(a float32) float32 {
	return float32(math.Sqrt(float64(a)))
}

func Pow32(a, b float32) float32 {
	return float32(math.Pow(float64(a), float64(b)))
}

func Vec3To32(vector mgl64.Vec3) mgl32.Vec3 {
	return mgl32.Vec3{
		float32(vector.X()),
		float32(vector.Y()),
		float32(vector.Z()),
	}
}

func Vec4To32(vector mgl64.Vec4) mgl32.Vec4 {
	return mgl32.Vec4{
		float32(vector.X()),
		float32(vector.Y()),
		float32(vector.Z()),
		float32(vector.W()),
	}
}

func Mat4To32(matrix mgl64.Mat4) mgl32.Mat4 {
	return mgl32.Mat4{
		float32(matrix[0]), float32(matrix[1]), float32(matrix[2]), float32(matrix[3]),
		float32(matrix[4]), float32(matrix[5]), float32(matrix[6]), float32(matrix[7]),
		float32(matrix[8]), float32(matrix[9]), float32(matrix[10]), float32(matrix[11]),
		float32(matrix[12]), float32(matrix[13]), float32(matrix[14]), float32(matrix[15]),
	}
}
