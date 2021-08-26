package auth

import "net/http"

type Handlers interface {
	Register() http.HandlerFunc
	Login() http.HandlerFunc
}
