package polyclip

import (
	"math"
	"testing"

	"github.com/ibinh/turf-go/geojson"
)

func TestUnionOverlappingSquares(t *testing.T) {
	a := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {2, 0}, {2, 2}, {0, 2}, {0, 0}},
	})
	b := geojson.NewPolygon([][]geojson.Position{
		{{1, 1}, {3, 1}, {3, 3}, {1, 3}, {1, 1}},
	})

	result, err := PolygonUnion(a, b)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal("union should not be nil")
	}

	poly, ok := result.Geometry.(*geojson.Polygon)
	if !ok {
		t.Fatalf("expected Polygon, got %T", result.Geometry)
	}
	area := polygonArea(poly.Coordinates[0])
	expected := 7.0
	if math.Abs(area-expected) > 0.01 {
		t.Errorf("expected area ~%.2f, got %.2f", expected, area)
	}
}

func TestIntersectOverlappingSquares(t *testing.T) {
	a := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {2, 0}, {2, 2}, {0, 2}, {0, 0}},
	})
	b := geojson.NewPolygon([][]geojson.Position{
		{{1, 1}, {3, 1}, {3, 3}, {1, 3}, {1, 1}},
	})

	result, err := PolygonIntersect(a, b)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal("intersect should not be nil")
	}

	poly, ok := result.Geometry.(*geojson.Polygon)
	if !ok {
		t.Fatalf("expected Polygon, got %T", result.Geometry)
	}
	area := polygonArea(poly.Coordinates[0])
	expected := 1.0
	if math.Abs(area-expected) > 0.01 {
		t.Errorf("expected area ~%.2f, got %.2f", expected, area)
	}
}

func TestDifferenceOverlappingSquares(t *testing.T) {
	a := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {2, 0}, {2, 2}, {0, 2}, {0, 0}},
	})
	b := geojson.NewPolygon([][]geojson.Position{
		{{1, 1}, {3, 1}, {3, 3}, {1, 3}, {1, 1}},
	})

	result, err := PolygonDifference(a, b)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal("difference should not be nil")
	}

	poly, ok := result.Geometry.(*geojson.Polygon)
	if !ok {
		t.Fatalf("expected Polygon, got %T", result.Geometry)
	}
	area := polygonArea(poly.Coordinates[0])
	expected := 3.0
	if math.Abs(area-expected) > 0.01 {
		t.Errorf("expected area ~%.2f, got %.2f", expected, area)
	}
}

func TestUnionDisjointSquares(t *testing.T) {
	a := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}},
	})
	b := geojson.NewPolygon([][]geojson.Position{
		{{2, 0}, {3, 0}, {3, 1}, {2, 1}, {2, 0}},
	})

	result, err := PolygonUnion(a, b)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal("union should not be nil")
	}

	mp, ok := result.Geometry.(*geojson.MultiPolygon)
	if !ok {
		poly, ok := result.Geometry.(*geojson.Polygon)
		if ok {
			a := polygonArea(poly.Coordinates[0])
			if a > 1.1 {
				t.Errorf("single polygon area should be ~1, got %.2f", a)
			}
			return
		}
		t.Fatalf("expected MultiPolygon or Polygon, got %T", result.Geometry)
	}
	if len(mp.Coordinates) < 2 {
		t.Error("disjoint squares should produce at least 2 polygons")
	}
	area := 0.0
	for _, poly := range mp.Coordinates {
		area += polygonArea(poly[0])
	}
	if math.Abs(area-2.0) > 0.01 {
		t.Errorf("expected total area 2.0, got %.2f", area)
	}
}

func TestIntersectDisjoint(t *testing.T) {
	a := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}},
	})
	b := geojson.NewPolygon([][]geojson.Position{
		{{2, 0}, {3, 0}, {3, 1}, {2, 1}, {2, 0}},
	})

	result, err := PolygonIntersect(a, b)
	if err != nil {
		t.Fatal(err)
	}
	if result != nil {
		t.Error("disjoint polygons should have nil intersection")
	}
}

func TestDifferenceNonOverlapping(t *testing.T) {
	a := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}},
	})
	b := geojson.NewPolygon([][]geojson.Position{
		{{2, 0}, {3, 0}, {3, 1}, {2, 1}, {2, 0}},
	})

	result, err := PolygonDifference(a, b)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal("difference of non-overlapping should return a")
	}
	poly, ok := result.Geometry.(*geojson.Polygon)
	if !ok {
		t.Fatalf("expected Polygon, got %T", result.Geometry)
	}
	area := polygonArea(poly.Coordinates[0])
	if math.Abs(area-1.0) > 0.01 {
		t.Errorf("expected area 1.0, got %.2f", area)
	}
}

func TestUnionIdentical(t *testing.T) {
	a := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {2, 0}, {2, 2}, {0, 2}, {0, 0}},
	})
	b := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {2, 0}, {2, 2}, {0, 2}, {0, 0}},
	})

	result, err := PolygonUnion(a, b)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal("union should not be nil")
	}
	poly, ok := result.Geometry.(*geojson.Polygon)
	if !ok {
		t.Fatalf("expected Polygon, got %T", result.Geometry)
	}
	area := polygonArea(poly.Coordinates[0])
	if math.Abs(area-4.0) > 0.01 {
		t.Errorf("expected area 4.0, got %.2f", area)
	}
}

func TestIntersectIdentical(t *testing.T) {
	a := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {2, 0}, {2, 2}, {0, 2}, {0, 0}},
	})

	result, err := PolygonIntersect(a, a)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal("intersect of identical should not be nil")
	}
	poly, ok := result.Geometry.(*geojson.Polygon)
	if !ok {
		t.Fatalf("expected Polygon, got %T", result.Geometry)
	}
	area := polygonArea(poly.Coordinates[0])
	if math.Abs(area-4.0) > 0.01 {
		t.Errorf("expected area 4.0, got %.2f", area)
	}
}

func TestUnionContained(t *testing.T) {
	outer := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {5, 0}, {5, 5}, {0, 5}, {0, 0}},
	})
	inner := geojson.NewPolygon([][]geojson.Position{
		{{1, 1}, {2, 1}, {2, 2}, {1, 2}, {1, 1}},
	})

	result, err := PolygonUnion(outer, inner)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal("union should not be nil")
	}
	poly, ok := result.Geometry.(*geojson.Polygon)
	if !ok {
		t.Fatalf("expected Polygon, got %T", result.Geometry)
	}
	area := polygonArea(poly.Coordinates[0])
	if math.Abs(area-25.0) > 0.01 {
		t.Errorf("expected area 25.0, got %.2f", area)
	}
}

func TestIntersectContained(t *testing.T) {
	outer := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {5, 0}, {5, 5}, {0, 5}, {0, 0}},
	})
	inner := geojson.NewPolygon([][]geojson.Position{
		{{1, 1}, {2, 1}, {2, 2}, {1, 2}, {1, 1}},
	})

	result, err := PolygonIntersect(outer, inner)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal("intersect should not be nil")
	}
	poly, ok := result.Geometry.(*geojson.Polygon)
	if !ok {
		t.Fatalf("expected Polygon, got %T", result.Geometry)
	}
	area := polygonArea(poly.Coordinates[0])
	if math.Abs(area-1.0) > 0.01 {
		t.Errorf("expected area 1.0, got %.2f", area)
	}
}

func TestDifferenceContained(t *testing.T) {
	outer := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {5, 0}, {5, 5}, {0, 5}, {0, 0}},
	})
	inner := geojson.NewPolygon([][]geojson.Position{
		{{1, 1}, {2, 1}, {2, 2}, {1, 2}, {1, 1}},
	})

	result, err := PolygonDifference(outer, inner)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal("difference should not be nil")
	}
}

func TestUnionTriangle(t *testing.T) {
	a := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {2, 0}, {1, 2}, {0, 0}},
	})
	b := geojson.NewPolygon([][]geojson.Position{
		{{2, 0}, {4, 0}, {3, 2}, {2, 0}},
	})

	result, err := PolygonUnion(a, b)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal("union should not be nil")
	}

	poly, ok := result.Geometry.(*geojson.Polygon)
	if !ok {
		mp, ok := result.Geometry.(*geojson.MultiPolygon)
		if ok {
			area := 0.0
			for _, p := range mp.Coordinates {
				area += polygonArea(p[0])
			}
			if math.Abs(area-4.0) > 0.01 {
				t.Errorf("expected area 4.0, got %.2f", area)
			}
			return
		}
		t.Fatalf("expected Polygon, got %T", result.Geometry)
	}
	area := polygonArea(poly.Coordinates[0])
	if math.Abs(area-4.0) > 0.01 {
		t.Errorf("expected area 4.0, got %.2f", area)
	}
}

func polygonArea(ring []geojson.Position) float64 {
	n := len(ring)
	area := 0.0
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		area += ring[i][0] * ring[j][1]
		area -= ring[j][0] * ring[i][1]
	}
	return math.Abs(area) / 2.0
}

func TestUnionTrianglePartiallyOverlapping(t *testing.T) {
	a := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {3, 0}, {1.5, 2}, {0, 0}},
	})
	b := geojson.NewPolygon([][]geojson.Position{
		{{1.5, 0}, {4, 0}, {2.5, 2}, {1.5, 0}},
	})

	result, err := PolygonUnion(a, b)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal("union should not be nil")
	}
}

func TestUnionFromFeatureCollection(t *testing.T) {
	a := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {2, 0}, {2, 2}, {0, 2}, {0, 0}},
	})
	b := geojson.NewPolygon([][]geojson.Position{
		{{1, 1}, {3, 1}, {3, 3}, {1, 3}, {1, 1}},
	})

	aFeat := geojson.NewFeature(a, nil)
	bFeat := geojson.NewFeature(b, nil)

	result, err := PolygonUnion(aFeat, bFeat)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal("union should not be nil")
	}
}

func TestPolygonWithHole(t *testing.T) {
	outer := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {5, 0}, {5, 5}, {0, 5}, {0, 0}},
	})
	inner := geojson.NewPolygon([][]geojson.Position{
		{{1, 1}, {4, 1}, {4, 4}, {1, 4}, {1, 1}},
	})

	result, err := PolygonDifference(outer, inner)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal("difference should not be nil")
	}
	poly, ok := result.Geometry.(*geojson.Polygon)
	if !ok {
		t.Fatalf("expected Polygon, got %T", result.Geometry)
	}
	if len(poly.Coordinates) < 2 {
		t.Error("difference with contained polygon should produce a hole")
	}
	outerArea := polygonArea(poly.Coordinates[0])
	if math.Abs(outerArea-25.0) > 0.01 {
		t.Errorf("expected outer area 25.0, got %.2f", outerArea)
	}
	if len(poly.Coordinates) > 1 {
		holeArea := polygonArea(poly.Coordinates[1])
		if math.Abs(holeArea-9.0) > 0.01 {
			t.Errorf("expected hole area 9.0, got %.2f", holeArea)
		}
	}
}

func TestPolygonTouching(t *testing.T) {
	a := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}},
	})
	b := geojson.NewPolygon([][]geojson.Position{
		{{1, 0}, {2, 0}, {2, 1}, {1, 1}, {1, 0}},
	})

	result, err := PolygonUnion(a, b)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal("union of touching should not be nil")
	}
	poly, ok := result.Geometry.(*geojson.Polygon)
	if !ok {
		t.Fatalf("expected Polygon, got %T", result.Geometry)
	}
	area := polygonArea(poly.Coordinates[0])
	if math.Abs(area-2.0) > 0.01 {
		t.Errorf("expected area 2.0, got %.2f", area)
	}
}

func TestInvalidInput(t *testing.T) {
	pt := geojson.NewPoint(geojson.Position{1, 2})
	_, err := PolygonUnion(pt, pt)
	if err == nil {
		t.Error("expected error for point input")
	}
}

func BenchmarkPolygonUnion(b *testing.B) {
	a := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {2, 0}, {2, 2}, {0, 2}, {0, 0}},
	})
	c := geojson.NewPolygon([][]geojson.Position{
		{{1, 1}, {3, 1}, {3, 3}, {1, 3}, {1, 1}},
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		PolygonUnion(a, c)
	}
}
