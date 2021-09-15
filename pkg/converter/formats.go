package converter

type ImageFormat string

const (
	PNG ImageFormat = "png"
	JPG ImageFormat = "jpg"
)

type VideoFormat string

const (
	MKV  VideoFormat = "matroska"
	WebM VideoFormat = "webm"
)

type AudioFormat string

const (
	MP3 AudioFormat = "mp3"
	WAV AudioFormat = "wav"
)
