package middleware

import (
	"net/http"

	"github.com/grigagod/compresso/internal/httper"
	"github.com/grigagod/compresso/internal/utils"
)

// JWTAuth is a middleware that extracts jwt claims from request and place them in the context.
func JWTAuth(jwtSecretKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) error {
			claims, err := utils.ExtractJWTFromRequest(r, jwtSecretKey)
			if err != nil {
				return httper.ParseJWTError(err)
			}
			next.ServeHTTP(w, r.WithContext(utils.ContextWithUserID(r.Context(), claims.ID)))

			return nil
		}

		return httper.HandlerWithError(fn)
	}
}
