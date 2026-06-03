package transform

import (
	"math"
	"testing"

	"github.com/ibinh/turf-go/geojson"
)

func TestTransformRotate(t *testing.T) {
	pt := geojson.NewFeature(geojson.NewPoint(geojson.Position{1, 0}), nil)
	rotated, err := TransformRotate(pt, 90, geojson.Position{0, 0})
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(rotated)
	if math.Abs(coord[0]-0) > 0.001 || math.Abs(coord[1]-1) > 0.001 {
		t.Errorf("expected ~[0,1] after 90 deg, got %v", coord)
	}
}

func TestTransformRotateNoPivot(t *testing.T) {
	pt := geojson.NewFeature(geojson.NewPoint(geojson.Position{1, 0}), nil)
	rotated, err := TransformRotate(pt, 180)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(rotated)
	if math.Abs(coord[0]-1) > 0.001 || math.Abs(coord[1]-0) > 0.001 {
		t.Errorf("point should stay same when rotating around itself, got %v", coord)
	}
}

func TestTransformRotateLine(t *testing.T) {
	line := geojson.NewFeature(
		geojson.NewLineString([]geojson.Position{{0, 0}, {2, 0}}),
		nil,
	)
	rotated, err := TransformRotate(line, 90, geojson.Position{0, 0})
	if err != nil {
		t.Fatal(err)
	}
	coords, _ := geojson.GetCoords(rotated)
	pts := coords.([]geojson.Position)
	if math.Abs(pts[0][0]) > 0.001 || math.Abs(pts[0][1]) > 0.001 {
		t.Errorf("first point should stay at origin, got %v", pts[0])
	}
	if math.Abs(pts[1][0]) > 0.001 || math.Abs(pts[1][1]-2) > 0.001 {
		t.Errorf("expected ~[0,2], got %v", pts[1])
	}
}

func TestTransformScale(t *testing.T) {
	pt := geojson.NewFeature(geojson.NewPoint(geojson.Position{2, 3}), nil)
	scaled, err := TransformScale(pt, 2, geojson.Position{0, 0})
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(scaled)
	if math.Abs(coord[0]-4) > 0.001 || math.Abs(coord[1]-6) > 0.001 {
		t.Errorf("expected [4,6], got %v", coord)
	}
}

func TestTransformScaleXY(t *testing.T) {
	pt := geojson.NewFeature(geojson.NewPoint(geojson.Position{2, 3}), nil)
	scaled, err := TransformScaleXY(pt, 2, 3, geojson.Position{0, 0})
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(scaled)
	if math.Abs(coord[0]-4) > 0.001 || math.Abs(coord[1]-9) > 0.001 {
		t.Errorf("expected [4,9], got %v", coord)
	}
}

func TestTransformTranslate(t *testing.T) {
	pt := geojson.NewFeature(geojson.NewPoint(geojson.Position{1, 2}), nil)
	translated, err := TransformTranslate(pt, 10, 20)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(translated)
	if math.Abs(coord[0]-11) > 0.001 || math.Abs(coord[1]-22) > 0.001 {
		t.Errorf("expected [11,22], got %v", coord)
	}
}

func TestTranslateMeters(t *testing.T) {
	pt := geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil)
	translated, err := TransformTranslate(pt, 111319.9, 0, &TranslateOptions{Units: "meters"})
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(translated)
	if math.Abs(coord[0]-1) > 0.01 {
		t.Errorf("expected ~1 degree east, got %f", coord[0])
	}
}

func TestFlip(t *testing.T) {
	pt := geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 20}), nil)
	flipped, err := Flip(pt)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(flipped)
	if coord[0] != 20 || coord[1] != 10 {
		t.Errorf("expected [20,10], got %v", coord)
	}
}

func TestTruncate(t *testing.T) {
	pt := geojson.NewFeature(geojson.NewPoint(geojson.Position{1.12345, 2.67890}), nil)
	truncated, err := Truncate(pt, 2)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(truncated)
	if math.Abs(coord[0]-1.12) > 0.001 || math.Abs(coord[1]-2.68) > 0.001 {
		t.Errorf("expected [1.12, 2.68], got %v", coord)
	}
}

func TestCleanCoords(t *testing.T) {
	line := geojson.NewFeature(
		geojson.NewLineString([]geojson.Position{{0, 0}, {0, 0}, {1, 1}, {1, 1}, {2, 2}}),
		nil,
	)
	cleaned, err := CleanCoords(line)
	if err != nil {
		t.Fatal(err)
	}
	coords, _ := geojson.GetCoords(cleaned)
	pts := coords.([]geojson.Position)
	if len(pts) != 3 {
		t.Errorf("expected 3 unique points, got %d", len(pts))
	}
}

func TestRewind(t *testing.T) {
	ring := []geojson.Position{{0, 0}, {1, 1}, {2, 0}, {0, 0}}
	poly := geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{ring}),
		nil,
	)
	rewound, err := Rewind(poly)
	if err != nil {
		t.Fatal(err)
	}
	p, err := geojson.GetGeometry(rewound)
	if err != nil {
		t.Fatal(err)
	}
	// After rewind, exterior should be CW
	pg, ok := p.(*geojson.Polygon)
	if !ok {
		t.Fatalf("expected Polygon")
	}
	if !isCWRing(pg.Coordinates[0]) {
		t.Error("exterior should be CW after default rewind")
	}
}

func TestRewindReverse(t *testing.T) {
	ring := []geojson.Position{{0, 0}, {2, 0}, {1, 1}, {0, 0}}
	poly := geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{ring}),
		nil,
	)
	rewound, err := Rewind(poly, true)
	if err != nil {
		t.Fatal(err)
	}
	p, err := geojson.GetGeometry(rewound)
	if err != nil {
		t.Fatal(err)
	}
	pg, ok := p.(*geojson.Polygon)
	if !ok {
		t.Fatalf("expected Polygon")
	}
	if isCWRing(pg.Coordinates[0]) {
		t.Error("exterior should be CCW with reverse=true")
	}
}

func TestToMercator(t *testing.T) {
	pt := geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil)
	merc, err := ToMercator(pt)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(merc)
	if math.Abs(coord[0]) > 0.001 || math.Abs(coord[1]) > 0.001 {
		t.Errorf("expected [0,0] at equator/prime, got %v", coord)
	}
}

func TestToWGS84(t *testing.T) {
	pt := geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil)
	wgs, err := ToWGS84(pt)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(wgs)
	if math.Abs(coord[0]) > 0.001 || math.Abs(coord[1]) > 0.001 {
		t.Errorf("expected [0,0], got %v", coord)
	}
}

func TestMercatorRoundtrip(t *testing.T) {
	original := geojson.NewFeature(geojson.NewPoint(geojson.Position{-73.935242, 40.73061}), nil)
	merc, err := ToMercator(original)
	if err != nil {
		t.Fatal(err)
	}
	wgs, err := ToWGS84(merc)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(wgs)
	if math.Abs(coord[0]+73.935242) > 0.01 || math.Abs(coord[1]-40.73061) > 0.01 {
		t.Errorf("roundtrip failed: expected ~[-73.94, 40.73], got %v", coord)
	}
}

func TestTransformMultiPoint(t *testing.T) {
	mp := geojson.NewFeature(geojson.NewMultiPoint([]geojson.Position{{1, 0}, {2, 0}}), nil)
	rotated, err := TransformRotate(mp, 90, geojson.Position{0, 0})
	if err != nil {
		t.Fatal(err)
	}
	coords, _ := geojson.GetCoords(rotated)
	pts := coords.([]geojson.Position)
	if math.Abs(pts[0][0]) > 0.001 || math.Abs(pts[0][1]-1) > 0.001 {
		t.Errorf("expected [0,1], got %v", pts[0])
	}
}

func TestTransformMultiLineString(t *testing.T) {
	mls := geojson.NewFeature(geojson.NewMultiLineString([][]geojson.Position{{{0, 0}, {1, 0}}}), nil)
	scaled, err := TransformScale(mls, 2, geojson.Position{0, 0})
	if err != nil {
		t.Fatal(err)
	}
	coords, _ := geojson.GetCoords(scaled)
	lines := coords.([][]geojson.Position)
	if math.Abs(lines[0][1][0]-2) > 0.001 {
		t.Errorf("expected x=2, got %f", lines[0][1][0])
	}
}

func TestTransformMultiPolygon(t *testing.T) {
	mp := geojson.NewFeature(geojson.NewMultiPolygon([][][]geojson.Position{{{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}}}}), nil)
	translated, err := TransformTranslate(mp, 10, 10)
	if err != nil {
		t.Fatal(err)
	}
	coords, _ := geojson.GetCoords(translated)
	polygons := coords.([][][]geojson.Position)
	if math.Abs(polygons[0][0][0][0]-10) > 0.001 {
		t.Errorf("expected x=10, got %f", polygons[0][0][0][0])
	}
}

func TestTransformGeometryCollection(t *testing.T) {
	gc := geojson.NewFeature(geojson.NewGeometryCollection([]geojson.Geometry{
		geojson.NewPoint(geojson.Position{1, 0}),
	}), nil)
	rotated, err := TransformRotate(gc, 90, geojson.Position{0, 0})
	if err != nil {
		t.Fatal(err)
	}
	g, err := geojson.GetGeometry(rotated)
	if err != nil {
		t.Fatal(err)
	}
	gc2, ok := g.(*geojson.GeometryCollection)
	if !ok {
		t.Fatalf("expected GeometryCollection, got %T", g)
	}
	pt := gc2.Geometries[0].(*geojson.Point)
	if math.Abs(pt.Coordinates[0]) > 0.001 || math.Abs(pt.Coordinates[1]-1) > 0.001 {
		t.Errorf("expected [0,1], got %v", pt.Coordinates)
	}
}

func TestTransformScaleNoPivot(t *testing.T) {
	poly := geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}}}),
		nil,
	)
	scaled, err := TransformScale(poly, 2)
	if err != nil {
		t.Fatal(err)
	}
	coords, _ := geojson.GetCoords(scaled)
	rings := coords.([][]geojson.Position)
	if math.Abs(rings[0][2][0]-16) > 0.001 {
		t.Errorf("expected x~16, got %f", rings[0][2][0])
	}
}

func TestTransformScaleXYNoPivot(t *testing.T) {
	poly := geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}}}),
		nil,
	)
	scaled, err := TransformScaleXY(poly, 2, 3)
	if err != nil {
		t.Fatal(err)
	}
	coords, _ := geojson.GetCoords(scaled)
	rings := coords.([][]geojson.Position)
	if math.Abs(rings[0][2][0]-16) > 0.001 || math.Abs(rings[0][2][1]-22) > 0.001 {
		t.Errorf("expected ~[16,22], got %v", rings[0][2])
	}
}

func TestTranslateMetersAtPole(t *testing.T) {
	pt := geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 89}), nil)
	translated, err := TransformTranslate(pt, 111319.9, 0, &TranslateOptions{Units: "meters"})
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(translated)
	if math.Abs(coord[1]-89) > 0.01 {
		t.Errorf("expected lat~89, got %f", coord[1])
	}
}

func TestCleanCoordsPolygon(t *testing.T) {
	poly := geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{{{0, 0}, {0, 0}, {1, 1}, {1, 1}, {0, 0}}}),
		nil,
	)
	cleaned, err := CleanCoords(poly)
	if err != nil {
		t.Fatal(err)
	}
	coords, _ := geojson.GetCoords(cleaned)
	rings := coords.([][]geojson.Position)
	if len(rings[0]) != 3 {
		t.Errorf("expected 3 unique points, got %d", len(rings[0]))
	}
}

func TestCleanCoordsMultiPolygon(t *testing.T) {
	mp := geojson.NewFeature(
		geojson.NewMultiPolygon([][][]geojson.Position{{{{0, 0}, {0, 0}, {1, 1}, {1, 1}, {0, 0}}}}),
		nil,
	)
	cleaned, err := CleanCoords(mp)
	if err != nil {
		t.Fatal(err)
	}
	_, err = geojson.GetCoords(cleaned)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCleanCoordsPoint(t *testing.T) {
	pt := geojson.NewFeature(geojson.NewPoint(geojson.Position{1, 2}), nil)
	cleaned, err := CleanCoords(pt)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(cleaned)
	if coord[0] != 1 || coord[1] != 2 {
		t.Errorf("expected [1,2], got %v", coord)
	}
}

func TestCleanCoordsMultiPoint(t *testing.T) {
	mp := geojson.NewFeature(
		geojson.NewMultiPoint([]geojson.Position{{1, 1}, {1, 1}, {2, 2}}),
		nil,
	)
	cleaned, err := CleanCoords(mp)
	if err != nil {
		t.Fatal(err)
	}
	coords, _ := geojson.GetCoords(cleaned)
	pts := coords.([]geojson.Position)
	if len(pts) != 2 {
		t.Errorf("expected 2 unique points, got %d", len(pts))
	}
}

func TestCleanCoordsMultiLineString(t *testing.T) {
	mls := geojson.NewFeature(
		geojson.NewMultiLineString([][]geojson.Position{{{0, 0}, {0, 0}, {1, 1}}}),
		nil,
	)
	cleaned, err := CleanCoords(mls)
	if err != nil {
		t.Fatal(err)
	}
	_, err = geojson.GetCoords(cleaned)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRewindMultiPolygon(t *testing.T) {
	mp := geojson.NewFeature(
		geojson.NewMultiPolygon([][][]geojson.Position{{{{0, 0}, {2, 0}, {1, 1}, {0, 0}}}}),
		nil,
	)
	rewound, err := Rewind(mp)
	if err != nil {
		t.Fatal(err)
	}
	g, err := geojson.GetGeometry(rewound)
	if err != nil {
		t.Fatal(err)
	}
	mpoly, ok := g.(*geojson.MultiPolygon)
	if !ok {
		t.Fatalf("expected MultiPolygon")
	}
	if !isCWRing(mpoly.Coordinates[0][0]) {
		t.Error("exterior should be CW after rewind")
	}
}

func TestTransformPolygon(t *testing.T) {
	ring := []geojson.Position{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}}
	poly := geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{ring}),
		nil,
	)
	translated, err := TransformTranslate(poly, 5, 5)
	if err != nil {
		t.Fatal(err)
	}
	coords, _ := geojson.GetCoords(translated)
	rings := coords.([][]geojson.Position)
	first := rings[0][0]
	if math.Abs(first[0]-5) > 0.001 || math.Abs(first[1]-5) > 0.001 {
		t.Errorf("first vertex should be [5,5], got %v", first)
	}
}
