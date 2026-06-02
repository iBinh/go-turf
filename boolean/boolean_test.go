package boolean

import (
	"testing"

	"github.com/ibinh/turf-go/geojson"
)

func pt(lng, lat float64) geojson.Position {
	return geojson.Position{lng, lat}
}

func TestClockwise(t *testing.T) {
	cwRing := []geojson.Position{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}}
	if !Clockwise(cwRing) {
		t.Error("expected (0,0)->(0,1)->(1,1)->(1,0) to be CW")
	}
	ccwRing := []geojson.Position{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}}
	if Clockwise(ccwRing) {
		t.Error("expected (0,0)->(1,0)->(1,1)->(0,1) to be CCW")
	}
}

func TestPointInPolygon(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}},
	})
	tests := []struct {
		pt       geojson.Position
		expected bool
	}{
		{pt(5, 5), true},
		{pt(0, 0), true},
		{pt(10, 10), true},
		{pt(-1, 5), false},
		{pt(5, -1), false},
		{pt(15, 5), false},
	}
	for _, tc := range tests {
		result, err := PointInPolygon(geojson.NewPoint(tc.pt), poly)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != tc.expected {
			t.Errorf("PointInPolygon(%v) = %v, want %v", tc.pt, result, tc.expected)
		}
	}
}

func TestPointInPolygonWithHole(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}},
		{{2, 2}, {2, 8}, {8, 8}, {8, 2}, {2, 2}},
	})
	tests := []struct {
		pt       geojson.Position
		expected bool
	}{
		{pt(5, 5), false},
		{pt(1, 1), true},
		{pt(5, 1), true},
	}
	for _, tc := range tests {
		result, err := PointInPolygon(geojson.NewPoint(tc.pt), poly)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != tc.expected {
			t.Errorf("PointInPolygon(%v) = %v, want %v", tc.pt, result, tc.expected)
		}
	}
}

func TestPointInMultiPolygon(t *testing.T) {
	mpoly := geojson.NewMultiPolygon([][][]geojson.Position{
		{{{0, 0}, {0, 5}, {5, 5}, {5, 0}, {0, 0}}},
		{{{10, 10}, {10, 15}, {15, 15}, {15, 10}, {10, 10}}},
	})
	tests := []struct {
		pt       geojson.Position
		expected bool
	}{
		{pt(2, 2), true},
		{pt(12, 12), true},
		{pt(7, 7), false},
	}
	for _, tc := range tests {
		result, err := PointInPolygon(geojson.NewPoint(tc.pt), mpoly)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != tc.expected {
			t.Errorf("PointInPolygon(%v) = %v, want %v", tc.pt, result, tc.expected)
		}
	}
}

func TestPointOnLine(t *testing.T) {
	line := geojson.NewLineString([]geojson.Position{{0, 0}, {10, 10}})
	tests := []struct {
		pt             geojson.Position
		ignoreEndpoints bool
		expected       bool
	}{
		{pt(0, 0), false, true},
		{pt(10, 10), false, true},
		{pt(5, 5), false, true},
		{pt(5, 6), false, false},
		{pt(0, 0), true, false},
		{pt(10, 10), true, false},
		{pt(5, 5), true, true},
	}
	for _, tc := range tests {
		result, err := PointOnLine(geojson.NewPoint(tc.pt), line, tc.ignoreEndpoints)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != tc.expected {
			t.Errorf("PointOnLine(%v, ignoreEndpoints=%v) = %v, want %v",
				tc.pt, tc.ignoreEndpoints, result, tc.expected)
		}
	}
}

func TestPointOnPolygonRing(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}},
	})
	result, err := PointOnLine(geojson.NewPoint(pt(0, 5)), poly, false)
	if err != nil {
		t.Fatal(err)
	}
	if !result {
		t.Error("point on polygon edge should be on line")
	}
}

func TestSegmentIntersect(t *testing.T) {
	tests := []struct {
		name     string
		a, b, c, d geojson.Position
		expected bool
	}{
		{"crossing", pt(0, 0), pt(10, 10), pt(0, 10), pt(10, 0), true},
		{"parallel", pt(0, 0), pt(10, 0), pt(0, 5), pt(10, 5), false},
		{"touching endpoint", pt(0, 0), pt(10, 0), pt(10, 0), pt(10, 10), true},
		{"disjoint", pt(0, 0), pt(5, 5), pt(6, 6), pt(10, 10), false},
		{"collinear overlap", pt(0, 0), pt(10, 10), pt(5, 5), pt(15, 15), true},
	}
	for _, tc := range tests {
		result := SegmentIntersect(tc.a, tc.b, tc.c, tc.d)
		if result != tc.expected {
			t.Errorf("%s: SegmentIntersect = %v, want %v", tc.name, result, tc.expected)
		}
	}
}

func TestContains(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}},
	})
	insidePt := geojson.NewPoint(pt(5, 5))
	outsidePt := geojson.NewPoint(pt(15, 15))

	ok, err := Contains(poly, insidePt)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("polygon should contain interior point")
	}

	ok, err = Contains(poly, outsidePt)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Error("polygon should not contain exterior point")
	}
}

func TestContainsLineString(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}},
	})
	insideLine := geojson.NewLineString([]geojson.Position{{1, 1}, {2, 2}, {3, 3}})
	ok, err := Contains(poly, insideLine)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("polygon should contain interior line")
	}
}

func TestContainsPolygon(t *testing.T) {
	outer := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}},
	})
	inner := geojson.NewPolygon([][]geojson.Position{
		{{1, 1}, {1, 9}, {9, 9}, {9, 1}, {1, 1}},
	})
	ok, err := Contains(outer, inner)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("outer polygon should contain inner polygon")
	}
}

func TestIntersectsPointPolygon(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}},
	})
	inside := geojson.NewPoint(pt(5, 5))
	outside := geojson.NewPoint(pt(15, 15))

	ok, err := Intersects(poly, inside)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("polygon should intersect interior point")
	}

	ok, err = Intersects(poly, outside)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Error("polygon should not intersect exterior point")
	}
}

func TestIntersectsLineLine(t *testing.T) {
	l1 := geojson.NewLineString([]geojson.Position{{0, 0}, {10, 10}})
	l2 := geojson.NewLineString([]geojson.Position{{0, 10}, {10, 0}})
	ok, err := Intersects(l1, l2)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("crossing lines should intersect")
	}

	l3 := geojson.NewLineString([]geojson.Position{{0, 0}, {5, 5}})
	ok, err = Intersects(l1, l3)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("overlapping lines should intersect")
	}
}

func TestDisjoint(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}},
	})
	pt := geojson.NewPoint(pt(20, 20))
	ok, err := Disjoint(poly, pt)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("far point should be disjoint from polygon")
	}
}

func TestTouches(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}},
	})
	edgePoint := geojson.NewPoint(pt(5, 0))
	touches, err := Touches(poly, edgePoint)
	if err != nil {
		t.Fatal(err)
	}
	if touches {
		t.Error("Touches with point should return false")
	}

	line := geojson.NewLineString([]geojson.Position{{-5, 5}, {5, 5}})
	touches, err = Touches(poly, line)
	if err != nil {
		t.Fatal(err)
	}
	if !touches {
		t.Error("line touching polygon boundary should be Touches=true")
	}
}

func TestCrosses(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}},
	})
	line := geojson.NewLineString([]geojson.Position{{-5, 5}, {15, 5}})
	crosses, err := Crosses(poly, line)
	if err != nil {
		t.Fatal(err)
	}
	if !crosses {
		t.Error("line crossing polygon should be Crosses=true")
	}
}

func TestOverlap(t *testing.T) {
	p1 := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}},
	})
	p2 := geojson.NewPolygon([][]geojson.Position{
		{{5, 5}, {5, 15}, {15, 15}, {15, 5}, {5, 5}},
	})
	overlap, err := Overlap(p1, p2)
	if err != nil {
		t.Fatal(err)
	}
	if !overlap {
		t.Error("overlapping polygons should return true")
	}

	contained := geojson.NewPolygon([][]geojson.Position{
		{{1, 1}, {1, 9}, {9, 9}, {9, 1}, {1, 1}},
	})
	overlap, err = Overlap(p1, contained)
	if err != nil {
		t.Fatal(err)
	}
	if overlap {
		t.Error("contained polygon should not be overlap")
	}
}

func TestValid(t *testing.T) {
	validPt := geojson.NewPoint(pt(1, 2))
	ok, err := Valid(validPt)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("valid point should be valid")
	}

	invalidRing := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {0, 10}, {10, 10}},
	})
	ok, err = Valid(invalidRing)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Error("unclosed ring should be invalid")
	}

	emptyLine := geojson.NewLineString([]geojson.Position{})
	ok, err = Valid(emptyLine)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Error("empty line should be invalid")
	}
}

func TestValidMultiPolygon(t *testing.T) {
	mp := geojson.NewMultiPolygon([][][]geojson.Position{
		{{{0, 0}, {0, 5}, {5, 5}, {5, 0}, {0, 0}}},
	})
	ok, err := Valid(mp)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("valid multipolygon should be valid")
	}
}

func TestConcave(t *testing.T) {
	convex := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}},
	})
	result, err := Concave(convex)
	if err != nil {
		t.Fatal(err)
	}
	if result {
		t.Error("rectangle should not be concave")
	}

	concavePoly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {0, 10}, {5, 5}, {10, 10}, {10, 0}, {0, 0}},
	})
	result, err = Concave(concavePoly)
	if err != nil {
		t.Fatal(err)
	}
	if !result {
		t.Error("polygon with indentation should be concave")
	}
}

func TestWithin(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}},
	})
	insidePt := geojson.NewPoint(pt(5, 5))
	outsidePt := geojson.NewPoint(pt(15, 15))

	ok, err := Within(insidePt, poly)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("point should be within polygon")
	}

	ok, err = Within(outsidePt, poly)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Error("outside point should not be within polygon")
	}
}

func TestPointInPolygonEdgeCases(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}},
	})
	onEdge := geojson.NewPoint(pt(5, 0))
	result, err := PointInPolygon(onEdge, poly)
	if err != nil {
		t.Fatal(err)
	}
	if !result {
		t.Error("point on edge should be considered inside")
	}
}

func TestLineInPolygonTouchingEdge(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}},
	})
	line := geojson.NewLineString([]geojson.Position{{-5, 5}, {5, 5}})
	ok, err := Contains(poly, line)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Error("line partly outside polygon should not be contained")
	}
}

func BenchmarkPointInPolygon(b *testing.B) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {0, 100}, {100, 100}, {100, 0}, {0, 0}},
	})
	pt := geojson.NewPoint(pt(50, 50))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		PointInPolygon(pt, poly)
	}
}

func TestContainsMultiPoint(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}},
	})
	mp := geojson.NewMultiPoint([]geojson.Position{{1, 1}, {2, 2}, {3, 3}})
	ok, err := Contains(poly, mp)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("polygon should contain all interior multipoints")
	}

	mp2 := geojson.NewMultiPoint([]geojson.Position{{1, 1}, {15, 15}})
	ok, err = Contains(poly, mp2)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Error("polygon should not contain multipoint with exterior point")
	}
}

func TestSegmentIntersectCollinear(t *testing.T) {
	result := SegmentIntersect(
		geojson.Position{0, 0}, geojson.Position{10, 10},
		geojson.Position{15, 15}, geojson.Position{20, 20},
	)
	if result {
		t.Error("collinear non-overlapping segments should not intersect")
	}

	result = SegmentIntersect(
		geojson.Position{0, 0}, geojson.Position{10, 10},
		geojson.Position{5, 5}, geojson.Position{15, 15},
	)
	if !result {
		t.Error("collinear overlapping segments should intersect")
	}
}

func TestSegmentIntersectEndpoint(t *testing.T) {
	result := SegmentIntersect(
		geojson.Position{0, 0}, geojson.Position{10, 0},
		geojson.Position{10, 0}, geojson.Position{10, 10},
	)
	if !result {
		t.Error("segments touching at endpoint should intersect")
	}
}

func TestClockwiseEdge(t *testing.T) {
	ring := []geojson.Position{{0, 0}, {1, 1}, {1, 0}, {0, 0}}
	if !Clockwise(ring) {
		t.Error("expected triangle (0,0)->(1,1)->(1,0) to be CW")
	}
}

func TestBooleanEqual(t *testing.T) {
	l1 := geojson.NewLineString([]geojson.Position{{0, 0}, {10, 10}})
	l2 := geojson.NewLineString([]geojson.Position{{0, 0}, {10, 10}})
	l3 := geojson.NewLineString([]geojson.Position{{0, 0}, {10, 0}})

	eq, err := BooleanEqual(l1, l2)
	if err != nil {
		t.Fatal(err)
	}
	if !eq {
		t.Error("identical lines should be equal")
	}

	eq, err = BooleanEqual(l1, l3)
	if err != nil {
		t.Fatal(err)
	}
	if eq {
		t.Error("different lines should not be equal")
	}
}

func TestBooleanEqualPoint(t *testing.T) {
	p1 := geojson.NewPoint(pt(5, 5))
	p2 := geojson.NewPoint(pt(5, 5))
	p3 := geojson.NewPoint(pt(6, 6))

	eq, err := BooleanEqual(p1, p2)
	if err != nil {
		t.Fatal(err)
	}
	if !eq {
		t.Error("identical points should be equal")
	}

	eq, err = BooleanEqual(p1, p3)
	if err != nil {
		t.Fatal(err)
	}
	if eq {
		t.Error("different points should not be equal")
	}
}

func TestBooleanEqualPolygon(t *testing.T) {
	p1 := geojson.NewPolygon([][]geojson.Position{{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}}})
	p2 := geojson.NewPolygon([][]geojson.Position{{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}}})
	eq, err := BooleanEqual(p1, p2)
	if err != nil {
		t.Fatal(err)
	}
	if !eq {
		t.Error("identical polygons should be equal")
	}
}

func TestBooleanParallel(t *testing.T) {
	l1 := geojson.NewLineString([]geojson.Position{{0, 0}, {10, 0}})
	l2 := geojson.NewLineString([]geojson.Position{{0, 5}, {10, 5}})
	parallel, err := BooleanParallel(l1, l2)
	if err != nil {
		t.Fatal(err)
	}
	if !parallel {
		t.Error("horizontal lines should be parallel")
	}
}

func TestBooleanParallelNotParallel(t *testing.T) {
	l1 := geojson.NewLineString([]geojson.Position{{0, 0}, {10, 0}})
	l2 := geojson.NewLineString([]geojson.Position{{0, 0}, {10, 10}})
	parallel, err := BooleanParallel(l1, l2)
	if err != nil {
		t.Fatal(err)
	}
	if parallel {
		t.Error("non-parallel lines should not be parallel")
	}
}

func TestBooleanParallelDiagonal(t *testing.T) {
	l1 := geojson.NewLineString([]geojson.Position{{0, 0}, {10, 10}})
	l2 := geojson.NewLineString([]geojson.Position{{0, 5}, {10, 15}})
	parallel, err := BooleanParallel(l1, l2)
	if err != nil {
		t.Fatal(err)
	}
	if !parallel {
		t.Error("diagonal lines with same slope should be parallel")
	}
}

func TestBooleanEqualFeature(t *testing.T) {
	f1 := geojson.NewFeature(geojson.NewPoint(pt(1, 2)), nil)
	f2 := geojson.NewFeature(geojson.NewPoint(pt(1, 2)), nil)
	eq, err := BooleanEqual(f1, f2)
	if err != nil {
		t.Fatal(err)
	}
	if !eq {
		t.Error("features with identical geometry should be equal")
	}
}

func TestConcaveNonPolygon(t *testing.T) {
	line := geojson.NewLineString([]geojson.Position{{0, 0}, {1, 1}})
	result, err := Concave(line)
	if err != nil {
		t.Fatal(err)
	}
	if result {
		t.Error("line should not be concave (not a polygon)")
	}
}

func TestValidPoint(t *testing.T) {
	pt := geojson.NewPoint(geojson.Position{1})
	ok, err := Valid(pt)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Error("point with 1 coordinate should be invalid")
	}
}

func TestWithinExactMatch(t *testing.T) {
	p1 := geojson.NewPoint(pt(5, 5))
	p2 := geojson.NewPoint(pt(5, 5))
	ok, err := Within(p1, p2)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("same point should be within")
	}
}

func TestTouchesLineLine(t *testing.T) {
	l1 := geojson.NewLineString([]geojson.Position{{0, 0}, {0, 10}})
	l2 := geojson.NewLineString([]geojson.Position{{0, 10}, {10, 10}})
	touches, err := Touches(l1, l2)
	if err != nil {
		t.Fatal(err)
	}
	if !touches {
		t.Error("lines touching at endpoint should be Touches=true")
	}
}

func TestLineInPolygonAllPoints(t *testing.T) {
	line := geojson.NewLineString([]geojson.Position{{1, 1}, {2, 2}})
	pt := geojson.NewPoint(pt(1, 1))
	in, err := PointInPolygon(pt, line)
	if err != nil {
		t.Fatal(err)
	}
	if !in {
		t.Error("point on line endpoint should be in")
	}
}

func TestMultiPointInPolygon(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}},
	})
	mp := geojson.NewMultiPoint([]geojson.Position{{5, 5}})
	in, err := PointInPolygon(mp, poly)
	if err != nil {
		t.Fatal(err)
	}
	if !in {
		t.Error("MultiPoint should be in polygon")
	}
}

func TestPointOnMultiLineString(t *testing.T) {
	ml := geojson.NewMultiLineString([][]geojson.Position{
		{{0, 0}, {0, 10}},
		{{10, 0}, {10, 10}},
	})
	on, err := PointOnLine(geojson.NewPoint(pt(0, 5)), ml, false)
	if err != nil {
		t.Fatal(err)
	}
	if !on {
		t.Error("point on first line of MLS should be on line")
	}

	on, err = PointOnLine(geojson.NewPoint(pt(10, 5)), ml, false)
	if err != nil {
		t.Fatal(err)
	}
	if !on {
		t.Error("point on second line of MLS should be on line")
	}
}

func TestIntersectsLinePolygon(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}},
	})
	line := geojson.NewLineString([]geojson.Position{{-5, 5}, {15, 5}})
	ok, err := Intersects(line, poly)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("line crossing polygon should intersect")
	}
}

func TestIntersectsPolygonMultiPolygon(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}},
	})
	mp := geojson.NewMultiPolygon([][][]geojson.Position{
		{{{5, 5}, {5, 15}, {15, 15}, {15, 5}, {5, 5}}},
	})
	ok, err := Intersects(poly, mp)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("overlapping polygon and multipolygon should intersect")
	}

	mp2 := geojson.NewMultiPolygon([][][]geojson.Position{
		{{{20, 20}, {20, 30}, {30, 30}, {30, 20}, {20, 20}}},
	})
	ok, err = Intersects(poly, mp2)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Error("non-overlapping polygon and multipolygon should not intersect")
	}
}

func TestContainsMultiPolygonInPolygon(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {0, 30}, {30, 30}, {30, 0}, {0, 0}},
	})
	mp := geojson.NewMultiPolygon([][][]geojson.Position{
		{{{1, 1}, {1, 5}, {5, 5}, {5, 1}, {1, 1}}},
		{{{10, 10}, {10, 15}, {15, 15}, {15, 10}, {10, 10}}},
	})
	ok, err := Contains(poly, mp)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("polygon should contain multipolygon")
	}
}

func TestCrossesLineLine(t *testing.T) {
	l1 := geojson.NewLineString([]geojson.Position{{0, 0}, {10, 10}})
	l2 := geojson.NewLineString([]geojson.Position{{0, 10}, {10, 0}})
	crosses, err := Crosses(l1, l2)
	if err != nil {
		t.Fatal(err)
	}
	if !crosses {
		t.Error("crossing lines should be Crosses=true")
	}
}

func TestOverlapSameTypeOnly(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}},
	})
	line := geojson.NewLineString([]geojson.Position{{-5, 5}, {15, 5}})
	overlap, err := Overlap(poly, line)
	if err != nil {
		t.Fatal(err)
	}
	if overlap {
		t.Error("Overlap should be false for different geometry types")
	}
}
