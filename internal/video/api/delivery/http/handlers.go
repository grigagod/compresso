package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/grigagod/compresso/internal/httper"
	"github.com/grigagod/compresso/internal/models"
	"github.com/grigagod/compresso/internal/utils"
	"github.com/grigagod/compresso/internal/video/api"
	"github.com/grigagod/compresso/pkg/converter"
)

type apiHandlers struct {
	apiUC api.UseCase
}

func NewAPIHandlers(apiUC api.UseCase) api.Handlers {
	return &apiHandlers{
		apiUC: apiUC,
	}
}

// Register godoc
// @Summary Create new video
// @Description Authorized users can upload their videos
// @Tags Video
// @Accept video/webm
// @Produce json
// @Success 201 {object} models.Video
// @Failure 400 {string} msg "Bad request msg"
// @Failure 401 {string} msg "Wrong credentials"
// @Failure 415 {string} msg "Provided media type is not allowed"
// @Security ApiKeyAuth
// @Router /videos [post].
func (h *apiHandlers) CreateVideo() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) error {
		userID, ok := utils.UserIDFromContext(r.Context())
		if !ok {
			return httper.NewWrongCredentialsMsg()
		}

		contentType, ok := utils.ContentTypeFromContext(r.Context())
		if !ok {
			return httper.NewNotAllowedContentType()
		}

		format, ok := utils.DetectVideoFormatFromHeader(contentType)
		if !ok {
			return httper.NewNotAllowedContentType()
		}

		video := models.Video{
			ID:       uuid.New(),
			AuthorID: userID,
			Format:   format,
		}

		v, err := h.apiUC.CreateVideo(r.Context(), &video, r.Body)
		if err != nil {
			return err
		}

		return utils.RespondWithJSON(w, http.StatusCreated, v)
	}

	return httper.HandlerWithError(fn)
}

// Register godoc
// @Summary Create new video ticket
// @Description Authorized users can create tickets for processing uploaded videos
// @Tags Video
// @Accept json
// @Produce json
// @Param req body CreateTicketRequest true "info for video processing"
// @Success 201 {object} models.VideoTicket
// @Failure 400 {string} msg "Bad request msg"
// @Failure 401 {string} msg "Wrong credentials"
// @Failure 415 {string} msg "Provided media type is not allowed"
// @Security ApiKeyAuth
// @Router /tickets [post].
func (h *apiHandlers) CreateTicket() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) error {
		userID, ok := utils.UserIDFromContext(r.Context())
		if !ok {
			return httper.NewWrongCredentialsMsg()
		}

		var req CreateTicketRequest

		if err := utils.StructScan(r.Body, &req); err != nil {
			return httper.NewBadRequestError(err)
		}

		if err := utils.ValidateStruct(&req); err != nil {
			return httper.ParseValidatorError(err)
		}

		if err := converter.ValidateCRF(req.CRF); err != nil {
			return httper.NewBadRequestError(err)
		}

		target_format, ok := utils.DetectVideoFormat(req.TargetFormat)
		if !ok {
			return httper.NewNotAllowedContentType()
		}

		ticket := models.VideoTicket{
			Ticket: models.Ticket{
				ID:       uuid.New(),
				AuthorID: userID,
			},
			VideoID:      req.VideoID,
			CRF:          req.CRF,
			TargetFormat: target_format,
		}

		t, err := h.apiUC.CreateTicket(r.Context(), &ticket)
		if err != nil {
			return err
		}

		return utils.RespondWithJSON(w, http.StatusCreated, t)
	}

	return httper.HandlerWithError(fn)
}

// Register godoc
// @Summary Get video by ID
// @Description Authorized users can get uploaded videos by ID.
// @Tags Video
// @Accept json
// @Produce json
// @Param id path string true "Video ID"
// @Success 200 {object} models.Video
// @Failure 400 {string} msg "Bad request msg"
// @Failure 401 {string} msg "Wrong credentials"
// @Failure 404 {string} msg "Not found"
// @Security ApiKeyAuth
// @Router /videos/{id} [get].
func (h *apiHandlers) GetVideoByID() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) error {
		userID, ok := utils.UserIDFromContext(r.Context())
		if !ok {
			return httper.NewWrongCredentialsMsg()
		}

		id, err := uuid.Parse(chi.URLParamFromCtx(r.Context(), "id"))
		if err != nil {
			return httper.NewBadRequestError(err)
		}

		video, err := h.apiUC.GetVideoByID(r.Context(), userID, id)
		if err != nil {
			return err
		}

		return utils.RespondWithJSON(w, http.StatusOK, video)
	}

	return httper.HandlerWithError(fn)
}

// Register godoc
// @Summary Get video by ID
// @Description Authorized users can get uploaded videos by ID.
// @Tags Video
// @Accept json
// @Produce json
// @Param id path string true "Video ID"
// @Success 200 {object} models.Video
// @Failure 400 {string} msg "Bad request msg"
// @Failure 401 {string} msg "Wrong credentials"
// @Failure 404 {string} msg "Not found"
// @Security ApiKeyAuth
// @Router /tickets/{id} [get].
func (h *apiHandlers) GetTicketByID() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) error {
		userID, ok := utils.UserIDFromContext(r.Context())
		if !ok {
			return httper.NewWrongCredentialsMsg()
		}

		id, err := uuid.Parse(chi.URLParamFromCtx(r.Context(), "id"))
		if err != nil {
			return httper.NewBadRequestError(err)
		}

		ticket, err := h.apiUC.GetTicketByID(r.Context(), userID, id)
		if err != nil {
			return err
		}

		return utils.RespondWithJSON(w, http.StatusOK, ticket)
	}

	return httper.HandlerWithError(fn)
}

// Register godoc
// @Summary Get videos
// @Description Authorized users can get all uploaded videos.
// @Tags Video
// @Accept json
// @Produce json
// @Success 200 {object} []models.Video
// @Failure 400 {string} msg "Bad request msg"
// @Failure 401 {string} msg "Wrong credentials"
// @Failure 404 {string} msg "Not found"
// @Security ApiKeyAuth
// @Router /videos [get].
func (h *apiHandlers) GetVideos() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) error {
		userID, ok := utils.UserIDFromContext(r.Context())
		if !ok {
			return httper.NewWrongCredentialsMsg()
		}

		videos, err := h.apiUC.GetVideos(r.Context(), userID)
		if err != nil {
			return err
		}

		return utils.RespondWithJSON(w, http.StatusOK, videos)
	}

	return httper.HandlerWithError(fn)
}

// Register godoc
// @Summary Get tickets
// @Description Authorized users can get all video tickets.
// @Tags Video
// @Accept json
// @Produce json
// @Success 200 {object} []models.VideoTicket
// @Failure 400 {string} msg "Bad request msg"
// @Failure 401 {string} msg "Wrong creadentials"
// @Failure 404 {string} msg "Not found"
// @Security ApiKeyAuth
// @Router /tickets [get].
func (h *apiHandlers) GetTickets() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) error {
		userID, ok := utils.UserIDFromContext(r.Context())
		if !ok {
			return httper.NewWrongCredentialsMsg()
		}

		tickets, err := h.apiUC.GetTickets(r.Context(), userID)
		if err != nil {
			return err
		}

		return utils.RespondWithJSON(w, http.StatusOK, tickets)
	}

	return httper.HandlerWithError(fn)
}
