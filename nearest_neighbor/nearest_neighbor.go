package nearestneighbor

import (
	"fmt"
	"math"

	"github.com/ibinh/turf-go/geojson"
	"github.com/ibinh/turf-go/meta"
)

type NearestNeighborResult struct {
	Observed float64
	Expected float64
	R        float64
	Z        float64
	P        float64
}

func NearestNeighborAnalysis(points *geojson.FeatureCollection) (*NearestNeighborResult, error) {
	if points == nil || len(points.Features) < 2 {
		return nil, fmt.Errorf("at least 2 points required")
	}

	var coords []geojson.Position
	meta.CoordEach(points, func(c geojson.Position, _ int) error {
		coords = append(coords, c)
		return nil
	})

	if len(coords) < 2 {
		return nil, fmt.Errorf("no valid coordinates")
	}

	n := float64(len(coords))

	minX, minY := coords[0][0], coords[0][1]
	maxX, maxY := minX, minY
	for _, c := range coords {
		if c[0] < minX {
			minX = c[0]
		}
		if c[0] > maxX {
			maxX = c[0]
		}
		if c[1] < minY {
			minY = c[1]
		}
		if c[1] > maxY {
			maxY = c[1]
		}
	}
	area := (maxX - minX) * (maxY - minY)
	if area <= 0 {
		area = 1
	}

	var sumDist float64
	for i, c := range coords {
		minDist := math.MaxFloat64
		for j, d := range coords {
			if i == j {
				continue
			}
			dx := c[0] - d[0]
			dy := c[1] - d[1]
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist < minDist {
				minDist = dist
			}
		}
		if minDist < math.MaxFloat64 {
			sumDist += minDist
		}
	}

	observed := sumDist / n
	expected := 0.5 * math.Sqrt(area/n)
	r := observed / expected

	se := 0.26136 / math.Sqrt(n*n/area)
	z := (observed - expected) / se
	if se < 1e-15 {
		z = 0
	}

	p := 0.5 * math.Erfc(math.Abs(z)/math.Sqrt2)

	return &NearestNeighborResult{
		Observed: observed,
		Expected: expected,
		R:        r,
		Z:        z,
		P:        p,
	}, nil
}
