package measurement

import (
	"math"
	"github.com/ibinh/turf-go/geojson"
)

func mathCos(v float64) float64  { return math.Cos(v) }
func mathSin(v float64) float64  { return math.Sin(v) }
func mathSqrt(v float64) float64 { return math.Sqrt(v) }
func mathAtan2(y, x float64) float64 { return math.Atan2(y, x) }

func Area(geom any) (float64, error) {
	switch g := geom.(type) {
	case *geojson.Feature:
		return areaForGeometry(g.Geometry)
	case geojson.Geometry:
		return areaForGeometry(g)
	case *geojson.FeatureCollection:
		var total float64
		for _, f := range g.Features {
			a, err := areaForGeometry(f.Geometry)
			if err != nil {
				return 0, err
			}
			total += a
		}
		return total, nil
	default:
		return 0, nil
	}
}

func areaForGeometry(geom geojson.Geometry) (float64, error) {
	if geom == nil {
		return 0, nil
	}
	switch g := geom.(type) {
	case *geojson.Polygon:
		return polygonArea(g.Coordinates), nil
	case *geojson.MultiPolygon:
		var total float64
		for _, poly := range g.Coordinates {
			total += polygonArea(poly)
		}
		return total, nil
	default:
		return 0, nil
	}
}

func polygonArea(rings [][]geojson.Position) float64 {
	if len(rings) == 0 {
		return 0
	}
	outer := ringArea(rings[0])
	var innerTotal float64
	for i := 1; i < len(rings); i++ {
		innerTotal += ringArea(rings[i])
	}
	return math.Abs(outer - innerTotal)
}

func ringArea(ring []geojson.Position) float64 {
	if len(ring) < 3 {
		return 0
	}

	var total float64
	for i := 0; i < len(ring)-1; i++ {
		p1 := ring[i]
		p2 := ring[i+1]
		total += (degToRad(p2[0]) - degToRad(p1[0])) *
			(2 + math.Sin(degToRad(p1[1])) + math.Sin(degToRad(p2[1])))
	}

	return total * EarthRadius * EarthRadius / 2
}
