package go_world

import (
	mgl "github.com/go-gl/mathgl/mgl64"
)

type Object struct {
	geometry             *Geometry
	position             mgl.Vec3
	rotation             mgl.Mat4
	scale                mgl.Vec3
	transformationBuffer mgl.Mat4
	dirty                bool
}

func NewObject(geometry *Geometry) *Object {
	object := new(Object)
	object.geometry = geometry
	object.position = mgl.Vec3{0, 0, 0}
	object.rotation = mgl.Ident4()
	object.scale = mgl.Vec3{1, 1, 1}
	object.transformationBuffer = mgl.Ident4()
	object.dirty = true

	return object
}

func (object Object) Geometry() *Geometry {
	return object.geometry
}

func (o *Object) SetPosition(position mgl.Vec3) *Object {
	o.position = position
	o.dirty = true
	return o
}

func (o *Object) Position() mgl.Vec3 {
	return o.position
}

func (o *Object) ShiftPosition(shift mgl.Vec3) *Object {
    o.position = o.Position().Add(shift)
    return o
}

func (object *Object) RotateX(angle float64) *Object {
	object.rotation = object.rotation.Mul4(mgl.HomogRotate3DX(angle))
	object.dirty = true
	return object
}

func (object *Object) RotateY(angle float64) *Object {
	object.rotation = object.rotation.Mul4(mgl.HomogRotate3DY(angle))
	object.dirty = true
	return object
}

func (object *Object) RotateZ(angle float64) *Object {
	object.rotation = object.rotation.Mul4(mgl.HomogRotate3DZ(angle))
	object.dirty = true
	return object
}

func (object *Object) Rotation() mgl.Mat4 {
	return object.rotation
}

func (object *Object) SetScale(scale float64) {
	object.scale = mgl.Vec3{scale, scale, scale}
	object.dirty = true
}

func (object *Object) Scale() mgl.Vec3{
	return object.scale
}

func (object *Object) TransformationMatrix() mgl.Mat4 {
	if object.dirty {
		mat := mgl.Ident4()
		trans := mgl.Translate3D(
			object.position[0],
			object.position[1],
			object.position[2],
		)
		scale := mgl.Scale3D(
			object.scale[0],
			object.scale[1],
			object.scale[2],
		)
		object.transformationBuffer = mat.Mul4(scale).Mul4(trans).Mul4(
			object.rotation,
		)
		object.dirty = false
	}
	return object.transformationBuffer
}
