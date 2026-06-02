package bbox

import (
	"fmt"
	"math"
	"github.com/ibinh/turf-go/geojson"
)

func BBox(obj any) ([]float64, error) {
	coords, err := getAllCoords(obj)
	if err != nil {
		return nil, err
	}
	if len(coords) == 0 {
		return nil, fmt.Errorf("bbox: no coordinates found")
	}

	minLng, minLat := coords[0][0], coords[0][1]
	maxLng, maxLat := minLng, minLat

	for _, c := range coords {
		if c[0] < minLng {
			minLng = c[0]
		}
		if c[0] > maxLng {
			maxLng = c[0]
		}
		if c[1] < minLat {
			minLat = c[1]
		}
		if c[1] > maxLat {
			maxLat = c[1]
		}
	}

	return []float64{minLng, minLat, maxLng, maxLat}, nil
}

func BBoxPolygon(bbox []float64) (*geojson.Feature, error) {
	if len(bbox) < 4 {
		return nil, fmt.Errorf("bbox-polygon: bbox must have 4 elements")
	}

	minLng, minLat := bbox[0], bbox[1]
	maxLng, maxLat := bbox[2], bbox[3]

	ring := []geojson.Position{
		{minLng, minLat},
		{maxLng, minLat},
		{maxLng, maxLat},
		{minLng, maxLat},
		{minLng, minLat},
	}

	return geojson.NewFeature(
		geojson.NewPolygon([][]geojson.Position{ring}),
		nil,
	), nil
}

func Envelope(obj any) (*geojson.Feature, error) {
	b, err := BBox(obj)
	if err != nil {
		return nil, err
	}
	return BBoxPolygon(b)
}

func Square(bbox []float64) ([]float64, error) {
	if len(bbox) < 4 {
		return nil, fmt.Errorf("square: bbox must have 4 elements")
	}

	width := bbox[2] - bbox[0]
	height := bbox[3] - bbox[1]

	if width >= height {
		cy := (bbox[1] + bbox[3]) / 2
		halfWidth := width / 2
		return []float64{
			bbox[0],
			cy - halfWidth,
			bbox[2],
			cy + halfWidth,
		}, nil
	}

	cx := (bbox[0] + bbox[2]) / 2
	halfHeight := height / 2
	return []float64{
		cx - halfHeight,
		bbox[1],
		cx + halfHeight,
		bbox[3],
	}, nil
}

func BBoxClip(obj any, bbox []float64) (*geojson.Feature, error) {
	if len(bbox) < 4 {
		return nil, fmt.Errorf("bbox must have 4 elements")
	}

	geom, err := geojson.GetGeometry(obj)
	if err != nil {
		return nil, err
	}

	minX, minY, maxX, maxY := bbox[0], bbox[1], bbox[2], bbox[3]

	switch g := geom.(type) {
	case *geojson.Point:
		if g.Coordinates[0] >= minX && g.Coordinates[0] <= maxX &&
			g.Coordinates[1] >= minY && g.Coordinates[1] <= maxY {
			return geojson.NewFeature(g, nil), nil
		}
		return nil, fmt.Errorf("point outside bbox")
	case *geojson.MultiPoint:
		var pts []geojson.Position
		for _, p := range g.Coordinates {
			if p[0] >= minX && p[0] <= maxX && p[1] >= minY && p[1] <= maxY {
				pts = append(pts, p)
			}
		}
		if len(pts) == 0 {
			return nil, fmt.Errorf("no points inside bbox")
		}
		return geojson.NewFeature(geojson.NewMultiPoint(pts), nil), nil
	case *geojson.LineString:
		clipped := clipLineString(g.Coordinates, minX, minY, maxX, maxY)
		if len(clipped) < 2 {
			return nil, fmt.Errorf("line clipped away")
		}
		longest := clipped[0]
		for _, seg := range clipped {
			if len(seg) > len(longest) {
				longest = seg
			}
		}
		return geojson.NewFeature(geojson.NewLineString(longest), nil), nil
	case *geojson.MultiLineString:
		var result [][]geojson.Position
		for _, line := range g.Coordinates {
			clipped := clipLineString(line, minX, minY, maxX, maxY)
			for _, seg := range clipped {
				if len(seg) >= 2 {
					result = append(result, seg)
				}
			}
		}
		if len(result) == 0 {
			return nil, fmt.Errorf("lines clipped away")
		}
		return geojson.NewFeature(geojson.NewMultiLineString(result), nil), nil
	case *geojson.Polygon:
		clippedRing := clipPolygonRing(g.Coordinates[0], minX, minY, maxX, maxY)
		if len(clippedRing) < 3 {
			return nil, fmt.Errorf("polygon clipped away")
		}
		clipped := [][]geojson.Position{clippedRing}
		for i := 1; i < len(g.Coordinates); i++ {
			hole := clipPolygonRing(g.Coordinates[i], minX, minY, maxX, maxY)
			if len(hole) >= 3 {
				clipped = append(clipped, hole)
			}
		}
		return geojson.NewFeature(geojson.NewPolygon(clipped), nil), nil
	case *geojson.MultiPolygon:
		var polys [][][]geojson.Position
		for _, poly := range g.Coordinates {
			clippedRing := clipPolygonRing(poly[0], minX, minY, maxX, maxY)
			if len(clippedRing) >= 3 {
				clipped := [][]geojson.Position{clippedRing}
				for i := 1; i < len(poly); i++ {
					hole := clipPolygonRing(poly[i], minX, minY, maxX, maxY)
					if len(hole) >= 3 {
						clipped = append(clipped, hole)
					}
				}
				polys = append(polys, clipped)
			}
		}
		if len(polys) == 0 {
			return nil, fmt.Errorf("polygons clipped away")
		}
		return geojson.NewFeature(geojson.NewMultiPolygon(polys), nil), nil
	default:
		return nil, fmt.Errorf("bbox-clip: unsupported geometry type %s", geom.Type())
	}
}

func clipLineString(line []geojson.Position, minX, minY, maxX, maxY float64) [][]geojson.Position {
	var segments [][]geojson.Position
	var current []geojson.Position

	for i := 0; i < len(line)-1; i++ {
		a, b := line[i], line[i+1]
		clipped := clipSegment(a, b, minX, minY, maxX, maxY)
		if len(clipped) == 0 {
			if len(current) > 0 {
				if len(current) >= 2 {
					segments = append(segments, current)
				}
				current = nil
			}
			continue
		}
		if len(current) == 0 {
			current = append(current, clipped[0])
		}
		current = append(current, clipped[1])
	}
	if len(current) >= 2 {
		segments = append(segments, current)
	}
	return segments
}

func clipSegment(a, b geojson.Position, minX, minY, maxX, maxY float64) []geojson.Position {
	outcode := func(p geojson.Position) int {
		c := 0
		if p[0] < minX {
			c |= 1
		}
		if p[0] > maxX {
			c |= 2
		}
		if p[1] < minY {
			c |= 4
		}
		if p[1] > maxY {
			c |= 8
		}
		return c
	}

	p0, p1 := a, b
	c0, c1 := outcode(p0), outcode(p1)

	for {
		if c0 == 0 && c1 == 0 {
			return []geojson.Position{p0, p1}
		}
		if c0&c1 != 0 {
			return nil
		}
		c := c0
		if c == 0 {
			c = c1
		}
		var p geojson.Position
		if c&1 != 0 {
			p = geojson.Position{minX, p0[1] + (p1[1]-p0[1])*(minX-p0[0])/(p1[0]-p0[0])}
		} else if c&2 != 0 {
			p = geojson.Position{maxX, p0[1] + (p1[1]-p0[1])*(maxX-p0[0])/(p1[0]-p0[0])}
		} else if c&4 != 0 {
			p = geojson.Position{p0[0] + (p1[0]-p0[0])*(minY-p0[1])/(p1[1]-p0[1]), minY}
		} else if c&8 != 0 {
			p = geojson.Position{p0[0] + (p1[0]-p0[0])*(maxY-p0[1])/(p1[1]-p0[1]), maxY}
		}
		if c == c0 {
			p0 = p
			c0 = outcode(p0)
		} else {
			p1 = p
			c1 = outcode(p1)
		}
	}
}

func clipPolygonRing(ring []geojson.Position, minX, minY, maxX, maxY float64) []geojson.Position {
	clipped := ring
	for _, clip := range []struct {
		test  func(geojson.Position) bool
		inter func(geojson.Position, geojson.Position) geojson.Position
	}{
		{
			func(p geojson.Position) bool { return p[0] >= minX },
			func(a, b geojson.Position) geojson.Position {
				t := (minX - a[0]) / (b[0] - a[0])
				return geojson.Position{minX, a[1] + t*(b[1]-a[1])}
			},
		},
		{
			func(p geojson.Position) bool { return p[0] <= maxX },
			func(a, b geojson.Position) geojson.Position {
				t := (maxX - a[0]) / (b[0] - a[0])
				return geojson.Position{maxX, a[1] + t*(b[1]-a[1])}
			},
		},
		{
			func(p geojson.Position) bool { return p[1] >= minY },
			func(a, b geojson.Position) geojson.Position {
				t := (minY - a[1]) / (b[1] - a[1])
				return geojson.Position{a[0] + t*(b[0]-a[0]), minY}
			},
		},
		{
			func(p geojson.Position) bool { return p[1] <= maxY },
			func(a, b geojson.Position) geojson.Position {
				t := (maxY - a[1]) / (b[1] - a[1])
				return geojson.Position{a[0] + t*(b[0]-a[0]), maxY}
			},
		},
	} {
		var next []geojson.Position
		for i := 0; i < len(clipped); i++ {
			cur := clipped[i]
			prev := clipped[(i+len(clipped)-1)%len(clipped)]
			curInside := clip.test(cur)
			prevInside := clip.test(prev)
			if curInside {
				if !prevInside {
					next = append(next, clip.inter(prev, cur))
				}
				next = append(next, cur)
			} else if prevInside {
				next = append(next, clip.inter(prev, cur))
			}
		}
		clipped = next
		if len(clipped) < 3 {
			return nil
		}
	}
	return clipped
}

func getAllCoords(obj any) ([]geojson.Position, error) {
	return geojson.CoordAll(obj)
}

var earthRadius = 6371008.0

func Distance(fromLng, fromLat, toLng, toLat float64) float64 {
	lat1 := fromLat * math.Pi / 180
	lat2 := toLat * math.Pi / 180
	dlat := lat2 - lat1
	dlon := (toLng - fromLng) * math.Pi / 180

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadius * c
}
