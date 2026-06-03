package distanceweight

import (
	"testing"

	"github.com/ibinh/turf-go/geojson"
)

func TestDistanceWeightTwoPoints(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{3, 4}), nil),
	})

	w, err := DistanceWeight(fc)
	if err != nil {
		t.Fatal(err)
	}
	if len(w) != 2 || len(w[0]) != 2 {
		t.Fatal("expected 2x2 matrix")
	}
	if w[0][0] != 0 {
		t.Errorf("expected diagonal 0, got %f", w[0][0])
	}
	if w[0][1] <= 0 {
		t.Errorf("expected positive weight between different points, got %f", w[0][1])
	}
}

func TestDistanceWeightSymmetry(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{1, 0}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 1}), nil),
	})

	w, err := DistanceWeight(fc)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if w[i][j] != w[j][i] {
				t.Errorf("expected symmetry at [%d][%d], got %f vs %f", i, j, w[i][j], w[j][i])
			}
		}
	}
}

func TestDistanceWeightThreshold(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{1, 0}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 0}), nil),
	})

	w, err := DistanceWeight(fc, 5)
	if err != nil {
		t.Fatal(err)
	}
	if w[0][2] != 0 {
		t.Errorf("expected 0 weight beyond threshold, got %f", w[0][2])
	}
	if w[0][1] <= 0 {
		t.Errorf("expected positive weight within threshold, got %f", w[0][1])
	}
}

func TestDistanceWeightMinPointsError(t *testing.T) {
	_, err := DistanceWeight(nil)
	if err == nil {
		t.Fatal("expected error for nil input")
	}

	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
	})
	_, err = DistanceWeight(fc)
	if err == nil {
		t.Fatal("expected error for < 2 points")
	}
}
