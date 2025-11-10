package storagehandler

import (
	"net/http"

	storagesvc "github.com/dooleyonline/backend/internal/service/storage"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	svc *storagesvc.Service
}

func New(svc *storagesvc.Service) *Handler {
	return &Handler{svc}
}

// PresignUpload godoc
//
//	@Summary	Generate presigned URL for item upload
//	@Tags		util
//	@Accept		json
//	@Produce	json
//	@Param		type	query		string	true	"Content type of the item to be uploaded"
//	@Param		bucket	query		string	true	"Storage bucket name"
//	@Success	200		{object}	storagesvc.PresignResult
//	@Router		/storage/presign [post]
func (h *Handler) PresignUpload(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)

	var params storagesvc.PresignParams
	if err := echo.QueryParamsBinder(c).
		String("type", &params.ContentType).
		String("bucket", &params.Bucket).
		BindError(); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if params.ContentType == "" || params.Bucket == "" {
		return echo.ErrBadRequest
	}

	res, err := h.svc.PresignUpload(ctx, params)
	if err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.JSON(http.StatusOK, res)
}
