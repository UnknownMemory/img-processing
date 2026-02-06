package process

import (
	"github.com/cshum/vipsgen/vips"
	"github.com/unknownmemory/img-processing/internal/shared"
)

func Transform(object []byte, operations shared.Transformations) ([]byte, string, error) {
	img, err := vips.NewImageFromBuffer(object, nil)
	if err != nil {
		return nil, "", err
	}
	defer img.Close()

	format := img.Format()

	if operations.Resize != nil && *operations.Resize > 0 {
		err := img.Resize(*operations.Resize, nil)
		if err != nil {
			return nil, "", err
		}
	}

	imgBuffer, mime, err := export(format, img)
	if err != nil {
		return nil, "", err
	}

	return imgBuffer, mime, nil
}

func export(format vips.ImageType, img *vips.Image) ([]byte, string, error) {
	var buff []byte
	var mime string
	var err error

	switch format {
	case vips.ImageTypeJpeg:
		buff, err = img.JpegsaveBuffer(nil)
		mime = "image/jpeg"
	case vips.ImageTypePng:
		buff, err = img.PngsaveBuffer(nil)
		mime = "image/png"
	default:
		buff, err = img.JpegsaveBuffer(nil)
		mime = "image/jpeg"
	}
	return buff, mime, err
}
