package transform

import (
	"math"
	"github.com/ibinh/turf-go/geojson"
)

const mercatorOffset = 20037508.34

func ToMercator(geom any) (*geojson.Feature, error) {
	toMerc := func(p geojson.Position) geojson.Position {
		lon := p[0]
		lat := p[1]
		x := lon * mercatorOffset / 180
		y := math.Log(math.Tan(math.Pi/4+lat*math.Pi/360)) * mercatorOffset / math.Pi
		result := geojson.Position{x, y}
		if len(p) > 2 {
			result = append(result, p[2:]...)
		}
		return result
	}

	result, err := applyToCoords(geom, toMerc)
	if err != nil {
		return nil, err
	}
	return geojson.NewFeature(result, nil), nil
}

func ToWGS84(geom any) (*geojson.Feature, error) {
	toWGS := func(p geojson.Position) geojson.Position {
		x := p[0]
		y := p[1]
		lon := x * 180 / mercatorOffset
		lat := (math.Atan(math.Exp(y*math.Pi/mercatorOffset))*360/math.Pi - 90)
		result := geojson.Position{lon, lat}
		if len(p) > 2 {
			result = append(result, p[2:]...)
		}
		return result
	}

	result, err := applyToCoords(geom, toWGS)
	if err != nil {
		return nil, err
	}
	return geojson.NewFeature(result, nil), nil
}
