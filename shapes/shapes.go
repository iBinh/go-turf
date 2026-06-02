package shapes

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/ibinh/turf-go/geojson"
	"github.com/ibinh/turf-go/measurement"
)

type CircleOptions struct {
	Steps     int
	Units     measurement.Unit
	Properties map[string]any
}

func Circle(center any, radius float64, options ...CircleOptions) (*geojson.Feature, error) {
	opts := CircleOptions{Steps: 64, Units: measurement.UnitKilometers, Properties: map[string]any{}}
	if len(options) > 0 {
		opts = options[0]
	}
	return ellipse(center, radius, radius, opts.Steps, opts.Units, opts.Properties)
}

type EllipseOptions struct {
	Steps     int
	Units     measurement.Unit
	Angle     float64
	Properties map[string]any
}

func Ellipse(center any, xSemiAxis, ySemiAxis float64, options ...EllipseOptions) (*geojson.Feature, error) {
	opts := EllipseOptions{Steps: 64, Units: measurement.UnitKilometers, Angle: 0, Properties: map[string]any{}}
	if len(options) > 0 {
		opts = options[0]
	}

	coord, err := geojson.GetCoord(center)
	if err != nil {
		return nil, err
	}

	steps := opts.Steps
	if steps < 3 {
		steps = 3
	}

	ring := make([]geojson.Position, steps+1)
	for i := 0; i < steps; i++ {
		angle := float64(i) * 360.0 / float64(steps)
		dest1, _ := measurement.Destination(geojson.NewPoint(coord), xSemiAxis, angle, opts.Units)
		c1, _ := geojson.GetCoord(dest1)
		dest2, _ := measurement.Destination(geojson.NewPoint(coord), ySemiAxis, angle+90, opts.Units)
		c2, _ := geojson.GetCoord(dest2)

		lng := c1[0] - coord[0] + c2[0] - coord[0] + coord[0]
		lat := c1[1] - coord[1] + c2[1] - coord[1] + coord[1]

		if opts.Angle != 0 {
			rotated := rotateCoord(geojson.Position{lng, lat}, coord, opts.Angle)
			lng, lat = rotated[0], rotated[1]
		}

		ring[i] = geojson.Position{lng, lat}
	}
	ring[steps] = geojson.Position{ring[0][0], ring[0][1]}

	return geojson.NewFeature(geojson.NewPolygon([][]geojson.Position{ring}), opts.Properties), nil
}

func ellipse(center any, xSemiAxis, ySemiAxis float64, steps int, units measurement.Unit, properties map[string]any) (*geojson.Feature, error) {
	return Ellipse(center, xSemiAxis, ySemiAxis, EllipseOptions{Steps: steps, Units: units, Properties: properties})
}

func rotateCoord(p, center geojson.Position, angleDeg float64) geojson.Position {
	rad := angleDeg * math.Pi / 180
	cos := math.Cos(rad)
	sin := math.Sin(rad)
	dx := p[0] - center[0]
	dy := p[1] - center[1]
	return geojson.Position{
		center[0] + dx*cos - dy*sin,
		center[1] + dx*sin + dy*cos,
	}
}

type BezierOptions struct {
	Properties   map[string]any
	Sharpness    float64
	Resolution   int
}

func BezierSpline(line any, options ...BezierOptions) (*geojson.Feature, error) {
	opts := BezierOptions{Sharpness: 0.85, Resolution: 10000, Properties: map[string]any{}}
	if len(options) > 0 {
		opts = options[0]
	}

	geom, err := geojson.GetGeometry(line)
	if err != nil {
		return nil, err
	}
	ls, ok := geom.(*geojson.LineString)
	if !ok {
		return nil, fmt.Errorf("expected LineString, got %s", geom.Type())
	}
	pts := ls.Coordinates
	if len(pts) < 2 {
		return nil, fmt.Errorf("line must have at least 2 points")
	}

	sharpness := opts.Sharpness
	resolution := opts.Resolution
	if resolution < 2 {
		resolution = 2
	}

	n := len(pts)
	if n == 2 {
		result := make([]geojson.Position, resolution+1)
		for i := 0; i <= resolution; i++ {
			t := float64(i) / float64(resolution)
			result[i] = geojson.Position{
				pts[0][0] + (pts[1][0]-pts[0][0])*t,
				pts[0][1] + (pts[1][1]-pts[0][1])*t,
			}
		}
		return geojson.NewFeature(geojson.NewLineString(result), opts.Properties), nil
	}

	controlPoints := make([]geojson.Position, 0, n*3)
	controlPoints = append(controlPoints, pts[0])
	for i := 1; i < n; i++ {
		prev := pts[i-1]
		curr := pts[i]

		var prevPrev, next geojson.Position
		if i > 1 {
			prevPrev = pts[i-2]
		} else {
			prevPrev = pts[0]
		}
		if i < n-1 {
			next = pts[i+1]
		} else {
			next = pts[n-1]
		}

		dx1 := (curr[0] - prevPrev[0]) * sharpness
		dy1 := (curr[1] - prevPrev[1]) * sharpness
		dx2 := (next[0] - prev[0]) * sharpness
		dy2 := (next[1] - prev[1]) * sharpness

		controlPoints = append(controlPoints,
			geojson.Position{prev[0] + dx1, prev[1] + dy1},
			geojson.Position{curr[0] - dx2, curr[1] - dy2},
			curr,
		)
	}

	result := make([]geojson.Position, 0, resolution+1)
	totalSegments := float64((len(controlPoints) - 1) / 3)
	for i := 0; i <= resolution; i++ {
		t := float64(i) / float64(resolution)
		segIdx := int(t * totalSegments)
		if segIdx >= int(totalSegments) {
			segIdx = int(totalSegments) - 1
		}
		localT := t*totalSegments - float64(segIdx)
		idx := segIdx * 3
		if idx+3 > len(controlPoints)-1 {
			idx = len(controlPoints) - 4
			if idx < 0 {
				idx = 0
			}
		}
		if idx+3 < len(controlPoints) {
			p0 := controlPoints[idx]
			p1 := controlPoints[idx+1]
			p2 := controlPoints[idx+2]
			p3 := controlPoints[idx+3]
			pt := cubicBezier(p0, p1, p2, p3, localT)
			result = append(result, pt)
		}
	}

	return geojson.NewFeature(geojson.NewLineString(result), opts.Properties), nil
}

func cubicBezier(p0, p1, p2, p3 geojson.Position, t float64) geojson.Position {
	u := 1 - t
	tt := t * t
	uu := u * u
	uuu := uu * u
	ttt := tt * t
	return geojson.Position{
		uuu*p0[0] + 3*uu*t*p1[0] + 3*u*tt*p2[0] + ttt*p3[0],
		uuu*p0[1] + 3*uu*t*p1[1] + 3*u*tt*p2[1] + ttt*p3[1],
	}
}

type RandomOptions struct {
	BBox       []float64
	NumVertices int
	MaxLength  float64
	MaxRotation float64
	Properties map[string]any
}

func RandomPosition(bbox []float64) geojson.Position {
	if len(bbox) < 4 {
		return geojson.Position{(rand.Float64() - 0.5) * 360, (rand.Float64() - 0.5) * 180}
	}
	return geojson.Position{
		bbox[0] + rand.Float64()*(bbox[2]-bbox[0]),
		bbox[1] + rand.Float64()*(bbox[3]-bbox[1]),
	}
}

func RandomPoint(count int, options ...RandomOptions) (*geojson.FeatureCollection, error) {
	opts := getRandomOpts(options)
	features := make([]*geojson.Feature, count)
	for i := 0; i < count; i++ {
		pt := RandomPosition(opts.BBox)
		features[i] = geojson.NewFeature(geojson.NewPoint(pt), opts.Properties)
	}
	return geojson.NewFeatureCollection(features), nil
}

func RandomLineString(count int, options ...RandomOptions) (*geojson.FeatureCollection, error) {
	opts := getRandomOpts(options)
	features := make([]*geojson.Feature, count)
	for i := 0; i < count; i++ {
		vertices := opts.NumVertices
		if vertices < 2 {
			vertices = 2
		}
		coords := make([]geojson.Position, vertices)
		coords[0] = RandomPosition(opts.BBox)
		for j := 1; j < vertices; j++ {
			coords[j] = randomOffset(coords[j-1], opts)
		}
		features[i] = geojson.NewFeature(geojson.NewLineString(coords), opts.Properties)
	}
	return geojson.NewFeatureCollection(features), nil
}

func RandomPolygon(count int, options ...RandomOptions) (*geojson.FeatureCollection, error) {
	opts := getRandomOpts(options)
	features := make([]*geojson.Feature, count)
	for i := 0; i < count; i++ {
		vertices := opts.NumVertices
		if vertices < 3 {
			vertices = 3
		}
		center := RandomPosition(opts.BBox)
		coords := make([]geojson.Position, vertices+1)
		radius := opts.MaxLength
		if radius <= 0 {
			radius = 1
		}
		for j := 0; j < vertices; j++ {
			angle := float64(j) * 360.0 / float64(vertices)
			angle += (rand.Float64() - 0.5) * opts.MaxRotation
			rad := angle * math.Pi / 180
			r := radius * (0.5 + rand.Float64()*0.5)
			coords[j] = geojson.Position{
				center[0] + r*math.Cos(rad),
				center[1] + r*math.Sin(rad),
			}
		}
		coords[vertices] = geojson.Position{coords[0][0], coords[0][1]}
		features[i] = geojson.NewFeature(geojson.NewPolygon([][]geojson.Position{coords}), opts.Properties)
	}
	return geojson.NewFeatureCollection(features), nil
}

func getRandomOpts(options []RandomOptions) RandomOptions {
	opts := RandomOptions{
		NumVertices: 10,
		MaxLength:   1,
		MaxRotation: math.Pi / 8,
		Properties:  map[string]any{},
	}
	if len(options) > 0 {
		if options[0].NumVertices > 0 {
			opts.NumVertices = options[0].NumVertices
		}
		if options[0].MaxLength > 0 {
			opts.MaxLength = options[0].MaxLength
		}
		if options[0].MaxRotation > 0 {
			opts.MaxRotation = options[0].MaxRotation
		}
		opts.BBox = options[0].BBox
		opts.Properties = options[0].Properties
	}
	return opts
}

func randomOffset(origin geojson.Position, opts RandomOptions) geojson.Position {
	angle := rand.Float64() * 2 * math.Pi
	length := opts.MaxLength * (0.1 + rand.Float64()*0.9)
	return geojson.Position{
		origin[0] + length*math.Cos(angle),
		origin[1] + length*math.Sin(angle),
	}
}
