package verify

import (
	verifysvc "github.com/dooleyonline/backend/internal/service/verify"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	svc *verifysvc.Service
}

func New(svc *verifysvc.Service) *Handler {
	return &Handler{svc}
}

// SendVerification godoc
//
//	@Summary		Send verification email
//	@Description	Sends (or resends) a verification link to the provided email. Returns 204 on success.
//	@Tags			verify
//	@Accept			json
//	@Produce		json
//	@Param			request	body	verifysvc.SendParams	true	"Email payload"
//	@Success		204
//	@Failure		400	{object}	map[string]string	"invalid request body"
//	@Failure		500	{object}	map[string]string	"internal error"
//	@Router			/auth/verification [post]
func (h *Handler) SendVerification(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)

	var params verifysvc.SendParams
	if err := c.Bind(&params); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	err := h.svc.Send(ctx, params)

	if err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return nil
}

// VerifyUser godoc
//
//	@Summary		Verify user by token
//	@Description	Consumes the verification token from the path and marks the user verified.
//	@Tags			verify
//	@Produce		json
//	@Param			id	path		string				true	"Verification token (UUID)"
//	@Success		200	{object}	map[string]string	"user verified"
//	@Failure		400	{object}	map[string]string	"invalid token format"
//	@Failure		401	{object}	map[string]string	"invalid or expired token"
//	@Failure		500	{object}	map[string]string	"internal error"
//	@Router			/auth/verification/{id} [post]
func (h *Handler) VerifyUser(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)

	var id string
	if err := echo.PathParamsBinder(c).String("id", &id).BindError(); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	err := h.svc.VerifyUserEmail(ctx, id)
	if err != nil {
		return echo.ErrUnauthorized.WithInternal(err)
	}

	return nil
}
