package interpolation

import (
	"math"
	"testing"

	"github.com/ibinh/turf-go/geojson"
	"github.com/ibinh/turf-go/measurement"
)

func TestSample(t *testing.T) {
	features := make([]*geojson.Feature, 10)
	for i := 0; i < 10; i++ {
		features[i] = geojson.NewFeature(geojson.NewPoint(geojson.Position{float64(i), float64(i)}), nil)
	}
	fc := geojson.NewFeatureCollection(features)

	sampled, err := Sample(fc, 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(sampled.Features) != 3 {
		t.Errorf("expected 3 samples, got %d", len(sampled.Features))
	}
}

func TestSampleEmpty(t *testing.T) {
	fc := geojson.NewFeatureCollection(nil)
	sampled, err := Sample(fc, 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(sampled.Features) != 0 {
		t.Error("expected empty result")
	}
}

func TestSampleMoreThanAvailable(t *testing.T) {
	features := make([]*geojson.Feature, 3)
	for i := 0; i < 3; i++ {
		features[i] = geojson.NewFeature(geojson.NewPoint(geojson.Position{float64(i), float64(i)}), nil)
	}
	fc := geojson.NewFeatureCollection(features)

	sampled, err := Sample(fc, 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(sampled.Features) != 3 {
		t.Errorf("expected 3 features (clamped), got %d", len(sampled.Features))
	}
}

func TestSampleNilInput(t *testing.T) {
	_, err := Sample(nil, 5)
	if err == nil {
		t.Error("expected error for nil input")
	}
}

func TestTin(t *testing.T) {
	features := make([]*geojson.Feature, 4)
	features[0] = geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil)
	features[1] = geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 0}), nil)
	features[2] = geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 10}), nil)
	features[3] = geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 10}), nil)
	fc := geojson.NewFeatureCollection(features)

	tin, err := Tin(fc)
	if err != nil {
		t.Fatal(err)
	}
	if len(tin.Features) < 2 {
		t.Errorf("expected at least 2 triangles, got %d", len(tin.Features))
	}
	for _, f := range tin.Features {
		if _, ok := f.Geometry.(*geojson.Polygon); !ok {
			t.Errorf("expected Polygon, got %T", f.Geometry)
		}
	}
}

func TestTinTooFewPoints(t *testing.T) {
	features := make([]*geojson.Feature, 2)
	features[0] = geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil)
	features[1] = geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 10}), nil)
	fc := geojson.NewFeatureCollection(features)

	_, err := Tin(fc)
	if err == nil {
		t.Error("expected error for fewer than 3 points")
	}
}

func TestTinNullInput(t *testing.T) {
	_, err := Tin(nil)
	if err == nil {
		t.Error("expected error for nil input")
	}
}

func TestTinNonTriangular(t *testing.T) {
	features := make([]*geojson.Feature, 3)
	features[0] = geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil)
	features[1] = geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 0}), nil)
	features[2] = geojson.NewFeature(geojson.NewPoint(geojson.Position{5, 10}), nil)
	fc := geojson.NewFeatureCollection(features)

	tin, err := Tin(fc)
	if err != nil {
		t.Fatal(err)
	}
	if len(tin.Features) != 1 {
		t.Errorf("expected 1 triangle, got %d", len(tin.Features))
	}
}

func TestInterpolate(t *testing.T) {
	features := make([]*geojson.Feature, 4)
	features[0] = geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), map[string]any{"z": 10.0})
	features[1] = geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 0}), map[string]any{"z": 20.0})
	features[2] = geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 10}), map[string]any{"z": 30.0})
	features[3] = geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 10}), map[string]any{"z": 40.0})
	fc := geojson.NewFeatureCollection(features)

	result, err := Interpolate(fc, 5, measurement.UnitDegrees, "z")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) == 0 {
		t.Error("expected interpolated points")
	}
	for _, f := range result.Features {
		if f.Properties == nil {
			t.Error("expected properties")
			continue
		}
		z, ok := f.Properties["z"]
		if !ok {
			t.Error("expected z property")
			continue
		}
		_, ok = z.(float64)
		if !ok {
			t.Errorf("expected float64 z, got %T", z)
		}
	}
}

func TestInterpolateTooFewPoints(t *testing.T) {
	features := make([]*geojson.Feature, 1)
	features[0] = geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), map[string]any{"z": 10.0})
	fc := geojson.NewFeatureCollection(features)

	_, err := Interpolate(fc, 5, measurement.UnitDegrees, "z")
	if err == nil {
		t.Error("expected error for too few points")
	}
}

func TestInterpolateExactMatch(t *testing.T) {
	features := make([]*geojson.Feature, 2)
	features[0] = geojson.NewFeature(geojson.NewPoint(geojson.Position{5, 5}), map[string]any{"z": 42.0})
	features[1] = geojson.NewFeature(geojson.NewPoint(geojson.Position{5.001, 5.001}), map[string]any{"z": 42.0})
	fc := geojson.NewFeatureCollection(features)

	result, err := Interpolate(fc, 1, measurement.UnitDegrees, "z")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) > 0 {
		z := result.Features[0].Properties["z"].(float64)
		if math.Abs(z-42) > 0.1 {
			t.Errorf("expected z ~42, got %f", z)
		}
	}
}

func TestInterpolateCustomWeight(t *testing.T) {
	features := make([]*geojson.Feature, 3)
	features[0] = geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), map[string]any{"v": 10.0})
	features[1] = geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 0}), map[string]any{"v": 20.0})
	features[2] = geojson.NewFeature(geojson.NewPoint(geojson.Position{5, 10}), map[string]any{"v": 30.0})
	fc := geojson.NewFeatureCollection(features)

	result, err := Interpolate(fc, 5, measurement.UnitDegrees, "v", InterpolateOptions{Weight: 3})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) == 0 {
		t.Error("expected results")
	}
}

func TestPlanarDistance(t *testing.T) {
	d := PlanarDistance(geojson.Position{0, 0}, geojson.Position{3, 4})
	if math.Abs(d-5) > 0.001 {
		t.Errorf("expected 5, got %f", d)
	}
}

func TestPlanarPointOnLine(t *testing.T) {
	on := PlanarPointOnLine(geojson.Position{5, 5}, geojson.Position{0, 0}, geojson.Position{10, 10})
	if !on {
		t.Error("point on line should be true")
	}

	off := PlanarPointOnLine(geojson.Position{5, 6}, geojson.Position{0, 0}, geojson.Position{10, 10})
	if off {
		t.Error("point off line should be false")
	}

	end := PlanarPointOnLine(geojson.Position{10, 10}, geojson.Position{0, 0}, geojson.Position{10, 10})
	if !end {
		t.Error("endpoint should be on line")
	}

	beyond := PlanarPointOnLine(geojson.Position{15, 15}, geojson.Position{0, 0}, geojson.Position{10, 10})
	if beyond {
		t.Error("point beyond segment should be off line")
	}
}

func TestInterpolateMissingProperty(t *testing.T) {
	features := make([]*geojson.Feature, 2)
	features[0] = geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), map[string]any{"a": 10.0})
	features[1] = geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 10}), map[string]any{"a": 20.0})
	fc := geojson.NewFeatureCollection(features)

	result, err := Interpolate(fc, 5, measurement.UnitDegrees, "b")
	if err == nil {
		t.Error("expected error for missing property")
		_ = result
	}
}

func TestSampleDeterministicCount(t *testing.T) {
	features := make([]*geojson.Feature, 100)
	for i := 0; i < 100; i++ {
		features[i] = geojson.NewFeature(geojson.NewPoint(geojson.Position{float64(i), float64(i)}), nil)
	}
	fc := geojson.NewFeatureCollection(features)

	sampled, err := Sample(fc, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(sampled.Features) != 10 {
		t.Errorf("expected 10, got %d", len(sampled.Features))
	}
}

func BenchmarkTin(b *testing.B) {
	features := make([]*geojson.Feature, 100)
	for i := 0; i < 100; i++ {
		x := float64(i % 10)
		y := float64(i / 10)
		features[i] = geojson.NewFeature(geojson.NewPoint(geojson.Position{x, y}), nil)
	}
	fc := geojson.NewFeatureCollection(features)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Tin(fc)
	}
}
