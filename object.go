package go_world

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Object struct {
	geometry *Geometry
	position mgl32.Vec3
	rotation mgl32.Vec3
	scale    mgl32.Vec3

	model        mgl32.Mat4
	modelUniform int32
	angle        float64
	vao          uint32
	vbo          uint32
}

func NewObject(geometry *Geometry) *Object {
	object := new(Object)
	object.geometry = geometry

	object.position = mgl32.Vec3{0, 0, 0}
	object.rotation = mgl32.Vec3{0, 0, 0}
	object.scale = mgl32.Vec3{1, 1, 1}

	object.angle = 0
	return object
}

func (object *Object) SetScale(scale float32) {
	object.setScale(scale)
}

func (object *Object) setScale(scale float32) {
	object.scale = mgl32.Vec3{scale, scale, scale}
}

func (o *Object) SetPosition(x, y, z float32) *Object {
	o.position = mgl32.Vec3{x, y, z}
	return o
}

func (o Object) Position() mgl32.Vec3 {
	return o.position
}

func (object *Object) Configure(program uint32) {
	object.configure(program)
}
func (object *Object) configure(program uint32) {

	gl.GenVertexArrays(1, &object.vao)
	gl.BindVertexArray(object.vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(object.geometry.vertices)*5, gl.Ptr(object.geometry.vertices), gl.STATIC_DRAW)
	object.vbo = vbo

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("position\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	// Configure global settings
    gl.Enable(gl.DEPTH_TEST)
    gl.DepthFunc(gl.LESS)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
}
