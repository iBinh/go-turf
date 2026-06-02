package turf

import (
	"github.com/ibinh/turf-go/bbox"
	"github.com/ibinh/turf-go/boolean"
	"github.com/ibinh/turf-go/buffer"
	"github.com/ibinh/turf-go/center"
	"github.com/ibinh/turf-go/clusters"
	"github.com/ibinh/turf-go/concave"
	"github.com/ibinh/turf-go/convert"
	"github.com/ibinh/turf-go/data"
	"github.com/ibinh/turf-go/directional_mean"
	"github.com/ibinh/turf-go/distance_weight"
	"github.com/ibinh/turf-go/geojson"
	"github.com/ibinh/turf-go/grid"
	"github.com/ibinh/turf-go/helpers"
	"github.com/ibinh/turf-go/interpolation"
	"github.com/ibinh/turf-go/isobands"
	"github.com/ibinh/turf-go/isolines"
	"github.com/ibinh/turf-go/kinks"
	"github.com/ibinh/turf-go/lineoffset"
	"github.com/ibinh/turf-go/lines"
	"github.com/ibinh/turf-go/mask"
	"github.com/ibinh/turf-go/measurement"
	"github.com/ibinh/turf-go/meta"
	"github.com/ibinh/turf-go/misc"
	"github.com/ibinh/turf-go/moran_index"
	"github.com/ibinh/turf-go/nearest_neighbor"
	"github.com/ibinh/turf-go/polyclip"
	"github.com/ibinh/turf-go/polygonize"
	"github.com/ibinh/turf-go/quadrat_analysis"
	"github.com/ibinh/turf-go/rbush"
	"github.com/ibinh/turf-go/shapes"
	"github.com/ibinh/turf-go/shortest_path"
	"github.com/ibinh/turf-go/simplify"
	"github.com/ibinh/turf-go/smooth"
	"github.com/ibinh/turf-go/standard_deviational_ellipse"
	"github.com/ibinh/turf-go/tangents"
	"github.com/ibinh/turf-go/transform"
	"github.com/ibinh/turf-go/unkink"
	"github.com/ibinh/turf-go/voronoi"
)

type Unit = measurement.Unit

const (
	UnitMeters       = measurement.UnitMeters
	UnitKilometers   = measurement.UnitKilometers
	UnitMiles        = measurement.UnitMiles
	UnitNauticalMiles = measurement.UnitNauticalMiles
	UnitDegrees      = measurement.UnitDegrees
	UnitRadians      = measurement.UnitRadians
	UnitFeet         = measurement.UnitFeet
)

const (
	TypePoint              = geojson.TypePoint
	TypeMultiPoint         = geojson.TypeMultiPoint
	TypeLineString         = geojson.TypeLineString
	TypeMultiLineString    = geojson.TypeMultiLineString
	TypePolygon            = geojson.TypePolygon
	TypeMultiPolygon       = geojson.TypeMultiPolygon
	TypeGeometryCollection = geojson.TypeGeometryCollection
	TypeFeature            = geojson.TypeFeature
	TypeFeatureCollection  = geojson.TypeFeatureCollection
)

func Point(coordinates geojson.Position, properties map[string]any, options ...geojson.FeatureOption) *geojson.Feature {
	return helpers.Point(coordinates, properties, options...)
}

func LineString(coordinates []geojson.Position, properties map[string]any, options ...geojson.FeatureOption) *geojson.Feature {
	return helpers.LineString(coordinates, properties, options...)
}

func Polygon(coordinates [][]geojson.Position, properties map[string]any, options ...geojson.FeatureOption) *geojson.Feature {
	return helpers.Polygon(coordinates, properties, options...)
}

func MultiPoint(coordinates []geojson.Position, properties map[string]any, options ...geojson.FeatureOption) *geojson.Feature {
	return helpers.MultiPoint(coordinates, properties, options...)
}

func MultiLineString(coordinates [][]geojson.Position, properties map[string]any, options ...geojson.FeatureOption) *geojson.Feature {
	return helpers.MultiLineString(coordinates, properties, options...)
}

func MultiPolygon(coordinates [][][]geojson.Position, properties map[string]any, options ...geojson.FeatureOption) *geojson.Feature {
	return helpers.MultiPolygon(coordinates, properties, options...)
}

func GeometryCollection(geometries []geojson.Geometry, properties map[string]any, options ...geojson.FeatureOption) *geojson.Feature {
	return helpers.GeometryCollection(geometries, properties, options...)
}

func WithBBox(bbox []float64) geojson.FeatureOption { return geojson.WithBBox(bbox) }

func WithID(id any) geojson.FeatureOption { return geojson.WithID(id) }

func GetType(obj geojson.GeoJSON) string { return geojson.GetType(obj) }

func GetGeometry(obj any) (geojson.Geometry, error) { return geojson.GetGeometry(obj) }

func GetCoord(obj any) (geojson.Position, error) { return geojson.GetCoord(obj) }

func GetCoords(obj any) (any, error) { return geojson.GetCoords(obj) }

func GetBBox(obj geojson.GeoJSON) []float64 { return geojson.GetBBox(obj) }

func CoordAll(obj any) ([]geojson.Position, error) { return geojson.CoordAll(obj) }

func CoordEach(obj any, fn func(coord geojson.Position, index int) error) error {
	return meta.CoordEach(obj, fn)
}

func CoordReduce[T any](obj any, fn func(acc T, coord geojson.Position, index int) (T, error), initial T) (T, error) {
	return meta.CoordReduce(obj, fn, initial)
}

func FeatureEach(obj any, fn func(f *geojson.Feature, index int) error) error {
	return meta.FeatureEach(obj, fn)
}

func FeatureReduce[T any](obj any, initial T, fn func(acc T, f *geojson.Feature, index int) (T, error)) (T, error) {
	return meta.FeatureReduce(obj, initial, fn)
}

func GeomEach(obj any, fn func(geom geojson.Geometry, index int) error) error {
	return meta.GeomEach(obj, fn)
}

func GeomReduce[T any](obj any, initial T, fn func(acc T, geom geojson.Geometry, index int) (T, error)) (T, error) {
	return meta.GeomReduce(obj, initial, fn)
}

func PropEach(obj any, fn func(props map[string]any, index int) error) error {
	return meta.PropEach(obj, fn)
}

func PropReduce[T any](obj any, initial T, fn func(acc T, props map[string]any, index int) (T, error)) (T, error) {
	return meta.PropReduce(obj, initial, fn)
}

func FlattenEach(obj any, fn func(f *geojson.Feature, index int) error) error {
	return meta.FlattenEach(obj, fn)
}

func FlattenReduce[T any](obj any, initial T, fn func(acc T, f *geojson.Feature, index int) (T, error)) (T, error) {
	return meta.FlattenReduce(obj, initial, fn)
}

func Distance(from, to any, units ...measurement.Unit) (float64, error) {
	return measurement.Distance(from, to, units...)
}

func Bearing(from, to any) (float64, error) {
	return measurement.Bearing(from, to)
}

func RhumbBearing(from, to any) (float64, error) {
	return measurement.RhumbBearing(from, to)
}

func Destination(origin any, distance float64, bearing float64, units ...measurement.Unit) (*geojson.Feature, error) {
	return measurement.Destination(origin, distance, bearing, units...)
}

func RhumbDistance(from, to any, units ...measurement.Unit) (float64, error) {
	return measurement.RhumbDistance(from, to, units...)
}

func RhumbDestination(origin any, distance float64, bearing float64, units ...measurement.Unit) (*geojson.Feature, error) {
	return measurement.RhumbDestination(origin, distance, bearing, units...)
}

func Midpoint(from, to any) (*geojson.Feature, error) {
	return measurement.Midpoint(from, to)
}

func Length(geom any, units ...measurement.Unit) (float64, error) {
	return measurement.Length(geom, units...)
}

func Area(geom any) (float64, error) {
	return measurement.Area(geom)
}

func Along(line any, distance float64, units ...measurement.Unit) (*geojson.Feature, error) {
	return measurement.Along(line, distance, units...)
}

func GreatCircle(from, to any, options ...any) (*geojson.Feature, error) {
	return measurement.GreatCircle(from, to, options...)
}

func PointToLineDistance(point, line any, units ...measurement.Unit) (float64, error) {
	return measurement.PointToLineDistance(point, line, units...)
}

func NearestPointOnLine(line any, point any) (*geojson.Feature, error) {
	return measurement.NearestPointOnLine(line, point)
}

func BBox(obj any) ([]float64, error) {
	return bbox.BBox(obj)
}

func BBoxPolygon(b []float64) (*geojson.Feature, error) {
	return bbox.BBoxPolygon(b)
}

func Envelope(obj any) (*geojson.Feature, error) {
	return bbox.Envelope(obj)
}

func Center(obj any) (*geojson.Feature, error) {
	return center.Center(obj)
}

func Centroid(obj any) (*geojson.Feature, error) {
	return center.Centroid(obj)
}

func CenterMean(obj any, properties map[string]any, weight ...string) (*geojson.Feature, error) {
	return center.CenterMean(obj, properties, weight...)
}

func CenterOfMass(obj any) (*geojson.Feature, error) {
	return center.CenterOfMass(obj)
}

type TranslateOptions = transform.TranslateOptions

func TransformRotate(geom any, angle float64, pivot ...geojson.Position) (*geojson.Feature, error) {
	return transform.TransformRotate(geom, angle, pivot...)
}

func TransformScale(geom any, factor float64, pivot ...geojson.Position) (*geojson.Feature, error) {
	return transform.TransformScale(geom, factor, pivot...)
}

func TransformScaleXY(geom any, xFactor, yFactor float64, pivot ...geojson.Position) (*geojson.Feature, error) {
	return transform.TransformScaleXY(geom, xFactor, yFactor, pivot...)
}

func TransformTranslate(geom any, dx, dy float64, options ...*transform.TranslateOptions) (*geojson.Feature, error) {
	return transform.TransformTranslate(geom, dx, dy, options...)
}

func Flip(geom any) (*geojson.Feature, error) {
	return transform.Flip(geom)
}

func Truncate(geom any, precision int, coordinates ...int) (*geojson.Feature, error) {
	return transform.Truncate(geom, precision, coordinates...)
}

func CleanCoords(geom any) (*geojson.Feature, error) {
	return transform.CleanCoords(geom)
}

func Rewind(geom any, reversed ...bool) (*geojson.Feature, error) {
	return transform.Rewind(geom, reversed...)
}

func ToMercator(geom any) (*geojson.Feature, error) {
	return transform.ToMercator(geom)
}

func ToWGS84(geom any) (*geojson.Feature, error) {
	return transform.ToWGS84(geom)
}

func Square(bbx []float64) ([]float64, error) {
	return bbox.Square(bbx)
}

func CenterMedian(obj any, properties map[string]any, weight ...string) (*geojson.Feature, error) {
	return center.CenterMedian(obj, properties, weight...)
}

func NearestPoint(targetPt any, points any) (*geojson.Feature, error) {
	return measurement.NearestPoint(targetPt, points)
}

func PointInPolygon(point any, polygon any) (bool, error) {
	return boolean.PointInPolygon(point, polygon)
}

func PointOnLine(point any, line any, ignoreEndpoints bool) (bool, error) {
	return boolean.PointOnLine(point, line, ignoreEndpoints)
}

func SegmentIntersect(a, b, c, d geojson.Position) bool {
	return boolean.SegmentIntersect(a, b, c, d)
}

func Clockwise(ring []geojson.Position) bool {
	return boolean.Clockwise(ring)
}

func Contains(geom1, geom2 any) (bool, error) {
	return boolean.Contains(geom1, geom2)
}

func Within(geom1, geom2 any) (bool, error) {
	return boolean.Within(geom1, geom2)
}

func Intersects(geom1, geom2 any) (bool, error) {
	return boolean.Intersects(geom1, geom2)
}

func Disjoint(geom1, geom2 any) (bool, error) {
	return boolean.Disjoint(geom1, geom2)
}

func Touches(geom1, geom2 any) (bool, error) {
	return boolean.Touches(geom1, geom2)
}

func Crosses(geom1, geom2 any) (bool, error) {
	return boolean.Crosses(geom1, geom2)
}

func Overlap(geom1, geom2 any) (bool, error) {
	return boolean.Overlap(geom1, geom2)
}

func Valid(geom any) (bool, error) {
	return boolean.Valid(geom)
}

func BooleanConcave(geom any) (bool, error) {
	return boolean.Concave(geom)
}

func BooleanEqual(geom1, geom2 any) (bool, error) {
	return boolean.BooleanEqual(geom1, geom2)
}

func BooleanParallel(line1, line2 any) (bool, error) {
	return boolean.BooleanParallel(line1, line2)
}

type GridOptions = grid.GridOptions
type CircleOptions = shapes.CircleOptions
type EllipseOptions = shapes.EllipseOptions
type BezierOptions = shapes.BezierOptions
type RandomOptions = shapes.RandomOptions

func HexGrid(bbox []float64, cellSide float64, units measurement.Unit, options ...grid.GridOptions) (*geojson.FeatureCollection, error) {
	return grid.HexGrid(bbox, cellSide, units, options...)
}

func PointGrid(bbox []float64, cellSide float64, units measurement.Unit, options ...grid.GridOptions) (*geojson.FeatureCollection, error) {
	return grid.PointGrid(bbox, cellSide, units, options...)
}

func SquareGrid(bbox []float64, cellSide float64, units measurement.Unit, options ...grid.GridOptions) (*geojson.FeatureCollection, error) {
	return grid.SquareGrid(bbox, cellSide, units, options...)
}

func TriangleGrid(bbox []float64, cellSide float64, units measurement.Unit, options ...grid.GridOptions) (*geojson.FeatureCollection, error) {
	return grid.TriangleGrid(bbox, cellSide, units, options...)
}

func Circle(center any, radius float64, options ...shapes.CircleOptions) (*geojson.Feature, error) {
	return shapes.Circle(center, radius, options...)
}

func Ellipse(center any, xSemiAxis, ySemiAxis float64, options ...shapes.EllipseOptions) (*geojson.Feature, error) {
	return shapes.Ellipse(center, xSemiAxis, ySemiAxis, options...)
}

func BezierSpline(line any, options ...shapes.BezierOptions) (*geojson.Feature, error) {
	return shapes.BezierSpline(line, options...)
}

func RandomPosition(bbox []float64) geojson.Position {
	return shapes.RandomPosition(bbox)
}

func RandomPoint(count int, options ...shapes.RandomOptions) (*geojson.FeatureCollection, error) {
	return shapes.RandomPoint(count, options...)
}

func RandomLineString(count int, options ...shapes.RandomOptions) (*geojson.FeatureCollection, error) {
	return shapes.RandomLineString(count, options...)
}

func RandomPolygon(count int, options ...shapes.RandomOptions) (*geojson.FeatureCollection, error) {
	return shapes.RandomPolygon(count, options...)
}

func ConvertLength(length float64, from, to measurement.Unit) float64 {
	return measurement.ConvertLength(length, from, to)
}

type TinOptions = interpolation.TinOptions
type InterpolateOptions = interpolation.InterpolateOptions
type DbscanOptions = clusters.DbscanOptions

func Sample(fc *geojson.FeatureCollection, n int) (*geojson.FeatureCollection, error) {
	return interpolation.Sample(fc, n)
}

func Tin(points *geojson.FeatureCollection, options ...interpolation.TinOptions) (*geojson.FeatureCollection, error) {
	return interpolation.Tin(points, options...)
}

func Interpolate(points *geojson.FeatureCollection, cellSide float64, units measurement.Unit, property string, options ...interpolation.InterpolateOptions) (*geojson.FeatureCollection, error) {
	return interpolation.Interpolate(points, cellSide, units, property, options...)
}

func PlanarDistance(a, b geojson.Position) float64 {
	return interpolation.PlanarDistance(a, b)
}

func ClustersKMeans(fc *geojson.FeatureCollection, k int) (*geojson.FeatureCollection, error) {
	return clusters.ClustersKMeans(fc, k)
}

func ClustersDbscan(fc *geojson.FeatureCollection, radius float64, options ...clusters.DbscanOptions) (*geojson.FeatureCollection, error) {
	return clusters.ClustersDbscan(fc, radius, options...)
}

func Dissolve(fc *geojson.FeatureCollection, property string) (*geojson.FeatureCollection, error) {
	return clusters.Dissolve(fc, property)
}

func Tag(points, polygons *geojson.FeatureCollection, field, outField string) (*geojson.FeatureCollection, error) {
	return data.Tag(points, polygons, field, outField)
}

func Collect(polygons, points *geojson.FeatureCollection, inField, outField string) (*geojson.FeatureCollection, error) {
	return data.Collect(polygons, points, inField, outField)
}

func Simplify(geom any, tolerance float64, highQuality bool) (*geojson.Feature, error) {
	return simplify.Simplify(geom, tolerance, highQuality)
}

func ConvexHull(geom any) (*geojson.Feature, error) {
	return simplify.ConvexHull(geom)
}

type OpType = polyclip.OpType

const (
	OpUnion       = polyclip.OpUnion
	OpIntersect   = polyclip.OpIntersect
	OpDifference  = polyclip.OpDifference
	OpXor         = polyclip.OpXor
)

func PolygonUnion(poly1, poly2 any) (*geojson.Feature, error) {
	return polyclip.PolygonUnion(poly1, poly2)
}

func PolygonIntersect(poly1, poly2 any) (*geojson.Feature, error) {
	return polyclip.PolygonIntersect(poly1, poly2)
}

func PolygonDifference(poly1, poly2 any) (*geojson.Feature, error) {
	return polyclip.PolygonDifference(poly1, poly2)
}

func PolygonXor(poly1, poly2 any) (*geojson.Feature, error) {
	return polyclip.PolygonXor(poly1, poly2)
}

type IsobandsOptions = isobands.IsobandsOptions

func Isobands(points *geojson.FeatureCollection, options ...isobands.IsobandsOptions) (*geojson.FeatureCollection, error) {
	return isobands.Isobands(points, options...)
}

type IsolinesOptions = isolines.IsolinesOptions

func Isolines(points *geojson.FeatureCollection, options ...isolines.IsolinesOptions) (*geojson.FeatureCollection, error) {
	return isolines.Isolines(points, options...)
}

// --- NEWLY ADDED RE-EXPORTS ---

func Clone(geom any) (*geojson.Feature, error) {
	return misc.Clone(geom)
}

func Combine(fc *geojson.FeatureCollection) (*geojson.FeatureCollection, error) {
	return misc.Combine(fc)
}

func Explode(geom any) (*geojson.FeatureCollection, error) {
	return misc.Explode(geom)
}

func Flatten(geom any) (*geojson.FeatureCollection, error) {
	return misc.Flatten(geom)
}

func Planepoint(point any, triangle any) (float64, error) {
	return misc.Planepoint(point, triangle)
}

func PointsWithinPolygon(points, polygon any) (*geojson.FeatureCollection, error) {
	return misc.PointsWithinPolygon(points, polygon)
}

func Tesselate(poly any) (*geojson.FeatureCollection, error) {
	return misc.Tesselate(poly)
}

func Buffer(geom any, radius float64, units measurement.Unit, steps ...int) (*geojson.Feature, error) {
	return buffer.Buffer(geom, radius, units, steps...)
}

func LineOffset(line any, distance float64, units measurement.Unit) (*geojson.Feature, error) {
	return lineoffset.LineOffset(line, distance, units)
}

func PolygonToLine(poly any) (*geojson.Feature, error) {
	return convert.PolygonToLine(poly)
}

func LineToPolygon(line any) (*geojson.Feature, error) {
	return convert.LineToPolygon(line)
}

func PolygonSmooth(poly any, iterations int) (*geojson.Feature, error) {
	return smooth.PolygonSmooth(poly, iterations)
}

func PolygonTangents(point any, polygon any) (*geojson.FeatureCollection, error) {
	return tangents.PolygonTangents(point, polygon)
}

func Kinks(geom any) (*geojson.FeatureCollection, error) {
	return kinks.Kinks(geom)
}

func Mask(outer any, inner any) (*geojson.Feature, error) {
	return mask.Mask(outer, inner)
}

func UnkinkPolygon(poly any) (*geojson.FeatureCollection, error) {
	return unkink.UnkinkPolygon(poly)
}

type VoronoiOptions = voronoi.VoronoiOptions

func Voronoi(points *geojson.FeatureCollection, options ...voronoi.VoronoiOptions) (*geojson.FeatureCollection, error) {
	return voronoi.Voronoi(points, options...)
}

func LineIntersect(line1, line2 any) (*geojson.FeatureCollection, error) {
	return lines.LineIntersect(line1, line2)
}

func LineSegment(geom any) (*geojson.FeatureCollection, error) {
	return lines.LineSegment(geom)
}

func LineOverlap(line1, line2 any) (*geojson.FeatureCollection, error) {
	return lines.LineOverlap(line1, line2)
}

func LineSlice(point1, point2, line any) (*geojson.Feature, error) {
	return lines.LineSlice(point1, point2, line)
}

func LineSliceAlong(line any, startDist, endDist float64, units measurement.Unit) (*geojson.Feature, error) {
	return lines.LineSliceAlong(line, startDist, endDist, units)
}

func LineChunk(line any, segmentLength float64, units measurement.Unit) (*geojson.FeatureCollection, error) {
	return lines.LineChunk(line, segmentLength, units)
}

func LineSplit(line any, point any) (*geojson.FeatureCollection, error) {
	return lines.LineSplit(line, point)
}

func LineArc(center any, radius float64, bearing1, bearing2 float64, steps int, units measurement.Unit) (*geojson.Feature, error) {
	return lines.LineArc(center, radius, bearing1, bearing2, steps, units)
}

func Sector(center any, radius float64, bearing1, bearing2 float64, steps int, units measurement.Unit) (*geojson.Feature, error) {
	return lines.Sector(center, radius, bearing1, bearing2, steps, units)
}

// --- NEW MODULES ---

func Angle(startPt, midPt, endPt any) (float64, error) {
	return measurement.Angle(startPt, midPt, endPt)
}

type ConcaveOptions = concave.ConcaveOptions

func ConcaveHull(points *geojson.FeatureCollection, options ...concave.ConcaveOptions) (*geojson.Feature, error) {
	return concave.ConcaveHull(points, options...)
}

func DirectionalMean(geom any) (float64, error) {
	return directionalmean.DirectionalMean(geom)
}

func DistanceWeight(points *geojson.FeatureCollection, options ...float64) ([][]float64, error) {
	return distanceweight.DistanceWeight(points, options...)
}

func MoranIndex(points *geojson.FeatureCollection, property string) (float64, error) {
	return moranindex.MoranIndex(points, property)
}

type NearestNeighborResult = nearestneighbor.NearestNeighborResult

func NearestNeighborAnalysis(points *geojson.FeatureCollection) (*nearestneighbor.NearestNeighborResult, error) {
	return nearestneighbor.NearestNeighborAnalysis(points)
}

func NearestPointToLine(line any, points *geojson.FeatureCollection, units ...measurement.Unit) (*geojson.Feature, error) {
	return measurement.NearestPointToLine(line, points, units...)
}

func PointOnFeature(geom any) (*geojson.Feature, error) {
	return center.PointOnFeature(geom)
}

func PointToPolygonDistance(point, polygon any, units ...measurement.Unit) (float64, error) {
	return measurement.PointToPolygonDistance(point, polygon, units...)
}

func Polygonize(geom any) (*geojson.FeatureCollection, error) {
	return polygonize.Polygonize(geom)
}

func RectangleGrid(bbox []float64, cellWidth, cellHeight float64, units measurement.Unit, options ...grid.GridOptions) (*geojson.FeatureCollection, error) {
	return grid.RectangleGrid(bbox, cellWidth, cellHeight, units, options...)
}

type QuadratAnalysisResult = quadratanalysis.QuadratAnalysisResult

func QuadratAnalysis(points *geojson.FeatureCollection, gridSize int) (*quadratanalysis.QuadratAnalysisResult, error) {
	return quadratanalysis.QuadratAnalysis(points, gridSize)
}

func NewRBush() *rbush.RBush {
	return rbush.NewRBush()
}

type ShortestPathOptions = shortestpath.ShortestPathOptions

func ShortestPath(start, end any, network *geojson.FeatureCollection, options ...shortestpath.ShortestPathOptions) (*geojson.Feature, error) {
	return shortestpath.ShortestPath(start, end, network, options...)
}

type StandardDeviationalEllipseResult = standarddeviationalellipse.StandardDeviationalEllipseResult

func StandardDeviationalEllipse(points *geojson.FeatureCollection, options ...int) (*standarddeviationalellipse.StandardDeviationalEllipseResult, error) {
	return standarddeviationalellipse.StandardDeviationalEllipse(points, options...)
}
