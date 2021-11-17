package utils

import (
	"github.com/grigagod/compresso/pkg/converter"
)

const (
	JSONContentType = "application/json"
	TextContentType = "text/plain; charset=utf-8"
)

var AllowedVideoContentTypes = map[string]converter.VideoFormat{
	"video/x-matroska": converter.MKV,
	"video/webm":       converter.WebM,
}

var AllowedContentTypes = []string{"video/x-matroska", "video/webm", "application/json"}

var AllowedVideoFormats = map[string]converter.VideoFormat{
	"matroska": converter.MKV,
	"webm":     converter.WebM,
}

var VideoFormatsToMIME = map[converter.VideoFormat]string{
	converter.WebM: "video/webm",
	converter.MKV:  "video/x-matroska",
}

// DetectVideoFormatFromHeader detect converter.VideoFormat from header.
func DetectVideoFormatFromHeader(header string) (converter.VideoFormat, bool) {
	f, ok := AllowedVideoContentTypes[header]
	return f, ok
}

// DetectVideoFormat detect converter.VideoFormat from string.
func DetectVideoFormat(format string) (converter.VideoFormat, bool) {
	f, ok := AllowedVideoFormats[format]
	return f, ok
}

func DetectVideoMIMEType(format converter.VideoFormat) (string, bool) {
	f, ok := VideoFormatsToMIME[format]
	return f, ok
}
