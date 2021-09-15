package converter

type imageFormat string

const (
	PNG imageFormat = "png"
	JPG imageFormat = "jpg"
)

type videoFormat string

const (
	MKV  videoFormat = "matroska"
	WebM videoFormat = "webm"
)

type audioFormat string

const (
	MP3 audioFormat = "mp3"
	WAV audioFormat = "wav"
)
