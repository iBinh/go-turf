package lineoffset

import (
	"fmt"
	"math"

	"github.com/ibinh/turf-go/geojson"
	"github.com/ibinh/turf-go/measurement"
)

func LineOffset(line any, distance float64, units measurement.Unit) (*geojson.Feature, error) {
	coords, err := geojson.GetCoords(line)
	if err != nil {
		return nil, fmt.Errorf("lineOffset: %w", err)
	}

	pts, ok := coords.([]geojson.Position)
	if !ok {
		return nil, fmt.Errorf("lineOffset: expected LineString coordinates")
	}

	if len(pts) < 2 {
		return nil, fmt.Errorf("lineOffset: line must have at least 2 points")
	}

	// Normalize the sign: positive moves right, negative moves left
	offset := -distance
	if units != measurement.UnitMeters {
		offset = measurement.ConvertLength(distance, units, measurement.UnitMeters)
		offset = -offset
	}

	offsetDeg := offset / (measurement.EarthRadius * math.Pi / 180)

	var result []geojson.Position
	if len(pts) == 2 {
		dx := pts[1][0] - pts[0][0]
		dy := pts[1][1] - pts[0][1]
		segLen := math.Sqrt(dx*dx + dy*dy)
		if segLen < 1e-15 {
			return nil, fmt.Errorf("lineOffset: degenerate segment")
		}
		nx := -dy / segLen * offsetDeg
		ny := dx / segLen * offsetDeg
		latRad := pts[0][1] * math.Pi / 180
		if math.Abs(latRad) < math.Pi/2*0.999 {
			nx /= math.Cos(latRad)
		}
		start := geojson.Position{pts[0][0] + nx, pts[0][1] + ny}
		latRad2 := pts[1][1] * math.Pi / 180
		nx2 := -dy / segLen * offsetDeg
		ny2 := dx / segLen * offsetDeg
		if math.Abs(latRad2) < math.Pi/2*0.999 {
			nx2 /= math.Cos(latRad2)
		}
		end := geojson.Position{pts[1][0] + nx2, pts[1][1] + ny2}
		result = []geojson.Position{start, end}
	} else {
		result = offsetPolyline(pts, offsetDeg)
	}

	return geojson.NewFeature(geojson.NewLineString(result), nil), nil
}

func offsetPolyline(pts []geojson.Position, offsetDeg float64) []geojson.Position {
	n := len(pts)
	type offsetSeg struct {
		start geojson.Position
		end   geojson.Position
	}
	segments := make([]offsetSeg, n-1)

	for i := 0; i < n-1; i++ {
		a, b := pts[i], pts[i+1]
		dx := b[0] - a[0]
		dy := b[1] - a[1]
		segLen := math.Sqrt(dx*dx + dy*dy)
		if segLen < 1e-15 {
			segments[i] = offsetSeg{start: a, end: a}
			continue
		}

		latRad := a[1] * math.Pi / 180
		lonScale := 1.0
		if math.Abs(latRad) < math.Pi/2*0.999 {
			lonScale = 1.0 / math.Cos(latRad)
		}

		nx := -dy / segLen * offsetDeg * lonScale
		ny := dx / segLen * offsetDeg

		startOff := geojson.Position{a[0] + nx, a[1] + ny}

		latRad2 := b[1] * math.Pi / 180
		lonScale2 := 1.0
		if math.Abs(latRad2) < math.Pi/2*0.999 {
			lonScale2 = 1.0 / math.Cos(latRad2)
		}
		nx2 := -dy / segLen * offsetDeg * lonScale2
		ny2 := dx / segLen * offsetDeg

		endOff := geojson.Position{b[0] + nx2, b[1] + ny2}

		segments[i] = offsetSeg{start: startOff, end: endOff}
	}

	result := make([]geojson.Position, 0, n)
	result = append(result, segments[0].start)

	for i := 0; i < n-2; i++ {
		intersection := miterIntersection(
			segments[i].start, segments[i].end,
			segments[i+1].start, segments[i+1].end,
		)
		if intersection != nil {
			// Check miter limit
			seg := segments[i]
			segVecLen := math.Sqrt(
				(seg.end[0]-seg.start[0])*(seg.end[0]-seg.start[0]) +
					(seg.end[1]-seg.start[1])*(seg.end[1]-seg.start[1]),
			)
			dist := math.Sqrt(
				((*intersection)[0]-seg.end[0])*((*intersection)[0]-seg.end[0]) +
					((*intersection)[1]-seg.end[1])*((*intersection)[1]-seg.end[1]),
			)
			miterLimit := math.Abs(offsetDeg) * 2.0
			if segVecLen > 1e-15 && dist/math.Max(segVecLen, 1e-15) > miterLimit {
				result = append(result, segments[i].end)
				result = append(result, segments[i+1].start)
			} else {
				result = append(result, *intersection)
			}
		} else {
			result = append(result, segments[i].end)
			result = append(result, segments[i+1].start)
		}
	}

	result = append(result, segments[n-2].end)

	return result
}

func miterIntersection(a1, a2, b1, b2 geojson.Position) *geojson.Position {
	dxA := a2[0] - a1[0]
	dyA := a2[1] - a1[1]
	dxB := b2[0] - b1[0]
	dyB := b2[1] - b1[1]

	den := dxA*dyB - dyA*dxB
	if math.Abs(den) < 1e-15 {
		return nil
	}

	t := ((b1[0]-a1[0])*dyB - (b1[1]-a1[1])*dxB) / den

	return &geojson.Position{a1[0] + t*dxA, a1[1] + t*dyA}
}
