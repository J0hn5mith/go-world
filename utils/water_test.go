package go_world_utils

import (
	"go-world/go-world"
	"testing"
)

type Pair struct {
	key, value float32
}

var PositionHeight = []Pair{
	Pair{0.5, 1.0},
	Pair{1.0, 1.5},
	Pair{1.5, 2.0},
	Pair{2.0, 2.5},
	Pair{2.5, 3.0},
	Pair{10.0, 3.0},
}

var PositionVelocity = []Pair{
	Pair{0.0, 0.0},
	Pair{1.0, 1.0},
	Pair{1.5, 1.5},
	Pair{10.0, 0.0},
}

var PositionVelocityAdvect = []Pair{
	Pair{0.0, 0.0},
	Pair{1.0, 0.0},
	Pair{1.5, 0.0},
	Pair{10.0, 0.0},
}

var PositionVelocityAdvect2 = []Pair{
	Pair{0.0, 0.0},
	Pair{1.0, 0.0},
	Pair{1.5, 0.5},
	Pair{2.0, 1.0},
	Pair{3.0, 1.0},
	Pair{10.0, 0.0},
}

var PositionHeightAdvect = []Pair{
	Pair{0.5, 1.0},
	Pair{1.0, 1.0},
	Pair{1.5, 1.0},
	Pair{2.0, 1.5},
	Pair{2.5, 2.0},
	Pair{10.0, 2.0},
}

func TestGetAt(t *testing.T) {
	world := new(go_world.World)
	water := CreateWater(0, 0, 3, 1, 1, world)
	for i, cell := range water.cells {
		cell.velocity = float32(i + 1)
		cell.height = float32(i + 1)
	}

	compareHeight(water, PositionHeight, t)
	compareVelocity(water, PositionVelocity, t)
}

func TestAdvect(t *testing.T) {
	world := new(go_world.World)
	water := CreateWater(0, 0, 3, 1, 1, world)
	for i, cell := range water.cells {
		cell.velocity = 0
		cell.height = float32(i + 1)
	}

	water.advectHeight(1)
	water.advectVelocities(1)
	compareHeight(water, PositionHeight, t)
	compareVelocity(water, PositionVelocityAdvect, t)
	for _, cell := range water.cells {
		cell.velocity = 1
		cell.height = cell.height
	}

	water.advectHeight(1)
	water.advectVelocities(1)
	compareHeight(water, PositionHeightAdvect, t)
	compareVelocity(water, PositionVelocityAdvect2, t)
}

func TestUpdateVelocity(t *testing.T) {
	world := new(go_world.World)
	water := CreateWater(0, 0, 3, 1, 1, world)
	for i, cell := range water.cells {
		cell.velocity = 0
		cell.height = float32(i + 1)
	}
	water.updateHeight(1)
	water.updateVelocity(1)
    velPairs := []Pair{
		Pair{0.0, 0.0},
		Pair{1.0, -9.81},
		Pair{2.0, -9.81},
		Pair{3.0, 0},
		Pair{10.0, 0.0},
	}
	compareVelocity(water, velPairs, t)
}

func TestUpdateHeight(t *testing.T) {
	world := new(go_world.World)
	water := CreateWater(0, 0, 3, 1, 1, world)
	for i, cell := range water.cells {
		cell.velocity = 1
		cell.height = float32(i + 1)
	}
	water.updateHeight(1)
    heightPairs := []Pair{
		Pair{0.0, 0.0},
		Pair{1.0, -9.81},
		Pair{2.0, -9.81},
		Pair{3.0, 0},
		Pair{10.0, 0.0},
	}
	compareHeight(water, heightPairs , t)
}

func compareHeight(water *Water, values []Pair, t *testing.T) {
	for _, entry := range values {
		if water.HeightAt(entry.key) != entry.value {
			t.Error("Height is incorrect", water.HeightAt(entry.key), entry.value, entry.key)
		}
	}
}

func compareVelocity(water *Water, values []Pair, t *testing.T) {
	for _, entry := range values {
		if water.VelocityAt(entry.key) != entry.value {
			t.Error("Velocity is incorrect", water.VelocityAt(entry.key), entry.value, entry.key)
		}
	}
}
