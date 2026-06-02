package directionalmean

import (
	"fmt"
	"math"

	"github.com/ibinh/turf-go/geojson"
	"github.com/ibinh/turf-go/meta"
)

func DirectionalMean(geom any) (float64, error) {
	var bearings []float64

	err := meta.GeomEach(geom, func(g geojson.Geometry, _ int) error {
		switch geo := g.(type) {
		case *geojson.LineString:
			if len(geo.Coordinates) >= 2 {
				dx := geo.Coordinates[len(geo.Coordinates)-1][0] - geo.Coordinates[0][0]
				dy := geo.Coordinates[len(geo.Coordinates)-1][1] - geo.Coordinates[0][1]
				b := math.Atan2(dx, dy) * 180 / math.Pi
				bearings = append(bearings, b)
			}
		case *geojson.MultiLineString:
			for _, line := range geo.Coordinates {
				if len(line) >= 2 {
					dx := line[len(line)-1][0] - line[0][0]
					dy := line[len(line)-1][1] - line[0][1]
					b := math.Atan2(dx, dy) * 180 / math.Pi
					bearings = append(bearings, b)
				}
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	if len(bearings) == 0 {
		return 0, fmt.Errorf("no line features found")
	}

	var sinSum, cosSum float64
	for _, b := range bearings {
		r := b * math.Pi / 180
		sinSum += math.Sin(r)
		cosSum += math.Cos(r)
	}

	mean := math.Atan2(sinSum/float64(len(bearings)), cosSum/float64(len(bearings))) * 180 / math.Pi
	if mean < 0 {
		mean += 360
	}
	return mean, nil
}
