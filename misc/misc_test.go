package misc

import (
	"testing"

	"github.com/ibinh/turf-go/geojson"
)

func pt(lng, lat float64) geojson.Position {
	return geojson.Position{lng, lat}
}

func TestClone(t *testing.T) {
	f := geojson.NewFeature(geojson.NewPoint(pt(10, 20)), map[string]any{"name": "test"})
	cloned, err := Clone(f)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(cloned)
	if coord[0] != 10 || coord[1] != 20 {
		t.Errorf("clone: expected (10,20), got (%v,%v)", coord[0], coord[1])
	}
	if cloned.Properties["name"] != "test" {
		t.Error("clone: properties not preserved")
	}
}

func TestCombinePoint(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(pt(1, 2)), nil),
		geojson.NewFeature(geojson.NewPoint(pt(3, 4)), nil),
	})
	result, err := Combine(fc)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 1 {
		t.Fatalf("expected 1 feature, got %d", len(result.Features))
	}
	mp, ok := result.Features[0].Geometry.(*geojson.MultiPoint)
	if !ok {
		t.Fatal("expected MultiPoint")
	}
	if len(mp.Coordinates) != 2 {
		t.Errorf("expected 2 coords, got %d", len(mp.Coordinates))
	}
}

func TestCombineLineString(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewLineString([]geojson.Position{{0, 0}, {1, 1}}), nil),
		geojson.NewFeature(geojson.NewLineString([]geojson.Position{{2, 2}, {3, 3}}), nil),
	})
	result, err := Combine(fc)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 1 {
		t.Fatalf("expected 1 feature, got %d", len(result.Features))
	}
	ml, ok := result.Features[0].Geometry.(*geojson.MultiLineString)
	if !ok {
		t.Fatal("expected MultiLineString")
	}
	if len(ml.Coordinates) != 2 {
		t.Errorf("expected 2 lines, got %d", len(ml.Coordinates))
	}
}

func TestCombinePolygon(t *testing.T) {
	p1 := geojson.NewPolygon([][]geojson.Position{{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}}})
	p2 := geojson.NewPolygon([][]geojson.Position{{{2, 2}, {2, 3}, {3, 3}, {3, 2}, {2, 2}}})
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(p1, nil),
		geojson.NewFeature(p2, nil),
	})
	result, err := Combine(fc)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 1 {
		t.Fatalf("expected 1 feature, got %d", len(result.Features))
	}
	mp, ok := result.Features[0].Geometry.(*geojson.MultiPolygon)
	if !ok {
		t.Fatal("expected MultiPolygon")
	}
	if len(mp.Coordinates) != 2 {
		t.Errorf("expected 2 polys, got %d", len(mp.Coordinates))
	}
}

func TestExplodeMultiPoint(t *testing.T) {
	mp := geojson.NewMultiPoint([]geojson.Position{{1, 2}, {3, 4}})
	f := geojson.NewFeature(mp, nil)
	result, err := Explode(f)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 2 {
		t.Fatalf("expected 2 features, got %d", len(result.Features))
	}
	for i, feat := range result.Features {
		if _, ok := feat.Geometry.(*geojson.Point); !ok {
			t.Errorf("feature %d: expected Point", i)
		}
	}
}

func TestExplodeMultiLineString(t *testing.T) {
	ml := geojson.NewMultiLineString([][]geojson.Position{
		{{0, 0}, {1, 1}},
		{{2, 2}, {3, 3}},
	})
	f := geojson.NewFeature(ml, nil)
	result, err := Explode(f)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 2 {
		t.Fatalf("expected 2 features, got %d", len(result.Features))
	}
	for i, feat := range result.Features {
		if _, ok := feat.Geometry.(*geojson.LineString); !ok {
			t.Errorf("feature %d: expected LineString", i)
		}
	}
}

func TestExplodeMultiPolygon(t *testing.T) {
	mp := geojson.NewMultiPolygon([][][]geojson.Position{
		{{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}}},
		{{{2, 2}, {2, 3}, {3, 3}, {3, 2}, {2, 2}}},
	})
	f := geojson.NewFeature(mp, nil)
	result, err := Explode(f)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 2 {
		t.Fatalf("expected 2 features, got %d", len(result.Features))
	}
	for i, feat := range result.Features {
		if _, ok := feat.Geometry.(*geojson.Polygon); !ok {
			t.Errorf("feature %d: expected Polygon", i)
		}
	}
}

func TestPointsWithinPolygon(t *testing.T) {
	poly := geojson.NewFeature(geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}},
	}), nil)
	pts := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(pt(5, 5)), nil),
		geojson.NewFeature(geojson.NewPoint(pt(15, 15)), nil),
		geojson.NewFeature(geojson.NewPoint(pt(1, 1)), nil),
	})
	result, err := PointsWithinPolygon(pts, poly)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 2 {
		t.Errorf("expected 2 points inside, got %d", len(result.Features))
	}
}

func TestPlanepoint(t *testing.T) {
	tri := geojson.NewPolygon([][]geojson.Position{
		{{0, 0, 0}, {1, 0, 10}, {0, 1, 20}, {0, 0, 0}},
	})
	z, err := Planepoint(geojson.NewPoint(pt(0.25, 0.25)), tri)
	if err != nil {
		t.Fatal(err)
	}
	if z < 7 || z > 8 {
		t.Errorf("expected z ~7.5, got %f", z)
	}
}

func TestTesselate(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}},
	})
	result, err := Tesselate(poly)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 2 {
		t.Errorf("expected 2 triangles, got %d", len(result.Features))
	}
	for i, f := range result.Features {
		if _, ok := f.Geometry.(*geojson.Polygon); !ok {
			t.Errorf("feature %d: expected Polygon", i)
		}
	}
}

func TestFlatten(t *testing.T) {
	mp := geojson.NewMultiPoint([]geojson.Position{{1, 2}, {3, 4}})
	f := geojson.NewFeature(mp, nil)
	result, err := Flatten(f)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 2 {
		t.Fatalf("expected 2 features, got %d", len(result.Features))
	}
	for i, feat := range result.Features {
		if _, ok := feat.Geometry.(*geojson.Point); !ok {
			t.Errorf("feature %d: expected Point", i)
		}
	}
}

func TestFlattenPolygon(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}}})
	f := geojson.NewFeature(poly, nil)
	result, err := Flatten(f)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) != 1 {
		t.Fatalf("expected 1 feature, got %d", len(result.Features))
	}
	if _, ok := result.Features[0].Geometry.(*geojson.Polygon); !ok {
		t.Error("expected Polygon")
	}
}

func TestTesselateConcave(t *testing.T) {
	poly := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {2, 0}, {2, 2}, {1, 1}, {0, 2}, {0, 0}},
	})
	result, err := Tesselate(poly)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Features) == 0 {
		t.Error("expected at least 1 triangle")
	}
}

func TestCloneGeometry(t *testing.T) {
	p := geojson.NewPoint(pt(-73, 40))
	cloned, err := Clone(p)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(cloned)
	if coord[0] != -73 || coord[1] != 40 {
		t.Error("clone geometry: coords mismatch")
	}
}

func TestPlanepointNoZ(t *testing.T) {
	tri := geojson.NewPolygon([][]geojson.Position{
		{{0, 0}, {1, 0}, {0, 1}, {0, 0}},
	})
	z, err := Planepoint(geojson.NewPoint(pt(0.25, 0.25)), tri)
	if err != nil {
		t.Fatal(err)
	}
	if z != 0 {
		t.Errorf("expected 0 when no Z values, got %f", z)
	}
}
