package http

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/grigagod/compresso/internal/auth"
	"github.com/grigagod/compresso/internal/auth/config"
	"github.com/grigagod/compresso/internal/models"
	"github.com/grigagod/compresso/pkg/utils"
)

type authHandlers struct {
	cfg    *config.Config
	authUC auth.UseCase
}

func NewAuthHandlers(cfg *config.Config, authUC auth.UseCase) auth.Handlers {
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
// @Router /register [post]
func (h *authHandlers) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AuthRequest

		if err := utils.StructScan(r, &req); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		user, err := h.authUC.Register(&models.User{
			ID:        uuid.New(),
			Username:  req.Username,
			Password:  req.Password,
			CreatedAt: time.Now(),
		})

		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.RespondWithJSON(w, http.StatusCreated, user)
	}
}

// Login godoc
// @Summary Login new user
// @Description login user, returns user and token
// @Tags Auth
// @Accept json
// @Produce json
// @Param creds body AuthRequest true "user credentials"
// @Success 200 {object} models.UserWithToken
// @Router /login [post]
func (h *authHandlers) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AuthRequest

		if err := utils.StructScan(r, &req); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		user, err := h.authUC.Login(&models.User{
			Username: req.Username,
			Password: req.Password,
		})
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, user)

	}
}
