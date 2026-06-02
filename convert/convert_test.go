package convert

import (
	"testing"

	"github.com/ibinh/turf-go/geojson"
)

func pt(lng, lat float64) geojson.Position {
	return geojson.Position{lng, lat}
}

func TestPolygonToLine(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}},
	})
	result, err := PolygonToLine(poly)
	if err != nil {
		t.Fatal(err)
	}
	ls, ok := result.Geometry.(*geojson.LineString)
	if !ok {
		t.Fatal("expected LineString")
	}
	if len(ls.Coordinates) != 5 {
		t.Errorf("expected 5 coords, got %d", len(ls.Coordinates))
	}
}

func TestPolygonToLineMultiPolygon(t *testing.T) {
	mp := geojson.NewMultiPolygon([][][]geojson.Position{
		{{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}}},
		{{{20, 20}, {30, 20}, {30, 30}, {20, 30}, {20, 20}}},
	})
	result, err := PolygonToLine(mp)
	if err != nil {
		t.Fatal(err)
	}
	mls, ok := result.Geometry.(*geojson.MultiLineString)
	if !ok {
		t.Fatal("expected MultiLineString for MultiPolygon with multiple polys")
	}
	if len(mls.Coordinates) != 2 {
		t.Errorf("expected 2 lines, got %d", len(mls.Coordinates))
	}
}

func TestLineToPolygon(t *testing.T) {
	line := geojson.NewLineString([]geojson.Position{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}})
	result, err := LineToPolygon(line)
	if err != nil {
		t.Fatal(err)
	}
	poly, ok := result.Geometry.(*geojson.Polygon)
	if !ok {
		t.Fatal("expected Polygon")
	}
	if len(poly.Coordinates) != 1 {
		t.Errorf("expected 1 ring, got %d", len(poly.Coordinates))
	}
	if len(poly.Coordinates[0]) != 5 {
		t.Errorf("expected 5 coords in closed ring, got %d", len(poly.Coordinates[0]))
	}
}

func TestLineToPolygonAutoClose(t *testing.T) {
	line := geojson.NewLineString([]geojson.Position{{0, 0}, {10, 0}, {10, 10}, {0, 10}})
	result, err := LineToPolygon(line)
	if err != nil {
		t.Fatal(err)
	}
	poly, ok := result.Geometry.(*geojson.Polygon)
	if !ok {
		t.Fatal("expected Polygon")
	}
	ring := poly.Coordinates[0]
	first := ring[0]
	last := ring[len(ring)-1]
	if first[0] != last[0] || first[1] != last[1] {
		t.Error("ring was not closed")
	}
}

func TestLineToPolygonMultiLineString(t *testing.T) {
	ml := geojson.NewMultiLineString([][]geojson.Position{
		{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}},
		{{2, 2}, {8, 2}, {8, 8}, {2, 8}, {2, 2}},
	})
	result, err := LineToPolygon(ml)
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

func TestLineToPolygonTooFewPoints(t *testing.T) {
	line := geojson.NewLineString([]geojson.Position{{0, 0}, {10, 0}})
	_, err := LineToPolygon(line)
	if err == nil {
		t.Error("expected error for line with < 3 points")
	}
}
