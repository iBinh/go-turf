package transform

import (
	"math"
	"github.com/ibinh/turf-go/geojson"
	"github.com/ibinh/turf-go/meta"
)

func Flip(geom any) (*geojson.Feature, error) {
	flipFn := func(p geojson.Position) geojson.Position {
		result := make(geojson.Position, len(p))
		if len(p) >= 2 {
			result[0] = p[1]
			result[1] = p[0]
		}
		for i := 2; i < len(p); i++ {
			result[i] = p[i]
		}
		return result
	}
	result, err := applyToCoords(geom, flipFn)
	if err != nil {
		return nil, err
	}
	return geojson.NewFeature(result, nil), nil
}

func Truncate(geom any, precision int, coordinates ...int) (*geojson.Feature, error) {
	maxCoords := 3
	if len(coordinates) > 0 {
		maxCoords = coordinates[0]
	}

	factor := math.Pow(10, float64(precision))

	truncate := func(p geojson.Position) geojson.Position {
		n := len(p)
		if n > maxCoords {
			n = maxCoords
		}
		result := make(geojson.Position, n)
		for i := 0; i < n; i++ {
			result[i] = math.Round(p[i]*factor) / factor
		}
		return result
	}

	result, err := applyToCoords(geom, truncate)
	if err != nil {
		return nil, err
	}
	return geojson.NewFeature(result, nil), nil
}

func CleanCoords(geom any) (*geojson.Feature, error) {
	var resultGeom geojson.Geometry

	err := meta.GeomEach(geom, func(g geojson.Geometry, _ int) error {
		resultGeom = dedupeGeometry(g)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return geojson.NewFeature(resultGeom, nil), nil
}

func dedupeGeometry(g geojson.Geometry) geojson.Geometry {
	switch v := g.(type) {
	case *geojson.Point:
		return v
	case *geojson.MultiPoint:
		return geojson.NewMultiPoint(dedupeRing(v.Coordinates))
	case *geojson.LineString:
		return geojson.NewLineString(dedupeRing(v.Coordinates))
	case *geojson.MultiLineString:
		lines := make([][]geojson.Position, len(v.Coordinates))
		for i, line := range v.Coordinates {
			lines[i] = dedupeRing(line)
		}
		return geojson.NewMultiLineString(lines)
	case *geojson.Polygon:
		rings := make([][]geojson.Position, len(v.Coordinates))
		for i, ring := range v.Coordinates {
			rings[i] = dedupeRing(ring)
		}
		return geojson.NewPolygon(rings)
	case *geojson.MultiPolygon:
		polygons := make([][][]geojson.Position, len(v.Coordinates))
		for i, poly := range v.Coordinates {
			rings := make([][]geojson.Position, len(poly))
			for j, ring := range poly {
				rings[j] = dedupeRing(ring)
			}
			polygons[i] = rings
		}
		return geojson.NewMultiPolygon(polygons)
	default:
		return g
	}
}

func dedupeRing(ring []geojson.Position) []geojson.Position {
	if len(ring) <= 1 {
		return ring
	}
	result := []geojson.Position{ring[0]}
	for i := 1; i < len(ring); i++ {
		last := result[len(result)-1]
		if ring[i][0] != last[0] || ring[i][1] != last[1] {
			result = append(result, ring[i])
		}
	}
	return result
}

func Rewind(geom any, reversed ...bool) (*geojson.Feature, error) {
	reverse := false
	if len(reversed) > 0 {
		reverse = reversed[0]
	}

	var result *geojson.Feature

	err := meta.GeomEach(geom, func(g geojson.Geometry, _ int) error {
		rewound, err := rewindGeometry(g, reverse)
		if err != nil {
			return err
		}
		result = geojson.NewFeature(rewound, nil)
		return nil
	})
	return result, err
}

func rewindGeometry(g geojson.Geometry, reverse bool) (geojson.Geometry, error) {
	switch v := g.(type) {
	case *geojson.Polygon:
		rings := make([][]geojson.Position, len(v.Coordinates))
		for i, ring := range v.Coordinates {
			if i == 0 {
				if isCWRing(ring) != reverse {
					rings[i] = ring
				} else {
					rings[i] = reverseRing(ring)
				}
			} else {
				if isCWRing(ring) == reverse {
					rings[i] = ring
				} else {
					rings[i] = reverseRing(ring)
				}
			}
		}
		return geojson.NewPolygon(rings), nil
	case *geojson.MultiPolygon:
		polygons := make([][][]geojson.Position, len(v.Coordinates))
		for i, poly := range v.Coordinates {
			rings := make([][]geojson.Position, len(poly))
			for j, ring := range poly {
				if j == 0 {
					if isCWRing(ring) != reverse {
						rings[j] = ring
					} else {
						rings[j] = reverseRing(ring)
					}
				} else {
					if isCWRing(ring) == reverse {
						rings[j] = ring
					} else {
						rings[j] = reverseRing(ring)
					}
				}
			}
			polygons[i] = rings
		}
		return geojson.NewMultiPolygon(polygons), nil
	default:
		return g, nil
	}
}

func isCWRing(ring []geojson.Position) bool {
	area := 0.0
	for i := 0; i < len(ring)-1; i++ {
		area += (ring[i+1][0] - ring[i][0]) * (ring[i+1][1] + ring[i][1])
	}
	return area > 0
}

func reverseRing(ring []geojson.Position) []geojson.Position {
	if len(ring) <= 2 {
		return ring
	}
	reversed := make([]geojson.Position, len(ring))
	reversed[0] = ring[0]
	for i := 1; i < len(ring)-1; i++ {
		reversed[i] = ring[len(ring)-1-i]
	}
	reversed[len(ring)-1] = ring[len(ring)-1]
	return reversed
}
