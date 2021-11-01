package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/grigagod/compresso/internal/httper"
	"github.com/grigagod/compresso/internal/middleware"
	"github.com/grigagod/compresso/internal/models"
	"github.com/grigagod/compresso/internal/utils"
	"github.com/grigagod/compresso/internal/video"
	"github.com/grigagod/compresso/pkg/converter"
)

type videoHandlers struct {
	videoUC video.UseCase
}

func NewVideoHandlers(videoUC video.UseCase) video.Handlers {
	return &videoHandlers{
		videoUC: videoUC,
	}
}

// Register godoc
// @Summary Create new video
// @Description Authorized users can upload their videos
// @Tags Video
// @Accept video/webm
// @Produce json
// @Success 201 {object} models.Video
// @Failure 401 {string} msg "Wrong creadentials"
// @Failure 409 {string} msg "Provided header is not allowed"
// @Failure 409 {string} msg "Provided media type is not allowed"
// @Security ApiKeyAuth
// @Router /videos [post]
func (h *videoHandlers) CreateVideo() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) error {
		userID, err := uuid.Parse(r.Context().Value(middleware.UserIDCtxKey{}).(string))
		if err != nil {
			return httper.NewStatusMsg(http.StatusUnauthorized, httper.WrongCredentialsMsg)
		}

		contentType, ok := r.Context().Value(middleware.ContentTypeCtxKey{}).(string)
		if !ok {
			return httper.NewBadRequestMsg(httper.NotAllowedHeaderMsg)
		}

		format, err := utils.DetectVideoFormatFromHeader(contentType)
		if err != nil {
			return httper.NewBadRequestMsg(httper.NotAllowedMediaTypeMsg)
		}

		video := models.Video{
			ID:       uuid.New(),
			AuthorID: userID,
			Format:   format,
		}

		v, err := h.videoUC.CreateVideo(r.Context(), &video, r.Body)
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
// @Failure 401 {string} msg "Wrong creadentials"
// @Failure 409 {string} msg "Provided header is not allowed"
// @Failure 409 {string} msg "Provided media type is not allowed"
// @Security ApiKeyAuth
// @Router /tickets [post]
func (h *videoHandlers) CreateTicket() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) error {
		userID, err := uuid.Parse(r.Context().Value(middleware.UserIDCtxKey{}).(string))
		if err != nil {
			return httper.NewStatusMsg(http.StatusUnauthorized, httper.WrongCredentialsMsg)
		}

		var req CreateTicketRequest

		if err := utils.StructScan(r, &req); err != nil {
			return httper.NewBadRequestError(err)
		}

		err = utils.ValidateStruct(&req)
		if err != nil {
			return httper.ParseValidatorError(err)
		}

		err = converter.ValidateCRF(req.CRF)
		if err != nil {
			return httper.NewBadRequestError(err)
		}

		target_format, err := utils.DetectVideoFormat(req.TargetFormat)
		if err != nil {
			return httper.NewBadRequestMsg(httper.NotAllowedMediaTypeMsg)
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

		t, err := h.videoUC.CreateTicket(r.Context(), &ticket)
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
// @Success 201 {object} models.Video
// @Failure 401 {string} msg "Wrong creadentials"
// @Failure 409 {string} msg "Provided header is not allowed"
// @Failure 409 {string} msg "Provided media type is not allowed"
// @Failure 409 {string} msg "Bad request"
// @Security ApiKeyAuth
// @Router /videos/{id} [get]
func (h *videoHandlers) GetVideoByID() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) error {
		userID, err := uuid.Parse(r.Context().Value(middleware.UserIDCtxKey{}).(string))
		if err != nil {
			return httper.NewStatusMsg(http.StatusUnauthorized, httper.WrongCredentialsMsg)
		}

		id, err := uuid.Parse(chi.URLParamFromCtx(r.Context(), "id"))
		if err != nil {
			return httper.NewBadRequestMsg(httper.BadRequestMsg)
		}

		video, err := h.videoUC.GetVideoByID(r.Context(), userID, id)
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
// @Success 201 {object} models.Video
// @Failure 401 {string} msg "Wrong creadentials"
// @Failure 409 {string} msg "Provided header is not allowed"
// @Failure 409 {string} msg "Provided media type is not allowed"
// @Failure 409 {string} msg "Bad request"
// @Security ApiKeyAuth
// @Router /tickets/{id} [get]
func (h *videoHandlers) GetTicketByID() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) error {
		userID, err := uuid.Parse(r.Context().Value(middleware.UserIDCtxKey{}).(string))
		if err != nil {
			return httper.NewStatusMsg(http.StatusUnauthorized, httper.WrongCredentialsMsg)
		}

		id, err := uuid.Parse(chi.URLParamFromCtx(r.Context(), "id"))
		if err != nil {
			return httper.NewBadRequestMsg(httper.BadRequestMsg)
		}

		ticket, err := h.videoUC.GetTicketByID(r.Context(), userID, id)
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
// @Success 201 {object} []models.Video
// @Failure 401 {string} msg "Wrong creadentials"
// @Failure 409 {string} msg "Provided header is not allowed"
// @Failure 409 {string} msg "Provided media type is not allowed"
// @Failure 409 {string} msg "Bad request"
// @Security ApiKeyAuth
// @Router /videos [get]
func (h *videoHandlers) GetVideos() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) error {
		userID, err := uuid.Parse(r.Context().Value(middleware.UserIDCtxKey{}).(string))
		if err != nil {
			return httper.NewStatusMsg(http.StatusUnauthorized, httper.WrongCredentialsMsg)
		}

		videos, err := h.videoUC.GetVideos(r.Context(), userID)
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
// @Success 201 {object} []models.VideoTicket
// @Failure 401 {string} msg "Wrong creadentials"
// @Failure 409 {string} msg "Provided header is not allowed"
// @Failure 409 {string} msg "Provided media type is not allowed"
// @Failure 409 {string} msg "Bad request"
// @Security ApiKeyAuth
// @Router /tickets [get]
func (h *videoHandlers) GetTickets() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) error {
		userID, err := uuid.Parse(r.Context().Value(middleware.UserIDCtxKey{}).(string))
		if err != nil {
			return httper.NewStatusMsg(http.StatusUnauthorized, httper.WrongCredentialsMsg)
		}

		tickets, err := h.videoUC.GetTickets(r.Context(), userID)
		if err != nil {
			return err
		}

		return utils.RespondWithJSON(w, http.StatusOK, tickets)
	}

	return httper.HandlerWithError(fn)
}
