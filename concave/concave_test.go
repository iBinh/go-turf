package concave

import (
	"testing"

	"github.com/ibinh/turf-go/geojson"
)

func TestConcaveHullTriangle(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 0}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{5, 10}), nil),
	})

	hull, err := ConcaveHull(fc)
	if err != nil {
		t.Fatal(err)
	}
	if hull == nil {
		t.Fatal("expected hull, got nil")
	}
	poly, ok := hull.Geometry.(*geojson.Polygon)
	if !ok {
		t.Fatal("expected polygon geometry")
	}
	if len(poly.Coordinates) == 0 || len(poly.Coordinates[0]) < 4 {
		t.Fatal("expected ring with closure")
	}
}

func TestConcaveHullMinPointsError(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 0}), nil),
	})

	_, err := ConcaveHull(fc)
	if err == nil {
		t.Fatal("expected error for < 3 points")
	}
}

func TestConcaveHullMaxEdge(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{6, 0}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{3, 6}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 6}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{6, 6}), nil),
	})

	hull, err := ConcaveHull(fc, ConcaveOptions{MaxEdge: 7})
	if err != nil {
		t.Fatal(err)
	}
	if hull == nil {
		t.Fatal("expected hull, got nil")
	}
}

func TestConcaveHullSquare(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 0}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 10}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 10}), nil),
	})

	hull, err := ConcaveHull(fc)
	if err != nil {
		t.Fatal(err)
	}
	if hull == nil {
		t.Fatal("expected hull, got nil")
	}
}
