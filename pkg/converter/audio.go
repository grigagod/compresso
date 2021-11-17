package converter

import (
	"context"
	"errors"
	"io"

	fluentffmpeg "github.com/modfy/fluent-ffmpeg"
)

// Package defined errors for audio processing.
var (
	ErrAudioFormat = errors.New("unsupported audio format")
)

// ProcessAudio process(convert) audio from source with given formats.
func ProcessAudio(ctx context.Context, src io.Reader, dst io.Writer, currentFormat, targetFormat AudioFormat) error {
	return fluentffmpeg.NewCommand("").
		PipeInput(src).
		PipeOutput(dst).
		InputOptions("-f", string(currentFormat)).
		OutputFormat(string(targetFormat)).
		RunWithContext(ctx)
}
