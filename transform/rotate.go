package transform

import (
	"math"
	"github.com/ibinh/turf-go/geojson"
)

func Clone(obj any) (*geojson.Feature, error) {
	// Works directly with GeoJSON types via JSON round-trip
	return nil, nil // Not implemented as a generic - use json.Marshal/Unmarshal directly
}

func TransformRotate(geom any, angle float64, pivot ...geojson.Position) (*geojson.Feature, error) {
	var pivotPt geojson.Position
	if len(pivot) > 0 {
		pivotPt = pivot[0]
	} else {
		var err error
		pivotPt, err = getPivot(geom)
		if err != nil {
			return nil, err
		}
	}

	rad := degToRad(angle)
	cos := math.Cos(rad)
	sin := math.Sin(rad)

	rotate := func(p geojson.Position) geojson.Position {
		dx := p[0] - pivotPt[0]
		dy := p[1] - pivotPt[1]
		return geojson.Position{
			pivotPt[0] + dx*cos - dy*sin,
			pivotPt[1] + dx*sin + dy*cos,
		}
	}

	result, err := applyToCoords(geom, rotate)
	if err != nil {
		return nil, err
	}
	return geojson.NewFeature(result, nil), nil
}
