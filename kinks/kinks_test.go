package kinks

import (
	"testing"

	"github.com/ibinh/turf-go/geojson"
)

func TestKinksLineStringNoKinks(t *testing.T) {
	line := geojson.NewLineString([]geojson.Position{{0, 0}, {10, 10}, {20, 0}})
	result, err := Kinks(line)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 0 {
		t.Errorf("expected no kinks, got %d", len(result.Features))
	}
}

func TestKinksLineStringSelfIntersect(t *testing.T) {
	line := geojson.NewLineString([]geojson.Position{
		{0, 0}, {10, 10}, {10, 0}, {0, 10},
	})
	result, err := Kinks(line)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 1 {
		t.Errorf("expected 1 kink, got %d", len(result.Features))
	}
	if len(result.Features) > 0 {
		coord, _ := geojson.GetCoord(result.Features[0])
		if coord[0] != 5 || coord[1] != 5 {
			t.Errorf("expected kink at (5,5), got (%v,%v)", coord[0], coord[1])
		}
	}
}

func TestKinksPolygonNoKinks(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}},
	})
	result, err := Kinks(poly)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 0 {
		t.Errorf("expected no kinks, got %d", len(result.Features))
	}
}

func TestKinksPolygonSelfIntersect(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {10, 10}, {10, 0}, {0, 10}, {0, 0}},
	})
	result, err := Kinks(poly)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 1 {
		t.Errorf("expected 1 kink, got %d", len(result.Features))
	}
}

func TestKinksMultiLineString(t *testing.T) {
	ml := geojson.NewMultiLineString([][]geojson.Position{
		{{0, 0}, {10, 10}, {10, 0}, {0, 10}, {0, 0}},
	})
	result, err := Kinks(ml)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 1 {
		t.Errorf("expected 1 kink from self-intersecting MLS, got %d", len(result.Features))
	}
}

func TestKinksMultiPolygon(t *testing.T) {
	mp := geojson.NewMultiPolygon([][][]geojson.Position{
		{{{0, 0}, {5, 0}, {5, 5}, {0, 5}, {0, 0}}},
		{{{3, 3}, {8, 3}, {8, 8}, {3, 8}, {3, 3}}},
	})
	result, err := Kinks(mp)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 0 {
		t.Errorf("expected 0 kinks, got %d", len(result.Features))
	}
}

func TestKinksFeature(t *testing.T) {
	f := geojson.NewFeature(geojson.NewLineString([]geojson.Position{
		{0, 0}, {10, 10}, {10, 0}, {0, 10},
	}), nil)
	result, err := Kinks(f)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 1 {
		t.Errorf("expected 1 kink from feature input, got %d", len(result.Features))
	}
}

func TestKinksAdjacentEdges(t *testing.T) {
	line := geojson.NewLineString([]geojson.Position{
		{0, 0}, {5, 5}, {10, 0},
	})
	result, err := Kinks(line)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 0 {
		t.Errorf("adjacent edges should not be reported as kinks, got %d", len(result.Features))
	}
}
