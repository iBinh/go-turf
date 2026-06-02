package measurement

import (
	"fmt"
	"math"

	"github.com/ibinh/turf-go/boolean"
	"github.com/ibinh/turf-go/geojson"
)

func PointToPolygonDistance(point any, polygon any, units ...Unit) (float64, error) {
	pt, err := geojson.GetCoord(point)
	if err != nil {
		return 0, fmt.Errorf("point: %w", err)
	}

	unit := UnitKilometers
	if len(units) > 0 {
		unit = units[0]
	}

	inside, err := boolean.PointInPolygon(point, polygon)
	if err != nil {
		return 0, err
	}
	if inside {
		return 0, nil
	}

	coords, err := geojson.GetCoords(polygon)
	if err != nil {
		return 0, fmt.Errorf("polygon: %w", err)
	}
	rings, ok := coords.([][]geojson.Position)
	if !ok {
		return 0, fmt.Errorf("expected polygon coords")
	}

	minDist := math.MaxFloat64
	for _, ring := range rings {
		for i := 0; i < len(ring)-1; i++ {
			d := pointSegmentDistance(pt, ring[i], ring[i+1])
			if d < minDist {
				minDist = d
			}
		}
	}

	return metersFromMeters(minDist, unit), nil
}
