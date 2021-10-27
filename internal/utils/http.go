// Package utils provides helper functions.
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
func RespondWithError(w http.ResponseWriter, code int, msg string) error {
	w.Header().Set("Content-type", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	_, err := w.Write([]byte(msg))

	return err
}

// RespondWithJSON marshal(into json) model and responds with given code.
func RespondWithJSON(w http.ResponseWriter, code int, model interface{}) error {
	resp, err := json.Marshal(model)

	if err != nil {
		return RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(resp)

	return err
}

var AllowedVideoContentTypes = map[string]converter.VideoFormat{
	"video/x-matroska": converter.MKV,
	"video/webm":       converter.WebM,
}

var AllowedVideoFormats = map[string]converter.VideoFormat{
	"matroska": converter.MKV,
	"webm":     converter.WebM,
}

// DetectVideoFormatFromHeader detect converter.VideoFormat from header.
func DetectVideoFormatFromHeader(header string) (converter.VideoFormat, error) {
	f, ok := AllowedVideoContentTypes[header]
	if !ok {
		return converter.VideoFormat(""), errors.New("this content type is not allowed")
	}

	return f, nil
}

// DetectVideoFormat detect converter.VideoFormat from string.
func DetectVideoFormat(format string) (converter.VideoFormat, error) {
	f, ok := AllowedVideoFormats[format]
	if !ok {
		return converter.VideoFormat(""), errors.New("this content type is not allowed")
	}

	return f, nil
}
