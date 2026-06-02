package tangents

import (
	"fmt"

	"github.com/ibinh/turf-go/geojson"
)

func PolygonTangents(point any, polygon any) (*geojson.FeatureCollection, error) {
	ptCoord, err := geojson.GetCoord(point)
	if err != nil {
		return nil, fmt.Errorf("polygonTangents: point: %w", err)
	}

	polyGeom, err := geojson.GetGeometry(polygon)
	if err != nil {
		return nil, fmt.Errorf("polygonTangents: polygon: %w", err)
	}

	var rings [][]geojson.Position
	switch v := polyGeom.(type) {
	case *geojson.Polygon:
		rings = v.Coordinates
	case *geojson.MultiPolygon:
		for _, poly := range v.Coordinates {
			rings = append(rings, poly[0])
		}
	default:
		return nil, fmt.Errorf("polygonTangents: expected Polygon or MultiPolygon, got %s", polyGeom.Type())
	}

	if len(rings) == 0 {
		return nil, fmt.Errorf("polygonTangents: polygon has no rings")
	}

	exterior := rings[0]
	if len(exterior) < 3 {
		return nil, fmt.Errorf("polygonTangents: polygon must have at least 3 vertices")
	}

	leftIdx := findLeftTangent(ptCoord, exterior)
	rightIdx := findRightTangent(ptCoord, exterior)

	leftTangent := geojson.NewFeature(geojson.NewPoint(exterior[leftIdx]), nil)
	rightTangent := geojson.NewFeature(geojson.NewPoint(exterior[rightIdx]), nil)

	return geojson.NewFeatureCollection([]*geojson.Feature{leftTangent, rightTangent}), nil
}

func cross(o, a, b geojson.Position) float64 {
	return (a[0]-o[0])*(b[1]-o[1]) - (a[1]-o[1])*(b[0]-o[0])
}

func findLeftTangent(pt geojson.Position, ring []geojson.Position) int {
	n := len(ring) - 1
	if n < 2 {
		return 0
	}

	prevLeft := cross(pt, ring[n-1], ring[0]) > 0
	for i := 0; i < n; i++ {
		currLeft := cross(pt, ring[i], ring[i+1]) > 0
		if prevLeft && !currLeft {
			leftIdx := i - 1
			if leftIdx < 0 {
				leftIdx = n - 1
			}
			return refineLeftTangent(pt, ring, leftIdx)
		}
		prevLeft = currLeft
	}

	return 0
}

func refineLeftTangent(pt geojson.Position, ring []geojson.Position, start int) int {
	n := len(ring) - 1
	best := start
	for i := 0; i < n; i++ {
		idx := (start - i + n) % n
		if cross(pt, ring[(idx+1)%n], ring[idx]) <= 0 {
			best = idx
		} else {
			break
		}
	}
	return best
}

func findRightTangent(pt geojson.Position, ring []geojson.Position) int {
	n := len(ring) - 1
	if n < 2 {
		return 0
	}

	prevLeft := cross(pt, ring[n-1], ring[0]) > 0
	for i := 0; i < n; i++ {
		currLeft := cross(pt, ring[i], ring[i+1]) > 0
		if !prevLeft && currLeft {
			rightIdx := i - 1
			if rightIdx < 0 {
				rightIdx = n - 1
			}
			return refineRightTangent(pt, ring, rightIdx)
		}
		prevLeft = currLeft
	}

	return 0
}

func refineRightTangent(pt geojson.Position, ring []geojson.Position, start int) int {
	n := len(ring) - 1
	best := start
	for i := 0; i < n; i++ {
		idx := (start + i) % n
		if cross(pt, ring[idx], ring[(idx+1)%n]) <= 0 {
			best = idx
		} else {
			break
		}
	}
	return best
}


