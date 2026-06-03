package lines

import (
	"fmt"
	"math"

	"github.com/ibinh/turf-go/boolean"
	"github.com/ibinh/turf-go/geojson"
	"github.com/ibinh/turf-go/measurement"
)

func LineIntersect(line1, line2 any) (*geojson.FeatureCollection, error) {
	segments1, err := extractSegments(line1)
	if err != nil {
		return nil, fmt.Errorf("lineIntersect: %w", err)
	}
	segments2, err := extractSegments(line2)
	if err != nil {
		return nil, fmt.Errorf("lineIntersect: %w", err)
	}

	var points []*geojson.Feature
	seen := make(map[string]bool)

	for _, segs1 := range segments1 {
		for i := 0; i < len(segs1)-1; i++ {
			a, b := segs1[i], segs1[i+1]
			for _, segs2 := range segments2 {
				for j := 0; j < len(segs2)-1; j++ {
					c, d := segs2[j], segs2[j+1]
					if pt, ok := segSegIntersect(a, b, c, d); ok {
						key := fmt.Sprintf("%.10f,%.10f", pt[0], pt[1])
						if !seen[key] {
							seen[key] = true
							points = append(points, geojson.NewFeature(geojson.NewPoint(pt), nil))
						}
					}
				}
			}
		}
	}

	return geojson.NewFeatureCollection(points), nil
}

func LineSegment(geom any) (*geojson.FeatureCollection, error) {
	g, err := geojson.GetGeometry(geom)
	if err != nil {
		return nil, err
	}

	collect := func(coords []geojson.Position) []*geojson.Feature {
		var segs []*geojson.Feature
		for i := 0; i < len(coords)-1; i++ {
			seg := geojson.NewLineString([]geojson.Position{coords[i], coords[i+1]})
			segs = append(segs, geojson.NewFeature(seg, nil))
		}
		return segs
	}

	var features []*geojson.Feature
	switch v := g.(type) {
	case *geojson.LineString:
		features = collect(v.Coordinates)
	case *geojson.MultiLineString:
		for _, line := range v.Coordinates {
			features = append(features, collect(line)...)
		}
	case *geojson.Polygon:
		for _, ring := range v.Coordinates {
			features = append(features, collect(ring)...)
		}
	case *geojson.MultiPolygon:
		for _, poly := range v.Coordinates {
			for _, ring := range poly {
				features = append(features, collect(ring)...)
			}
		}
	default:
		return nil, fmt.Errorf("lineSegment: unsupported geometry type %s", g.Type())
	}

	return geojson.NewFeatureCollection(features), nil
}

func LineOverlap(line1, line2 any) (*geojson.FeatureCollection, error) {
	coords1, err := geojson.GetCoords(line1)
	if err != nil {
		return nil, fmt.Errorf("lineOverlap: %w", err)
	}
	coords2, err := geojson.GetCoords(line2)
	if err != nil {
		return nil, fmt.Errorf("lineOverlap: %w", err)
	}
	pts1, ok1 := coords1.([]geojson.Position)
	pts2, ok2 := coords2.([]geojson.Position)
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("lineOverlap: expected LineString coordinates")
	}

	var features []*geojson.Feature
	seen := make(map[string]bool)

	for i := 0; i < len(pts1)-1; i++ {
		a, b := pts1[i], pts1[i+1]
		for j := 0; j < len(pts2)-1; j++ {
			c, d := pts2[j], pts2[j+1]
			seg := overlappingSegment(a, b, c, d)
			if seg == nil {
				continue
			}
			key := fmt.Sprintf("%.10f,%.10f-%.10f,%.10f", seg[0][0], seg[0][1], seg[1][0], seg[1][1])
			if !seen[key] {
				seen[key] = true
				revKey := fmt.Sprintf("%.10f,%.10f-%.10f,%.10f", seg[1][0], seg[1][1], seg[0][0], seg[0][1])
				seen[revKey] = true
				features = append(features, geojson.NewFeature(geojson.NewLineString(seg), nil))
			}
		}
	}

	return geojson.NewFeatureCollection(features), nil
}

func overlappingSegment(a, b, c, d geojson.Position) []geojson.Position {
	vx := b[0] - a[0]
	vy := b[1] - a[1]

	cross1 := (c[0]-a[0])*vy - (c[1]-a[1])*vx
	cross2 := (d[0]-a[0])*vy - (d[1]-a[1])*vx
	if math.Abs(cross1) > 1e-10 || math.Abs(cross2) > 1e-10 {
		return nil
	}

	lenSq := vx*vx + vy*vy
	if lenSq < 1e-15 {
		return nil
	}

	t1 := ((c[0]-a[0])*vx + (c[1]-a[1])*vy) / lenSq
	t2 := ((d[0]-a[0])*vx + (d[1]-a[1])*vy) / lenSq

	lo := math.Max(math.Min(t1, t2), 0.0)
	hi := math.Min(math.Max(t1, t2), 1.0)

	if lo >= hi || math.Abs(hi-lo) < 1e-15 {
		return nil
	}

	start := geojson.Position{a[0] + lo*vx, a[1] + lo*vy}
	end := geojson.Position{a[0] + hi*vx, a[1] + hi*vy}

	return []geojson.Position{start, end}
}

func LineSlice(point1, point2, line any) (*geojson.Feature, error) {
	coords, err := geojson.GetCoords(line)
	if err != nil {
		return nil, fmt.Errorf("lineSlice: %w", err)
	}
	pts, ok := coords.([]geojson.Position)
	if !ok {
		return nil, fmt.Errorf("lineSlice: expected LineString")
	}

	startNearest, err := measurement.NearestPointOnLine(line, point1)
	if err != nil {
		return nil, err
	}
	endNearest, err := measurement.NearestPointOnLine(line, point2)
	if err != nil {
		return nil, err
	}

	startPt, _ := geojson.GetCoord(startNearest)
	endPt, _ := geojson.GetCoord(endNearest)

	startIdx := toInt(startNearest.Properties["index"])
	endIdx := toInt(endNearest.Properties["index"])

	var result []geojson.Position
	if startIdx < endIdx {
		result = append(result, startPt)
		for i := startIdx + 1; i <= endIdx; i++ {
			result = append(result, pts[i])
		}
		result = append(result, endPt)
	} else if startIdx > endIdx {
		result = append(result, startPt)
		for i := startIdx + 1; i < len(pts); i++ {
			result = append(result, pts[i])
		}
		for i := 0; i <= endIdx; i++ {
			result = append(result, pts[i])
		}
		result = append(result, endPt)
	} else {
		dist1 := haversineDist(pts[startIdx], startPt)
		dist2 := haversineDist(pts[startIdx], endPt)
		segLen := haversineDist(pts[startIdx], pts[startIdx+1])
		if math.Abs(segLen) < 1e-15 {
			result = []geojson.Position{startPt, endPt}
		} else {
			if dist1 < dist2 {
				result = []geojson.Position{startPt, endPt}
			} else {
				result = []geojson.Position{endPt, startPt}
			}
		}
	}

	return geojson.NewFeature(geojson.NewLineString(result), nil), nil
}

func haversineDist(a, b geojson.Position) float64 {
	lat1 := a[1] * math.Pi / 180
	lat2 := b[1] * math.Pi / 180
	dlat := lat2 - lat1
	dlon := (b[0] - a[0]) * math.Pi / 180
	sinDlat := math.Sin(dlat / 2)
	sinDlon := math.Sin(dlon / 2)
	sinDLat2 := sinDlat * sinDlat
	sinDLon2 := sinDlon * sinDlon
	aVal := sinDLat2 + math.Cos(lat1)*math.Cos(lat2)*sinDLon2
	return 2 * math.Atan2(math.Sqrt(aVal), math.Sqrt(1-aVal)) * measurement.EarthRadius
}

func LineSliceAlong(line any, startDist, endDist float64, units measurement.Unit) (*geojson.Feature, error) {
	coords, err := geojson.GetCoords(line)
	if err != nil {
		return nil, fmt.Errorf("lineSliceAlong: %w", err)
	}
	pts, ok := coords.([]geojson.Position)
	if !ok {
		return nil, fmt.Errorf("lineSliceAlong: expected LineString")
	}

	startFeat, err := measurement.Along(line, startDist, units)
	if err != nil {
		return nil, err
	}
	endFeat, err := measurement.Along(line, endDist, units)
	if err != nil {
		return nil, err
	}

	startPt, _ := geojson.GetCoord(startFeat)
	endPt, _ := geojson.GetCoord(endFeat)

	targetStart := startDist
	targetEnd := endDist
	if units != measurement.UnitMeters {
		targetStart = measurement.ConvertLength(startDist, units, measurement.UnitMeters)
		targetEnd = measurement.ConvertLength(endDist, units, measurement.UnitMeters)
	}

	if targetStart > targetEnd {
		targetStart, targetEnd = targetEnd, targetStart
	}

	travelled := 0.0
	var lineCoords []geojson.Position
	lineCoords = append(lineCoords, startPt)

	for i := 0; i < len(pts)-1; i++ {
		segLen := haversineDist(pts[i], pts[i+1])
		segStart := travelled
		segEnd := travelled + segLen

		if segEnd <= targetStart || segStart >= targetEnd {
			travelled += segLen
			continue
		}

		if segStart >= targetStart && segEnd <= targetEnd {
			if !equalPos(pts[i+1], startPt) && !equalPos(pts[i+1], endPt) {
				lineCoords = append(lineCoords, pts[i+1])
			}
		}

		travelled += segLen
	}

	if !equalPos(lineCoords[len(lineCoords)-1], endPt) {
		lineCoords = append(lineCoords, endPt)
	}

	lineLen, _ := measurement.Length(line, measurement.UnitMeters)
	if lineLen < 1e-10 {
		return geojson.NewFeature(geojson.NewLineString([]geojson.Position{startPt, endPt}), nil), nil
	}

	return geojson.NewFeature(geojson.NewLineString(lineCoords), nil), nil
}

func equalPos(a, b geojson.Position) bool {
	return a[0] == b[0] && a[1] == b[1]
}

func LineChunk(line any, segmentLength float64, units measurement.Unit) (*geojson.FeatureCollection, error) {
	coords, err := geojson.GetCoords(line)
	if err != nil {
		return nil, fmt.Errorf("lineChunk: %w", err)
	}
	pts, ok := coords.([]geojson.Position)
	if !ok {
		return nil, fmt.Errorf("lineChunk: expected LineString")
	}
	if len(pts) < 2 {
		return nil, fmt.Errorf("lineChunk: line must have at least 2 vertices")
	}

	totalLen, err := measurement.Length(line, measurement.UnitMeters)
	if err != nil {
		return nil, err
	}

	segLenMeters := measurement.ConvertLength(segmentLength, units, measurement.UnitMeters)

	var chunks []*geojson.Feature
	currentDist := 0.0

	for currentDist < totalLen-0.01 {
		nextDist := currentDist + segLenMeters
		if nextDist > totalLen {
			nextDist = totalLen
		}
		chunk, err := LineSliceAlong(line, currentDist, nextDist, measurement.UnitMeters)
		if err != nil {
			return nil, err
		}
		chunks = append(chunks, chunk)
		currentDist = nextDist
		if math.Abs(currentDist-totalLen) < 0.01 {
			break
		}
	}

	return geojson.NewFeatureCollection(chunks), nil
}

func LineSplit(line any, point any) (*geojson.FeatureCollection, error) {
	coords, err := geojson.GetCoords(line)
	if err != nil {
		return nil, fmt.Errorf("lineSplit: %w", err)
	}
	pts, ok := coords.([]geojson.Position)
	if !ok {
		return nil, fmt.Errorf("lineSplit: expected LineString")
	}

	ptCoord, err := geojson.GetCoord(point)
	if err != nil {
		return nil, err
	}

	onLine, err := boolean.PointOnLine(point, line, false)
	if err != nil {
		return nil, err
	}
	if !onLine {
		return nil, fmt.Errorf("lineSplit: point must be on the line")
	}

	if equalPos(pts[0], ptCoord) || equalPos(pts[len(pts)-1], ptCoord) {
		f := geojson.NewFeature(geojson.NewLineString(pts), nil)
		return geojson.NewFeatureCollection([]*geojson.Feature{f}), nil
	}

	for i := 0; i < len(pts)-1; i++ {
		if equalPos(pts[i], ptCoord) {
			line1 := make([]geojson.Position, i+1)
			copy(line1, pts[:i+1])
			line2 := make([]geojson.Position, len(pts)-i)
			copy(line2, pts[i:])
			return geojson.NewFeatureCollection([]*geojson.Feature{
				geojson.NewFeature(geojson.NewLineString(line1), nil),
				geojson.NewFeature(geojson.NewLineString(line2), nil),
			}), nil
		}
	}

	for i := 0; i < len(pts)-1; i++ {
		if isPointOnSegment(ptCoord, pts[i], pts[i+1]) {
			line1 := make([]geojson.Position, 0, i+2)
			line1 = append(line1, pts[:i+1]...)
			line1 = append(line1, ptCoord)
			line2 := make([]geojson.Position, 0, len(pts)-i+1)
			line2 = append(line2, ptCoord)
			line2 = append(line2, pts[i+1:]...)
			return geojson.NewFeatureCollection([]*geojson.Feature{
				geojson.NewFeature(geojson.NewLineString(line1), nil),
				geojson.NewFeature(geojson.NewLineString(line2), nil),
			}), nil
		}
	}

	return nil, fmt.Errorf("lineSplit: could not split line")
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
	return dot-lenSq <= 1e-10
}

func LineArc(center any, radius float64, bearing1, bearing2 float64, steps int, units measurement.Unit) (*geojson.Feature, error) {
	if steps < 1 {
		steps = 64
	}

	var pts []geojson.Position
	if bearing1 <= bearing2 {
		for i := 0; i <= steps; i++ {
			bearing := bearing1 + (bearing2-bearing1)*float64(i)/float64(steps)
			pt, err := measurement.Destination(center, radius, bearing, units)
			if err != nil {
				return nil, err
			}
			coord, _ := geojson.GetCoord(pt)
			pts = append(pts, coord)
		}
	} else {
		for i := 0; i <= steps; i++ {
			bearing := bearing1 + (360+bearing2-bearing1)*float64(i)/float64(steps)
			if bearing >= 360 {
				bearing -= 360
			}
			pt, err := measurement.Destination(center, radius, bearing, units)
			if err != nil {
				return nil, err
			}
			coord, _ := geojson.GetCoord(pt)
			pts = append(pts, coord)
		}
	}

	return geojson.NewFeature(geojson.NewLineString(pts), nil), nil
}

func Sector(center any, radius float64, bearing1, bearing2 float64, steps int, units measurement.Unit) (*geojson.Feature, error) {
	if steps < 1 {
		steps = 64
	}

	centerCoord, err := geojson.GetCoord(center)
	if err != nil {
		return nil, fmt.Errorf("sector: %w", err)
	}

	var pts []geojson.Position
	if bearing1 <= bearing2 {
		for i := 0; i <= steps; i++ {
			bearing := bearing1 + (bearing2-bearing1)*float64(i)/float64(steps)
			pt, err := measurement.Destination(center, radius, bearing, units)
			if err != nil {
				return nil, err
			}
			coord, _ := geojson.GetCoord(pt)
			pts = append(pts, coord)
		}
	} else {
		for i := 0; i <= steps; i++ {
			bearing := bearing1 + (360+bearing2-bearing1)*float64(i)/float64(steps)
			if bearing >= 360 {
				bearing -= 360
			}
			pt, err := measurement.Destination(center, radius, bearing, units)
			if err != nil {
				return nil, err
			}
			coord, _ := geojson.GetCoord(pt)
			pts = append(pts, coord)
		}
	}

	ring := append([]geojson.Position{centerCoord}, pts...)
	ring = append(ring, centerCoord)

	return geojson.NewFeature(geojson.NewPolygon([][]geojson.Position{ring}), nil), nil
}

func toInt(v any) int {
	switch i := v.(type) {
	case int:
		return i
	case float64:
		return int(i)
	default:
		return 0
	}
}

func extractSegments(geom any) ([][]geojson.Position, error) {
	g, err := geojson.GetGeometry(geom)
	if err != nil {
		return nil, err
	}
	switch v := g.(type) {
	case *geojson.LineString:
		return [][]geojson.Position{v.Coordinates}, nil
	case *geojson.MultiLineString:
		return v.Coordinates, nil
	case *geojson.Polygon:
		return v.Coordinates, nil
	case *geojson.MultiPolygon:
		var all [][]geojson.Position
		for _, poly := range v.Coordinates {
			all = append(all, poly...)
		}
		return all, nil
	default:
		return nil, fmt.Errorf("unsupported geometry type %s", g.Type())
	}
}

func segSegIntersect(a, b, c, d geojson.Position) (geojson.Position, bool) {
	den := (b[0]-a[0])*(d[1]-c[1]) - (b[1]-a[1])*(d[0]-c[0])
	if math.Abs(den) < 1e-15 {
		return geojson.Position{}, false
	}
	t := ((c[0]-a[0])*(d[1]-c[1]) - (c[1]-a[1])*(d[0]-c[0])) / den
	u := ((c[0]-a[0])*(b[1]-a[1]) - (c[1]-a[1])*(b[0]-a[0])) / den
	if t < 0 || t > 1 || u < 0 || u > 1 {
		return geojson.Position{}, false
	}
	return geojson.Position{a[0] + t*(b[0]-a[0]), a[1] + t*(b[1]-a[1])}, true
}
