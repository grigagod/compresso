// Package utils provides helper functions.
package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

// StructScan decodes(as json) into model.
func StructScan(r io.Reader, model interface{}) error {
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(model); err != nil {
		return err
	}

	return nil
}

// RespondWithText responds with plain text and given code.
func RespondWithText(w http.ResponseWriter, code int, msg string) error {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	_, err := w.Write([]byte(msg))

	return err
}

// RespondWithJSON marshal(into json) model and responds with given code.
func RespondWithJSON(w http.ResponseWriter, code int, model interface{}) error {
	resp, err := json.Marshal(model)

	if err != nil {
		return RespondWithText(w, http.StatusInternalServerError, err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(resp)

	return err
}
