package directionalmean

import (
	"math"
	"testing"

	"github.com/ibinh/turf-go/geojson"
)

func TestDirectionalMeanSingleLine(t *testing.T) {
	line := geojson.NewFeature(
		geojson.NewLineString([]geojson.Position{{0, 0}, {10, 10}}),
		nil,
	)

	mean, err := DirectionalMean(line)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(mean-45) > 0.01 {
		t.Errorf("expected ~45°, got %f", mean)
	}
}

func TestDirectionalMeanEastWestLine(t *testing.T) {
	line := geojson.NewFeature(
		geojson.NewLineString([]geojson.Position{{0, 0}, {10, 0}}),
		nil,
	)

	mean, err := DirectionalMean(line)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(mean-90) > 0.01 {
		t.Errorf("expected ~90°, got %f", mean)
	}
}

func TestDirectionalMeanMultiLineString(t *testing.T) {
	ml := geojson.NewFeature(
		geojson.NewMultiLineString([][]geojson.Position{
			{{0, 0}, {10, 10}},
			{{0, 0}, {10, 0}},
		}),
		nil,
	)

	mean, err := DirectionalMean(ml)
	if err != nil {
		t.Fatal(err)
	}
	if mean < 0 || mean > 360 {
		t.Errorf("expected bearing in [0, 360], got %f", mean)
	}
}

func TestDirectionalMeanNoLines(t *testing.T) {
	pt := geojson.NewFeature(geojson.NewPoint(geojson.Position{0, 0}), nil)

	_, err := DirectionalMean(pt)
	if err == nil {
		t.Fatal("expected error for no line features")
	}
}
