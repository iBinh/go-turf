package tangents

import (
	"testing"

	"github.com/ibinh/turf-go/geojson"
)

func pt(lng, lat float64) geojson.Position {
	return geojson.Position{lng, lat}
}

func TestPolygonTangentsSquare(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}},
	})
	point := geojson.NewPoint(pt(-5, 5))
	result, err := PolygonTangents(point, poly)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 2 {
		t.Fatalf("expected 2 features, got %d", len(result.Features))
	}
}

func TestPolygonTangentsOutsidePoint(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}},
	})
	point := geojson.NewPoint(pt(5, 20))
	result, err := PolygonTangents(point, poly)
	if err != nil {
		t.Fatal(err)
	}
	left := result.Features[0]
	right := result.Features[1]
	lCoord, _ := geojson.GetCoord(left)
	rCoord, _ := geojson.GetCoord(right)
	if lCoord[0] == rCoord[0] && lCoord[1] == rCoord[1] {
		t.Error("left and right tangents should be different points")
	}
}

func TestPolygonTangentsMultiPolygon(t *testing.T) {
	mp := geojson.NewMultiPolygon([][][]geojson.Position{
		{{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}}},
	})
	point := geojson.NewPoint(pt(-5, 5))
	result, err := PolygonTangents(point, mp)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 2 {
		t.Errorf("expected 2 features, got %d", len(result.Features))
	}
}

func TestPolygonTangentsInvalidInput(t *testing.T) {
	line := geojson.NewLineString([]geojson.Position{{0, 0}, {10, 10}})
	point := geojson.NewPoint(pt(0, 0))
	_, err := PolygonTangents(point, line)
	if err == nil {
		t.Error("expected error for LineString input")
	}
}

func TestPolygonTangentsInsidePoint(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}},
	})
	point := geojson.NewPoint(pt(5, 5))
	result, err := PolygonTangents(point, poly)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 2 {
		t.Errorf("expected 2 features, got %d", len(result.Features))
	}
}
