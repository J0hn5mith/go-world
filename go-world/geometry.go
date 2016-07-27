package go_world

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	mgl "github.com/go-gl/mathgl/mgl64"
)

type Geometry struct {
	vertices    []float32
	drawMethod uint32
	color       mgl.Vec3
	program     uint32
	vao         uint32
	vbo         uint32
}

func NewGeometry(vertices []float64) *Geometry {
	geometry := new(Geometry)
	geometry.vertices = ToFloat32Array(vertices)
	geometry.color = mgl.Vec3{0, 0, 0}
    geometry.program = 0
	return geometry
}

/*
Sets the color of the geometry object.
*/
func (geometry *Geometry) SetColorRGB(r, g, b float64) *Geometry {
	geometry.color = mgl.Vec3{r, g, b}
	return geometry
}

func (geometry *Geometry) Color() mgl.Vec3 {
    return geometry.color
}

/*
   Loads the geometry data to the GPU memory
*/
func (geometry *Geometry) Load(program uint32) *Geometry {
    geometry.program = program
	gl.GenVertexArrays(1, &geometry.vao)
	gl.BindVertexArray(geometry.vao)
	gl.GenBuffers(1, &geometry.vbo)
	geometry.UpdateVertices(geometry.vertices)

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("position\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
	return geometry
}

/*
    Frees the GPU resources bound by this geometry object.
*/
func (geometry *Geometry) Delete() *Geometry {
    if geometry.program != 0 {
        gl.DeleteVertexArrays(1, &geometry.vao)
        gl.DeleteBuffers(1, &geometry.vbo)
    }
    return geometry
}

func (geometry *Geometry) DrawMethod() uint32 {
	return geometry.drawMethod
}

func (geometry *Geometry) SetDrawMethod(method uint32) *Geometry {
	geometry.drawMethod = method
	return geometry
}

func (geometry *Geometry) Vertices() []float32 {
	return geometry.vertices
}

func (geometry *Geometry) Vao() uint32 {
	return geometry.vao
}

func (geometry *Geometry) UpdateVertices(vertices []float32) *Geometry {
	geometry.vertices = vertices
	gl.BindVertexArray(geometry.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, geometry.vbo)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		len(vertices)*4,
		gl.Ptr(vertices),
		gl.DYNAMIC_DRAW,
	)
	return geometry
}
