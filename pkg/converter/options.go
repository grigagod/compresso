package converter

type ImageOpts struct {
	CurrentFormat    ImageFormat
	TargetFormat     ImageFormat
	CompressionRatio int
}

func NewImageOpts(current, target ImageFormat, ratio int) ImageOpts {
	return ImageOpts{
		CurrentFormat:    current,
		TargetFormat:     target,
		CompressionRatio: ratio,
	}
}
