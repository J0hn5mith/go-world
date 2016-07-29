package go_world_physics

import (
    mgl "github.com/go-gl/mathgl/mgl64"
    "testing"
)

func TestOuterProductSum(t *testing.T) {
    //t.Skip()
    points := []mgl.Vec3{{1, 1, 1}, {1, 1, 1}, {1, 1, 1}}
    check := outerProductSum(points, points)
    for _, value := range check {
        if value != 3 {
            t.Error("Value is not 3")
        }
    }
}

type Test struct {
    name    string
    pointsA []mgl.Vec3
    pointsB []mgl.Vec3
    result  mgl.Mat3
}

type ExtractRotationTest struct {
    name    string
    matrix  mgl.Mat3
    result  mgl.Mat3
}

func TestExtractRotationMatrix(t *testing.T) {
    //t.Skip()
    //Simplest test
    matrix := mgl.Mat3{
        0, 1, 0,
        0, 0, 0,
        0, 0, 0,
    }
    check := ExtractRotationFromMatrix(matrix)
    for _, value := range check.Diag() {
        if value != 1 {
            t.Error("1 != ", value)
        }
    }


    // Test for simple case
    points_a := []mgl.Vec3{{-1, 0, 0}, {1, 0, 0}}
    points_b := []mgl.Vec3{{0, 1, 0}, {0, -1, 0}}
    matrix_result := mgl.Mat3{
        0, -1, 0,
        -1, 0, 0,
        0, 0, 1,
    }
    tmp := outerProductSum(points_a, points_b)
    check_2 := ExtractRotationFromMatrix(tmp)
    for i, value := range check_2 {
        if value != matrix_result[i] {
            t.Error(matrix_result[i], " != ", value)
        }
    }

    // Test for slightly instable case
    points_c := []mgl.Vec3{{-0.19999996, -0.20000005, 0}}
    points_d := []mgl.Vec3{{-0.20000002, -0.20000005, 0}}
    tmp_3 := outerProductSum(points_c, points_d)
    check_3 := ExtractRotationFromMatrix(tmp_3)
    for _, value := range check_3.Diag() {
        if value != 1 {
            t.Error("1 != ", value)
        }
    }

    extractionTestCases := []ExtractRotationTest {
        {
            "Basic Decomposition",
            mgl.Mat3{
                3, 8, 2,
                2,5,7,
                1,4,6,
            },
            mgl.Mat3{
                0.3019, 0.9175,  -0.2589,
                0.6774, -0.0153, 0.7355,
                -0.6708, 0.3974, 0.6261,
            },
        },
    }
    for _, test := range extractionTestCases  {
        check := ExtractRotationFromMatrix(test.matrix)
        for i, value := range check {
            value = mgl.Round(value, 4) 
            if value != test.result[i] {
                t.Error("\nTest \"", test.name, "\" failed\n", test.result[i], " != ", value, "\n", check)
            }
        }
    }
    tests := []Test{
        {
            "Z-Axis Rotation",
            []mgl.Vec3{{1.0, 1.0, 0}, {1.0, 1, 0}},
            []mgl.Vec3{{1.0, 1.0, 0}, {1.0, 1.0, 0}},
            mgl.Mat3{
                1, 0, 0,
                0, 1, 0,
                0, 0, 1,
            },
        },
        {
            "X-Axis Rotation",
            []mgl.Vec3{{0.0, 1.0, 1.0}, {0.0, 1.0, 1.0}},
            []mgl.Vec3{{0.0, 1.0, 1.0}, {0.0, 1.0, 1.0}},
            mgl.Mat3{
                1, 0, 0,
                0, 1, 0,
                0, 0, 1,
            },
        },
        {
            "Y-Axis Rotation",
            []mgl.Vec3{{1.0, 0, 1.0}, {1.0, 0, 1.0}},
            []mgl.Vec3{{1.0, 0, 1.0}, {1.0, 0, 1.0}},
            mgl.Mat3{
                1, 0, 0,
                0, 1, 0,
                0, 0, 1,
            },
        },
    }
    for _, test := range tests {
        tmp := outerProductSum(test.pointsA, test.pointsB)
        check := ExtractRotationFromMatrix(tmp)
        for i, value := range check {
            result := mgl.Round(test.result[i], 4)
            value := mgl.Round(value, 4)
            if value !=  result {
                t.Error("\nTest \"", test.name, "\" failed\n", result, " != ", value, "\n", check)
            }
        }
    }
}
