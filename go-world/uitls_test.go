package go_world

import (
    "testing"
    "math"
)

type AngleToXY struct {
    angle, x, y float32
}

func TestAngleToCoordinates32(t *testing.T) {
    var x, y float32

    configs := []AngleToXY{
        {0,1,0},
        {math.Pi/2,0,1},
        {math.Pi,-1,0},
        {math.Pi*1.5,0,-0.1},
    }

    for _, config := range configs {
        x, y = AngleToCoordinates32(config.angle)
        if Abs32(x - config.x) > 0.000001 && Abs32(y - config.y)  > 0.0000001{
            t.Error("Incorrect angle to coordinates conversion.", x, y)
        }

    }
}
