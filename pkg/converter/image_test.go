package converter

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestConvertImage(t *testing.T) {
	t.Run("JpgToPng", func(t *testing.T) {
		file, err := os.Open("./testdata/test.jpg")
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				log.Fatal(err)
			}
		}()

		reader := bufio.NewReader(file)

		_, err = ProcessImage(reader, ioutil.Discard, JPG, 75)
		if err != nil {
			log.Fatal(err)
		}
	})
	t.Run("PngToJpg", func(t *testing.T) {
		file, err := os.Open("./testdata/test.png")
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				log.Fatal(err)
			}
		}()

		reader := bufio.NewReader(file)

		_, err = ProcessImage(reader, ioutil.Discard, PNG, 75)
		if err != nil {
			log.Fatal(err)
		}
	})
}
