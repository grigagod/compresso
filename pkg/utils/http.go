package utils

import (
	"encoding/json"
	"net/http"
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
