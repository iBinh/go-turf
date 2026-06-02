package geojson

type FeatureOption func(*Feature)

func WithBBox(bbox []float64) FeatureOption {
	return func(f *Feature) {
		f.SetBBox(bbox)
	}
}

func WithID(id any) FeatureOption {
	return func(f *Feature) {
		f.ID = id
	}
}
