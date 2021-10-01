package converter

import (
	"context"
	"errors"
	"io"

	fluentffmpeg "github.com/modfy/fluent-ffmpeg"
)

var (
	ErrVideoCRF = errors.New("CRF value is out of range")
)

// ProcessVideo process video from source with given options.
func ProcessVideo(ctx context.Context, src io.Reader, dst io.Writer, targetFormat VideoFormat, crf int) error {
	if crf > 51 || crf < 0 {
		return ErrVideoCRF
	}

	cmd := fluentffmpeg.NewCommand("").
		PipeInput(src).
		ConstantRateFactor(crf).
		OutputFormat(string(targetFormat))

	err := cmd.PipeOutput(dst).RunWithContext(ctx)
	return err
}
