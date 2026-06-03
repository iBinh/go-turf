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

func TestCoordReduceFC(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{1, 2}), nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{3, 4}), nil),
	})
	sum, err := CoordReduce(fc, func(acc float64, c geojson.Position, idx int) (float64, error) {
		return acc + c[0] + c[1], nil
	}, 0.0)
	if err != nil {
		t.Fatal(err)
	}
	if sum != 10 {
		t.Errorf("expected 10, got %f", sum)
	}
}

func TestCoordEachFCNilGeometry(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(nil, nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{5, 6}), nil),
	})
	count := 0
	err := CoordEach(fc, func(c geojson.Position, idx int) error {
		count++
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("expected 1 coord (nil skipped), got %d", count)
	}
}

func TestCoordEachNonFeatureErr(t *testing.T) {
	err := CoordEach("not a geometry", func(c geojson.Position, idx int) error {
		return nil
	})
	if err == nil {
		t.Error("expected error for non-geometry input")
	}
}

func TestCoordReduceMultiPoint(t *testing.T) {
	mp := geojson.NewFeature(geojson.NewMultiPoint([]geojson.Position{{0, 0}, {1, 1}, {2, 2}}), nil)
	count, err := CoordReduce(mp, func(acc int, c geojson.Position, idx int) (int, error) {
		return acc + 1, nil
	}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if count != 3 {
		t.Errorf("expected 3 coords from MultiPoint, got %d", count)
	}
}

func TestCoordReduceMultiLineString(t *testing.T) {
	mls := geojson.NewFeature(geojson.NewMultiLineString([][]geojson.Position{{{0, 0}, {1, 1}}, {{2, 2}, {3, 3}}}), nil)
	count, err := CoordReduce(mls, func(acc int, c geojson.Position, idx int) (int, error) {
		return acc + 1, nil
	}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Errorf("expected 4 coords from MultiLineString, got %d", count)
	}
}

func TestCoordReduceMultiPolygon(t *testing.T) {
	mpoly := geojson.NewFeature(geojson.NewMultiPolygon([][][]geojson.Position{{{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}}}}), nil)
	count, err := CoordReduce(mpoly, func(acc int, c geojson.Position, idx int) (int, error) {
		return acc + 1, nil
	}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if count != 5 {
		t.Errorf("expected 5 coords from MultiPolygon, got %d", count)
	}
}

func TestCoordReduceGeometryCollection(t *testing.T) {
	gc := geojson.NewFeature(geojson.NewGeometryCollection([]geojson.Geometry{
		geojson.NewPoint(geojson.Position{1, 2}),
		geojson.NewLineString([]geojson.Position{{3, 4}, {5, 6}}),
	}), nil)
	count, err := CoordReduce(gc, func(acc int, c geojson.Position, idx int) (int, error) {
		return acc + 1, nil
	}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if count != 3 {
		t.Errorf("expected 3 coords from GeometryCollection, got %d", count)
	}
}

func TestGeomReduceGeometry(t *testing.T) {
	p := geojson.NewPoint(geojson.Position{1, 2})
	count, err := GeomReduce(p, 0, func(acc int, g geojson.Geometry, idx int) (int, error) {
		return acc + 1, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("expected 1, got %d", count)
	}
}

func TestGeomReduceNilGeometry(t *testing.T) {
	f := geojson.NewFeature(nil, nil)
	count, err := GeomReduce(f, 0, func(acc int, g geojson.Geometry, idx int) (int, error) {
		return acc + 1, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Errorf("expected 0 (nil geometry), got %d", count)
	}
}

func TestGeomReduceErr(t *testing.T) {
	_, err := GeomReduce("bad", 0, func(acc int, g geojson.Geometry, idx int) (int, error) {
		return acc, nil
	})
	if err == nil {
		t.Error("expected error for invalid type")
	}
}

func TestGeomReduceFCEmpty(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{})
	count, err := GeomReduce(fc, 0, func(acc int, g geojson.Geometry, idx int) (int, error) {
		return acc + 1, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Errorf("expected 0, got %d", count)
	}
}

func TestGeomReduceFCNilFeature(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(nil, nil),
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
	})
	count, err := GeomReduce(fc, 0, func(acc int, g geojson.Geometry, idx int) (int, error) {
		return acc + 1, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("expected 1 (nil geom skipped), got %d", count)
	}
}

func TestPropReduceSingleFeature(t *testing.T) {
	f := geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), map[string]any{"k": "v"})
	count, err := PropReduce(f, 0, func(acc int, p map[string]any, idx int) (int, error) {
		return acc + len(p), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("expected 1 prop, got %d", count)
	}
}

func TestPropReduceNilProps(t *testing.T) {
	f := geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil)
	count, err := PropReduce(f, 0, func(acc int, p map[string]any, idx int) (int, error) {
		return acc + 1, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("expected 1 even with nil props, got %d", count)
	}
}

func TestPropReduceFCNilProps(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
	})
	count, err := PropReduce(fc, 0, func(acc int, p map[string]any, idx int) (int, error) {
		return acc + 1, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("expected 1, got %d", count)
	}
}

func TestPropReduceErr(t *testing.T) {
	_, err := PropReduce("bad", 0, func(acc int, p map[string]any, idx int) (int, error) {
		return acc, nil
	})
	if err == nil {
		t.Error("expected error")
	}
}

func TestFlattenReduceMultiLineString(t *testing.T) {
	mls := geojson.NewFeature(geojson.NewMultiLineString([][]geojson.Position{{{0, 0}, {1, 1}}, {{2, 2}, {3, 3}}}), nil)
	count, err := FlattenReduce(mls, 0, func(acc int, f *geojson.Feature, idx int) (int, error) {
		return acc + 1, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Errorf("expected 2 features from MLS, got %d", count)
	}
}

func TestFlattenReduceMultiPolygon(t *testing.T) {
	mp := geojson.NewFeature(geojson.NewMultiPolygon([][][]geojson.Position{{{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}}}}), nil)
	count, err := FlattenReduce(mp, 0, func(acc int, f *geojson.Feature, idx int) (int, error) {
		return acc + 1, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("expected 1, got %d", count)
	}
}

func TestFlattenReduceGeometryCollection(t *testing.T) {
	// GC in a Feature: Point carries through but accumulation for subsequent simple geometries has a known limitation
	gc := geojson.NewFeature(geojson.NewGeometryCollection([]geojson.Geometry{
		geojson.NewMultiPoint([]geojson.Position{{0, 0}, {1, 1}}),
	}), nil)
	count, err := FlattenReduce(gc, 0, func(acc int, f *geojson.Feature, idx int) (int, error) {
		return acc + 1, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Errorf("expected 2, got %d", count)
	}
}

func TestFlattenReduceSingleGeometry(t *testing.T) {
	p := geojson.NewPoint(geojson.Position{0, 0})
	count, err := FlattenReduce(p, 0, func(acc int, f *geojson.Feature, idx int) (int, error) {
		return acc + 1, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("expected 1, got %d", count)
	}
}

func TestFlattenReduceSingleFeature(t *testing.T) {
	f := geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), map[string]any{"k": "v"})
	count, err := FlattenReduce(f, 0, func(acc int, f2 *geojson.Feature, idx int) (int, error) {
		if _, ok := f2.Properties["k"]; !ok {
			t.Error("expected property to carry through")
		}
		return acc + 1, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("expected 1, got %d", count)
	}
}

func TestFlattenReduceErr(t *testing.T) {
	_, err := FlattenReduce("bad", 0, func(acc int, f *geojson.Feature, idx int) (int, error) {
		return acc, nil
	})
	if err == nil {
		t.Error("expected error")
	}
}

func TestFlattenReduceNilGeom(t *testing.T) {
	f := geojson.NewFeature(nil, nil)
	count, err := FlattenReduce(f, 0, func(acc int, f2 *geojson.Feature, idx int) (int, error) {
		return acc + 1, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Errorf("expected 0 for nil geometry, got %d", count)
	}
}

func TestGeomCount(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  int
	}{
		{
			name: "feature collection with mixed geometry types",
			input: geojson.NewFeatureCollection([]*geojson.Feature{
				geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
				geojson.NewFeature(geojson.NewLineString([]geojson.Position{{0, 0}, {1, 1}}), nil),
				geojson.NewFeature(geojson.NewPolygon([][]geojson.Position{{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}}}), nil),
			}),
			want: 3,
		},
		{
			name:  "single feature with point geometry",
			input: geojson.NewFeature(geojson.NewPoint(geojson.Position{1, 2}), nil),
			want:  1,
		},
		{
			name:  "single geometry (point directly)",
			input: geojson.NewPoint(geojson.Position{3, 4}),
			want:  1,
		},
		{
			name: "feature collection with empty slice",
			input: geojson.NewFeatureCollection([]*geojson.Feature{}),
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GeomCount(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("GeomCount() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestFeatureCount(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  int
	}{
		{
			name: "feature collection with multiple features",
			input: geojson.NewFeatureCollection([]*geojson.Feature{
				geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
				geojson.NewFeature(geojson.NewPoint(geojson.Position{1, 1}), nil),
				geojson.NewFeature(geojson.NewPoint(geojson.Position{2, 2}), nil),
			}),
			want: 3,
		},
		{
			name:  "single feature",
			input: geojson.NewFeature(geojson.NewPoint(geojson.Position{1, 2}), nil),
			want:  1,
		},
		{
			name: "feature collection with single feature",
			input: geojson.NewFeatureCollection([]*geojson.Feature{
				geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil),
			}),
			want: 1,
		},
		{
			name: "feature collection with empty slice",
			input: geojson.NewFeatureCollection([]*geojson.Feature{}),
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FeatureCount(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("FeatureCount() = %d, want %d", got, tt.want)
			}
		})
	}
}
