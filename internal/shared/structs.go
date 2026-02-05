package shared

type ImageTransform struct {
	ImageID         string          `json:"image_id"`
	Transformations Transformations `json:"transformations"`
}

type Transformations struct {
	Resize *struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"resize"`
}
