package bbox

import (
	"testing"
	"github.com/ibinh/turf-go/geojson"
)

func TestBBox(t *testing.T) {
	f := geojson.NewFeature(
		geojson.NewPoint(geojson.Position{1, 2}),
		nil,
	)
	b, err := BBox(f)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 4 || b[0] != 1 || b[1] != 2 || b[2] != 1 || b[3] != 2 {
		t.Errorf("unexpected bbox: %v", b)
	}
}

func TestBBoxMultiPoint(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{1, 2}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{3, 5}), nil),
	})
	b, err := BBox(fc)
	if err != nil {
		t.Fatal(err)
	}
	if b[0] != 1 || b[1] != 2 || b[2] != 3 || b[3] != 5 {
		t.Errorf("unexpected bbox: %v", b)
	}
}

func TestBBoxPolygon(t *testing.T) {
	bbox := []float64{0, 0, 10, 10}
	f, err := BBoxPolygon(bbox)
	if err != nil {
		t.Fatal(err)
	}
	poly, ok := f.Geometry.(*geojson.Polygon)
	if !ok {
		t.Fatalf("expected Polygon, got %T", f.Geometry)
	}
	if len(poly.Coordinates[0]) != 5 {
		t.Errorf("expected 5 vertices, got %d", len(poly.Coordinates[0]))
	}
}

func TestEnvelope(t *testing.T) {
	f := geojson.NewFeature(
		geojson.NewPoint(geojson.Position{5, 5}),
		nil,
	)
	env, err := Envelope(f)
	if err != nil {
		t.Fatal(err)
	}
	if env == nil {
		t.Fatal("expected envelope feature")
	}
}

func TestSquare(t *testing.T) {
	bbox := []float64{0, 0, 2, 1}
	sq, err := Square(bbox)
	if err != nil {
		t.Fatal(err)
	}
	width := sq[2] - sq[0]
	height := sq[3] - sq[1]
	if mathAbs(width-height) > 0.001 {
		t.Errorf("square should have equal width/height, got %f x %f", width, height)
	}
}

func TestSquareTall(t *testing.T) {
	bbox := []float64{0, 0, 1, 3}
	sq, err := Square(bbox)
	if err != nil {
		t.Fatal(err)
	}
	width := sq[2] - sq[0]
	height := sq[3] - sq[1]
	if mathAbs(width-height) > 0.001 {
		t.Errorf("square should have equal width/height, got %f x %f", width, height)
	}
}

func TestBBoxPolygonRoundtrip(t *testing.T) {
	f := geojson.NewFeature(
		geojson.NewPoint(geojson.Position{1, 2}),
		nil,
	)
	b, err := BBox(f)
	if err != nil {
		t.Fatal(err)
	}
	poly, err := BBoxPolygon(b)
	if err != nil {
		t.Fatal(err)
	}
	if poly == nil {
		t.Fatal("expected polygon")
	}
}

func mathAbs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
