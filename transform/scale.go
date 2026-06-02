package transform

import (
	"github.com/ibinh/turf-go/geojson"
)

func TransformScale(geom any, factor float64, pivot ...geojson.Position) (*geojson.Feature, error) {
	return TransformScaleXY(geom, factor, factor, pivot...)
}

func TransformScaleXY(geom any, xFactor, yFactor float64, pivot ...geojson.Position) (*geojson.Feature, error) {
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

	scale := func(p geojson.Position) geojson.Position {
		return geojson.Position{
			pivotPt[0] + (p[0]-pivotPt[0])*xFactor,
			pivotPt[1] + (p[1]-pivotPt[1])*yFactor,
		}
	}

	result, err := applyToCoords(geom, scale)
	if err != nil {
		return nil, err
	}
	return geojson.NewFeature(result, nil), nil
}
