package concave

import (
	"fmt"
	"math"

	"github.com/ibinh/turf-go/geojson"
	"github.com/ibinh/turf-go/interpolation"
)

type point struct {
	x, y float64
}

type edgeKey struct {
	a, b point
}

type edge struct {
	a, b  point
	count int
}

type ConcaveOptions struct {
	MaxEdge float64
}

func ConcaveHull(points *geojson.FeatureCollection, options ...ConcaveOptions) (*geojson.Feature, error) {
	opts := ConcaveOptions{MaxEdge: math.MaxFloat64}
	if len(options) > 0 {
		opts = options[0]
	}

	if points == nil || len(points.Features) < 3 {
		return nil, fmt.Errorf("at least 3 points required")
	}

	tin, err := interpolation.Tin(points)
	if err != nil {
		return nil, err
	}

	edgeMap := map[edgeKey]*edge{}
	makeKey := func(a, b point) edgeKey {
		if a.x < b.x || (a.x == b.x && a.y < b.y) {
			return edgeKey{a, b}
		}
		return edgeKey{b, a}
	}

	posToPoint := func(p geojson.Position) point {
		return point{p[0], p[1]}
	}

	for _, f := range tin.Features {
		poly, ok := f.Geometry.(*geojson.Polygon)
		if !ok || len(poly.Coordinates) == 0 {
			continue
		}
		ring := poly.Coordinates[0]
		if len(ring) < 4 {
			continue
		}
		for i := 0; i < 3; i++ {
			a, b := ring[i], ring[(i+1)%3]
			if opts.MaxEdge < math.MaxFloat64 {
				dx := a[0] - b[0]
				dy := a[1] - b[1]
				planarDist := math.Sqrt(dx*dx + dy*dy)
				if planarDist > opts.MaxEdge {
					continue
				}
			}
			pa, pb := posToPoint(a), posToPoint(b)
			key := makeKey(pa, pb)
			if e, exists := edgeMap[key]; exists {
				e.count++
			} else {
				edgeMap[key] = &edge{a: pa, b: pb, count: 1}
			}
		}
	}

	var boundaryEdges []struct{ a, b point }
	for _, e := range edgeMap {
		if e.count == 1 {
			boundaryEdges = append(boundaryEdges, struct{ a, b point }{e.a, e.b})
		}
	}

	if len(boundaryEdges) < 3 {
		return nil, fmt.Errorf("concave hull: insufficient boundary edges")
	}

	adj := map[point][]point{}
	for _, e := range boundaryEdges {
		adj[e.a] = append(adj[e.a], e.b)
		adj[e.b] = append(adj[e.b], e.a)
	}

	var ring []geojson.Position
	start := boundaryEdges[0].a
	cur := start
	visited := map[point]bool{}
	limit := len(boundaryEdges) * 2

	for len(ring) < limit {
		if visited[cur] {
			break
		}
		visited[cur] = true
		ring = append(ring, geojson.Position{cur.x, cur.y})
		neighbors := adj[cur]
		var next point
		found := false
		for _, nb := range neighbors {
			if !visited[nb] {
				next = nb
				found = true
				break
			}
		}
		if !found {
			break
		}
		cur = next
	}

	if len(ring) < 3 {
		return nil, fmt.Errorf("concave hull: failed to form ring")
	}

	ring = append(ring, ring[0])

	return geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{ring}),
		nil,
	), nil
}
