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
