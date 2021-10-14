package http

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/grigagod/compresso/internal/httper"
	"github.com/grigagod/compresso/internal/middleware"
	"github.com/grigagod/compresso/internal/models"
	"github.com/grigagod/compresso/internal/video"
	"github.com/grigagod/compresso/pkg/converter"
	"github.com/grigagod/compresso/pkg/utils"
)

type videoHandlers struct {
	videoUC video.UseCase
}

func NewVideoHandlers(videoUC video.UseCase) video.Handlers {
	return &videoHandlers{
		videoUC: videoUC,
	}
}

func (h *videoHandlers) UploadVideo() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) error {
		userID := r.Context().Value(middleware.UserIDCtxKey{}).(uuid.UUID)

		format, err := utils.DetectVideoFormatFromHeader(r.Header["Content-Type"][0])
		if err != nil {
			return httper.NewBadRequestMsg(httper.NotAllowedMediaTypeMsg)
		}

		video := models.Video{
			ID:       uuid.New(),
			AuthorID: userID,
			Format:   format,
		}
		video.URL = utils.GenerateURL(video.AuthorID, video.ID)

		v, err := h.videoUC.UploadVideo(r.Context(), &video, r.Body)
		if err != nil {
			return err
		}

		utils.RespondWithJSON(w, http.StatusCreated, &v)

		return nil
	}

	return httper.HandlerWithError(fn)
}

func (h *videoHandlers) CreateTicket() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) error {
		userId := r.Context().Value(middleware.UserIDCtxKey{}).(uuid.UUID)
		var req CreateTicketRequest

		if err := utils.StructScan(r, &req); err != nil {
			log.Fatal(err)
			return httper.NewBadRequestError(err)
		}

		err := utils.ValidateStruct(&req)
		if err != nil {
			log.Fatal(err)
			return httper.ParseValidatorError(err)
		}

		err = converter.ValidateCRF(req.CRF)
		if err != nil {
			log.Fatal(err)
			return httper.NewBadRequestError(err)
		}

		target_format, err := utils.DetectVideoFormat(req.TargetFormat)
		if err != nil {
			return httper.NewBadRequestMsg(httper.NotAllowedMediaTypeMsg)
		}

		ticket := models.VideoTicket{
			Ticket: models.Ticket{
				ID:       uuid.New(),
				AuthorID: userId,
			},
			VideoID:      req.VideoID,
			CRF:          req.CRF,
			TargetFormat: target_format,
		}

		t, err := h.videoUC.CreateTicket(r.Context(), &ticket)
		if err != nil {
			return err
		}

		utils.RespondWithJSON(w, http.StatusCreated, &t)

		return nil
	}

	return httper.HandlerWithError(fn)
}
