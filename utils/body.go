package go_world_utils

import (
	mgl "github.com/go-gl/mathgl/mgl64"
	"go-world/go-world"
	"go-world/physics"
)

func CreateBoxBody(physics *go_world_physics.Physics, boxDimension mgl.Vec3) *go_world_physics.RigidBody {
	body := go_world_physics.CreateDynamicBody()
	go_world_physics.AddMassParticle(body, boxDimension, 1)
	physics.RegisterBody(body)
	return body
}

func CreateBoxTower(physics *go_world_physics.Physics, world *go_world.World, height float64, origin mgl.Vec3) []*go_world_physics.RigidBody {
	var bodies []*go_world_physics.RigidBody
	towerElementSize := mgl.Vec3{8, 8, 4}
	for i := 0.0; i < height; i++ {
		CreateTowerElement(
			physics,
			world,
			towerElementSize,
			origin.Add(mgl.Vec3{0, 2 + i*towerElementSize.Y(), 0}),
		)
	}

	return bodies
}

func CreateTowerElement(physics *go_world_physics.Physics, world *go_world.World, dimension mgl.Vec3, origin mgl.Vec3) {
	gridSize := 4.0
	unitBox := dimension.Mul(1.0 / gridSize)
	boxSmall := mgl.Vec3{unitBox.X(), unitBox.Y(), dimension.Z()}
	boxBig := mgl.Vec3{gridSize * unitBox.X(), unitBox.Y(), dimension.Z()}
	for i := 0.0; i < gridSize; i++ {
		if i != 0 {
			box(
				physics,
				world, boxSmall,
				origin.Add(mgl.Vec3{-unitBox.X(), (i + 0.5) * unitBox.Y(), 0}),
			)
			box(
				physics,
				world,
				boxSmall,
				origin.Add(mgl.Vec3{unitBox.X(), (i + 0.5) * unitBox.Y(), 0}),
			)
		} else {
			box(physics, world, boxBig, origin.Add(mgl.Vec3{0, (i + 0.5) * unitBox.Y(), 0}))
		}

	}
}

func box(physics *go_world_physics.Physics, world *go_world.World, dimension mgl.Vec3, origin mgl.Vec3) {

	object := world.NewObject(
		CreateBoxGeometry(dimension).
			Load(world.Program()).
			SetColorRGB(0, .5, .5),
	)
	world.Scene.AddObject(object)
	body := physics.RegisterObject(object)
	go_world_physics.AddMassParticle(body, dimension, 1)
	body.SetPosition(origin)
}
