package geojson

import (
	"encoding/json"
	"fmt"
)

type Position []float64

func (p Position) Lat() float64 {
	if len(p) < 2 {
		return 0
	}
	return p[1]
}

func (p Position) Lng() float64 {
	if len(p) < 2 {
		return 0
	}
	return p[0]
}

func (p Position) Elevation() float64 {
	if len(p) < 3 {
		return 0
	}
	return p[2]
}

const (
	TypePoint              = "Point"
	TypeMultiPoint         = "MultiPoint"
	TypeLineString         = "LineString"
	TypeMultiLineString    = "MultiLineString"
	TypePolygon            = "Polygon"
	TypeMultiPolygon       = "MultiPolygon"
	TypeGeometryCollection = "GeometryCollection"
	TypeFeature            = "Feature"
	TypeFeatureCollection  = "FeatureCollection"
)

type Geometry interface {
	Type() string
	BBox() []float64
	SetBBox([]float64)
	Dimensions() int
}

type Point struct {
	Coordinates Position `json:"coordinates"`
	bbox        []float64
}

func NewPoint(coordinates Position) *Point {
	return &Point{Coordinates: coordinates}
}

func (p *Point) Type() string          { return TypePoint }
func (p *Point) BBox() []float64       { return p.bbox }
func (p *Point) SetBBox(b []float64)   { p.bbox = b }
func (p *Point) Dimensions() int       { return 0 }

type pointJSON struct {
	Type        string    `json:"type"`
	Coordinates Position  `json:"coordinates"`
	BBox        []float64 `json:"bbox,omitempty"`
}

func (p *Point) MarshalJSON() ([]byte, error) {
	return json.Marshal(pointJSON{
		Type:        TypePoint,
		Coordinates: p.Coordinates,
		BBox:        p.bbox,
	})
}

func (p *Point) UnmarshalJSON(data []byte) error {
	var j pointJSON
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	p.Coordinates = j.Coordinates
	p.bbox = j.BBox
	return nil
}

type MultiPoint struct {
	Coordinates []Position `json:"coordinates"`
	bbox        []float64
}

func NewMultiPoint(coordinates []Position) *MultiPoint {
	return &MultiPoint{Coordinates: coordinates}
}

func (mp *MultiPoint) Type() string        { return TypeMultiPoint }
func (mp *MultiPoint) BBox() []float64     { return mp.bbox }
func (mp *MultiPoint) SetBBox(b []float64) { mp.bbox = b }
func (mp *MultiPoint) Dimensions() int     { return 0 }

type multiPointJSON struct {
	Type        string      `json:"type"`
	Coordinates []Position  `json:"coordinates"`
	BBox        []float64   `json:"bbox,omitempty"`
}

func (mp *MultiPoint) MarshalJSON() ([]byte, error) {
	return json.Marshal(multiPointJSON{
		Type:        TypeMultiPoint,
		Coordinates: mp.Coordinates,
		BBox:        mp.bbox,
	})
}

func (mp *MultiPoint) UnmarshalJSON(data []byte) error {
	var j multiPointJSON
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	mp.Coordinates = j.Coordinates
	mp.bbox = j.BBox
	return nil
}

type LineString struct {
	Coordinates []Position `json:"coordinates"`
	bbox        []float64
}

func NewLineString(coordinates []Position) *LineString {
	return &LineString{Coordinates: coordinates}
}

func (ls *LineString) Type() string        { return TypeLineString }
func (ls *LineString) BBox() []float64     { return ls.bbox }
func (ls *LineString) SetBBox(b []float64) { ls.bbox = b }
func (ls *LineString) Dimensions() int     { return 1 }

type lineStringJSON struct {
	Type        string      `json:"type"`
	Coordinates []Position  `json:"coordinates"`
	BBox        []float64   `json:"bbox,omitempty"`
}

func (ls *LineString) MarshalJSON() ([]byte, error) {
	return json.Marshal(lineStringJSON{
		Type:        TypeLineString,
		Coordinates: ls.Coordinates,
		BBox:        ls.bbox,
	})
}

func (ls *LineString) UnmarshalJSON(data []byte) error {
	var j lineStringJSON
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	ls.Coordinates = j.Coordinates
	ls.bbox = j.BBox
	return nil
}

type MultiLineString struct {
	Coordinates [][]Position `json:"coordinates"`
	bbox        []float64
}

func NewMultiLineString(coordinates [][]Position) *MultiLineString {
	return &MultiLineString{Coordinates: coordinates}
}

func (mls *MultiLineString) Type() string        { return TypeMultiLineString }
func (mls *MultiLineString) BBox() []float64     { return mls.bbox }
func (mls *MultiLineString) SetBBox(b []float64) { mls.bbox = b }
func (mls *MultiLineString) Dimensions() int     { return 1 }

type multiLineStringJSON struct {
	Type        string        `json:"type"`
	Coordinates [][]Position  `json:"coordinates"`
	BBox        []float64     `json:"bbox,omitempty"`
}

func (mls *MultiLineString) MarshalJSON() ([]byte, error) {
	return json.Marshal(multiLineStringJSON{
		Type:        TypeMultiLineString,
		Coordinates: mls.Coordinates,
		BBox:        mls.bbox,
	})
}

func (mls *MultiLineString) UnmarshalJSON(data []byte) error {
	var j multiLineStringJSON
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	mls.Coordinates = j.Coordinates
	mls.bbox = j.BBox
	return nil
}

type Polygon struct {
	Coordinates [][]Position `json:"coordinates"`
	bbox        []float64
}

func NewPolygon(coordinates [][]Position) *Polygon {
	return &Polygon{Coordinates: coordinates}
}

func (p *Polygon) Type() string        { return TypePolygon }
func (p *Polygon) BBox() []float64     { return p.bbox }
func (p *Polygon) SetBBox(b []float64) { p.bbox = b }
func (p *Polygon) Dimensions() int     { return 2 }

type polygonJSON struct {
	Type        string        `json:"type"`
	Coordinates [][]Position  `json:"coordinates"`
	BBox        []float64     `json:"bbox,omitempty"`
}

func (p *Polygon) MarshalJSON() ([]byte, error) {
	return json.Marshal(polygonJSON{
		Type:        TypePolygon,
		Coordinates: p.Coordinates,
		BBox:        p.bbox,
	})
}

func (p *Polygon) UnmarshalJSON(data []byte) error {
	var j polygonJSON
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	p.Coordinates = j.Coordinates
	p.bbox = j.BBox
	return nil
}

type MultiPolygon struct {
	Coordinates [][][]Position `json:"coordinates"`
	bbox        []float64
}

func NewMultiPolygon(coordinates [][][]Position) *MultiPolygon {
	return &MultiPolygon{Coordinates: coordinates}
}

func (mp *MultiPolygon) Type() string        { return TypeMultiPolygon }
func (mp *MultiPolygon) BBox() []float64     { return mp.bbox }
func (mp *MultiPolygon) SetBBox(b []float64) { mp.bbox = b }
func (mp *MultiPolygon) Dimensions() int     { return 2 }

type multiPolygonJSON struct {
	Type        string          `json:"type"`
	Coordinates [][][]Position  `json:"coordinates"`
	BBox        []float64       `json:"bbox,omitempty"`
}

func (mp *MultiPolygon) MarshalJSON() ([]byte, error) {
	return json.Marshal(multiPolygonJSON{
		Type:        TypeMultiPolygon,
		Coordinates: mp.Coordinates,
		BBox:        mp.bbox,
	})
}

func (mp *MultiPolygon) UnmarshalJSON(data []byte) error {
	var j multiPolygonJSON
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	mp.Coordinates = j.Coordinates
	mp.bbox = j.BBox
	return nil
}

type GeometryCollection struct {
	Geometries []Geometry `json:"-"`
	bbox       []float64
}

func NewGeometryCollection(geometries []Geometry) *GeometryCollection {
	return &GeometryCollection{Geometries: geometries}
}

func (gc *GeometryCollection) Type() string        { return TypeGeometryCollection }
func (gc *GeometryCollection) BBox() []float64     { return gc.bbox }
func (gc *GeometryCollection) SetBBox(b []float64) { gc.bbox = b }
func (gc *GeometryCollection) Dimensions() int     { return 0 }

type Feature struct {
	typeName   string
	Geometry   Geometry               `json:"-"`
	Properties map[string]any         `json:"properties,omitempty"`
	ID         any                    `json:"id,omitempty"`
	bbox       []float64
}

func NewFeature(geometry Geometry, properties map[string]any) *Feature {
	return &Feature{
		typeName:   TypeFeature,
		Geometry:   geometry,
		Properties: properties,
	}
}

func (f *Feature) Type() string           { return f.typeName }
func (f *Feature) SetType(t string)       { f.typeName = t }
func (f *Feature) BBox() []float64        { return f.bbox }
func (f *Feature) SetBBox(b []float64)    { f.bbox = b }

type FeatureCollection struct {
	typeName string
	Features []*Feature `json:"features"`
	bbox     []float64
}

func NewFeatureCollection(features []*Feature) *FeatureCollection {
	return &FeatureCollection{
		typeName: TypeFeatureCollection,
		Features: features,
	}
}

func (fc *FeatureCollection) Type() string           { return fc.typeName }
func (fc *FeatureCollection) SetType(t string)       { fc.typeName = t }
func (fc *FeatureCollection) BBox() []float64        { return fc.bbox }
func (fc *FeatureCollection) SetBBox(b []float64)    { fc.bbox = b }

func validatePosition(p Position) error {
	if len(p) < 2 {
		return fmt.Errorf("position must have at least 2 elements, got %d", len(p))
	}
	return nil
}

func validateLinearRing(ring []Position) error {
	if len(ring) < 4 {
		return fmt.Errorf("linear ring must have at least 4 positions, got %d", len(ring))
	}
	first, last := ring[0], ring[len(ring)-1]
	for i := range first {
		if first[i] != last[i] {
			return fmt.Errorf("linear ring is not closed: first %v, last %v", first, last)
		}
	}
	return nil
}

type GeoJSON interface {
	Type() string
	BBox() []float64
	SetBBox([]float64)
}

var (
	_ Geometry = (*Point)(nil)
	_ Geometry = (*MultiPoint)(nil)
	_ Geometry = (*LineString)(nil)
	_ Geometry = (*MultiLineString)(nil)
	_ Geometry = (*Polygon)(nil)
	_ Geometry = (*MultiPolygon)(nil)
	_ Geometry = (*GeometryCollection)(nil)
	_ GeoJSON  = (*Feature)(nil)
	_ GeoJSON  = (*FeatureCollection)(nil)
)
