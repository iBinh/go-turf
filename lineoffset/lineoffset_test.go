package lineoffset

import (
	"math"
	"testing"

	"github.com/ibinh/turf-go/geojson"
	"github.com/ibinh/turf-go/measurement"
)

func pt(lng, lat float64) geojson.Position {
	return geojson.Position{lng, lat}
}

func TestLineOffsetSimple(t *testing.T) {
	line := geojson.NewLineString([]geojson.Position{{0, 0}, {10, 0}})
	result, err := LineOffset(line, 1000, measurement.UnitMeters)
	if err != nil {
		t.Fatal(err)
	}
	coords, _ := geojson.GetCoords(result)
	pts := coords.([]geojson.Position)
	if len(pts) != 2 {
		t.Fatalf("expected 2 points, got %d", len(pts))
	}
}

func TestLineOffsetNegative(t *testing.T) {
	line := geojson.NewLineString([]geojson.Position{{0, 0}, {10, 0}})
	result, err := LineOffset(line, -1000, measurement.UnitMeters)
	if err != nil {
		t.Fatal(err)
	}
	coords, _ := geojson.GetCoords(result)
	pts := coords.([]geojson.Position)
	if len(pts) != 2 {
		t.Fatalf("expected 2 points, got %d", len(pts))
	}
}

func TestLineOffsetThreePoints(t *testing.T) {
	line := geojson.NewLineString([]geojson.Position{{0, 0}, {5, 5}, {10, 0}})
	result, err := LineOffset(line, 500, measurement.UnitMeters)
	if err != nil {
		t.Fatal(err)
	}
	coords, _ := geojson.GetCoords(result)
	pts := coords.([]geojson.Position)
	if len(pts) < 3 {
		t.Errorf("expected at least 3 points, got %d", len(pts))
	}
}

func TestLineOffsetHorizontal(t *testing.T) {
	line := geojson.NewLineString([]geojson.Position{{0, 0}, {1, 0}})
	result, err := LineOffset(line, 100, measurement.UnitKilometers)
	if err != nil {
		t.Fatal(err)
	}
	coords, _ := geojson.GetCoords(result)
	pts := coords.([]geojson.Position)
	dist, _ := measurement.Distance(geojson.NewPoint(pts[0]), geojson.NewPoint(pt(0, 0)))
	if math.Abs(dist-100) > 10 {
		t.Errorf("expected offset ~100km, got %f", dist)
	}
}

func TestLineOffsetVertical(t *testing.T) {
	line := geojson.NewLineString([]geojson.Position{{0, 0}, {0, 1}})
	result, err := LineOffset(line, 100, measurement.UnitKilometers)
	if err != nil {
		t.Fatal(err)
	}
	coords, _ := geojson.GetCoords(result)
	pts := coords.([]geojson.Position)
	dist, _ := measurement.Distance(geojson.NewPoint(pts[0]), geojson.NewPoint(pt(0, 0)))
	if math.Abs(dist-100) > 10 {
		t.Errorf("expected offset ~100km, got %f", dist)
	}
}

func TestLineOffsetInvalidInput(t *testing.T) {
	pt := geojson.NewPoint(pt(0, 0))
	_, err := LineOffset(pt, 100, measurement.UnitMeters)
	if err == nil {
		t.Error("expected error for Point input")
	}
}

func TestLineOffsetUnits(t *testing.T) {
	line := geojson.NewLineString([]geojson.Position{{0, 0}, {1, 0}})
	result, err := LineOffset(line, 10, measurement.UnitKilometers)
	if err != nil {
		t.Fatal(err)
	}
	coords, _ := geojson.GetCoords(result)
	pts := coords.([]geojson.Position)
	if len(pts) != 2 {
		t.Errorf("expected 2 points, got %d", len(pts))
	}
}
