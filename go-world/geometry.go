package go_world

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Geometry struct {
	vertices    []float32
	draw_method uint32
	color       mgl32.Vec3

	vao uint32
	vbo uint32
}

func NewGeometry(vertices []float32) *Geometry {
	geometry := new(Geometry)
	geometry.vertices = vertices
	geometry.color = mgl32.Vec3{0, 0, 0}
	return geometry
}

func (geometry *Geometry) SetColorRGB(r, g, b float32) *Geometry {
	geometry.color = mgl32.Vec3{r, g, b}
	return geometry
}

/*
Loads the geometry data to the GPU memory
*/
func (geometry *Geometry) Load(program uint32) *Geometry {
	gl.GenVertexArrays(1, &geometry.vao)
	gl.BindVertexArray(geometry.vao)
	gl.GenBuffers(1, &geometry.vbo)
	geometry.UpdateVertices(geometry.vertices)

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("position\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
	return geometry
}

func (geometry *Geometry) DrawMethod() uint32 {
    return geometry.draw_method
}
func (geometry *Geometry) SetDrawMethod(method uint32) *Geometry {
	geometry.draw_method = method
	return geometry
}

func (geometry *Geometry) Vertices() []float32 {
    return geometry.vertices
}

func (geometry *Geometry) UpdateVertices(vertices []float32) *Geometry {
	geometry.vertices = vertices
	gl.BindVertexArray(geometry.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, geometry.vbo)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		len(geometry.vertices)*4,
		gl.Ptr(geometry.vertices),
		gl.DYNAMIC_DRAW,
	)
	return geometry
}
