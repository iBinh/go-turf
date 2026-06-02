package buffer

import (
	"testing"

	"github.com/ibinh/turf-go/geojson"
	"github.com/ibinh/turf-go/measurement"
)

func isPolygon(t *testing.T, geomType string) {
	t.Helper()
	if geomType != geojson.TypePolygon && geomType != geojson.TypeMultiPolygon {
		t.Errorf("expected Polygon or MultiPolygon, got %s", geomType)
	}
}

func TestBufferPoint(t *testing.T) {
	pt := geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil)

	result, err := Buffer(pt, 1, measurement.UnitDegrees)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal("expected a feature")
	}
	isPolygon(t, result.Geometry.Type())
}

func TestBufferPointWithSteps(t *testing.T) {
	pt := geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 20}), nil)

	result, err := Buffer(pt, 5, measurement.UnitKilometers, 8)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal("expected a feature")
	}
	isPolygon(t, result.Geometry.Type())
}

func TestBufferLineString(t *testing.T) {
	line := geojson.NewFeature(
		geojson.NewLineString([]geojson.Position{{0, 0}, {1, 0}}),
		nil,
	)

	result, err := Buffer(line, 0.1, measurement.UnitDegrees)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal("expected a feature")
	}
	isPolygon(t, result.Geometry.Type())
}

func TestBufferLineStringMultiSegment(t *testing.T) {
	line := geojson.NewFeature(
		geojson.NewLineString([]geojson.Position{{0, 0}, {1, 0}, {1, 1}}),
		nil,
	)

	result, err := Buffer(line, 0.05, measurement.UnitDegrees)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal("expected a feature")
	}
	isPolygon(t, result.Geometry.Type())
}

func TestBufferPolygon(t *testing.T) {
	poly := geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{
			{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}},
		}),
		nil,
	)

	result, err := Buffer(poly, 0.1, measurement.UnitDegrees)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal("expected a feature")
	}
	isPolygon(t, result.Geometry.Type())
}

func TestBufferLineStringZeroRadius(t *testing.T) {
	line := geojson.NewFeature(
		geojson.NewLineString([]geojson.Position{{0, 0}, {1, 0}}),
		nil,
	)

	_, err := Buffer(line, 0, measurement.UnitDegrees)
	if err == nil {
		t.Error("expected error for zero radius")
	}
}

func TestSegmentBuffer(t *testing.T) {
	p1 := geojson.Position{0, 0}
	p2 := geojson.Position{1, 0}

	poly := segmentBuffer(p1, p2, 0.5)
	if poly == nil {
		t.Fatal("expected non-nil polygon")
	}
	if len(poly.Coordinates) != 1 {
		t.Errorf("expected 1 ring, got %d", len(poly.Coordinates))
	}
	if len(poly.Coordinates[0]) != 5 {
		t.Errorf("expected 5 points in ring, got %d", len(poly.Coordinates[0]))
	}
}

func TestVertexCircle(t *testing.T) {
	center := geojson.Position{0, 0}
	poly := vertexCircle(center, 0.5, 8)
	if poly == nil {
		t.Fatal("expected non-nil polygon")
	}
	if len(poly.Coordinates[0]) != 9 {
		t.Errorf("expected 9 points (8 steps + closing), got %d", len(poly.Coordinates[0]))
	}
}
