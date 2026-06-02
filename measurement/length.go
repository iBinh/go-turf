package measurement

import (
	"fmt"
	"github.com/ibinh/turf-go/geojson"
)

func Length(geom any, units ...Unit) (float64, error) {
	coords, err := geojson.GetCoords(geom)
	if err != nil {
		return 0, fmt.Errorf("length: %w", err)
	}

	unit := UnitKilometers
	if len(units) > 0 {
		unit = units[0]
	}

	var total float64
	switch c := coords.(type) {
	case []geojson.Position:
		for i := 0; i < len(c)-1; i++ {
			total += HaversineDistance(c[i][0], c[i][1], c[i+1][0], c[i+1][1])
		}
	case geojson.Position:
		return 0, nil
	default:
		return 0, fmt.Errorf("length: unsupported coordinate type %T", coords)
	}

	return metersFromMeters(total, unit), nil
}

func Midpoint(from, to any) (*geojson.Feature, error) {
	fromCoord, err := geojson.GetCoord(from)
	if err != nil {
		return nil, fmt.Errorf("from: %w", err)
	}
	toCoord, err := geojson.GetCoord(to)
	if err != nil {
		return nil, fmt.Errorf("to: %w", err)
	}

	lat1 := degToRad(fromCoord[1])
	lon1 := degToRad(fromCoord[0])
	lat2 := degToRad(toCoord[1])
	lon2 := degToRad(toCoord[0])

	dlon := lon2 - lon1

	x := mathCos(lat2) * mathCos(dlon)
	y := mathCos(lat2) * mathSin(dlon)

	lat3 := mathAtan2(mathSin(lat1)+mathSin(lat2),
		mathSqrt((mathCos(lat1)+x)*(mathCos(lat1)+x)+y*y))
	lon3 := lon1 + mathAtan2(y, mathCos(lat1)+x)

	return geojson.NewFeature(
		geojson.NewPoint(geojson.Position{radToDeg(lon3), radToDeg(lat3)}),
		nil,
	), nil
}
