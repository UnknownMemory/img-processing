package process

import "github.com/cshum/vipsgen/vips"

func Transform(object []byte, operations interface{}) (*vips.Image, error) {
	img, err := vips.NewImageFromBuffer(object, nil)
	if err != nil {
		return nil, err
	}
	return img, nil
}
