package measurement

import (
	"fmt"
	"math"
	"github.com/ibinh/turf-go/geojson"
)

func Along(line any, distance float64, units ...Unit) (*geojson.Feature, error) {
	coords, err := geojson.GetCoords(line)
	if err != nil {
		return nil, fmt.Errorf("along: %w", err)
	}
	pts, ok := coords.([]geojson.Position)
	if !ok {
		return nil, fmt.Errorf("along: expected LineString coordinates")
	}

	unit := UnitKilometers
	if len(units) > 0 {
		unit = units[0]
	}

	travelled := 0.0
	targetDist := distance
	if unit != UnitMeters {
		targetDist = convertLength(distance, unit, UnitMeters)
	}

	for i := 0; i < len(pts)-1; i++ {
		segDist := HaversineDistance(pts[i][0], pts[i][1], pts[i+1][0], pts[i+1][1])
		if travelled+segDist >= targetDist {
			remain := targetDist - travelled
			fraction := remain / segDist
			lng := pts[i][0] + (pts[i+1][0]-pts[i][0])*fraction
			lat := pts[i][1] + (pts[i+1][1]-pts[i][1])*fraction
			return geojson.NewFeature(
				geojson.NewPoint(geojson.Position{lng, lat}),
				nil,
			), nil
		}
		travelled += segDist
	}

	last := pts[len(pts)-1]
	return geojson.NewFeature(
		geojson.NewPoint(geojson.Position{last[0], last[1]}),
		nil,
	), nil
}

func GreatCircle(from, to any, options ...any) (*geojson.Feature, error) {
	fromCoord, err := geojson.GetCoord(from)
	if err != nil {
		return nil, fmt.Errorf("from: %w", err)
	}
	toCoord, err := geojson.GetCoord(to)
	if err != nil {
		return nil, fmt.Errorf("to: %w", err)
	}

	npoints := 100
	if len(options) > 0 {
		if v, ok := options[0].(int); ok && v > 1 {
			npoints = v
		}
	}

	lon1 := degToRad(fromCoord[0])
	lat1 := degToRad(fromCoord[1])
	lon2 := degToRad(toCoord[0])
	lat2 := degToRad(toCoord[1])

	d := 2 * math.Asin(math.Sqrt(
		math.Pow(math.Sin((lat2-lat1)/2), 2)+
			math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin((lon2-lon1)/2), 2)))

	pts := make([]geojson.Position, npoints)
	for i := 0; i < npoints; i++ {
		f := float64(i) / float64(npoints-1)
		a := math.Sin((1-f)*d) / math.Sin(d)
		b := math.Sin(f*d) / math.Sin(d)
		x := a*math.Cos(lat1)*math.Cos(lon1) + b*math.Cos(lat2)*math.Cos(lon2)
		y := a*math.Cos(lat1)*math.Sin(lon1) + b*math.Cos(lat2)*math.Sin(lon2)
		z := a*math.Sin(lat1) + b*math.Sin(lat2)
		pts[i] = geojson.Position{radToDeg(math.Atan2(y, x)), radToDeg(math.Atan2(z, math.Sqrt(x*x+y*y)))}
	}

	return geojson.NewFeature(
		geojson.NewLineString(pts),
		nil,
	), nil
}
