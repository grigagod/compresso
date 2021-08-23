package http

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/grigagod/compresso/internal/auth"
	"github.com/grigagod/compresso/internal/auth/config"
	"github.com/grigagod/compresso/internal/auth/models"
	"github.com/grigagod/compresso/pkg/utils"
)

type authHandlers struct {
	cfg    *config.Config
	authUC auth.UseCase
}

func NewAuthHandlers(cfg *config.Config, authUC auth.UseCase) *authHandlers {
	return &authHandlers{
		cfg:    cfg,
		authUC: authUC,
	}
}

func (h *authHandlers) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AuthRequest

		if err := utils.StructScan(r, &req); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		}

		user, err := h.authUC.Register(&models.User{
			ID:        uuid.New(),
			Username:  req.Username,
			Password:  req.Password,
			CreatedAt: time.Now(),
		})

		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		}

		utils.RespondWithJSON(w, http.StatusCreated, user)

	}
}

func (h *authHandlers) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AuthRequest

		if err := utils.StructScan(r, &req); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		}

		user, err := h.authUC.Login(&models.User{
			Username: req.Username,
			Password: req.Password,
		})
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		}

		utils.RespondWithJSON(w, http.StatusOK, user)

	}
}
