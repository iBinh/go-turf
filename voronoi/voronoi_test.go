package voronoi

import (
	"testing"

	"github.com/ibinh/turf-go/geojson"
)

func TestVoronoiFourPoints(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{1, 0}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 1}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{1, 1}), nil),
	})

	result, err := Voronoi(fc, VoronoiOptions{
		BBox: []float64{-1, -1, 2, 2},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if len(result.Features) != 4 {
		t.Errorf("expected 4 cells, got %d", len(result.Features))
	}
	for i, f := range result.Features {
		if f.Geometry.Type() != geojson.TypePolygon {
			t.Errorf("cell %d: expected Polygon, got %s", i, f.Geometry.Type())
		}
	}
}

func TestVoronoiSinglePoint(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0.5, 0.5}), nil),
	})

	result, err := Voronoi(fc, VoronoiOptions{
		BBox: []float64{0, 0, 1, 1},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 1 {
		t.Errorf("expected 1 cell, got %d", len(result.Features))
	}
}

func TestVoronoiNoOptions(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
	})

	_, err := Voronoi(fc)
	if err == nil {
		t.Error("expected error when no bbox provided")
	}
}

func TestVoronoiNoPoints(t *testing.T) {
	fc := geojson.NewFeatureCollection(nil)

	_, err := Voronoi(fc, VoronoiOptions{
		BBox: []float64{0, 0, 1, 1},
	})
	if err == nil {
		t.Error("expected error when no points")
	}
}

func TestClipPolygonByHalfPlane(t *testing.T) {
	square := []geojson.Position{
		{0, 0},
		{1, 0},
		{1, 1},
		{0, 1},
		{0, 0},
	}

	a := geojson.Position{0.5, 0.5}
	b := geojson.Position{0.5, 1.5}

	result := clipPolygonByHalfPlane(square, a, b)
	if len(result) < 3 {
		t.Fatal("expected at least 3 points after clipping")
	}

	poly := geojson.NewPolygon([][]geojson.Position{result})
	if poly.Type() != geojson.TypePolygon {
		t.Errorf("expected Polygon type, got %s", poly.Type())
	}
}
