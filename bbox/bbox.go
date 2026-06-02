package bbox

import (
	"fmt"
	"math"
	"github.com/ibinh/turf-go/geojson"
)

func BBox(obj any) ([]float64, error) {
	coords, err := getAllCoords(obj)
	if err != nil {
		return nil, err
	}
	if len(coords) == 0 {
		return nil, fmt.Errorf("bbox: no coordinates found")
	}

	minLng, minLat := coords[0][0], coords[0][1]
	maxLng, maxLat := minLng, minLat

	for _, c := range coords {
		if c[0] < minLng {
			minLng = c[0]
		}
		if c[0] > maxLng {
			maxLng = c[0]
		}
		if c[1] < minLat {
			minLat = c[1]
		}
		if c[1] > maxLat {
			maxLat = c[1]
		}
	}

	return []float64{minLng, minLat, maxLng, maxLat}, nil
}

func BBoxPolygon(bbox []float64) (*geojson.Feature, error) {
	if len(bbox) < 4 {
		return nil, fmt.Errorf("bbox-polygon: bbox must have 4 elements")
	}

	minLng, minLat := bbox[0], bbox[1]
	maxLng, maxLat := bbox[2], bbox[3]

	ring := []geojson.Position{
		{minLng, minLat},
		{maxLng, minLat},
		{maxLng, maxLat},
		{minLng, maxLat},
		{minLng, minLat},
	}

	return geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{ring}),
		nil,
	), nil
}

func Envelope(obj any) (*geojson.Feature, error) {
	b, err := BBox(obj)
	if err != nil {
		return nil, err
	}
	return BBoxPolygon(b)
}

func Square(bbox []float64) ([]float64, error) {
	if len(bbox) < 4 {
		return nil, fmt.Errorf("square: bbox must have 4 elements")
	}

	width := bbox[2] - bbox[0]
	height := bbox[3] - bbox[1]

	if width >= height {
		cy := (bbox[1] + bbox[3]) / 2
		halfWidth := width / 2
		return []float64{
			bbox[0],
			cy - halfWidth,
			bbox[2],
			cy + halfWidth,
		}, nil
	}

	cx := (bbox[0] + bbox[2]) / 2
	halfHeight := height / 2
	return []float64{
		cx - halfHeight,
		bbox[1],
		cx + halfHeight,
		bbox[3],
	}, nil
}

func BBoxClip(obj any, bbox []float64) (*geojson.Feature, error) {
	bboxPoly, err := BBoxPolygon(bbox)
	if err != nil {
		return nil, err
	}
	result := geojson.NewFeature(bboxPoly.Geometry, nil)
	return result, nil
}

func getAllCoords(obj any) ([]geojson.Position, error) {
	return geojson.CoordAll(obj)
}

var earthRadius = 6371008.0

func Distance(fromLng, fromLat, toLng, toLat float64) float64 {
	lat1 := fromLat * math.Pi / 180
	lat2 := toLat * math.Pi / 180
	dlat := lat2 - lat1
	dlon := (toLng - fromLng) * math.Pi / 180

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadius * c
}
