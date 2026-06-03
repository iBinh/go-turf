package unkink

import (
	"testing"

	"github.com/ibinh/turf-go/geojson"
)

func TestUnkinkPolygonNoKinks(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}},
	})
	result, err := UnkinkPolygon(poly)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 1 {
		t.Errorf("expected 1 feature, got %d", len(result.Features))
	}
}

func TestUnkinkPolygonSelfIntersecting(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {10, 10}, {10, 0}, {0, 10}, {0, 0}},
	})
	result, err := UnkinkPolygon(poly)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) < 1 {
		t.Error("expected at least 1 feature from self-intersecting polygon")
	}
}

func TestUnkinkPolygonBowTie(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {10, 10}, {10, 0}, {0, 10}, {0, 0}},
		{{2, 2}, {8, 2}, {8, 8}, {2, 8}, {2, 2}},
	})
	result, err := UnkinkPolygon(poly)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) < 1 {
		t.Error("expected at least 1 feature")
	}
	for i, f := range result.Features {
		_, ok := f.Geometry.(*geojson.Polygon)
		if !ok {
			t.Errorf("feature %d: expected Polygon", i)
		}
	}
}

func TestUnkinkPolygonMultiPolygon(t *testing.T) {
	mp := geojson.NewMultiPolygon([][][]geojson.Position{
		{{{0, 0}, {10, 10}, {10, 0}, {0, 10}, {0, 0}}},
	})
	result, err := UnkinkPolygon(mp)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) < 1 {
		t.Error("expected at least 1 feature")
	}
}

func TestUnkinkPolygonInvalidInput(t *testing.T) {
	line := geojson.NewLineString([]geojson.Position{{0, 0}, {10, 10}})
	_, err := UnkinkPolygon(line)
	if err == nil {
		t.Error("expected error for LineString input")
	}
}

func TestUnkinkPolygonTriangle(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {10, 0}, {5, 10}, {0, 0}},
	})
	result, err := UnkinkPolygon(poly)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 1 {
		t.Errorf("expected 1 feature for simple triangle, got %d", len(result.Features))
	}
}
