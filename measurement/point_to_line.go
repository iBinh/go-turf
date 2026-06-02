package measurement

import (
	"fmt"
	"math"
	"github.com/ibinh/turf-go/geojson"
)

func PointToLineDistance(point, line any, units ...Unit) (float64, error) {
	pt, err := geojson.GetCoord(point)
	if err != nil {
		return 0, fmt.Errorf("point: %w", err)
	}
	coords, err := geojson.GetCoords(line)
	if err != nil {
		return 0, fmt.Errorf("line: %w", err)
	}
	pts, ok := coords.([]geojson.Position)
	if !ok {
		return 0, fmt.Errorf("point-to-line: expected LineString")
	}

	unit := UnitKilometers
	if len(units) > 0 {
		unit = units[0]
	}

	minDist := math.MaxFloat64
	for i := 0; i < len(pts)-1; i++ {
		d := pointSegmentDistance(pt, pts[i], pts[i+1])
		if d < minDist {
			minDist = d
		}
	}

	return metersFromMeters(minDist, unit), nil
}

func pointSegmentDistance(p, a, b geojson.Position) float64 {
	dx := b[0] - a[0]
	dy := b[1] - a[1]
	length2 := dx*dx + dy*dy

	if length2 == 0 {
		return HaversineDistance(p[0], p[1], a[0], a[1])
	}

	t := ((p[0]-a[0])*dx + (p[1]-a[1])*dy) / length2
	t = math.Max(0, math.Min(1, t))

	proj := geojson.Position{a[0] + t*dx, a[1] + t*dy}
	return HaversineDistance(p[0], p[1], proj[0], proj[1])
}

func NearestPointOnLine(line any, point any) (*geojson.Feature, error) {
	pt, err := geojson.GetCoord(point)
	if err != nil {
		return nil, fmt.Errorf("point: %w", err)
	}
	coords, err := geojson.GetCoords(line)
	if err != nil {
		return nil, fmt.Errorf("line: %w", err)
	}
	pts, ok := coords.([]geojson.Position)
	if !ok {
		return nil, fmt.Errorf("nearest-point-on-line: expected LineString")
	}

	minDist := math.MaxFloat64
	var nearest geojson.Position
	nearestIndex := 0

	for i := 0; i < len(pts)-1; i++ {
		a, b := pts[i], pts[i+1]
		dx := b[0] - a[0]
		dy := b[1] - a[1]
		length2 := dx*dx + dy*dy

		var proj geojson.Position
		var d float64

		if length2 == 0 {
			proj = a
			d = HaversineDistance(pt[0], pt[1], a[0], a[1])
		} else {
			t := ((pt[0]-a[0])*dx + (pt[1]-a[1])*dy) / length2
			t = math.Max(0, math.Min(1, t))
			proj = geojson.Position{a[0] + t*dx, a[1] + t*dy}
			d = HaversineDistance(pt[0], pt[1], proj[0], proj[1])
		}

		if d < minDist {
			minDist = d
			nearest = proj
			nearestIndex = i
		}
	}

	f := geojson.NewFeature(
		geojson.NewPoint(nearest),
		map[string]any{
			"dist":     minDist,
			"index":    nearestIndex,
			"location": minDist,
		},
	)
	return f, nil
}

func NearestPoint(targetPoint any, points any) (*geojson.Feature, error) {
	target, err := geojson.GetCoord(targetPoint)
	if err != nil {
		return nil, fmt.Errorf("target: %w", err)
	}
	features, err := extractFeatures(points)
	if err != nil {
		return nil, err
	}

	var nearest *geojson.Feature
	minDist := math.MaxFloat64

	for _, f := range features {
		coord, err := geojson.GetCoord(f)
		if err != nil {
			continue
		}
		d := HaversineDistance(target[0], target[1], coord[0], coord[1])
		if d < minDist {
			minDist = d
			nearest = f
		}
	}
	return nearest, nil
}

func NearestPointToLine(line any, points *geojson.FeatureCollection, units ...Unit) (*geojson.Feature, error) {
	if points == nil {
		return nil, fmt.Errorf("points are required")
	}
	unit := UnitKilometers
	if len(units) > 0 {
		unit = units[0]
	}

	var nearest *geojson.Feature
	minDist := math.MaxFloat64

	for _, f := range points.Features {
		d, err := PointToLineDistance(f, line, unit)
		if err != nil {
			continue
		}
		if d < minDist {
			minDist = d
			nearest = f
		}
	}
	if nearest == nil {
		return nil, fmt.Errorf("no nearest point found")
	}

	if nearest.Properties == nil {
		nearest.Properties = map[string]any{}
	}
	nearest.Properties["dist"] = minDist
	return nearest, nil
}

func extractFeatures(obj any) ([]*geojson.Feature, error) {
	switch v := obj.(type) {
	case *geojson.FeatureCollection:
		return v.Features, nil
	case *geojson.Feature:
		return []*geojson.Feature{v}, nil
	default:
		return nil, fmt.Errorf("expected FeatureCollection or Feature")
	}
}
