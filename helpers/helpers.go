package helpers

import "github.com/ibinh/turf-go/geojson"

func Point(coordinates geojson.Position, properties map[string]any, options ...geojson.FeatureOption) *geojson.Feature {
	f := geojson.NewFeature(geojson.NewPoint(coordinates), properties)
	for _, opt := range options {
		opt(f)
	}
	return f
}

func MultiPoint(coordinates []geojson.Position, properties map[string]any, options ...geojson.FeatureOption) *geojson.Feature {
	f := geojson.NewFeature(geojson.NewMultiPoint(coordinates), properties)
	for _, opt := range options {
		opt(f)
	}
	return f
}

func LineString(coordinates []geojson.Position, properties map[string]any, options ...geojson.FeatureOption) *geojson.Feature {
	f := geojson.NewFeature(geojson.NewLineString(coordinates), properties)
	for _, opt := range options {
		opt(f)
	}
	return f
}

func MultiLineString(coordinates [][]geojson.Position, properties map[string]any, options ...geojson.FeatureOption) *geojson.Feature {
	f := geojson.NewFeature(geojson.NewMultiLineString(coordinates), properties)
	for _, opt := range options {
		opt(f)
	}
	return f
}

func Polygon(coordinates [][]geojson.Position, properties map[string]any, options ...geojson.FeatureOption) *geojson.Feature {
	f := geojson.NewFeature(geojson.NewPolygon(coordinates), properties)
	for _, opt := range options {
		opt(f)
	}
	return f
}

func MultiPolygon(coordinates [][][]geojson.Position, properties map[string]any, options ...geojson.FeatureOption) *geojson.Feature {
	f := geojson.NewFeature(geojson.NewMultiPolygon(coordinates), properties)
	for _, opt := range options {
		opt(f)
	}
	return f
}

func GeometryCollection(geometries []geojson.Geometry, properties map[string]any, options ...geojson.FeatureOption) *geojson.Feature {
	f := geojson.NewFeature(geojson.NewGeometryCollection(geometries), properties)
	for _, opt := range options {
		opt(f)
	}
	return f
}
