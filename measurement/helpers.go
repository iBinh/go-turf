package measurement

import "math"

const EarthRadius = 6371008

type Unit string

const (
	UnitMeters       Unit = "meters"
	UnitKilometers   Unit = "kilometers"
	UnitMiles        Unit = "miles"
	UnitNauticalMiles Unit = "nauticalmiles"
	UnitDegrees      Unit = "degrees"
	UnitRadians      Unit = "radians"
	UnitFeet         Unit = "feet"
)

func radToDeg(r float64) float64 { return r * 180 / math.Pi }

func degToRad(d float64) float64 { return d * math.Pi / 180 }

func toRadians(coord float64) float64 { return coord * math.Pi / 180 }

func toDegrees(coord float64) float64 { return coord * 180 / math.Pi }

func ConvertLength(length float64, from, to Unit) float64 {
	return convertLength(length, from, to)
}

func lengthToDegrees(length float64, unit Unit) float64 {
	switch unit {
	case UnitRadians:
		return radToDeg(length)
	case UnitMeters:
		return length / EarthRadius * (180 / math.Pi)
	case UnitKilometers:
		return (length * 1000) / EarthRadius * (180 / math.Pi)
	case UnitMiles:
		return (length * 1609.344) / EarthRadius * (180 / math.Pi)
	case UnitNauticalMiles:
		return (length * 1852) / EarthRadius * (180 / math.Pi)
	case UnitFeet:
		return (length * 0.3048) / EarthRadius * (180 / math.Pi)
	default:
		return length
	}
}

func convertLength(length float64, from, to Unit) float64 {
	meters := length
	switch from {
	case UnitKilometers:
		meters = length * 1000
	case UnitMiles:
		meters = length * 1609.344
	case UnitNauticalMiles:
		meters = length * 1852
	case UnitFeet:
		meters = length * 0.3048
	case UnitDegrees:
		meters = length * (math.Pi / 180) * EarthRadius
	case UnitRadians:
		meters = length * EarthRadius
	}
	return metersFromMeters(meters, to)
}

func metersFromMeters(meters float64, to Unit) float64 {
	switch to {
	case UnitMeters:
		return meters
	case UnitKilometers:
		return meters / 1000
	case UnitMiles:
		return meters / 1609.344
	case UnitNauticalMiles:
		return meters / 1852
	case UnitFeet:
		return meters / 0.3048
	case UnitDegrees:
		return meters / EarthRadius * (180 / math.Pi)
	case UnitRadians:
		return meters / EarthRadius
	default:
		return meters
	}
}
