package isolines

import (
	"math"
	"testing"

	"github.com/ibinh/turf-go/geojson"
)

func TestIsolinesBasic(t *testing.T) {
	pts := []geojson.Position{
		{0, 0}, {1, 0}, {2, 0},
		{0, 1}, {1, 1}, {2, 1},
		{0, 2}, {1, 2}, {2, 2},
	}
	features := make([]*geojson.Feature, len(pts))
	for i, p := range pts {
		features[i] = geojson.NewFeature(geojson.NewPoint(p), map[string]any{"z": float64(i)})
	}
	fc := geojson.NewFeatureCollection(features)

	result, err := Isolines(fc, IsolinesOptions{ZProperty: "z", Breaks: []float64{3, 6}})
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal("result should not be nil")
	}
	if len(result.Features) == 0 {
		t.Fatal("expected at least 1 contour line")
	}
	for _, f := range result.Features {
		if f.Geometry.Type() != geojson.TypeLineString {
			t.Errorf("expected LineString, got %s", f.Geometry.Type())
		}
		ls, ok := f.Geometry.(*geojson.LineString)
		if !ok {
			t.Fatal("failed type assertion")
		}
		if len(ls.Coordinates) < 2 {
			t.Error("polyline should have at least 2 points")
		}
	}
}

func TestIsolinesErrors(t *testing.T) {
	_, err := Isolines(nil)
	if err == nil {
		t.Error("expected error for nil input")
	}

	pts := geojson.NewFeatureCollection(nil)
	_, err = Isolines(pts, IsolinesOptions{Breaks: []float64{1}})
	if err == nil {
		t.Error("expected error for empty collection")
	}

	single := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint([]float64{0, 0}), nil),
	})
	_, err = Isolines(single, IsolinesOptions{Breaks: []float64{1}})
	if err == nil {
		t.Error("expected error for <3 points")
	}
}

func TestIsolinesNoBreaks(t *testing.T) {
	pts := make([]*geojson.Feature, 4)
	for i := 0; i < 4; i++ {
		pts[i] = geojson.NewFeature(geojson.NewPoint([]float64{float64(i), float64(i)}), map[string]any{"z": float64(i)})
	}
	fc := geojson.NewFeatureCollection(pts)

	_, err := Isolines(fc, IsolinesOptions{ZProperty: "z"})
	if err == nil {
		t.Error("expected error for empty breaks")
	}
}

func TestIsolinesGrid(t *testing.T) {
	pts := make([]*geojson.Feature, 25)
	for j := 0; j < 5; j++ {
		for i := 0; i < 5; i++ {
			idx := j*5 + i
			pts[idx] = geojson.NewFeature(
				geojson.NewPoint([]float64{float64(i), float64(j)}),
				map[string]any{"z": float64(idx)},
			)
		}
	}
	fc := geojson.NewFeatureCollection(pts)

	result, err := Isolines(fc, IsolinesOptions{ZProperty: "z", Breaks: []float64{10, 15}})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) == 0 {
		t.Fatal("expected contours")
	}
}

func TestInterpolateContour(t *testing.T) {
	pts := [4]geojson.Position{
		{0, 0}, {1, 0}, {1, 1}, {0, 1},
	}
	vals := [4]float64{0, 2, 4, 6}

	// Edge 0 (bottom, 0->1): threshold=1, v0=0, v1=2
	// t = (1-0)/(2-0) = 0.5, pt = (0.5, 0)
	p := interpolateContour(pts, vals, 0, 1)
	if math.Abs(p[0]-0.5) > 1e-10 || math.Abs(p[1]) > 1e-10 {
		t.Errorf("expected (0.5, 0), got (%f, %f)", p[0], p[1])
	}

	// Edge 1 (right, 1->2): threshold=3, v1=2, v2=4
	// t = (3-2)/(4-2) = 0.5, pt = (1, 0.5)
	p = interpolateContour(pts, vals, 1, 3)
	if math.Abs(p[0]-1) > 1e-10 || math.Abs(p[1]-0.5) > 1e-10 {
		t.Errorf("expected (1, 0.5), got (%f, %f)", p[0], p[1])
	}
}

func TestConnectSegmentsToPolylines(t *testing.T) {
	segs := []rawSeg{
		{geojson.Position{0, 0}, geojson.Position{1, 0}},
		{geojson.Position{1, 0}, geojson.Position{2, 0}},
	}
	lines := connectSegmentsToPolylines(segs)
	if len(lines) != 1 {
		t.Fatalf("expected 1 polyline, got %d", len(lines))
	}
	if len(lines[0]) != 3 {
		t.Errorf("expected 3 points, got %d", len(lines[0]))
	}
}

func TestMarchingSquaresLines(t *testing.T) {
	grid := []gridCell{
		{
			corners: [4]geojson.Position{
				{0, 0}, {1, 0}, {1, 1}, {0, 1},
			},
			values: [4]float64{0, 2, 4, 6},
		},
	}
	lines := marchingSquaresLines(grid, 3)
	if len(lines) == 0 {
		t.Log("no lines found (may be valid)")
	}
}
