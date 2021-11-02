package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/grigagod/compresso/internal/httper"
)

type ContentTypeCtxKey struct{}

func ContentType(contentTypes ...string) func(next http.Handler) http.Handler {
	allowedContentTypes := make(map[string]struct{}, len(contentTypes))
	for _, ctype := range contentTypes {
		allowedContentTypes[strings.TrimSpace(strings.ToLower(ctype))] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) error {
			if r.ContentLength == 0 {
				// skip check for empty content body
				next.ServeHTTP(w, r)
				return nil
			}

			s := strings.ToLower(strings.TrimSpace(r.Header.Get("Content-Type")))
			if i := strings.Index(s, ";"); i > -1 {
				s = s[0:i]
			}

			if _, ok := allowedContentTypes[s]; ok {
				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ContentTypeCtxKey{}, s)))
				return nil
			}

			return httper.NewNotAllowedMediaMsg()
		}
		return httper.HandlerWithError(fn)
	}
}
