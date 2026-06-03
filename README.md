# turf-go

<div align="center">

**Go port of [Turf.js](https://turfjs.org/) — advanced geospatial analysis for Go.**

[![Go Version](https://img.shields.io/github/go-mod/go-version/ibinh/go-turf)](https://github.com/ibinh/go-turf)
[![CI](https://github.com/ibinh/go-turf/actions/workflows/ci.yml/badge.svg)](https://github.com/ibinh/go-turf/actions/workflows/ci.yml)
[![Coverage](https://codecov.io/gh/ibinh/go-turf/branch/master/graph/badge.svg)](https://codecov.io/gh/ibinh/go-turf)
[![Go Report Card](https://goreportcard.com/badge/github.com/ibinh/go-turf)](https://goreportcard.com/report/github.com/ibinh/go-turf)
[![License](https://img.shields.io/github/license/ibinh/go-turf)](LICENSE)

Zero external dependencies. Single Go module. 39 sub-packages covering the full Turf.js v7.3.5 API surface.

</div>

---

## Installation

```bash
go get github.com/ibinh/go-turf
```

Requires Go 1.21+.

## Quick Start

```go
import turf "github.com/ibinh/go-turf"

// Distance between two points
from := turf.Point([]float64{-75.343, 39.984})
to   := turf.Point([]float64{-75.534, 39.123})
km, _ := turf.Distance(from, to, turf.UnitKilometers) // ~97.3 km

// Buffer a point (100 meters)
buffered, _ := turf.Buffer(from, 100, turf.UnitMeters)

// Point-in-polygon check
poly := turf.Polygon([][][]float64{{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}}})
inside, _ := turf.BooleanPointInPolygon(turf.Point([]float64{5, 5}), poly)

// Clip a polygon to a bounding box
bbox := []float64{0, 0, 5, 5}
clipped, _ := turf.BBoxClip(poly, bbox)

// Contour lines from point data
fc := turf.RandomPoint(100)
isolines, _ := turf.Isolines(fc, turf.IsolinesOptions{
    ZProperty: "z",
    Breaks:    []float64{5, 10, 15, 20},
})
```

## Packages

### Geometry & Types
| Package | Functions |
|---------|-----------|
| `geojson` | RFC 7946 compliant types — `Point`, `LineString`, `Polygon`, `Multi*`, `Feature`, `FeatureCollection`, `GeometryCollection` with constructors, invariants, JSON marshal/unmarshal |
| `helpers` | `Point()`, `LineString()`, `Polygon()`, `MultiPoint()`, `MultiLineString()`, `MultiPolygon()`, `GeometryCollection()` |

### Measurement
| Package | Functions |
|---------|-----------|
| `measurement` | `Distance` (Haversine/Rhumb), `Bearing`, `Destination`, `Length`, `Area` (geodesic), `Midpoint`, `Along`, `GreatCircle`, `PointToLineDistance`, `NearestPointOnLine`, `Angle`, `ConvertLength` |

### Spatial Analysis
| Package | Functions |
|---------|-----------|
| `boolean` | `Contains`, `Within`, `Intersects`, `Disjoint`, `Touches`, `Crosses`, `Overlap`, `PointInPolygon`, `PointOnLine`, `SegmentIntersect`, `Clockwise`, `Valid`, `Concave`, `BooleanEqual`, `BooleanParallel` |
| `polyclip` | `PolygonUnion`, `PolygonIntersect`, `PolygonDifference`, `PolygonXor` |
| `buffer` | `Buffer` — distance-based polygon buffering |
| `mask` | `Mask` — outer ring with hole masking |
| `unkink` | `UnkinkPolygon` — fix self-intersecting polygons |
| `kinks` | `Kinks` — self-intersection detection |
| `voronoi` | `Voronoi` — Voronoi diagram generation |
| `centroid` | `Center`, `Centroid`, `CenterMean`, `CenterMedian`, `CenterOfMass`, `PointOnFeature` |

### Transformations
| Package | Functions |
|---------|-----------|
| `transform` | `TransformRotate`, `TransformScale`/`ScaleXY`, `TransformTranslate`, `Flip`, `Rewind`, `Truncate`, `CleanCoords`, `ToMercator`, `ToWGS84` |
| `simplify` | `Simplify` (Ramer-Douglas-Peucker), `ConvexHull` (Andrew's monotone chain) |
| `smooth` | `PolygonSmooth` (Chaikin's subdivision) |
| `convert` | `PolygonToLine`, `LineToPolygon` |

### Lines & Shapes
| Package | Functions |
|---------|-----------|
| `lines` | `LineIntersect`, `LineSegment`, `LineOverlap`, `LineSlice`, `LineSliceAlong`, `LineChunk`, `LineSplit`, `LineArc`, `Sector` |
| `lineoffset` | `LineOffset` — parallel line with miter joins |
| `tangents` | `PolygonTangents` — rotating calipers |
| `shapes` | `Circle`, `Ellipse`, `BezierSpline`, `RandomPoint`, `RandomLineString`, `RandomPolygon`, `RandomPosition` |

### Grids & Statistics
| Package | Functions |
|---------|-----------|
| `grid` | `HexGrid`, `PointGrid`, `SquareGrid`, `TriangleGrid`, `RectangleGrid` |
| `clusters` | `ClustersKMeans`, `ClustersDbscan`, `Dissolve` |
| `interpolation` | `Sample`, `Tin` (Delaunay), `Interpolate` (IDW), `PlanarDistance`, `PlanarPointOnLine` |
| `nearest_neighbor` | `NearestNeighborAnalysis` — R statistic, z-score, p-value |
| `quadrat_analysis` | `QuadratAnalysis` — chi-squared quadrat test |
| `moran_index` | `MoranIndex` — spatial autocorrelation (Moran's I) |
| `directional_mean` | `DirectionalMean` — circular mean of line bearings |
| `standard_deviational_ellipse` | `StandardDeviationalEllipse` — 2-sigma directional distribution |
| `distance_weight` | `DistanceWeight` — inverse-distance weight matrix |

### Iterators & Data
| Package | Functions |
|---------|-----------|
| `meta` | `CoordEach`/`Reduce`, `FeatureEach`/`Reduce`, `GeomEach`/`Reduce`, `PropEach`/`Reduce`, `FlattenEach`/`Reduce`, `CoordCount`, `FeatureCount`, `GeomCount` |
| `data` | `Tag`, `Collect` |
| `misc` | `Clone`, `Combine`, `Explode`, `Flatten`, `PointsWithinPolygon`, `Planepoint`, `Tesselate` |

### Contours & Surface
| Package | Functions |
|---------|-----------|
| `isobands` | `Isobands` — marching squares contour band polygons |
| `isolines` | `Isolines` — marching squares contour polylines |
| `concave` | `ConcaveHull` — alpha-shapes concave hull via Delaunay triangulation |
| `polygonize` | `Polygonize` — polygon extraction from linework via graph walk |

### Spatial Index & Routing
| Package | Functions |
|---------|-----------|
| `rbush` | `RBush` — flat-array spatial index (`Insert`, `Search`, `Remove`, `Load`, `Collisions`, `Nearest`) |
| `shortest_path` | `ShortestPath` — Dijkstra shortest path with nearest-node snapping |

### Bounding Box
| Package | Functions |
|---------|-----------|
| `bbox` | `BBox`, `BBoxPolygon`, `Envelope`, `Square`, `BBoxClip` (Sutherland-Hodgman) |

## Development

```bash
# Run all tests
make test

# Run with race detector and coverage
make cover

# Benchmarks
make bench

# Lint
make lint
```

## License

MIT — see [LICENSE](LICENSE)
