package measurement

import (
	"fmt"
	"math"
	"github.com/ibinh/turf-go/geojson"
)

func Destination(origin any, distance float64, bearing float64, units ...Unit) (*geojson.Feature, error) {
	coord, err := geojson.GetCoord(origin)
	if err != nil {
		return nil, fmt.Errorf("origin: %w", err)
	}

	unit := UnitKilometers
	if len(units) > 0 {
		unit = units[0]
	}

	d := lengthToDegrees(distance, unit)
	rad := degToRad(d)
	bearingRad := degToRad(bearing)

	lat1 := degToRad(coord[1])
	lon1 := degToRad(coord[0])

	lat2 := math.Asin(math.Sin(lat1)*math.Cos(rad) + math.Cos(lat1)*math.Sin(rad)*math.Cos(bearingRad))
	lon2 := lon1 + math.Atan2(math.Sin(bearingRad)*math.Sin(rad)*math.Cos(lat1), math.Cos(rad)-math.Sin(lat1)*math.Sin(lat2))

	lon2 = math.Mod(lon2+math.Pi, 2*math.Pi) - math.Pi

	return geojson.NewFeature(
		geojson.NewPoint(geojson.Position{radToDeg(lon2), radToDeg(lat2)}),
		nil,
	), nil
}

func RhumbDestination(origin any, distance float64, bearing float64, units ...Unit) (*geojson.Feature, error) {
	coord, err := geojson.GetCoord(origin)
	if err != nil {
		return nil, fmt.Errorf("origin: %w", err)
	}

	unit := UnitKilometers
	if len(units) > 0 {
		unit = units[0]
	}

	d := lengthToDegrees(distance, unit)
	delta := degToRad(d)
	bearingRad := degToRad(bearing)

	lat1 := degToRad(coord[1])
	lon1 := degToRad(coord[0])

	lat2 := lat1 + delta*math.Cos(bearingRad)
	dPhi := math.Log(math.Tan(lat2/2+math.Pi/4) / math.Tan(lat1/2+math.Pi/4))
	q := (lat2 - lat1) / dPhi
	if math.IsInf(q, 0) || math.IsNaN(q) {
		q = math.Cos(lat1)
	}
	dlon := delta * math.Sin(bearingRad) / q
	lon2 := lon1 + dlon

	lat2 = math.Mod(lat2+math.Pi, 2*math.Pi) - math.Pi
	lat2 = math.Max(-math.Pi/2, math.Min(math.Pi/2, lat2))
	lon2 = math.Mod(lon2+math.Pi, 2*math.Pi) - math.Pi

	return geojson.NewFeature(
		geojson.NewPoint(geojson.Position{radToDeg(lon2), radToDeg(lat2)}),
		nil,
	), nil
}
