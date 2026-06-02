package smooth

import (
	"testing"

	"github.com/ibinh/turf-go/geojson"
)

func pt(lng, lat float64) geojson.Position {
	return geojson.Position{lng, lat}
}

func TestPolygonSmooth(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}},
	})
	result, err := PolygonSmooth(poly, 1)
	if err != nil {
		t.Fatal(err)
	}
	p, ok := result.Geometry.(*geojson.Polygon)
	if !ok {
		t.Fatal("expected Polygon")
	}
	ring := p.Coordinates[0]
	if len(ring) < 5 {
		t.Errorf("smoothed ring should have at least 5 points, got %d", len(ring))
	}
}

func TestPolygonSmoothMultiPolygon(t *testing.T) {
	mp := geojson.NewMultiPolygon([][][]geojson.Position{
		{{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}}},
	})
	result, err := PolygonSmooth(mp, 1)
	if err != nil {
		t.Fatal(err)
	}
	_, ok := result.Geometry.(*geojson.MultiPolygon)
	if !ok {
		t.Fatal("expected MultiPolygon")
	}
}

func TestPolygonSmoothMultipleIterations(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}},
	})
	result1, _ := PolygonSmooth(poly, 1)
	result2, _ := PolygonSmooth(poly, 2)
	p1 := result1.Geometry.(*geojson.Polygon)
	p2 := result2.Geometry.(*geojson.Polygon)
	if len(p1.Coordinates[0]) == len(p2.Coordinates[0]) {
		t.Error("iterations should produce different number of points")
	}
}

func TestPolygonSmoothClosedRing(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}},
	})
	result, err := PolygonSmooth(poly, 1)
	if err != nil {
		t.Fatal(err)
	}
	p := result.Geometry.(*geojson.Polygon)
	ring := p.Coordinates[0]
	first := ring[0]
	last := ring[len(ring)-1]
	if first[0] != last[0] || first[1] != last[1] {
		t.Error("smoothed ring should still be closed")
	}
}

func TestPolygonSmoothWithHole(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {20, 0}, {20, 20}, {0, 20}, {0, 0}},
		{{5, 5}, {15, 5}, {15, 15}, {5, 15}, {5, 5}},
	})
	result, err := PolygonSmooth(poly, 1)
	if err != nil {
		t.Fatal(err)
	}
	p := result.Geometry.(*geojson.Polygon)
	if len(p.Coordinates) != 2 {
		t.Errorf("expected 2 rings, got %d", len(p.Coordinates))
	}
}

func TestPolygonSmoothInvalidInput(t *testing.T) {
	pt := geojson.NewPoint(pt(0, 0))
	_, err := PolygonSmooth(pt, 1)
	if err == nil {
		t.Error("expected error for Point input")
	}
}
