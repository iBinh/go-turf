package polyclip

import (
	"fmt"
	"math"
	"sort"

	"github.com/ibinh/turf-go/geojson"
)

type OpType int

const (
	OpUnion OpType = iota
	OpIntersect
	OpDifference
	OpXor
)

func PolygonUnion(poly1, poly2 any) (*geojson.Feature, error) {
	return operate(poly1, poly2, OpUnion)
}

func PolygonIntersect(poly1, poly2 any) (*geojson.Feature, error) {
	return operate(poly1, poly2, OpIntersect)
}

func PolygonDifference(poly1, poly2 any) (*geojson.Feature, error) {
	return operate(poly1, poly2, OpDifference)
}

func PolygonXor(poly1, poly2 any) (*geojson.Feature, error) {
	return operate(poly1, poly2, OpXor)
}

func operate(poly1, poly2 any, op OpType) (*geojson.Feature, error) {
	rings1, err := getRings(poly1)
	if err != nil {
		return nil, err
	}
	rings2, err := getRings(poly2)
	if err != nil {
		return nil, err
	}
	if len(rings1) == 0 || len(rings2) == 0 {
		return nil, nil
	}

	polys1 := extractPolygons(rings1)
	polys2 := extractPolygons(rings2)

	var allResultRings [][][]geojson.Position
	for _, p1 := range polys1 {
		for _, p2 := range polys2 {
			result := clipPolygon(p1.exterior, p2.exterior, op)
			allResultRings = append(allResultRings, result...)
		}
	}

	if len(allResultRings) == 0 {
		return nil, nil
	}

	merged := mergeRings(allResultRings)
	if len(merged) == 0 {
		return nil, nil
	}
	if len(merged) == 1 {
		return geojson.NewFeature(geojson.NewPolygon(merged[0]), nil), nil
	}
	return geojson.NewFeature(geojson.NewMultiPolygon(merged), nil), nil
}

func getRings(geom any) ([][][]geojson.Position, error) {
	g, err := geojson.GetGeometry(geom)
	if err != nil {
		return nil, err
	}
	switch v := g.(type) {
	case *geojson.Polygon:
		return [][][]geojson.Position{v.Coordinates}, nil
	case *geojson.MultiPolygon:
		return v.Coordinates, nil
	}
	return nil, fmt.Errorf("unsupported geometry type for polygon operation")
}

func extractPolygons(rings [][][]geojson.Position) []*polygon {
	var polys []*polygon
	for _, ringset := range rings {
		if len(ringset) == 0 {
			continue
		}
		p := &polygon{exterior: ringset[0]}
		if len(ringset) > 1 {
			p.holes = ringset[1:]
		}
		polys = append(polys, p)
	}
	return polys
}

type polygon struct {
	exterior []geojson.Position
	holes    [][]geojson.Position
}

type ptNode struct {
	point geojson.Position
	next  *ptNode
	prev  *ptNode
	isect bool
	used  bool
	neigh *ptNode
}

type entry struct {
	onA     *ptNode
	onB     *ptNode
	isEntry bool
	used    bool
}

func clipPolygon(ring1, ring2 []geojson.Position, op OpType) [][][]geojson.Position {
	pts1 := ringToSlice(ring1)
	pts2 := ringToSlice(ring2)

	allIsects := findAllIntersections(pts1, pts2)
	var isects []isectInfo
	for _, is := range allIsects {
		if !is.collin {
			isects = append(isects, is)
		}
	}
	if len(isects) == 0 {
		return noIntersectionResult(pts1, pts2, op)
	}

	nodes1 := buildRing(pts1)
	nodes2 := buildRing(pts2)

	entries := insertIntersections(nodes1, nodes2, isects)

	markEntryExit(nodes1, entries)

	return buildResultRings(nodes1, nodes2, entries, op)
}

func ringToSlice(ring []geojson.Position) []geojson.Position {
	if len(ring) < 3 {
		return nil
	}
	n := len(ring)
	if ring[0][0] == ring[n-1][0] && ring[0][1] == ring[n-1][1] {
		ring = ring[:n-1]
	}
	if len(ring) < 3 {
		return nil
	}
	return ring
}

func noIntersectionResult(pts1, pts2 []geojson.Position, op OpType) [][][]geojson.Position {
	if len(pts1) < 3 || len(pts2) < 3 {
		return nil
	}

	// Check for identical rings
	if ringsEqual(pts1, pts2) {
		switch op {
		case OpUnion, OpIntersect:
			return [][][]geojson.Position{{closeRing(pts1)}}
		case OpDifference, OpXor:
			return nil
		}
	}

	p1In2 := pointInRing(pts1[0], pts2)
	p2In1 := pointInRing(pts2[0], pts1)

	switch op {
	case OpUnion:
		if p2In1 {
			return [][][]geojson.Position{{closeRing(pts1)}}
		}
		if p1In2 {
			return [][][]geojson.Position{{closeRing(pts2)}}
		}
		// Check if touching (share a boundary edge)
		if ringsTouching(pts1, pts2) {
			merged := mergeTouchingRings(pts1, pts2)
			if merged != nil {
				return [][][]geojson.Position{{merged}}
			}
		}
		return [][][]geojson.Position{{closeRing(pts1)}, {closeRing(pts2)}}
	case OpIntersect:
		if p2In1 {
			return [][][]geojson.Position{{closeRing(pts2)}}
		}
		if p1In2 {
			return [][][]geojson.Position{{closeRing(pts1)}}
		}
		return nil
	case OpDifference:
		if p2In1 {
			return [][][]geojson.Position{{closeRing(pts1)}, {closeRing(pts2)}}
		}
		if p1In2 {
			return nil
		}
		return [][][]geojson.Position{{closeRing(pts1)}}
	case OpXor:
		if p1In2 || p2In1 {
			return nil
		}
		return [][][]geojson.Position{{closeRing(pts1)}, {closeRing(pts2)}}
	}
	return nil
}

func ringsEqual(a, b []geojson.Position) bool {
	if len(a) != len(b) {
		return false
	}
	n := len(a)
	// Find matching start point
	start := -1
	for i := 0; i < n; i++ {
		if pointEqual(a[0], b[i]) {
			start = i
			break
		}
	}
	if start < 0 {
		return false
	}
	for i := 0; i < n; i++ {
		if !pointEqual(a[i], b[(start+i)%n]) {
			return false
		}
	}
	return true
}

func ringsTouching(a, b []geojson.Position) bool {
	na := len(a)
	nb := len(b)
	for i := 0; i < na; i++ {
		a1, a2 := a[i], a[(i+1)%na]
		for j := 0; j < nb; j++ {
			b1, b2 := b[j], b[(j+1)%nb]
			if edgesCollinear(a1, a2, b1, b2) && edgesOverlap(a1, a2, b1, b2) {
				return true
			}
		}
	}
	return false
}

func edgesCollinear(a, b, c, d geojson.Position) bool {
	abx := b[0] - a[0]
	aby := b[1] - a[1]
	cross1 := abx*(c[1]-a[1]) - aby*(c[0]-a[0])
	cross2 := abx*(d[1]-a[1]) - aby*(d[0]-a[0])
	return math.Abs(cross1) < 1e-10 && math.Abs(cross2) < 1e-10
}

func edgesOverlap(a, b, c, d geojson.Position) bool {
	abx := b[0] - a[0]
	aby := b[1] - a[1]
	dot := abx*abx + aby*aby
	if dot < 1e-15 {
		return false
	}
	tc := ((c[0]-a[0])*abx + (c[1]-a[1])*aby) / dot
	td := ((d[0]-a[0])*abx + (d[1]-a[1])*aby) / dot
	if tc > td {
		tc, td = td, tc
	}
	overlapStart := math.Max(tc, 0.0)
	overlapEnd := math.Min(td, 1.0)
	return overlapStart < overlapEnd-1e-12
}

func mergeTouchingRings(a, b []geojson.Position) []geojson.Position {
	na := len(a)
	nb := len(b)
	// Find a shared edge pair
	for i := 0; i < na; i++ {
		a1, a2 := a[i], a[(i+1)%na]
		for j := 0; j < nb; j++ {
			b1, b2 := b[j], b[(j+1)%nb]
			if edgesCollinear(a1, a2, b1, b2) && edgesOverlap(a1, a2, b1, b2) {
				// Merge by joining the two rings
				// Walk from a1 to a2, then from b1 to a1 (the other way)
				return buildMergedRing(a, i, b, j)
			}
		}
	}
	return nil
}

func buildMergedRing(a []geojson.Position, edgeA int, b []geojson.Position, edgeB int) []geojson.Position {
	na := len(a)
	nb := len(b)

	// The shared edge is a[edgeA]→a[edgeA+1] and b[edgeB]→b[edgeB+1]
	// We want to walk from a[edgeA+1] along a to a[edgeA], then
	// from b[edgeB+1] along b to b[edgeB] (skipping the shared edge twice)

	// Walk a forward from edgeA+1 back to edgeA
	var ring []geojson.Position
	idx := (edgeA + 1) % na
	for idx != edgeA {
		ring = append(ring, a[idx])
		idx = (idx + 1) % na
	}
	ring = append(ring, a[edgeA])

	// Walk b forward from (edgeB+1) back to edgeB
	idx = (edgeB + 1) % nb
	for idx != edgeB {
		ring = append(ring, b[idx])
		idx = (idx + 1) % nb
	}
	ring = append(ring, b[edgeB])

	if len(ring) < 3 {
		return nil
	}
	return ring
}

func closeRing(pts []geojson.Position) []geojson.Position {
	if len(pts) < 2 {
		return pts
	}
	last := pts[len(pts)-1]
	first := pts[0]
	if last[0] == first[0] && last[1] == first[1] {
		return pts
	}
	return append(pts, first)
}

func pointInRing(pt geojson.Position, ring []geojson.Position) bool {
	inside := false
	n := len(ring)
	j := n - 1
	for i := 0; i < n; i++ {
		if ((ring[i][1] > pt[1]) != (ring[j][1] > pt[1])) &&
			(pt[0] < (ring[j][0]-ring[i][0])*(pt[1]-ring[i][1])/(ring[j][1]-ring[i][1])+ring[i][0]) {
			inside = !inside
		}
		j = i
	}
	return inside
}

func pointEqual(a, b geojson.Position) bool {
	return math.Abs(a[0]-b[0]) < 1e-10 && math.Abs(a[1]-b[1]) < 1e-10
}

type isectInfo struct {
	pt      geojson.Position
	idxA    int
	paramA  float64
	idxB    int
	paramB  float64
	collin  bool
}

func findAllIntersections(ptsA, ptsB []geojson.Position) []isectInfo {
	var result []isectInfo
	nA := len(ptsA)
	nB := len(ptsB)
	for i := 0; i < nA; i++ {
		a, b := ptsA[i], ptsA[(i+1)%nA]
		for j := 0; j < nB; j++ {
			c, d := ptsB[j], ptsB[(j+1)%nB]

			// Check collinear overlap
			if co := collinearOverlap(a, b, c, d); len(co) > 0 {
				for _, ci := range co {
					result = append(result, ci)
				}
				continue
			}

			pt, t, u, ok := segSegIntersect(a, b, c, d)
			if ok {
				// Skip only if the intersection is at a shared vertex of both edges
				if !(pointEqual(pt, a) && pointEqual(pt, c)) &&
					!(pointEqual(pt, a) && pointEqual(pt, d)) &&
					!(pointEqual(pt, b) && pointEqual(pt, c)) &&
					!(pointEqual(pt, b) && pointEqual(pt, d)) {
					result = append(result, isectInfo{pt: pt, idxA: i, paramA: t, idxB: j, paramB: u})
				}
			}
		}
	}
	return result
}

func collinearOverlap(a, b, c, d geojson.Position) []isectInfo {
	abx := b[0] - a[0]
	aby := b[1] - a[1]
	acx := c[0] - a[0]
	acy := c[1] - a[1]
	adx := d[0] - a[0]
	ady := d[1] - a[1]

	cross1 := abx*acy - aby*acx
	cross2 := abx*ady - aby*adx
	if math.Abs(cross1) > 1e-10 || math.Abs(cross2) > 1e-10 {
		return nil
	}

	if math.Abs(abx) < 1e-10 && math.Abs(aby) < 1e-10 {
		return nil
	}

	dot := abx*abx + aby*aby

	tc := (acx*abx + acy*aby) / dot
	td := (adx*abx + ady*aby) / dot

	if tc > td {
		tc, td = td, tc
	}

	overlapStart := math.Max(tc, 0.0)
	overlapEnd := math.Min(td, 1.0)

	if overlapStart >= overlapEnd-1e-12 {
		return nil
	}

	// For determining which edge of A this overlap belongs to:
	// The overlap is on A's current edge, idxA and idxB are known from caller
	// But since we iterate in findAllIntersections, we can use the provided indices
	// Return the midpoint as a usable intersection point
	mid := (overlapStart + overlapEnd) / 2
	midPt := geojson.Position{a[0] + mid*abx, a[1] + mid*aby}

	return []isectInfo{{
		pt:     midPt,
		idxA:   -1,
		paramA: 0,
		idxB:   -1,
		paramB: 0,
		collin: true,
	}}
}

func segSegIntersect(a, b, c, d geojson.Position) (geojson.Position, float64, float64, bool) {
	den := (b[0]-a[0])*(d[1]-c[1]) - (b[1]-a[1])*(d[0]-c[0])
	if math.Abs(den) < 1e-15 {
		return geojson.Position{}, 0, 0, false
	}
	t := ((c[0]-a[0])*(d[1]-c[1]) - (c[1]-a[1])*(d[0]-c[0])) / den
	u := ((c[0]-a[0])*(b[1]-a[1]) - (c[1]-a[1])*(b[0]-a[0])) / den
	if t < -1e-12 || t > 1+1e-12 || u < -1e-12 || u > 1+1e-12 {
		return geojson.Position{}, 0, 0, false
	}
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	if u < 0 {
		u = 0
	}
	if u > 1 {
		u = 1
	}
	return geojson.Position{
		a[0] + t*(b[0]-a[0]),
		a[1] + t*(b[1]-a[1]),
	}, t, u, true
}

func buildRing(pts []geojson.Position) []*ptNode {
	n := len(pts)
	nodes := make([]*ptNode, n)
	for i := 0; i < n; i++ {
		nodes[i] = &ptNode{point: pts[i]}
	}
	for i := 0; i < n; i++ {
		nodes[i].next = nodes[(i+1)%n]
		nodes[i].prev = nodes[(i+n-1)%n]
	}
	return nodes
}

func insertIntersections(nodesA, nodesB []*ptNode, isects []isectInfo) []*entry {
	aEdgeIsects := make([][]struct {
		param float64
		node  *ptNode
	}, len(nodesA))
	bEdgeIsects := make([][]struct {
		param float64
		node  *ptNode
	}, len(nodesB))

	var entries []*entry

	for _, is := range isects {
		if is.collin {
			continue
		}
		nA := &ptNode{point: is.pt, isect: true}
		nB := &ptNode{point: is.pt, isect: true}
		nA.neigh = nB
		nB.neigh = nA
		entries = append(entries, &entry{onA: nA, onB: nB})
		aEdgeIsects[is.idxA] = append(aEdgeIsects[is.idxA], struct {
			param float64
			node  *ptNode
		}{is.paramA, nA})
		bEdgeIsects[is.idxB] = append(bEdgeIsects[is.idxB], struct {
			param float64
			node  *ptNode
		}{is.paramB, nB})
	}

	for edgeIdx, lst := range aEdgeIsects {
		if len(lst) == 0 {
			continue
		}
		sort.Slice(lst, func(i, j int) bool { return lst[i].param < lst[j].param })
		cur := nodesA[edgeIdx]
		for _, item := range lst {
			n := item.node
			nxt := cur.next
			n.next = nxt
			n.prev = cur
			cur.next = n
			nxt.prev = n
			cur = n
		}
	}

	for edgeIdx, lst := range bEdgeIsects {
		if len(lst) == 0 {
			continue
		}
		sort.Slice(lst, func(i, j int) bool { return lst[i].param < lst[j].param })
		cur := nodesB[edgeIdx]
		for _, item := range lst {
			n := item.node
			nxt := cur.next
			n.next = nxt
			n.prev = cur
			cur.next = n
			nxt.prev = n
			cur = n
		}
	}

	return entries
}

func markEntryExit(nodesA []*ptNode, entries []*entry) {
	inside := false
	start := nodesA[0]
	cur := start
	for {
		if cur.isect {
			for _, e := range entries {
				if e.onA == cur {
					e.isEntry = !inside
					break
				}
			}
			inside = !inside
		}
		cur = cur.next
		if cur == start {
			break
		}
	}
}

func buildResultRings(nodesA, nodesB []*ptNode, entries []*entry, op OpType) [][][]geojson.Position {
	var rings [][][]geojson.Position

	for {
		var start *entry
		for _, e := range entries {
			if !e.used && shouldUseEntry(e, op) {
				start = e
				break
			}
		}
		if start == nil {
			break
		}

		result := traverseRing(start, entries, op)
		if len(result) > 2 {
			result = closeRing(result)
			rings = append(rings, [][]geojson.Position{result})
		}
	}

	return rings
}

func shouldUseEntry(e *entry, op OpType) bool {
	switch op {
	case OpIntersect:
		return e.isEntry
	case OpUnion, OpDifference:
		return !e.isEntry
	case OpXor:
		return true
	}
	return false
}

func traverseRing(start *entry, entries []*entry, op OpType) []geojson.Position {
	var ring []geojson.Position
	curEntry := start
	onA := true

	for {
		curEntry.used = true

		var p *ptNode
		if onA {
			p = curEntry.onA
		} else {
			p = curEntry.onB
		}
		ring = append(ring, p.point)

		// Determine traversal direction:
		// Intersection: forward on A → forward on B
		// Union: forward on A → forward on B
		// Difference (A\B): forward on A → BACKWARD on B
		useNext := true
		if !onA && op == OpDifference {
			useNext = false
		}

		cur := p.next
		if !useNext {
			cur = p.prev
		}

		for cur != nil && !cur.isect && !cur.used {
			ring = append(ring, cur.point)
			cur.used = true
			if useNext {
				cur = cur.next
			} else {
				cur = cur.prev
			}
		}

		if cur == nil || cur.used {
			break
		}

		var nextEntry *entry
		for _, e := range entries {
			if e.used {
				continue
			}
			if e.onA == cur || e.onB == cur {
				nextEntry = e
				break
			}
		}

		if nextEntry == nil || nextEntry == start {
			break
		}

		onA = !onA
		curEntry = nextEntry
	}

	return ring
}

func mergeRings(rings [][][]geojson.Position) [][][]geojson.Position {
	if len(rings) <= 1 {
		return rings
	}

	var result [][][]geojson.Position
	used := make([]bool, len(rings))

	for i := 0; i < len(rings); i++ {
		if used[i] {
			continue
		}
		exterior := rings[i][0]
		var holes [][]geojson.Position

		for j := 0; j < len(rings); j++ {
			if i == j || used[j] {
				continue
			}
			r := rings[j][0]
			if len(r) < 3 {
				continue
			}
			if !pointInRing(r[0], exterior) {
				continue
			}
			isHole := false
			for _, h := range holes {
				if pointInRing(r[0], h) {
					isHole = true
					break
				}
			}
			if !isHole {
				holes = append(holes, r)
				used[j] = true
			}
		}

		poly := [][]geojson.Position{exterior}
		if len(holes) > 0 {
			poly = append(poly, holes...)
		}
		result = append(result, poly)
		used[i] = true
	}

	return result
}
