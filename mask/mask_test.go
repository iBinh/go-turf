package mask

import (
	"testing"

	"github.com/ibinh/turf-go/geojson"
)

func pt(lng, lat float64) geojson.Position {
	return geojson.Position{lng, lat}
}

func TestMaskBasic(t *testing.T) {
	outer := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {20, 0}, {20, 20}, {0, 20}, {0, 0}},
	})
	inner := geojson.NewPolygon([][]geojson.Position{
		{{5, 5}, {15, 5}, {15, 15}, {5, 15}, {5, 5}},
	})
	result, err := Mask(outer, inner)
	if err != nil {
		t.Fatal(err)
	}
	poly, ok := result.Geometry.(*geojson.Polygon)
	if !ok {
		t.Fatal("expected Polygon")
	}
	if len(poly.Coordinates) != 2 {
		t.Errorf("expected 2 rings, got %d", len(poly.Coordinates))
	}
}

func TestMaskMultiPolygonInner(t *testing.T) {
	outer := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {30, 0}, {30, 30}, {0, 30}, {0, 0}},
	})
	mp := geojson.NewMultiPolygon([][][]geojson.Position{
		{{{5, 5}, {10, 5}, {10, 10}, {5, 10}, {5, 5}}},
		{{{20, 20}, {25, 20}, {25, 25}, {20, 25}, {20, 20}}},
	})
	result, err := Mask(outer, mp)
	if err != nil {
		t.Fatal(err)
	}
	poly, ok := result.Geometry.(*geojson.Polygon)
	if !ok {
		t.Fatal("expected Polygon")
	}
	if len(poly.Coordinates) != 3 {
		t.Errorf("expected 3 rings (1 outer + 2 holes), got %d", len(poly.Coordinates))
	}
}

func TestMaskOuterIsFeature(t *testing.T) {
	outer := geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{
			{{0, 0}, {20, 0}, {20, 20}, {0, 20}, {0, 0}},
		}),
		nil,
	)
	inner := geojson.NewPolygon([][]geojson.Position{
		{{5, 5}, {15, 5}, {15, 15}, {5, 15}, {5, 5}},
	})
	result, err := Mask(outer, inner)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestMaskNoInnerInside(t *testing.T) {
	outer := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {20, 0}, {20, 20}, {0, 20}, {0, 0}},
	})
	inner := geojson.NewPolygon([][]geojson.Position{
		{{100, 100}, {110, 100}, {110, 110}, {100, 110}, {100, 100}},
	})
	result, err := Mask(outer, inner)
	if err != nil {
		t.Fatal(err)
	}
	poly, ok := result.Geometry.(*geojson.Polygon)
	if !ok {
		t.Fatal("expected Polygon")
	}
	if len(poly.Coordinates) != 1 {
		t.Errorf("expected 1 ring when inner is outside, got %d", len(poly.Coordinates))
	}
}

func TestMaskInvalidInput(t *testing.T) {
	outer := geojson.NewPoint(pt(0, 0))
	inner := geojson.NewPolygon([][]geojson.Position{
		{{5, 5}, {15, 5}, {15, 15}, {5, 15}, {5, 5}},
	})
	_, err := Mask(outer, inner)
	if err == nil {
		t.Error("expected error for Point outer")
	}
}
