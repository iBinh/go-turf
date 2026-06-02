package grid

import (
	"fmt"
	"math"
	"github.com/ibinh/turf-go/geojson"
	"github.com/ibinh/turf-go/measurement"
)

type GridOptions struct {
	Properties map[string]any
}

func HexGrid(bbox []float64, cellSide float64, units measurement.Unit, options ...GridOptions) (*geojson.FeatureCollection, error) {
	if err := validBBox(bbox); err != nil {
		return nil, err
	}
	cellSide = measurement.ConvertLength(cellSide, units, measurement.UnitDegrees)
	return hexGrid(bbox, cellSide, options...)
}

func hexGrid(bbox []float64, cellSide float64, options ...GridOptions) (*geojson.FeatureCollection, error) {
	opts := getOpts(options)
	minX, minY, maxX, maxY := bbox[0], bbox[1], bbox[2], bbox[3]

	centerY := (minY + maxY) / 2
	cellSide = adjustCellSide(cellSide, centerY)
	if cellSide == 0 {
		cellSide = math.Max(maxX-minX, 0.001)
	}

	xFraction := cellSide * 2
	yFraction := cellSide * math.Sqrt(3)

	var features []*geojson.Feature
	rown := 0

	for y := minY; y < maxY; y += yFraction * 3 / 4 {
		startX := minX
		if rown%2 == 1 {
			startX = minX - xFraction/2
		}
		for x := startX; x < maxX; x += xFraction {
			features = append(features, hexCell(x, y, cellSide, opts.Properties))
		}
		rown++
	}
	if len(features) == 0 {
		features = append(features, hexCell(minX, minY, cellSide, opts.Properties))
	}

	return geojson.NewFeatureCollection(features), nil
}

func hexCell(cx, cy, side float64, props map[string]any) *geojson.Feature {
	h := side * math.Sqrt(3) / 2
	ring := []geojson.Position{
		{cx + side/2, cy + h},
		{cx + side, cy},
		{cx + side/2, cy - h},
		{cx - side/2, cy - h},
		{cx - side, cy},
		{cx - side/2, cy + h},
		{cx + side/2, cy + h},
	}
	return geojson.NewFeature(geojson.NewPolygon([][]geojson.Position{ring}), props)
}

func PointGrid(bbox []float64, cellSide float64, units measurement.Unit, options ...GridOptions) (*geojson.FeatureCollection, error) {
	if err := validBBox(bbox); err != nil {
		return nil, err
	}
	cellSide = measurement.ConvertLength(cellSide, units, measurement.UnitDegrees)
	return pointGrid(bbox, cellSide, options...)
}

func pointGrid(bbox []float64, cellSide float64, options ...GridOptions) (*geojson.FeatureCollection, error) {
	opts := getOpts(options)
	minX, minY, maxX, maxY := bbox[0], bbox[1], bbox[2], bbox[3]

	centerY := (minY + maxY) / 2
	cellSide = adjustCellSide(cellSide, centerY)
	if cellSide == 0 {
		cellSide = 0.001
	}

	var features []*geojson.Feature
	for x := minX; x <= maxX; x += cellSide {
		for y := minY; y <= maxY; y += cellSide {
			features = append(features, geojson.NewFeature(geojson.NewPoint(geojson.Position{x, y}), opts.Properties))
		}
	}

	return geojson.NewFeatureCollection(features), nil
}

func SquareGrid(bbox []float64, cellSide float64, units measurement.Unit, options ...GridOptions) (*geojson.FeatureCollection, error) {
	if err := validBBox(bbox); err != nil {
		return nil, err
	}
	cellSide = measurement.ConvertLength(cellSide, units, measurement.UnitDegrees)
	return squareGrid(bbox, cellSide, options...)
}

func squareGrid(bbox []float64, cellSide float64, options ...GridOptions) (*geojson.FeatureCollection, error) {
	opts := getOpts(options)
	minX, minY, maxX, maxY := bbox[0], bbox[1], bbox[2], bbox[3]

	centerY := (minY + maxY) / 2
	cellSide = adjustCellSide(cellSide, centerY)
	if cellSide == 0 {
		cellSide = 0.001
	}

	var features []*geojson.Feature
	for x := minX; x < maxX; x += cellSide {
		for y := minY; y < maxY; y += cellSide {
			ring := []geojson.Position{
				{x, y},
				{x + cellSide, y},
				{x + cellSide, y + cellSide},
				{x, y + cellSide},
				{x, y},
			}
			features = append(features, geojson.NewFeature(geojson.NewPolygon([][]geojson.Position{ring}), opts.Properties))
		}
	}

	return geojson.NewFeatureCollection(features), nil
}

func TriangleGrid(bbox []float64, cellSide float64, units measurement.Unit, options ...GridOptions) (*geojson.FeatureCollection, error) {
	if err := validBBox(bbox); err != nil {
		return nil, err
	}
	cellSide = measurement.ConvertLength(cellSide, units, measurement.UnitDegrees)
	return triangleGrid(bbox, cellSide, options...)
}

func triangleGrid(bbox []float64, cellSide float64, options ...GridOptions) (*geojson.FeatureCollection, error) {
	opts := getOpts(options)
	minX, minY, maxX, maxY := bbox[0], bbox[1], bbox[2], bbox[3]

	centerY := (minY + maxY) / 2
	cellSide = adjustCellSide(cellSide, centerY)
	if cellSide == 0 {
		cellSide = 0.001
	}

	var features []*geojson.Feature
	for x := minX; x < maxX; x += cellSide {
		for y := minY; y < maxY; y += cellSide {
			x2, y2 := x+cellSide, y+cellSide
			features = append(features,
				geojson.NewFeature(geojson.NewPolygon([][]geojson.Position{{{x, y}, {x2, y2}, {x, y2}, {x, y}}}), opts.Properties),
				geojson.NewFeature(geojson.NewPolygon([][]geojson.Position{{{x, y}, {x2, y}, {x2, y2}, {x, y}}}), opts.Properties),
			)
		}
	}

	return geojson.NewFeatureCollection(features), nil
}

func validBBox(bbox []float64) error {
	if len(bbox) < 4 {
		return fmt.Errorf("bbox must have at least 4 elements")
	}
	return nil
}

func getOpts(options []GridOptions) GridOptions {
	if len(options) > 0 {
		return options[0]
	}
	return GridOptions{Properties: map[string]any{}}
}

func adjustCellSide(cellSide, centerY float64) float64 {
	cellSide = cellSide / math.Cos(centerY*math.Pi/180)
	if math.IsInf(cellSide, 0) || math.IsNaN(cellSide) {
		cellSide = 0.001
	}
	return cellSide
}
