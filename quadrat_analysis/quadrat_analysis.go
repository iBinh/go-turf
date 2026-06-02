package quadratanalysis

import (
	"fmt"
	"math"

	"github.com/ibinh/turf-go/geojson"
	"github.com/ibinh/turf-go/meta"
)

type QuadratAnalysisResult struct {
	Chi2 float64
	DF   int
	P    float64
}

func QuadratAnalysis(points *geojson.FeatureCollection, gridSize int) (*QuadratAnalysisResult, error) {
	if points == nil || len(points.Features) < 2 {
		return nil, fmt.Errorf("at least 2 points required")
	}
	if gridSize < 2 {
		gridSize = 10
	}

	var coords []geojson.Position
	meta.CoordEach(points, func(c geojson.Position, _ int) error {
		coords = append(coords, c)
		return nil
	})

	if len(coords) < 2 {
		return nil, fmt.Errorf("no valid coordinates")
	}

	minX, minY := coords[0][0], coords[0][1]
	maxX, maxY := minX, minY
	for _, c := range coords {
		if c[0] < minX {
			minX = c[0]
		}
		if c[0] > maxX {
			maxX = c[0]
		}
		if c[1] < minY {
			minY = c[1]
		}
		if c[1] > maxY {
			maxY = c[1]
		}
	}

	cellW := (maxX - minX) / float64(gridSize)
	cellH := (maxY - minY) / float64(gridSize)
	if cellW <= 0 || cellH <= 0 {
		return nil, fmt.Errorf("points are collinear")
	}

	grid := gridSize * gridSize
	counts := make([]int, grid)

	for _, c := range coords {
		col := int((c[0] - minX) / cellW)
		row := int((c[1] - minY) / cellH)
		if col >= gridSize {
			col = gridSize - 1
		}
		if row >= gridSize {
			row = gridSize - 1
		}
		if col < 0 {
			col = 0
		}
		if row < 0 {
			row = 0
		}
		counts[row*gridSize+col]++
	}

	n := float64(len(coords))
	mean := n / float64(grid)

	var chi2 float64
	for _, c := range counts {
		d := float64(c) - mean
		chi2 += d * d / mean
	}

	df := grid - 1
	p := 1.0
	if chi2 > 0 && df > 0 {
		p = chiSquaredPValue(chi2, df)
	}

	return &QuadratAnalysisResult{
		Chi2: chi2,
		DF:   df,
		P:    p,
	}, nil
}

func chiSquaredPValue(chi2 float64, df int) float64 {
	return 1.0 - chiSquaredCDF(chi2, df)
}

func chiSquaredCDF(x float64, k int) float64 {
	if x <= 0 {
		return 0
	}
	return 1.0 - regularizedGammaP(float64(k)/2, x/2)
}

func regularizedGammaP(a, x float64) float64 {
	const epsilon = 1e-10
	if x < a+1 {
		var sum, term float64 = 1, 1
		for i := 1; i <= 100; i++ {
			term *= x / (a + float64(i))
			sum += term
			if term < epsilon {
				break
			}
		}
		return math.Exp(-x + a*math.Log(x) - logGamma(a)) * sum
	}
	var f float64 = 1 + (1-a)/x
	if math.Abs(f) < 1e-15 {
		f = 1e-15
	}
	var c float64 = 1 / f
	var d float64 = 1 / f
	for i := 1; i <= 100; i++ {
		f = 1 + float64(i)*(1-a)/x
		if math.Abs(f) < 1e-15 {
			f = 1e-15
		}
		c = f
		d = 1 / f
		f = 1 + float64(i)/x
		if math.Abs(f) < 1e-15 {
			f = 1e-15
		}
		c *= f
		d *= f
		if math.Abs(d) < 1e-15 {
			d = 1e-15
		}
		d = 1 / d
		delta := c * d
		if math.Abs(delta-1) < epsilon {
			break
		}
	}
	return 1.0 - math.Exp(-x+a*math.Log(x)-logGamma(a))/math.Pi*f
}

func logGamma(x float64) float64 {
	if x <= 0 {
		return 0
	}
	coeff := []float64{
		76.18009172947146, -86.50532032941677,
		24.01409824083091, -1.231739572450155,
		0.1208650973866179e-2, -0.5395239384953e-5,
	}
	y := x
	tmp := x + 5.5
	tmp -= (x + 0.5) * math.Log(tmp)
	ser := 1.000000000190015
	for j := 0; j < 6; j++ {
		y++
		ser += coeff[j] / y
	}
	return -tmp + math.Log(2.5066282746310005*ser/x)
}
