package voronoi

import (
	"fmt"
	"math"
	"sort"

	"github.com/ibinh/turf-go/geojson"
)

type VoronoiOptions struct {
	BBox []float64
}

func Voronoi(points *geojson.FeatureCollection, options ...VoronoiOptions) (*geojson.FeatureCollection, error) {
	opts := VoronoiOptions{}
	if len(options) > 0 {
		opts = options[0]
	}
	if len(opts.BBox) < 4 {
		return nil, fmt.Errorf("bbox is required with [minX, minY, maxX, maxY]")
	}

	var seeds []geojson.Position
	for _, f := range points.Features {
		if f.Geometry != nil && f.Geometry.Type() == geojson.TypePoint {
			pt := f.Geometry.(*geojson.Point)
			seeds = append(seeds, pt.Coordinates)
		}
	}
	if len(seeds) == 0 {
		return nil, fmt.Errorf("no point features found")
	}

	sort.Slice(seeds, func(i, j int) bool {
		return seeds[i][0] < seeds[j][0]
	})

	features := make([]*geojson.Feature, 0, len(seeds))
	for _, seed := range seeds {
		cell := createCell(seed, seeds, opts.BBox)
		if cell != nil && len(cell.Coordinates) > 0 && len(cell.Coordinates[0]) >= 3 {
			features = append(features, geojson.NewFeature(cell, nil))
		}
	}

	return geojson.NewFeatureCollection(features), nil
}

func createCell(seed geojson.Position, allSeeds []geojson.Position, bbox []float64) *geojson.Polygon {
	poly := []geojson.Position{
		{bbox[0], bbox[1]},
		{bbox[2], bbox[1]},
		{bbox[2], bbox[3]},
		{bbox[0], bbox[3]},
		{bbox[0], bbox[1]},
	}

	for _, other := range allSeeds {
		if math.Abs(seed[0]-other[0]) < 1e-15 && math.Abs(seed[1]-other[1]) < 1e-15 {
			continue
		}

		mx := (seed[0] + other[0]) / 2
		my := (seed[1] + other[1]) / 2
		dx := other[0] - seed[0]
		dy := other[1] - seed[1]

		a := geojson.Position{mx, my}
		b := geojson.Position{mx - dy, my + dx}

		poly = clipPolygonByHalfPlane(poly, a, b)

		if len(poly) < 3 {
			return nil
		}
	}

	return geojson.NewPolygon([][]geojson.Position{poly})
}

func clipPolygonByHalfPlane(poly []geojson.Position, a, b geojson.Position) []geojson.Position {
	n := len(poly)
	if n < 3 {
		return nil
	}

	var result []geojson.Position
	for i := 0; i < n; i++ {
		cur := poly[i]
		nxt := poly[(i+1)%n]

		curCross := cross(a, b, cur)
		nxtCross := cross(a, b, nxt)

		curInside := curCross >= -1e-10
		nxtInside := nxtCross >= -1e-10

		if curInside {
			result = append(result, cur)
		}

		if curInside != nxtInside {
			inter := intersectSegLine(a, b, cur, nxt)
			result = append(result, inter)
		}
	}

	return result
}

func cross(a, b, p geojson.Position) float64 {
	return (b[0]-a[0])*(p[1]-a[1]) - (b[1]-a[1])*(p[0]-a[0])
}

func intersectSegLine(a, b, p1, p2 geojson.Position) geojson.Position {
	dxL := b[0] - a[0]
	dyL := b[1] - a[1]
	dxS := p2[0] - p1[0]
	dyS := p2[1] - p1[1]

	den := dxL*dyS - dyL*dxS
	if math.Abs(den) < 1e-15 {
		return p1
	}

	t := ((p1[0]-a[0])*dyS - (p1[1]-a[1])*dxS) / den
	return geojson.Position{
		a[0] + t*dxL,
		a[1] + t*dyL,
	}
}
