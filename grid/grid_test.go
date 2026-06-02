package grid

import (
	"testing"
	"github.com/ibinh/turf-go/geojson"
	"github.com/ibinh/turf-go/measurement"
)

func TestHexGrid(t *testing.T) {
	bbox := []float64{-10, -10, 10, 10}
	fc, err := HexGrid(bbox, 5, measurement.UnitDegrees)
	if err != nil {
		t.Fatal(err)
	}
	if len(fc.Features) == 0 {
		t.Error("expected at least one hex cell")
	}
	for _, f := range fc.Features {
		poly, ok := f.Geometry.(*geojson.Polygon)
		if !ok {
			t.Errorf("expected Polygon, got %T", f.Geometry)
			continue
		}
		if len(poly.Coordinates[0]) != 7 {
			t.Errorf("expected 7 vertices (closed hex), got %d", len(poly.Coordinates[0]))
		}
	}
}

func TestHexGridSmallBBox(t *testing.T) {
	bbox := []float64{0, 0, 1, 1}
	fc, err := HexGrid(bbox, 0.5, measurement.UnitDegrees)
	if err != nil {
		t.Fatal(err)
	}
	if len(fc.Features) == 0 {
		t.Error("expected hex cells for small bbox")
	}
}

func TestPointGrid(t *testing.T) {
	bbox := []float64{0, 0, 10, 10}
	fc, err := PointGrid(bbox, 2, measurement.UnitDegrees)
	if err != nil {
		t.Fatal(err)
	}
	if len(fc.Features) == 0 {
		t.Error("expected at least one point")
	}
	for _, f := range fc.Features {
		if _, ok := f.Geometry.(*geojson.Point); !ok {
			t.Errorf("expected Point, got %T", f.Geometry)
		}
	}
}

func TestPointGridCount(t *testing.T) {
	bbox := []float64{0, 0, 10, 10}
	fc, err := PointGrid(bbox, 5, measurement.UnitDegrees)
	if err != nil {
		t.Fatal(err)
	}
	if len(fc.Features) == 0 {
		t.Error("expected points")
	}
}

func TestSquareGrid(t *testing.T) {
	bbox := []float64{0, 0, 10, 10}
	fc, err := SquareGrid(bbox, 2, measurement.UnitDegrees)
	if err != nil {
		t.Fatal(err)
	}
	if len(fc.Features) == 0 {
		t.Error("expected at least one square")
	}
	for _, f := range fc.Features {
		poly, ok := f.Geometry.(*geojson.Polygon)
		if !ok {
			t.Errorf("expected Polygon, got %T", f.Geometry)
			continue
		}
		if len(poly.Coordinates[0]) != 5 {
			t.Errorf("expected 5 vertices (closed square), got %d", len(poly.Coordinates[0]))
		}
	}
}

func TestSquareGridProperties(t *testing.T) {
	bbox := []float64{0, 0, 5, 5}
	props := map[string]any{"name": "test"}
	fc, err := SquareGrid(bbox, 1, measurement.UnitDegrees, GridOptions{Properties: props})
	if err != nil {
		t.Fatal(err)
	}
	if len(fc.Features) == 0 {
		t.Fatal("expected features")
	}
	for _, f := range fc.Features {
		if f.Properties["name"] != "test" {
			t.Errorf("expected property name=test, got %v", f.Properties["name"])
		}
	}
}

func TestTriangleGrid(t *testing.T) {
	bbox := []float64{0, 0, 10, 10}
	fc, err := TriangleGrid(bbox, 2, measurement.UnitDegrees)
	if err != nil {
		t.Fatal(err)
	}
	if len(fc.Features) == 0 {
		t.Error("expected at least one triangle")
	}
	for _, f := range fc.Features {
		poly, ok := f.Geometry.(*geojson.Polygon)
		if !ok {
			t.Errorf("expected Polygon, got %T", f.Geometry)
			continue
		}
		if len(poly.Coordinates[0]) != 4 {
			t.Errorf("expected 4 vertices (closed triangle), got %d", len(poly.Coordinates[0]))
		}
	}
}

func TestTriangleGridEvenCount(t *testing.T) {
	bbox := []float64{0, 0, 10, 10}
	fc, err := TriangleGrid(bbox, 5, measurement.UnitDegrees)
	if err != nil {
		t.Fatal(err)
	}
	if len(fc.Features)%2 != 0 {
		t.Errorf("expected even number of triangles, got %d", len(fc.Features))
	}
}

func TestGridInvalidBBox(t *testing.T) {
	_, err := HexGrid([]float64{0, 0}, 1, measurement.UnitDegrees)
	if err == nil {
		t.Error("expected error for invalid bbox")
	}
	_, err = PointGrid([]float64{0, 0}, 1, measurement.UnitDegrees)
	if err == nil {
		t.Error("expected error for invalid bbox")
	}
	_, err = SquareGrid([]float64{0, 0}, 1, measurement.UnitDegrees)
	if err == nil {
		t.Error("expected error for invalid bbox")
	}
	_, err = TriangleGrid([]float64{0, 0}, 1, measurement.UnitDegrees)
	if err == nil {
		t.Error("expected error for invalid bbox")
	}
}

func TestHexGridDegenerate(t *testing.T) {
	bbox := []float64{180, -90, -180, 90}
	_, err := HexGrid(bbox, 100, measurement.UnitKilometers)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPointGridKilometers(t *testing.T) {
	bbox := []float64{0, 0, 1, 1}
	fc, err := PointGrid(bbox, 50, measurement.UnitKilometers)
	if err != nil {
		t.Fatal(err)
	}
	if len(fc.Features) == 0 {
		t.Error("expected points with kilometer units")
	}
}

func BenchmarkHexGrid(b *testing.B) {
	bbox := []float64{-10, -10, 10, 10}
	for i := 0; i < b.N; i++ {
		HexGrid(bbox, 5, measurement.UnitDegrees)
	}
}

func BenchmarkPointGrid(b *testing.B) {
	bbox := []float64{0, 0, 10, 10}
	for i := 0; i < b.N; i++ {
		PointGrid(bbox, 1, measurement.UnitDegrees)
	}
}
