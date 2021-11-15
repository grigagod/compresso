package http

import (
	"net/http"

	"github.com/grigagod/compresso/internal/auth"
	"github.com/grigagod/compresso/internal/httper"
	"github.com/grigagod/compresso/internal/utils"
)

type AuthHandlers struct {
	authUC auth.UseCase
}

func NewAuthHandlers(authUC auth.UseCase) *AuthHandlers {
	return &AuthHandlers{
		authUC: authUC,
	}
}

// Register godoc
// @Summary Register new user
// @Description register new user, returns user and token
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body AuthRequest true "user credentials"
// @Success 201 {object} auth.UserWithToken
// @Failure 400 {string} string "Bad request msg"
// @Router /register [post].
func (h *AuthHandlers) Register() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) error {
		var req AuthRequest

		if err := utils.StructScan(r.Body, &req); err != nil {
			return httper.NewBadRequestError(err)
		}

		err := utils.ValidateStruct(&req)
		if err != nil {
			return httper.ParseValidatorError(err)
		}

		user, err := h.authUC.Register(r.Context(), &auth.User{
			Username: req.Username,
			Password: req.Password,
		})
		if err != nil {
			return err
		}

		return utils.RespondWithJSON(w, http.StatusCreated, user)
	}

	return httper.HandlerWithError(fn)
}

// Login godoc
// @Summary Login new user
// @Description login user, returns user and token
// @Tags Auth
// @Accept json
// @Produce json
// @Param creds body AuthRequest true "user credentials"
// @Success 200 {object} auth.UserWithToken
// @Failure 400 {string} string "Bad request msg"
// @Failure 401 {string} string "Wrong credentials"
// @Failure 404 {string} string "User with such username is not found"
// @Router /login [post].
func (h *AuthHandlers) Login() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) error {
		var req AuthRequest

		if err := utils.StructScan(r.Body, &req); err != nil {
			return httper.NewBadRequestError(err)
		}

		user, err := h.authUC.Login(r.Context(), &auth.User{
			Username: req.Username,
			Password: req.Password,
		})
		if err != nil {
			return err
		}

		return utils.RespondWithJSON(w, http.StatusOK, user)
	}

	return httper.HandlerWithError(fn)
}
