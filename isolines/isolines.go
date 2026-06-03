package isolines

import (
	"fmt"
	"math"

	"github.com/ibinh/turf-go/geojson"
	"github.com/ibinh/turf-go/meta"
)

type IsolinesOptions struct {
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
var contourSegments = [16][][2]int{
	0:  {},
	1:  {{0, 3}},
	2:  {{0, 1}},
	3:  {{1, 3}},
	4:  {{1, 2}},
	5:  {{0, 1}, {2, 3}},
	6:  {{0, 2}},
	7:  {{2, 3}},
	8:  {{2, 3}},
	9:  {{0, 2}},
	10: {{0, 1}, {2, 3}},
	11: {{1, 2}},
	12: {{1, 3}},
	13: {{0, 1}},
	14: {{0, 3}},
	15: {},
}

var edgeVerts = [4][2]int{{0, 1}, {1, 2}, {2, 3}, {3, 0}}

func Isolines(points *geojson.FeatureCollection, options ...IsolinesOptions) (*geojson.FeatureCollection, error) {
	opts := IsolinesOptions{ZProperty: "z"}
	if len(options) > 0 {
		opts = options[0]
	}
	if len(opts.Breaks) == 0 {
		return nil, fmt.Errorf("at least 1 break required")
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
	for _, contourVal := range opts.Breaks {
		polylines := marchingSquaresLines(grid, contourVal)
		for _, polyline := range polylines {
			if len(polyline) >= 2 {
				ls := geojson.NewLineString(polyline)
				features = append(features, geojson.NewFeature(ls, map[string]any{
					"break": contourVal,
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

func marchingSquaresLines(grid []gridCell, threshold float64) [][]geojson.Position {
	var segs []rawSeg

	for _, cell := range grid {
		pts, vals := cell.corners, cell.values

		mask := 0
		for i := 0; i < 4; i++ {
			if vals[i] >= threshold {
				mask |= 1 << i
			}
		}
		if mask == 0 || mask == 15 {
			continue
		}

		pairs := contourSegments[mask]
		for _, pair := range pairs {
			e1, e2 := pair[0], pair[1]
			p1 := interpolateContour(pts, vals, e1, threshold)
			p2 := interpolateContour(pts, vals, e2, threshold)
			segs = append(segs, rawSeg{p1, p2})
		}
	}

	return connectSegmentsToPolylines(segs)
}

func interpolateContour(pts [4]geojson.Position, vals [4]float64, edgeIdx int, threshold float64) geojson.Position {
	vi := edgeVerts[edgeIdx][0]
	vj := edgeVerts[edgeIdx][1]

	pi, pj := pts[vi], pts[vj]
	vai, vaj := vals[vi], vals[vj]

	t := (threshold - vai) / (vaj - vai)
	if t < 0 || t > 1 {
		t = 0.5
	}
	return geojson.Position{
		pi[0] + t*(pj[0]-pi[0]),
		pi[1] + t*(pj[1]-pi[1]),
	}
}

func connectSegmentsToPolylines(segs []rawSeg) [][]geojson.Position {
	if len(segs) == 0 {
		return nil
	}

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
	var polylines [][]geojson.Position

	for len(adj) > 0 {
		var start posKey
		var foundStart bool
		for k := range adj {
			if len(adj[k]) == 1 {
				start = k
				foundStart = true
				break
			}
		}
		if !foundStart {
			for k := range adj {
				start = k
				break
			}
		}

		var line []geojson.Position
		cur := start

		for {
			if used[cur] {
				break
			}
			used[cur] = true
			line = append(line, ptMap[cur])

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

		if len(line) >= 2 {
			polylines = append(polylines, line)
		}
	}

	return polylines
}
