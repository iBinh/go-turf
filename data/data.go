package data

import (
	"fmt"

	"github.com/ibinh/turf-go/boolean"
	"github.com/ibinh/turf-go/geojson"
)

func Tag(points *geojson.FeatureCollection, polygons *geojson.FeatureCollection, field, outField string) (*geojson.FeatureCollection, error) {
	if points == nil || polygons == nil {
		return nil, fmt.Errorf("points and polygons are required")
	}

	result := make([]*geojson.Feature, 0, len(points.Features))

	for _, pt := range points.Features {
		coord, err := geojson.GetCoord(pt)
		if err != nil {
			continue
		}
		props := make(map[string]any)
		for k, v := range pt.Properties {
			props[k] = v
		}

		point := geojson.NewPoint(coord)
		for _, poly := range polygons.Features {
			inside, err := boolean.PointInPolygon(point, poly)
			if err != nil {
				continue
			}
			if inside {
				if val, ok := poly.Properties[field]; ok {
					props[outField] = val
				}
				break
			}
		}

		result = append(result, geojson.NewFeature(point, props))
	}

	return geojson.NewFeatureCollection(result), nil
}

func Collect(polygons *geojson.FeatureCollection, points *geojson.FeatureCollection, inField, outField string) (*geojson.FeatureCollection, error) {
	if polygons == nil || points == nil {
		return nil, fmt.Errorf("polygons and points are required")
	}

	result := make([]*geojson.Feature, 0, len(polygons.Features))

	for _, poly := range polygons.Features {
		props := make(map[string]any)
		for k, v := range poly.Properties {
			props[k] = v
		}

		var collected []any
		for _, pt := range points.Features {
			inside, err := boolean.PointInPolygon(pt, poly)
			if err != nil {
				continue
			}
			if inside {
				if val, ok := pt.Properties[inField]; ok {
					collected = append(collected, val)
				}
			}
		}
		props[outField] = collected

		geom, err := geojson.GetGeometry(poly)
		if err != nil {
			continue
		}
		result = append(result, geojson.NewFeature(geom, props))
	}

	return geojson.NewFeatureCollection(result), nil
}
