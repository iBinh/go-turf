package shapes

import (
	"math"
	"testing"

	"github.com/ibinh/turf-go/geojson"
	"github.com/ibinh/turf-go/measurement"
)

func TestCircle(t *testing.T) {
	center := geojson.NewPoint(geojson.Position{0, 0})
	circle, err := Circle(center, 100)
	if err != nil {
		t.Fatal(err)
	}
	if circle == nil {
		t.Fatal("expected circle feature")
	}
	poly, ok := circle.Geometry.(*geojson.Polygon)
	if !ok {
		t.Fatalf("expected Polygon, got %T", circle.Geometry)
	}
	ring := poly.Coordinates[0]
	if len(ring) < 4 {
		t.Errorf("expected at least 4 vertices, got %d", len(ring))
	}
	first, last := ring[0], ring[len(ring)-1]
	if math.Abs(first[0]-last[0]) > 0.001 || math.Abs(first[1]-last[1]) > 0.001 {
		t.Error("circle ring is not closed")
	}
}

func TestCircleSteps(t *testing.T) {
	center := geojson.NewPoint(geojson.Position{0, 0})
	circle, err := Circle(center, 100, CircleOptions{Steps: 8})
	if err != nil {
		t.Fatal(err)
	}
	poly := circle.Geometry.(*geojson.Polygon)
	if len(poly.Coordinates[0]) != 9 {
		t.Errorf("expected 9 vertices (8 + close), got %d", len(poly.Coordinates[0]))
	}
}

func TestCircleProperties(t *testing.T) {
	center := geojson.NewPoint(geojson.Position{0, 0})
	circle, err := Circle(center, 100, CircleOptions{Properties: map[string]any{"name": "test"}})
	if err != nil {
		t.Fatal(err)
	}
	if circle.Properties["name"] != "test" {
		t.Errorf("expected property name=test, got %v", circle.Properties["name"])
	}
}

func TestEllipse(t *testing.T) {
	center := geojson.NewPoint(geojson.Position{0, 0})
	ellipse, err := Ellipse(center, 200, 100)
	if err != nil {
		t.Fatal(err)
	}
	poly, ok := ellipse.Geometry.(*geojson.Polygon)
	if !ok {
		t.Fatalf("expected Polygon, got %T", ellipse.Geometry)
	}
	ring := poly.Coordinates[0]
	if len(ring) < 4 {
		t.Errorf("expected at least 4 vertices, got %d", len(ring))
	}
}

func TestEllipseSteps(t *testing.T) {
	center := geojson.NewPoint(geojson.Position{0, 0})
	ellipse, err := Ellipse(center, 100, 50, EllipseOptions{Steps: 16})
	if err != nil {
		t.Fatal(err)
	}
	poly := ellipse.Geometry.(*geojson.Polygon)
	if len(poly.Coordinates[0]) != 17 {
		t.Errorf("expected 17 vertices (16 + close), got %d", len(poly.Coordinates[0]))
	}
}

func TestEllipseMinSteps(t *testing.T) {
	center := geojson.NewPoint(geojson.Position{0, 0})
	ellipse, err := Ellipse(center, 100, 50, EllipseOptions{Steps: 1})
	if err != nil {
		t.Fatal(err)
	}
	poly := ellipse.Geometry.(*geojson.Polygon)
	if len(poly.Coordinates[0]) < 4 {
		t.Errorf("expected at least 4 vertices, got %d", len(poly.Coordinates[0]))
	}
}

func TestEllipseAngled(t *testing.T) {
	center := geojson.NewPoint(geojson.Position{0, 0})
	ellipse, err := Ellipse(center, 100, 50, EllipseOptions{Angle: 45})
	if err != nil {
		t.Fatal(err)
	}
	if ellipse == nil {
		t.Fatal("expected ellipse feature")
	}
}

func TestBezierSpline(t *testing.T) {
	line := geojson.NewLineString([]geojson.Position{
		{0, 0}, {10, 10}, {20, 0},
	})
	bezier, err := BezierSpline(line)
	if err != nil {
		t.Fatal(err)
	}
	ls, ok := bezier.Geometry.(*geojson.LineString)
	if !ok {
		t.Fatalf("expected LineString, got %T", bezier.Geometry)
	}
	if len(ls.Coordinates) < 2 {
		t.Error("bezier should have at least 2 points")
	}
	first := ls.Coordinates[0]
	last := ls.Coordinates[len(ls.Coordinates)-1]
	if math.Abs(first[0]) > 0.001 || math.Abs(first[1]) > 0.001 {
		t.Errorf("expected first point at origin, got %v", first)
	}
	if math.Abs(last[0]-20) > 0.001 || math.Abs(last[1]) > 0.001 {
		t.Errorf("expected last point at [20,0], got %v", last)
	}
}

func TestBezierTwoPoints(t *testing.T) {
	line := geojson.NewLineString([]geojson.Position{
		{0, 0}, {10, 10},
	})
	bezier, err := BezierSpline(line)
	if err != nil {
		t.Fatal(err)
	}
	if bezier == nil {
		t.Fatal("expected bezier feature")
	}
}

func TestBezierOptions(t *testing.T) {
	line := geojson.NewLineString([]geojson.Position{
		{0, 0}, {5, 10}, {10, 0},
	})
	bezier, err := BezierSpline(line, BezierOptions{Sharpness: 0.5, Resolution: 100})
	if err != nil {
		t.Fatal(err)
	}
	ls := bezier.Geometry.(*geojson.LineString)
	if len(ls.Coordinates) != 101 {
		t.Errorf("expected 101 points, got %d", len(ls.Coordinates))
	}
}

func TestBezierInvalidInput(t *testing.T) {
	pt := geojson.NewPoint(geojson.Position{0, 0})
	_, err := BezierSpline(pt)
	if err == nil {
		t.Error("expected error for Point input")
	}
}

func TestBezierSinglePoint(t *testing.T) {
	line := geojson.NewLineString([]geojson.Position{{0, 0}})
	_, err := BezierSpline(line)
	if err == nil {
		t.Error("expected error for single point")
	}
}

func TestRandomPosition(t *testing.T) {
	pos := RandomPosition(nil)
	if len(pos) < 2 {
		t.Error("expected position with at least 2 coords")
	}
}

func TestRandomPositionBBox(t *testing.T) {
	bbox := []float64{0, 0, 10, 10}
	for i := 0; i < 100; i++ {
		pos := RandomPosition(bbox)
		if pos[0] < 0 || pos[0] > 10 || pos[1] < 0 || pos[1] > 10 {
			t.Errorf("position %v outside bbox %v", pos, bbox)
		}
	}
}

func TestRandomPoint(t *testing.T) {
	fc, err := RandomPoint(5)
	if err != nil {
		t.Fatal(err)
	}
	if len(fc.Features) != 5 {
		t.Errorf("expected 5 features, got %d", len(fc.Features))
	}
	for _, f := range fc.Features {
		if _, ok := f.Geometry.(*geojson.Point); !ok {
			t.Errorf("expected Point, got %T", f.Geometry)
		}
	}
}

func TestRandomPointInBBox(t *testing.T) {
	bbox := []float64{0, 0, 10, 10}
	fc, err := RandomPoint(10, RandomOptions{BBox: bbox})
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range fc.Features {
		coord, _ := geojson.GetCoord(f)
		if coord[0] < 0 || coord[0] > 10 || coord[1] < 0 || coord[1] > 10 {
			t.Errorf("point %v outside bbox %v", coord, bbox)
		}
	}
}

func TestRandomPointProperties(t *testing.T) {
	fc, err := RandomPoint(1, RandomOptions{Properties: map[string]any{"type": "random"}})
	if err != nil {
		t.Fatal(err)
	}
	if fc.Features[0].Properties["type"] != "random" {
		t.Error("expected property type=random")
	}
}

func TestRandomLineString(t *testing.T) {
	fc, err := RandomLineString(3)
	if err != nil {
		t.Fatal(err)
	}
	if len(fc.Features) != 3 {
		t.Errorf("expected 3 features, got %d", len(fc.Features))
	}
	for _, f := range fc.Features {
		if _, ok := f.Geometry.(*geojson.LineString); !ok {
			t.Errorf("expected LineString, got %T", f.Geometry)
		}
	}
}

func TestRandomLineStringVertices(t *testing.T) {
	fc, err := RandomLineString(1, RandomOptions{NumVertices: 5})
	if err != nil {
		t.Fatal(err)
	}
	ls := fc.Features[0].Geometry.(*geojson.LineString)
	if len(ls.Coordinates) != 5 {
		t.Errorf("expected 5 vertices, got %d", len(ls.Coordinates))
	}
}

func TestRandomLineStringMinVertices(t *testing.T) {
	fc, err := RandomLineString(1, RandomOptions{NumVertices: 0})
	if err != nil {
		t.Fatal(err)
	}
	ls := fc.Features[0].Geometry.(*geojson.LineString)
	if len(ls.Coordinates) < 2 {
		t.Errorf("expected at least 2 vertices, got %d", len(ls.Coordinates))
	}
}

func TestRandomPolygon(t *testing.T) {
	fc, err := RandomPolygon(3)
	if err != nil {
		t.Fatal(err)
	}
	if len(fc.Features) != 3 {
		t.Errorf("expected 3 features, got %d", len(fc.Features))
	}
	for _, f := range fc.Features {
		if _, ok := f.Geometry.(*geojson.Polygon); !ok {
			t.Errorf("expected Polygon, got %T", f.Geometry)
		}
	}
}

func TestRandomPolygonVertices(t *testing.T) {
	fc, err := RandomPolygon(1, RandomOptions{NumVertices: 6})
	if err != nil {
		t.Fatal(err)
	}
	poly := fc.Features[0].Geometry.(*geojson.Polygon)
	expected := 7
	if len(poly.Coordinates[0]) != expected {
		t.Errorf("expected %d vertices (6 + close), got %d", expected, len(poly.Coordinates[0]))
	}
}

func TestRandomPolygonInBBox(t *testing.T) {
	bbox := []float64{0, 0, 10, 10}
	fc, err := RandomPolygon(1, RandomOptions{BBox: bbox, MaxLength: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(fc.Features) != 1 {
		t.Fatal("expected 1 feature")
	}
}

func TestRandomPolygonClosed(t *testing.T) {
	fc, err := RandomPolygon(1)
	if err != nil {
		t.Fatal(err)
	}
	poly := fc.Features[0].Geometry.(*geojson.Polygon)
	ring := poly.Coordinates[0]
	first, last := ring[0], ring[len(ring)-1]
	if first[0] != last[0] || first[1] != last[1] {
		t.Error("polygon ring is not closed")
	}
}

func TestCircleDegreeUnits(t *testing.T) {
	center := geojson.NewPoint(geojson.Position{0, 0})
	circle, err := Circle(center, 10, CircleOptions{Units: measurement.UnitDegrees})
	if err != nil {
		t.Fatal(err)
	}
	if circle == nil {
		t.Fatal("expected circle feature")
	}
}

func TestEllipseProperties(t *testing.T) {
	center := geojson.NewPoint(geojson.Position{0, 0})
	ellipse, err := Ellipse(center, 100, 50, EllipseOptions{Properties: map[string]any{"name": "test"}})
	if err != nil {
		t.Fatal(err)
	}
	if ellipse.Properties["name"] != "test" {
		t.Errorf("expected property name=test, got %v", ellipse.Properties["name"])
	}
}

func BenchmarkCircle(b *testing.B) {
	center := geojson.NewPoint(geojson.Position{0, 0})
	for i := 0; i < b.N; i++ {
		Circle(center, 100)
	}
}

func BenchmarkBezierSpline(b *testing.B) {
	line := geojson.NewLineString([]geojson.Position{
		{0, 0}, {5, 10}, {10, 5}, {15, 10}, {20, 0},
	})
	for i := 0; i < b.N; i++ {
		BezierSpline(line)
	}
}
