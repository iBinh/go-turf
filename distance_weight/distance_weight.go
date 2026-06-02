package distanceweight

import (
	"fmt"
	"math"

	"github.com/ibinh/turf-go/geojson"
	"github.com/ibinh/turf-go/meta"
)

func DistanceWeight(points *geojson.FeatureCollection, options ...float64) ([][]float64, error) {
	threshold := 0.0
	if len(options) > 0 {
		threshold = options[0]
	}

	if points == nil || len(points.Features) < 2 {
		return nil, fmt.Errorf("at least 2 points required")
	}

	var coords []geojson.Position
	meta.CoordEach(points, func(c geojson.Position, _ int) error {
		coords = append(coords, c)
		return nil
	})

	n := len(coords)
	weights := make([][]float64, n)
	for i := 0; i < n; i++ {
		weights[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			if i == j {
				weights[i][j] = 0
				continue
			}
			dx := coords[i][0] - coords[j][0]
			dy := coords[i][1] - coords[j][1]
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist < 1e-15 {
				weights[i][j] = 0
			} else {
				w := 1.0 / dist
				if threshold > 0 && dist > threshold {
					w = 0
				}
				weights[i][j] = w
			}
		}
	}

	return weights, nil
}
