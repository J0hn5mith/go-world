package go_world

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Object struct {
	geometry *Geometry
	position mgl32.Vec3
	rotation mgl32.Mat4
	scale    mgl32.Vec3

	model        mgl32.Mat4
	modelUniform int32
	vao          uint32
	vbo          uint32
}

func NewObject(geometry *Geometry) *Object {
	object := new(Object)
	object.geometry = geometry

	object.position = mgl32.Vec3{0, 0, 0}
	object.rotation = mgl32.Ident4()
    object.scale = mgl32.Vec3{1, 1, 1}

	return object
}

func (o *Object) SetPosition(x, y, z float32) *Object {
	o.position = mgl32.Vec3{x, y, z}
	return o
}

func (o Object) Position() mgl32.Vec3 {
	return o.position
}

func (object *Object) RotateX(angle float32) *Object{
    object.rotation = object.rotation.Mul4(mgl32.HomogRotate3DX(angle))
	return object
}

func (object *Object) RotateY(angle float32) *Object{
    object.rotation = object.rotation.Mul4(mgl32.HomogRotate3DY(angle))
	return object
}

func (object *Object) RotateZ(angle float32) *Object{
    object.rotation = object.rotation.Mul4(mgl32.HomogRotate3DZ(angle))
	return object
}

func (object *Object) SetScale(scale float32) {
	object.setScale(scale)
}

func (object *Object) setScale(scale float32) {
	object.scale = mgl32.Vec3{scale, scale, scale}
}

