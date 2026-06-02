package measurement

import (
	"fmt"
	"math"
	"github.com/ibinh/turf-go/geojson"
)

func HaversineDistance(fromLng, fromLat, toLng, toLat float64) float64 {
	lat1 := degToRad(fromLat)
	lat2 := degToRad(toLat)
	dlat := lat2 - lat1
	dlon := degToRad(toLng - fromLng)

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EarthRadius * c
}

func Distance(from, to any, units ...Unit) (float64, error) {
	fromCoord, err := geojson.GetCoord(from)
	if err != nil {
		return 0, fmt.Errorf("from: %w", err)
	}
	toCoord, err := geojson.GetCoord(to)
	if err != nil {
		return 0, fmt.Errorf("to: %w", err)
	}

	unit := UnitKilometers
	if len(units) > 0 {
		unit = units[0]
	}

	meters := HaversineDistance(fromCoord[0], fromCoord[1], toCoord[0], toCoord[1])
	return metersFromMeters(meters, unit), nil
}

func RhumbDistance(from, to any, units ...Unit) (float64, error) {
	fromCoord, err := geojson.GetCoord(from)
	if err != nil {
		return 0, fmt.Errorf("from: %w", err)
	}
	toCoord, err := geojson.GetCoord(to)
	if err != nil {
		return 0, fmt.Errorf("to: %w", err)
	}

	lat1 := degToRad(fromCoord[1])
	lat2 := degToRad(toCoord[1])
	dlon := degToRad(toCoord[0] - fromCoord[0])

	dLat := lat2 - lat1
	dPhi := math.Log(math.Tan(lat2/2+math.Pi/4) / math.Tan(lat1/2+math.Pi/4))

	q := dLat / dPhi
	if math.IsInf(q, 0) || math.IsNaN(q) {
		q = math.Cos(lat1)
	}

	d := math.Sqrt(dLat*dLat + q*q*dlon*dlon) * EarthRadius
	unit := UnitKilometers
	if len(units) > 0 {
		unit = units[0]
	}
	return metersFromMeters(d, unit), nil
}
