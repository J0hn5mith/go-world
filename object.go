package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Object struct{
    geometry *Geometry
    position mgl32.Vec3
    rotation mgl32.Vec3
    scale mgl32.Vec3

    model mgl32.Mat4
    modelUniform int32
	angle float64
    vao uint32
    vbo uint32
}

func NewObject(geometry *Geometry) *Object {
    object := new(Object)
    object.geometry = geometry

    object.position = mgl32.Vec3{0,0,0}
    object.rotation = mgl32.Vec3{0,0,0}
    object.scale = mgl32.Vec3{1,1,1}

    object.angle = 0
    return object
}

func (object *Object) setScale(scale float32) {
    object.scale = mgl32.Vec3{scale,scale,scale}
}

func (object *Object)configure(program uint32){
    // Create the glsl uniforms for communicating with the shader
	// Configure the vertex data
	//object.modelUniform = gl.GetUniformLocation(program, gl.Str("model\x00"))
	//gl.UniformMatrix4fv(object.modelUniform, 1, false, &object.model[0])
	//gl.BindFragDataLocation(program, 171, gl.Str("outputColor\x00"))

	gl.GenVertexArrays(1, &object.vao)
	gl.BindVertexArray(object.vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(object.geometry.vertices)*6, gl.Ptr(object.geometry.vertices), gl.STATIC_DRAW)
    object.vbo = vbo

    vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("position\x00")))
    gl.EnableVertexAttribArray(vertAttrib)
    gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	// Configure global settings
	//gl.Enable(gl.DEPTH_TEST)
	//gl.DepthFunc(gl.LESS)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
}
