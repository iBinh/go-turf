package turf

import (
	"math"
	"testing"

	"github.com/ibinh/turf-go/geojson"
)

const eps = 1e-6

func TestIntegrationDistance(t *testing.T) {
	from := geojson.NewPoint([]float64{-75.343, 39.984})
	to := geojson.NewPoint([]float64{-75.534, 39.123})
	d, err := Distance(from, to, UnitKilometers)
	if err != nil {
		t.Fatal(err)
	}
	expected := 97.1
	if math.Abs(d-expected) > 0.2 {
		t.Errorf("expected ~%.1f km, got %.6f", expected, d)
	}
}

func TestIntegrationArea(t *testing.T) {
	pts := [][]geojson.Position{
		{{-180, -90}, {180, -90}, {180, 90}, {-180, 90}, {-180, -90}},
	}
	poly := geojson.NewPolygon(pts)
	a, err := Area(poly)
	if err != nil {
		t.Fatal(err)
	}
	if a <= 0 {
		t.Error("area should be positive")
	}
	areaFraction := a / (4 * math.Pi * 6371008.8 * 6371008.8)
	if math.Abs(areaFraction-1) > 0.01 {
		t.Errorf("world polygon should cover ~half of sphere, got %.4f", areaFraction)
	}
}

func TestIntegrationBearing(t *testing.T) {
	from := geojson.NewPoint([]float64{0, 0})
	to := geojson.NewPoint([]float64{0, 10})
	b, err := Bearing(from, to)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(b) > eps {
		t.Errorf("expected 0 bearing north, got %.6f", b)
	}
}

func TestIntegrationDestination(t *testing.T) {
	origin := geojson.NewPoint([]float64{0, 0})
	dest, err := Destination(origin, 111.32, 90)
	if err != nil {
		t.Fatal(err)
	}
	coord, err := geojson.GetCoord(dest)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(coord[0]-1) > 0.01 || math.Abs(coord[1]) > 0.01 {
		t.Errorf("expected ~(1, 0), got (%.4f, %.4f)", coord[0], coord[1])
	}
}

func TestIntegrationMidpoint(t *testing.T) {
	from := geojson.NewPoint([]float64{0, 0})
	to := geojson.NewPoint([]float64{10, 0})
	mid, err := Midpoint(from, to)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(mid)
	if math.Abs(coord[0]-5) > 0.01 || math.Abs(coord[1]) > 0.01 {
		t.Errorf("expected ~(5, 0), got (%.4f, %.4f)", coord[0], coord[1])
	}
}

func TestIntegrationLength(t *testing.T) {
	ls := geojson.NewLineString([]geojson.Position{
		{-77.031669, 38.878605},
		{-77.029609, 38.881946},
	})
	l, err := Length(ls, UnitKilometers)
	if err != nil {
		t.Fatal(err)
	}
	if l <= 0 {
		t.Error("length should be positive")
	}
}

func TestIntegrationBBox(t *testing.T) {
	ls := geojson.NewLineString([]geojson.Position{
		{-10, -10}, {10, 10},
	})
	bb, err := BBox(ls)
	if err != nil {
		t.Fatal(err)
	}
	if len(bb) != 4 {
		t.Fatal("bbox should have 4 values")
	}
	if math.Abs(bb[0]+10) > eps || math.Abs(bb[1]+10) > eps {
		t.Error("bbox min values incorrect")
	}
	if math.Abs(bb[2]-10) > eps || math.Abs(bb[3]-10) > eps {
		t.Error("bbox max values incorrect")
	}
}

func TestIntegrationCenter(t *testing.T) {
	ls := geojson.NewLineString([]geojson.Position{
		{-10, -10}, {10, 10},
	})
	c, err := Center(ls)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(c)
	if math.Abs(coord[0]) > eps || math.Abs(coord[1]) > eps {
		t.Errorf("expected (0,0), got (%.4f, %.4f)", coord[0], coord[1])
	}
}

func TestIntegrationCentroid(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}},
	})
	c, err := Centroid(poly)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(c)
	if math.Abs(coord[0]-4) > 0.01 || math.Abs(coord[1]-4) > 0.01 {
		t.Errorf("expected (4,4), got (%.4f, %.4f)", coord[0], coord[1])
	}
}

func TestIntegrationCenterFunc(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}},
	})
	c, err := Center(poly)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(c)
	if math.Abs(coord[0]-5) > 0.01 || math.Abs(coord[1]-5) > 0.01 {
		t.Errorf("expected (5,5), got (%.4f, %.4f)", coord[0], coord[1])
	}
}

func TestIntegrationPointGrid(t *testing.T) {
	grid, err := PointGrid([]float64{0, 0, 10, 10}, 2, UnitMiles)
	if err != nil {
		t.Fatal(err)
	}
	if len(grid.Features) < 2 {
		t.Error("point grid should produce points")
	}
}

func TestIntegrationHexGrid(t *testing.T) {
	grid, err := HexGrid([]float64{0, 0, 5, 5}, 2, UnitMiles)
	if err != nil {
		t.Fatal(err)
	}
	if len(grid.Features) < 1 {
		t.Error("hex grid should produce features")
	}
}

func TestIntegrationCircle(t *testing.T) {
	center := geojson.NewPoint([]float64{0, 0})
	circle, err := Circle(center, 100, CircleOptions{Steps: 64})
	if err != nil {
		t.Fatal(err)
	}
	poly, ok := circle.Geometry.(*geojson.Polygon)
	if !ok {
		t.Fatal("expected Polygon")
	}
	if len(poly.Coordinates[0]) < 4 {
		t.Error("circle should have multiple vertices")
	}
}

func TestIntegrationFlip(t *testing.T) {
	pt := geojson.NewPoint([]float64{10, 20})
	flipped, err := Flip(pt)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(flipped)
	if math.Abs(coord[0]-20) > eps || math.Abs(coord[1]-10) > eps {
		t.Errorf("expected (20,10), got (%.4f, %.4f)", coord[0], coord[1])
	}
}

func TestIntegrationPointInPolygon(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}},
	})
	pt := geojson.NewPoint([]float64{5, 5})
	inside, err := PointInPolygon(pt, poly)
	if err != nil {
		t.Fatal(err)
	}
	if !inside {
		t.Error("point should be inside polygon")
	}

	outsidePt := geojson.NewPoint([]float64{20, 20})
	outside, err := PointInPolygon(outsidePt, poly)
	if err != nil {
		t.Fatal(err)
	}
	if outside {
		t.Error("point should be outside polygon")
	}
}

func TestIntegrationIntersects(t *testing.T) {
	poly1 := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}},
	})
	poly2 := geojson.NewPolygon([][]geojson.Position{
		{{5, 5}, {15, 5}, {15, 15}, {5, 15}, {5, 5}},
	})
	intersects, err := Intersects(poly1, poly2)
	if err != nil {
		t.Fatal(err)
	}
	if !intersects {
		t.Error("overlapping polygons should intersect")
	}
}

func TestIntegrationContains(t *testing.T) {
	outer := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}},
	})
	inner := geojson.NewPolygon([][]geojson.Position{
		{{1, 1}, {9, 1}, {9, 9}, {1, 9}, {1, 1}},
	})
	contains, err := Contains(outer, inner)
	if err != nil {
		t.Fatal(err)
	}
	if !contains {
		t.Error("outer should contain inner")
	}
}

func TestIntegrationWithin(t *testing.T) {
	outer := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}},
	})
	inner := geojson.NewPolygon([][]geojson.Position{
		{{1, 1}, {9, 1}, {9, 9}, {1, 9}, {1, 1}},
	})
	within, err := Within(inner, outer)
	if err != nil {
		t.Fatal(err)
	}
	if !within {
		t.Error("inner should be within outer")
	}
}

func TestIntegrationTransformTranslate(t *testing.T) {
	pt := geojson.NewPoint([]float64{0, 0})
	translated, err := TransformTranslate(pt, 10, 5)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(translated)
	if math.Abs(coord[0]-10) > eps || math.Abs(coord[1]-5) > eps {
		t.Errorf("expected (10,5), got (%.4f, %.4f)", coord[0], coord[1])
	}
}

func TestIntegrationTransformRotate(t *testing.T) {
	pt := geojson.NewPoint([]float64{1, 0})
	origin := geojson.Position{0, 0}
	rotated, err := TransformRotate(pt, 90, origin)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(rotated)
	if math.Abs(coord[0]) > 0.01 || math.Abs(coord[1]-1) > 0.01 {
		t.Errorf("expected ~(0,1), got (%.4f, %.4f)", coord[0], coord[1])
	}
}

func TestIntegrationPolygonUnion(t *testing.T) {
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
}

func TestIntegrationSimplify(t *testing.T) {
	line := geojson.NewLineString([]geojson.Position{
		{0, 0}, {1, 0.1}, {2, 0}, {3, 0.1}, {4, 0},
	})
	result, err := Simplify(line, 0.3, false)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal("simplify should not be nil")
	}
	ls, ok := result.Geometry.(*geojson.LineString)
	if !ok {
		t.Fatal("expected LineString")
	}
	if len(ls.Coordinates) < 2 {
		t.Error("simplified should have at least 2 points")
	}
}

func TestIntegrationConvexHull(t *testing.T) {
	pts := []geojson.Position{
		{0, 0}, {10, 0}, {10, 10}, {0, 10}, {5, 5},
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
	if hull == nil {
		t.Fatal("convex hull should not be nil")
	}
}

func TestIntegrationKMeans(t *testing.T) {
	pts := []geojson.Position{
		{0, 0}, {1, 0}, {0, 1}, {10, 10}, {11, 10}, {10, 11},
	}
	features := make([]*geojson.Feature, len(pts))
	for i, p := range pts {
		features[i] = geojson.NewFeature(geojson.NewPoint(p), nil)
	}
	fc := geojson.NewFeatureCollection(features)
	clustered, err := ClustersKMeans(fc, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(clustered.Features) != len(pts) {
		t.Errorf("expected %d features, got %d", len(pts), len(clustered.Features))
	}
}

func TestIntegrationDbscan(t *testing.T) {
	pts := []geojson.Position{
		{0, 0}, {1, 0}, {0, 1}, {10, 10}, {11, 10}, {10, 11},
	}
	features := make([]*geojson.Feature, len(pts))
	for i, p := range pts {
		features[i] = geojson.NewFeature(geojson.NewPoint(p), nil)
	}
	fc := geojson.NewFeatureCollection(features)
	clustered, err := ClustersDbscan(fc, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(clustered.Features) != len(pts) {
		t.Errorf("expected %d features, got %d", len(pts), len(clustered.Features))
	}
}

func TestIntegrationTag(t *testing.T) {
	pts := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint([]float64{1, 1}), map[string]any{"name": "a"}),
		geojson.NewFeature(geojson.NewPoint([]float64{5, 5}), map[string]any{"name": "b"}),
	})
	polys := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPolygon([][]geojson.Position{
			{{0, 0}, {3, 0}, {3, 3}, {0, 3}, {0, 0}},
		}), map[string]any{"zone": "A"}),
	})
	result, err := Tag(pts, polys, "zone", "zone")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 2 {
		t.Fatal("tag should return same number of points")
	}
}

func TestIntegrationCollect(t *testing.T) {
	polys := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPolygon([][]geojson.Position{
			{{0, 0}, {3, 0}, {3, 3}, {0, 3}, {0, 0}},
		}), nil),
	})
	pts := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint([]float64{1, 1}), map[string]any{"val": "x"}),
	})
	result, err := Collect(polys, pts, "val", "values")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 1 {
		t.Fatal("collect should return polys")
	}
}

func TestIntegrationGreatCircle(t *testing.T) {
	from := geojson.NewPoint([]float64{0, 0})
	to := geojson.NewPoint([]float64{10, 10})
	gc, err := GreatCircle(from, to)
	if err != nil {
		t.Fatal(err)
	}
	if gc == nil {
		t.Fatal("great circle should not be nil")
	}
}

func TestIntegrationRhumbDistance(t *testing.T) {
	from := geojson.NewPoint([]float64{0, 0})
	to := geojson.NewPoint([]float64{10, 0})
	d, err := RhumbDistance(from, to, UnitKilometers)
	if err != nil {
		t.Fatal(err)
	}
	if d <= 0 {
		t.Error("rhumb distance should be positive")
	}
}

func TestIntegrationIsobands(t *testing.T) {
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
		t.Error("expected at least 1 band")
	}
	for _, f := range result.Features {
		if f.Geometry.Type() != geojson.TypePolygon {
			t.Errorf("expected Polygon band, got %s", f.Geometry.Type())
		}
	}
}

func TestIntegrationIsolines(t *testing.T) {
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
	result, err := Isolines(fc, IsolinesOptions{ZProperty: "z", Breaks: []float64{5, 10, 15}})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) == 0 {
		t.Error("expected at least 1 contour line")
	}
	for _, f := range result.Features {
		if f.Geometry.Type() != geojson.TypeLineString {
			t.Errorf("expected LineString contour, got %s", f.Geometry.Type())
		}
	}
}

func TestIntegrationRhumbDestination(t *testing.T) {
	origin := geojson.NewPoint([]float64{0, 0})
	dest, err := RhumbDestination(origin, 111.32, 90)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(dest)
	if math.Abs(coord[0]-1) > 0.1 || math.Abs(coord[1]) > 0.1 {
		t.Errorf("unexpected destination (%.4f, %.4f)", coord[0], coord[1])
	}
}
