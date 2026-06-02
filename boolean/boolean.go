package boolean

import (
	"fmt"
	"math"
	"reflect"
	"github.com/ibinh/turf-go/geojson"
)

func Clockwise(ring []geojson.Position) bool {
	area := 0.0
	for i := 0; i < len(ring)-1; i++ {
		area += (ring[i+1][0] - ring[i][0]) * (ring[i+1][1] + ring[i][1])
	}
	return area > 0
}

func PointInPolygon(point any, polygon any) (bool, error) {
	poly, err := geojson.GetGeometry(polygon)
	if err != nil {
		return false, err
	}
	var coords []geojson.Position
	switch p := point.(type) {
	case *geojson.Feature:
		return PointInPolygon(p.Geometry, polygon)
	case *geojson.Point:
		coords = []geojson.Position{p.Coordinates}
	case *geojson.MultiPoint:
		coords = p.Coordinates
	case geojson.Geometry:
		pt, err := geojson.GetCoord(point)
		if err != nil {
			return false, err
		}
		coords = []geojson.Position{pt}
	default:
		pt, err := geojson.GetCoord(point)
		if err != nil {
			return false, err
		}
		coords = []geojson.Position{pt}
	}
	checks := 0
	for _, pt := range coords {
		switch p := poly.(type) {
		case *geojson.Polygon:
			if pointInPolygonCoords(pt, p.Coordinates) {
				checks++
			}
		case *geojson.MultiPolygon:
			for _, polyCoords := range p.Coordinates {
				if pointInPolygonCoords(pt, polyCoords) {
					checks++
					break
				}
			}
		default:
			return pointInPolygon(pt, poly)
		}
	}
	return checks == len(coords), nil
}

func pointInPolygonCoords(pt geojson.Position, rings [][]geojson.Position) bool {
	if pointOnRingBoundary(pt, rings[0]) {
		return true
	}
	inside := rayCast(pt, rings[0])
	if !inside {
		return false
	}
	for i := 1; i < len(rings); i++ {
		if pointOnRingBoundary(pt, rings[i]) {
			return true
		}
		if rayCast(pt, rings[i]) {
			return false
		}
	}
	return true
}

func pointOnRingBoundary(pt geojson.Position, ring []geojson.Position) bool {
	for i := 0; i < len(ring)-1; i++ {
		if isPointOnSegment(pt, ring[i], ring[i+1]) {
			return true
		}
	}
	return false
}

func rayCast(pt geojson.Position, ring []geojson.Position) bool {
	x, y := pt[0], pt[1]
	inside := false
	n := len(ring)
	j := n - 1
	for i := 0; i < n; i++ {
		xi, yi := ring[i][0], ring[i][1]
		xj, yj := ring[j][0], ring[j][1]
		if ((yi > y) != (yj > y)) && (x < (xj-xi)*(y-yi)/(yj-yi)+xi) {
			inside = !inside
		}
		j = i
	}
	return inside
}

func PointOnLine(point any, line any, ignoreEndpoints bool) (bool, error) {
	pt, err := geojson.GetCoord(point)
	if err != nil {
		return false, err
	}
	lineGeom, err := geojson.GetGeometry(line)
	if err != nil {
		return false, err
	}
	switch l := lineGeom.(type) {
	case *geojson.LineString:
		return pointOnLineCoords(pt, l.Coordinates, ignoreEndpoints), nil
	case *geojson.MultiLineString:
		for _, coords := range l.Coordinates {
			if pointOnLineCoords(pt, coords, ignoreEndpoints) {
				return true, nil
			}
		}
		return false, nil
	case *geojson.Polygon:
		for _, ring := range l.Coordinates {
			if pointOnLineCoords(pt, ring, ignoreEndpoints) {
				return true, nil
			}
		}
		return false, nil
	case *geojson.MultiPolygon:
		for _, poly := range l.Coordinates {
			for _, ring := range poly {
				if pointOnLineCoords(pt, ring, ignoreEndpoints) {
					return true, nil
				}
			}
		}
		return false, nil
	}
	return false, nil
}

func pointOnLineCoords(pt geojson.Position, line []geojson.Position, ignoreEndpoints bool) bool {
	for i := 0; i < len(line)-1; i++ {
		a := line[i]
		b := line[i+1]
		on := isPointOnSegment(pt, a, b)
		if on {
			if ignoreEndpoints && ((pt[0] == a[0] && pt[1] == a[1]) || (pt[0] == b[0] && pt[1] == b[1])) {
				continue
			}
			return true
		}
	}
	return false
}

func isPointOnSegment(p, a, b geojson.Position) bool {
	cross := (p[0]-a[0])*(b[1]-a[1]) - (p[1]-a[1])*(b[0]-a[0])
	if math.Abs(cross) > 1e-10 {
		return false
	}
	dot := (p[0]-a[0])*(b[0]-a[0]) + (p[1]-a[1])*(b[1]-a[1])
	if dot < 0 {
		return false
	}
	lenSq := (b[0]-a[0])*(b[0]-a[0]) + (b[1]-a[1])*(b[1]-a[1])
	if dot-lenSq > 1e-10 {
		return false
	}
	return true
}

func SegmentIntersect(a, b, c, d geojson.Position) bool {
	denom := (b[0]-a[0])*(d[1]-c[1]) - (b[1]-a[1])*(d[0]-c[0])
	if math.Abs(denom) < 1e-10 {
		return collinearSegmentIntersect(a, b, c, d)
	}
	t := ((c[0]-a[0])*(d[1]-c[1]) - (c[1]-a[1])*(d[0]-c[0])) / denom
	u := -((b[0]-a[0])*(c[1]-a[1]) - (b[1]-a[1])*(c[0]-a[0])) / denom
	return t >= -1e-10 && t <= 1+1e-10 && u >= -1e-10 && u <= 1+1e-10
}

func collinearSegmentIntersect(a, b, c, d geojson.Position) bool {
	crossAC := (c[0]-a[0])*(b[1]-a[1]) - (c[1]-a[1])*(b[0]-a[0])
	crossAD := (d[0]-a[0])*(b[1]-a[1]) - (d[1]-a[1])*(b[0]-a[0])
	if math.Abs(crossAC) > 1e-10 || math.Abs(crossAD) > 1e-10 {
		return false
	}
	abx := b[0] - a[0]
	aby := b[1] - a[1]
	acx := c[0] - a[0]
	acy := c[1] - a[1]
	adx := d[0] - a[0]
	ady := d[1] - a[1]
	dotAC := acx*abx + acy*aby
	dotAD := adx*abx + ady*aby
	lenSq := abx*abx + aby*aby

	t1 := dotAC / lenSq
	t2 := dotAD / lenSq

	if t1 > t2 {
		t1, t2 = t2, t1
	}

	if t2 < 0 || t1 > 1.0 {
		return false
	}
	return true
}

func Contains(geom1, geom2 any) (bool, error) {
	g1, err := geojson.GetGeometry(geom1)
	if err != nil {
		return false, err
	}
	g2, err := geojson.GetGeometry(geom2)
	if err != nil {
		return false, err
	}
	return contains(g1, g2)
}

func contains(g1, g2 geojson.Geometry) (bool, error) {
	switch g1 := g1.(type) {
	case *geojson.Point:
		return containsPoint(g1.Coordinates, g2)
	case *geojson.MultiPoint:
		for _, pt := range g1.Coordinates {
			ok, err := containsPoint(pt, g2)
			if !ok || err != nil {
				return false, err
			}
		}
		return true, nil
	case *geojson.LineString:
		return containsGeomInLine(g1.Coordinates, g2)
	case *geojson.MultiLineString:
		for _, line := range g1.Coordinates {
			ok, err := containsGeomInLine(line, g2)
			if !ok || err != nil {
				return false, err
			}
		}
		return true, nil
	case *geojson.Polygon:
		switch g2 := g2.(type) {
		case *geojson.Point:
			return PointInPolygon(g2, g1)
		case *geojson.MultiPoint:
			for _, pt := range g2.Coordinates {
				ok, err := PointInPolygon(geojson.NewPoint(pt), g1)
				if err != nil || !ok {
					return false, err
				}
			}
			return true, nil
		case *geojson.LineString:
			return lineInPolygon(g2.Coordinates, g1.Coordinates)
		case *geojson.MultiLineString:
			for _, line := range g2.Coordinates {
				ok, err := lineInPolygon(line, g1.Coordinates)
				if err != nil || !ok {
					return false, err
				}
			}
			return true, nil
		case *geojson.Polygon:
			return polygonInPolygon(g2.Coordinates, g1.Coordinates)
		case *geojson.MultiPolygon:
			for _, p := range g2.Coordinates {
				ok, err := contains(g1, geojson.NewPolygon(p))
				if !ok || err != nil {
					return false, err
				}
			}
			return true, nil
		}
	case *geojson.MultiPolygon:
		switch g2 := g2.(type) {
		case *geojson.Point:
			return PointInPolygon(g2, g1)
		case *geojson.MultiPoint:
			for _, pt := range g2.Coordinates {
				ok, err := PointInPolygon(geojson.NewPoint(pt), g1)
				if err != nil || !ok {
					return false, err
				}
			}
			return true, nil
		case *geojson.LineString:
			for _, poly := range g1.Coordinates {
				ok, err := contains(geojson.NewPolygon(poly), g2)
				if ok {
					return true, nil
				}
				if err != nil {
					return false, err
				}
			}
			return false, nil
		case *geojson.MultiLineString:
			for _, poly := range g1.Coordinates {
				ok, err := contains(geojson.NewPolygon(poly), g2)
				if ok {
					return true, nil
				}
				if err != nil {
					return false, err
				}
			}
			return false, nil
		case *geojson.Polygon:
			for _, poly := range g1.Coordinates {
				ok, err := contains(geojson.NewPolygon(poly), g2)
				if ok {
					return true, nil
				}
				if err != nil {
					return false, err
				}
			}
			return false, nil
		case *geojson.MultiPolygon:
			for _, poly := range g2.Coordinates {
				ok, err := contains(g1, geojson.NewPolygon(poly))
				if !ok || err != nil {
					return false, err
				}
			}
			return true, nil
		}
	}
	return false, nil
}

func containsPoint(pt geojson.Position, g2 geojson.Geometry) (bool, error) {
	switch g2 := g2.(type) {
	case *geojson.Point:
		return pt[0] == g2.Coordinates[0] && pt[1] == g2.Coordinates[1], nil
	case *geojson.MultiPoint:
		for _, p := range g2.Coordinates {
			if pt[0] == p[0] && pt[1] == p[1] {
				return true, nil
			}
		}
		return false, nil
	}
	return false, nil
}

func containsGeomInLine(line []geojson.Position, g2 geojson.Geometry) (bool, error) {
	switch g2 := g2.(type) {
	case *geojson.Point:
		return PointOnLine(g2, geojson.NewLineString(line), false)
	case *geojson.MultiPoint:
		for _, pt := range g2.Coordinates {
			ok, err := PointOnLine(geojson.NewPoint(pt), geojson.NewLineString(line), false)
			if err != nil || !ok {
				return false, err
			}
		}
		return true, nil
	case *geojson.LineString:
		for _, pt := range g2.Coordinates {
			ok, err := PointOnLine(geojson.NewPoint(pt), geojson.NewLineString(line), false)
			if err != nil || !ok {
				return false, err
			}
		}
		return true, nil
	}
	return false, nil
}

func lineInPolygon(line []geojson.Position, rings [][]geojson.Position) (bool, error) {
	coords := line
	if isRing(line) {
		coords = line[:len(line)-1]
	}
	for _, pt := range coords {
		onBoundary := false
		for _, ring := range rings {
			if pointOnRingBoundary(pt, ring) {
				onBoundary = true
				break
			}
		}
		if !onBoundary {
			in, err := PointInPolygon(geojson.NewPoint(pt), geojson.NewPolygon(rings))
			if err != nil || !in {
				return false, err
			}
		}
	}
	return true, nil
}

func isRing(coords []geojson.Position) bool {
	if len(coords) < 2 {
		return false
	}
	first := coords[0]
	last := coords[len(coords)-1]
	return first[0] == last[0] && first[1] == last[1]
}

func allPointsInPolygon(pts []geojson.Position, rings [][]geojson.Position) bool {
	for _, pt := range pts {
		if !pointInPolygonCoords(pt, rings) {
			return false
		}
	}
	return true
}

func polygonInPolygon(innerRings, outerRings [][]geojson.Position) (bool, error) {
	if !allPointsInPolygon(innerRings[0][:len(innerRings[0])-1], outerRings) {
		return false, nil
	}
	for i := 1; i < len(innerRings); i++ {
		for _, pt := range innerRings[i] {
			onBoundary := false
			for _, ring := range outerRings {
				if pointOnRingBoundary(pt, ring) {
					onBoundary = true
					break
				}
			}
			if !onBoundary {
				inside, _ := PointInPolygon(geojson.NewPoint(pt), geojson.NewPolygon(outerRings))
				if inside {
					return false, nil
				}
			}
		}
	}
	return true, nil
}

func Within(geom1, geom2 any) (bool, error) {
	return Contains(geom2, geom1)
}

func Intersects(geom1, geom2 any) (bool, error) {
	g1, err := geojson.GetGeometry(geom1)
	if err != nil {
		return false, err
	}
	g2, err := geojson.GetGeometry(geom2)
	if err != nil {
		return false, err
	}
	return intersects(g1, g2)
}

func intersects(g1, g2 geojson.Geometry) (bool, error) {
	t1, t2 := g1.Type(), g2.Type()

	if t1 == geojson.TypePoint || t1 == geojson.TypeMultiPoint {
		return pointOrMultiPointIntersects(g1, g2)
	}
	if t2 == geojson.TypePoint || t2 == geojson.TypeMultiPoint {
		return pointOrMultiPointIntersects(g2, g1)
	}

	if t1 == geojson.TypeLineString || t1 == geojson.TypeMultiLineString {
		return lineOrMultiLineIntersects(g1, g2)
	}
	if t2 == geojson.TypeLineString || t2 == geojson.TypeMultiLineString {
		return lineOrMultiLineIntersects(g2, g1)
	}

	if t1 == geojson.TypePolygon || t1 == geojson.TypeMultiPolygon {
		return polygonOrMultiPolygonIntersects(g1, g2)
	}
	return false, nil
}

func pointOrMultiPointIntersects(pointGeom, other geojson.Geometry) (bool, error) {
	var pts []geojson.Position
	switch p := pointGeom.(type) {
	case *geojson.Point:
		pts = []geojson.Position{p.Coordinates}
	case *geojson.MultiPoint:
		pts = p.Coordinates
	}

	for _, pt := range pts {
		switch o := other.(type) {
		case *geojson.Point:
			if pt[0] == o.Coordinates[0] && pt[1] == o.Coordinates[1] {
				return true, nil
			}
		case *geojson.MultiPoint:
			for _, op := range o.Coordinates {
				if pt[0] == op[0] && pt[1] == op[1] {
					return true, nil
				}
			}
		case *geojson.LineString:
			on, _ := PointOnLine(geojson.NewPoint(pt), o, false)
			if on {
				return true, nil
			}
		case *geojson.MultiLineString:
			on, _ := PointOnLine(geojson.NewPoint(pt), o, false)
			if on {
				return true, nil
			}
		case *geojson.Polygon:
			inside, _ := PointInPolygon(geojson.NewPoint(pt), o)
			if inside {
				return true, nil
			}
		case *geojson.MultiPolygon:
			inside, _ := PointInPolygon(geojson.NewPoint(pt), o)
			if inside {
				return true, nil
			}
		}
	}
	return false, nil
}

func lineOrMultiLineIntersects(lineGeom, other geojson.Geometry) (bool, error) {
	var lines [][]geojson.Position
	switch l := lineGeom.(type) {
	case *geojson.LineString:
		lines = [][]geojson.Position{l.Coordinates}
	case *geojson.MultiLineString:
		lines = l.Coordinates
	}

	for _, line := range lines {
		for i := 0; i < len(line)-1; i++ {
			a, b := line[i], line[i+1]

			switch o := other.(type) {
			case *geojson.LineString:
				for j := 0; j < len(o.Coordinates)-1; j++ {
					if SegmentIntersect(a, b, o.Coordinates[j], o.Coordinates[j+1]) {
						return true, nil
					}
				}
			case *geojson.MultiLineString:
				for _, l2 := range o.Coordinates {
					for j := 0; j < len(l2)-1; j++ {
						if SegmentIntersect(a, b, l2[j], l2[j+1]) {
							return true, nil
						}
					}
				}
			case *geojson.Polygon:
				for _, ring := range o.Coordinates {
					for j := 0; j < len(ring)-1; j++ {
						if SegmentIntersect(a, b, ring[j], ring[j+1]) {
							return true, nil
						}
					}
				}
			case *geojson.MultiPolygon:
				for _, poly := range o.Coordinates {
					for _, ring := range poly {
						for j := 0; j < len(ring)-1; j++ {
							if SegmentIntersect(a, b, ring[j], ring[j+1]) {
								return true, nil
							}
						}
					}
				}
			}
		}
	}
	return false, nil
}

func polygonOrMultiPolygonIntersects(polyGeom, other geojson.Geometry) (bool, error) {
	var polys [][][]geojson.Position
	switch p := polyGeom.(type) {
	case *geojson.Polygon:
		polys = [][][]geojson.Position{p.Coordinates}
	case *geojson.MultiPolygon:
		polys = p.Coordinates
	}

	for _, rings := range polys {
		switch o := other.(type) {
		case *geojson.LineString:
			for _, ring := range rings {
				for i := 0; i < len(ring)-1; i++ {
					for j := 0; j < len(o.Coordinates)-1; j++ {
						if SegmentIntersect(ring[i], ring[i+1], o.Coordinates[j], o.Coordinates[j+1]) {
							return true, nil
						}
					}
				}
			}
		case *geojson.MultiLineString:
			for _, line := range o.Coordinates {
				for _, ring := range rings {
					for i := 0; i < len(ring)-1; i++ {
						for j := 0; j < len(line)-1; j++ {
							if SegmentIntersect(ring[i], ring[i+1], line[j], line[j+1]) {
								return true, nil
							}
						}
					}
				}
			}
		case *geojson.Polygon:
			for _, ring := range rings {
				for _, otherRing := range o.Coordinates {
					for i := 0; i < len(ring)-1; i++ {
						for j := 0; j < len(otherRing)-1; j++ {
							if SegmentIntersect(ring[i], ring[i+1], otherRing[j], otherRing[j+1]) {
								return true, nil
							}
						}
					}
				}
			}
			inside, _ := contains(polyGeom, other)
			if inside {
				return true, nil
			}
			return false, nil
		case *geojson.MultiPolygon:
			for _, otherRings := range o.Coordinates {
				ok, err := polygonOrMultiPolygonIntersects(geojson.NewPolygon(rings), geojson.NewPolygon(otherRings))
				if err != nil {
					return false, err
				}
				if ok {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

func Disjoint(geom1, geom2 any) (bool, error) {
	intersects, err := Intersects(geom1, geom2)
	if err != nil {
		return false, err
	}
	return !intersects, nil
}

func Touches(geom1, geom2 any) (bool, error) {
	g1, err := geojson.GetGeometry(geom1)
	if err != nil {
		return false, err
	}
	g2, err := geojson.GetGeometry(geom2)
	if err != nil {
		return false, err
	}
	t1, t2 := g1.Type(), g2.Type()

	if t1 == geojson.TypePoint || t2 == geojson.TypePoint {
		return false, nil
	}

	intersects, err := Intersects(geom1, geom2)
	if err != nil || !intersects {
		return false, err
	}

	if t1 == geojson.TypePolygon || t1 == geojson.TypeMultiPolygon ||
		t2 == geojson.TypePolygon || t2 == geojson.TypeMultiPolygon {
		within, err := Within(geom1, geom2)
		if err != nil {
			return false, err
		}
		if within {
			return false, nil
		}
		within, err = Within(geom2, geom1)
		if err != nil {
			return false, err
		}
		if within {
			return false, nil
		}
	}

	return true, nil
}

func Crosses(geom1, geom2 any) (bool, error) {
	g1, err := geojson.GetGeometry(geom1)
	if err != nil {
		return false, err
	}
	g2, err := geojson.GetGeometry(geom2)
	if err != nil {
		return false, err
	}
	t1, t2 := g1.Type(), g2.Type()

	if t1 == t2 && t1 == geojson.TypePoint {
		return false, nil
	}

	intersects, err := Intersects(geom1, geom2)
	if err != nil || !intersects {
		return false, err
	}

	within, _ := Within(geom1, geom2)
	if within {
		return false, nil
	}
	within, _ = Within(geom2, geom1)
	if within {
		return false, nil
	}

	return true, nil
}

func Overlap(geom1, geom2 any) (bool, error) {
	g1, err := geojson.GetGeometry(geom1)
	if err != nil {
		return false, err
	}
	g2, err := geojson.GetGeometry(geom2)
	if err != nil {
		return false, err
	}
	t1, t2 := g1.Type(), g2.Type()

	if t1 != t2 {
		return false, nil
	}

	if t1 != geojson.TypePolygon && t1 != geojson.TypeMultiPolygon &&
		t1 != geojson.TypeLineString && t1 != geojson.TypeMultiLineString {
		return false, nil
	}

	intersects, err := Intersects(geom1, geom2)
	if err != nil || !intersects {
		return false, err
	}

	within, _ := Within(geom1, geom2)
	if within {
		return false, nil
	}
	within, _ = Within(geom2, geom1)
	if within {
		return false, nil
	}

	return true, nil
}

func pointInPolygon(pt geojson.Position, geom geojson.Geometry) (bool, error) {
	switch g := geom.(type) {
	case *geojson.Point:
		return pt[0] == g.Coordinates[0] && pt[1] == g.Coordinates[1], nil
	case *geojson.MultiPoint:
		for _, p := range g.Coordinates {
			if pt[0] == p[0] && pt[1] == p[1] {
				return true, nil
			}
		}
		return false, nil
	case *geojson.LineString:
		return pointOnLineCoords(pt, g.Coordinates, false), nil
	case *geojson.MultiLineString:
		for _, coords := range g.Coordinates {
			if pointOnLineCoords(pt, coords, false) {
				return true, nil
			}
		}
		return false, nil
	default:
		return false, nil
	}
}

func Valid(geom any) (bool, error) {
	g, err := geojson.GetGeometry(geom)
	if err != nil {
		return false, nil
	}
	switch v := g.(type) {
	case *geojson.Point:
		return len(v.Coordinates) >= 2, nil
	case *geojson.MultiPoint:
		if len(v.Coordinates) == 0 {
			return false, nil
		}
		for _, p := range v.Coordinates {
			if len(p) < 2 {
				return false, nil
			}
		}
		return true, nil
	case *geojson.LineString:
		return len(v.Coordinates) >= 2, nil
	case *geojson.MultiLineString:
		if len(v.Coordinates) == 0 {
			return false, nil
		}
		for _, line := range v.Coordinates {
			if len(line) < 2 {
				return false, nil
			}
		}
		return true, nil
	case *geojson.Polygon:
		return validPolygon(v.Coordinates), nil
	case *geojson.MultiPolygon:
		if len(v.Coordinates) == 0 {
			return false, nil
		}
		for _, poly := range v.Coordinates {
			if !validPolygon(poly) {
				return false, nil
			}
		}
		return true, nil
	}
	return false, nil
}

func validPolygon(rings [][]geojson.Position) bool {
	if len(rings) == 0 {
		return false
	}
	for _, ring := range rings {
		if len(ring) < 4 {
			return false
		}
		first := ring[0]
		last := ring[len(ring)-1]
		if first[0] != last[0] || first[1] != last[1] {
			return false
		}
	}
	return true
}

func Concave(geom any) (bool, error) {
	g, err := geojson.GetGeometry(geom)
	if err != nil {
		return false, err
	}
	poly, ok := g.(*geojson.Polygon)
	if !ok {
		return false, nil
	}
	ring := poly.Coordinates[0]
	n := len(ring)
	if n < 4 {
		return false, nil
	}
	cw := Clockwise(ring)
	for i := 0; i < n-1; i++ {
		a := ring[i]
		b := ring[(i+1)%(n-1)]
		c := ring[(i+2)%(n-1)]
		cross := (b[0]-a[0])*(c[1]-b[1]) - (b[1]-a[1])*(c[0]-b[0])
		if cw {
			if cross > 1e-10 {
				return true, nil
			}
		} else {
			if cross < -1e-10 {
				return true, nil
			}
		}
	}
	return false, nil
}

func BooleanEqual(geom1, geom2 any) (bool, error) {
	coords1, err := geojson.GetCoords(geom1)
	if err != nil {
		return false, err
	}
	coords2, err := geojson.GetCoords(geom2)
	if err != nil {
		return false, err
	}
	return reflect.DeepEqual(coords1, coords2), nil
}

func BooleanParallel(line1, line2 any) (bool, error) {
	coords1, err := geojson.GetCoords(line1)
	if err != nil {
		return false, err
	}
	coords2, err := geojson.GetCoords(line2)
	if err != nil {
		return false, err
	}

	var pts1, pts2 []geojson.Position
	switch c := coords1.(type) {
	case []geojson.Position:
		pts1 = c
	default:
		return false, fmt.Errorf("booleanParallel: expected LineString coordinates")
	}
	switch c := coords2.(type) {
	case []geojson.Position:
		pts2 = c
	default:
		return false, fmt.Errorf("booleanParallel: expected LineString coordinates")
	}

	if len(pts1) < 2 || len(pts2) < 2 {
		return false, fmt.Errorf("booleanParallel: each line must have at least 2 vertices")
	}

	for i := 0; i < len(pts1)-1; i++ {
		v1 := geojson.Position{pts1[i+1][0] - pts1[i][0], pts1[i+1][1] - pts1[i][1]}
		v2 := geojson.Position{pts2[i+1][0] - pts2[i][0], pts2[i+1][1] - pts2[i][1]}
		cross := v1[0]*v2[1] - v1[1]*v2[0]
		if math.Abs(cross) > 1e-10 {
			return false, nil
		}
	}
	return true, nil
}
