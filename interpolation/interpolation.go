package interpolation

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/ibinh/turf-go/geojson"
	"github.com/ibinh/turf-go/measurement"
	"github.com/ibinh/turf-go/meta"
)

func Sample(fc *geojson.FeatureCollection, n int) (*geojson.FeatureCollection, error) {
	if fc == nil {
		return nil, fmt.Errorf("feature collection is required")
	}
	total := len(fc.Features)
	if total == 0 {
		return geojson.NewFeatureCollection(nil), nil
	}
	if n >= total {
		n = total
	}
	indices := rand.Perm(total)
	result := make([]*geojson.Feature, n)
	for i := 0; i < n; i++ {
		result[i] = fc.Features[indices[i]]
	}
	return geojson.NewFeatureCollection(result), nil
}

type TinOptions struct {
	Properties map[string]any
}

func Tin(points *geojson.FeatureCollection, options ...TinOptions) (*geojson.FeatureCollection, error) {
	opts := TinOptions{Properties: map[string]any{}}
	if len(options) > 0 {
		opts = options[0]
	}
	if points == nil || len(points.Features) < 3 {
		return nil, fmt.Errorf("need at least 3 points for TIN")
	}

	pts := make([]geojson.Position, len(points.Features))
	for i, f := range points.Features {
		coord, err := geojson.GetCoord(f)
		if err != nil {
			return nil, err
		}
		pts[i] = coord
	}

	triangles := delaunayTriangulation(pts)

	props := opts.Properties
	features := make([]*geojson.Feature, len(triangles))
	for i, tri := range triangles {
		ring := []geojson.Position{
			tri[0], tri[1], tri[2], tri[0],
		}
		features[i] = geojson.NewFeature(geojson.NewPolygon([][]geojson.Position{ring}), props)
	}
	return geojson.NewFeatureCollection(features), nil
}

func delaunayTriangulation(pts []geojson.Position) [][3]geojson.Position {
	minX, minY := pts[0][0], pts[0][1]
	maxX, maxY := minX, minY
	for _, p := range pts {
		if p[0] < minX {
			minX = p[0]
		}
		if p[0] > maxX {
			maxX = p[0]
		}
		if p[1] < minY {
			minY = p[1]
		}
		if p[1] > maxY {
			maxY = p[1]
		}
	}

	dx := maxX - minX
	dy := maxY - minY
	if dx == 0 {
		dx = 1
	}
	if dy == 0 {
		dy = 1
	}
	d := math.Max(dx, dy) * 20

	super := [3]geojson.Position{
		{minX - d, minY - d},
		{maxX + d, minY - d},
		{minX + dx/2, maxY + d},
	}

	triangles := [][3]geojson.Position{super}

	for _, pt := range pts {
		var badTriangles []int
		for i, tri := range triangles {
			if inCircumcircle(pt, tri) {
				badTriangles = append(badTriangles, i)
			}
		}

		polygon := []edge{}
		for _, idx := range badTriangles {
			tri := triangles[idx]
			for i := 0; i < 3; i++ {
				e := edge{tri[i], tri[(i+1)%3]}
				shared := false
				for j, pe := range polygon {
					if e.equals(pe) {
						polygon = append(polygon[:j], polygon[j+1:]...)
						shared = true
						break
					}
				}
				if !shared {
					polygon = append(polygon, e)
				}
			}
		}

		newTriangles := make([][3]geojson.Position, 0, len(badTriangles)*2)
		for i, tri := range triangles {
			keep := true
			for _, b := range badTriangles {
				if i == b {
					keep = false
					break
				}
			}
			if keep {
				newTriangles = append(newTriangles, tri)
			}
		}

		for _, e := range polygon {
			newTriangles = append(newTriangles, [3]geojson.Position{e.a, e.b, pt})
		}

		triangles = newTriangles
	}

	result := make([][3]geojson.Position, 0, len(triangles))
	for _, tri := range triangles {
		if !hasSuperVertex(tri, super) {
			result = append(result, tri)
		}
	}
	return result
}

type edge struct {
	a, b geojson.Position
}

func (e edge) equals(other edge) bool {
	return (e.a[0] == other.a[0] && e.a[1] == other.a[1] && e.b[0] == other.b[0] && e.b[1] == other.b[1]) ||
		(e.a[0] == other.b[0] && e.a[1] == other.b[1] && e.b[0] == other.a[0] && e.b[1] == other.a[1])
}

func inCircumcircle(p geojson.Position, tri [3]geojson.Position) bool {
	ax, ay := tri[0][0], tri[0][1]
	bx, by := tri[1][0], tri[1][1]
	cx, cy := tri[2][0], tri[2][1]
	dx, dy := p[0], p[1]

	d := 2 * (ax*(by-cy) + bx*(cy-ay) + cx*(ay-by))
	if math.Abs(d) < 1e-15 {
		return false
	}

	ux := ((ax*ax+ay*ay)*(by-cy) + (bx*bx+by*by)*(cy-ay) + (cx*cx+cy*cy)*(ay-by)) / d
	uy := ((ax*ax+ay*ay)*(cx-bx) + (bx*bx+by*by)*(ax-cx) + (cx*cx+cy*cy)*(bx-ax)) / d

	r2 := (ax-ux)*(ax-ux) + (ay-uy)*(ay-uy)
	pd2 := (dx-ux)*(dx-ux) + (dy-uy)*(dy-uy)

	return pd2 <= r2+1e-10
}

func hasSuperVertex(tri [3]geojson.Position, super [3]geojson.Position) bool {
	for _, v := range tri {
		for _, sv := range super {
			if v[0] == sv[0] && v[1] == sv[1] {
				return true
			}
		}
	}
	return false
}

type InterpolateOptions struct {
	Weight      float64
	Properties  map[string]any
}

func Interpolate(points *geojson.FeatureCollection, cellSide float64, units measurement.Unit, property string, options ...InterpolateOptions) (*geojson.FeatureCollection, error) {
	opts := InterpolateOptions{Weight: 2, Properties: map[string]any{}}
	if len(options) > 0 {
		opts = options[0]
	}
	if points == nil || len(points.Features) < 2 {
		return nil, fmt.Errorf("need at least 2 points for interpolation")
	}

	bbox, err := pointsBBox(points)
	if err != nil {
		return nil, err
	}

	cellSide = measurement.ConvertLength(cellSide, units, measurement.UnitDegrees)
	if cellSide <= 0 {
		return nil, fmt.Errorf("cellSide must be positive")
	}

	centerY := (bbox[1] + bbox[3]) / 2
	cellSide = cellSide / math.Cos(centerY*math.Pi/180)
	if math.IsInf(cellSide, 0) || math.IsNaN(cellSide) || cellSide <= 0 {
		cellSide = 0.01
	}

	var data []struct {
		pos geojson.Position
		val float64
	}
	err = meta.FeatureEach(points, func(f *geojson.Feature, _ int) error {
		coord, err := geojson.GetCoord(f)
		if err != nil {
			return nil
		}
		val, ok := f.Properties[property]
		if !ok {
			return nil
		}
		v, ok := val.(float64)
		if !ok {
			return nil
		}
		data = append(data, struct {
			pos geojson.Position
			val float64
		}{coord, v})
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(data) < 2 {
		return nil, fmt.Errorf("at least 2 points with property %q required", property)
	}

	var features []*geojson.Feature
	for x := bbox[0]; x <= bbox[2]; x += cellSide {
		for y := bbox[1]; y <= bbox[3]; y += cellSide {
			val := idw(geojson.Position{x, y}, data, opts.Weight)
			props := make(map[string]any)
			for k, v := range opts.Properties {
				props[k] = v
			}
			props[property] = val
			features = append(features, geojson.NewFeature(geojson.NewPoint(geojson.Position{x, y}), props))
		}
	}

	return geojson.NewFeatureCollection(features), nil
}

func idw(pt geojson.Position, data []struct {
	pos geojson.Position
	val float64
}, weight float64) float64 {
	var sumWeight, sumValue float64
	for _, d := range data {
		dist := math.Sqrt((pt[0]-d.pos[0])*(pt[0]-d.pos[0]) + (pt[1]-d.pos[1])*(pt[1]-d.pos[1]))
		if dist < 1e-15 {
			return d.val
		}
		w := 1.0 / math.Pow(dist, weight)
		sumWeight += w
		sumValue += w * d.val
	}
	if sumWeight < 1e-15 {
		return 0
	}
	return sumValue / sumWeight
}

func pointsBBox(fc *geojson.FeatureCollection) ([]float64, error) {
	var minX, minY, maxX, maxY float64
	first := true
	err := meta.CoordEach(fc, func(c geojson.Position, _ int) error {
		if first {
			minX, minY, maxX, maxY = c[0], c[1], c[0], c[1]
			first = false
		} else {
			if c[0] < minX {
				minX = c[0]
			}
			if c[0] > maxX {
				maxX = c[0]
			}
			if c[1] < minY {
				minY = c[1]
			}
			if c[1] > maxY {
				maxY = c[1]
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if first {
		return nil, fmt.Errorf("no coordinates found")
	}
	return []float64{minX, minY, maxX, maxY}, nil
}

func PlanarDistance(a, b geojson.Position) float64 {
	dx := a[0] - b[0]
	dy := a[1] - b[1]
	return math.Sqrt(dx*dx + dy*dy)
}

func PlanarPointOnLine(pt, start, end geojson.Position) bool {
	dx := end[0] - start[0]
	dy := end[1] - start[1]
	cross := (pt[0]-start[0])*dy - (pt[1]-start[1])*dx
	if math.Abs(cross) > 1e-10 {
		return false
	}
	dot := (pt[0]-start[0])*dx + (pt[1]-start[1])*dy
	if dot < 0 {
		return false
	}
	lenSq := dx*dx + dy*dy
	return dot <= lenSq
}
