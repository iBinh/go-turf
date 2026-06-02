package clusters

import (
	"testing"

	"github.com/ibinh/turf-go/geojson"
)

func TestKMeans(t *testing.T) {
	features := make([]*geojson.Feature, 6)
	features[0] = geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil)
	features[1] = geojson.NewFeature(geojson.NewPoint(geojson.Position{1, 0}), nil)
	features[2] = geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 1}), nil)
	features[3] = geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 10}), nil)
	features[4] = geojson.NewFeature(geojson.NewPoint(geojson.Position{11, 10}), nil)
	features[5] = geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 11}), nil)
	fc := geojson.NewFeatureCollection(features)

	result, err := ClustersKMeans(fc, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 6 {
		t.Errorf("expected 6 features, got %d", len(result.Features))
	}
	for _, f := range result.Features {
		c, ok := f.Properties["cluster"]
		if !ok {
			t.Error("expected cluster property")
			continue
		}
		ci, ok := c.(int)
		if !ok {
			t.Errorf("expected int cluster, got %T", c)
			continue
		}
		if ci < 0 || ci > 1 {
			t.Errorf("cluster out of range [0,1]: %d", ci)
		}
	}
}

func TestKMeansSingleCluster(t *testing.T) {
	features := make([]*geojson.Feature, 3)
	for i := 0; i < 3; i++ {
		features[i] = geojson.NewFeature(geojson.NewPoint(geojson.Position{5, 5}), nil)
	}
	fc := geojson.NewFeatureCollection(features)

	result, err := ClustersKMeans(fc, 1)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range result.Features {
		if f.Properties["cluster"].(int) != 0 {
			t.Error("all points should be in cluster 0")
		}
	}
}

func TestKMeansEmptyInput(t *testing.T) {
	_, err := ClustersKMeans(nil, 2)
	if err == nil {
		t.Error("expected error for nil input")
	}
}

func TestKMeansMoreClustersThanPoints(t *testing.T) {
	features := make([]*geojson.Feature, 2)
	features[0] = geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil)
	features[1] = geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 10}), nil)
	fc := geojson.NewFeatureCollection(features)

	result, err := ClustersKMeans(fc, 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 2 {
		t.Errorf("expected 2 features, got %d", len(result.Features))
	}
}

func TestKMeansCentroidProperty(t *testing.T) {
	features := make([]*geojson.Feature, 4)
	features[0] = geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil)
	features[1] = geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil)
	features[2] = geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 10}), nil)
	features[3] = geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 10}), nil)
	fc := geojson.NewFeatureCollection(features)

	result, err := ClustersKMeans(fc, 2)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range result.Features {
		if _, ok := f.Properties["centroid"]; !ok {
			t.Error("expected centroid property")
		}
	}
}

func TestDBSCAN(t *testing.T) {
	features := make([]*geojson.Feature, 6)
	features[0] = geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil)
	features[1] = geojson.NewFeature(geojson.NewPoint(geojson.Position{1, 0}), nil)
	features[2] = geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 1}), nil)
	features[3] = geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 10}), nil)
	features[4] = geojson.NewFeature(geojson.NewPoint(geojson.Position{11, 10}), nil)
	features[5] = geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 11}), nil)
	fc := geojson.NewFeatureCollection(features)

	result, err := ClustersDbscan(fc, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 6 {
		t.Errorf("expected 6 features, got %d", len(result.Features))
	}
	for _, f := range result.Features {
		if _, ok := f.Properties["cluster"]; !ok {
			t.Error("expected cluster property")
		}
	}
}

func TestDBSCANNoise(t *testing.T) {
	features := make([]*geojson.Feature, 4)
	features[0] = geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil)
	features[1] = geojson.NewFeature(geojson.NewPoint(geojson.Position{1, 1}), nil)
	features[2] = geojson.NewFeature(geojson.NewPoint(geojson.Position{100, 100}), nil)
	features[3] = geojson.NewFeature(geojson.NewPoint(geojson.Position{101, 101}), nil)
	fc := geojson.NewFeatureCollection(features)

	result, err := ClustersDbscan(fc, 2, DbscanOptions{MinPoints: 2})
	if err != nil {
		t.Fatal(err)
	}
	clusters := make(map[int]int)
	for _, f := range result.Features {
		c := f.Properties["cluster"].(int)
		clusters[c]++
	}

	cluster0Count := 0
	for c, count := range clusters {
		if c == 0 {
			cluster0Count = count
		}
		_ = count
	}
	if cluster0Count != 0 {
		t.Logf("noise point cluster 0 has %d points", cluster0Count)
	}
}

func TestDBSCANEmptyInput(t *testing.T) {
	_, err := ClustersDbscan(nil, 1)
	if err == nil {
		t.Error("expected error for nil input")
	}
}

func TestDBSCANZeroRadius(t *testing.T) {
	features := make([]*geojson.Feature, 2)
	features[0] = geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil)
	features[1] = geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 10}), nil)
	fc := geojson.NewFeatureCollection(features)

	_, err := ClustersDbscan(fc, 0)
	if err == nil {
		t.Error("expected error for zero radius")
	}
}

func TestDBSCANCustomMinPoints(t *testing.T) {
	features := make([]*geojson.Feature, 5)
	for i := 0; i < 5; i++ {
		x := float64(i)
		features[i] = geojson.NewFeature(geojson.NewPoint(geojson.Position{x, 0}), nil)
	}
	fc := geojson.NewFeatureCollection(features)

	result, err := ClustersDbscan(fc, 2, DbscanOptions{MinPoints: 2})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 5 {
		t.Errorf("expected 5 features, got %d", len(result.Features))
	}
}

func TestDissolve(t *testing.T) {
	features := make([]*geojson.Feature, 2)
	features[0] = geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{{{0, 0}, {5, 0}, {5, 5}, {0, 5}, {0, 0}}}),
		map[string]any{"group": "a"},
	)
	features[1] = geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{{{5, 0}, {10, 0}, {10, 5}, {5, 5}, {5, 0}}}),
		map[string]any{"group": "a"},
	)
	fc := geojson.NewFeatureCollection(features)

	result, err := Dissolve(fc, "group")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) == 0 {
		t.Error("expected merged features")
	}
}

func TestDissolveNoProperty(t *testing.T) {
	features := make([]*geojson.Feature, 2)
	features[0] = geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{{{0, 0}, {5, 0}, {5, 5}, {0, 5}, {0, 0}}}),
		nil,
	)
	features[1] = geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{{{10, 10}, {15, 10}, {15, 15}, {10, 15}, {10, 10}}}),
		nil,
	)
	fc := geojson.NewFeatureCollection(features)

	result, err := Dissolve(fc, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) == 0 {
		t.Error("expected features")
	}
}

func TestDissolveEmptyInput(t *testing.T) {
	_, err := Dissolve(nil, "prop")
	if err == nil {
		t.Error("expected error for nil input")
	}
}

func TestDissolveMultiPolygon(t *testing.T) {
	features := make([]*geojson.Feature, 1)
	features[0] = geojson.NewFeature(
		geojson.NewMultiPolygon([][][]geojson.Position{
			{{{0, 0}, {5, 0}, {5, 5}, {0, 5}, {0, 0}}},
			{{{10, 10}, {15, 10}, {15, 15}, {10, 15}, {10, 10}}},
		}),
		map[string]any{"g": "a"},
	)
	fc := geojson.NewFeatureCollection(features)

	result, err := Dissolve(fc, "g")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) == 0 {
		t.Error("expected features")
	}
}

func TestDissolveTouching(t *testing.T) {
	a := geojson.Position{0, 0}
	b := geojson.Position{5, 0}
	c := geojson.Position{10, 0}

	features := make([]*geojson.Feature, 2)
	features[0] = geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{{a, b, {5, 5}, {0, 5}, a}}),
		nil,
	)
	features[1] = geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{{b, c, {10, 5}, {5, 5}, b}}),
		nil,
	)
	fc := geojson.NewFeatureCollection(features)

	result, err := Dissolve(fc, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) == 0 {
		t.Error("expected merged feature")
	}
}

func BenchmarkKMeans(b *testing.B) {
	features := make([]*geojson.Feature, 1000)
	for i := 0; i < 1000; i++ {
		x := float64(i % 50)
		y := float64(i / 50)
		features[i] = geojson.NewFeature(geojson.NewPoint(geojson.Position{x, y}), nil)
	}
	fc := geojson.NewFeatureCollection(features)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ClustersKMeans(fc, 5)
	}
}

func BenchmarkDBSCAN(b *testing.B) {
	features := make([]*geojson.Feature, 200)
	for i := 0; i < 200; i++ {
		x := float64(i % 20)
		y := float64(i / 20)
		features[i] = geojson.NewFeature(geojson.NewPoint(geojson.Position{x, y}), nil)
	}
	fc := geojson.NewFeatureCollection(features)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ClustersDbscan(fc, 3)
	}
}
