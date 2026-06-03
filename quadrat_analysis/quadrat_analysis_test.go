package quadratanalysis

import (
	"math"
	"testing"

	"github.com/ibinh/turf-go/geojson"
)

func TestQuadratAnalysisRegular(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 0}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 10}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 10}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{5, 5}), nil),
	})

	result, err := QuadratAnalysis(fc, 3)
	if err != nil {
		t.Fatal(err)
	}
	if result.Chi2 < 0 {
		t.Errorf("expected non-negative chi2, got %f", result.Chi2)
	}
	if result.DF != 8 {
		t.Errorf("expected 8 degrees of freedom, got %d", result.DF)
	}
	if result.P < 0 || result.P > 1 {
		t.Errorf("expected p-value in [0,1], got %f", result.P)
	}
}

func TestQuadratAnalysisMinPointsError(t *testing.T) {
	_, err := QuadratAnalysis(nil, 3)
	if err == nil {
		t.Fatal("expected error for nil input")
	}

	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
	})
	_, err = QuadratAnalysis(fc, 3)
	if err == nil {
		t.Fatal("expected error for < 2 points")
	}
}

func TestQuadratAnalysisDefaultGridSize(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{1, 2}), nil),
	})

	result, err := QuadratAnalysis(fc, 0)
	if err != nil {
		t.Fatal(err)
	}
	if math.IsNaN(result.Chi2) {
		t.Errorf("expected valid chi2, got NaN")
	}
}
