# turf-go

**Go port of [Turf.js](https://turfjs.org/)** — a modular geospatial analysis engine.

Zero external dependencies. Single module, 29 sub-packages covering the full Turf.js API surface.

## Packages

| Package | Description |
|---------|-------------|
| `geojson` | GeoJSON types (Point, LineString, Polygon, Multi\*, Feature, FeatureCollection, GeometryCollection) with custom JSON marshal/unmarshal, constructors, invariants |
| `helpers` | Convenience constructors: `Point()`, `LineString()`, `Polygon()`, `MultiPoint()`, etc. |
| `meta` | Iterators: `CoordEach/Reduce`, `FeatureEach/Reduce`, `GeomEach/Reduce`, `PropEach/Reduce`, `FlattenEach/Reduce` |
| `measurement` | Haversine/rhumb distance, bearing, destination, length, area, midpoint, along, great-circle, point-to-line-distance, nearest-point-on-line |
| `bbox` | `BBox`, `BBoxPolygon`, `Envelope`, `Square` |
| `center` | `Center`, `Centroid`, `CenterMean`, `CenterMedian`, `CenterOfMass` |
| `transform` | `Rotate`, `Scale`, `ScaleXY`, `Translate`, `Rewind`, `Flip`, `Truncate`, `CleanCoords`, `ToMercator`, `ToWGS84` |
| `boolean` | `Clockwise`, `PointInPolygon`, `PointOnLine`, `SegmentIntersect`, `Contains`, `Within`, `Intersects`, `Disjoint`, `Touches`, `Crosses`, `Overlap`, `Valid`, `Concave`, `BooleanEqual`, `BooleanParallel` |
| `grid` | `HexGrid`, `PointGrid`, `SquareGrid`, `TriangleGrid` |
| `shapes` | `Circle`, `Ellipse`, `BezierSpline`, `RandomPosition`, `RandomPoint`, `RandomLineString`, `RandomPolygon` |
| `interpolation` | `Sample`, `Tin` (Delaunay), `Interpolate` (IDW), `PlanarDistance`, `PlanarPointOnLine` |
| `clusters` | `ClustersKMeans`, `ClustersDbscan`, `Dissolve` |
| `data` | `Tag`, `Collect` |
| `simplify` | `Simplify` (RDP), `ConvexHull` |
| `polyclip` | `PolygonUnion`, `PolygonIntersect`, `PolygonDifference`, `PolygonXor` |
| `lines` | `LineIntersect`, `LineSegment`, `LineOverlap`, `LineSlice`, `LineSliceAlong`, `LineChunk`, `LineSplit`, `LineArc`, `Sector` |
| `convert` | `PolygonToLine`, `LineToPolygon` |
| `smooth` | `PolygonSmooth` (Chaikin) |
| `tangents` | `PolygonTangents` |
| `lineoffset` | `LineOffset` |
| `mask` | `Mask` |
| `unkink` | `UnkinkPolygon` |
| `voronoi` | `Voronoi` |
| `buffer` | `Buffer` |
| `kinks` | `Kinks` (self-intersection detection) |
| `misc` | `Clone`, `Combine`, `Explode`, `PointsWithinPolygon`, `Planepoint`, `Tesselate`, `Flatten` |
| `isobands` | `Isobands` — marching squares contour bands |
| `isolines` | `Isolines` — marching squares contour lines |
| `turf` (root) | Umbrella re-export of all packages |

## Usage

```go
import "github.com/ibinh/turf-go"

// Distance between two points
from := turf.Point([]float64{-75.343, 39.984})
to := turf.Point([]float64{-75.534, 39.123})
d, _ := turf.Distance(from, to, turf.UnitKilometers)

// Buffer a point
buffered, _ := turf.Buffer(from, 100, turf.UnitMeters)

// Polygon boolean ops
a := turf.Polygon([][][]float64{{{0, 0}, {2, 0}, {2, 2}, {0, 2}, {0, 0}}})
b := turf.Polygon([][][]float64{{{1, 1}, {3, 1}, {3, 3}, {1, 3}, {1, 1}}})
union, _ := turf.PolygonUnion(a, b)

// Isobands (contour polygons)
fc := turf.RandomPoint(100)
bands, _ := turf.Isobands(fc, turf.IsobandsOptions{
    ZProperty: "z",
    Breaks:    []float64{0, 10, 20, 30},
})

// Isolines (contour lines)
lines, _ := turf.Isolines(fc, turf.IsolinesOptions{
    ZProperty: "z",
    Breaks:    []float64{5, 10, 15},
})
```

## Development

```bash
make test    # run all tests
make bench   # benchmarks
make lint    # golangci-lint
make cover   # coverage report
```

Go 1.23+. Run with `GOROOT=/opt/homebrew/opt/go@1.23/libexec go test ./...` if using Homebrew's Go 1.23.

## License

MIT
