package shared

type ImageTransform struct {
	ImageID         string          `json:"image_id"`
	Transformations Transformations `json:"transformations"`
}

type Transformations struct {
	Resize *float64 `json:"resize"`
}
