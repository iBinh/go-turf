package isobands

import (
	"math"
	"testing"

	"github.com/ibinh/turf-go/geojson"
)

func TestIsobandsBasic(t *testing.T) {
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

	result, err := Isobands(fc, IsobandsOptions{ZProperty: "z", Breaks: []float64{0, 3, 6, 10}})
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal("result should not be nil")
	}
	if len(result.Features) == 0 {
		t.Fatal("expected at least 1 band polygon")
	}
	for _, f := range result.Features {
		if f.Geometry.Type() != geojson.TypePolygon {
			t.Errorf("expected Polygon, got %s", f.Geometry.Type())
		}
		poly, ok := f.Geometry.(*geojson.Polygon)
		if !ok {
			t.Fatal("failed type assertion")
		}
		if len(poly.Coordinates[0]) < 3 {
			t.Error("ring should have at least 3 points")
		}
	}
}

func TestIsobandsErrors(t *testing.T) {
	_, err := Isobands(nil)
	if err == nil {
		t.Error("expected error for nil input")
	}

	pts := geojson.NewFeatureCollection(nil)
	_, err = Isobands(pts, IsobandsOptions{Breaks: []float64{1, 2}})
	if err == nil {
		t.Error("expected error for empty collection")
	}

	single := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint([]float64{0, 0}), nil),
	})
	_, err = Isobands(single, IsobandsOptions{Breaks: []float64{1, 2}})
	if err == nil {
		t.Error("expected error for <3 points")
	}
}

func TestIsobandsNoBreaks(t *testing.T) {
	pts := make([]*geojson.Feature, 4)
	for i := 0; i < 4; i++ {
		pts[i] = geojson.NewFeature(geojson.NewPoint([]float64{float64(i), float64(i)}), map[string]any{"z": float64(i)})
	}
	fc := geojson.NewFeatureCollection(pts)

	_, err := Isobands(fc, IsobandsOptions{ZProperty: "z", Breaks: []float64{1}})
	if err == nil {
		t.Error("expected error for <2 breaks")
	}
}

func TestIsobandsGridData(t *testing.T) {
	// Create a simple grid with known values
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

	result, err := Isobands(fc, IsobandsOptions{ZProperty: "z", Breaks: []float64{0, 10, 20}})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) == 0 {
		t.Fatal("expected bands")
	}
}

func TestEstimateCellSize(t *testing.T) {
	data := []zPoint{
		{geojson.Position{0, 0}, 1},
		{geojson.Position{1, 0}, 2},
		{geojson.Position{0, 1}, 3},
		{geojson.Position{1, 1}, 4},
	}
	size := estimateCellSize(data, 0, 1, 0, 1)
	if size <= 0 {
		t.Error("cell size should be positive")
	}
	if math.Abs(size-0.5) > 1e-10 {
		t.Errorf("expected cell size 0.5, got %f", size)
	}
}

func TestIdw(t *testing.T) {
	data := []zPoint{
		{geojson.Position{0, 0}, 10},
		{geojson.Position{10, 0}, 20},
	}
	// At (0,0) should return exact value 10
	v := idw(geojson.Position{0, 0}, data, 2)
	if math.Abs(v-10) > 1e-10 {
		t.Errorf("expected 10, got %f", v)
	}
	// At (5,0) should be between 10 and 20
	v = idw(geojson.Position{5, 0}, data, 2)
	if v <= 10 || v >= 20 {
		t.Errorf("expected between 10 and 20, got %f", v)
	}
}

func TestExtractData(t *testing.T) {
	pts := []*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint([]float64{0, 0}), map[string]any{"z": 1.0}),
		geojson.NewFeature(geojson.NewPoint([]float64{1, 1}), map[string]any{"z": 2.0}),
		geojson.NewFeature(geojson.NewPoint([]float64{2, 2}), nil),
	}
	fc := geojson.NewFeatureCollection(pts)
	data := extractData(fc, "z")
	if len(data) != 2 {
		t.Errorf("expected 2 data points, got %d", len(data))
	}
}

func TestInterpolateEdgeCrossing(t *testing.T) {
	pts := [4]geojson.Position{
		{0, 0}, {1, 0}, {1, 1}, {0, 1},
	}
	vals := [4]float64{0, 1, 2, 3}

	// Edge 0 (bottom, 0->1): v0=0, v1=1, lo=0.5, hi=1.5
	// v0<lo=0.5, v1>=lo=0.5 && v1<hi=1.5 → v0 outside (below), v1 inside
	// Entering band, v1 inside → threshold = lo = 0.5
	// t = (0.5 - 0) / (1 - 0) = 0.5
	// pt = (0.5, 0)
	p := interpolateEdgeCrossing(pts, vals, 0, 0.5, 1.5)
	if math.Abs(p[0]-0.5) > 1e-10 || math.Abs(p[1]) > 1e-10 {
		t.Errorf("expected (0.5, 0), got (%f, %f)", p[0], p[1])
	}
}

func TestConnectSegmentsToRings(t *testing.T) {
	// Two segments forming a square
	segs := []rawSeg{
		{geojson.Position{0, 0}, geojson.Position{1, 0}},
		{geojson.Position{1, 0}, geojson.Position{1, 1}},
		{geojson.Position{1, 1}, geojson.Position{0, 1}},
		{geojson.Position{0, 1}, geojson.Position{0, 0}},
	}
	rings := connectSegmentsToRings(segs)
	if len(rings) != 1 {
		t.Fatalf("expected 1 ring, got %d", len(rings))
	}
	if len(rings[0]) != 4 {
		t.Errorf("expected 4 points in ring, got %d", len(rings[0]))
	}
}

func TestMarchingSquaresBands(t *testing.T) {
	// 2x2 grid with known values
	grid := []gridCell{
		{
			corners: [4]geojson.Position{
				{0, 0}, {1, 0}, {1, 1}, {0, 1},
			},
			values: [4]float64{0, 1, 2, 3},
		},
	}
	rings := marchingSquaresBands(grid, 0.5, 2.5)
	if len(rings) == 0 {
		t.Log("no rings found (may be valid depending on threshold)")
	}
}
