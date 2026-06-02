package moranindex

import (
	"fmt"
	"math"

	"github.com/ibinh/turf-go/geojson"
	"github.com/ibinh/turf-go/meta"
)

func MoranIndex(points *geojson.FeatureCollection, property string) (float64, error) {
	if points == nil || len(points.Features) < 3 {
		return 0, fmt.Errorf("at least 3 points required")
	}

	var vals []float64
	var coords []geojson.Position

	meta.FeatureEach(points, func(f *geojson.Feature, _ int) error {
		v, ok := f.Properties[property].(float64)
		if !ok {
			return nil
		}
		coord, err := geojson.GetCoord(f)
		if err != nil {
			return nil
		}
		vals = append(vals, v)
		coords = append(coords, coord)
		return nil
	})

	if len(vals) < 3 {
		return 0, fmt.Errorf("not enough valid features with property %q", property)
	}

	n := float64(len(vals))

	var sum float64
	for _, v := range vals {
		sum += v
	}
	mean := sum / n

	deviations := make([]float64, len(vals))
	var sumSq float64
	for i, v := range vals {
		d := v - mean
		deviations[i] = d
		sumSq += d * d
	}

	if sumSq < 1e-15 {
		return 0, nil
	}

	var num, den float64
	for i := 0; i < len(vals); i++ {
		for j := 0; j < len(vals); j++ {
			if i == j {
				continue
			}
			dx := coords[i][0] - coords[j][0]
			dy := coords[i][1] - coords[j][1]
			dist := math.Sqrt(dx*dx + dy*dy)
			w := 1.0
			if dist > 1e-15 {
				w = 1.0 / dist
			}
			num += w * deviations[i] * deviations[j]
		}
	}
	den = sumSq / n

	return num / den / n, nil
}
