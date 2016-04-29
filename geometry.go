package main

import (
    "github.com/go-gl/mathgl/mgl32"
    "github.com/go-gl/gl/v4.1-core/gl"
    "math"
    "fmt"
    "reflect"
)

type Geometry struct {
    vertices []float32
    draw_method  uint32
}

func createSphereGeometry(radius float32, rings, sectors float64) *Geometry {
    var R float64 = float64(1.0 / (rings - 1))
    var S float64 = float64(1.0 / (sectors - 1))

    num_vertices := uint32(rings * sectors * 3)
    vertices := make([]float32, num_vertices)

    var r, s float64
    var x, y, z float64
    i := 0

    for r = 0; r < rings; r++ {
        for s = 0; s < sectors; s++ {
            y = math.Sin(-math.Pi/2.0 + math.Pi*r*R)
            x = math.Cos(2.0*math.Pi*s*S) * math.Sin(math.Pi*r*R)
            z = math.Sin(2.0*math.Pi*s*S) * math.Sin(math.Pi*r*R)

            vertices[i] = float32(x)
            i += 1
            vertices[i] = float32(y)
            i += 1
            vertices[i] = float32(z)
            i += 1
        }
    }

    geometry := new(Geometry)
    geometry.vertices = vertices
    return geometry
}

func createDiamondGeometry(side_length float32) *Geometry {
    var vertices = []float32{
        0.5, 0.0, 0.0,
        0.0, 0.5, 0.0,

        0.0, 0.5, 0.0,
        -0.5, 0.0, 0.0,

        -0.5, 0.0, 0.0,
        0.0, -0.5, 0.0,

        0.0, -0.5, 0.0,
        0.5, 0.0, 0.0,
    }
    geometry := new(Geometry)
    geometry.vertices = vertices
    return geometry
}

func createTriangle2DGeometry(side_length float32) *Geometry {
    side_length = side_length
    var vertices = []float32{
        -side_length / 2.0, -side_length / 2, 0.0,
        side_length / 2.0, -side_length / 2, 0.0,
        0, side_length / 2.0, 0.0,
        0, side_length / 2.0, 0.0,
        //side_length/2, -side_length/2, 0,
        //side_length/2, side_length/2, 0,
        //0, side_length/2, 0,
    }
    geometry := new(Geometry)
    geometry.vertices = vertices
    return geometry
}

func createCubeGeometry(side_length float32) *Geometry {
    side_length = side_length
    var vertices = []float32{
        // Bottom
        -side_length, -side_length, -side_length,
        side_length, -side_length, -side_length,
        -side_length, -side_length, side_length,
        side_length, -side_length, -side_length,
        side_length, -side_length, side_length,
        -side_length, -side_length, side_length,

        // Top
        -side_length, side_length, -side_length,
        -side_length, side_length, side_length,
        side_length, side_length, -side_length,
        side_length, side_length, -side_length,
        -side_length, side_length, side_length,
        side_length, side_length, side_length,

        // Front
        -side_length, -side_length, side_length,
        side_length, -side_length, side_length,
        -side_length, side_length, side_length,
        side_length, -side_length, side_length,
        side_length, side_length, side_length,
        -side_length, side_length, side_length,

        // Back
        -side_length, -side_length, -side_length,
        -side_length, side_length, -side_length,
        side_length, -side_length, -side_length,
        side_length, -side_length, -side_length,
        -side_length, side_length, -side_length,
        side_length, side_length, -side_length,

        // Left
        -side_length, -side_length, side_length,
        -side_length, side_length, -side_length,
        -side_length, -side_length, -side_length,
        -side_length, -side_length, side_length,
        -side_length, side_length, side_length,
        -side_length, side_length, -side_length,

        // Right
        side_length, -side_length, side_length,
        side_length, -side_length, -side_length,
        side_length, side_length, -side_length,
        side_length, -side_length, side_length,
        side_length, side_length, -side_length,
        side_length, side_length, side_length,
    }
    geometry := new(Geometry)
    geometry.vertices = vertices
    return geometry
}

func createCircleGeometry(num_vertices int) *Geometry {
    vertices := []float32{}

    px, py := angleToCoords(0)

    for i := 1; i < int(num_vertices); i++ {
        var angle float64 = (2 * math.Pi * float64(i) / float64(num_vertices))
        x, y := angleToCoords(angle)
        vertices = append(
            vertices, 
            float32(x), float32(y), 0.0,
            float32(px), float32(py), 0.0,
        )
        px = x
        py = y
    }

    x, y := angleToCoords(0)
    vertices = append(
        vertices, 
        float32(x), float32(y), 0.0,
        float32(px), float32(py), 0.0,
    )

    geometry := new(Geometry)
    geometry.vertices = vertices
    geometry.draw_method = gl.LINES
    fmt.Println(reflect.TypeOf(gl.LINES))
    return geometry
}

func angleToCoords(angle float64) (float64, float64) {
    x := math.Sin(angle)
    y := math.Cos(angle)
    return x, y
}
func createOctahedronGeometry() *Geometry {
    p0 := mgl32.Vec3{0, 1, 0}
    p1 := mgl32.Vec3{-1, 0, 1}
    p2 := mgl32.Vec3{1, 0, 1}
    p3 := mgl32.Vec3{1, 0, -1}
    p4 := mgl32.Vec3{-1, 0, -1}
    p5 := mgl32.Vec3{0, -1, 0}

    t1 := Triangle{p1, p2, p0}
    t2 := Triangle{p2, p3, p0}
    t3 := Triangle{p3, p4, p0}
    t4 := Triangle{p4, p1, p0}
    t5 := Triangle{p5, p1, p2}
    t6 := Triangle{p5, p2, p3}
    t7 := Triangle{p5, p3, p4}
    t8 := Triangle{p5, p4, p1}

    geometry := new(Geometry)
    geometry.vertices = t_to_array(t1, t2, t3, t4, t5, t6, t7, t8)
    return geometry
}

func to_array(vertices ...mgl32.Vec3) []float32 {
    var array []float32
    for _, v := range vertices {
        array = append(array, v[0], v[1], v[2])
    }
    return array
}

type Triangle struct {
    v1, v2, v3 mgl32.Vec3
}

func t_to_array(triangles ...Triangle) []float32 {
    var array []float32
    for _, t := range triangles {
        values := to_array(t.v1, t.v2, t.v3)
        array = append(array, values...)
    }
    return array
}
