package helpers

import (
	"encoding/json"
	"testing"

	"github.com/ibinh/turf-go/geojson"
)

func TestHelpersPoint(t *testing.T) {
	f := Point(geojson.Position{1, 2}, map[string]any{"name": "test"})
	if f.Type() != "Feature" {
		t.Errorf("expected Feature type")
	}
	p, ok := f.Geometry.(*geojson.Point)
	if !ok {
		t.Fatalf("expected *Point, got %T", f.Geometry)
	}
	if p.Coordinates[0] != 1 || p.Coordinates[1] != 2 {
		t.Errorf("unexpected coords: %v", p.Coordinates)
	}
	if f.Properties["name"] != "test" {
		t.Errorf("unexpected props: %v", f.Properties)
	}
}

func TestHelpersPointWithOptions(t *testing.T) {
	f := Point(geojson.Position{1, 2}, nil, geojson.WithBBox([]float64{0, 0, 2, 2}), geojson.WithID("abc"))
	if f.ID != "abc" {
		t.Errorf("expected id abc, got %v", f.ID)
	}
	if len(f.BBox()) != 4 {
		t.Errorf("expected bbox")
	}
}

func TestHelpersLineString(t *testing.T) {
	coords := []geojson.Position{{0, 0}, {1, 1}, {2, 2}}
	f := LineString(coords, nil)
	ls, ok := f.Geometry.(*geojson.LineString)
	if !ok {
		t.Fatalf("expected *LineString, got %T", f.Geometry)
	}
	if len(ls.Coordinates) != 3 {
		t.Errorf("expected 3 coords")
	}
}

func TestHelpersPolygon(t *testing.T) {
	ring := []geojson.Position{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}}
	f := Polygon([][]geojson.Position{ring}, nil)
	poly, ok := f.Geometry.(*geojson.Polygon)
	if !ok {
		t.Fatalf("expected *Polygon, got %T", f.Geometry)
	}
	if len(poly.Coordinates) != 1 {
		t.Errorf("expected 1 ring")
	}
}

func TestHelpersMultiPoint(t *testing.T) {
	f := MultiPoint([]geojson.Position{{0, 0}, {1, 1}}, nil)
	mp, ok := f.Geometry.(*geojson.MultiPoint)
	if !ok {
		t.Fatalf("expected *MultiPoint, got %T", f.Geometry)
	}
	if len(mp.Coordinates) != 2 {
		t.Errorf("expected 2 coords")
	}
}

func TestHelpersGeometryCollection(t *testing.T) {
	f := GeometryCollection([]geojson.Geometry{
		geojson.NewPoint(geojson.Position{1, 2}),
		geojson.NewLineString([]geojson.Position{{0, 0}, {1, 1}}),
	}, nil)
	gc, ok := f.Geometry.(*geojson.GeometryCollection)
	if !ok {
		t.Fatalf("expected *GeometryCollection, got %T", f.Geometry)
	}
	if len(gc.Geometries) != 2 {
		t.Errorf("expected 2 geometries")
	}
}

func TestHelpersJSONRoundtrip(t *testing.T) {
	f := Point(geojson.Position{1, 2}, nil, geojson.WithID("p1"))
	data, err := json.Marshal(f)
	if err != nil {
		t.Fatal(err)
	}

	var g geojson.Feature
	if err := json.Unmarshal(data, &g); err != nil {
		t.Fatal(err)
	}
	if g.ID != "p1" {
		t.Errorf("expected id p1, got %v", g.ID)
	}
}
