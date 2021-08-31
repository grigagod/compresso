package auth

import "net/http"

type Handlers interface {
	Register() http.Handler
	Login() http.Handler
}
