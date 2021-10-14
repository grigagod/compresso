package utils

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/grigagod/compresso/pkg/converter"
)

// StructScan decodes(as json) request body into model.
func StructScan(r *http.Request, model interface{}) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(model); err != nil {
		return err
	}

	return nil
}

// RespondWithError responds with plain text error and given code.
func RespondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-type", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	w.Write([]byte(msg))
}

// RespondWithJSON marshal(into json) model and responds with given code.
func RespondWithJSON(w http.ResponseWriter, code int, model interface{}) {
	resp, err := json.Marshal(model)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)
}

var allowedVideoContentTypes = map[string]converter.VideoFormat{
	"video/x-matroska": converter.MKV,
	"video/webm":       converter.WebM,
}

var allowedVideoFormats = map[string]converter.VideoFormat{
	"matroska": converter.MKV,
	"webm":     converter.WebM,
}

func DetectVideoFormatFromHeader(header string) (converter.VideoFormat, error) {
	f, ok := allowedVideoContentTypes[header]
	if !ok {
		return converter.VideoFormat(""), errors.New("this content type is not allowed")
	}

	return f, nil
}

func DetectVideoFormat(format string) (converter.VideoFormat, error) {
	f, ok := allowedVideoFormats[format]
	if !ok {
		return converter.VideoFormat(""), errors.New("this content type is not allowed")
	}

	return f, nil
}
