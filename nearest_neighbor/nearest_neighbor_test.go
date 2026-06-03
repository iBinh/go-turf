package nearestneighbor

import (
	"testing"

	"github.com/ibinh/turf-go/geojson"
)

func TestNearestNeighborClustered(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0.1, 0}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0.1}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 10}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{10.1, 10}), nil),
	})

	result, err := NearestNeighborAnalysis(fc)
	if err != nil {
		t.Fatal(err)
	}
	if result.R <= 0 {
		t.Errorf("expected positive R statistic, got %f", result.R)
	}
	if result.Z > 0 {
		t.Errorf("clustered pattern should have negative z-score, got %f", result.Z)
	}
}

func TestNearestNeighborRegular(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 0}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 10}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 10}), nil),
	})

	result, err := NearestNeighborAnalysis(fc)
	if err != nil {
		t.Fatal(err)
	}
	if result.R <= 0 {
		t.Errorf("expected positive R statistic, got %f", result.R)
	}
}

func TestNearestNeighborMinPointsError(t *testing.T) {
	_, err := NearestNeighborAnalysis(nil)
	if err == nil {
		t.Fatal("expected error for nil input")
	}

	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
	})
	_, err = NearestNeighborAnalysis(fc)
	if err == nil {
		t.Fatal("expected error for < 2 points")
	}
}
