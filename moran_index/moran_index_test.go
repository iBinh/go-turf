package moranindex

import (
	"math"
	"testing"

	"github.com/ibinh/turf-go/geojson"
)

func TestMoranIndexClustered(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), map[string]any{"v": float64(10)}),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{1, 0}), map[string]any{"v": float64(10)}),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 1}), map[string]any{"v": float64(10)}),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 10}), map[string]any{"v": float64(0)}),
	})

	mi, err := MoranIndex(fc, "v")
	if err != nil {
		t.Fatal(err)
	}
	if mi <= 0 {
		t.Errorf("expected positive Moran's I for clustered values, got %f", mi)
	}
}

func TestMoranIndexRandom(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), map[string]any{"v": float64(10)}),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 0}), map[string]any{"v": float64(0)}),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 10}), map[string]any{"v": float64(10)}),
	})

	mi, err := MoranIndex(fc, "v")
	if err != nil {
		t.Fatal(err)
	}
	if math.IsNaN(mi) {
		t.Errorf("expected valid Moran's I, got NaN")
	}
}

func TestMoranIndexConstantValues(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), map[string]any{"v": float64(5)}),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 0}), map[string]any{"v": float64(5)}),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 10}), map[string]any{"v": float64(5)}),
	})

	mi, err := MoranIndex(fc, "v")
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(mi) > 0.001 {
		t.Errorf("expected ~0 Moran's I for constant values, got %f", mi)
	}
}

func TestMoranIndexMinPointsError(t *testing.T) {
	_, err := MoranIndex(nil, "v")
	if err == nil {
		t.Fatal("expected error for nil input")
	}

	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), map[string]any{"v": float64(1)}),
	})
	_, err = MoranIndex(fc, "v")
	if err == nil {
		t.Fatal("expected error for < 3 points")
	}
}

func TestMoranIndexMissingProperty(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), map[string]any{"v": float64(1)}),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 0}), map[string]any{"v": float64(2)}),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 10}), map[string]any{"x": float64(3)}),
	})

	_, err := MoranIndex(fc, "v")
	if err == nil {
		t.Fatal("expected error for features missing property")
	}
}
