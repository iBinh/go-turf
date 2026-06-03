package center

import (
	"math"
	"testing"
	"github.com/ibinh/turf-go/geojson"
)

func TestCenter(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{2, 0}), nil),
	})

	c, err := Center(fc)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(c)
	if math.Abs(coord[0]-1) > 0.001 {
		t.Errorf("expected center x=1, got %f", coord[0])
	}
}

func TestCenterSinglePoint(t *testing.T) {
	f := geojson.NewFeature(
		geojson.NewPoint(geojson.Position{10, 20}),
		nil,
	)
	c, err := Center(f)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(c)
	if coord[0] != 10 || coord[1] != 20 {
		t.Errorf("expected center at [10,20], got %v", coord)
	}
}

func TestCentroid(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 10}), nil),
	})

	c, err := Centroid(fc)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(c)
	if math.Abs(coord[1]-5) > 0.001 {
		t.Errorf("expected centroid y=5, got %f", coord[1])
	}
}

func TestCenterMean(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), map[string]any{"weight": float64(1)}),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 0}), map[string]any{"weight": float64(3)}),
	})

	c, err := CenterMean(fc, nil, "weight")
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(c)
	if math.Abs(coord[0]-7.5) > 0.001 {
		t.Errorf("expected weighted center x=7.5, got %f", coord[0])
	}
}

func TestCenterOfMass(t *testing.T) {
	f := geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{{{0, 0}, {0, 2}, {2, 2}, {2, 0}, {0, 0}}}),
		nil,
	)
	com, err := CenterOfMass(f)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(com)
	if math.Abs(coord[0]-0.8) > 0.01 || math.Abs(coord[1]-0.8) > 0.01 {
		t.Errorf("expected center of mass at [0.8,0.8], got %v", coord)
	}
}

func TestPointOnFeature_EmptyGeometry(t *testing.T) {
	gc := geojson.NewGeometryCollection([]geojson.Geometry{})
	_, err := PointOnFeature(gc)
	if err == nil {
		t.Error("expected error for unknown geometry type")
	}
}

func TestPointOnFeature_Point(t *testing.T) {
	pt := geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 20}), nil)
	result, err := PointOnFeature(pt)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(result)
	if coord[0] != 10 || coord[1] != 20 {
		t.Errorf("expected [10,20], got %v", coord)
	}
}

func TestPointOnFeature_MultiPoint(t *testing.T) {
	mp := geojson.NewFeature(geojson.NewMultiPoint([]geojson.Position{{1, 2}, {3, 4}}), nil)
	result, err := PointOnFeature(mp)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(result)
	if coord[0] != 1 || coord[1] != 2 {
		t.Errorf("expected first point [1,2], got %v", coord)
	}
}

func TestPointOnFeature_MultiPointEmpty(t *testing.T) {
	mp := geojson.NewFeature(geojson.NewMultiPoint([]geojson.Position{}), nil)
	_, err := PointOnFeature(mp)
	if err == nil {
		t.Error("expected error for empty MultiPoint")
	}
}

func TestPointOnFeature_LineString(t *testing.T) {
	ls := geojson.NewFeature(geojson.NewLineString([]geojson.Position{{0, 0}, {10, 10}, {20, 20}}), nil)
	result, err := PointOnFeature(ls)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(result)
	if coord[0] != 10 || coord[1] != 10 {
		t.Errorf("expected midpoint [10,10], got %v", coord)
	}
}

func TestPointOnFeature_LineStringShort(t *testing.T) {
	ls := geojson.NewFeature(geojson.NewLineString([]geojson.Position{{0, 0}}), nil)
	_, err := PointOnFeature(ls)
	if err == nil {
		t.Error("expected error for line too short")
	}
}

func TestPointOnFeature_MultiLineString(t *testing.T) {
	mls := geojson.NewFeature(geojson.NewMultiLineString([][]geojson.Position{{{0, 0}, {10, 10}, {20, 20}}}), nil)
	result, err := PointOnFeature(mls)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(result)
	if coord[0] != 10 || coord[1] != 10 {
		t.Errorf("expected midpoint [10,10], got %v", coord)
	}
}

func TestPointOnFeature_MultiLineStringEmpty(t *testing.T) {
	mls := geojson.NewFeature(geojson.NewMultiLineString([][]geojson.Position{}), nil)
	_, err := PointOnFeature(mls)
	if err == nil {
		t.Error("expected error for empty MultiLineString")
	}
}

func TestPointOnFeature_Polygon(t *testing.T) {
	poly := geojson.NewFeature(geojson.NewPolygon([][]geojson.Position{{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}}}), nil)
	result, err := PointOnFeature(poly)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(result)
	if coord[0] != 4 || coord[1] != 4 {
		t.Errorf("expected centroid [4,4], got %v", coord)
	}
}

func TestPointOnFeature_MultiPolygon(t *testing.T) {
	mp := geojson.NewFeature(geojson.NewMultiPolygon([][][]geojson.Position{{{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}}}}), nil)
	result, err := PointOnFeature(mp)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(result)
	if coord[0] != 4 || coord[1] != 4 {
		t.Errorf("expected centroid [4,4], got %v", coord)
	}
}

func TestPointOnFeature_Default(t *testing.T) {
	ls := geojson.NewLineString([]geojson.Position{{0, 0}, {1, 1}})
	result, err := PointOnFeature(ls)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(result)
	if coord[0] != 1 || coord[1] != 1 {
		t.Errorf("expected midpoint [1,1], got %v", coord)
	}
}

func TestCenterOfMassEmpty(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{})
	result, err := CenterOfMass(fc)
	if err != nil {
		t.Fatal(err)
	}
	if result != nil {
		t.Error("expected nil for empty FC")
	}
}

func TestCentroidEmpty(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{})
	result, err := Centroid(fc)
	if err != nil {
		t.Fatal(err)
	}
	if result != nil {
		t.Error("expected nil for empty FC")
	}
}

func TestCenterMeanNoWeight(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 10}), nil),
	})
	c, err := CenterMean(fc, nil)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(c)
	if coord[0] != 5 || coord[1] != 5 {
		t.Errorf("expected [5,5], got %v", coord)
	}
}

func TestCenterMeanWithProperties(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
	})
	c, err := CenterMean(fc, map[string]any{"name": "test"})
	if err != nil {
		t.Fatal(err)
	}
	if c.Properties["name"] != "test" {
		t.Errorf("expected properties to carry through, got %v", c.Properties)
	}
}

func TestCenterMeanZeroWeight(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), map[string]any{"w": float64(0)}),
	})
	c, err := CenterMean(fc, nil, "w")
	if err != nil {
		t.Fatal(err)
	}
	if c != nil {
		t.Error("expected nil for zero total weight")
	}
}

func TestCenterMedian(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 10}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 20}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 30}), nil),
	})

	c, err := CenterMedian(fc, nil)
	if err != nil {
		t.Fatal(err)
	}
	coord, _ := geojson.GetCoord(c)
	if math.Abs(coord[1]-20) > 0.001 {
		t.Errorf("expected median lat 20, got %f", coord[1])
	}
}
