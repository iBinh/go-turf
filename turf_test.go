package turf

import (
	"encoding/json"
	"testing"

	"github.com/ibinh/turf-go/geojson"
)

func TestUmbrellaPoint(t *testing.T) {
	f := Point(geojson.Position{1, 2}, map[string]any{"a": 1})
	data, err := json.Marshal(f)
	if err != nil {
		t.Fatal(err)
	}
	var g geojson.Feature
	if err := json.Unmarshal(data, &g); err != nil {
		t.Fatal(err)
	}
	coord, err := GetCoord(&g)
	if err != nil {
		t.Fatal(err)
	}
	if coord[0] != 1 || coord[1] != 2 {
		t.Errorf("unexpected coord: %v", coord)
	}
}

func TestUmbrellaLineString(t *testing.T) {
	f := LineString([]geojson.Position{{0, 0}, {1, 1}, {2, 2}}, nil)
	coords := make([]geojson.Position, 0)
	err := CoordEach(f, func(coord geojson.Position, idx int) error {
		coords = append(coords, coord)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(coords) != 3 {
		t.Errorf("expected 3 coords, got %d", len(coords))
	}
}

func TestUmbrellaFeatureCollection(t *testing.T) {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{
		Point(geojson.Position{0, 0}, map[string]any{"v": float64(1)}),
		Point(geojson.Position{1, 1}, map[string]any{"v": float64(2)}),
	})
	sum, err := FeatureReduce(fc, 0, func(acc int, f *geojson.Feature, idx int) (int, error) {
		v, _ := f.Properties["v"].(float64)
		return acc + int(v), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if sum != 3 {
		t.Errorf("expected 3, got %d", sum)
	}
}

func TestUmbrellaTypes(t *testing.T) {
	if TypePoint != "Point" {
		t.Error("TypePoint mismatch")
	}
	if TypeFeatureCollection != "FeatureCollection" {
		t.Error("TypeFeatureCollection mismatch")
	}
}

func TestUmbrellaWithOptions(t *testing.T) {
	f := Point(geojson.Position{1, 2}, nil, WithBBox([]float64{0, 0, 2, 2}), WithID("x1"))
	if f.BBox()[0] != 0 {
		t.Errorf("bbox missing")
	}
	if f.ID != "x1" {
		t.Errorf("expected id x1, got %v", f.ID)
	}
}

func TestUmbrellaCoords(t *testing.T) {
	p := Point(geojson.Position{1, 2}, nil)
	coord, err := GetCoord(p)
	if err != nil {
		t.Fatal(err)
	}
	if coord[0] != 1 || coord[1] != 2 {
		t.Errorf("unexpected coord: %v", coord)
	}
}
