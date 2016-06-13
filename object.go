package go_world

import (
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
