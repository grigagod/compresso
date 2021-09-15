package converter

import (
	"context"
	"errors"
	"io"

	fluentffmpeg "github.com/modfy/fluent-ffmpeg"
)

var (
	ErrVideoCRF           = errors.New("CRF value is out of range")
	ErrProcessInterrupted = errors.New("process interrupted")
)

// ProcessVideo process video from source with given options.
func ProcessVideo(ctx context.Context, src io.Reader, dst io.Writer, opts VideoOpts) error {
	if opts.CRF > 51 || opts.CRF < 0 {
		return ErrVideoCRF
	}

	cmd := fluentffmpeg.NewCommand("").
		PipeInput(src).
		ConstantRateFactor(opts.CRF).
		OutputFormat(string(opts.TargetFormat))

	err := cmd.PipeOutput(dst).RunWithContext(ctx)
	return err
}
