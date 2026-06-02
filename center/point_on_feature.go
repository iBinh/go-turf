package center

import (
	"fmt"

	"github.com/ibinh/turf-go/geojson"
)

func PointOnFeature(geom any) (*geojson.Feature, error) {
	g, err := geojson.GetGeometry(geom)
	if err != nil {
		return nil, err
	}

	switch geo := g.(type) {
	case *geojson.Point:
		return geojson.NewFeature(geo, nil), nil
	case *geojson.MultiPoint:
		if len(geo.Coordinates) == 0 {
			return nil, fmt.Errorf("empty geometry")
		}
		return geojson.NewFeature(geojson.NewPoint(geo.Coordinates[0]), nil), nil
	case *geojson.LineString:
		if len(geo.Coordinates) < 2 {
			return nil, fmt.Errorf("line too short")
		}
		mid := len(geo.Coordinates) / 2
		return geojson.NewFeature(geojson.NewPoint(geo.Coordinates[mid]), nil), nil
	case *geojson.MultiLineString:
		if len(geo.Coordinates) == 0 || len(geo.Coordinates[0]) < 2 {
			return nil, fmt.Errorf("empty geometry")
		}
		line := geo.Coordinates[0]
		mid := len(line) / 2
		return geojson.NewFeature(geojson.NewPoint(line[mid]), nil), nil
	case *geojson.Polygon:
		return pointOnPolygon(geo)
	case *geojson.MultiPolygon:
		return pointOnPolygon(geojson.NewPolygon(geo.Coordinates[0]))
	default:
		return nil, fmt.Errorf("unsupported geometry type")
	}
}

func pointOnPolygon(poly *geojson.Polygon) (*geojson.Feature, error) {
	var sumX, sumY float64
	count := 0
	for _, ring := range poly.Coordinates {
		for _, p := range ring {
			sumX += p[0]
			sumY += p[1]
			count++
		}
	}
	if count == 0 {
		return nil, fmt.Errorf("empty polygon")
	}
	cx, cy := sumX/float64(count), sumY/float64(count)
	return geojson.NewFeature(geojson.NewPoint([]float64{cx, cy}), nil), nil
}
