package go_world_utils

import (
	"go-world/go-world"
	"math"
)

var INITIAL_VELOCITY float32 = 0.0

type Water struct {
	cells                []*WaterCell
	width, height, depth float32
	cellWidth            float32
}

func CreateWater(x, y, width, height, interval float32, world *go_world.World) *Water {
	water := new(Water)
	water.width, water.height, water.depth = width, height, 3.0

	numCells := int(width / interval)
	cellWidth := width / float32(numCells)
	water.cellWidth = cellWidth
	var i int
	for i = 0; i < numCells; i++ {
		cell := CreateWaterCell(cellWidth, water)
		water.cells = append(water.cells, cell)
		height := float32(i)*4.0/float32(numCells) + water.height*0.5
		cell.SetHeight(height)
	}

	return water
}

func (water *Water) SetHeights(heights []float32) {
	for i, cell := range water.cells {
		cell.SetHeight(heights[i])
	}
}

func (water *Water) Update(timeDelta float32) {
    //water.advectHeight(timeDelta)
    //water.advectVelocities(timeDelta)

	water.updateHeight(timeDelta)
	water.updateVelocity(timeDelta)
}

func (water *Water) HeightAt(position float32) float32 {
	return water.interpolateHeight(position - water.cellWidth*0.5)
}

func (water *Water) VelocityAt(position float32) float32 {
	return water.interpolateVelocity(position - water.cellWidth)
}

func (water *Water) GetCell(index int32) *WaterCell {
	var cell *WaterCell
	if index < 0 {
		cell = new(WaterCell)
		cell.velocity = 0
		cell.height = water.cells[0].height
	} else if index >= int32(len(water.cells)) {
		cell = new(WaterCell)
		cell.velocity = 0
		cell.height = water.cells[len(water.cells)-1].height
	} else {
		cell = water.cells[index]
	}
	return cell
}

func (water *Water) advectHeight(timeDelta float32) {
	advectHeights := []float32{}
	for i, _ := range water.cells {
		position := (float32(i) + 0.5) * water.cellWidth
		velocity := water.VelocityAt(position)
		prevPosition := position - velocity*timeDelta
		newHeight := water.HeightAt(prevPosition)
		advectHeights = append(advectHeights, newHeight)
	}

	for i, cell := range water.cells {
		cell.SetHeight(advectHeights[i])
	}
}

func (water *Water) advectVelocities(timeDelta float32) {
	advectVelocities := []float32{}
	for i, _ := range water.cells {
		position := (float32(i) + 1) * water.cellWidth
		velocity := water.VelocityAt(position)
		prevPosition := position - timeDelta*velocity
		newVelocity := water.VelocityAt(prevPosition)
		advectVelocities = append(advectVelocities, newVelocity)
	}
	for i, cell := range water.cells {
		cell.velocity = advectVelocities[i]
	}
}

func (water *Water) cellsForInterpolation(position float32) (*WaterCell, *WaterCell, float32) {
	position = position / water.cellWidth
	positionBefore := int32(math.Floor(float64(position)))
	positionNext := int32(math.Ceil(float64(position)))
	offset := position - float32(positionBefore)
	prevCell := water.GetCell(positionBefore)
	nextCell := water.GetCell(positionNext)
	return prevCell, nextCell, offset
}

func (water *Water) interpolateHeight(position float32) float32 {
	prevCell, nextCell, offset := water.cellsForInterpolation(position)
	newHeight := prevCell.height*offset + nextCell.height*(1-offset)
	return newHeight
}

func (water *Water) interpolateVelocity(position float32) float32 {
	prevCell, nextCell, offset := water.cellsForInterpolation(position)
	newVelocity := prevCell.velocity*offset + nextCell.velocity*(1-offset)
	return newVelocity
}

func (water *Water) updateHeight(timeDelta float32) {
	for i, cell := range water.cells {
			nextCell := water.GetCell(int32(i + 1))
			velocity := (nextCell.velocity - cell.velocity) / water.cellWidth
			cell.SetHeight(cell.height - cell.height*velocity*timeDelta)
	}
}

func (water *Water) updateVelocity(timeDelta float32) {
	var prevCell *WaterCell
    prevCell = water.GetCell(-1)
	for _, cell := range water.cells {
		if prevCell != nil {
			difference := (prevCell.height - cell.height) / water.cellWidth
			acceleration := 9.81 * difference * timeDelta
			cell.velocity += acceleration
		}
		prevCell = cell
	}
}

type WaterCell struct {
	water    *Water
	height   float32
	velocity float32
	width    float32
}

func CreateWaterCell(width float32, water *Water) *WaterCell {
	cell := new(WaterCell)
	cell.water = water
	cell.width = width
	cell.velocity = INITIAL_VELOCITY
	return cell
}

func (waterCell *WaterCell) SetHeight(height float32) *WaterCell {
	waterCell.height = height
	return waterCell
}

type WaterObject struct {
	water       *Water
	object      *go_world.Object
	cellObjects []*go_world.Object
	world       *go_world.World
}

func CreateWaterObject(x, y, width, height, interval float32, world *go_world.World) *WaterObject {
	waterObject := new(WaterObject)
	waterObject.world = world
	waterObject.water = CreateWater(x, y, width, height, interval, world)
	waterObject.object = waterObject.world.NewObject(
		CreateWireBoxGeometry(
			waterObject.water.width,
			waterObject.water.height,
			waterObject.water.depth,
		).
			Load(waterObject.world.Program()).
			SetColorRGB(0, 0, 0),
	)
	waterObject.createWaterCellObjects()
	waterObject.updateCellObjects()
	return waterObject
}

func (waterObject *WaterObject) createWaterCellObjects() {
	for i, cell := range waterObject.water.cells {
		cellObject := waterObject.world.NewObject(
			CreateBoxGeometry(
				cell.width,
				waterObject.water.height,
				waterObject.water.depth,
			).
				Load(waterObject.world.Program()).
				SetColorRGB(0, 0, 1),
		)
		cellWidth := waterObject.water.cellWidth
		xPosition := float32(i)*cellWidth + (cellWidth / 2) - waterObject.water.width/2.0
		cellObject.SetPosition(xPosition, 0, 0)
		waterObject.cellObjects = append(waterObject.cellObjects, cellObject)
	}
}

//[>
//Returns all objects realted to the water (also the ones of the cells)
//*/
func (waterObject *WaterObject) Objects() []*go_world.Object {
	objects := []*go_world.Object{waterObject.object}
	for _, object := range waterObject.cellObjects {
		objects = append(objects, object)
	}
	return objects
}

func (waterObject *WaterObject) Update(timeDelta float32) {
	waterObject.water.Update(timeDelta)
	waterObject.updateCellObjects()
}

func (waterObject *WaterObject) updateCellObjects() {
	for i, cell := range waterObject.water.cells {
		object := waterObject.cellObjects[i]
		newVertices := CreateBoxVertices(
			cell.width,
			cell.height,
			cell.water.depth,
		)
		object.Geometry().UpdateVertices(newVertices)
		position := object.Position()
		object.SetPosition(
			position.X(),
			-(cell.water.height-cell.height)/2,
			position.Z(),
		)
	}
}
