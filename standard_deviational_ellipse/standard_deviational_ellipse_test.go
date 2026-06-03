package standarddeviationalellipse

import (
	"math"
	"testing"

	"github.com/ibinh/turf-go/geojson"
)

func TestStandardDeviationalEllipseFourPoints(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 0}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 10}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 10}), nil),
	})

	result, err := StandardDeviationalEllipse(fc)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}

	if math.Abs(result.Center[0]-5) > 0.01 || math.Abs(result.Center[1]-5) > 0.01 {
		t.Errorf("expected center at [5,5], got %v", result.Center)
	}
	if result.XAxis <= 0 || result.YAxis <= 0 {
		t.Errorf("expected positive axes, got x=%f y=%f", result.XAxis, result.YAxis)
	}
	if result.Polygon == nil {
		t.Fatal("expected polygon feature")
	}
}

func TestStandardDeviationalEllipseElongated(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{5, 0}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 0}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{15, 0}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{5, 5}), nil),
	})

	result, err := StandardDeviationalEllipse(fc)
	if err != nil {
		t.Fatal(err)
	}
	if result.XAxis <= result.YAxis {
		t.Errorf("expected x-axis larger than y-axis for horizontally elongated points")
	}
}

func TestStandardDeviationalEllipseMinPointsError(t *testing.T) {
	_, err := StandardDeviationalEllipse(nil)
	if err == nil {
		t.Fatal("expected error for nil input")
	}

	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
	})
	_, err = StandardDeviationalEllipse(fc)
	if err == nil {
		t.Fatal("expected error for < 3 points")
	}
}

func TestStandardDeviationalEllipseCustomSteps(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 0}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{5, 10}), nil),
	})

	result, err := StandardDeviationalEllipse(fc, 32)
	if err != nil {
		t.Fatal(err)
	}
	if result.Polygon == nil {
		t.Fatal("expected polygon feature")
	}
}
