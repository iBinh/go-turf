package unkink

import (
	"fmt"
	"math"
	"sort"

	"github.com/ibinh/turf-go/geojson"
)

func UnkinkPolygon(poly any) (*geojson.FeatureCollection, error) {
	g, err := geojson.GetGeometry(poly)
	if err != nil {
		return nil, fmt.Errorf("unkinkPolygon: %w", err)
	}

	var allRings [][]geojson.Position
	switch v := g.(type) {
	case *geojson.Polygon:
		if len(v.Coordinates) == 0 {
			return nil, fmt.Errorf("unkinkPolygon: polygon has no rings")
		}
		allRings = v.Coordinates
	case *geojson.MultiPolygon:
		for _, polyCoords := range v.Coordinates {
			allRings = append(allRings, polyCoords...)
		}
	default:
		return nil, fmt.Errorf("unkinkPolygon: expected Polygon or MultiPolygon, got %s", g.Type())
	}

	if len(allRings) == 0 {
		return nil, fmt.Errorf("unkinkPolygon: no rings found")
	}

	exterior := allRings[0]
	if len(exterior) < 4 {
		return geojson.NewFeatureCollection([]*geojson.Feature{
			geojson.NewFeature(geojson.NewPolygon(allRings), nil),
		}), nil
	}

	isects := findAllSelfIntersections(exterior)

	if len(isects) == 0 {
		return geojson.NewFeatureCollection([]*geojson.Feature{
			geojson.NewFeature(geojson.NewPolygon(allRings), nil),
		}), nil
	}

	polys := splitAtIntersections(exterior, isects)

	var features []*geojson.Feature
	for _, ring := range polys {
		if len(ring) >= 4 {
			features = append(features, geojson.NewFeature(geojson.NewPolygon([][]geojson.Position{ring}), nil))
		}
	}

	if len(features) == 0 {
		return geojson.NewFeatureCollection([]*geojson.Feature{
			geojson.NewFeature(geojson.NewPolygon(allRings), nil),
		}), nil
	}

	return geojson.NewFeatureCollection(features), nil
}

type isectPoint struct {
	point geojson.Position
	idxA  int
	tA    float64
	idxB  int
	tB    float64
}

func findAllSelfIntersections(ring []geojson.Position) []isectPoint {
	n := len(ring) - 1
	if n < 4 {
		return nil
	}

	var result []isectPoint
	for i := 0; i < n; i++ {
		a, b := ring[i], ring[(i+1)%n]
		for j := i + 2; j < n; j++ {
			if i == 0 && j == n-1 {
				continue
			}
			c, d := ring[j], ring[(j+1)%n]

			if shareVertex(a, b, c, d) {
				continue
			}

			pt, t, u, ok := segIntersect(a, b, c, d)
			if ok {
				result = append(result, isectPoint{
					point: pt,
					idxA:  i,
					tA:    t,
					idxB:  j,
					tB:    u,
				})
			}
		}
	}
	return result
}

func shareVertex(a, b, c, d geojson.Position) bool {
	return pointEqual(a, c) || pointEqual(a, d) || pointEqual(b, c) || pointEqual(b, d)
}

func pointEqual(a, b geojson.Position) bool {
	return math.Abs(a[0]-b[0]) < 1e-10 && math.Abs(a[1]-b[1]) < 1e-10
}

func segIntersect(a, b, c, d geojson.Position) (geojson.Position, float64, float64, bool) {
	den := (b[0]-a[0])*(d[1]-c[1]) - (b[1]-a[1])*(d[0]-c[0])
	if math.Abs(den) < 1e-15 {
		return geojson.Position{}, 0, 0, false
	}
	t := ((c[0]-a[0])*(d[1]-c[1]) - (c[1]-a[1])*(d[0]-c[0])) / den
	u := ((c[0]-a[0])*(b[1]-a[1]) - (c[1]-a[1])*(b[0]-a[0])) / den
	if t <= 0 || t >= 1 || u <= 0 || u >= 1 {
		return geojson.Position{}, 0, 0, false
	}
	return geojson.Position{
		a[0] + t*(b[0]-a[0]),
		a[1] + t*(b[1]-a[1]),
	}, t, u, true
}

type node struct {
	point geojson.Position
	next  *node
	prev  *node
	isect bool
	used  bool
}

type edgeIsect struct {
	param float64
	node  *node
}

func splitAtIntersections(ring []geojson.Position, isects []isectPoint) [][]geojson.Position {
	n := len(ring) - 1

	nodes := make([]*node, n)
	for i := 0; i < n; i++ {
		nodes[i] = &node{point: ring[i]}
	}
	for i := 0; i < n; i++ {
		nodes[i].next = nodes[(i+1)%n]
		nodes[i].prev = nodes[(i+n-1)%n]
	}

	edgeIsects := make([][]edgeIsect, n)
	for _, is := range isects {
		nA := &node{point: is.point, isect: true}
		nB := &node{point: is.point, isect: true}
		nA.next = nB
		nB.next = nA

		edgeIsects[is.idxA] = append(edgeIsects[is.idxA], edgeIsect{is.tA, nA})
		edgeIsects[is.idxB] = append(edgeIsects[is.idxB], edgeIsect{is.tB, nB})
	}

	for edgeIdx, lst := range edgeIsects {
		if len(lst) == 0 {
			continue
		}
		sort.Slice(lst, func(i, j int) bool { return lst[i].param < lst[j].param })
		cur := nodes[edgeIdx]
		for _, item := range lst {
			nn := item.node
			nxt := cur.next
			nn.next = nxt
			nn.prev = cur
			cur.next = nn
			nxt.prev = nn
			cur = nn
		}
	}

	var result [][]geojson.Position

	// Walk the graph to find all cycles
	start := nodes[0]
	cur := start
	for cur != nil {
		if cur.used {
			cur = cur.next
			if cur == start {
				break
			}
			continue
		}

		var ringPts []geojson.Position
		walk := cur
		for walk != nil && !walk.used {
			walk.used = true
			ringPts = append(ringPts, walk.point)
			walk = walk.next
			if walk == cur {
				break
			}
		}
		if len(ringPts) >= 3 {
			ringPts = append(ringPts, ringPts[0])
			result = append(result, ringPts)
		}

		cur = cur.next
		if cur == start {
			break
		}
	}

	if len(result) == 0 {
		result = append(result, ring)
	}

	return result
}


