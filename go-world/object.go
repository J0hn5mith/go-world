package go_world

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Object struct {
	geometry             *Geometry
	position             mgl32.Vec3
	rotation             mgl32.Mat4
	scale                mgl32.Vec3
	transformationBuffer mgl32.Mat4
	dirty                bool
}

func NewObject(geometry *Geometry) *Object {
	object := new(Object)
	object.geometry = geometry
	object.position = mgl32.Vec3{0, 0, 0}
	object.rotation = mgl32.Ident4()
	object.scale = mgl32.Vec3{1, 1, 1}
	object.transformationBuffer = mgl32.Ident4()
	object.dirty = true

	return object
}

func (object Object) Geometry() *Geometry {
	return object.geometry
}

func (o *Object) SetPosition(x, y, z float32) *Object {
	o.position = mgl32.Vec3{x, y, z}
    o.dirty = true
	return o
}

func (o Object) Position() mgl32.Vec3 {
	return o.position
}

func (object *Object) RotateX(angle float32) *Object {
	object.rotation = object.rotation.Mul4(mgl32.HomogRotate3DX(angle))
    object.dirty = true
	return object
}

func (object *Object) RotateY(angle float32) *Object {
	object.rotation = object.rotation.Mul4(mgl32.HomogRotate3DY(angle))
    object.dirty = true
	return object
}

func (object *Object) RotateZ(angle float32) *Object {
	object.rotation = object.rotation.Mul4(mgl32.HomogRotate3DZ(angle))
    object.dirty = true
	return object
}

func (object *Object) Rotation() mgl32.Mat4 {
	return object.rotation
}

func (object *Object) SetScale(scale float32) {
	object.scale = mgl32.Vec3{scale, scale, scale}
    object.dirty = true
}

func (object *Object) TransformationMatrix() mgl32.Mat4 {
	if object.dirty {
		mat := mgl32.Ident4()
		trans := mgl32.Translate3D(
			object.position[0],
			object.position[1],
			object.position[2],
		)
		scale := mgl32.Scale3D(
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
