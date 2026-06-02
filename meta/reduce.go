package meta

import (
	"github.com/ibinh/turf-go/geojson"
)

func CoordCount(obj any) (int, error) {
	count := 0
	err := CoordEach(obj, func(coord geojson.Position, index int) error {
		count++
		return nil
	})
	return count, err
}

func GeomCount(obj any) (int, error) {
	count := 0
	err := GeomEach(obj, func(geom geojson.Geometry, index int) error {
		count++
		return nil
	})
	return count, err
}

func FeatureCount(obj any) (int, error) {
	count := 0
	err := FeatureEach(obj, func(f *geojson.Feature, index int) error {
		count++
		return nil
	})
	return count, err
}

func GetFirstCoord(obj any) (geojson.Position, error) {
	var found geojson.Position
	err := CoordEach(obj, func(coord geojson.Position, index int) error {
		if index == 0 {
			found = coord
		}
		return nil
	})
	return found, err
}

func GetLastCoord(obj any) (geojson.Position, error) {
	var found geojson.Position
	err := CoordEach(obj, func(coord geojson.Position, index int) error {
		found = coord
		return nil
	})
	return found, err
}
