package geojson

import (
	"encoding/json"
	"fmt"
)

type geoJSONObject struct {
	Type        string          `json:"type"`
	BBox        []float64       `json:"bbox,omitempty"`
	Coordinates json.RawMessage `json:"coordinates,omitempty"`
	Geometries  json.RawMessage `json:"geometries,omitempty"`
	Geometry    json.RawMessage `json:"geometry,omitempty"`
	Properties  json.RawMessage `json:"properties,omitempty"`
	ID          json.RawMessage `json:"id,omitempty"`
	Features    json.RawMessage `json:"features,omitempty"`
}

func dispatchGeometry(raw json.RawMessage) (Geometry, error) {
	var obj geoJSONObject
	if err := json.Unmarshal(raw, &obj); err != nil {
		return nil, err
	}
	switch obj.Type {
	case TypePoint:
		var p Point
		if err := json.Unmarshal(raw, &p); err != nil {
			return nil, err
		}
		return &p, nil
	case TypeMultiPoint:
		var mp MultiPoint
		if err := json.Unmarshal(raw, &mp); err != nil {
			return nil, err
		}
		return &mp, nil
	case TypeLineString:
		var ls LineString
		if err := json.Unmarshal(raw, &ls); err != nil {
			return nil, err
		}
		return &ls, nil
	case TypeMultiLineString:
		var mls MultiLineString
		if err := json.Unmarshal(raw, &mls); err != nil {
			return nil, err
		}
		return &mls, nil
	case TypePolygon:
		var poly Polygon
		if err := json.Unmarshal(raw, &poly); err != nil {
			return nil, err
		}
		return &poly, nil
	case TypeMultiPolygon:
		var mp MultiPolygon
		if err := json.Unmarshal(raw, &mp); err != nil {
			return nil, err
		}
		return &mp, nil
	case TypeGeometryCollection:
		var gc GeometryCollection
		if err := json.Unmarshal(raw, &gc); err != nil {
			return nil, err
		}
		return &gc, nil
	default:
		return nil, fmt.Errorf("unknown geometry type: %s", obj.Type)
	}
}

func marshalGeometry(geom Geometry) (json.RawMessage, error) {
	return json.Marshal(geom)
}

func (f *Feature) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type":       TypeFeature,
		"properties": f.Properties,
	}
	if f.Geometry != nil {
		m["geometry"] = f.Geometry
	} else {
		m["geometry"] = nil
	}
	if f.ID != nil {
		m["id"] = f.ID
	}
	if len(f.bbox) > 0 {
		m["bbox"] = f.bbox
	}
	return json.Marshal(m)
}

func (f *Feature) UnmarshalJSON(data []byte) error {
	var obj geoJSONObject
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	if obj.Type != TypeFeature {
		return fmt.Errorf("expected type Feature, got %s", obj.Type)
	}
	f.SetType(TypeFeature)
	f.bbox = obj.BBox

	if obj.Properties != nil {
		if err := json.Unmarshal(obj.Properties, &f.Properties); err != nil {
			return err
		}
	}
	if obj.ID != nil {
		if err := json.Unmarshal(obj.ID, &f.ID); err != nil {
			return err
		}
	}
	if obj.Geometry != nil {
		if string(obj.Geometry) == "null" {
			f.Geometry = nil
		} else {
			geom, err := dispatchGeometry(obj.Geometry)
			if err != nil {
				return fmt.Errorf("feature geometry: %w", err)
			}
			f.Geometry = geom
		}
	}
	return nil
}

func (gc *GeometryCollection) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type": TypeGeometryCollection,
	}
	geoms := make([]any, len(gc.Geometries))
	for i, g := range gc.Geometries {
		if g == nil {
			geoms[i] = nil
			continue
		}
		raw, err := marshalGeometry(g)
		if err != nil {
			return nil, err
		}
		var v any
		if err := json.Unmarshal(raw, &v); err != nil {
			return nil, err
		}
		geoms[i] = v
	}
	m["geometries"] = geoms
	if len(gc.bbox) > 0 {
		m["bbox"] = gc.bbox
	}
	return json.Marshal(m)
}

func (gc *GeometryCollection) UnmarshalJSON(data []byte) error {
	var obj geoJSONObject
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	if obj.Type != TypeGeometryCollection {
		return fmt.Errorf("expected type GeometryCollection, got %s", obj.Type)
	}
	gc.bbox = obj.BBox

	if obj.Geometries != nil {
		var rawGeoms []json.RawMessage
		if err := json.Unmarshal(obj.Geometries, &rawGeoms); err != nil {
			return err
		}
		gc.Geometries = make([]Geometry, len(rawGeoms))
		for i, raw := range rawGeoms {
			geom, err := dispatchGeometry(raw)
			if err != nil {
				return fmt.Errorf("geometry collection[%d]: %w", i, err)
			}
			gc.Geometries[i] = geom
		}
	}
	return nil
}

func (fc *FeatureCollection) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type":     TypeFeatureCollection,
		"features": fc.Features,
	}
	if len(fc.bbox) > 0 {
		m["bbox"] = fc.bbox
	}
	return json.Marshal(m)
}

func (fc *FeatureCollection) UnmarshalJSON(data []byte) error {
	var obj geoJSONObject
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	if obj.Type != TypeFeatureCollection {
		return fmt.Errorf("expected type FeatureCollection, got %s", obj.Type)
	}
	fc.SetType(TypeFeatureCollection)
	fc.bbox = obj.BBox

	if obj.Features != nil {
		var rawFeatures []json.RawMessage
		if err := json.Unmarshal(obj.Features, &rawFeatures); err != nil {
			return err
		}
		fc.Features = make([]*Feature, len(rawFeatures))
		for i, raw := range rawFeatures {
			var f Feature
			if err := json.Unmarshal(raw, &f); err != nil {
				return fmt.Errorf("feature collection[%d]: %w", i, err)
			}
			fc.Features[i] = &f
		}
	}
	return nil
}

func UnmarshalGeoJSON(data []byte) (GeoJSON, error) {
	var obj geoJSONObject
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}
	switch obj.Type {
	case TypeFeature:
		var f Feature
		if err := json.Unmarshal(data, &f); err != nil {
			return nil, err
		}
		return &f, nil
	case TypeFeatureCollection:
		var fc FeatureCollection
		if err := json.Unmarshal(data, &fc); err != nil {
			return nil, err
		}
		return &fc, nil
	case TypePoint, TypeMultiPoint, TypeLineString, TypeMultiLineString,
		TypePolygon, TypeMultiPolygon, TypeGeometryCollection:
		return dispatchGeometry(data)
	default:
		return nil, fmt.Errorf("unsupported GeoJSON type: %s", obj.Type)
	}
}
