package go_world_utils

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	mgl "github.com/go-gl/mathgl/mgl64"
	"go-world/go-world"
	"math"
)

/*
Geometry Creators
*/
func CreateBoxGeometry(dimension mgl.Vec3) *go_world.Geometry {
	return createBoxGeometry(dimension)
}

func CreateWireBoxGeometry(width, height, depth float64) *go_world.Geometry {
	vertices := CreateWireBoxVertices(width, height, depth)
	geometry := go_world.NewGeometry(vertices)
	geometry.SetDrawMethod(gl.LINES)
	return geometry
}

func CreateCubeGeometry(sideLength float64) *go_world.Geometry {
	return createCubeGeometry(sideLength)
}

func CreateLineGeometry(start, end mgl.Vec3) *go_world.Geometry {
	var vertices = []float64{
		start.X(), start.Y(), start.Z(),
		end.X(), end.Y(), end.Z(),
	}

	geometry := go_world.NewGeometry(vertices)
	geometry.SetDrawMethod(gl.LINES)
	return geometry
}

func CreateLineLoopGeometry(points ...mgl.Vec3) *go_world.Geometry {
	vertices := []float64{}
	for i, point := range points {
		pn := points[(i+1)%len(points)]
		vertices = append(vertices, point.X(), point.Y(), point.Z(), pn.X(), pn.Y(), pn.Z())
	}
	geometry := go_world.NewGeometry(vertices)
	geometry.SetDrawMethod(gl.LINES)

	return geometry
}

func CreateCircleLineGeometry(num_vertices int, radius float64) *go_world.Geometry {
	return CreateCircleLineGeometry(num_vertices, radius)
}
func CreateCircleGeometry(num_vertices int, radius float64) *go_world.Geometry {
	return createCircleGeometry(num_vertices, radius)
}

func CreatePlaneGeometry(sideLength float64) *go_world.Geometry {
	halfSideLength := sideLength / 2.0
	geometry := go_world.NewGeometry([]float64{
		-halfSideLength, 0, halfSideLength,
		halfSideLength, 0, halfSideLength,
		-halfSideLength, 0, -halfSideLength,
		halfSideLength, 0, -halfSideLength,
	})
	geometry.SetDrawMethod(gl.TRIANGLE_STRIP)

	return geometry
}

func createSphereGeometry(radius float64, rings, sectors float64) *go_world.Geometry {
	var R float64 = float64(1.0 / (rings - 1))
	var S float64 = float64(1.0 / (sectors - 1))

	num_vertices := uint32(rings * sectors * 3)
	vertices := make([]float64, num_vertices)

	var r, s float64
	var x, y, z float64
	i := 0

	for r = 0; r < rings; r++ {
		for s = 0; s < sectors; s++ {
			y = math.Sin(-math.Pi/2.0 + math.Pi*r*R)
			x = math.Cos(2.0*math.Pi*s*S) * math.Sin(math.Pi*r*R)
			z = math.Sin(2.0*math.Pi*s*S) * math.Sin(math.Pi*r*R)

			vertices[i] = float64(x)
			i += 1
			vertices[i] = float64(y)
			i += 1
			vertices[i] = float64(z)
			i += 1
		}
	}

	geometry := go_world.NewGeometry(vertices)
	return geometry
}

func createDiamondGeometry(side_length float64) *go_world.Geometry {
	var vertices = []float64{
		0.5, 0.0, 0.0,
		0.0, 0.5, 0.0,

		0.0, 0.5, 0.0,
		-0.5, 0.0, 0.0,

		-0.5, 0.0, 0.0,
		0.0, -0.5, 0.0,

		0.0, -0.5, 0.0,
		0.5, 0.0, 0.0,
	}
	geometry := go_world.NewGeometry(vertices)
	return geometry
}

func createTriangle2DGeometry(side_length float64) *go_world.Geometry {
	side_length = side_length
	var vertices = []float64{
		-side_length / 2.0, -side_length / 2, 0.0,
		side_length / 2.0, -side_length / 2, 0.0,
		0, side_length / 2.0, 0.0,
		0, side_length / 2.0, 0.0,
	}
	geometry := go_world.NewGeometry(vertices)
	return geometry
}

func createCubeGeometry(sideLength float64) *go_world.Geometry {
	return createBoxGeometry(mgl.Vec3{sideLength, sideLength, sideLength})
}

func createBoxGeometry(dimension mgl.Vec3) *go_world.Geometry {
	vertices := CreateBoxVertices(dimension.X(), dimension.Y(), dimension.Z())
	geometry := go_world.NewGeometry(vertices)
	geometry.SetDrawMethod(gl.TRIANGLES)
	return geometry
}

func createCircleLineGeometry(num_vertices int, radius float64) *go_world.Geometry {
	vertices := []float64{}

	px, py := angleToCoords(0)

	for i := 1; i < int(num_vertices); i++ {
		var angle float64 = (2 * math.Pi * float64(i) / float64(num_vertices))
		x, y := angleToCoords(angle)
		vertices = append(
			vertices,
			float64(x)*radius, float64(y)*radius, 0.0,
			float64(px)*radius, float64(py)*radius, 0.0,
		)
		px = x
		py = y
	}

	x, y := angleToCoords(0)
	vertices = append(
		vertices,
		float64(x)*radius, float64(y)*radius, 0.0,
		float64(px)*radius, float64(py)*radius, 0.0,
	)

	geometry := go_world.NewGeometry(vertices)
	geometry.SetDrawMethod(gl.LINE_LOOP)
	return geometry
}

func createCircleGeometry(num_vertices int, radius float64) *go_world.Geometry {
	vertices := []float64{}

	px, py := angleToCoords(0)

	for i := 1; i < int(num_vertices); i++ {
		var angle float64 = (2 * math.Pi * float64(i) / float64(num_vertices))
		x, y := angleToCoords(angle)
		vertices = append(
			vertices,
			0, 0, 0,
			float64(x)*radius, float64(y)*radius, 0.0,
			float64(px)*radius, float64(py)*radius, 0.0,
		)
		px = x
		py = y
	}

	x, y := angleToCoords(0)
	vertices = append(
		vertices,
		0, 0, 0,
		float64(x)*radius, float64(y)*radius, 0.0,
		float64(px)*radius, float64(py)*radius, 0.0,
	)

	geometry := go_world.NewGeometry(vertices)
	geometry.SetDrawMethod(gl.TRIANGLES)
	return geometry
}

func angleToCoords(angle float64) (float64, float64) {
	x := math.Sin(angle)
	y := math.Cos(angle)
	return x, y
}

func createOctahedronGeometry() *go_world.Geometry {
	p0 := mgl.Vec3{0, 1, 0}
	p1 := mgl.Vec3{-1, 0, 1}
	p2 := mgl.Vec3{1, 0, 1}
	p3 := mgl.Vec3{1, 0, -1}
	p4 := mgl.Vec3{-1, 0, -1}
	p5 := mgl.Vec3{0, -1, 0}

	t1 := Triangle{p1, p2, p0}
	t2 := Triangle{p2, p3, p0}
	t3 := Triangle{p3, p4, p0}
	t4 := Triangle{p4, p1, p0}
	t5 := Triangle{p5, p1, p2}
	t6 := Triangle{p5, p2, p3}
	t7 := Triangle{p5, p3, p4}
	t8 := Triangle{p5, p4, p1}

	geometry := go_world.NewGeometry(t_to_array(t1, t2, t3, t4, t5, t6, t7, t8))
	return geometry
}

/*
Vertex Creators
*/
func CreateBoxVertices(width, height, depth float64) []float64 {
	xOffset := width / 2
	yOffset := height / 2
	zOffset := depth / 2

	var vertices = []float64{
		// Bottom
		-xOffset, -yOffset, -zOffset,
		xOffset, -yOffset, -zOffset,
		-xOffset, -yOffset, zOffset,
		xOffset, -yOffset, -zOffset,
		xOffset, -yOffset, zOffset,
		-xOffset, -yOffset, zOffset,

		// Top
		-xOffset, yOffset, -zOffset,
		-xOffset, yOffset, zOffset,
		xOffset, yOffset, -zOffset,
		xOffset, yOffset, -zOffset,
		-xOffset, yOffset, zOffset,
		xOffset, yOffset, zOffset,

		// Front
		-xOffset, -yOffset, zOffset,
		xOffset, -yOffset, zOffset,
		-xOffset, yOffset, zOffset,
		xOffset, -yOffset, zOffset,
		xOffset, yOffset, zOffset,
		-xOffset, yOffset, zOffset,

		// Back
		-xOffset, -yOffset, -zOffset,
		-xOffset, yOffset, -zOffset,
		xOffset, -yOffset, -zOffset,
		xOffset, -yOffset, -zOffset,
		-xOffset, yOffset, -zOffset,
		xOffset, yOffset, -zOffset,

		// Left
		-xOffset, -yOffset, zOffset,
		-xOffset, yOffset, -zOffset,
		-xOffset, -yOffset, -zOffset,
		-xOffset, -yOffset, zOffset,
		-xOffset, yOffset, zOffset,
		-xOffset, yOffset, -zOffset,

		// Right
		xOffset, -yOffset, zOffset,
		xOffset, -yOffset, -zOffset,
		xOffset, yOffset, -zOffset,
		xOffset, -yOffset, zOffset,
		xOffset, yOffset, -zOffset,
		xOffset, yOffset, zOffset,
	}
	return vertices
}
func CreateWireBoxVertices(width, height, depth float64) []float64 {
	width = width / 2
	height = height / 2
	depth = depth / 2
	var vertices = []float64{
		//Front
		-width, -height, depth,
		width, -height, depth,

		width, -height, depth,
		width, height, depth,

		width, height, depth,
		-width, height, depth,

		-width, height, depth,
		-width, -height, depth,

		//Middle
		-width, -height, depth,
		-width, -height, -depth,

		width, -height, depth,
		width, -height, -depth,

		width, height, depth,
		width, height, -depth,

		-width, height, depth,
		-width, height, -depth,

		//Back
		-width, -height, -depth,
		width, -height, -depth,

		width, -height, -depth,
		width, height, -depth,

		width, height, -depth,
		-width, height, -depth,

		-width, height, -depth,
		-width, -height, -depth,
	}
	return vertices
}

func to_array(vertices ...mgl.Vec3) []float64 {
	var array []float64
	for _, v := range vertices {
		array = append(array, v[0], v[1], v[2])
	}
	return array
}

type Triangle struct {
	v1, v2, v3 mgl.Vec3
}

func t_to_array(triangles ...Triangle) []float64 {
	var array []float64
	for _, t := range triangles {
		values := to_array(t.v1, t.v2, t.v3)
		array = append(array, values...)
	}
	return array
}
