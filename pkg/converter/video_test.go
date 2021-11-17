//go:build integration

package converter

import (
	"bufio"
	"context"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

func TestIntegration_ProcessVideo(t *testing.T) {
	t.Run("MKVtoMKV", func(t *testing.T) {
		file, err := os.Open("./testdata/lorenz.mkv")
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				log.Fatal(err)
			}
		}()

		reader := bufio.NewReader(file)

		err = ProcessVideo(context.Background(), reader, ioutil.Discard, MKV, 18)
		if err != nil {
			log.Fatal(err)
		}
	})
	t.Run("MKVtoMKVWithTimeout", func(t *testing.T) {
		file, err := os.Open("./testdata/lorenz.mkv")
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				log.Fatal(err)
			}
		}()

		reader := bufio.NewReader(file)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		err = ProcessVideo(ctx, reader, ioutil.Discard, MKV, 17)
		if ctx.Err() != context.DeadlineExceeded {
			log.Print(err)
			log.Fatal(ctx.Err())
		}

	})
}
