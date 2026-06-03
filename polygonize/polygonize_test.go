package polygonize

import (
	"testing"

	"github.com/ibinh/turf-go/geojson"
)

func TestPolygonizeSquare(t *testing.T) {
	line := geojson.NewFeature(
		geojson.NewLineString([]geojson.Position{
			{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0},
		}),
		nil,
	)

	result, err := Polygonize(line)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 1 {
		t.Errorf("expected 1 polygon, got %d", len(result.Features))
	}
}

func TestPolygonizeNotEnoughEdges(t *testing.T) {
	line := geojson.NewFeature(
		geojson.NewLineString([]geojson.Position{{0, 0}, {10, 0}}),
		nil,
	)
	_, err := Polygonize(line)
	if err == nil {
		t.Fatal("expected error for not enough edges")
	}
}

func TestPolygonizePolygonInput(t *testing.T) {
	poly := geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}}}),
		nil,
	)

	result, err := Polygonize(poly)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 1 {
		t.Errorf("expected 1 polygon, got %d", len(result.Features))
	}
}
