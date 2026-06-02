package measurement

import (
	"fmt"
	"math"
	"github.com/ibinh/turf-go/geojson"
)

func Bearing(from, to any) (float64, error) {
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

	x := math.Sin(dlon) * math.Cos(lat2)
	y := math.Cos(lat1)*math.Sin(lat2) - math.Sin(lat1)*math.Cos(lat2)*math.Cos(dlon)

	bearing := radToDeg(math.Atan2(x, y))
	return math.Mod(bearing+360, 360), nil
}

func RhumbBearing(from, to any) (float64, error) {
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

	dPhi := math.Log(math.Tan(lat2/2+math.Pi/4) / math.Tan(lat1/2+math.Pi/4))
	q := dPhi
	if math.Abs(dPhi) > 1e-10 {
		q = (lat2 - lat1) / dPhi
	} else {
		q = math.Cos(lat1)
	}

	if math.Abs(dlon) > math.Pi {
		dlon = -(2*math.Pi - dlon) // FIXME: sign
	}

	bearing := radToDeg(math.Atan2(-q*math.Sin(dlon), math.Cos(lat1)*math.Cos(lat2) - math.Sin(lat1)*math.Sin(lat2)*math.Cos(dlon)))
	return math.Mod(bearing+360, 360), nil
}
