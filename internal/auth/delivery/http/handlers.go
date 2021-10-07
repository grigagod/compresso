package http

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/grigagod/compresso/internal/auth"
	"github.com/grigagod/compresso/internal/auth/config"
	"github.com/grigagod/compresso/internal/httper"
	"github.com/grigagod/compresso/pkg/utils"
)

type authHandlers struct {
	cfg    *config.Auth
	authUC auth.UseCase
}

func NewAuthHandlers(cfg *config.Auth, authUC auth.UseCase) auth.Handlers {
	return &authHandlers{
		cfg:    cfg,
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
// @Success 201 {object} models.UserWithToken
// @Failure 400 {string} string "Provided credentials don't match requirements"
// @Failure 409 {string} string "User with such username already exists"
// @Router /register [post]
func (h *authHandlers) Register() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) error {
		var req AuthRequest

		if err := utils.StructScan(r, &req); err != nil {
			return httper.NewBadRequestError(err)
		}

		err := utils.ValidateStruct(&req)
		if err != nil {
			return httper.ParseValidatorError(err)
		}

		user, err := h.authUC.Register(&auth.User{
			ID:        uuid.New(),
			Username:  req.Username,
			Password:  req.Password,
			CreatedAt: time.Now(),
		})

		if err != nil {
			return err
		}

		utils.RespondWithJSON(w, http.StatusCreated, user)
		return nil
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
// @Success 200 {object} models.UserWithToken
// @Failure 400 {string} string "User with such username is not found"
// @Failure 400 {string} string "Provided password is wrong"
// @Router /login [post]
func (h *authHandlers) Login() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) error {
		var req AuthRequest

		if err := utils.StructScan(r, &req); err != nil {
			return httper.NewBadRequestError(err)
		}

		user, err := h.authUC.Login(&auth.User{
			Username: req.Username,
			Password: req.Password,
		})
		if err != nil {
			return err
		}

		utils.RespondWithJSON(w, http.StatusOK, user)
		return nil
	}

	return httper.HandlerWithError(fn)
}
