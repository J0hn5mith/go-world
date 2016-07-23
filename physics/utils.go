package go_world_physics

import (
    "github.com/go-gl/mathgl/mgl32"
    "github.com/gonum/matrix/mat64"
    "github.com/gonum/matrix"
)

var NUM_OF_ROUNDING_DIGITS int = 7


/*
    Extracts the rotation matrix regarding the point (0,0)
*/
func ExtractRotationFromPositions(old_positions, new_positions [] mgl32.Vec3) mgl32.Mat3 {
    a := outerProductSum(new_positions, old_positions)
    rotation := ExtractRotationFromMatrix(a)
    return rotation
}

func ExtractRotationFromMatrix(in_matrix mgl32.Mat3)mgl32.Mat3{
    m := mat64.NewDense(3, 3, toFloat64Array(matrixToArray((in_matrix))))
    var svd mat64.SVD
    _ = svd.Factorize(m, matrix.SVDFull)

    var u,v,rotation mat64.Dense
    u.UFromSVD(&svd)
    v.VFromSVD(&svd)
    rotation.Mul(&u, &v)

    return arrayToMatrix(toFloat32Array(rotation.RawMatrix().Data))
}

func outerProductSum(old_positions, new_positions [] mgl32.Vec3) mgl32.Mat3 {
    a  := mgl32.Mat3{}

    for i, old := range(old_positions) {
        a = a.Add(old.OuterProd3(new_positions[i]))
    }
    for i,value := range(a) {
        a[i] = mgl32.Round(value, NUM_OF_ROUNDING_DIGITS)
    }
    return a
}

func matrixToArray(matrix mgl32.Mat3) []float32{
    var array []float32
    for _, entry := range(matrix) {
        array = append(array, mgl32.Round(entry, NUM_OF_ROUNDING_DIGITS))
    }
    return array
}

func arrayToMatrix(array []float32) mgl32.Mat3{
    var matrix mgl32.Mat3
    for i, entry := range(array) {
        matrix[i] = mgl32.Round(entry, NUM_OF_ROUNDING_DIGITS)
    }
    return matrix
}

func toFloat64Array(values32 []float32) []float64 {
    var values64 []float64
    for _, value := range values32 {
        values64 = append(values64, float64(mgl32.Round(value, NUM_OF_ROUNDING_DIGITS)))
    }
    return values64
}

func toFloat32Array(values64 []float64) []float32 {
    var values32 []float32
    for _, value := range values64 {
        values32 = append(values32, mgl32.Round(float32(value), NUM_OF_ROUNDING_DIGITS))
    }
    return values32
}
