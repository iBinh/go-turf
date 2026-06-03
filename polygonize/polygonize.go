package polygonize

import (
	"fmt"
	"math"
	"sort"

	"github.com/ibinh/turf-go/geojson"
	"github.com/ibinh/turf-go/lines"
	"github.com/ibinh/turf-go/meta"
)

type position struct {
	x, y float64
}

func pos(p geojson.Position) position { return position{p[0], p[1]} }

type segment struct {
	a, b position
}

func (s segment) key() string {
	if s.a.x < s.b.x || (s.a.x == s.b.x && s.a.y < s.b.y) {
		return fmt.Sprintf("%.10f,%.10f-%.10f,%.10f", s.a.x, s.a.y, s.b.x, s.b.y)
	}
	return fmt.Sprintf("%.10f,%.10f-%.10f,%.10f", s.b.x, s.b.y, s.a.x, s.a.y)
}

func Polygonize(geom any) (*geojson.FeatureCollection, error) {
	// Collect all segments from the input geometry
	segments, err := collectSegments(geom)
	if err != nil {
		return nil, err
	}
	if len(segments) < 3 {
		return nil, fmt.Errorf("not enough edges to form polygons")
	}

	// Build adjacency graph
	adj := map[position][]position{}
	segSet := map[string]bool{}

	for _, s := range segments {
		key := s.key()
		if segSet[key] {
			continue // deduplicate
		}
		segSet[key] = true

		adj[s.a] = append(adj[s.a], s.b)
		adj[s.b] = append(adj[s.b], s.a)
	}

	// Find all simple cycles by walking the graph
	// For each node, take the next clockwise edge and walk
	var rings [][]geojson.Position
	usedEdges := map[string]bool{}

	for start, neighbors := range adj {
		for _, nb := range neighbors {
			key := edgeKey(start, nb)
			if usedEdges[key] {
				continue
			}

			// Walk the polygon
			ring := walkPolygon(start, nb, adj, usedEdges)
			if len(ring) >= 3 {
				// Close the ring
				ring = append(ring, ring[0])
				rings = append(rings, ring)
			}
		}
	}

	if len(rings) == 0 {
		return nil, fmt.Errorf("no polygons found")
	}

	// Filter to minimal polygons and create features
	// Sort by area descending, keep outer rings and discard duplicates
	sort.Slice(rings, func(i, j int) bool {
		return ringArea(rings[i]) > ringArea(rings[j])
	})

	var features []*geojson.Feature
	used := make([]bool, len(rings))

	for i, ring := range rings {
		if used[i] {
			continue
		}
		// Check if this ring contains any other ring (hole)
		var holes [][]geojson.Position
		for j, other := range rings {
			if i == j || used[j] {
				continue
			}
			if ringContainsRing(ring, other) {
				holes = append(holes, other)
				used[j] = true
			}
		}

		// Build polygon with holes
		coords := [][]geojson.Position{ring}
		for _, h := range holes {
			coords = append(coords, h)
		}

		poly := geojson.NewPolygon(coords)
		features = append(features, geojson.NewFeature(poly, nil))
	}

	return geojson.NewFeatureCollection(features), nil
}

func collectSegments(geom any) ([]segment, error) {
	// Use LineSegment to extract all segments from the geometry
	segs, err := lines.LineSegment(geom)
	if err != nil || segs == nil {
		// Fallback: extract from coordinates manually
		return collectSegmentsManual(geom)
	}

	var result []segment
	for _, f := range segs.Features {
		coords, err := geojson.GetCoords(f)
		if err != nil {
			continue
		}
		pts, ok := coords.([]geojson.Position)
		if !ok || len(pts) < 2 {
			continue
		}
		result = append(result, segment{pos(pts[0]), pos(pts[1])})
	}
	return result, nil
}

func collectSegmentsManual(geom any) ([]segment, error) {
	var result []segment

	err := meta.GeomEach(geom, func(g geojson.Geometry, _ int) error {
		switch geo := g.(type) {
		case *geojson.LineString:
			for i := 0; i < len(geo.Coordinates)-1; i++ {
				result = append(result, segment{pos(geo.Coordinates[i]), pos(geo.Coordinates[i+1])})
			}
		case *geojson.MultiLineString:
			for _, line := range geo.Coordinates {
				for i := 0; i < len(line)-1; i++ {
					result = append(result, segment{pos(line[i]), pos(line[i+1])})
				}
			}
		case *geojson.Polygon:
			for _, ring := range geo.Coordinates {
				for i := 0; i < len(ring)-1; i++ {
					result = append(result, segment{pos(ring[i]), pos(ring[i+1])})
				}
			}
		case *geojson.MultiPolygon:
			for _, poly := range geo.Coordinates {
				for _, ring := range poly {
					for i := 0; i < len(ring)-1; i++ {
						result = append(result, segment{pos(ring[i]), pos(ring[i+1])})
					}
				}
			}
		}
		return nil
	})

	return result, err
}

func edgeKey(a, b position) string {
	if a.x < b.x || (a.x == b.x && a.y < b.y) {
		return fmt.Sprintf("%.10f,%.10f-%.10f,%.10f", a.x, a.y, b.x, b.y)
	}
	return fmt.Sprintf("%.10f,%.10f-%.10f,%.10f", b.x, b.y, a.x, a.y)
}

func walkPolygon(start, next position, adj map[position][]position, usedEdges map[string]bool) []geojson.Position {
	var ring []geojson.Position
	cur := start
	prev := next

	for {
		key := edgeKey(cur, prev)
		if usedEdges[key] {
			break
		}
		usedEdges[key] = true
		ring = append(ring, geojson.Position{cur.x, cur.y})

		nxt, foundNext := findNextEdge(prev, cur, adj, usedEdges)
		if !foundNext {
			break
		}

		cur = prev
		prev = nxt

		// Check if we're back to the start
		if cur == start && len(ring) > 2 {
			break
		}

		// Safety limit
		if len(ring) > len(usedEdges)+1 {
			break
		}
	}

	return ring
}

func findNextEdge(from, comingFrom position, adj map[position][]position, usedEdges map[string]bool) (position, bool) {
	neighbors := adj[from]
	if len(neighbors) == 0 {
		return position{}, false
	}

	angle := math.Atan2(comingFrom.y-from.y, comingFrom.x-from.x)

	type cand struct {
		nb    position
		angle float64
	}
	var candidates []cand

	for _, nb := range neighbors {
		if nb == comingFrom {
			continue
		}
		a := math.Atan2(nb.y-from.y, nb.x-from.x)
		diff := angle - a
		if diff < 0 {
			diff += 2 * math.Pi
		}
		candidates = append(candidates, cand{nb, diff})
	}

	if len(candidates) == 0 {
		return position{}, false
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].angle < candidates[j].angle
	})

	return candidates[0].nb, true
}

func ringArea(ring []geojson.Position) float64 {
	var area float64
	for i := 0; i < len(ring)-1; i++ {
		area += ring[i][0]*ring[i+1][1] - ring[i+1][0]*ring[i][1]
	}
	return math.Abs(area) / 2
}

func pointInRing(p geojson.Position, ring []geojson.Position) bool {
	inside := false
	for i, j := 0, len(ring)-1; i < len(ring); j, i = i, i+1 {
		if (ring[j][1] > p[1]) != (ring[i][1] > p[1]) &&
			p[0] < (ring[i][0]-ring[j][0])*(p[1]-ring[j][1])/(ring[i][1]-ring[j][1])+ring[j][0] {
			inside = !inside
		}
	}
	return inside
}

func ringContainsRing(outer, inner []geojson.Position) bool {
	if len(outer) == 0 || len(inner) == 0 {
		return false
	}
	// Check if first point of inner is inside outer
	return pointInRing(inner[0], outer)
}
