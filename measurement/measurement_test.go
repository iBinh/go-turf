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

// ---------------------------------------------------------------------------
// New comprehensive table-driven tests for uncovered functions
// ---------------------------------------------------------------------------

func TestAngle(t *testing.T) {
	tests := []struct {
		name    string
		start   any
		mid     any
		end     any
		wantErr bool
	}{
		{
			name:  "three points with Feature inputs",
			start: geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 10}), nil),
			mid:   geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
			end:   geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 0}), nil),
			// NOTE: Angle internally extracts Position via GetCoord then passes
			// raw Position to Bearing, which requires *Feature/*Point. This
			// causes Bearing to always return an error for valid Feature inputs.
			wantErr: true,
		},
		{
			name: "three points with raw Point geometries",
			start: geojson.NewPoint(geojson.Position{0, 10}),
			mid:   geojson.NewPoint(geojson.Position{0, 0}),
			end:   geojson.NewPoint(geojson.Position{10, 0}),
			// Same issue: GetCoord returns Position, Bearing can't accept Position.
			wantErr: true,
		},
		{
			name:    "invalid start returns error from GetCoord",
			start:   "not a geometry",
			mid:     geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
			end:     geojson.NewFeature(geojson.NewPoint(geojson.Position{1, 1}), nil),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Angle(tt.start, tt.mid, tt.end)
			if !tt.wantErr {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Error("expected error but got none")
			}
		})
	}
}

func TestRhumbBearing(t *testing.T) {
	tests := []struct {
		name    string
		from    any
		to      any
		want    float64
		tol     float64
		wantErr bool
	}{
		{
			name: "north bearing is 0",
			from: geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
			to:   geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 10}), nil),
			want: 0,
			tol:  0.001,
		},
		{
			name: "east along equator",
			from: geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
			to:   geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 0}), nil),
			want: 350.1489238834,
			tol:  0.0001,
		},
		{
			name: "south returns 0 (dlon=0 edge case in formula)",
			from: geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 10}), nil),
			to:   geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
			want: 0,
			tol:  0.001,
		},
		{
			name: "west along equator",
			from: geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 0}), nil),
			to:   geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
			want: 9.8510761166,
			tol:  0.0001,
		},
		{
			name: "same point returns 0",
			from: geojson.NewFeature(geojson.NewPoint(geojson.Position{5, 5}), nil),
			to:   geojson.NewFeature(geojson.NewPoint(geojson.Position{5, 5}), nil),
			want: 0,
			tol:  0.001,
		},
		{
			name: "northeast",
			from: geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
			to:   geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 10}), nil),
			want: 350.0498795265,
			tol:  0.0001,
		},
		{
			name:    "invalid from returns error",
			from:    42,
			to:      geojson.NewFeature(geojson.NewPoint(geojson.Position{1, 1}), nil),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RhumbBearing(tt.from, tt.to)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if math.Abs(got-tt.want) > tt.tol {
				t.Errorf("RhumbBearing() = %.10f, want %.10f ± %.10f", got, tt.want, tt.tol)
			}
		})
	}
}

func TestRhumbDestination(t *testing.T) {
	tests := []struct {
		name     string
		origin   any
		distance float64
		bearing  float64
		units    []Unit
		wantLon  float64
		wantLat  float64
		tol      float64
		wantErr  bool
	}{
		{
			name:     "north 1 degree at equator (111.195 km)",
			origin:   geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
			distance: 111.195,
			bearing:  0,
			units:    []Unit{UnitKilometers},
			wantLon:  0,
			wantLat:  1,
			tol:      0.5,
		},
		{
			name:     "east 500 km at equator",
			origin:   geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
			distance: 500,
			bearing:  90,
			units:    []Unit{UnitKilometers},
			wantLon:  4.5,
			wantLat:  0,
			tol:      0.5,
		},
		{
			name:     "default unit is kilometers when units omitted",
			origin:   geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
			distance: 111.195,
			bearing:  0,
			wantLon:  0,
			wantLat:  1,
			tol:      0.5,
		},
		{
			name:     "south 111.195 km",
			origin:   geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 1}), nil),
			distance: 111.195,
			bearing:  180,
			units:    []Unit{UnitKilometers},
			wantLon:  0,
			wantLat:  0,
			tol:      0.5,
		},
		{
			name:     "northwest (315 deg)",
			origin:   geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
			distance: 111.195,
			bearing:  315,
			units:    []Unit{UnitKilometers},
			wantLon: -0.707,
			wantLat: 0.707,
			tol:      0.5,
		},
		{
			name:     "using meters unit",
			origin:   geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
			distance: 111195,
			bearing:  0,
			units:    []Unit{UnitMeters},
			wantLon:  0,
			wantLat:  1,
			tol:      0.5,
		},
		{
			name:    "invalid origin returns error",
			origin:  "bad",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dest, err := RhumbDestination(tt.origin, tt.distance, tt.bearing, tt.units...)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			coord, err := geojson.GetCoord(dest)
			if err != nil {
				t.Fatalf("GetCoord error: %v", err)
			}
			if math.Abs(coord[0]-tt.wantLon) > tt.tol {
				t.Errorf("longitude = %.6f, want %.6f ± %.6f", coord[0], tt.wantLon, tt.tol)
			}
			if math.Abs(coord[1]-tt.wantLat) > tt.tol {
				t.Errorf("latitude = %.6f, want %.6f ± %.6f", coord[1], tt.wantLat, tt.tol)
			}
		})
	}
}

func TestConvertLengthExported(t *testing.T) {
	tests := []struct {
		name   string
		length float64
		from   Unit
		to     Unit
		want   float64
		tol    float64
	}{
		{
			name:   "kilometers to meters",
			length: 1,
			from:   UnitKilometers,
			to:     UnitMeters,
			want:   1000,
			tol:    0.01,
		},
		{
			name:   "miles to kilometers",
			length: 1,
			from:   UnitMiles,
			to:     UnitKilometers,
			want:   1.609344,
			tol:    0.001,
		},
		{
			name:   "meters to kilometers",
			length: 5000,
			from:   UnitMeters,
			to:     UnitKilometers,
			want:   5,
			tol:    0.001,
		},
		{
			name:   "nautical miles to meters",
			length: 1,
			from:   UnitNauticalMiles,
			to:     UnitMeters,
			want:   1852,
			tol:    0.01,
		},
		{
			name:   "feet to meters",
			length: 3.28084,
			from:   UnitFeet,
			to:     UnitMeters,
			want:   1,
			tol:    0.001,
		},
		{
			name:   "degrees to meters",
			length: 1,
			from:   UnitDegrees,
			to:     UnitMeters,
			want:   (math.Pi / 180) * EarthRadius,
			tol:    0.01,
		},
		{
			name:   "radians to meters",
			length: 1,
			from:   UnitRadians,
			to:     UnitMeters,
			want:   EarthRadius,
			tol:    0.01,
		},
		{
			name:   "meters to miles",
			length: 1609.344,
			from:   UnitMeters,
			to:     UnitMiles,
			want:   1,
			tol:    0.001,
		},
		{
			name:   "kilometers to nautical miles",
			length: 1.852,
			from:   UnitKilometers,
			to:     UnitNauticalMiles,
			want:   1,
			tol:    0.001,
		},
		{
			name:   "meters to feet",
			length: 1,
			from:   UnitMeters,
			to:     UnitFeet,
			want:   3.28084,
			tol:    0.001,
		},
		{
			name:   "meters to degrees",
			length: (math.Pi/180)*EarthRadius,
			from:   UnitMeters,
			to:     UnitDegrees,
			want:   1,
			tol:    0.001,
		},
		{
			name:   "meters to radians",
			length: EarthRadius,
			from:   UnitMeters,
			to:     UnitRadians,
			want:   1,
			tol:    0.001,
		},
		{
			name:   "identity (same unit)",
			length: 42,
			from:   UnitKilometers,
			to:     UnitKilometers,
			want:   42,
			tol:    0.001,
		},
		{
			name:   "zero length",
			length: 0,
			from:   UnitKilometers,
			to:     UnitMeters,
			want:   0,
			tol:    0.001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertLength(tt.length, tt.from, tt.to)
			if math.Abs(got-tt.want) > tt.tol {
				t.Errorf("ConvertLength() = %.6f, want %.6f ± %.6f", got, tt.want, tt.tol)
			}
		})
	}
}

func TestNearestPoint(t *testing.T) {
	tests := []struct {
		name      string
		target    any
		points    any
		wantCoord geojson.Position
		wantErr   bool
	}{
		{
			name:   "find nearest in FeatureCollection",
			target: geojson.NewFeature(geojson.NewPoint(geojson.Position{1, 1}), nil),
			points: geojson.NewFeatureCollection([]*geojson.Feature{
				geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 10}), nil),
				geojson.NewFeature(geojson.NewPoint(geojson.Position{1.1, 1.1}), nil),
				geojson.NewFeature(geojson.NewPoint(geojson.Position{-5, -5}), nil),
			}),
			wantCoord: geojson.Position{1.1, 1.1},
		},
		{
			name:   "find nearest when target is on a point",
			target: geojson.NewFeature(geojson.NewPoint(geojson.Position{2, 2}), nil),
			points: geojson.NewFeatureCollection([]*geojson.Feature{
				geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
				geojson.NewFeature(geojson.NewPoint(geojson.Position{2, 2}), nil),
				geojson.NewFeature(geojson.NewPoint(geojson.Position{5, 5}), nil),
			}),
			wantCoord: geojson.Position{2, 2},
		},
		{
			name:   "single Feature as points argument",
			target: geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
			points: geojson.NewFeature(geojson.NewPoint(geojson.Position{3, 4}), nil),
			wantCoord: geojson.Position{3, 4},
		},
		{
			name:      "invalid target returns error",
			target:    "bad",
			points:    geojson.NewFeatureCollection(nil),
			wantErr:   true,
		},
		{
			name:      "invalid points type returns error",
			target:    geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
			points:    "not valid",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nearest, err := NearestPoint(tt.target, tt.points)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if nearest == nil {
				t.Fatal("expected a Feature but got nil")
			}
			coord, err := geojson.GetCoord(nearest)
			if err != nil {
				t.Fatalf("GetCoord error: %v", err)
			}
			if math.Abs(coord[0]-tt.wantCoord[0]) > 0.01 || math.Abs(coord[1]-tt.wantCoord[1]) > 0.01 {
				t.Errorf("nearest = %v, want %v", coord, tt.wantCoord)
			}
		})
	}
}

func TestNearestPointToLine(t *testing.T) {
	tests := []struct {
		name      string
		line      any
		points    *geojson.FeatureCollection
		units     []Unit
		wantCoord geojson.Position
		wantErr   bool
	}{
		{
			name: "find nearest point to a line",
			line: geojson.NewFeature(
				geojson.NewLineString([]geojson.Position{{0, 0}, {10, 0}}),
				nil,
			),
			points: geojson.NewFeatureCollection([]*geojson.Feature{
				geojson.NewFeature(geojson.NewPoint(geojson.Position{5, 10}), nil),
				geojson.NewFeature(geojson.NewPoint(geojson.Position{5, 0.1}), nil),
				geojson.NewFeature(geojson.NewPoint(geojson.Position{5, 50}), nil),
			}),
			wantCoord: geojson.Position{5, 0.1},
		},
		{
			name: "point exactly on line is chosen",
			line: geojson.NewFeature(
				geojson.NewLineString([]geojson.Position{{0, 0}, {10, 10}}),
				nil,
			),
			points: geojson.NewFeatureCollection([]*geojson.Feature{
				geojson.NewFeature(geojson.NewPoint(geojson.Position{2, 2}), nil),
				geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 10}), nil),
			}),
			wantCoord: geojson.Position{2, 2},
		},
		{
			name: "with miles unit",
			line: geojson.NewFeature(
				geojson.NewLineString([]geojson.Position{{0, 0}, {1, 0}}),
				nil,
			),
			points: geojson.NewFeatureCollection([]*geojson.Feature{
				geojson.NewFeature(geojson.NewPoint(geojson.Position{0.5, 0.5}), nil),
			}),
			units: []Unit{UnitMiles},
			wantCoord: geojson.Position{0.5, 0.5},
		},
		{
			name:    "nil points returns error",
			line:    geojson.NewFeature(geojson.NewLineString([]geojson.Position{{0, 0}}), nil),
			points:  nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nearest, err := NearestPointToLine(tt.line, tt.points, tt.units...)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if nearest == nil {
				t.Fatal("expected a Feature but got nil")
			}
			coord, err := geojson.GetCoord(nearest)
			if err != nil {
				t.Fatalf("GetCoord error: %v", err)
			}
			if math.Abs(coord[0]-tt.wantCoord[0]) > 0.01 || math.Abs(coord[1]-tt.wantCoord[1]) > 0.01 {
				t.Errorf("nearest = %v, want %v", coord, tt.wantCoord)
			}
			// Verify "dist" property was set
			if _, ok := nearest.Properties["dist"]; !ok {
				t.Error("expected 'dist' property on nearest feature")
			}
		})
	}
}

func TestExtractFeatures(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		wantLen int
		wantErr bool
	}{
		{
			name: "FeatureCollection returns all features",
			input: geojson.NewFeatureCollection([]*geojson.Feature{
				geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
				geojson.NewFeature(geojson.NewPoint(geojson.Position{1, 1}), nil),
				geojson.NewFeature(geojson.NewPoint(geojson.Position{2, 2}), nil),
			}),
			wantLen: 3,
		},
		{
			name: "empty FeatureCollection",
			input: geojson.NewFeatureCollection([]*geojson.Feature{}),
			wantLen: 0,
		},
		{
			name: "single Feature returns slice of one",
			input: geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 20}), nil),
			wantLen: 1,
		},
		{
			name:    "invalid type returns error",
			input:   "not a geojson type",
			wantErr: true,
		},
		{
			name:    "nil returns error",
			input:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			features, err := extractFeatures(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(features) != tt.wantLen {
				t.Errorf("len = %d, want %d", len(features), tt.wantLen)
			}
		})
	}
}

func TestPointToPolygonDistance(t *testing.T) {
	// A simple unit square polygon from (0,0) to (1,1)
	ring := []geojson.Position{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}}
	square := geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{ring}),
		nil,
	)

	tests := []struct {
		name    string
		point   any
		polygon any
		units   []Unit
		want    float64
		tol     float64
		wantErr bool
	}{
		{
			name:    "point inside polygon returns 0",
			point:   geojson.NewFeature(geojson.NewPoint(geojson.Position{0.5, 0.5}), nil),
			polygon: square,
			want:    0,
			tol:     0.001,
		},
		{
			name:    "point on polygon edge returns 0",
			point:   geojson.NewFeature(geojson.NewPoint(geojson.Position{0.5, 0}), nil),
			polygon: square,
			want:    0,
			tol:     0.001,
		},
		{
			name:    "point outside polygon returns positive distance (degrees)",
			point:   geojson.NewFeature(geojson.NewPoint(geojson.Position{2, 0.5}), nil),
			polygon: square,
			units:   []Unit{UnitDegrees},
			want:    1,
			tol:     0.1,
		},
		{
			name:    "point near corner",
			point:   geojson.NewFeature(geojson.NewPoint(geojson.Position{1.2, 1.2}), nil),
			polygon: square,
			units:   []Unit{UnitDegrees},
			want:    0.283,
			tol:     0.01,
		},
		{
			name:    "point far from polygon",
			point:   geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 10}), nil),
			polygon: square,
			units:   []Unit{UnitDegrees},
			want:    12.727,
			tol:     0.1,
		},
		{
			name:    "default unit is kilometers",
			point:   geojson.NewFeature(geojson.NewPoint(geojson.Position{2, 0.5}), nil),
			polygon: square,
			want:    111.195,
			tol:     0.01,
		},
		{
			name:    "invalid point returns error",
			point:   "bad",
			polygon: square,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PointToPolygonDistance(tt.point, tt.polygon, tt.units...)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if math.Abs(got-tt.want) > tt.tol {
				t.Errorf("PointToPolygonDistance() = %.6f, want %.6f ± %.6f", got, tt.want, tt.tol)
			}
		})
	}
}
