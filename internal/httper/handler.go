package httper

import (
	"net/http"

	"github.com/grigagod/compresso/internal/utils"
)

type HandlerWithError func(w http.ResponseWriter, r *http.Request) error

func (fn HandlerWithError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil {
		err, ok := e.(Error)
		if !ok {
			utils.RespondWithText(w, http.StatusInternalServerError, e.Error())
			return
		}
		utils.RespondWithText(w, err.Status(), err.Error())
	}
}
