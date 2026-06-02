package standarddeviationalellipse

import (
	"fmt"
	"math"

	"github.com/ibinh/turf-go/geojson"
	"github.com/ibinh/turf-go/meta"
)

type StandardDeviationalEllipseResult struct {
	Center   geojson.Position
	XAxis    float64
	YAxis    float64
	Rotation float64
	Polygon  *geojson.Feature
}

func StandardDeviationalEllipse(points *geojson.FeatureCollection, options ...int) (*StandardDeviationalEllipseResult, error) {
	steps := 64
	if len(options) > 0 && options[0] > 0 {
		steps = options[0]
	}

	if points == nil || len(points.Features) < 3 {
		return nil, fmt.Errorf("at least 3 points required")
	}

	var coords []geojson.Position
	meta.CoordEach(points, func(c geojson.Position, _ int) error {
		coords = append(coords, c)
		return nil
	})

	if len(coords) < 3 {
		return nil, fmt.Errorf("no valid coordinates")
	}

	n := float64(len(coords))

	var cx, cy float64
	for _, c := range coords {
		cx += c[0]
		cy += c[1]
	}
	cx /= n
	cy /= n

	var sumXX, sumYY, sumXY float64
	for _, c := range coords {
		dx := c[0] - cx
		dy := c[1] - cy
		sumXX += dx * dx
		sumYY += dy * dy
		sumXY += dx * dy
	}

	theta := 0.5 * math.Atan2(2*sumXY, sumXX-sumYY)

	sinT := math.Sin(theta)
	cosT := math.Cos(theta)

	var sumU, sumV float64
	for _, c := range coords {
		dx := c[0] - cx
		dy := c[1] - cy
		u := dx*cosT + dy*sinT
		v := -dx*sinT + dy*cosT
		sumU += u * u
		sumV += v * v
	}

	sigmaX := math.Sqrt(sumU / n)
	sigmaY := math.Sqrt(sumV / n)

	var ring []geojson.Position
	for i := 0; i <= steps; i++ {
		angle := float64(i) * 2 * math.Pi / float64(steps)
		ex := sigmaX * 2 * math.Cos(angle)
		ey := sigmaY * 2 * math.Sin(angle)
		px := cx + ex*cosT - ey*sinT
		py := cy + ex*sinT + ey*cosT
		ring = append(ring, geojson.Position{px, py})
	}

	poly := geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{ring}),
		nil,
	)

	return &StandardDeviationalEllipseResult{
		Center:   geojson.Position{cx, cy},
		XAxis:    sigmaX * 2,
		YAxis:    sigmaY * 2,
		Rotation: theta * 180 / math.Pi,
		Polygon:  poly,
	}, nil
}
