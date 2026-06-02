package data

import (
	"testing"

	"github.com/ibinh/turf-go/geojson"
)

func TestTag(t *testing.T) {
	poly := geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}}}),
		map[string]any{"zone": "A"},
	)
	polygons := geojson.NewFeatureCollection([]*geojson.Feature{poly})

	pt := geojson.NewFeature(geojson.NewPoint(geojson.Position{5, 5}), nil)
	points := geojson.NewFeatureCollection([]*geojson.Feature{pt})

	result, err := Tag(points, polygons, "zone", "tagged")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 1 {
		t.Fatalf("expected 1 feature, got %d", len(result.Features))
	}
	if result.Features[0].Properties["tagged"] != "A" {
		t.Errorf("expected tagged=A, got %v", result.Features[0].Properties["tagged"])
	}
}

func TestTagOutsidePoint(t *testing.T) {
	poly := geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}}}),
		map[string]any{"zone": "A"},
	)
	polygons := geojson.NewFeatureCollection([]*geojson.Feature{poly})
	pt := geojson.NewFeature(geojson.NewPoint(geojson.Position{20, 20}), nil)
	points := geojson.NewFeatureCollection([]*geojson.Feature{pt})

	result, err := Tag(points, polygons, "zone", "tagged")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := result.Features[0].Properties["tagged"]; ok {
		t.Error("outside point should not be tagged")
	}
}

func TestTagMultiplePolygons(t *testing.T) {
	poly1 := geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{{{0, 0}, {0, 5}, {5, 5}, {5, 0}, {0, 0}}}),
		map[string]any{"zone": "A"},
	)
	poly2 := geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{{{5, 0}, {5, 5}, {10, 5}, {10, 0}, {5, 0}}}),
		map[string]any{"zone": "B"},
	)
	polygons := geojson.NewFeatureCollection([]*geojson.Feature{poly1, poly2})
	pt := geojson.NewFeature(geojson.NewPoint(geojson.Position{7, 2}), nil)
	points := geojson.NewFeatureCollection([]*geojson.Feature{pt})

	result, err := Tag(points, polygons, "zone", "tagged")
	if err != nil {
		t.Fatal(err)
	}
	if result.Features[0].Properties["tagged"] != "B" {
		t.Errorf("expected tagged=B, got %v", result.Features[0].Properties["tagged"])
	}
}

func TestTagNillInput(t *testing.T) {
	_, err := Tag(nil, nil, "", "")
	if err == nil {
		t.Error("expected error for nil input")
	}
}

func TestTagPreservesProperties(t *testing.T) {
	poly := geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}}}),
		map[string]any{"zone": "A"},
	)
	polygons := geojson.NewFeatureCollection([]*geojson.Feature{poly})
	pt := geojson.NewFeature(
		geojson.NewPoint(geojson.Position{5, 5}),
		map[string]any{"name": "test"},
	)
	points := geojson.NewFeatureCollection([]*geojson.Feature{pt})

	result, err := Tag(points, polygons, "zone", "tagged")
	if err != nil {
		t.Fatal(err)
	}
	if result.Features[0].Properties["name"] != "test" {
		t.Error("original properties should be preserved")
	}
}

func TestCollect(t *testing.T) {
	poly := geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}}}),
		nil,
	)
	polygons := geojson.NewFeatureCollection([]*geojson.Feature{poly})

	pts := make([]*geojson.Feature, 3)
	pts[0] = geojson.NewFeature(
		geojson.NewPoint(geojson.Position{1, 1}),
		map[string]any{"val": 10},
	)
	pts[1] = geojson.NewFeature(
		geojson.NewPoint(geojson.Position{5, 5}),
		map[string]any{"val": 20},
	)
	pts[2] = geojson.NewFeature(
		geojson.NewPoint(geojson.Position{9, 9}),
		map[string]any{"val": 30},
	)
	points := geojson.NewFeatureCollection(pts)

	result, err := Collect(polygons, points, "val", "values")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 1 {
		t.Fatalf("expected 1 feature, got %d", len(result.Features))
	}
	values := result.Features[0].Properties["values"]
	vals, ok := values.([]any)
	if !ok {
		t.Fatalf("expected []any, got %T", values)
	}
	if len(vals) != 3 {
		t.Errorf("expected 3 collected values, got %d", len(vals))
	}
}

func TestCollectOutsidePoint(t *testing.T) {
	poly := geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}}}),
		nil,
	)
	polygons := geojson.NewFeatureCollection([]*geojson.Feature{poly})

	pt := geojson.NewFeature(
		geojson.NewPoint(geojson.Position{20, 20}),
		map[string]any{"val": 100},
	)
	points := geojson.NewFeatureCollection([]*geojson.Feature{pt})

	result, err := Collect(polygons, points, "val", "values")
	if err != nil {
		t.Fatal(err)
	}
	values := result.Features[0].Properties["values"].([]any)
	if len(values) != 0 {
		t.Errorf("expected 0 collected values, got %d", len(values))
	}
}

func TestCollectMultiplePolygons(t *testing.T) {
	poly1 := geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{{{0, 0}, {0, 5}, {5, 5}, {5, 0}, {0, 0}}}),
		nil,
	)
	poly2 := geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{{{5, 5}, {5, 10}, {10, 10}, {10, 5}, {5, 5}}}),
		nil,
	)
	polygons := geojson.NewFeatureCollection([]*geojson.Feature{poly1, poly2})

	pts := make([]*geojson.Feature, 2)
	pts[0] = geojson.NewFeature(
		geojson.NewPoint(geojson.Position{2, 2}),
		map[string]any{"v": "a"},
	)
	pts[1] = geojson.NewFeature(
		geojson.NewPoint(geojson.Position{7, 7}),
		map[string]any{"v": "b"},
	)
	points := geojson.NewFeatureCollection(pts)

	result, err := Collect(polygons, points, "v", "vs")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 2 {
		t.Fatalf("expected 2 features, got %d", len(result.Features))
	}

	v0 := result.Features[0].Properties["vs"].([]any)
	if len(v0) != 1 || v0[0] != "a" {
		t.Errorf("expected ['a'], got %v", v0)
	}
	v1 := result.Features[1].Properties["vs"].([]any)
	if len(v1) != 1 || v1[0] != "b" {
		t.Errorf("expected ['b'], got %v", v1)
	}
}

func TestCollectNilInput(t *testing.T) {
	_, err := Collect(nil, nil, "", "")
	if err == nil {
		t.Error("expected error for nil input")
	}
}

func TestCollectPreservesProperties(t *testing.T) {
	poly := geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}}}),
		map[string]any{"name": "test"},
	)
	polygons := geojson.NewFeatureCollection([]*geojson.Feature{poly})
	pts := make([]*geojson.Feature, 1)
	pts[0] = geojson.NewFeature(
		geojson.NewPoint(geojson.Position{5, 5}),
		map[string]any{"val": 42},
	)
	points := geojson.NewFeatureCollection(pts)

	result, err := Collect(polygons, points, "val", "values")
	if err != nil {
		t.Fatal(err)
	}
	if result.Features[0].Properties["name"] != "test" {
		t.Error("polygon properties should be preserved")
	}
}

func TestCollectEmptyResult(t *testing.T) {
	poly := geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}}}),
		nil,
	)
	polygons := geojson.NewFeatureCollection([]*geojson.Feature{poly})

	pt := geojson.NewFeature(
		geojson.NewPoint(geojson.Position{5, 5}),
		map[string]any{"notval": 42},
	)
	points := geojson.NewFeatureCollection([]*geojson.Feature{pt})

	result, err := Collect(polygons, points, "val", "values")
	if err != nil {
		t.Fatal(err)
	}
	vals := result.Features[0].Properties["values"].([]any)
	if len(vals) != 0 {
		t.Errorf("expected 0 values when field not found, got %d", len(vals))
	}
}

func TestTagMissingField(t *testing.T) {
	poly := geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}}}),
		map[string]any{"zone": "A"},
	)
	polygons := geojson.NewFeatureCollection([]*geojson.Feature{poly})
	pt := geojson.NewFeature(geojson.NewPoint(geojson.Position{5, 5}), nil)
	points := geojson.NewFeatureCollection([]*geojson.Feature{pt})

	result, err := Tag(points, polygons, "nonexistent", "tagged")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := result.Features[0].Properties["tagged"]; ok {
		t.Error("should not set property for missing field")
	}
}
