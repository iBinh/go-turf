package smooth

import (
	"fmt"

	"github.com/ibinh/turf-go/geojson"
)

func PolygonSmooth(poly any, iterations int) (*geojson.Feature, error) {
	if iterations < 1 {
		iterations = 1
	}
	if iterations > 5 {
		iterations = 5
	}

	g, err := geojson.GetGeometry(poly)
	if err != nil {
		return nil, fmt.Errorf("polygonSmooth: %w", err)
	}

	switch v := g.(type) {
	case *geojson.Polygon:
		rings := smoothRings(v.Coordinates, iterations)
		return geojson.NewFeature(geojson.NewPolygon(rings), nil), nil

	case *geojson.MultiPolygon:
		result := make([][][]geojson.Position, len(v.Coordinates))
		for i, polyCoords := range v.Coordinates {
			result[i] = smoothRings(polyCoords, iterations)
		}
		return geojson.NewFeature(geojson.NewMultiPolygon(result), nil), nil

	default:
		return nil, fmt.Errorf("polygonSmooth: expected Polygon or MultiPolygon, got %s", g.Type())
	}
}

func smoothRings(rings [][]geojson.Position, iterations int) [][]geojson.Position {
	result := make([][]geojson.Position, len(rings))
	for i, ring := range rings {
		result[i] = smoothRing(ring, iterations)
	}
	return result
}

func smoothRing(ring []geojson.Position, iterations int) []geojson.Position {
	current := ring
	for iter := 0; iter < iterations; iter++ {
		if len(current) < 3 {
			break
		}
		current = chaikinStep(current)
	}
	return current
}

func chaikinStep(ring []geojson.Position) []geojson.Position {
	n := len(ring)
	if n < 3 {
		return ring
	}

	closed := ring[0][0] == ring[n-1][0] && ring[0][1] == ring[n-1][1]

	var result []geojson.Position

	last := n - 1
	if closed {
		last = n - 1
	} else {
		last = n - 1
	}

	for i := 0; i < last; i++ {
		curr := ring[i]
		next := ring[(i+1)%n]

		qx := curr[0] + 0.25*(next[0]-curr[0])
		qy := curr[1] + 0.25*(next[1]-curr[1])
		rx := curr[0] + 0.75*(next[0]-curr[0])
		ry := curr[1] + 0.75*(next[1]-curr[1])

		result = append(result, geojson.Position{qx, qy})
		result = append(result, geojson.Position{rx, ry})
	}

	if closed {
		if len(result) > 0 {
			result = append(result, result[0])
		}
	}

	return result
}


