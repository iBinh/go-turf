package transform

import (
	"math"
	"testing"

	"github.com/ibinh/turf-go/geojson"
)

func TestTransformRotate(t *testing.T) {
	pt := geojson.NewFeature(geojson.NewPoint(geojson.Position{1, 0}), nil)
	rotated, err := TransformRotate(pt, 90, geojson.Position{0, 0})
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(rotated)
	if math.Abs(coord[0]-0) > 0.001 || math.Abs(coord[1]-1) > 0.001 {
		t.Errorf("expected ~[0,1] after 90 deg, got %v", coord)
	}
}

func TestTransformRotateNoPivot(t *testing.T) {
	pt := geojson.NewFeature(geojson.NewPoint(geojson.Position{1, 0}), nil)
	rotated, err := TransformRotate(pt, 180)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(rotated)
	if math.Abs(coord[0]-1) > 0.001 || math.Abs(coord[1]-0) > 0.001 {
		t.Errorf("point should stay same when rotating around itself, got %v", coord)
	}
}

func TestTransformRotateLine(t *testing.T) {
	line := geojson.NewFeature(
		geojson.NewLineString([]geojson.Position{{0, 0}, {2, 0}}),
		nil,
	)
	rotated, err := TransformRotate(line, 90, geojson.Position{0, 0})
	if err != nil {
		t.Fatal(err)
	}
	coords, _ := geojson.GetCoords(rotated)
	pts := coords.([]geojson.Position)
	if math.Abs(pts[0][0]) > 0.001 || math.Abs(pts[0][1]) > 0.001 {
		t.Errorf("first point should stay at origin, got %v", pts[0])
	}
	if math.Abs(pts[1][0]) > 0.001 || math.Abs(pts[1][1]-2) > 0.001 {
		t.Errorf("expected ~[0,2], got %v", pts[1])
	}
}

func TestTransformScale(t *testing.T) {
	pt := geojson.NewFeature(geojson.NewPoint(geojson.Position{2, 3}), nil)
	scaled, err := TransformScale(pt, 2, geojson.Position{0, 0})
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(scaled)
	if math.Abs(coord[0]-4) > 0.001 || math.Abs(coord[1]-6) > 0.001 {
		t.Errorf("expected [4,6], got %v", coord)
	}
}

func TestTransformScaleXY(t *testing.T) {
	pt := geojson.NewFeature(geojson.NewPoint(geojson.Position{2, 3}), nil)
	scaled, err := TransformScaleXY(pt, 2, 3, geojson.Position{0, 0})
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(scaled)
	if math.Abs(coord[0]-4) > 0.001 || math.Abs(coord[1]-9) > 0.001 {
		t.Errorf("expected [4,9], got %v", coord)
	}
}

func TestTransformTranslate(t *testing.T) {
	pt := geojson.NewFeature(geojson.NewPoint(geojson.Position{1, 2}), nil)
	translated, err := TransformTranslate(pt, 10, 20)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(translated)
	if math.Abs(coord[0]-11) > 0.001 || math.Abs(coord[1]-22) > 0.001 {
		t.Errorf("expected [11,22], got %v", coord)
	}
}

func TestTranslateMeters(t *testing.T) {
	pt := geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil)
	translated, err := TransformTranslate(pt, 111319.9, 0, &TranslateOptions{Units: "meters"})
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(translated)
	if math.Abs(coord[0]-1) > 0.01 {
		t.Errorf("expected ~1 degree east, got %f", coord[0])
	}
}

func TestFlip(t *testing.T) {
	pt := geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 20}), nil)
	flipped, err := Flip(pt)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(flipped)
	if coord[0] != 20 || coord[1] != 10 {
		t.Errorf("expected [20,10], got %v", coord)
	}
}

func TestTruncate(t *testing.T) {
	pt := geojson.NewFeature(geojson.NewPoint(geojson.Position{1.12345, 2.67890}), nil)
	truncated, err := Truncate(pt, 2)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(truncated)
	if math.Abs(coord[0]-1.12) > 0.001 || math.Abs(coord[1]-2.68) > 0.001 {
		t.Errorf("expected [1.12, 2.68], got %v", coord)
	}
}

func TestCleanCoords(t *testing.T) {
	line := geojson.NewFeature(
		geojson.NewLineString([]geojson.Position{{0, 0}, {0, 0}, {1, 1}, {1, 1}, {2, 2}}),
		nil,
	)
	cleaned, err := CleanCoords(line)
	if err != nil {
		t.Fatal(err)
	}
	coords, _ := geojson.GetCoords(cleaned)
	pts := coords.([]geojson.Position)
	if len(pts) != 3 {
		t.Errorf("expected 3 unique points, got %d", len(pts))
	}
}

func TestRewind(t *testing.T) {
	ring := []geojson.Position{{0, 0}, {1, 1}, {2, 0}, {0, 0}}
	poly := geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{ring}),
		nil,
	)
	rewound, err := Rewind(poly)
	if err != nil {
		t.Fatal(err)
	}
	p, err := geojson.GetGeometry(rewound)
	if err != nil {
		t.Fatal(err)
	}
	// After rewind, exterior should be CW
	pg, ok := p.(*geojson.Polygon)
	if !ok {
		t.Fatalf("expected Polygon")
	}
	if !isCWRing(pg.Coordinates[0]) {
		t.Error("exterior should be CW after default rewind")
	}
}

func TestRewindReverse(t *testing.T) {
	ring := []geojson.Position{{0, 0}, {2, 0}, {1, 1}, {0, 0}}
	poly := geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{ring}),
		nil,
	)
	rewound, err := Rewind(poly, true)
	if err != nil {
		t.Fatal(err)
	}
	p, err := geojson.GetGeometry(rewound)
	if err != nil {
		t.Fatal(err)
	}
	pg, ok := p.(*geojson.Polygon)
	if !ok {
		t.Fatalf("expected Polygon")
	}
	if isCWRing(pg.Coordinates[0]) {
		t.Error("exterior should be CCW with reverse=true")
	}
}

func TestToMercator(t *testing.T) {
	pt := geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil)
	merc, err := ToMercator(pt)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(merc)
	if math.Abs(coord[0]) > 0.001 || math.Abs(coord[1]) > 0.001 {
		t.Errorf("expected [0,0] at equator/prime, got %v", coord)
	}
}

func TestToWGS84(t *testing.T) {
	pt := geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil)
	wgs, err := ToWGS84(pt)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(wgs)
	if math.Abs(coord[0]) > 0.001 || math.Abs(coord[1]) > 0.001 {
		t.Errorf("expected [0,0], got %v", coord)
	}
}

func TestMercatorRoundtrip(t *testing.T) {
	original := geojson.NewFeature(geojson.NewPoint(geojson.Position{-73.935242, 40.73061}), nil)
	merc, err := ToMercator(original)
	if err != nil {
		t.Fatal(err)
	}
	wgs, err := ToWGS84(merc)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(wgs)
	if math.Abs(coord[0]+73.935242) > 0.01 || math.Abs(coord[1]-40.73061) > 0.01 {
		t.Errorf("roundtrip failed: expected ~[-73.94, 40.73], got %v", coord)
	}
}

func TestTransformPolygon(t *testing.T) {
	ring := []geojson.Position{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}}
	poly := geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{ring}),
		nil,
	)
	translated, err := TransformTranslate(poly, 5, 5)
	if err != nil {
		t.Fatal(err)
	}
	coords, _ := geojson.GetCoords(translated)
	rings := coords.([][]geojson.Position)
	first := rings[0][0]
	if math.Abs(first[0]-5) > 0.001 || math.Abs(first[1]-5) > 0.001 {
		t.Errorf("first vertex should be [5,5], got %v", first)
	}
}
