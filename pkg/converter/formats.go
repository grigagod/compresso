// Package converter implements media converting functionality.
package converter

// Media formats supported by pkg definitions

// ImageFormat used for representing image files format.
type ImageFormat string

// Allowed image formats.
const (
	PNG ImageFormat = "png"
	JPG ImageFormat = "jpg"
)

// VideoFormat used for representing image files format.
type VideoFormat string

// Allowed video formats.
const (
	MKV  VideoFormat = "matroska"
	WebM VideoFormat = "webm"
)

// AudioFormat used for representing image files format.
type AudioFormat string

// Allowed audio formats.
const (
	MP3 AudioFormat = "mp3"
	WAV AudioFormat = "wav"
)
