package kinks

import (
	"fmt"
	"math"
	"github.com/ibinh/turf-go/geojson"
)

func Kinks(geom any) (*geojson.FeatureCollection, error) {
	g, err := geojson.GetGeometry(geom)
	if err != nil {
		return nil, err
	}

	var rings [][]geojson.Position
	switch v := g.(type) {
	case *geojson.LineString:
		rings = [][]geojson.Position{v.Coordinates}
	case *geojson.Polygon:
		rings = v.Coordinates
	case *geojson.MultiLineString:
		rings = v.Coordinates
	case *geojson.MultiPolygon:
		for _, poly := range v.Coordinates {
			rings = append(rings, poly...)
		}
	default:
		return nil, fmt.Errorf("kinks: unsupported geometry type %s", g.Type())
	}

	var points []*geojson.Feature
	seen := make(map[string]bool)

	for _, ring := range rings {
		m := len(ring)
		for i := 0; i < m-1; i++ {
			a, b := ring[i], ring[i+1]
			for j := i + 1; j < m-1; j++ {
				c, d := ring[j], ring[j+1]
				if shareVertex(a, b, c, d) {
					continue
				}
				if pt, ok := segSegIntersect(a, b, c, d); ok {
					key := fmt.Sprintf("%.10f,%.10f", pt[0], pt[1])
					if !seen[key] {
						seen[key] = true
						points = append(points, geojson.NewFeature(geojson.NewPoint(pt), nil))
					}
				}
			}
		}
	}

	return geojson.NewFeatureCollection(points), nil
}

func shareVertex(a, b, c, d geojson.Position) bool {
	return (equalPos(a, c) || equalPos(a, d) || equalPos(b, c) || equalPos(b, d))
}

func equalPos(a, b geojson.Position) bool {
	return a[0] == b[0] && a[1] == b[1]
}

func segSegIntersect(a, b, c, d geojson.Position) (geojson.Position, bool) {
	den := (b[0]-a[0])*(d[1]-c[1]) - (b[1]-a[1])*(d[0]-c[0])
	if math.Abs(den) < 1e-15 {
		return geojson.Position{}, false
	}
	t := ((c[0]-a[0])*(d[1]-c[1]) - (c[1]-a[1])*(d[0]-c[0])) / den
	u := ((c[0]-a[0])*(b[1]-a[1]) - (c[1]-a[1])*(b[0]-a[0])) / den
	if t < 0 || t > 1 || u < 0 || u > 1 {
		return geojson.Position{}, false
	}
	return geojson.Position{a[0] + t*(b[0]-a[0]), a[1] + t*(b[1]-a[1])}, true
}
