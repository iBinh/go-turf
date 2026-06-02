package transform

import (
	"math"
	"github.com/ibinh/turf-go/geojson"
)

type TranslateOptions struct {
	Units string
}

func TransformTranslate(geom any, dx, dy float64, options ...*TranslateOptions) (*geojson.Feature, error) {
	opts := &TranslateOptions{Units: "degrees"}
	if len(options) > 0 && options[0] != nil {
		opts = options[0]
	}

	var lngOffset, latOffset float64

	switch opts.Units {
	case "meters":
		centroid, _ := getPivot(geom)
		lat := centroid[1]
		latOffset = dy / 111319.9
		cosLat := math.Cos(degToRad(lat))
		if cosLat == 0 {
			lngOffset = dx / 111319.9
		} else {
			lngOffset = dx / (111319.9 * cosLat)
		}
	default:
		lngOffset = dx
		latOffset = dy
	}

	translate := func(p geojson.Position) geojson.Position {
		return geojson.Position{
			p[0] + lngOffset,
			p[1] + latOffset,
		}
	}

	result, err := applyToCoords(geom, translate)
	if err != nil {
		return nil, err
	}
	return geojson.NewFeature(result, nil), nil
}
