package shortestpath

import (
	"testing"

	"github.com/ibinh/turf-go/geojson"
)

func TestShortestPathSimple(t *testing.T) {
	start := geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil)
	end := geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 10}), nil)

	// Simple L-shaped network
	network := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(
			geojson.NewLineString([]geojson.Position{{0, 0}, {5, 0}, {10, 0}, {10, 10}}),
			nil,
		),
	})

	path, err := ShortestPath(start, end, network)
	if err != nil {
		t.Fatal(err)
	}
	if path == nil {
		t.Fatal("expected path, got nil")
	}

	coords, err := geojson.GetCoords(path)
	if err != nil {
		t.Fatal(err)
	}
	pts, ok := coords.([]geojson.Position)
	if !ok || len(pts) < 2 {
		t.Fatal("expected path with at least 2 points")
	}
}

func TestShortestPathDirectLine(t *testing.T) {
	start := geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil)
	end := geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 10}), nil)

	network := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(
			geojson.NewLineString([]geojson.Position{{0, 0}, {10, 10}}),
			nil,
		),
	})

	path, err := ShortestPath(start, end, network)
	if err != nil {
		t.Fatal(err)
	}
	if path == nil {
		t.Fatal("expected path, got nil")
	}
}

func TestShortestPathNoNetwork(t *testing.T) {
	start := geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil)
	end := geojson.NewFeature(geojson.NewPoint(geojson.Position{10, 0}), nil)

	_, err := ShortestPath(start, end, nil)
	if err == nil {
		t.Fatal("expected error for nil network")
	}
}

func TestShortestPathDisconnected(t *testing.T) {
	start := geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil)
	end := geojson.NewFeature(geojson.NewPoint(geojson.Position{100, 100}), nil)

	network := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(
			geojson.NewLineString([]geojson.Position{{0, 0}, {5, 0}}),
			nil,
		),
		geojson.NewFeature(
			geojson.NewLineString([]geojson.Position{{50, 50}, {100, 100}}),
			nil,
		),
	})

	_, err := ShortestPath(start, end, network)
	if err == nil {
		t.Fatal("expected error for disconnected network")
	}
}

func TestShortestPathMaxDistance(t *testing.T) {
	start := geojson.NewFeature(geojson.NewPoint(geojson.Position{100, 100}), nil)
	end := geojson.NewFeature(geojson.NewPoint(geojson.Position{200, 200}), nil)

	network := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(
			geojson.NewLineString([]geojson.Position{{0, 0}, {10, 10}}),
			nil,
		),
	})

	_, err := ShortestPath(start, end, network, ShortestPathOptions{MaxDistance: 0.1})
	if err == nil {
		t.Fatal("expected error when start/end too far from network")
	}
}
