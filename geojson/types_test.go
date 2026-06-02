package geojson

import (
	"encoding/json"
	"testing"
)

func TestPointRoundtrip(t *testing.T) {
	f := NewFeature(NewPoint(Position{1.0, 2.0}), nil)
	f.SetBBox([]float64{0, 0, 2, 2})
	f.ID = "test1"

	data, err := json.Marshal(f)
	if err != nil {
		t.Fatal(err)
	}

	var got Feature
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.Type() != TypeFeature {
		t.Errorf("expected type Feature, got %s", got.Type())
	}
	p, ok := got.Geometry.(*Point)
	if !ok {
		t.Fatalf("expected *Point geometry, got %T", got.Geometry)
	}
	if p.Coordinates[0] != 1.0 || p.Coordinates[1] != 2.0 {
		t.Errorf("unexpected coordinates: %v", p.Coordinates)
	}
	if got.ID != "test1" {
		t.Errorf("unexpected id: %v", got.ID)
	}
}

func TestLineStringRoundtrip(t *testing.T) {
	coords := []Position{{0, 0}, {1, 1}, {2, 2}}
	f := NewFeature(NewLineString(coords), map[string]any{"name": "route"})
	data, err := json.Marshal(f)
	if err != nil {
		t.Fatal(err)
	}

	var got Feature
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	ls, ok := got.Geometry.(*LineString)
	if !ok {
		t.Fatalf("expected *LineString geometry, got %T", got.Geometry)
	}
	if len(ls.Coordinates) != 3 {
		t.Errorf("expected 3 coords, got %d", len(ls.Coordinates))
	}
	if got.Properties["name"] != "route" {
		t.Errorf("unexpected properties: %v", got.Properties)
	}
}

func TestPolygonRoundtrip(t *testing.T) {
	ring := []Position{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}}
	f := NewFeature(NewPolygon([][]Position{ring}), nil)
	data, err := json.Marshal(f)
	if err != nil {
		t.Fatal(err)
	}

	var got Feature
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	poly, ok := got.Geometry.(*Polygon)
	if !ok {
		t.Fatalf("expected *Polygon, got %T", got.Geometry)
	}
	if len(poly.Coordinates) != 1 {
		t.Errorf("expected 1 ring, got %d", len(poly.Coordinates))
	}
}

func TestFeatureCollectionRoundtrip(t *testing.T) {
	f1 := NewFeature(NewPoint(Position{1, 2}), nil)
	f2 := NewFeature(NewPoint(Position{3, 4}), map[string]any{"label": "B"})
	fc := NewFeatureCollection([]*Feature{f1, f2})
	data, err := json.Marshal(fc)
	if err != nil {
		t.Fatal(err)
	}

	var got FeatureCollection
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if len(got.Features) != 2 {
		t.Errorf("expected 2 features, got %d", len(got.Features))
	}
	if got.Features[1].Properties["label"] != "B" {
		t.Errorf("unexpected property")
	}
}

func TestGeometryCollectionRoundtrip(t *testing.T) {
	gc := NewGeometryCollection([]Geometry{
		NewPoint(Position{1, 2}),
		NewLineString([]Position{{0, 0}, {1, 1}}),
	})
	data, err := json.Marshal(gc)
	if err != nil {
		t.Fatal(err)
	}

	var got GeometryCollection
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if len(got.Geometries) != 2 {
		t.Errorf("expected 2 geometries, got %d", len(got.Geometries))
	}
	if got.Geometries[0].Type() != TypePoint {
		t.Errorf("expected Point, got %s", got.Geometries[0].Type())
	}
	if got.Geometries[1].Type() != TypeLineString {
		t.Errorf("expected LineString, got %s", got.Geometries[1].Type())
	}
}

func TestMultiGeometryRoundtrips(t *testing.T) {
	t.Run("MultiPoint", func(t *testing.T) {
		f := NewFeature(NewMultiPoint([]Position{{0, 0}, {1, 1}}), nil)
		data, _ := json.Marshal(f)
		var got Feature
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatal(err)
		}
		if got.Geometry.Type() != TypeMultiPoint {
			t.Errorf("expected MultiPoint, got %s", got.Geometry.Type())
		}
	})

	t.Run("MultiLineString", func(t *testing.T) {
		f := NewFeature(NewMultiLineString([][]Position{{{0, 0}, {1, 1}}, {{2, 2}, {3, 3}}}), nil)
		data, _ := json.Marshal(f)
		var got Feature
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatal(err)
		}
		if got.Geometry.Type() != TypeMultiLineString {
			t.Errorf("expected MultiLineString, got %s", got.Geometry.Type())
		}
	})

	t.Run("MultiPolygon", func(t *testing.T) {
		f := NewFeature(NewMultiPolygon([][][]Position{{{{0, 0}, {0, 5}, {5, 5}, {5, 0}, {0, 0}}}}), nil)
		data, _ := json.Marshal(f)
		var got Feature
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatal(err)
		}
		if got.Geometry.Type() != TypeMultiPolygon {
			t.Errorf("expected MultiPolygon, got %s", got.Geometry.Type())
		}
	})
}

func TestNullGeometry(t *testing.T) {
	f := NewFeature(nil, nil)
	data, err := json.Marshal(f)
	if err != nil {
		t.Fatal(err)
	}

	var got Feature
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.Geometry != nil {
		t.Errorf("expected nil geometry, got %T", got.Geometry)
	}
}

func TestUnmarshalGeoJSON(t *testing.T) {
	t.Run("bare Point", func(t *testing.T) {
		data := []byte(`{"type":"Point","coordinates":[1,2]}`)
		gj, err := UnmarshalGeoJSON(data)
		if err != nil {
			t.Fatal(err)
		}
		if gj.Type() != TypePoint {
			t.Errorf("expected Point, got %s", gj.Type())
		}
	})

	t.Run("Feature", func(t *testing.T) {
		data := []byte(`{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2]},"properties":{"a":1}}`)
		gj, err := UnmarshalGeoJSON(data)
		if err != nil {
			t.Fatal(err)
		}
		f, ok := gj.(*Feature)
		if !ok {
			t.Fatalf("expected *Feature, got %T", gj)
		}
		if f.Geometry.Type() != TypePoint {
			t.Errorf("expected Point geometry, got %s", f.Geometry.Type())
		}
		if f.Properties["a"] != float64(1) {
			t.Errorf("unexpected property")
		}
	})

	t.Run("FeatureCollection", func(t *testing.T) {
		data := []byte(`{"type":"FeatureCollection","features":[{"type":"Feature","geometry":{"type":"Point","coordinates":[0,0]}}]}`)
		gj, err := UnmarshalGeoJSON(data)
		if err != nil {
			t.Fatal(err)
		}
		fc, ok := gj.(*FeatureCollection)
		if !ok {
			t.Fatalf("expected *FeatureCollection, got %T", gj)
		}
		if len(fc.Features) != 1 {
			t.Errorf("expected 1 feature, got %d", len(fc.Features))
		}
	})
}

func TestInvariantGetCoord(t *testing.T) {
	f := NewFeature(NewPoint(Position{10, 20}), nil)
	coord, err := GetCoord(f)
	if err != nil {
		t.Fatal(err)
	}
	if coord[0] != 10 || coord[1] != 20 {
		t.Errorf("unexpected coord: %v", coord)
	}

	coord, err = GetCoord(f.Geometry)
	if err != nil {
		t.Fatal(err)
	}
	if coord[0] != 10 || coord[1] != 20 {
		t.Errorf("unexpected coord from bare geometry: %v", coord)
	}
}

func TestInvariantGetCoords(t *testing.T) {
	coords := []Position{{0, 0}, {1, 1}, {2, 2}}
	f := NewFeature(NewLineString(coords), nil)
	got, err := GetCoords(f)
	if err != nil {
		t.Fatal(err)
	}
	gotCoords, ok := got.([]Position)
	if !ok {
		t.Fatalf("expected []Position, got %T", got)
	}
	if len(gotCoords) != 3 {
		t.Errorf("expected 3 coords, got %d", len(gotCoords))
	}
}

func TestInvariantCoordAll(t *testing.T) {
	f := NewFeature(NewPolygon([][]Position{{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}}}), nil)
	coords, err := CoordAll(f)
	if err != nil {
		t.Fatal(err)
	}
	if len(coords) != 5 {
		t.Errorf("expected 5 coords, got %d", len(coords))
	}
}

func TestInvariantCollectionOf(t *testing.T) {
	fc := NewFeatureCollection([]*Feature{
		NewFeature(NewPoint(Position{0, 0}), nil),
		NewFeature(NewPoint(Position{1, 1}), nil),
	})
	if err := CollectionOf(fc, TypePoint); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	fc.Features = append(fc.Features, NewFeature(NewLineString([]Position{{0, 0}, {1, 1}}), nil))
	if err := CollectionOf(fc, TypePoint); err == nil {
		t.Errorf("expected error for mixed types")
	}
}

func TestPositionAccessors(t *testing.T) {
	p := Position{1.0, 2.0, 3.0}
	if p.Lng() != 1.0 {
		t.Errorf("expected lng=1.0, got %f", p.Lng())
	}
	if p.Lat() != 2.0 {
		t.Errorf("expected lat=2.0, got %f", p.Lat())
	}
	if p.Elevation() != 3.0 {
		t.Errorf("expected elevation=3.0, got %f", p.Elevation())
	}
}

func TestBBox(t *testing.T) {
	f := NewFeature(NewPoint(Position{1, 2}), nil)
	f.SetBBox([]float64{0, 0, 2, 2})
	data, err := json.Marshal(f)
	if err != nil {
		t.Fatal(err)
	}
	var got Feature
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	bbox := got.BBox()
	if len(bbox) != 4 || bbox[0] != 0 || bbox[3] != 2 {
		t.Errorf("unexpected bbox: %v", bbox)
	}
}
