package simplify

import (
	"fmt"
	"sort"

	"github.com/ibinh/turf-go/geojson"
	"github.com/ibinh/turf-go/meta"
)

func Simplify(geom any, tolerance float64, highQuality bool) (*geojson.Feature, error) {
	if tolerance <= 0 {
		return nil, fmt.Errorf("tolerance must be positive")
	}

	var resultGeom geojson.Geometry
	err := meta.GeomEach(geom, func(g geojson.Geometry, _ int) error {
		var err error
		resultGeom, err = simplifyGeometry(g, tolerance, highQuality)
		return err
	})
	if err != nil {
		return nil, err
	}
	return geojson.NewFeature(resultGeom, nil), nil
}

func simplifyGeometry(g geojson.Geometry, tolerance float64, highQuality bool) (geojson.Geometry, error) {
	switch v := g.(type) {
	case *geojson.Point, *geojson.MultiPoint:
		return g, nil
	case *geojson.LineString:
		return geojson.NewLineString(rdpSimplify(v.Coordinates, tolerance)), nil
	case *geojson.MultiLineString:
		lines := make([][]geojson.Position, len(v.Coordinates))
		for i, line := range v.Coordinates {
			lines[i] = rdpSimplify(line, tolerance)
		}
		return geojson.NewMultiLineString(lines), nil
	case *geojson.Polygon:
		rings := make([][]geojson.Position, len(v.Coordinates))
		for i, ring := range v.Coordinates {
			simplified := rdpSimplify(ring[:len(ring)-1], tolerance)
			rings[i] = append(simplified, simplified[0])
		}
		return geojson.NewPolygon(rings), nil
	case *geojson.MultiPolygon:
		polygons := make([][][]geojson.Position, len(v.Coordinates))
		for i, poly := range v.Coordinates {
			rings := make([][]geojson.Position, len(poly))
			for j, ring := range poly {
				simplified := rdpSimplify(ring[:len(ring)-1], tolerance)
				rings[j] = append(simplified, simplified[0])
			}
			polygons[i] = rings
		}
		return geojson.NewMultiPolygon(polygons), nil
	default:
		return g, nil
	}
}

func rdpSimplify(pts []geojson.Position, tolerance float64) []geojson.Position {
	if len(pts) <= 2 {
		return pts
	}

	sqTol := tolerance * tolerance

	maxDist := 0.0
	maxIdx := 0
	first, last := pts[0], pts[len(pts)-1]

	for i := 1; i < len(pts)-1; i++ {
		dist := sqSegmentDist(pts[i], first, last)
		if dist > maxDist {
			maxDist = dist
			maxIdx = i
		}
	}

	if maxDist > sqTol {
		left := rdpSimplify(pts[:maxIdx+1], tolerance)
		right := rdpSimplify(pts[maxIdx:], tolerance)
		return append(left[:len(left)-1], right...)
	}

	return []geojson.Position{first, last}
}

func sqSegmentDist(p, a, b geojson.Position) float64 {
	abx := b[0] - a[0]
	aby := b[1] - a[1]
	apx := p[0] - a[0]
	apy := p[1] - a[1]

	dot := apx*abx + apy*aby
	lenSq := abx*abx + aby*aby

	if lenSq == 0 {
		return apx*apx + apy*apy
	}

	t := dot / lenSq
	if t < 0 {
		return apx*apx + apy*apy
	}
	if t > 1 {
		return (p[0]-b[0])*(p[0]-b[0]) + (p[1]-b[1])*(p[1]-b[1])
	}

	projX := a[0] + t*abx
	projY := a[1] + t*aby
	dx := p[0] - projX
	dy := p[1] - projY
	return dx*dx + dy*dy
}

func ConvexHull(geom any) (*geojson.Feature, error) {
	var pts []geojson.Position
	err := meta.CoordEach(geom, func(c geojson.Position, _ int) error {
		pts = append(pts, c)
		return nil
	})
	if err != nil {
		return nil, err
	}

	if len(pts) < 3 {
		return nil, fmt.Errorf("need at least 3 points for convex hull")
	}

	hull := monotoneChain(pts)
	if len(hull) < 2 {
		pt := geojson.NewPoint(pts[0])
		return geojson.NewFeature(pt, nil), nil
	}
	if len(hull) < 3 {
		ls := geojson.NewLineString([]geojson.Position{hull[0], hull[1]})
		return geojson.NewFeature(ls, nil), nil
	}

	ring := make([]geojson.Position, len(hull)+1)
	copy(ring, hull)
	ring[len(hull)] = hull[0]

	return geojson.NewFeature(geojson.NewPolygon([][]geojson.Position{ring}), nil), nil
}

func monotoneChain(pts []geojson.Position) []geojson.Position {
	sorted := make([]geojson.Position, len(pts))
	copy(sorted, pts)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i][0] != sorted[j][0] {
			return sorted[i][0] < sorted[j][0]
		}
		return sorted[i][1] < sorted[j][1]
	})

	if len(sorted) <= 1 {
		return sorted
	}

	lower := []geojson.Position{}
	for _, p := range sorted {
		for len(lower) >= 2 && cross(lower[len(lower)-2], lower[len(lower)-1], p) <= 0 {
			lower = lower[:len(lower)-1]
		}
		lower = append(lower, p)
	}

	upper := []geojson.Position{}
	for i := len(sorted) - 1; i >= 0; i-- {
		p := sorted[i]
		for len(upper) >= 2 && cross(upper[len(upper)-2], upper[len(upper)-1], p) <= 0 {
			upper = upper[:len(upper)-1]
		}
		upper = append(upper, p)
	}

	result := lower[:len(lower)-1]
	result = append(result, upper[:len(upper)-1]...)
	return result
}

func cross(o, a, b geojson.Position) float64 {
	return (a[0]-o[0])*(b[1]-o[1]) - (a[1]-o[1])*(b[0]-o[0])
}
