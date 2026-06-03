package meta

import (
	"fmt"

	"github.com/ibinh/turf-go/geojson"
)

func CoordEach(obj any, fn func(coord geojson.Position, index int) error) error {
	switch v := obj.(type) {
	case *geojson.FeatureCollection:
		idx := 0
		for _, f := range v.Features {
			if f.Geometry == nil {
				continue
			}
			_, err := coordReduce(f.Geometry, 0, func(acc int, coord geojson.Position, index int) (int, error) {
				return acc, fn(coord, idx)
			})
			if err != nil {
				return err
			}
			idx++
		}
		return nil
	default:
		geom, err := geojson.GetGeometry(obj)
		if err != nil {
			return err
		}
		_, err = coordReduce(geom, 0, func(acc int, coord geojson.Position, index int) (int, error) {
			return acc, fn(coord, index)
		})
		return err
	}
}

func CoordReduce[T any](obj any, fn func(acc T, coord geojson.Position, index int) (T, error), initial T) (T, error) {
	switch v := obj.(type) {
	case *geojson.FeatureCollection:
		idx := 0
		acc := initial
		for _, f := range v.Features {
			if f.Geometry == nil {
				continue
			}
			result, err := coordReduce(f.Geometry, acc, func(a T, c geojson.Position, _ int) (T, error) {
				return fn(a, c, idx)
			})
			if err != nil {
				return acc, err
			}
			acc = result
			idx++
		}
		return acc, nil
	default:
		geom, err := geojson.GetGeometry(obj)
		if err != nil {
			return initial, err
		}
		return coordReduce(geom, initial, fn)
	}
}

func coordReduce[T any](geom geojson.Geometry, initial T, fn func(T, geojson.Position, int) (T, error)) (T, error) {
	acc := initial
	var idx int

	switch g := geom.(type) {
	case *geojson.Point:
		acc, _ = callFn(acc, g.Coordinates, 0, fn)
		return acc, nil
	case *geojson.MultiPoint:
		for _, c := range g.Coordinates {
			acc, idx = callFn(acc, c, idx, fn)
		}
		return acc, nil
	case *geojson.LineString:
		for _, c := range g.Coordinates {
			acc, idx = callFn(acc, c, idx, fn)
		}
		return acc, nil
	case *geojson.MultiLineString:
		for _, line := range g.Coordinates {
			for _, c := range line {
				acc, idx = callFn(acc, c, idx, fn)
			}
		}
		return acc, nil
	case *geojson.Polygon:
		for _, ring := range g.Coordinates {
			for _, c := range ring {
				acc, idx = callFn(acc, c, idx, fn)
			}
		}
		return acc, nil
	case *geojson.MultiPolygon:
		for _, poly := range g.Coordinates {
			for _, ring := range poly {
				for _, c := range ring {
					acc, idx = callFn(acc, c, idx, fn)
				}
			}
		}
		return acc, nil
	case *geojson.GeometryCollection:
		for _, sub := range g.Geometries {
			var err error
			acc, err = coordReduce(sub, acc, fn)
			if err != nil {
				return acc, err
			}
		}
		return acc, nil
	default:
		return acc, fmt.Errorf("coordReduce: unexpected geometry type %T", geom)
	}
}

func callFn[T any](acc T, coord geojson.Position, index int, fn func(T, geojson.Position, int) (T, error)) (T, int) {
	result, err := fn(acc, coord, index)
	if err != nil {
		return result, index
	}
	return result, index + 1
}

func FeatureEach(obj any, fn func(f *geojson.Feature, index int) error) error {
	_, err := FeatureReduce(obj, 0, func(acc int, f *geojson.Feature, index int) (int, error) {
		return acc, fn(f, index)
	})
	return err
}

func FeatureReduce[T any](obj any, initial T, fn func(acc T, f *geojson.Feature, index int) (T, error)) (T, error) {
	acc := initial
	idx := 0

	switch v := obj.(type) {
	case *geojson.Feature:
		result, err := fn(acc, v, idx)
		if err != nil {
			return acc, err
		}
		return result, nil
	case *geojson.FeatureCollection:
		for _, f := range v.Features {
			result, err := fn(acc, f, idx)
			if err != nil {
				return acc, err
			}
			acc = result
			idx++
		}
		return acc, nil
	default:
		return acc, fmt.Errorf("FeatureReduce: expected Feature or FeatureCollection, got %T", obj)
	}
}

func GeomEach(obj any, fn func(geom geojson.Geometry, index int) error) error {
	_, err := GeomReduce(obj, 0, func(acc int, geom geojson.Geometry, index int) (int, error) {
		return acc, fn(geom, index)
	})
	return err
}

func GeomReduce[T any](obj any, initial T, fn func(acc T, geom geojson.Geometry, index int) (T, error)) (T, error) {
	idx := 0

	switch v := obj.(type) {
	case geojson.Geometry:
		result, err := fn(initial, v, idx)
		if err != nil {
			return initial, err
		}
		return result, nil
	case *geojson.Feature:
		if v.Geometry != nil {
			result, err := fn(initial, v.Geometry, idx)
			if err != nil {
				return initial, err
			}
			return result, nil
		}
		return initial, nil
	case *geojson.FeatureCollection:
		acc := initial
		for _, f := range v.Features {
			if f.Geometry == nil {
				continue
			}
			result, err := fn(acc, f.Geometry, idx)
			if err != nil {
				return acc, err
			}
			acc = result
			idx++
		}
		return acc, nil
	default:
		return initial, fmt.Errorf("GeomReduce: unexpected type %T", obj)
	}
}

func PropEach(obj any, fn func(props map[string]any, index int) error) error {
	_, err := PropReduce(obj, 0, func(acc int, props map[string]any, index int) (int, error) {
		return acc, fn(props, index)
	})
	return err
}

func PropReduce[T any](obj any, initial T, fn func(acc T, props map[string]any, index int) (T, error)) (T, error) {
	idx := 0

	processFeature := func(f *geojson.Feature) (T, int, error) {
		p := f.Properties
		if p == nil {
			p = map[string]any{}
		}
		result, err := fn(initial, p, idx)
		return result, idx + 1, err
	}

	switch v := obj.(type) {
	case *geojson.Feature:
		result, _, err := processFeature(v)
		return result, err
	case *geojson.FeatureCollection:
		acc := initial
		for _, f := range v.Features {
			p := f.Properties
			if p == nil {
				p = map[string]any{}
			}
			result, err := fn(acc, p, idx)
			if err != nil {
				return acc, err
			}
			acc = result
			idx++
		}
		return acc, nil
	default:
		return initial, fmt.Errorf("PropReduce: expected Feature or FeatureCollection, got %T", obj)
	}
}

func FlattenEach(obj any, fn func(f *geojson.Feature, index int) error) error {
	_, err := FlattenReduce(obj, 0, func(acc int, f *geojson.Feature, index int) (int, error) {
		return acc, fn(f, index)
	})
	return err
}

func FlattenReduce[T any](obj any, initial T, fn func(acc T, f *geojson.Feature, index int) (T, error)) (T, error) {
	idx := 0

	var process func(geom geojson.Geometry, props map[string]any) (T, error)
	process = func(geom geojson.Geometry, props map[string]any) (T, error) {
		if geom == nil {
			return initial, nil
		}

		makeFeature := func(g geojson.Geometry) *geojson.Feature {
			f := geojson.NewFeature(g, props)
			if bbox := g.BBox(); len(bbox) > 0 {
				f.SetBBox(bbox)
			}
			return f
		}

		switch g := geom.(type) {
		case *geojson.Point:
			return fn(initial, makeFeature(g), idx)
		case *geojson.LineString:
			return fn(initial, makeFeature(g), idx)
		case *geojson.Polygon:
			return fn(initial, makeFeature(g), idx)
		case *geojson.MultiPoint:
			acc := initial
			for _, coord := range g.Coordinates {
				f := geojson.NewFeature(geojson.NewPoint(coord), props)
				result, err := fn(acc, f, idx)
				if err != nil {
					return acc, err
				}
				acc = result
				idx++
			}
			return acc, nil
		case *geojson.MultiLineString:
			acc := initial
			for _, line := range g.Coordinates {
				f := geojson.NewFeature(geojson.NewLineString(line), props)
				result, err := fn(acc, f, idx)
				if err != nil {
					return acc, err
				}
				acc = result
				idx++
			}
			return acc, nil
		case *geojson.MultiPolygon:
			acc := initial
			for _, poly := range g.Coordinates {
				f := geojson.NewFeature(geojson.NewPolygon(poly), props)
				result, err := fn(acc, f, idx)
				if err != nil {
					return acc, err
				}
				acc = result
				idx++
			}
			return acc, nil
		case *geojson.GeometryCollection:
			acc := initial
			for _, sub := range g.Geometries {
				result, err := process(sub, props)
				if err != nil {
					return acc, err
				}
				acc = result
			}
			return acc, nil
		default:
			return initial, fmt.Errorf("FlattenReduce: unexpected geometry type %T", geom)
		}
	}

	switch v := obj.(type) {
	case *geojson.Feature:
		props := v.Properties
		if props == nil {
			props = map[string]any{}
		}
		return process(v.Geometry, props)
	case *geojson.FeatureCollection:
		acc := initial
		for _, f := range v.Features {
			props := f.Properties
			if props == nil {
				props = map[string]any{}
			}
			result, err := process(f.Geometry, props)
			if err != nil {
				return acc, err
			}
			acc = result
		}
		return acc, nil
	case geojson.Geometry:
		return process(v, map[string]any{})
	default:
		return initial, fmt.Errorf("FlattenReduce: expected Feature, FeatureCollection, or Geometry, got %T", obj)
	}
}
