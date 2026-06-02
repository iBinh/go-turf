package lines

import (
	"testing"

	"github.com/ibinh/turf-go/geojson"
	"github.com/ibinh/turf-go/measurement"
)

func pt(lng, lat float64) geojson.Position {
	return geojson.Position{lng, lat}
}

func TestLineIntersectCrossing(t *testing.T) {
	l1 := geojson.NewLineString([]geojson.Position{{0, 0}, {10, 10}})
	l2 := geojson.NewLineString([]geojson.Position{{0, 10}, {10, 0}})
	result, err := LineIntersect(l1, l2)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 1 {
		t.Fatalf("expected 1 intersection, got %d", len(result.Features))
	}
	coord, _ := geojson.GetCoord(result.Features[0])
	if coord[0] != 5 || coord[1] != 5 {
		t.Errorf("expected (5,5), got (%v,%v)", coord[0], coord[1])
	}
}

func TestLineIntersectParallel(t *testing.T) {
	l1 := geojson.NewLineString([]geojson.Position{{0, 0}, {10, 0}})
	l2 := geojson.NewLineString([]geojson.Position{{0, 5}, {10, 5}})
	result, err := LineIntersect(l1, l2)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 0 {
		t.Errorf("expected 0 intersections, got %d", len(result.Features))
	}
}

func TestLineIntersectNoIntersection(t *testing.T) {
	l1 := geojson.NewLineString([]geojson.Position{{0, 0}, {5, 5}})
	l2 := geojson.NewLineString([]geojson.Position{{6, 6}, {10, 10}})
	result, err := LineIntersect(l1, l2)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 0 {
		t.Errorf("expected 0 intersections, got %d", len(result.Features))
	}
}

func TestLineIntersectFeatureInput(t *testing.T) {
	l1 := geojson.NewFeature(geojson.NewLineString([]geojson.Position{{0, 0}, {10, 10}}), nil)
	l2 := geojson.NewFeature(geojson.NewLineString([]geojson.Position{{0, 10}, {10, 0}}), nil)
	result, err := LineIntersect(l1, l2)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 1 {
		t.Errorf("expected 1 intersection with feature input, got %d", len(result.Features))
	}
}

func TestLineSegment(t *testing.T) {
	line := geojson.NewLineString([]geojson.Position{{0, 0}, {5, 5}, {10, 0}})
	result, err := LineSegment(line)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 2 {
		t.Fatalf("expected 2 segments, got %d", len(result.Features))
	}
	for i, seg := range result.Features {
		ls, ok := seg.Geometry.(*geojson.LineString)
		if !ok {
			t.Errorf("segment %d: expected LineString", i)
			continue
		}
		if len(ls.Coordinates) != 2 {
			t.Errorf("segment %d: expected 2 coordinates, got %d", i, len(ls.Coordinates))
		}
	}
}

func TestLineSegmentPolygon(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}},
	})
	result, err := LineSegment(poly)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 4 {
		t.Errorf("expected 4 segments for a quad polygon, got %d", len(result.Features))
	}
}

func TestLineSegmentMultiLineString(t *testing.T) {
	ml := geojson.NewMultiLineString([][]geojson.Position{
		{{0, 0}, {1, 1}},
		{{2, 2}, {3, 3}, {4, 4}},
	})
	result, err := LineSegment(ml)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 3 {
		t.Errorf("expected 3 segments, got %d", len(result.Features))
	}
}

func TestLineOverlapFullOverlap(t *testing.T) {
	l1 := geojson.NewLineString([]geojson.Position{{0, 0}, {10, 0}})
	l2 := geojson.NewLineString([]geojson.Position{{0, 0}, {10, 0}})
	result, err := LineOverlap(l1, l2)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) == 0 {
		t.Error("expected at least 1 overlapping segment")
	}
}

func TestLineOverlapPartial(t *testing.T) {
	l1 := geojson.NewLineString([]geojson.Position{{0, 0}, {10, 0}})
	l2 := geojson.NewLineString([]geojson.Position{{3, 0}, {7, 0}})
	result, err := LineOverlap(l1, l2)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) == 0 {
		t.Error("expected at least 1 overlapping segment")
	}
}

func TestLineOverlapNoOverlap(t *testing.T) {
	l1 := geojson.NewLineString([]geojson.Position{{0, 0}, {5, 0}})
	l2 := geojson.NewLineString([]geojson.Position{{6, 0}, {10, 0}})
	result, err := LineOverlap(l1, l2)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 0 {
		t.Errorf("expected 0 overlapping segments, got %d", len(result.Features))
	}
}

func TestLineSlice(t *testing.T) {
	line := geojson.NewLineString([]geojson.Position{{0, 0}, {10, 0}})
	p1 := geojson.NewPoint(pt(0, 0))
	p2 := geojson.NewPoint(pt(10, 0))
	result, err := LineSlice(p1, p2, line)
	if err != nil {
		t.Fatal(err)
	}
	coords, _ := geojson.GetCoords(result)
	pts := coords.([]geojson.Position)
	if len(pts) < 2 {
		t.Fatal("expected at least 2 coordinates")
	}
}

func TestLineSliceAlong(t *testing.T) {
	line := geojson.NewLineString([]geojson.Position{{0, 0}, {10, 0}})
	startDist := 0.0
	endDist := 50000.0
	result, err := LineSliceAlong(line, startDist, endDist, measurement.UnitMeters)
	if err != nil {
		t.Fatal(err)
	}
	coords, _ := geojson.GetCoords(result)
	pts := coords.([]geojson.Position)
	if len(pts) < 2 {
		t.Fatal("expected at least 2 coordinates")
	}
}

func TestLineChunk(t *testing.T) {
	line := geojson.NewLineString([]geojson.Position{{0, 0}, {10, 0}})
	result, err := LineChunk(line, 50000, measurement.UnitMeters)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) == 0 {
		t.Error("expected at least 1 chunk")
	}
}

func TestLineChunkExactDivisions(t *testing.T) {
	line := geojson.NewLineString([]geojson.Position{{0, 0}, {0, 1}})
	lenKm, _ := measurement.Length(line, measurement.UnitKilometers)
	segLen := lenKm / 4
	result, err := LineChunk(line, segLen, measurement.UnitKilometers)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) < 3 || len(result.Features) > 5 {
		t.Errorf("expected ~4 chunks, got %d", len(result.Features))
	}
}

func TestLineSplit(t *testing.T) {
	line := geojson.NewLineString([]geojson.Position{{0, 0}, {10, 0}, {20, 0}})
	point := geojson.NewPoint(pt(10, 0))
	result, err := LineSplit(line, point)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 2 {
		t.Fatalf("expected 2 features, got %d", len(result.Features))
	}
	coords1, _ := geojson.GetCoords(result.Features[0])
	pts1 := coords1.([]geojson.Position)
	if pts1[len(pts1)-1][0] != 10 || pts1[len(pts1)-1][1] != 0 {
		t.Error("split point not at end of first line")
	}
}

func TestLineSplitAtVertex(t *testing.T) {
	line := geojson.NewLineString([]geojson.Position{{0, 0}, {5, 0}, {10, 0}})
	point := geojson.NewPoint(pt(5, 0))
	result, err := LineSplit(line, point)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 2 {
		t.Fatalf("expected 2 features, got %d", len(result.Features))
	}
	coords1, _ := geojson.GetCoords(result.Features[0])
	pts1 := coords1.([]geojson.Position)
	if len(pts1) != 2 {
		t.Errorf("first part should have 2 coords, got %d", len(pts1))
	}
}

func TestLineArc(t *testing.T) {
	center := geojson.NewPoint(pt(0, 0))
	result, err := LineArc(center, 100000, 0, 90, 10, measurement.UnitMeters)
	if err != nil {
		t.Fatal(err)
	}
	coords, _ := geojson.GetCoords(result)
	pts := coords.([]geojson.Position)
	if len(pts) != 11 {
		t.Errorf("expected 11 coords for 10-step arc, got %d", len(pts))
	}
}

func TestLineArcWrapAround(t *testing.T) {
	center := geojson.NewPoint(pt(0, 0))
	result, err := LineArc(center, 100000, 270, 90, 10, measurement.UnitMeters)
	if err != nil {
		t.Fatal(err)
	}
	coords, _ := geojson.GetCoords(result)
	pts := coords.([]geojson.Position)
	if len(pts) != 11 {
		t.Errorf("expected 11 coords for wrap-around arc, got %d", len(pts))
	}
}

func TestSector(t *testing.T) {
	center := geojson.NewPoint(pt(0, 0))
	result, err := Sector(center, 100000, 0, 90, 10, measurement.UnitMeters)
	if err != nil {
		t.Fatal(err)
	}
	poly, ok := result.Geometry.(*geojson.Polygon)
	if !ok {
		t.Fatal("expected Polygon")
	}
	ring := poly.Coordinates[0]
	if len(ring) < 3 {
		t.Fatalf("expected at least 3 points in ring, got %d", len(ring))
	}
	first := ring[0]
	last := ring[len(ring)-1]
	if first[0] != last[0] || first[1] != last[1] {
		t.Error("sector polygon should be closed (first==last)")
	}
	if first[0] != 0 || first[1] != 0 {
		t.Error("first point should be the center")
	}
}

func TestSectorWrapAround(t *testing.T) {
	center := geojson.NewPoint(pt(0, 0))
	result, err := Sector(center, 100000, 270, 90, 10, measurement.UnitMeters)
	if err != nil {
		t.Fatal(err)
	}
	poly, ok := result.Geometry.(*geojson.Polygon)
	if !ok {
		t.Fatal("expected Polygon")
	}
	ring := poly.Coordinates[0]
	if len(ring) < 3 {
		t.Fatalf("expected at least 3 points, got %d", len(ring))
	}
}

func TestLineSplitAtEndpoint(t *testing.T) {
	line := geojson.NewLineString([]geojson.Position{{0, 0}, {10, 0}})
	point := geojson.NewPoint(pt(0, 0))
	result, err := LineSplit(line, point)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 1 {
		t.Errorf("expected 1 feature when splitting at endpoint, got %d", len(result.Features))
	}
}
