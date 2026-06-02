package mask

import (
	"fmt"

	"github.com/ibinh/turf-go/geojson"
	"github.com/ibinh/turf-go/boolean"
)

func Mask(outer any, inner any) (*geojson.Feature, error) {
	outerGeom, err := geojson.GetGeometry(outer)
	if err != nil {
		return nil, fmt.Errorf("mask: outer: %w", err)
	}

	var outerRing []geojson.Position
	switch v := outerGeom.(type) {
	case *geojson.Polygon:
		if len(v.Coordinates) == 0 {
			return nil, fmt.Errorf("mask: outer polygon has no rings")
		}
		outerRing = v.Coordinates[0]
	case *geojson.MultiPolygon:
		if len(v.Coordinates) == 0 || len(v.Coordinates[0]) == 0 {
			return nil, fmt.Errorf("mask: outer multipolygon has no rings")
		}
		outerRing = v.Coordinates[0][0]
	default:
		return nil, fmt.Errorf("mask: outer must be Polygon or MultiPolygon, got %s", outerGeom.Type())
	}

	var innerRings [][]geojson.Position
	innerGeom, err := geojson.GetGeometry(inner)
	if err != nil {
		return nil, fmt.Errorf("mask: inner: %w", err)
	}

	switch v := innerGeom.(type) {
	case *geojson.Polygon:
		for _, ring := range v.Coordinates {
			if isRingInside(ring, outerRing) {
				innerRings = append(innerRings, reorientRing(ring, false))
			}
		}
	case *geojson.MultiPolygon:
		for _, poly := range v.Coordinates {
			for _, ring := range poly {
				if isRingInside(ring, outerRing) {
					innerRings = append(innerRings, reorientRing(ring, false))
				}
			}
		}
	default:
		return nil, fmt.Errorf("mask: inner must be Polygon or MultiPolygon, got %s", innerGeom.Type())
	}

	resultRings := make([][]geojson.Position, 0, 1+len(innerRings))
	resultRings = append(resultRings, outerRing)
	resultRings = append(resultRings, innerRings...)

	return geojson.NewFeature(geojson.NewPolygon(resultRings), nil), nil
}

func isRingInside(ring []geojson.Position, outer []geojson.Position) bool {
	if len(ring) == 0 {
		return false
	}
	pt := ring[0]
	inside, _ := boolean.PointInPolygon(
		geojson.NewPoint(pt),
		geojson.NewPolygon([][]geojson.Position{outer}),
	)
	return inside
}

func reorientRing(ring []geojson.Position, clockwise bool) []geojson.Position {
	if len(ring) < 3 {
		return ring
	}
	isCW := isClockwise(ring)
	if isCW != clockwise {
		result := make([]geojson.Position, len(ring))
		for i, p := range ring {
			result[len(ring)-1-i] = p
		}
		return result
	}
	result := make([]geojson.Position, len(ring))
	copy(result, ring)
	return result
}

func isClockwise(ring []geojson.Position) bool {
	area := 0.0
	for i := 0; i < len(ring)-1; i++ {
		area += (ring[i+1][0] - ring[i][0]) * (ring[i+1][1] + ring[i][1])
	}
	return area > 0
}
