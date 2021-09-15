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

type VideoOpts struct {
	CurrentFormat VideoFormat
	TargetFormat  VideoFormat
	CRF           int
}

func NewVideoOpts(current, target VideoFormat, crf int) VideoOpts {
	return VideoOpts{

		CurrentFormat: current,
		TargetFormat:  target,
		CRF:           crf,
	}
}
