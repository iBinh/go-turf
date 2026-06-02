package misc

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/ibinh/turf-go/boolean"
	"github.com/ibinh/turf-go/geojson"
	"github.com/ibinh/turf-go/meta"
)

func Clone(geom any) (*geojson.Feature, error) {
	f, err := asFeature(geom)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}
	var cloned geojson.Feature
	if err := json.Unmarshal(data, &cloned); err != nil {
		return nil, err
	}
	return &cloned, nil
}

func Combine(fc *geojson.FeatureCollection) (*geojson.FeatureCollection, error) {
	if fc == nil {
		return nil, fmt.Errorf("combine: feature collection is nil")
	}
	groups := make(map[string][]*geojson.Feature)
	for _, f := range fc.Features {
		if f == nil || f.Geometry == nil {
			continue
		}
		t := f.Geometry.Type()
		groups[t] = append(groups[t], f)
	}
	var out []*geojson.Feature
	for t, features := range groups {
		switch t {
		case geojson.TypePoint:
			var pts []geojson.Position
			for _, f := range features {
				p := f.Geometry.(*geojson.Point)
				pts = append(pts, p.Coordinates)
			}
			out = append(out, geojson.NewFeature(geojson.NewMultiPoint(pts), nil))
		case geojson.TypeLineString:
			var lines [][]geojson.Position
			for _, f := range features {
				l := f.Geometry.(*geojson.LineString)
				lines = append(lines, l.Coordinates)
			}
			out = append(out, geojson.NewFeature(geojson.NewMultiLineString(lines), nil))
		case geojson.TypePolygon:
			var polys [][][]geojson.Position
			for _, f := range features {
				p := f.Geometry.(*geojson.Polygon)
				polys = append(polys, p.Coordinates)
			}
			out = append(out, geojson.NewFeature(geojson.NewMultiPolygon(polys), nil))
		default:
			for _, f := range features {
				out = append(out, f)
			}
		}
	}
	return geojson.NewFeatureCollection(out), nil
}

func Explode(geom any) (*geojson.FeatureCollection, error) {
	var features []*geojson.Feature
	addFn := func(f *geojson.Feature, _ int) error {
		features = append(features, f)
		return nil
	}
	if err := meta.FlattenEach(geom, addFn); err != nil {
		return nil, err
	}
	return geojson.NewFeatureCollection(features), nil
}

func PointsWithinPolygon(points, polygon any) (*geojson.FeatureCollection, error) {
	var features []*geojson.Feature
	addFn := func(f *geojson.Feature, _ int) error {
		if f == nil || f.Geometry == nil {
			return nil
		}
		inside, err := boolean.PointInPolygon(f, polygon)
		if err != nil {
			return err
		}
		if inside {
			features = append(features, f)
		}
		return nil
	}
	if err := meta.FeatureEach(points, addFn); err != nil {
		return nil, err
	}
	return geojson.NewFeatureCollection(features), nil
}

func Planepoint(point any, triangle any) (float64, error) {
	pt, err := geojson.GetCoord(point)
	if err != nil {
		return 0, err
	}
	geom, err := geojson.GetGeometry(triangle)
	if err != nil {
		return 0, err
	}
	poly, ok := geom.(*geojson.Polygon)
	if !ok {
		return 0, fmt.Errorf("planepoint: expected Polygon geometry")
	}
	ring := poly.Coordinates[0]
	if len(ring) < 4 {
		return 0, fmt.Errorf("planepoint: polygon must have at least 3 vertices")
	}
	a, b, c := ring[0], ring[1], ring[2]

	v0 := geojson.Position{c[0] - a[0], c[1] - a[1]}
	v1 := geojson.Position{b[0] - a[0], b[1] - a[1]}
	v2 := geojson.Position{pt[0] - a[0], pt[1] - a[1]}

	dot00 := v0[0]*v0[0] + v0[1]*v0[1]
	dot01 := v0[0]*v1[0] + v0[1]*v1[1]
	dot02 := v0[0]*v2[0] + v0[1]*v2[1]
	dot11 := v1[0]*v1[0] + v1[1]*v1[1]
	dot12 := v1[0]*v2[0] + v1[1]*v2[1]

	denom := dot00*dot11 - dot01*dot01
	if math.Abs(denom) < 1e-15 {
		return 0, fmt.Errorf("planepoint: degenerate triangle")
	}

	u := (dot11*dot02 - dot01*dot12) / denom
	v := (dot00*dot12 - dot01*dot02) / denom
	w := 1 - u - v

	zA := 0.0
	zB := 0.0
	zC := 0.0
	if len(a) > 2 {
		zA = a[2]
	}
	if len(b) > 2 {
		zB = b[2]
	}
	if len(c) > 2 {
		zC = c[2]
	}

	return w*zA + u*zB + v*zC, nil
}

func Tesselate(poly any) (*geojson.FeatureCollection, error) {
	geom, err := geojson.GetGeometry(poly)
	if err != nil {
		return nil, err
	}
	var rings [][]geojson.Position
	switch g := geom.(type) {
	case *geojson.Polygon:
		rings = g.Coordinates
	default:
		return nil, fmt.Errorf("tesselate: expected Polygon geometry")
	}
	cw := isClockwise(rings[0])
	triangles := earClip(rings[0], cw)
	var features []*geojson.Feature
	for _, tri := range triangles {
		features = append(features, geojson.NewFeature(geojson.NewPolygon([][]geojson.Position{tri}), nil))
	}
	return geojson.NewFeatureCollection(features), nil
}

func isClockwise(ring []geojson.Position) bool {
	area := 0.0
	for i := 0; i < len(ring)-1; i++ {
		area += (ring[i+1][0] - ring[i][0]) * (ring[i+1][1] + ring[i][1])
	}
	return area > 0
}

func earClip(ring []geojson.Position, cw bool) [][]geojson.Position {
	n := len(ring) - 1
	if n < 3 {
		return nil
	}
	pts := make([]geojson.Position, n)
	copy(pts, ring[:n])

	type node struct {
		idx      int
		prev, next int
	}
	nodes := make([]node, n)
	for i := 0; i < n; i++ {
		nodes[i] = node{idx: i, prev: (i - 1 + n) % n, next: (i + 1) % n}
	}

	var triangles [][]geojson.Position
	remaining := n

	for remaining > 3 {
		found := false
		for i := 0; i < n; i++ {
			ni := nodes[i]
			if ni.prev == -1 {
				continue
			}
			p := pts[nodes[ni.prev].idx]
			q := pts[ni.idx]
			r := pts[nodes[ni.next].idx]

			cross := (q[0]-p[0])*(r[1]-q[1]) - (q[1]-p[1])*(r[0]-q[0])
			isConvex := false
			if cw {
				isConvex = cross <= 1e-10
			} else {
				isConvex = cross >= -1e-10
			}
			if !isConvex {
				continue
			}

			inside := false
			for j := 0; j < n; j++ {
				if j == nodes[ni.prev].idx || j == ni.idx || j == nodes[ni.next].idx || nodes[j].prev == -1 {
					continue
				}
				if pointInTriangle(pts[j], p, q, r) {
					inside = true
					break
				}
			}
			if inside {
				continue
			}

			triangles = append(triangles, []geojson.Position{p, q, r, p})
			nodes[ni.prev].next = ni.next
			nodes[ni.next].prev = ni.prev
			nodes[i].prev = -1
			nodes[i].next = -1
			remaining--
			found = true
			break
		}
		if !found {
			break
		}
	}

	var last []int
	for i := 0; i < n; i++ {
		if nodes[i].prev != -1 {
			last = append(last, nodes[i].idx)
		}
	}
	if len(last) == 3 {
		triangles = append(triangles, []geojson.Position{
			pts[last[0]], pts[last[1]], pts[last[2]], pts[last[0]],
		})
	}

	return triangles
}

func pointInTriangle(p, a, b, c geojson.Position) bool {
	d1 := sign(p, a, b)
	d2 := sign(p, b, c)
	d3 := sign(p, c, a)
	hasNeg := (d1 < -1e-10) || (d2 < -1e-10) || (d3 < -1e-10)
	hasPos := (d1 > 1e-10) || (d2 > 1e-10) || (d3 > 1e-10)
	return !(hasNeg && hasPos)
}

func sign(p1, p2, p3 geojson.Position) float64 {
	return (p1[0]-p3[0])*(p2[1]-p3[1]) - (p2[0]-p3[0])*(p1[1]-p3[1])
}

func Flatten(geom any) (*geojson.FeatureCollection, error) {
	var features []*geojson.Feature
	addFn := func(f *geojson.Feature, _ int) error {
		features = append(features, f)
		return nil
	}
	if err := meta.FlattenEach(geom, addFn); err != nil {
		return nil, err
	}
	return geojson.NewFeatureCollection(features), nil
}

func asFeature(geom any) (*geojson.Feature, error) {
	switch v := geom.(type) {
	case *geojson.Feature:
		return v, nil
	case geojson.Geometry:
		return geojson.NewFeature(v, nil), nil
	default:
		return nil, fmt.Errorf("expected Feature or Geometry, got %T", geom)
	}
}
