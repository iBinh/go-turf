package measurement

import (
	"math"
	"testing"

	"github.com/ibinh/turf-go/geojson"
)

func TestHaversineDistance(t *testing.T) {
	d := HaversineDistance(0, 0, 0, 0)
	if d != 0 {
		t.Errorf("expected 0 at same point, got %f", d)
	}

	from := geojson.NewFeature(geojson.NewPoint(geojson.Position{-74.006, 40.7128}), nil)
	to := geojson.NewFeature(geojson.NewPoint(geojson.Position{-73.935242, 40.73061}), nil)

	dist, err := Distance(from, to, UnitMeters)
	if err != nil {
		t.Fatal(err)
	}
	if dist < 5000 || dist > 10000 {
		t.Errorf("NYC distance seems off: %f meters", dist)
	}
}

func TestDistanceAcrossAtlantic(t *testing.T) {
	from := geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 51.5}), nil)
	to := geojson.NewFeature(geojson.NewPoint(geojson.Position{-74.006, 40.7128}), nil)

	dist, err := Distance(from, to, UnitKilometers)
	if err != nil {
		t.Fatal(err)
	}
	if dist < 5000 || dist > 6000 {
		t.Errorf("London-NYC distance seems off: %f km", dist)
	}
}

func TestBearing(t *testing.T) {
	from := geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil)
	to := geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 10}), nil)

	bearing, err := Bearing(from, to)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(bearing-0) > 1 && math.Abs(bearing-360) > 1 {
		t.Errorf("expected bearing ~0 N, got %f", bearing)
	}

	to2 := geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 0}), nil)
	bearing2, err := Bearing(from, to2)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(bearing2-90) > 2 {
		t.Errorf("expected bearing ~90 E, got %f", bearing2)
	}
}

func TestDestination(t *testing.T) {
	origin := geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil)
	dest, err := Destination(origin, 111.195, 0, UnitKilometers)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(dest)
	if math.Abs(coord[1]-1) > 0.5 {
		t.Errorf("expected lat ~1, got %f", coord[1])
	}

	dest2, err := Destination(origin, 500, 90, UnitKilometers)
	if err != nil {
		t.Fatal(err)
	}
	coord2, _ := geojson.GetCoord(dest2)
	if coord2[0] <= 0 {
		t.Errorf("expected positive longitude for eastward destination, got %f", coord2[0])
	}
}

func TestLength(t *testing.T) {
	line := geojson.NewFeature(
		geojson.NewLineString([]geojson.Position{{0, 0}, {1, 0}, {2, 0}}),
		nil,
	)
	length, err := Length(line, UnitDegrees)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(length-2) > 0.01 {
		t.Errorf("expected ~2 degree length, got %f", length)
	}

	length2, err := Length(line, UnitMeters)
	if err != nil {
		t.Fatal(err)
	}
	if length2 <= 0 {
		t.Errorf("expected positive length in meters, got %f", length2)
	}
}

func TestArea(t *testing.T) {
	ring := []geojson.Position{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}}
	poly := geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{ring}),
		nil,
	)
	area, err := Area(poly)
	if err != nil {
		t.Fatal(err)
	}
	if area <= 0 {
		t.Errorf("expected positive area, got %f", area)
	}
}

func TestAreaMultiPolygon(t *testing.T) {
	poly1 := [][]geojson.Position{{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}}}
	poly2 := [][]geojson.Position{{{2, 2}, {2, 3}, {3, 3}, {3, 2}, {2, 2}}}
	mp := geojson.NewFeature(
		geojson.NewMultiPolygon([][][]geojson.Position{poly1, poly2}),
		nil,
	)
	area, err := Area(mp)
	if err != nil {
		t.Fatal(err)
	}
	if area <= 0 {
		t.Errorf("expected positive area for multipolygon, got %f", area)
	}
}

func TestMidpoint(t *testing.T) {
	from := geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil)
	to := geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 10}), nil)

	mid, err := Midpoint(from, to)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(mid)
	if math.Abs(coord[1]-5) > 0.01 {
		t.Errorf("expected lat ~5, got %f", coord[1])
	}
}

func TestAlong(t *testing.T) {
	line := geojson.NewFeature(
		geojson.NewLineString([]geojson.Position{{0, 0}, {0, 10}}),
		nil,
	)
	pt, err := Along(line, 5, UnitDegrees)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(pt)
	if math.Abs(coord[1]-5) > 0.5 {
		t.Errorf("expected lat ~5, got %f", coord[1])
	}
}

func TestGreatCircle(t *testing.T) {
	from := geojson.NewFeature(geojson.NewPoint(geojson.Position{-74.006, 40.7128}), nil)
	to := geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 51.5}), nil)

	gc, err := GreatCircle(from, to)
	if err != nil {
		t.Fatal(err)
	}
	if gc == nil {
		t.Fatal("expected great circle feature")
	}
}

func TestRhumbDistance(t *testing.T) {
	from := geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil)
	to := geojson.NewFeature(geojson.NewPoint(geojson.Position{1, 0}), nil)

	dist, err := RhumbDistance(from, to, UnitKilometers)
	if err != nil {
		t.Fatal(err)
	}
	if dist <= 0 {
		t.Errorf("expected positive rhumb distance, got %f", dist)
	}
}

func TestPointToLineDistance(t *testing.T) {
	pt := geojson.NewFeature(geojson.NewPoint(geojson.Position{0.5, 0.5}), nil)
	line := geojson.NewFeature(
		geojson.NewLineString([]geojson.Position{{0, 0}, {1, 1}}),
		nil,
	)
	d, err := PointToLineDistance(pt, line, UnitDegrees)
	if err != nil {
		t.Fatal(err)
	}
	if d > 0.1 {
		t.Errorf("expected small distance, got %f", d)
	}
}

func TestNearestPointOnLine(t *testing.T) {
	pt := geojson.NewFeature(geojson.NewPoint(geojson.Position{0.5, 0.6}), nil)
	line := geojson.NewFeature(
		geojson.NewLineString([]geojson.Position{{0, 0}, {1, 1}}),
		nil,
	)
	nearest, err := NearestPointOnLine(line, pt)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(nearest)
	if coord[0] < 0.4 || coord[0] > 0.6 {
		t.Errorf("expected nearest ~0.55, got %v", coord)
	}
}

func TestConvertLength(t *testing.T) {
	result := convertLength(1, UnitKilometers, UnitMeters)
	if math.Abs(result-1000) > 0.01 {
		t.Errorf("expected 1000m, got %f", result)
	}

	result = convertLength(1, UnitMiles, UnitKilometers)
	if result < 1.6 || result > 1.7 {
		t.Errorf("1 mile ~1.609km, got %f", result)
	}
}

func TestLengthToDegrees(t *testing.T) {
	d := lengthToDegrees(111.195, UnitKilometers)
	if math.Abs(d-1) > 0.01 {
		t.Errorf("111.195km ~1 degree, got %f", d)
	}
}
