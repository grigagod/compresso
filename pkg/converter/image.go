package converter

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
)

var (
	ErrDecodeImage    = errors.New("can't decode image")
	ErrEncodeImage    = errors.New("can't encode image")
	ErrImageFormat    = errors.New("unsupported image format")
	ErrImageConvesion = errors.New("unsupported image conversion")
)

// Global encoder to reuse buffer pool
var enc png.Encoder

// ProcessImage process image from source with given options.
func ProcessImage(src io.Reader, dst io.Writer, currentFormat imageFormat, ratio int) error {
	img, err := DecodeImage(src, currentFormat)
	if err != nil {
		return err
	}

	switch currentFormat {
	case PNG:
		err = jpeg.Encode(dst, img, &jpeg.Options{
			Quality: ratio,
		})
		if err != nil {
			return ErrEncodeImage
		}

		return nil
	case JPG:
		enc.CompressionLevel = ratioToCompression(ratio)

		err = enc.Encode(dst, img)
		if err != nil {
			return ErrEncodeImage
		}

		return nil
	default:
		return ErrImageFormat
	}
}

// DecodeImage decode image of supported formats.
func DecodeImage(src io.Reader, current imageFormat) (image.Image, error) {
	switch current {
	case PNG:
		img, err := png.Decode(src)
		if err != nil {
			return nil, ErrDecodeImage
		}
		return img, err
	case JPG:
		img, err := jpeg.Decode(src)
		if err != nil {
			return nil, ErrDecodeImage
		}
		return img, nil
	default:
		return nil, ErrImageFormat
	}
}

func ratioToCompression(ratio int) png.CompressionLevel {
	switch {
	case ratio < 25:
		return png.BestCompression
	case ratio > 24 && ratio < 50:
		return png.DefaultCompression
	case ratio > 49 && ratio < 75:
		return png.BestSpeed
	default:
		return png.NoCompression
	}
}
