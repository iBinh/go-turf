package meta

import (
	"testing"

	"github.com/ibinh/turf-go/geojson"
)

func TestCoordEach(t *testing.T) {
	p := geojson.NewFeature(geojson.NewPoint(geojson.Position{1, 2}), nil)

	count := 0
	err := CoordEach(p, func(coord geojson.Position, index int) error {
		count++
		if coord[0] != 1 || coord[1] != 2 {
			t.Errorf("unexpected coord: %v", coord)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("expected 1 coord, got %d", count)
	}
}

func TestCoordEachMulti(t *testing.T) {
	ls := geojson.NewFeature(geojson.NewLineString([]geojson.Position{{0, 0}, {1, 1}, {2, 2}}), nil)

	count := 0
	err := CoordEach(ls, func(coord geojson.Position, index int) error {
		count++
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if count != 3 {
		t.Errorf("expected 3 coords, got %d", count)
	}
}

func TestCoordReduce(t *testing.T) {
	p := geojson.NewFeature(geojson.NewPoint(geojson.Position{5, 10}), nil)

	sum, err := CoordReduce(p, func(acc float64, coord geojson.Position, index int) (float64, error) {
		return acc + coord[0] + coord[1], nil
	}, 0.0)
	if err != nil {
		t.Fatal(err)
	}
	if sum != 15 {
		t.Errorf("expected sum 15, got %f", sum)
	}
}

func TestFeatureEach(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{1, 1}), nil),
	})

	count := 0
	err := FeatureEach(fc, func(f *geojson.Feature, index int) error {
		count++
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Errorf("expected 2 features, got %d", count)
	}
}

func TestFeatureEachSingle(t *testing.T) {
	f := geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil)

	count := 0
	err := FeatureEach(f, func(f *geojson.Feature, index int) error {
		count++
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("expected 1 feature, got %d", count)
	}
}

func TestFeatureReduce(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), map[string]any{"val": float64(1)}),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{1, 1}), map[string]any{"val": float64(2)}),
	})

	sum, err := FeatureReduce(fc, 0, func(acc int, f *geojson.Feature, index int) (int, error) {
		v, _ := f.Properties["val"].(float64)
		return acc + int(v), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if sum != 3 {
		t.Errorf("expected 3, got %d", sum)
	}
}

func TestGeomEach(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
		geojson.NewFeature(geojson.NewLineString([]geojson.Position{{0, 0}, {1, 1}}), nil),
	})

	types := []string{}
	err := GeomEach(fc, func(geom geojson.Geometry, index int) error {
		types = append(types, geom.Type())
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(types) != 2 || types[0] != "Point" || types[1] != "LineString" {
		t.Errorf("unexpected types: %v", types)
	}
}

func TestPropEach(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), map[string]any{"a": 1}),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{1, 1}), map[string]any{"b": 2}),
	})

	var keys []string
	err := PropEach(fc, func(props map[string]any, index int) error {
		for k := range props {
			keys = append(keys, k)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 2 {
		t.Errorf("expected 2 prop keys, got %d", len(keys))
	}
}

func TestFlattenMultiPoint(t *testing.T) {
	mp := geojson.NewFeature(geojson.NewMultiPoint([]geojson.Position{{0, 0}, {1, 1}, {2, 2}}), nil)

	count := 0
	err := FlattenEach(mp, func(f *geojson.Feature, index int) error {
		count++
		if f.Geometry.Type() != "Point" {
			t.Errorf("expected Point, got %s", f.Geometry.Type())
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if count != 3 {
		t.Errorf("expected 3 features, got %d", count)
	}
}

func TestCoordCount(t *testing.T) {
	p := geojson.NewFeature(geojson.NewPolygon([][]geojson.Position{{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}}}), nil)

	count, err := CoordCount(p)
	if err != nil {
		t.Fatal(err)
	}
	if count != 5 {
		t.Errorf("expected 5 coords, got %d", count)
	}
}

func TestGetFirstLastCoord(t *testing.T) {
	ls := geojson.NewFeature(geojson.NewLineString([]geojson.Position{{1, 2}, {3, 4}, {5, 6}}), nil)

	first, err := GetFirstCoord(ls)
	if err != nil {
		t.Fatal(err)
	}
	if first[0] != 1 || first[1] != 2 {
		t.Errorf("expected [1,2], got %v", first)
	}

	last, err := GetLastCoord(ls)
	if err != nil {
		t.Fatal(err)
	}
	if last[0] != 5 || last[1] != 6 {
		t.Errorf("expected [5,6], got %v", last)
	}
}
