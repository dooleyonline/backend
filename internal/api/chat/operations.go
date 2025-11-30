package chathandler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/dooleyonline/backend/internal/api/shared"
	chatsvc "github.com/dooleyonline/backend/internal/service/chat"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// GetLatest godoc
//
//	@Summary	Get latest messages
//	@Tags		chat
//	@Produce	json
//	@Param		roomID	path	string	true	"Room ID"
//	@Success	200		{array}	model.ChatMessage
//	@Router		/chat/{roomID}/messages [get]
func (h *Handler) GetMessages(c echo.Context) error {
	var (
		req    = c.Request()
		ctx    = req.Context()
		userID = c.(shared.Context).UserID
	)

	var (
		roomID string
		page   int32
	)

	if err := echo.PathParamsBinder(c).
		String("roomID", &roomID).BindError(); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	if err := echo.QueryParamsBinder(c).
		Int32("page", &page).BindError(); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if page < 1 {
		page = 1
	}

	isParticipant, err := h.svc.IsParticipant(ctx, roomID, userID)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if !isParticipant {
		return echo.ErrForbidden.WithInternal(errors.New("user is not a participant"))
	}

	res, err := h.svc.GetLatestMessages(ctx, roomID, page)
	if err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.JSON(http.StatusOK, res)
}

// EditMessage godoc
//
//	@Summary	Edit an existing message
//	@Tags		chat
//	@Param		messageID	path	int64	true	"Message ID"
//	@Param		body		body	string	true	"New message body"
//	@Success	204
//	@Router		/chat/messages/{messageID} [patch]
func (h *Handler) EditMessage(c echo.Context) error {
	var (
		req    = c.Request()
		ctx    = req.Context()
		userID = c.(shared.Context).UserID
	)

	var params chatsvc.EditMessageParams
	params.UserID = userID
	if err := echo.PathParamsBinder(c).Int64("messageID", &params.MessageID).BindError(); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if err := c.Bind(&params.Body); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if err := h.svc.EditMessage(ctx, &params); err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// DeleteMessage godoc
//
//	@Summary	Delete a message
//	@Tags		chat
//	@Param		messageID	path	int64	true	"Message ID"
//	@Success	204
//	@Router		/chat/messages/{messageID} [delete]
func (h *Handler) DeleteMessage(c echo.Context) error {
	var (
		req    = c.Request()
		ctx    = req.Context()
		userID = c.(shared.Context).UserID
	)

	var msgID int64
	if err := echo.PathParamsBinder(c).Int64("messageID", &msgID).BindError(); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	msg, err := h.svc.GetMessageByID(ctx, msgID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.ErrNotFound
		}
		return echo.ErrInternalServerError.WithInternal(err)
	}
	if msg.SentBy != userID {
		return echo.ErrForbidden.WithInternal(errors.New("user is not the sender of this message"))
	}

	if err := h.svc.DeleteMessage(ctx, msgID); err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// CreateRoom godoc
//
//	@Summary	Create a new room
//	@Tags		chat
//	@Produce	json
//	@Param		userIDs	body		[]string	true	"User IDs"
//	@Success	201		{string}	string
//	@Router		/chat/rooms [post]
func (h *Handler) CreateRoom(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)

	var userIDs []string
	if err := c.Bind(&userIDs); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if len(userIDs) < 2 {
		return echo.NewHTTPError(http.StatusBadRequest, "at least two user IDs are required to create a chat room")
	}

	roomID, err := h.svc.CreateRoom(ctx, userIDs)
	if err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.JSON(http.StatusCreated, roomID)
}

// DeleteRoom godoc
//
//	@Summary	Delete a chat room
//	@Tags		chat
//	@Param		roomID	path	string	true	"Room ID"
//	@Success	204
//	@Router		/chat/{roomID} [delete]
func (h *Handler) DeleteRoom(c echo.Context) error {
	var (
		req    = c.Request()
		ctx    = req.Context()
		userID = c.(shared.Context).UserID
	)

	var roomID string
	if err := echo.PathParamsBinder(c).String("roomID", &roomID).BindError(); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	isParticipant, err := h.svc.IsParticipant(ctx, roomID, userID)
	if err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}
	if !isParticipant {
		return echo.ErrForbidden.WithInternal(errors.New("user is not a participant of this room"))
	}

	if err := h.svc.DeleteRoom(ctx, roomID); err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// RemoveParticipant godoc
//
//	@Summary	Remove a user from a chat room
//	@Tags		chat
//	@Param		roomID	path	string	true	"Room ID"
//	@Param		userID	path	string	true	"User ID"
//	@Success	204
//	@Router		/chat/{roomID}/participants/{userID} [delete]
func (h *Handler) RemoveParticipant(c echo.Context) error {
	var (
		req    = c.Request()
		ctx    = req.Context()
		userID = c.(shared.Context).UserID
	)

	var roomID, participantID string
	if err := echo.PathParamsBinder(c).
		String("roomID", &roomID).
		String("userID", &participantID).
		BindError(); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	isParticipant, err := h.svc.IsParticipant(ctx, roomID, userID)
	if err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	if !isParticipant {
		return echo.ErrForbidden.WithInternal(errors.New("user is not a participant and does not have authority to remove another user"))
	}

	if err := h.svc.RemoveParticipant(ctx, &chatsvc.ParticipantParams{
		RoomID: roomID,
		UserID: participantID,
	}); err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// GetRooms godoc
//
//	@Summary	Get chat rooms for a user
//	@Tags		chat
//	@Produce	json
//	@Success	200	{array}	chatsvc.GetRoomsResult
//	@Router		/chat/rooms [get]
func (h *Handler) GetRooms(c echo.Context) error {
	var (
		req    = c.Request()
		ctx    = req.Context()
		userID = c.(shared.Context).UserID
		res    = c.Response()
	)

	res.Header().Set("Content-Type", "text/event-stream")
	res.Header().Set("Cache-Control", "no-cache")
	res.Header().Set("Connection", "keep-alive")

	flusher, ok := res.Writer.(http.Flusher)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "streaming unsupported")
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	lastPayload := []byte{}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			rooms, err := h.svc.GetRooms(ctx, userID)
			if err != nil {
				return echo.ErrInternalServerError.WithInternal(err)
			}

			payload, err := json.Marshal(rooms)
			if err != nil {
				c.Logger().Errorf("marshal rooms failed: %v", err)
				continue
			}

			// Send SSE event
			if !bytes.Equal(payload, lastPayload) {
				_, _ = res.Write([]byte("event: rooms\n"))
				_, _ = res.Write([]byte("data: "))
				_, _ = res.Write(payload)
				_, _ = res.Write([]byte("\n\n"))

				flusher.Flush()
				lastPayload = payload
			}
		}
	}
}

// GetParticipants godoc
//
//	@Summary	Get participants in a chat room
//	@Tags		chat
//	@Produce	json
//	@Param		roomID	path	string	true	"Room ID"
//	@Success	200		{array}	model.ChatParticipant
//	@Router		/chat/{roomID}/participants [get]
func (h *Handler) GetParticipants(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)

	var roomID string
	if err := echo.PathParamsBinder(c).String("roomID", &roomID).BindError(); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	participants, err := h.svc.GetParticipants(ctx, roomID)
	if err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.JSON(http.StatusOK, participants)
}
