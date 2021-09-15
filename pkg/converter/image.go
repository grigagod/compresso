package converter

import (
	"bytes"
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

// ProcessImage process image with specified options.
func ProcessImage(src io.Reader, dst io.Writer, currentFormat ImageFormat, ratio int) (io.Reader, error) {
	img, err := DecodeImage(src, currentFormat)
	if err != nil {
		return nil, err
	}

	switch currentFormat {
	case PNG:
		buf := new(bytes.Buffer)
		err = jpeg.Encode(buf, img, &jpeg.Options{
			Quality: ratio,
		})
		if err != nil {
			return nil, ErrEncodeImage
		}

		return bytes.NewReader(buf.Bytes()), nil
	case JPG:
		buf := new(bytes.Buffer)
		enc.CompressionLevel = ratioToCompression(ratio)

		err = enc.Encode(buf, img)
		if err != nil {
			return nil, ErrEncodeImage
		}

		return bytes.NewReader(buf.Bytes()), nil
	default:
		return nil, ErrImageFormat
	}
}

// DecodeImage decode image of supported formats.
func DecodeImage(input io.Reader, current ImageFormat) (image.Image, error) {
	switch current {
	case PNG:
		img, err := png.Decode(input)
		if err != nil {
			return nil, ErrDecodeImage
		}
		return img, err
	case JPG:
		img, err := jpeg.Decode(input)
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
