package buffer

import (
	"fmt"
	"math"

	"github.com/ibinh/turf-go/geojson"
	"github.com/ibinh/turf-go/measurement"
	"github.com/ibinh/turf-go/polyclip"
	"github.com/ibinh/turf-go/shapes"
)

func Buffer(geom any, radius float64, units measurement.Unit, steps ...int) (*geojson.Feature, error) {
	if radius <= 0 {
		return nil, fmt.Errorf("radius must be positive")
	}

	nSteps := 16
	if len(steps) > 0 && steps[0] > 2 {
		nSteps = steps[0]
	}

	g, err := geojson.GetGeometry(geom)
	if err != nil {
		return nil, err
	}

	switch g := g.(type) {
	case *geojson.Point:
		return shapes.Circle(g, radius, shapes.CircleOptions{
			Steps: nSteps,
			Units: units,
		})

	case *geojson.MultiPoint:
		var all []*geojson.Polygon
		for _, pt := range g.Coordinates {
			c, err := shapes.Circle(geojson.NewPoint(pt), radius, shapes.CircleOptions{
				Steps: nSteps,
				Units: units,
			})
			if err != nil {
				return nil, err
			}
			all = append(all, c.Geometry.(*geojson.Polygon))
		}
		return unionAll(all)

	case *geojson.LineString:
		radiusDeg := measurement.ConvertLength(radius, units, measurement.UnitDegrees)
		return bufferLine(g.Coordinates, radiusDeg, nSteps)

	case *geojson.MultiLineString:
		radiusDeg := measurement.ConvertLength(radius, units, measurement.UnitDegrees)
		var all []*geojson.Polygon
		for _, coords := range g.Coordinates {
			f, err := bufferLine(coords, radiusDeg, nSteps)
			if err != nil {
				return nil, err
			}
			all = append(all, f.Geometry.(*geojson.Polygon))
		}
		return unionAll(all)

	case *geojson.Polygon:
		radiusDeg := measurement.ConvertLength(radius, units, measurement.UnitDegrees)
		return bufferPolygon(g, radiusDeg, nSteps)

	case *geojson.MultiPolygon:
		radiusDeg := measurement.ConvertLength(radius, units, measurement.UnitDegrees)
		var all []*geojson.Polygon
		for _, rings := range g.Coordinates {
			sub := geojson.NewPolygon(rings)
			f, err := bufferPolygon(sub, radiusDeg, nSteps)
			if err != nil {
				return nil, err
			}
			all = append(all, f.Geometry.(*geojson.Polygon))
		}
		return unionAll(all)

	default:
		return nil, fmt.Errorf("buffer not supported for geometry type %s", g.Type())
	}
}

func bufferLine(coords []geojson.Position, radius float64, steps int) (*geojson.Feature, error) {
	if len(coords) < 2 {
		return nil, fmt.Errorf("line must have at least 2 coordinates")
	}

	var pieces []*geojson.Polygon

	for i := 0; i < len(coords)-1; i++ {
		p1, p2 := coords[i], coords[i+1]
		rect := segmentBuffer(p1, p2, radius)
		if rect != nil {
			pieces = append(pieces, rect)
		}
	}

	for _, pt := range coords {
		circle := vertexCircle(pt, radius, steps)
		if circle != nil {
			pieces = append(pieces, circle)
		}
	}

	if len(pieces) == 0 {
		return nil, nil
	}
	return unionAll(pieces)
}

func bufferPolygon(poly *geojson.Polygon, radius float64, steps int) (*geojson.Feature, error) {
	var pieces []*geojson.Polygon

	for _, ring := range poly.Coordinates {
		for i := 0; i < len(ring)-1; i++ {
			p1, p2 := ring[i], ring[i+1]
			rect := segmentBuffer(p1, p2, radius)
			if rect != nil {
				pieces = append(pieces, rect)
			}
		}
		for _, pt := range ring {
			circle := vertexCircle(pt, radius, steps)
			if circle != nil {
				pieces = append(pieces, circle)
			}
		}
	}

	pieces = append(pieces, poly)

	if len(pieces) == 0 {
		return nil, nil
	}
	return unionAll(pieces)
}

func segmentBuffer(p1, p2 geojson.Position, radius float64) *geojson.Polygon {
	dx := p2[0] - p1[0]
	dy := p2[1] - p1[1]
	length := math.Sqrt(dx*dx + dy*dy)
	if length < 1e-15 {
		return nil
	}
	nx := -dy / length * radius
	ny := dx / length * radius
	return geojson.NewPolygon([][]geojson.Position{
		{
			{p1[0] + nx, p1[1] + ny},
			{p2[0] + nx, p2[1] + ny},
			{p2[0] - nx, p2[1] - ny},
			{p1[0] - nx, p1[1] - ny},
			{p1[0] + nx, p1[1] + ny},
		},
	})
}

func vertexCircle(center geojson.Position, radius float64, steps int) *geojson.Polygon {
	if steps < 3 {
		steps = 3
	}
	pts := make([]geojson.Position, steps+1)
	for i := 0; i <= steps; i++ {
		angle := float64(i) * 2 * math.Pi / float64(steps)
		pts[i] = geojson.Position{
			center[0] + radius*math.Cos(angle),
			center[1] + radius*math.Sin(angle),
		}
	}
	return geojson.NewPolygon([][]geojson.Position{pts})
}

func unionAll(pieces []*geojson.Polygon) (*geojson.Feature, error) {
	if len(pieces) == 0 {
		return nil, nil
	}

	result := geojson.NewFeature(pieces[0], nil)
	for i := 1; i < len(pieces); i++ {
		var err error
		result, err = polyclip.PolygonUnion(result, geojson.NewFeature(pieces[i], nil))
		if err != nil {
			return nil, err
		}
		if result == nil {
			result = geojson.NewFeature(pieces[i], nil)
		}
	}
	return result, nil
}
