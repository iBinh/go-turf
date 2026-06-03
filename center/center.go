package center

import (
	"github.com/ibinh/turf-go/bbox"
	"github.com/ibinh/turf-go/geojson"
	"github.com/ibinh/turf-go/meta"
)

func Center(obj any) (*geojson.Feature, error) {
	bb, err := bbox.BBox(obj)
	if err != nil {
		return nil, err
	}
	return geojson.NewFeature(
		geojson.NewPoint(geojson.Position{(bb[0] + bb[2]) / 2, (bb[1] + bb[3]) / 2}),
		nil,
	), nil
}

func Centroid(obj any) (*geojson.Feature, error) {
	var totalLng, totalLat float64
	count := 0
	err := meta.CoordEach(obj, func(coord geojson.Position, _ int) error {
		totalLng += coord[0]
		totalLat += coord[1]
		count++
		return nil
	})
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, nil
	}
	return geojson.NewFeature(
		geojson.NewPoint(geojson.Position{totalLng / float64(count), totalLat / float64(count)}),
		nil,
	), nil
}

func CenterMean(obj any, properties map[string]any, weight ...string) (*geojson.Feature, error) {
	var totalLng, totalLat, totalWeight float64

	err := meta.FeatureEach(obj, func(f *geojson.Feature, _ int) error {
		coord, err := geojson.GetCoord(f)
		if err != nil {
			return nil
		}
		w := 1.0
		if len(weight) > 0 {
			if v, ok := f.Properties[weight[0]].(float64); ok {
				w = v
			}
		}
		totalLng += coord[0] * w
		totalLat += coord[1] * w
		totalWeight += w
		return nil
	})
	if err != nil {
		return nil, err
	}
	if totalWeight == 0 {
		return nil, nil
	}
	return geojson.NewFeature(
		geojson.NewPoint(geojson.Position{totalLng / totalWeight, totalLat / totalWeight}),
		properties,
	), nil
}

func CenterMedian(obj any, properties map[string]any, weight ...string) (*geojson.Feature, error) {
	type weightedPoint struct {
		coord  geojson.Position
		weight float64
	}

	var pts []weightedPoint
	err := meta.FeatureEach(obj, func(f *geojson.Feature, _ int) error {
		coord, err := geojson.GetCoord(f)
		if err != nil {
			return nil
		}
		w := 1.0
		if len(weight) > 0 {
			if v, ok := f.Properties[weight[0]].(float64); ok {
				w = v
			}
		}
		pts = append(pts, weightedPoint{coord, w})
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(pts) == 0 {
		return geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), properties), nil
	}

	var totalWeight float64
	var medianLng, medianLat float64
	for _, p := range pts {
		totalWeight += p.weight
	}

	cumulative := 0.0
	halfWeight := totalWeight / 2

	for _, p := range pts {
		cumulative += p.weight
		if cumulative >= halfWeight {
			medianLng = p.coord[0]
			medianLat = p.coord[1]
			break
		}
	}

	return geojson.NewFeature(
		geojson.NewPoint(geojson.Position{medianLng, medianLat}),
		properties,
	), nil
}

func CenterOfMass(obj any) (*geojson.Feature, error) {
	var totalX, totalY float64
	vertices := 0

	err := meta.CoordEach(obj, func(coord geojson.Position, _ int) error {
		totalX += coord[0]
		totalY += coord[1]
		vertices++
		return nil
	})
	if err != nil {
		return nil, err
	}
	if vertices == 0 {
		return nil, nil
	}

	return geojson.NewFeature(
		geojson.NewPoint(geojson.Position{totalX / float64(vertices), totalY / float64(vertices)}),
		nil,
	), nil
}

type Position = geojson.Position


