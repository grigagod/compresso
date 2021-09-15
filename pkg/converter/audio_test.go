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

func TestProcessAudio(t *testing.T) {
	t.Run("MP3toWAV", func(t *testing.T) {
		file, err := os.Open("./testdata/lorenz.mp3")
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				log.Fatal(err)
			}
		}()

		reader := bufio.NewReader(file)

		err = ProcessAudio(context.Background(), reader, ioutil.Discard, MP3, WAV)
		if err != nil {
			log.Fatal(err)
		}
	})
	t.Run("WAVtoMP3", func(t *testing.T) {
		file, err := os.Open("./testdata/lorenz.wav")
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				log.Fatal(err)
			}
		}()

		reader := bufio.NewReader(file)

		err = ProcessAudio(context.Background(), reader, ioutil.Discard, WAV, MP3)
		if err != nil {
			log.Fatal(err)
		}

	})
	t.Run("WithTimeout", func(t *testing.T) {
		file, err := os.Open("./testdata/lorenz.mp3")
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				log.Fatal(err)
			}
		}()

		reader := bufio.NewReader(file)
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1)
		defer cancel()
		err = ProcessAudio(ctx, reader, ioutil.Discard, WAV, MP3)
		if ctx.Err() != context.DeadlineExceeded {
			log.Print(err)
			log.Fatal(ctx.Err())
		}

	})

}
