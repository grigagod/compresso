package converter

import (
	"bytes"
	"errors"
	"image/jpeg"
	"image/png"
	"io"
)

var (
	ErrDecodeImage = errors.New("can't decode image")
	ErrEncodeImage = errors.New("can't encode image")
	ErrImageFormat = errors.New("unsupported image format")
)

// Global encoder to reuse buffer pool
var enc png.Encoder

// ConvertImage convert image with specified quality(compression) ratio.
func ConvertImage(reader io.Reader, currentFormat Format, ratio int) (io.Reader, error) {
	switch currentFormat {
	case PNG:
		source, err := png.Decode(reader)
		if err != nil {
			return nil, ErrDecodeImage
		}

		buf := new(bytes.Buffer)
		err = jpeg.Encode(buf, source, &jpeg.Options{
			Quality: ratio,
		})
		if err != nil {
			return nil, ErrEncodeImage
		}

		return bytes.NewReader(buf.Bytes()), nil
	case JPG:
		source, err := jpeg.Decode(reader)
		if err != nil {
			return nil, ErrDecodeImage
		}

		buf := new(bytes.Buffer)
		enc.CompressionLevel = ratioToCompression(ratio)
		err = enc.Encode(buf, source)
		if err != nil {
			return nil, ErrEncodeImage
		}

		return bytes.NewReader(buf.Bytes()), nil
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
