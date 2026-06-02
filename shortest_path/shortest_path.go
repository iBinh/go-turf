package shortestpath

import (
	"container/heap"
	"fmt"
	"math"

	"github.com/ibinh/turf-go/geojson"
	"github.com/ibinh/turf-go/meta"
)

type ShortestPathOptions struct {
	MaxDistance float64
}

func ShortestPath(start, end any, network *geojson.FeatureCollection, options ...ShortestPathOptions) (*geojson.Feature, error) {
	opts := ShortestPathOptions{MaxDistance: math.MaxFloat64}
	if len(options) > 0 {
		opts = options[0]
	}

	startCoord, err := geojson.GetCoord(start)
	if err != nil {
		return nil, fmt.Errorf("start: %w", err)
	}
	endCoord, err := geojson.GetCoord(end)
	if err != nil {
		return nil, fmt.Errorf("end: %w", err)
	}

	if network == nil || len(network.Features) == 0 {
		return nil, fmt.Errorf("network must have at least 1 line")
	}

	graph := newGraph()
	var nodeCoords []geojson.Position

	meta.FeatureEach(network, func(f *geojson.Feature, _ int) error {
		coords, err := geojson.GetCoords(f)
		if err != nil {
			return nil
		}
		pts, ok := coords.([]geojson.Position)
		if !ok || len(pts) < 2 {
			return nil
		}
		for i := 0; i < len(pts)-1; i++ {
			a, b := pts[i], pts[i+1]
			d := planarDist(a, b)
			graph.addEdge(a, b, d)
		}
		nodeCoords = append(nodeCoords, pts...)
		return nil
	})

	if graph.size() < 2 {
		return nil, fmt.Errorf("network graph too small")
	}

	startNode := findNearest(startCoord, nodeCoords, opts.MaxDistance)
	endNode := findNearest(endCoord, nodeCoords, opts.MaxDistance)

	if startNode == nil || endNode == nil {
		return nil, fmt.Errorf("start or end point not connected to network")
	}

	path := graph.shortestPath(*startNode, *endNode)
	if len(path) < 2 {
		return nil, fmt.Errorf("no path found")
	}

	return geojson.NewFeature(geojson.NewLineString(path), nil), nil
}

type posKey struct {
	x, y int64
}

func toKey(p geojson.Position) posKey {
	return posKey{int64(math.Round(p[0] * 1e10)), int64(math.Round(p[1] * 1e10))}
}

type edge struct {
	to     posKey
	weight float64
}

type graph struct {
	adj map[posKey][]edge
}

func newGraph() *graph {
	return &graph{adj: make(map[posKey][]edge)}
}

func (g *graph) addEdge(a, b geojson.Position, weight float64) {
	ka, kb := toKey(a), toKey(b)
	g.adj[ka] = append(g.adj[ka], edge{kb, weight})
	g.adj[kb] = append(g.adj[kb], edge{ka, weight})
}

func (g *graph) size() int {
	return len(g.adj)
}

func (g *graph) shortestPath(from, to geojson.Position) []geojson.Position {
	fk, tk := toKey(from), toKey(to)

	if fk == tk {
		return []geojson.Position{from, to}
	}

	dist := map[posKey]float64{}
	prev := map[posKey]posKey{}
	pq := &priorityQueue{}
	heap.Init(pq)

	dist[fk] = 0
	heap.Push(pq, &item{key: fk, priority: 0})

	for pq.Len() > 0 {
		cur := heap.Pop(pq).(*item)
		if cur.key == tk {
			break
		}
		if cur.priority > dist[cur.key] {
			continue
		}
		for _, e := range g.adj[cur.key] {
			nd := cur.priority + e.weight
			if d, ok := dist[e.to]; !ok || nd < d {
				dist[e.to] = nd
				prev[e.to] = cur.key
				heap.Push(pq, &item{key: e.to, priority: nd})
			}
		}
	}

	if _, ok := prev[tk]; !ok {
		return nil
	}

	var path []geojson.Position
	for cur := tk; ; cur = prev[cur] {
		path = append(path, geojson.Position{
			float64(cur.x) / 1e10,
			float64(cur.y) / 1e10,
		})
		if cur == fk {
			break
		}
	}

	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

func findNearest(pt geojson.Position, candidates []geojson.Position, maxDist float64) *geojson.Position {
	var nearest *geojson.Position
	minDist := maxDist
	for _, c := range candidates {
		d := planarDist(pt, c)
		if d < minDist {
			minDist = d
			c := c
			nearest = &c
		}
	}
	return nearest
}

func planarDist(a, b geojson.Position) float64 {
	dx := a[0] - b[0]
	dy := a[1] - b[1]
	return math.Sqrt(dx*dx + dy*dy)
}

type item struct {
	key      posKey
	priority float64
	index    int
}

type priorityQueue []*item

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[:n-1]
	return item
}
