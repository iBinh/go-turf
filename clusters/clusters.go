package clusters

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/ibinh/turf-go/geojson"
)

func ClustersKMeans(fc *geojson.FeatureCollection, k int) (*geojson.FeatureCollection, error) {
	if fc == nil || len(fc.Features) == 0 {
		return nil, fmt.Errorf("feature collection is required")
	}
	if k < 1 {
		return nil, fmt.Errorf("k must be at least 1")
	}
	if k > len(fc.Features) {
		k = len(fc.Features)
	}

	pts := make([]geojson.Position, len(fc.Features))
	for i, f := range fc.Features {
		coord, err := geojson.GetCoord(f)
		if err != nil {
			return nil, err
		}
		pts[i] = coord
	}

	centroids := make([]geojson.Position, k)
	for i := 0; i < k; i++ {
		centroids[i] = pts[rand.Intn(len(pts))]
	}

	assignments := make([]int, len(pts))
	iterations := 100

	for iter := 0; iter < iterations; iter++ {
		changed := false
		for i, pt := range pts {
			best := 0
			bestDist := math.MaxFloat64
			for j, c := range centroids {
				d := sqDist(pt, c)
				if d < bestDist {
					bestDist = d
					best = j
				}
			}
			if assignments[i] != best {
				assignments[i] = best
				changed = true
			}
		}

		if !changed {
			break
		}

		newCentroids := make([]geojson.Position, k)
		for j := range newCentroids {
			newCentroids[j] = geojson.Position{0, 0}
		}
		counts := make([]int, k)
		for i, pt := range pts {
			c := assignments[i]
			newCentroids[c][0] += pt[0]
			newCentroids[c][1] += pt[1]
			counts[c]++
		}
		for j := 0; j < k; j++ {
			if counts[j] > 0 {
				newCentroids[j][0] /= float64(counts[j])
				newCentroids[j][1] /= float64(counts[j])
			} else {
				newCentroids[j] = pts[rand.Intn(len(pts))]
			}
		}
		centroids = newCentroids
	}

	result := make([]*geojson.Feature, len(fc.Features))
	for i, f := range fc.Features {
		props := make(map[string]any)
		for k, v := range f.Properties {
			props[k] = v
		}
		props["cluster"] = assignments[i]
		centroid := centroids[assignments[i]]
		props["centroid"] = []float64{centroid[0], centroid[1]}
		result[i] = geojson.NewFeature(f.Geometry, props)
	}

	return geojson.NewFeatureCollection(result), nil
}

type DbscanOptions struct {
	MinPoints int
}

func ClustersDbscan(fc *geojson.FeatureCollection, radius float64, options ...DbscanOptions) (*geojson.FeatureCollection, error) {
	opts := DbscanOptions{MinPoints: 3}
	if len(options) > 0 {
		opts = options[0]
	}
	if fc == nil || len(fc.Features) == 0 {
		return nil, fmt.Errorf("feature collection is required")
	}
	if radius <= 0 {
		return nil, fmt.Errorf("radius must be positive")
	}

	pts := make([]geojson.Position, len(fc.Features))
	for i, f := range fc.Features {
		coord, err := geojson.GetCoord(f)
		if err != nil {
			return nil, err
		}
		pts[i] = coord
	}

	n := len(pts)
	labels := make([]int, n)
	for i := range labels {
		labels[i] = -1
	}

	clusterID := 0

	for i := 0; i < n; i++ {
		if labels[i] != -1 {
			continue
		}
		neighbors := regionQuery(pts, i, radius)
		if len(neighbors) < opts.MinPoints {
			labels[i] = 0
			continue
		}
		clusterID++
		labels[i] = clusterID
		seeds := neighbors
		for _, seedIdx := range seeds {
			if labels[seedIdx] == 0 {
				labels[seedIdx] = clusterID
			} else if labels[seedIdx] != -1 {
				continue
			} else {
				labels[seedIdx] = clusterID
				seedNeighbors := regionQuery(pts, seedIdx, radius)
				if len(seedNeighbors) >= opts.MinPoints {
					for _, sn := range seedNeighbors {
						found := false
						for _, s := range seeds {
							if s == sn {
								found = true
								break
							}
						}
						if !found {
							seeds = append(seeds, sn)
						}
					}
				}
			}
		}
	}

	result := make([]*geojson.Feature, n)
	for i, f := range fc.Features {
		props := make(map[string]any)
		for k, v := range f.Properties {
			props[k] = v
		}
		props["cluster"] = labels[i]
		result[i] = geojson.NewFeature(f.Geometry, props)
	}

	return geojson.NewFeatureCollection(result), nil
}

func regionQuery(pts []geojson.Position, idx int, radius float64) []int {
	var neighbors []int
	for j := 0; j < len(pts); j++ {
		if j == idx {
			continue
		}
		if sqDist(pts[idx], pts[j]) <= radius*radius {
			neighbors = append(neighbors, j)
		}
	}
	return neighbors
}

func sqDist(a, b geojson.Position) float64 {
	dx := a[0] - b[0]
	dy := a[1] - b[1]
	return dx*dx + dy*dy
}

func Dissolve(fc *geojson.FeatureCollection, property string) (*geojson.FeatureCollection, error) {
	if fc == nil || len(fc.Features) == 0 {
		return nil, fmt.Errorf("feature collection is required")
	}

	groups := make(map[string][]*geojson.Feature)
	for _, f := range fc.Features {
		key := ""
		if property != "" {
			if v, ok := f.Properties[property]; ok {
				key = fmt.Sprintf("%v", v)
			}
		}
		groups[key] = append(groups[key], f)
	}

	var result []*geojson.Feature
	for key, group := range groups {
		if len(group) == 1 {
			result = append(result, group[0])
			continue
		}
		props := map[string]any{}
		if key != "" {
			if property != "" {
				props[property] = group[0].Properties[property]
			}
		}
		merged := unionGroup(group)
		result = append(result, merged...)
	}

	return geojson.NewFeatureCollection(result), nil
}

func unionGroup(features []*geojson.Feature) []*geojson.Feature {
	var polys []*geojson.Polygon
	for _, f := range features {
		geom, err := geojson.GetGeometry(f)
		if err != nil {
			continue
		}
		switch g := geom.(type) {
		case *geojson.Polygon:
			polys = append(polys, g)
		case *geojson.MultiPolygon:
			for _, p := range g.Coordinates {
				polys = append(polys, geojson.NewPolygon(p))
			}
		}
	}

	var result []*geojson.Feature
	used := make([]bool, len(polys))

	for i := 0; i < len(polys); i++ {
		if used[i] {
			continue
		}
		merged := polys[i]
		changed := true
		for changed {
			changed = false
			for j := i + 1; j < len(polys); j++ {
				if used[j] {
					continue
				}
				if ringsTouch(merged.Coordinates[0], polys[j].Coordinates[0]) {
					newRings := append(merged.Coordinates, polys[j].Coordinates...)
					merged = geojson.NewPolygon(newRings)
					used[j] = true
					changed = true
				}
			}
		}
		result = append(result, geojson.NewFeature(merged, nil))
	}

	return result
}

func ringsTouch(a, b []geojson.Position) bool {
	for _, pa := range a {
		for _, pb := range b {
			if math.Abs(pa[0]-pb[0]) < 1e-10 && math.Abs(pa[1]-pb[1]) < 1e-10 {
				return true
			}
		}
	}
	return false
}

func ClustersDbscanWithin(fc *geojson.FeatureCollection, radius float64, options ...DbscanOptions) (*geojson.FeatureCollection, error) {
	return ClustersDbscan(fc, radius, options...)
}
