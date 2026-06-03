package bbox

import (
	"math"
	"testing"

	"github.com/ibinh/turf-go/geojson"
)

// ---------------------------------------------------------------------------
// Existing tests (unchanged)
// ---------------------------------------------------------------------------

func TestBBox(t *testing.T) {
	f := geojson.NewFeature(
		geojson.NewPoint(geojson.Position{1, 2}),
		nil,
	)
	b, err := BBox(f)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 4 || b[0] != 1 || b[1] != 2 || b[2] != 1 || b[3] != 2 {
		t.Errorf("unexpected bbox: %v", b)
	}
}

func TestBBoxMultiPoint(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{1, 2}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{3, 5}), nil),
	})
	b, err := BBox(fc)
	if err != nil {
		t.Fatal(err)
	}
	if b[0] != 1 || b[1] != 2 || b[2] != 3 || b[3] != 5 {
		t.Errorf("unexpected bbox: %v", b)
	}
}

func TestBBoxPolygon(t *testing.T) {
	bbox := []float64{0, 0, 10, 10}
	f, err := BBoxPolygon(bbox)
	if err != nil {
		t.Fatal(err)
	}
	poly, ok := f.Geometry.(*geojson.Polygon)
	if !ok {
		t.Fatalf("expected Polygon, got %T", f.Geometry)
	}
	if len(poly.Coordinates[0]) != 5 {
		t.Errorf("expected 5 vertices, got %d", len(poly.Coordinates[0]))
	}
}

func TestEnvelope(t *testing.T) {
	f := geojson.NewFeature(
		geojson.NewPoint(geojson.Position{5, 5}),
		nil,
	)
	env, err := Envelope(f)
	if err != nil {
		t.Fatal(err)
	}
	if env == nil {
		t.Fatal("expected envelope feature")
	}
}

func TestSquare(t *testing.T) {
	bbox := []float64{0, 0, 2, 1}
	sq, err := Square(bbox)
	if err != nil {
		t.Fatal(err)
	}
	width := sq[2] - sq[0]
	height := sq[3] - sq[1]
	if mathAbs(width-height) > 0.001 {
		t.Errorf("square should have equal width/height, got %f x %f", width, height)
	}
}

func TestSquareTall(t *testing.T) {
	bbox := []float64{0, 0, 1, 3}
	sq, err := Square(bbox)
	if err != nil {
		t.Fatal(err)
	}
	width := sq[2] - sq[0]
	height := sq[3] - sq[1]
	if mathAbs(width-height) > 0.001 {
		t.Errorf("square should have equal width/height, got %f x %f", width, height)
	}
}

func TestBBoxPolygonRoundtrip(t *testing.T) {
	f := geojson.NewFeature(
		geojson.NewPoint(geojson.Position{1, 2}),
		nil,
	)
	b, err := BBox(f)
	if err != nil {
		t.Fatal(err)
	}
	poly, err := BBoxPolygon(b)
	if err != nil {
		t.Fatal(err)
	}
	if poly == nil {
		t.Fatal("expected polygon")
	}
}

func mathAbs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// ---------------------------------------------------------------------------
// clipSegment tests (unexported, same package)
// ---------------------------------------------------------------------------

func TestClipSegment_TrivialAccept(t *testing.T) {
	// Both endpoints inside the bbox → returned as-is
	a := geojson.Position{2, 2}
	b := geojson.Position{8, 8}
	result := clipSegment(a, b, 0, 0, 10, 10)
	if len(result) != 2 {
		t.Fatalf("expected 2 points, got %d", len(result))
	}
	if result[0][0] != 2 || result[0][1] != 2 {
		t.Errorf("expected (2,2), got %v", result[0])
	}
	if result[1][0] != 8 || result[1][1] != 8 {
		t.Errorf("expected (8,8), got %v", result[1])
	}
}

func TestClipSegment_TrivialReject(t *testing.T) {
	tests := []struct {
		name string
		a, b geojson.Position
	}{
		{"both left", geojson.Position{-5, 5}, geojson.Position{-2, 5}},
		{"both right", geojson.Position{12, 5}, geojson.Position{15, 5}},
		{"both below", geojson.Position{5, -5}, geojson.Position{5, -2}},
		{"both above", geojson.Position{5, 12}, geojson.Position{5, 15}},
		{"both left-bottom", geojson.Position{-5, -5}, geojson.Position{-2, -2}},
		{"both right-top", geojson.Position{12, 12}, geojson.Position{15, 15}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := clipSegment(tc.a, tc.b, 0, 0, 10, 10)
			if result != nil {
				t.Errorf("expected nil, got %v", result)
			}
		})
	}
}

func TestClipSegment_LeftClip(t *testing.T) {
	// A is left of bbox, B is inside → clips to left edge
	a := geojson.Position{-5, 5}
	b := geojson.Position{5, 5}
	result := clipSegment(a, b, 0, 0, 10, 10)
	if len(result) != 2 {
		t.Fatalf("expected 2 points, got %d", len(result))
	}
	expectPoint(t, result[0], 0, 5)
	expectPoint(t, result[1], 5, 5)
}

func TestClipSegment_RightClip(t *testing.T) {
	// A is inside, B is right of bbox → clips to right edge
	a := geojson.Position{5, 5}
	b := geojson.Position{15, 5}
	result := clipSegment(a, b, 0, 0, 10, 10)
	if len(result) != 2 {
		t.Fatalf("expected 2 points, got %d", len(result))
	}
	expectPoint(t, result[0], 5, 5)
	expectPoint(t, result[1], 10, 5)
}

func TestClipSegment_BottomClip(t *testing.T) {
	a := geojson.Position{5, -5}
	b := geojson.Position{5, 5}
	result := clipSegment(a, b, 0, 0, 10, 10)
	if len(result) != 2 {
		t.Fatalf("expected 2 points, got %d", len(result))
	}
	expectPoint(t, result[0], 5, 0)
	expectPoint(t, result[1], 5, 5)
}

func TestClipSegment_TopClip(t *testing.T) {
	a := geojson.Position{5, 5}
	b := geojson.Position{5, 15}
	result := clipSegment(a, b, 0, 0, 10, 10)
	if len(result) != 2 {
		t.Fatalf("expected 2 points, got %d", len(result))
	}
	expectPoint(t, result[0], 5, 5)
	expectPoint(t, result[1], 5, 10)
}

func TestClipSegment_BothOutsideDifferentRegions(t *testing.T) {
	// A is left of bbox, B is right of bbox → needs two clips
	a := geojson.Position{-5, 5}
	b := geojson.Position{15, 5}
	result := clipSegment(a, b, 0, 0, 10, 10)
	if len(result) != 2 {
		t.Fatalf("expected 2 points, got %d", len(result))
	}
	expectPoint(t, result[0], 0, 5)
	expectPoint(t, result[1], 10, 5)
}

func TestClipSegment_LeftBottomToInside(t *testing.T) {
	// A is left+bottom, B is inside → clips two edges iteratively
	a := geojson.Position{-5, -5}
	b := geojson.Position{5, 5}
	result := clipSegment(a, b, 0, 0, 10, 10)
	if len(result) != 2 {
		t.Fatalf("expected 2 points, got %d", len(result))
	}
	expectPoint(t, result[0], 0, 0)
	expectPoint(t, result[1], 5, 5)
}

func TestClipSegment_CrossesCorner(t *testing.T) {
	// A is left of bbox, B is above bbox → crosses top-left corner
	a := geojson.Position{-5, 15}
	b := geojson.Position{15, -5}
	result := clipSegment(a, b, 0, 0, 10, 10)
	if len(result) != 2 {
		t.Fatalf("expected 2 points, got %d", len(result))
	}
	// The clipped segment should go through the bbox from left to bottom
	expectPoint(t, result[0], 0, 10)
	expectPoint(t, result[1], 10, 0)
}

func TestClipSegment_VerticalLine(t *testing.T) {
	// Vertical line crossing bottom boundary
	a := geojson.Position{5, -5}
	b := geojson.Position{5, 15}
	result := clipSegment(a, b, 0, 0, 10, 10)
	if len(result) != 2 {
		t.Fatalf("expected 2 points, got %d", len(result))
	}
	expectPoint(t, result[0], 5, 0)
	expectPoint(t, result[1], 5, 10)
}

func TestClipSegment_HorizontalLine(t *testing.T) {
	// Horizontal line crossing left and right boundaries
	a := geojson.Position{-5, 5}
	b := geojson.Position{15, 5}
	result := clipSegment(a, b, 0, 0, 10, 10)
	if len(result) != 2 {
		t.Fatalf("expected 2 points, got %d", len(result))
	}
	expectPoint(t, result[0], 0, 5)
	expectPoint(t, result[1], 10, 5)
}

// ---------------------------------------------------------------------------
// clipLineString tests (unexported, same package)
// ---------------------------------------------------------------------------

func TestClipLineString_AllInside(t *testing.T) {
	line := []geojson.Position{{2, 2}, {5, 5}, {8, 8}}
	segments := clipLineString(line, 0, 0, 10, 10)
	if len(segments) != 1 {
		t.Fatalf("expected 1 segment, got %d", len(segments))
	}
	if len(segments[0]) != 3 {
		t.Fatalf("expected 3 points, got %d", len(segments[0]))
	}
	expectPoint(t, segments[0][0], 2, 2)
	expectPoint(t, segments[0][1], 5, 5)
	expectPoint(t, segments[0][2], 8, 8)
}

func TestClipLineString_CrossingBoundary(t *testing.T) {
	// Line enters from left, crosses through, exits right
	line := []geojson.Position{{-5, 5}, {5, 5}, {15, 5}}
	segments := clipLineString(line, 0, 0, 10, 10)
	if len(segments) != 1 {
		t.Fatalf("expected 1 segment, got %d", len(segments))
	}
	if len(segments[0]) != 3 {
		t.Fatalf("expected 3 points, got %d", len(segments[0]))
	}
	expectPoint(t, segments[0][0], 0, 5)
	expectPoint(t, segments[0][1], 5, 5)
	expectPoint(t, segments[0][2], 10, 5)
}

func TestClipLineString_AllOutside(t *testing.T) {
	line := []geojson.Position{{-5, -5}, {-2, -2}}
	segments := clipLineString(line, 0, 0, 10, 10)
	if len(segments) != 0 {
		t.Fatalf("expected 0 segments, got %d", len(segments))
	}
}

func TestClipLineString_DisconnectedSegments(t *testing.T) {
	// Line goes: inside → outside segment → inside (produces 2 segments)
	line := []geojson.Position{
		{-5, 5},    // outside left
		{5, 5},     // inside
		{15, -5},   // outside right+bottom (trivial reject from next segment)
		{15, 15},   // outside right+top
		{5, 5},     // inside again
		{-5, 5},    // outside left
	}
	segments := clipLineString(line, 0, 0, 10, 10)
	// Should produce 2 segments: the first going from (0,5) to (5,5) to (10,0),
	// and the second from (10,10) to (5,5) to (0,5)
	if len(segments) != 2 {
		t.Fatalf("expected 2 segments, got %d", len(segments))
	}
	// First segment: enters at left edge, exits through right+bottom
	if len(segments[0]) < 2 {
		t.Fatalf("first segment too short: %d", len(segments[0]))
	}
	if len(segments[1]) < 2 {
		t.Fatalf("second segment too short: %d", len(segments[1]))
	}
}

func TestClipLineString_SingleSegmentClipped(t *testing.T) {
	// Single segment crossing left boundary
	line := []geojson.Position{{-5, 5}, {5, 5}}
	segments := clipLineString(line, 0, 0, 10, 10)
	if len(segments) != 1 {
		t.Fatalf("expected 1 segment, got %d", len(segments))
	}
	if len(segments[0]) != 2 {
		t.Fatalf("expected 2 points, got %d", len(segments[0]))
	}
	expectPoint(t, segments[0][0], 0, 5)
	expectPoint(t, segments[0][1], 5, 5)
}

func TestClipLineString_EntersExitsReenters(t *testing.T) {
	// Line enters bbox, exits, re-enters (multiple segments)
	line := []geojson.Position{
		{-5, 5},   // outside
		{5, 5},    // inside
		{15, 5},   // outside
		{25, 5},   // still outside
	}
	segments := clipLineString(line, 0, 0, 10, 10)
	if len(segments) != 1 {
		t.Fatalf("expected 1 segment, got %d", len(segments))
	}
	// Should have [(0,5), (5,5), (10,5)]
	if len(segments[0]) != 3 {
		t.Fatalf("expected 3 points, got %d", len(segments[0]))
	}
	expectPoint(t, segments[0][0], 0, 5)
	expectPoint(t, segments[0][1], 5, 5)
	expectPoint(t, segments[0][2], 10, 5)
}

// ---------------------------------------------------------------------------
// clipPolygonRing tests (unexported, same package)
// ---------------------------------------------------------------------------

func TestClipPolygonRing_AllInside(t *testing.T) {
	ring := []geojson.Position{{2, 2}, {8, 2}, {8, 8}, {2, 8}, {2, 2}}
	result := clipPolygonRing(ring, 0, 0, 10, 10)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if len(result) != 5 {
		t.Fatalf("expected 5 vertices, got %d", len(result))
	}
	// Vertices should be unchanged
	expectPoint(t, result[0], 2, 2)
	expectPoint(t, result[1], 8, 2)
	expectPoint(t, result[2], 8, 8)
	expectPoint(t, result[3], 2, 8)
	expectPoint(t, result[4], 2, 2)
}

func TestClipPolygonRing_PartiallyClipped(t *testing.T) {
	// Ring larger than bbox → gets clipped to bbox boundary
	ring := []geojson.Position{{-5, -5}, {15, -5}, {15, 15}, {-5, 15}, {-5, -5}}
	result := clipPolygonRing(ring, 0, 0, 10, 10)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	// The resulting ring should represent the intersection
	if len(result) < 3 {
		t.Fatalf("expected at least 3 vertices, got %d", len(result))
	}
	// All points should be within the bbox
	for i, p := range result {
		if p[0] < 0-1e-10 || p[0] > 10+1e-10 || p[1] < 0-1e-10 || p[1] > 10+1e-10 {
			t.Errorf("point %d (%v) is outside bbox [0,0,10,10]", i, p)
		}
	}
}

func TestClipPolygonRing_AllOutside(t *testing.T) {
	ring := []geojson.Position{{-20, -20}, {-15, -20}, {-15, -15}, {-20, -15}, {-20, -20}}
	result := clipPolygonRing(ring, 0, 0, 10, 10)
	if result != nil {
		t.Fatalf("expected nil, got %v", result)
	}
}

func TestClipPolygonRing_OnBoundary(t *testing.T) {
	// Ring exactly on the bbox boundary
	ring := []geojson.Position{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}}
	result := clipPolygonRing(ring, 0, 0, 10, 10)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if len(result) < 3 {
		t.Fatalf("expected at least 3 vertices, got %d", len(result))
	}
}

func TestClipPolygonRing_ClippedToRectangle(t *testing.T) {
	// Large CCW ring covering the full bbox → expected clipped result
	ring := []geojson.Position{{-5, -5}, {15, -5}, {15, 15}, {-5, 15}, {-5, -5}}
	result := clipPolygonRing(ring, 0, 0, 10, 10)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	// Verify the result represents the bbox rectangle
	// Expected vertices (in CCW order after Sutherland-Hodgman):
	// (0,10), (0,0), (10,0), (10,10)
	if len(result) != 4 {
		t.Fatalf("expected 4 vertices, got %d: %v", len(result), result)
	}
	expectPoint(t, result[0], 0, 10)
	expectPoint(t, result[1], 0, 0)
	expectPoint(t, result[2], 10, 0)
	expectPoint(t, result[3], 10, 10)
}

func TestClipPolygonRing_RingFullyInsideSmall(t *testing.T) {
	// Small ring fully inside bbox
	ring := []geojson.Position{{3, 3}, {7, 3}, {7, 7}, {3, 7}, {3, 3}}
	result := clipPolygonRing(ring, 0, 0, 10, 10)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if len(result) != 5 {
		t.Fatalf("expected 5 vertices, got %d", len(result))
	}
}

// ---------------------------------------------------------------------------
// BBoxClip tests
// ---------------------------------------------------------------------------

func TestBBoxClip_PointInside(t *testing.T) {
	f := geojson.NewFeature(geojson.NewPoint(geojson.Position{5, 5}), nil)
	bbox := []float64{0, 0, 10, 10}

	result, err := BBoxClip(f, bbox)
	if err != nil {
		t.Fatal(err)
	}
	p, ok := result.Geometry.(*geojson.Point)
	if !ok {
		t.Fatalf("expected *Point, got %T", result.Geometry)
	}
	expectPoint(t, p.Coordinates, 5, 5)
}

func TestBBoxClip_PointOutside(t *testing.T) {
	f := geojson.NewFeature(geojson.NewPoint(geojson.Position{-5, -5}), nil)
	bbox := []float64{0, 0, 10, 10}

	_, err := BBoxClip(f, bbox)
	if err == nil {
		t.Fatal("expected error for point outside bbox")
	}
}

func TestBBoxClip_PointOnBoundary(t *testing.T) {
	f := geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 5}), nil)
	bbox := []float64{0, 0, 10, 10}

	result, err := BBoxClip(f, bbox)
	if err != nil {
		t.Fatal(err)
	}
	p, ok := result.Geometry.(*geojson.Point)
	if !ok {
		t.Fatalf("expected *Point, got %T", result.Geometry)
	}
	expectPoint(t, p.Coordinates, 0, 5)
}

func TestBBoxClip_MultiPoint(t *testing.T) {
	// Mix of inside and outside points
	pts := []geojson.Position{{5, 5}, {-5, -5}, {15, 15}, {2, 8}}
	mp := geojson.NewMultiPoint(pts)
	f := geojson.NewFeature(mp, nil)
	bbox := []float64{0, 0, 10, 10}

	result, err := BBoxClip(f, bbox)
	if err != nil {
		t.Fatal(err)
	}
	rp, ok := result.Geometry.(*geojson.MultiPoint)
	if !ok {
		t.Fatalf("expected *MultiPoint, got %T", result.Geometry)
	}
	// Only (5,5) and (2,8) are inside the bbox
	if len(rp.Coordinates) != 2 {
		t.Fatalf("expected 2 points, got %d", len(rp.Coordinates))
	}
}

func TestBBoxClip_MultiPoint_AllOutside(t *testing.T) {
	pts := []geojson.Position{{-5, -5}, {-2, -2}}
	mp := geojson.NewMultiPoint(pts)
	f := geojson.NewFeature(mp, nil)
	bbox := []float64{0, 0, 10, 10}

	_, err := BBoxClip(f, bbox)
	if err == nil {
		t.Fatal("expected error for all points outside bbox")
	}
}

func TestBBoxClip_LineString_Inside(t *testing.T) {
	// Line that enters, travels through, exits, and re-enters the bbox,
	// producing 2+ clipped segments (required by current implementation).
	// The longest segment is returned.
	line := []geojson.Position{
		{-5, 5},   // outside left
		{5, 5},    // inside
		{15, 5},   // outside right
		{15, 15},  // outside right+top → disconnect from prev segment
		{5, 5},    // inside
		{-5, 5},   // outside left
	}
	ls := geojson.NewLineString(line)
	f := geojson.NewFeature(ls, nil)
	bbox := []float64{0, 0, 10, 10}

	result, err := BBoxClip(f, bbox)
	if err != nil {
		t.Fatal(err)
	}
	rls, ok := result.Geometry.(*geojson.LineString)
	if !ok {
		t.Fatalf("expected *LineString, got %T", result.Geometry)
	}
	if len(rls.Coordinates) < 2 {
		t.Fatalf("expected at least 2 points, got %d", len(rls.Coordinates))
	}
}

func TestBBoxClip_LineString_Clipped(t *testing.T) {
	// Line crossing left and right boundaries with a disconnect
	// to ensure 2+ clipped segments.
	line := []geojson.Position{
		{-5, 5},   // outside left
		{5, 5},    // inside
		{15, 5},   // outside right
		{15, 15},  // outside right+top → disconnect from prev
		{5, 5},    // inside
	}
	ls := geojson.NewLineString(line)
	f := geojson.NewFeature(ls, nil)
	bbox := []float64{0, 0, 10, 10}

	result, err := BBoxClip(f, bbox)
	if err != nil {
		t.Fatal(err)
	}
	rls, ok := result.Geometry.(*geojson.LineString)
	if !ok {
		t.Fatalf("expected *LineString, got %T", result.Geometry)
	}
	// Picks the longest segment (3 points: enters at left, goes through, exits at right)
	if len(rls.Coordinates) != 3 {
		t.Fatalf("expected 3 points, got %d: %v", len(rls.Coordinates), rls.Coordinates)
	}
	expectPoint(t, rls.Coordinates[0], 0, 5)
	expectPoint(t, rls.Coordinates[1], 5, 5)
	expectPoint(t, rls.Coordinates[2], 10, 5)
}

func TestBBoxClip_LineString_Outside(t *testing.T) {
	line := []geojson.Position{{-5, -5}, {-2, -2}}
	ls := geojson.NewLineString(line)
	f := geojson.NewFeature(ls, nil)
	bbox := []float64{0, 0, 10, 10}

	_, err := BBoxClip(f, bbox)
	if err == nil {
		t.Fatal("expected error for line outside bbox")
	}
}

func TestBBoxClip_LineString_OnlyOneEndpointInside(t *testing.T) {
	// Line entering bbox from left, exiting right, with a disconnect
	// to produce 2+ clipped segments.
	line := []geojson.Position{
		{-5, 5},   // outside left
		{5, 5},    // inside
		{15, 5},   // outside right
		{15, 15},  // outside right+top → disconnect
		{5, 5},    // re-enter at right
	}
	ls := geojson.NewLineString(line)
	f := geojson.NewFeature(ls, nil)
	bbox := []float64{0, 0, 10, 10}

	result, err := BBoxClip(f, bbox)
	if err != nil {
		t.Fatal(err)
	}
	rls, ok := result.Geometry.(*geojson.LineString)
	if !ok {
		t.Fatalf("expected *LineString, got %T", result.Geometry)
	}
	// The longest segment is the first one through the bbox:
	// enters at (0,5), goes to (5,5), exits at (10,5)
	if len(rls.Coordinates) != 3 {
		t.Fatalf("expected 3 points, got %d", len(rls.Coordinates))
	}
	expectPoint(t, rls.Coordinates[0], 0, 5)
	expectPoint(t, rls.Coordinates[1], 5, 5)
	expectPoint(t, rls.Coordinates[2], 10, 5)
}

func TestBBoxClip_MultiLineString(t *testing.T) {
	lines := [][]geojson.Position{
		{{-5, 5}, {5, 5}},
		{{15, 5}, {25, 5}},
	}
	mls := geojson.NewMultiLineString(lines)
	f := geojson.NewFeature(mls, nil)
	bbox := []float64{0, 0, 10, 10}

	result, err := BBoxClip(f, bbox)
	if err != nil {
		t.Fatal(err)
	}
	rls, ok := result.Geometry.(*geojson.MultiLineString)
	if !ok {
		t.Fatalf("expected *MultiLineString, got %T", result.Geometry)
	}
	// Only the first line has a clipped portion
	if len(rls.Coordinates) != 1 {
		t.Fatalf("expected 1 line, got %d", len(rls.Coordinates))
	}
}

func TestBBoxClip_MultiLineString_AllOutside(t *testing.T) {
	lines := [][]geojson.Position{
		{{-15, 5}, {-5, 5}},
		{{15, 5}, {25, 5}},
	}
	mls := geojson.NewMultiLineString(lines)
	f := geojson.NewFeature(mls, nil)
	bbox := []float64{0, 0, 10, 10}

	_, err := BBoxClip(f, bbox)
	if err == nil {
		t.Fatal("expected error for all lines outside bbox")
	}
}

func TestBBoxClip_Polygon_Inside(t *testing.T) {
	ring := []geojson.Position{{2, 2}, {8, 2}, {8, 8}, {2, 8}, {2, 2}}
	poly := geojson.NewPolygon([][]geojson.Position{ring})
	f := geojson.NewFeature(poly, nil)
	bbox := []float64{0, 0, 10, 10}

	result, err := BBoxClip(f, bbox)
	if err != nil {
		t.Fatal(err)
	}
	rp, ok := result.Geometry.(*geojson.Polygon)
	if !ok {
		t.Fatalf("expected *Polygon, got %T", result.Geometry)
	}
	if len(rp.Coordinates) != 1 {
		t.Fatalf("expected 1 ring, got %d", len(rp.Coordinates))
	}
	if len(rp.Coordinates[0]) != 5 {
		t.Fatalf("expected 5 vertices, got %d", len(rp.Coordinates[0]))
	}
}

func TestBBoxClip_Polygon_Clipped(t *testing.T) {
	// Large polygon covering area larger than bbox → gets clipped
	ring := []geojson.Position{{-5, -5}, {15, -5}, {15, 15}, {-5, 15}, {-5, -5}}
	poly := geojson.NewPolygon([][]geojson.Position{ring})
	f := geojson.NewFeature(poly, nil)
	bbox := []float64{0, 0, 10, 10}

	result, err := BBoxClip(f, bbox)
	if err != nil {
		t.Fatal(err)
	}
	rp, ok := result.Geometry.(*geojson.Polygon)
	if !ok {
		t.Fatalf("expected *Polygon, got %T", result.Geometry)
	}
	// Should be a single clipped ring
	if len(rp.Coordinates) != 1 {
		t.Fatalf("expected 1 ring, got %d", len(rp.Coordinates))
	}
	ringResult := rp.Coordinates[0]
	// All vertices must be inside bbox
	for i, p := range ringResult {
		if p[0] < -1e-10 || p[0] > 10+1e-10 || p[1] < -1e-10 || p[1] > 10+1e-10 {
			t.Errorf("vertex %d (%v) is outside bbox [0,0,10,10]", i, p)
		}
	}
}

func TestBBoxClip_Polygon_Outside(t *testing.T) {
	ring := []geojson.Position{{-20, -20}, {-15, -20}, {-15, -15}, {-20, -15}, {-20, -20}}
	poly := geojson.NewPolygon([][]geojson.Position{ring})
	f := geojson.NewFeature(poly, nil)
	bbox := []float64{0, 0, 10, 10}

	_, err := BBoxClip(f, bbox)
	if err == nil {
		t.Fatal("expected error for polygon outside bbox")
	}
}

func TestBBoxClip_PolygonWithHole(t *testing.T) {
	// Polygon with one hole, all inside the bbox
	outer := []geojson.Position{{1, 1}, {9, 1}, {9, 9}, {1, 9}, {1, 1}}
	inner := []geojson.Position{{3, 3}, {7, 3}, {7, 7}, {3, 7}, {3, 3}}
	poly := geojson.NewPolygon([][]geojson.Position{outer, inner})
	f := geojson.NewFeature(poly, nil)
	bbox := []float64{0, 0, 10, 10}

	result, err := BBoxClip(f, bbox)
	if err != nil {
		t.Fatal(err)
	}
	rp, ok := result.Geometry.(*geojson.Polygon)
	if !ok {
		t.Fatalf("expected *Polygon, got %T", result.Geometry)
	}
	// Should have outer ring + hole ring
	if len(rp.Coordinates) != 2 {
		t.Fatalf("expected 2 rings, got %d", len(rp.Coordinates))
	}
}

func TestBBoxClip_MultiPolygon(t *testing.T) {
	ring1 := []geojson.Position{{1, 1}, {4, 1}, {4, 4}, {1, 4}, {1, 1}}
	ring2 := []geojson.Position{{-10, -10}, {-5, -10}, {-5, -5}, {-10, -5}, {-10, -10}}
	polys := [][][]geojson.Position{
		{ring1},
		{ring2},
	}
	mp := geojson.NewMultiPolygon(polys)
	f := geojson.NewFeature(mp, nil)
	bbox := []float64{0, 0, 10, 10}

	result, err := BBoxClip(f, bbox)
	if err != nil {
		t.Fatal(err)
	}
	rmp, ok := result.Geometry.(*geojson.MultiPolygon)
	if !ok {
		t.Fatalf("expected *MultiPolygon, got %T", result.Geometry)
	}
	// Only ring1 is inside the bbox
	if len(rmp.Coordinates) != 1 {
		t.Fatalf("expected 1 polygon, got %d", len(rmp.Coordinates))
	}
}

func TestBBoxClip_MultiPolygon_AllOutside(t *testing.T) {
	ring1 := []geojson.Position{{-20, -20}, {-15, -20}, {-15, -15}, {-20, -15}, {-20, -20}}
	ring2 := []geojson.Position{{20, 20}, {25, 20}, {25, 25}, {20, 25}, {20, 20}}
	polys := [][][]geojson.Position{
		{ring1},
		{ring2},
	}
	mp := geojson.NewMultiPolygon(polys)
	f := geojson.NewFeature(mp, nil)
	bbox := []float64{0, 0, 10, 10}

	_, err := BBoxClip(f, bbox)
	if err == nil {
		t.Fatal("expected error for all polygons outside bbox")
	}
}

func TestBBoxClip_WithGeometryDirect(t *testing.T) {
	// Pass a geometry directly (not wrapped in Feature)
	pt := geojson.NewPoint(geojson.Position{5, 5})
	bbox := []float64{0, 0, 10, 10}

	result, err := BBoxClip(pt, bbox)
	if err != nil {
		t.Fatal(err)
	}
	rp, ok := result.Geometry.(*geojson.Point)
	if !ok {
		t.Fatalf("expected *Point, got %T", result.Geometry)
	}
	expectPoint(t, rp.Coordinates, 5, 5)
}

func TestBBoxClip_InvalidBBox(t *testing.T) {
	f := geojson.NewFeature(geojson.NewPoint(geojson.Position{5, 5}), nil)

	// BBox with fewer than 4 elements
	_, err := BBoxClip(f, []float64{0, 0})
	if err == nil {
		t.Fatal("expected error for invalid bbox")
	}

	_, err = BBoxClip(f, []float64{})
	if err == nil {
		t.Fatal("expected error for empty bbox")
	}
}

func TestBBoxClip_UnsupportedType(t *testing.T) {
	gc := geojson.NewGeometryCollection([]geojson.Geometry{
		geojson.NewPoint(geojson.Position{5, 5}),
	})
	f := geojson.NewFeature(gc, nil)
	bbox := []float64{0, 0, 10, 10}

	_, err := BBoxClip(f, bbox)
	if err == nil {
		t.Fatal("expected error for unsupported geometry type")
	}
}

func TestBBoxClip_ErrorMessages(t *testing.T) {
	tests := []struct {
		name     string
		obj      any
		bbox     []float64
		wantErr  bool
		contains string
	}{
		{
			name:    "empty bbox",
			obj:     geojson.NewFeature(geojson.NewPoint(geojson.Position{1, 1}), nil),
			bbox:    []float64{},
			wantErr: true,
		},
		{
			name:    "short bbox",
			obj:     geojson.NewFeature(geojson.NewPoint(geojson.Position{1, 1}), nil),
			bbox:    []float64{1, 2, 3},
			wantErr: true,
		},
		{
			name:    "point outside",
			obj:     geojson.NewFeature(geojson.NewPoint(geojson.Position{-1, -1}), nil),
			bbox:    []float64{0, 0, 10, 10},
			wantErr: true,
		},
		{
			name:    "no feature geometry",
			obj:     geojson.NewFeature(nil, nil),
			bbox:    []float64{0, 0, 10, 10},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := BBoxClip(tc.obj, tc.bbox)
			if tc.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestBBoxClip_LineString_AllSegmentsDisconnected(t *testing.T) {
	// Line with two disconnected segments inside bbox → picks the longest
	line := []geojson.Position{
		{-5, 5},  // enters bbox from left
		{5, 5},   // inside
		{15, -5}, // exits right+bottom (trivial reject with next)
		{15, 15}, // outside right+top
		{5, 5},   // re-enters at right
		{-5, 5},  // exits left
	}
	ls := geojson.NewLineString(line)
	f := geojson.NewFeature(ls, nil)
	bbox := []float64{0, 0, 10, 10}

	result, err := BBoxClip(f, bbox)
	if err != nil {
		t.Fatal(err)
	}
	// Should return a LineString (picks the longest clipped segment)
	rls, ok := result.Geometry.(*geojson.LineString)
	if !ok {
		t.Fatalf("expected *LineString, got %T", result.Geometry)
	}
	if len(rls.Coordinates) < 2 {
		t.Fatalf("expected at least 2 points, got %d", len(rls.Coordinates))
	}
}

// ---------------------------------------------------------------------------
// Distance tests
// ---------------------------------------------------------------------------

func TestDistance_Zero(t *testing.T) {
	d := Distance(0, 0, 0, 0)
	if d != 0 {
		t.Errorf("expected 0, got %f", d)
	}

	d = Distance(10, 20, 10, 20)
	if d != 0 {
		t.Errorf("expected 0, got %f", d)
	}
}

func TestDistance_OneDegreeAlongEquator(t *testing.T) {
	// Distance from (0,0) to (0,1) along the equator
	// Should be approximately 1° in radians * earth radius
	expected := earthRadius * math.Pi / 180 // ≈ 111195 m
	d := Distance(0, 0, 0, 1)
	if mathAbs(d-expected) > 0.5 {
		t.Errorf("expected ~%f, got %f (diff %f)", expected, d, mathAbs(d-expected))
	}
}

func TestDistance_OneDegreeAlongMeridian(t *testing.T) {
	// Distance from (0,0) to (1,0) along the prime meridian (same as equator)
	expected := earthRadius * math.Pi / 180
	d := Distance(0, 0, 1, 0)
	if mathAbs(d-expected) > 0.5 {
		t.Errorf("expected ~%f, got %f (diff %f)", expected, d, mathAbs(d-expected))
	}
}

func TestDistance_HalfDegree(t *testing.T) {
	// Distance from (0,0) to (0,0.5) → half of 1° distance
	expected := earthRadius * math.Pi / 360
	d := Distance(0, 0, 0, 0.5)
	if mathAbs(d-expected) > 0.5 {
		t.Errorf("expected ~%f, got %f (diff %f)", expected, d, mathAbs(d-expected))
	}
}

func TestDistance_Symmetric(t *testing.T) {
	// Distance should be symmetric (commutative)
	d1 := Distance(10, 20, 30, 40)
	d2 := Distance(30, 40, 10, 20)
	if mathAbs(d1-d2) > 0.001 {
		t.Errorf("distance should be symmetric: %f vs %f", d1, d2)
	}
}

func TestDistance_LargeDistance(t *testing.T) {
	// Distance between two distant points
	// New York (lat 40.7128, lng -74.0060) to London (lat 51.5074, lng -0.1278)
	// Approximate distance: ~5570 km
	d := Distance(-74.0060, 40.7128, -0.1278, 51.5074)
	if d < 5000000 || d > 6000000 {
		t.Errorf("expected distance ~5570km, got %f km", d/1000)
	}
}

func TestDistance_SameLongitude(t *testing.T) {
	// Points with same longitude → pure latitudinal distance
	// 1° latitude ≈ 111195 m
	d := Distance(0, 30, 0, 30.5)
	expected := earthRadius * (0.5 * math.Pi / 180) // half degree
	if mathAbs(d-expected) > 0.5 {
		t.Errorf("expected ~%f, got %f (diff %f)", expected, d, mathAbs(d-expected))
	}
}

func TestDistance_SameLatitude(t *testing.T) {
	// Points with same latitude → pure longitudinal distance
	// At equator, 1° longitude ≈ 111195 m
	d := Distance(0, 0, 1, 0)
	expected := earthRadius * math.Pi / 180
	if mathAbs(d-expected) > 0.5 {
		t.Errorf("expected ~%f, got %f (diff %f)", expected, d, mathAbs(d-expected))
	}
}

// ---------------------------------------------------------------------------
// BBoxClip table-driven tests
// ---------------------------------------------------------------------------

func TestBBoxClip_TableDriven(t *testing.T) {
	bbox := []float64{0, 0, 10, 10}

	tests := []struct {
		name       string
		obj        any
		wantErr    bool
		wantType   string // geojson type name
		checkCoords func(t *testing.T, geom geojson.Geometry)
	}{
		{
			name:     "point inside",
			obj:      geojson.NewFeature(geojson.NewPoint(geojson.Position{5, 5}), nil),
			wantType: geojson.TypePoint,
			checkCoords: func(t *testing.T, geom geojson.Geometry) {
				p := geom.(*geojson.Point)
				expectPoint(t, p.Coordinates, 5, 5)
			},
		},
		{
			name:     "point on left edge",
			obj:      geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 5}), nil),
			wantType: geojson.TypePoint,
			checkCoords: func(t *testing.T, geom geojson.Geometry) {
				p := geom.(*geojson.Point)
				expectPoint(t, p.Coordinates, 0, 5)
			},
		},
		{
			name:     "point on top-right corner",
			obj:      geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 10}), nil),
			wantType: geojson.TypePoint,
			checkCoords: func(t *testing.T, geom geojson.Geometry) {
				p := geom.(*geojson.Point)
				expectPoint(t, p.Coordinates, 10, 10)
			},
		},
		{
			name:    "point outside",
			obj:     geojson.NewFeature(geojson.NewPoint(geojson.Position{-5, 5}), nil),
			wantErr: true,
		},
		{
			name:     "multipoint partial",
			obj:      geojson.NewFeature(geojson.NewMultiPoint([]geojson.Position{{5, 5}, {-5, -5}}), nil),
			wantType: geojson.TypeMultiPoint,
			checkCoords: func(t *testing.T, geom geojson.Geometry) {
				mp := geom.(*geojson.MultiPoint)
				if len(mp.Coordinates) != 1 {
					t.Fatalf("expected 1 point, got %d", len(mp.Coordinates))
				}
			},
		},
		{
			name:    "multipoint all outside",
			obj:     geojson.NewFeature(geojson.NewMultiPoint([]geojson.Position{{-5, -5}, {-2, -2}}), nil),
			wantErr: true,
		},
		{
			name: "linestring inside",
			obj: geojson.NewFeature(geojson.NewLineString([]geojson.Position{
				{-5, 5}, {5, 5}, {15, 5}, {15, 15}, {5, 5}, {-5, 5},
			}), nil),
			wantType: geojson.TypeLineString,
		},
		{
			name: "linestring clipped across vertical",
			obj: geojson.NewFeature(geojson.NewLineString([]geojson.Position{
				{-5, 5}, {5, 5}, {15, 5}, {15, 15}, {5, 5},
			}), nil),
			wantType: geojson.TypeLineString,
			checkCoords: func(t *testing.T, geom geojson.Geometry) {
				ls := geom.(*geojson.LineString)
				// Longest segment: enters at (0,5), crosses to (10,5)
				if len(ls.Coordinates) < 2 {
					t.Fatalf("expected >= 2 points, got %d", len(ls.Coordinates))
				}
				expectPoint(t, ls.Coordinates[0], 0, 5)
			},
		},
		{
			name:     "linestring clipped across horizontal",
			obj: geojson.NewFeature(geojson.NewLineString([]geojson.Position{
				{5, -5}, {5, 5}, {5, 15}, {15, 15}, {5, 5},
			}), nil),
			wantType: geojson.TypeLineString,
			checkCoords: func(t *testing.T, geom geojson.Geometry) {
				ls := geom.(*geojson.LineString)
				if len(ls.Coordinates) < 2 {
					t.Fatalf("expected >= 2 points, got %d", len(ls.Coordinates))
				}
				expectPoint(t, ls.Coordinates[0], 5, 0)
			},
		},
		{
			name:    "linestring all outside",
			obj:     geojson.NewFeature(geojson.NewLineString([]geojson.Position{{-5, -5}, {-2, -2}}), nil),
			wantErr: true,
		},
		{
			name:     "polygon inside",
			obj:      geojson.NewFeature(geojson.NewPolygon([][]geojson.Position{{{2, 2}, {8, 2}, {8, 8}, {2, 8}, {2, 2}}}), nil),
			wantType: geojson.TypePolygon,
		},
		{
			name:    "polygon outside",
			obj:     geojson.NewFeature(geojson.NewPolygon([][]geojson.Position{{{20, 20}, {25, 20}, {25, 25}, {20, 25}, {20, 20}}}), nil),
			wantErr: true,
		},
		{
			name:    "unsupported geometry",
			obj:     geojson.NewFeature(geojson.NewGeometryCollection([]geojson.Geometry{geojson.NewPoint(geojson.Position{5, 5})}), nil),
			wantErr: true,
		},
		{
			name:    "invalid bbox",
			obj:     geojson.NewFeature(geojson.NewPoint(geojson.Position{5, 5}), nil),
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tb := bbox
			if tc.name == "invalid bbox" {
				tb = []float64{0, 0}
			}
			result, err := BBoxClip(tc.obj, tb)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.wantType != "" && result.Geometry.Type() != tc.wantType {
				t.Errorf("expected type %s, got %s", tc.wantType, result.Geometry.Type())
			}
			if tc.checkCoords != nil {
				tc.checkCoords(t, result.Geometry)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// expectPoint checks that a geojson.Position has the expected x, y values
// within a small tolerance.
func expectPoint(t *testing.T, p geojson.Position, x, y float64) {
	t.Helper()
	const tol = 1e-10
	if mathAbs(p[0]-x) > tol || mathAbs(p[1]-y) > tol {
		t.Errorf("expected (%.2f, %.2f), got (%.2f, %.2f)", x, y, p[0], p[1])
	}
}
