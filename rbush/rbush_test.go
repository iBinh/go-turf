package rbush

import (
	"testing"

	"github.com/ibinh/turf-go/geojson"
)

func pointFeature(x, y float64) *geojson.Feature {
	p := geojson.NewPolygon([][]geojson.Position{
		{{x, y}, {x + 1, y}, {x + 1, y + 1}, {x, y + 1}, {x, y}},
	})
	return geojson.NewFeature(p, nil)
}

func TestRBushInsertSearch(t *testing.T) {
	r := NewRBush()
	f := pointFeature(0, 0)
	r.Insert(f)

	results := r.Search([]float64{-1, -1, 2, 2})
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

func TestRBushSearchNoMatch(t *testing.T) {
	r := NewRBush()
	r.Insert(pointFeature(0, 0))

	results := r.Search([]float64{10, 10, 20, 20})
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestRBushRemove(t *testing.T) {
	r := NewRBush()
	f1 := pointFeature(0, 0)
	f2 := pointFeature(10, 10)
	r.Insert(f1)
	r.Insert(f2)
	r.Remove(f1)

	results := r.Search([]float64{-1, -1, 2, 2})
	if len(results) != 0 {
		t.Errorf("expected 0 results after remove, got %d", len(results))
	}
	if len(r.All()) != 1 {
		t.Errorf("expected 1 item remaining, got %d", len(r.All()))
	}
}

func TestRBushNearest(t *testing.T) {
	r := NewRBush()
	r.Insert(pointFeature(0, 0))
	r.Insert(pointFeature(10, 0))
	r.Insert(pointFeature(0, 10))
	r.Insert(pointFeature(100, 100))

	nearest := r.Nearest(geojson.Position{0, 0}, 2)
	if len(nearest) != 2 {
		t.Errorf("expected 2 nearest, got %d", len(nearest))
	}
}

func TestRBushLoad(t *testing.T) {
	r := NewRBush()
	features := []*geojson.Feature{
		pointFeature(0, 0),
		pointFeature(1, 1),
		pointFeature(2, 2),
	}
	r.Load(features)

	if len(r.All()) != 3 {
		t.Errorf("expected 3 items, got %d", len(r.All()))
	}
}

func TestRBushClear(t *testing.T) {
	r := NewRBush()
	r.Insert(pointFeature(0, 0))
	r.Clear()

	if len(r.All()) != 0 {
		t.Errorf("expected 0 items after clear, got %d", len(r.All()))
	}
}

func TestRBushInsertNil(t *testing.T) {
	r := NewRBush()
	r.Insert(nil)
	if len(r.All()) != 0 {
		t.Errorf("expected 0 items after nil insert, got %d", len(r.All()))
	}
}

func TestRBushCollisionsQuery(t *testing.T) {
	r := NewRBush()
	r.Insert(pointFeature(0, 0))
	r.Insert(pointFeature(10, 10))

	results := r.CollisionsQuery(geojson.Position{0, 0}, 0.5)
	if len(results) != 1 {
		t.Errorf("expected 1 collision, got %d", len(results))
	}
}

func TestRBushSearchInvalidBBox(t *testing.T) {
	r := NewRBush()
	r.Insert(pointFeature(0, 0))

	results := r.Search([]float64{1, 2})
	if len(results) != 0 {
		t.Errorf("expected 0 results for invalid bbox, got %d", len(results))
	}
}
