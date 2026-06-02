package simplify

import (
	"math"
	"math/rand"
	"testing"

	"github.com/ibinh/turf-go/geojson"
)

func pt(lng, lat float64) geojson.Position {
	return geojson.Position{lng, lat}
}

func TestSimplifyLine(t *testing.T) {
	line := geojson.NewFeature(
		geojson.NewLineString([]geojson.Position{
			{0, 0}, {1, 0.1}, {2, 0}, {3, 0.1}, {4, 0}, {5, 0.1}, {6, 0},
		}),
		nil,
	)
	result, err := Simplify(line, 0.5, false)
	if err != nil {
		t.Fatal(err)
	}
	ls := result.Geometry.(*geojson.LineString)
	if len(ls.Coordinates) < 2 {
		t.Error("simplified line should have at least 2 points")
	}
	first := ls.Coordinates[0]
	last := ls.Coordinates[len(ls.Coordinates)-1]
	if first[0] != 0 || first[1] != 0 {
		t.Errorf("first point should remain, got %v", first)
	}
	if last[0] != 6 || last[1] != 0 {
		t.Errorf("last point should remain, got %v", last)
	}
}

func TestSimplifyLineLowTolerance(t *testing.T) {
	line := geojson.NewLineString([]geojson.Position{
		{0, 0}, {1, 1}, {2, 0}, {3, 1}, {4, 0},
	})
	result, err := Simplify(line, 0.01, false)
	if err != nil {
		t.Fatal(err)
	}
	ls := result.Geometry.(*geojson.LineString)
	if len(ls.Coordinates) != 5 {
		t.Errorf("low tolerance should keep all points, got %d", len(ls.Coordinates))
	}
}

func TestSimplifyStraightLine(t *testing.T) {
	line := geojson.NewLineString([]geojson.Position{
		{0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 0},
	})
	result, err := Simplify(line, 0.1, false)
	if err != nil {
		t.Fatal(err)
	}
	ls := result.Geometry.(*geojson.LineString)
	if len(ls.Coordinates) != 2 {
		t.Errorf("straight line should simplify to 2 points, got %d", len(ls.Coordinates))
	}
}

func TestSimplifyPolygon(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {1, 0.1}, {2, 0}, {2, 1}, {2, 2}, {1, 2}, {0, 2}, {0, 1}, {0, 0}},
	})
	result, err := Simplify(poly, 0.3, false)
	if err != nil {
		t.Fatal(err)
	}
	p := result.Geometry.(*geojson.Polygon)
	if len(p.Coordinates[0]) < 4 {
		t.Errorf("simplified polygon should have at least 4 points, got %d", len(p.Coordinates[0]))
	}
	first := p.Coordinates[0][0]
	last := p.Coordinates[0][len(p.Coordinates[0])-1]
	if first[0] != last[0] || first[1] != last[1] {
		t.Error("simplified polygon should be closed")
	}
}

func TestSimplifyInvalidTolerance(t *testing.T) {
	line := geojson.NewLineString([]geojson.Position{{0, 0}, {1, 1}})
	_, err := Simplify(line, 0, false)
	if err == nil {
		t.Error("expected error for zero tolerance")
	}
}

func TestSimplifyPoint(t *testing.T) {
	pt := geojson.NewPoint(geojson.Position{1, 2})
	result, err := Simplify(pt, 0.5, false)
	if err != nil {
		t.Fatal(err)
	}
	_, ok := result.Geometry.(*geojson.Point)
	if !ok {
		t.Error("point should remain unchanged")
	}
}

func TestSimplifyTwoPoints(t *testing.T) {
	line := geojson.NewLineString([]geojson.Position{{0, 0}, {1, 1}})
	result, err := Simplify(line, 0.5, false)
	if err != nil {
		t.Fatal(err)
	}
	ls := result.Geometry.(*geojson.LineString)
	if len(ls.Coordinates) != 2 {
		t.Error("2-point line should remain unchanged")
	}
}

func TestConvexHull(t *testing.T) {
	features := make([]*geojson.Feature, 6)
	features[0] = geojson.NewFeature(geojson.NewPoint(pt(0, 0)), nil)
	features[1] = geojson.NewFeature(geojson.NewPoint(pt(10, 0)), nil)
	features[2] = geojson.NewFeature(geojson.NewPoint(pt(10, 10)), nil)
	features[3] = geojson.NewFeature(geojson.NewPoint(pt(0, 10)), nil)
	features[4] = geojson.NewFeature(geojson.NewPoint(pt(5, 5)), nil)
	features[5] = geojson.NewFeature(geojson.NewPoint(pt(3, 3)), nil)
	fc := geojson.NewFeatureCollection(features)

	hull, err := ConvexHull(fc)
	if err != nil {
		t.Fatal(err)
	}
	poly := hull.Geometry.(*geojson.Polygon)
	if len(poly.Coordinates[0]) < 5 {
		t.Errorf("expected at least 5 points (4+close), got %d", len(poly.Coordinates[0]))
	}
}

func TestConvexHullCollinear(t *testing.T) {
	features := make([]*geojson.Feature, 3)
	features[0] = geojson.NewFeature(geojson.NewPoint(pt(0, 0)), nil)
	features[1] = geojson.NewFeature(geojson.NewPoint(pt(5, 0)), nil)
	features[2] = geojson.NewFeature(geojson.NewPoint(pt(10, 0)), nil)
	fc := geojson.NewFeatureCollection(features)

	hull, err := ConvexHull(fc)
	if err != nil {
		t.Fatal(err)
	}
	ls, ok := hull.Geometry.(*geojson.LineString)
	if !ok {
		t.Fatal("collinear points should return a LineString")
	}
	if len(ls.Coordinates) != 2 {
		t.Errorf("expected 2 points, got %d", len(ls.Coordinates))
	}
}

func TestConvexHullTooFew(t *testing.T) {
	features := make([]*geojson.Feature, 2)
	features[0] = geojson.NewFeature(geojson.NewPoint(pt(0, 0)), nil)
	features[1] = geojson.NewFeature(geojson.NewPoint(pt(1, 1)), nil)
	fc := geojson.NewFeatureCollection(features)

	_, err := ConvexHull(fc)
	if err == nil {
		t.Error("expected error for fewer than 3 points")
	}
}

func TestConvexHullFromPolygon(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}},
	})
	hull, err := ConvexHull(poly)
	if err != nil {
		t.Fatal(err)
	}
	p := hull.Geometry.(*geojson.Polygon)
	if len(p.Coordinates[0]) < 5 {
		t.Errorf("expected at least 5 points, got %d", len(p.Coordinates[0]))
	}
}

func TestConvexHullRectangle(t *testing.T) {
	pts := []geojson.Position{
		{0, 0}, {10, 0}, {10, 10}, {0, 10},
		{2, 2}, {8, 2}, {8, 8}, {2, 8},
	}
	features := make([]*geojson.Feature, len(pts))
	for i, p := range pts {
		features[i] = geojson.NewFeature(geojson.NewPoint(p), nil)
	}
	fc := geojson.NewFeatureCollection(features)

	hull, err := ConvexHull(fc)
	if err != nil {
		t.Fatal(err)
	}
	p := hull.Geometry.(*geojson.Polygon)
	ring := p.Coordinates[0]
	if len(ring) < 5 {
		t.Errorf("expected at least 5 points, got %d", len(ring))
	}
}

func TestSimplifyMultiLineString(t *testing.T) {
	mls := geojson.NewMultiLineString([][]geojson.Position{
		{{0, 0}, {1, 0.1}, {2, 0}},
		{{3, 0}, {4, 0.1}, {5, 0}},
	})
	result, err := Simplify(mls, 0.3, false)
	if err != nil {
		t.Fatal(err)
	}
	_, ok := result.Geometry.(*geojson.MultiLineString)
	if !ok {
		t.Error("expected MultiLineString")
	}
}

func TestSimplifyMultiPolygon(t *testing.T) {
	mp := geojson.NewMultiPolygon([][][]geojson.Position{
		{{{0, 0}, {1, 0.1}, {2, 0}, {2, 1}, {2, 2}, {0, 2}, {0, 0}}},
	})
	result, err := Simplify(mp, 0.3, false)
	if err != nil {
		t.Fatal(err)
	}
	_, ok := result.Geometry.(*geojson.MultiPolygon)
	if !ok {
		t.Error("expected MultiPolygon")
	}
}

func BenchmarkSimplify(b *testing.B) {
	pts := make([]geojson.Position, 1000)
	for i := 0; i < 1000; i++ {
		pts[i] = pt(float64(i), math.Sin(float64(i)*0.1)*10)
	}
	line := geojson.NewLineString(pts)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Simplify(line, 0.5, false)
	}
}

func BenchmarkConvexHull(b *testing.B) {
	pts := make([]geojson.Position, 1000)
	for i := 0; i < 1000; i++ {
		pts[i] = pt(rand.Float64()*100, rand.Float64()*100)
	}
	features := make([]*geojson.Feature, 1000)
	for i, p := range pts {
		features[i] = geojson.NewFeature(geojson.NewPoint(p), nil)
	}
	fc := geojson.NewFeatureCollection(features)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ConvexHull(fc)
	}
}

func TestSimplifyPreservesLine(t *testing.T) {
	line := geojson.NewLineString([]geojson.Position{
		{0, 0}, {1, 1}, {2, 2}, {3, 3},
	})
	result, err := Simplify(line, 0.5, false)
	if err != nil {
		t.Fatal(err)
	}
	ls := result.Geometry.(*geojson.LineString)
	if len(ls.Coordinates) != 2 {
		t.Errorf("collinear points should simplify to 2, got %d", len(ls.Coordinates))
	}
}
