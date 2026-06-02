package geojson

import "fmt"

func GetType(obj GeoJSON) string {
	return obj.Type()
}

func GetGeometry(obj any) (Geometry, error) {
	switch v := obj.(type) {
	case *Feature:
		if v.Geometry == nil {
			return nil, fmt.Errorf("feature has no geometry")
		}
		return v.Geometry, nil
	case Geometry:
		return v, nil
	default:
		return nil, fmt.Errorf("expected a Feature or Geometry, got %T", obj)
	}
}

func GetCoordinates(obj any) (any, error) {
	geom, err := GetGeometry(obj)
	if err != nil {
		return nil, err
	}
	return geom, nil
}

func GetCoord(obj any) (Position, error) {
	geom, err := GetGeometry(obj)
	if err != nil {
		return nil, err
	}
	p, ok := geom.(*Point)
	if !ok {
		return nil, fmt.Errorf("expected Point geometry, got %s", geom.Type())
	}
	return p.Coordinates, nil
}

func GetCoords(obj any) (any, error) {
	geom, err := GetGeometry(obj)
	if err != nil {
		return nil, err
	}
	switch g := geom.(type) {
	case *Point:
		return g.Coordinates, nil
	case *MultiPoint:
		return g.Coordinates, nil
	case *LineString:
		return g.Coordinates, nil
	case *MultiLineString:
		return g.Coordinates, nil
	case *Polygon:
		return g.Coordinates, nil
	case *MultiPolygon:
		return g.Coordinates, nil
	default:
		return nil, fmt.Errorf("GetCoords not supported for %T", geom)
	}
}

func GetBBox(obj GeoJSON) []float64 {
	return obj.BBox()
}

func CoordAll(obj any) ([]Position, error) {
	switch v := obj.(type) {
	case *FeatureCollection:
		var result []Position
		for _, f := range v.Features {
			if f.Geometry == nil {
				continue
			}
			result = append(result, collectCoords(f.Geometry)...)
		}
		return result, nil
	case *Feature:
		if v.Geometry == nil {
			return nil, fmt.Errorf("feature has no geometry")
		}
		return collectCoords(v.Geometry), nil
	case Geometry:
		return collectCoords(v), nil
	default:
		return nil, fmt.Errorf("CoordAll: expected Feature, FeatureCollection, or Geometry, got %T", obj)
	}
}

func collectCoords(geom Geometry) []Position {
	switch g := geom.(type) {
	case *Point:
		return []Position{g.Coordinates}
	case *MultiPoint:
		return g.Coordinates
	case *LineString:
		return g.Coordinates
	case *MultiLineString:
		var result []Position
		for _, line := range g.Coordinates {
			result = append(result, line...)
		}
		return result
	case *Polygon:
		var result []Position
		for _, ring := range g.Coordinates {
			result = append(result, ring...)
		}
		return result
	case *MultiPolygon:
		var result []Position
		for _, poly := range g.Coordinates {
			for _, ring := range poly {
				result = append(result, ring...)
			}
		}
		return result
	case *GeometryCollection:
		var result []Position
		for _, sub := range g.Geometries {
			result = append(result, collectCoords(sub)...)
		}
		return result
	default:
		return nil
	}
}

func GetGeom(obj any) (Geometry, error) {
	return GetGeometry(obj)
}

func CollectionOf(fc *FeatureCollection, geomType string) error {
	for i, f := range fc.Features {
		if f.Geometry == nil {
			return fmt.Errorf("feature %d has no geometry", i)
		}
		if f.Geometry.Type() != geomType {
			return fmt.Errorf("feature %d: expected geometry type %s, got %s", i, geomType, f.Geometry.Type())
		}
	}
	return nil
}
