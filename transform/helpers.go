package transform

import (
	"math"
	"github.com/ibinh/turf-go/geojson"
	"github.com/ibinh/turf-go/meta"
)

func applyToCoords(geom any, fn func(p geojson.Position) geojson.Position) (geojson.Geometry, error) {
	var result geojson.Geometry

	err := meta.GeomEach(geom, func(g geojson.Geometry, _ int) error {
		var err error
		result, err = transformGeometry(g, fn)
		return err
	})
	return result, err
}

func transformGeometry(g geojson.Geometry, fn func(geojson.Position) geojson.Position) (geojson.Geometry, error) {
	switch v := g.(type) {
	case *geojson.Point:
		return geojson.NewPoint(fn(v.Coordinates)), nil
	case *geojson.MultiPoint:
		pts := make([]geojson.Position, len(v.Coordinates))
		for i, c := range v.Coordinates {
			pts[i] = fn(c)
		}
		return geojson.NewMultiPoint(pts), nil
	case *geojson.LineString:
		pts := make([]geojson.Position, len(v.Coordinates))
		for i, c := range v.Coordinates {
			pts[i] = fn(c)
		}
		return geojson.NewLineString(pts), nil
	case *geojson.MultiLineString:
		lines := make([][]geojson.Position, len(v.Coordinates))
		for i, line := range v.Coordinates {
			lines[i] = make([]geojson.Position, len(line))
			for j, c := range line {
				lines[i][j] = fn(c)
			}
		}
		return geojson.NewMultiLineString(lines), nil
	case *geojson.Polygon:
		rings := make([][]geojson.Position, len(v.Coordinates))
		for i, ring := range v.Coordinates {
			rings[i] = make([]geojson.Position, len(ring))
			for j, c := range ring {
				rings[i][j] = fn(c)
			}
		}
		return geojson.NewPolygon(rings), nil
	case *geojson.MultiPolygon:
		polygons := make([][][]geojson.Position, len(v.Coordinates))
		for i, poly := range v.Coordinates {
			polygons[i] = make([][]geojson.Position, len(poly))
			for j, ring := range poly {
				polygons[i][j] = make([]geojson.Position, len(ring))
				for k, c := range ring {
					polygons[i][j][k] = fn(c)
				}
			}
		}
		return geojson.NewMultiPolygon(polygons), nil
	case *geojson.GeometryCollection:
		geoms := make([]geojson.Geometry, len(v.Geometries))
		for i, sub := range v.Geometries {
			var err error
			geoms[i], err = transformGeometry(sub, fn)
			if err != nil {
				return nil, err
			}
		}
		return geojson.NewGeometryCollection(geoms), nil
	}
	return nil, nil
}

func getPivot(obj any) (geojson.Position, error) {
	var sumLng, sumLat float64
	count := 0
	err := meta.CoordEach(obj, func(c geojson.Position, _ int) error {
		sumLng += c[0]
		sumLat += c[1]
		count++
		return nil
	})
	if err != nil || count == 0 {
		return geojson.Position{0, 0}, nil
	}
	return geojson.Position{sumLng / float64(count), sumLat / float64(count)}, nil
}

func degToRad(d float64) float64 { return d * math.Pi / 180 }
