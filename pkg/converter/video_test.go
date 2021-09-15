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

func TestProcessVideo(t *testing.T) {
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

		err = ProcessVideo(context.Background(), reader, ioutil.Discard, NewVideoOpts(MKV, MKV, 18))
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
		err = ProcessVideo(ctx, reader, ioutil.Discard, NewVideoOpts(MKV, MKV, 17))
		defer cancel()
		if ctx.Err() != context.DeadlineExceeded {
			log.Print(err)
			log.Fatal(ctx.Err())
		}

	})
}
