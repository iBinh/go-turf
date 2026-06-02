package convert

import (
	"fmt"

	"github.com/ibinh/turf-go/geojson"
)

func PolygonToLine(poly any) (*geojson.Feature, error) {
	g, err := geojson.GetGeometry(poly)
	if err != nil {
		return nil, fmt.Errorf("polygonToLine: %w", err)
	}

	switch v := g.(type) {
	case *geojson.Polygon:
		if len(v.Coordinates) == 0 {
			return nil, fmt.Errorf("polygonToLine: polygon has no rings")
		}
		exterior := make([]geojson.Position, len(v.Coordinates[0]))
		copy(exterior, v.Coordinates[0])
		return geojson.NewFeature(geojson.NewLineString(exterior), nil), nil

	case *geojson.MultiPolygon:
		var lines [][]geojson.Position
		for _, polyCoords := range v.Coordinates {
			if len(polyCoords) == 0 {
				continue
			}
			exterior := make([]geojson.Position, len(polyCoords[0]))
			copy(exterior, polyCoords[0])
			lines = append(lines, exterior)
		}
		if len(lines) == 0 {
			return nil, fmt.Errorf("polygonToLine: no rings found")
		}
		if len(lines) == 1 {
			return geojson.NewFeature(geojson.NewLineString(lines[0]), nil), nil
		}
		return geojson.NewFeature(geojson.NewMultiLineString(lines), nil), nil

	default:
		return nil, fmt.Errorf("polygonToLine: expected Polygon or MultiPolygon, got %s", g.Type())
	}
}

func LineToPolygon(line any) (*geojson.Feature, error) {
	g, err := geojson.GetGeometry(line)
	if err != nil {
		return nil, fmt.Errorf("lineToPolygon: %w", err)
	}

	switch v := g.(type) {
	case *geojson.LineString:
		if len(v.Coordinates) < 3 {
			return nil, fmt.Errorf("lineToPolygon: line must have at least 3 coordinates")
		}
		ring := closeRing(v.Coordinates)
		return geojson.NewFeature(geojson.NewPolygon([][]geojson.Position{ring}), nil), nil

	case *geojson.MultiLineString:
		if len(v.Coordinates) == 0 {
			return nil, fmt.Errorf("lineToPolygon: no lines provided")
		}
		rings := make([][]geojson.Position, len(v.Coordinates))
		for i, line := range v.Coordinates {
			if len(line) < 3 {
				return nil, fmt.Errorf("lineToPolygon: line %d must have at least 3 coordinates", i)
			}
			rings[i] = closeRing(line)
		}
		return geojson.NewFeature(geojson.NewPolygon(rings), nil), nil

	default:
		return nil, fmt.Errorf("lineToPolygon: expected LineString or MultiLineString, got %s", g.Type())
	}
}

func closeRing(coords []geojson.Position) []geojson.Position {
	if len(coords) == 0 {
		return coords
	}
	first := coords[0]
	last := coords[len(coords)-1]
	if first[0] == last[0] && first[1] == last[1] {
		result := make([]geojson.Position, len(coords))
		copy(result, coords)
		return result
	}
	result := make([]geojson.Position, len(coords)+1)
	copy(result, coords)
	result[len(coords)] = first
	return result
}
