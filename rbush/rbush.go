package rbush

import (
	"math"
	"sort"

	"github.com/ibinh/turf-go/geojson"
)

type RBush struct {
	items []*entry
}

type entry struct {
	bbox    []float64
	feature *geojson.Feature
}

func NewRBush() *RBush {
	return &RBush{}
}

func (r *RBush) Insert(feature *geojson.Feature) {
	if feature == nil {
		return
	}
	bb := feature.BBox()
	if bb == nil {
		bb = computeBBox(feature)
		if bb == nil {
			return
		}
	}
	r.items = append(r.items, &entry{bbox: bb, feature: feature})
}

func (r *RBush) Load(features []*geojson.Feature) {
	for _, f := range features {
		r.Insert(f)
	}
}

func (r *RBush) Remove(feature *geojson.Feature) {
	for i, e := range r.items {
		if e.feature == feature {
			r.items = append(r.items[:i], r.items[i+1:]...)
			return
		}
	}
}

func (r *RBush) Search(bbox []float64) []*geojson.Feature {
	if len(bbox) < 4 {
		return nil
	}
	var result []*geojson.Feature
	for _, e := range r.items {
		if bboxesOverlap(e.bbox, bbox) {
			result = append(result, e.feature)
		}
	}
	return result
}

func (r *RBush) CollisionsQuery(point geojson.Position, maxDist float64) []*geojson.Feature {
	var result []*geojson.Feature
	for _, e := range r.items {
		d := pointBBoxDist(point, e.bbox)
		if d <= maxDist {
			result = append(result, e.feature)
		}
	}
	return result
}

func (r *RBush) Nearest(point geojson.Position, n int) []*geojson.Feature {
	if n <= 0 {
		n = 1
	}
	type distEntry struct {
		dist    float64
		feature *geojson.Feature
	}
	var entries []distEntry
	for _, e := range r.items {
		d := pointBBoxDist(point, e.bbox)
		entries = append(entries, distEntry{dist: d, feature: e.feature})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].dist < entries[j].dist
	})
	if n > len(entries) {
		n = len(entries)
	}
	result := make([]*geojson.Feature, n)
	for i := 0; i < n; i++ {
		result[i] = entries[i].feature
	}
	return result
}

func (r *RBush) All() []*geojson.Feature {
	result := make([]*geojson.Feature, len(r.items))
	for i, e := range r.items {
		result[i] = e.feature
	}
	return result
}

func (r *RBush) Clear() {
	r.items = nil
}

func computeBBox(f *geojson.Feature) []float64 {
	coords, err := geojson.CoordAll(f)
	if err != nil || len(coords) == 0 {
		return nil
	}
	minX, minY := coords[0][0], coords[0][1]
	maxX, maxY := minX, minY
	for _, c := range coords {
		if c[0] < minX {
			minX = c[0]
		}
		if c[0] > maxX {
			maxX = c[0]
		}
		if c[1] < minY {
			minY = c[1]
		}
		if c[1] > maxY {
			maxY = c[1]
		}
	}
	return []float64{minX, minY, maxX, maxY}
}

func bboxesOverlap(a, b []float64) bool {
	return a[0] <= b[2] && a[2] >= b[0] && a[1] <= b[3] && a[3] >= b[1]
}

func pointBBoxDist(p geojson.Position, bbox []float64) float64 {
	dx := math.Max(bbox[0]-p[0], math.Max(0, p[0]-bbox[2]))
	dy := math.Max(bbox[1]-p[1], math.Max(0, p[1]-bbox[3]))
	return math.Sqrt(dx*dx + dy*dy)
}
