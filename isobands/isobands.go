package isobands

import (
	"fmt"
	"math"

	"github.com/ibinh/turf-go/geojson"
	"github.com/ibinh/turf-go/meta"
)

type IsobandsOptions struct {
	ZProperty string
	Breaks    []float64
}

type zPoint struct {
	pos geojson.Position
	val float64
}

type gridCell struct {
	corners [4]geojson.Position
	values  [4]float64
}

// edges: 0=bottom(0->1), 1=right(1->2), 2=top(2->3), 3=left(3->0)
// each entry is pairs of edge indices forming segments within the cell
var bandSegments = [16][][2]int{
	0:  {},
	1:  {{0, 3}},     // SW in
	2:  {{0, 1}},     // SE in
	3:  {{1, 3}},     // SW,SE in
	4:  {{1, 2}},     // NE in
	5:  {{0, 1}, {2, 3}}, // SW,NE in (saddle)
	6:  {{0, 2}},     // SE,NE in
	7:  {{2, 3}},     // SW,SE,NE in (NW out)
	8:  {{2, 3}},     // NW in
	9:  {{0, 2}},     // SW,NW in
	10: {{0, 1}, {2, 3}}, // SE,NW in (saddle)
	11: {{1, 2}},    // SW,NE,NW in (SE out)
	12: {{1, 3}},    // NE,NW in (SW,SE out)
	13: {{0, 1}},    // SE,NE,NW in (SW out)
	14: {{0, 3}},    // SW,SE,NW in (NE out)
	15: {},
}

func Isobands(points *geojson.FeatureCollection, options ...IsobandsOptions) (*geojson.FeatureCollection, error) {
	opts := IsobandsOptions{ZProperty: "z"}
	if len(options) > 0 {
		opts = options[0]
	}
	if len(opts.Breaks) < 2 {
		return nil, fmt.Errorf("at least 2 breaks required")
	}
	if points == nil || len(points.Features) < 3 {
		return nil, fmt.Errorf("at least 3 points required")
	}

	data := extractData(points, opts.ZProperty)
	if len(data) < 3 {
		return nil, fmt.Errorf("not enough points with property %q", opts.ZProperty)
	}

	grid := createGrid(data)

	var features []*geojson.Feature
	for i := 0; i < len(opts.Breaks)-1; i++ {
		lo, hi := opts.Breaks[i], opts.Breaks[i+1]
		rings := marchingSquaresBands(grid, lo, hi)
		for _, ring := range rings {
			if len(ring) >= 3 {
				poly := geojson.NewPolygon([][]geojson.Position{ring})
				features = append(features, geojson.NewFeature(poly, map[string]any{
					"break": lo,
				}))
			}
		}
	}

	return geojson.NewFeatureCollection(features), nil
}

func extractData(fc *geojson.FeatureCollection, prop string) []zPoint {
	var data []zPoint
	meta.FeatureEach(fc, func(f *geojson.Feature, _ int) error {
		coord, err := geojson.GetCoord(f)
		if err != nil {
			return nil
		}
		v, ok := f.Properties[prop].(float64)
		if !ok {
			return nil
		}
		data = append(data, zPoint{coord, v})
		return nil
	})
	return data
}

func createGrid(data []zPoint) []gridCell {
	if len(data) < 2 {
		return nil
	}
	minX, minY := data[0].pos[0], data[0].pos[1]
	maxX, maxY := minX, minY
	for _, d := range data {
		if d.pos[0] < minX {
			minX = d.pos[0]
		}
		if d.pos[0] > maxX {
			maxX = d.pos[0]
		}
		if d.pos[1] < minY {
			minY = d.pos[1]
		}
		if d.pos[1] > maxY {
			maxY = d.pos[1]
		}
	}

	cellSize := estimateCellSize(data, minX, maxX, minY, maxY)
	if cellSize <= 0 {
		cellSize = math.Max((maxX-minX)/20, (maxY-minY)/20)
		if cellSize <= 0 {
			cellSize = 0.01
		}
	}

	nx := int(math.Ceil((maxX - minX) / cellSize))
	ny := int(math.Ceil((maxY - minY) / cellSize))
	if nx < 2 {
		nx = 2
	}
	if ny < 2 {
		ny = 2
	}

	var cells []gridCell
	for j := 0; j < ny; j++ {
		for i := 0; i < nx; i++ {
			x0 := minX + float64(i)*cellSize
			x1 := x0 + cellSize
			y0 := minY + float64(j)*cellSize
			y1 := y0 + cellSize

			pts := [4]geojson.Position{
				{x0, y0}, {x1, y0}, {x1, y1}, {x0, y1},
			}
			vals := [4]float64{
				idw(pts[0], data, 2),
				idw(pts[1], data, 2),
				idw(pts[2], data, 2),
				idw(pts[3], data, 2),
			}
			cells = append(cells, gridCell{corners: pts, values: vals})
		}
	}
	return cells
}

func estimateCellSize(data []zPoint, minX, maxX, minY, maxY float64) float64 {
	if len(data) < 2 {
		return 0
	}
	dx := maxX - minX
	dy := maxY - minY
	area := dx * dy
	if area <= 0 {
		return 0
	}
	return math.Sqrt(area / float64(len(data)))
}

func idw(pt geojson.Position, data []zPoint, weight float64) float64 {
	var sumW, sumV float64
	for _, d := range data {
		dx := pt[0] - d.pos[0]
		dy := pt[1] - d.pos[1]
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist < 1e-15 {
			return d.val
		}
		w := 1.0 / math.Pow(dist, weight)
		sumW += w
		sumV += w * d.val
	}
	if sumW < 1e-15 {
		return 0
	}
	return sumV / sumW
}

type rawSeg struct {
	a, b geojson.Position
}

func marchingSquaresBands(grid []gridCell, lo, hi float64) [][]geojson.Position {
	var segs []rawSeg

	for _, cell := range grid {
		pts, vals := cell.corners, cell.values

		// compute inside-band mask
		mask := 0
		for i := 0; i < 4; i++ {
			if vals[i] >= lo && vals[i] < hi {
				mask |= 1 << i
			}
		}
		if mask == 0 || mask == 15 {
			continue
		}

		pairs := bandSegments[mask]
		for _, pair := range pairs {
			e1, e2 := pair[0], pair[1]
			// For each pair of edges, get the interpolated crossing points
			p1 := interpolateEdgeCrossing(pts, vals, e1, lo, hi)
			p2 := interpolateEdgeCrossing(pts, vals, e2, lo, hi)
			segs = append(segs, rawSeg{p1, p2})
		}
	}

	return connectSegmentsToRings(segs)
}

// edge vertex indices: 0=SW, 1=SE, 2=NE, 3=NW
var edgeVerts = [4][2]int{{0, 1}, {1, 2}, {2, 3}, {3, 0}}

// interpolateEdgeCrossing finds the crossing point on a cell edge for the band boundary.
// The crossing threshold is lo or hi depending on which side of the band is crossed.
func interpolateEdgeCrossing(pts [4]geojson.Position, vals [4]float64, edgeIdx int, lo, hi float64) geojson.Position {
	vi := edgeVerts[edgeIdx][0]
	vj := edgeVerts[edgeIdx][1]

	vai, vaj := vals[vi], vals[vj]
	pi, pj := pts[vi], pts[vj]

	inI := vai >= lo && vai < hi
	inJ := vaj >= lo && vaj < hi

	threshold := lo
	if inI && !inJ {
		if vaj >= hi {
			threshold = hi
		}
	} else if !inI && inJ {
		if vai >= hi {
			threshold = hi
		}
	} else {
		// both outside but straddling entire band
		// Pick whichever boundary gives a t in [0,1]
		tLo := (lo - vai) / (vaj - vai)
		if tLo >= 0 && tLo <= 1 {
			threshold = lo
		} else {
			threshold = hi
		}
	}

	t := (threshold - vai) / (vaj - vai)
	if t < 0 || t > 1 {
		t = 0.5
	}
	return geojson.Position{
		pi[0] + t*(pj[0]-pi[0]),
		pi[1] + t*(pj[1]-pi[1]),
	}
}

func connectSegmentsToRings(segs []rawSeg) [][]geojson.Position {
	if len(segs) == 0 {
		return nil
	}

	// Build adjacency: each point maps to its neighbors
	type posKey struct{ x, y int64 }
	toKey := func(p geojson.Position) posKey {
		return posKey{int64(math.Round(p[0] * 1e10)), int64(math.Round(p[1] * 1e10))}
	}

	adj := map[posKey][]geojson.Position{}
	ptMap := map[posKey]geojson.Position{}

	addEdge := func(a, b geojson.Position) {
		ka, kb := toKey(a), toKey(b)
		adj[ka] = append(adj[ka], b)
		adj[kb] = append(adj[kb], a)
		ptMap[ka] = a
		ptMap[kb] = b
	}

	for _, s := range segs {
		addEdge(s.a, s.b)
	}

	used := map[posKey]bool{}
	var rings [][]geojson.Position

	for len(adj) > 0 {
		// find start point (prefer an endpoint with degree != 2)
		var start posKey
		for k := range adj {
			if len(adj[k]) != 2 {
				start = k
				break
			}
		}
		if start == (posKey{}) {
			for k := range adj {
				start = k
				break
			}
		}

		var ring []geojson.Position
		cur := start

		for {
			if used[cur] {
				break
			}
			used[cur] = true
			ring = append(ring, ptMap[cur])

			neighbors := adj[cur]
			delete(adj, cur)

			var next posKey
			found := false
			for _, nb := range neighbors {
				nk := toKey(nb)
				if !used[nk] {
					next = nk
					found = true
					break
				}
			}
			if !found {
				break
			}
			cur = next
		}

		if len(ring) > 2 {
			rings = append(rings, ring)
		}
	}

	return rings
}
