package httper

import "net/http"

type HandlerWithError func(w http.ResponseWriter, r *http.Request) error

func (fn HandlerWithError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil {
		err, ok := e.(Error)
		if !ok {
			http.Error(w, e.Error(), http.StatusInternalServerError)
			return
		}

		http.Error(w, err.Error(), err.Status())
	}
}
