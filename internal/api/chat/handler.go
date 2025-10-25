package chathandler

import (
	"net/http"

	"github.com/dooleyonline/backend/internal/api/shared"
	"github.com/dooleyonline/backend/internal/model"
	chatsvc "github.com/dooleyonline/backend/internal/service/chat"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	svc *chatsvc.Service
	hub *Hub
}

func New(svc *chatsvc.Service) *Handler {
	hub := &Hub{
		rooms:      make(map[string]map[*Client]bool),
		broadcast:  make(chan *model.ChatMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	go hub.run()
	return &Handler{svc, hub}
}

type Client struct {
	conn   *websocket.Conn
	roomID string
	userID string
}

type Hub struct {
	rooms      map[string]map[*Client]bool
	broadcast  chan *model.ChatMessage
	register   chan *Client
	unregister chan *Client
}

// GetLatest godoc
//
//	@Summary	Get latest messages
//	@Tags		chat
//	@Produce	json
//	@Param		roomID	path	string	true	"Room ID"
//	@Success	200		{array}	model.ChatMessage
//	@Router		/chat/rooms/{roomID}/messages [get]
func (h *Handler) GetLatest(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)

	var roomID string
	if err := echo.PathParamsBinder(c).String("roomID", &roomID).BindError(); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	res, err := h.svc.GetLatestMessages(ctx, roomID)
	if err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.JSON(http.StatusOK, res)
}

// CreateMessage godoc
//
//	@Summary	Post a new message
//	@Tags		chat
//	@Param		roomID	path	string	true	"Room ID"
//	@Param		body	body	string	true	"Message body"
//	@Success	201
//	@Router		/chat/rooms/{roomID}/messages [post]
func (h *Handler) CreateMessage(c echo.Context) error {
	var (
		req    = c.Request()
		ctx    = req.Context()
		userID = c.(shared.Context).UserID
	)

	params := &chatsvc.CreateMessageParams{}

	if err := echo.PathParamsBinder(c).
		String("roomID", &params.RoomID).
		BindError(); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	params.UserID = userID

	if err := c.Bind(&params.Message); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if err := h.svc.CreateMessage(ctx, params); err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.NoContent(http.StatusCreated)
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
		req = c.Request()
		ctx = req.Context()
	)

	var msgID int64
	if err := echo.PathParamsBinder(c).Int64("messageID", &msgID).BindError(); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	var body string
	if err := c.Bind(&body); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if err := h.svc.EditMessage(ctx, &chatsvc.EditMessageParams{
		ID:   msgID,
		Body: body,
	}); err != nil {
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
		req = c.Request()
		ctx = req.Context()
	)

	var msgID int64
	if err := echo.PathParamsBinder(c).Int64("messageID", &msgID).BindError(); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
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
//	@Success	201		{object}	string
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

	return c.JSON(http.StatusCreated, map[string]string{"room_id": roomID})
}

// DeleteRoom godoc
//
//	@Summary	Delete a chat room
//	@Tags		chat
//	@Param		roomID	path	string	true	"Room ID"
//	@Success	204
//	@Router		/chat/rooms/{roomID} [delete]
func (h *Handler) DeleteRoom(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)

	var roomID string
	if err := echo.PathParamsBinder(c).String("roomID", &roomID).BindError(); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if err := h.svc.DeleteRoom(ctx, roomID); err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// AddParticipant godoc
//
//	@Summary	Add a user to a chat room
//	@Tags		chat
//	@Param		roomID	path	string	true	"Room ID"
//	@Param		userID	body	string	true	"User ID"
//	@Success	204
//	@Router		/chat/rooms/{roomID}/participants [post]
func (h *Handler) AddParticipant(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)

	var roomID string
	if err := echo.PathParamsBinder(c).String("roomID", &roomID).BindError(); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	var userID string
	if err := c.Bind(&userID); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	params := &chatsvc.ParticipantParams{
		RoomID: roomID,
		UserID: userID,
	}
	if err := h.svc.AddParticipant(ctx, params); err != nil {
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
//	@Router		/chat/rooms/{roomID}/participants/{userID} [delete]
func (h *Handler) RemoveParticipant(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)

	var roomID, userID string
	if err := echo.PathParamsBinder(c).
		String("roomID", &roomID).
		String("userID", &userID).
		BindError(); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if err := h.svc.RemoveParticipant(ctx, &chatsvc.ParticipantParams{
		RoomID: roomID,
		UserID: userID,
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
//	@Success	200	{array}	model.ChatParticipant
//	@Router		/chat/rooms [get]
func (h *Handler) GetRooms(c echo.Context) error {
	var (
		req    = c.Request()
		ctx    = req.Context()
		userID = c.(shared.Context).UserID
	)

	rooms, err := h.svc.GetRooms(ctx, userID)
	if err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.JSON(http.StatusOK, rooms)
}

// GetParticipants godoc
//
//	@Summary	Get participants in a chat room
//	@Tags		chat
//	@Produce	json
//	@Param		roomID	path	string	true	"Room ID"
//	@Success	200		{array}	model.ChatParticipant
//	@Router		/chat/rooms/{roomID}/participants [get]
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

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// HandleConnections godoc
//
//	@Summary	Handle WebSocket connections
//	@Tags		chat
//	@Param		roomID	path	string	true	"Room ID"
//	@Success	200
//	@Router		/chat/rooms/{roomID}/ws [get]
func (h *Handler) HandleConnections(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)

	var roomID string
	if err := echo.PathParamsBinder(c).String("roomID", &roomID).BindError(); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	if roomID == "" {
		return echo.ErrBadRequest
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	defer conn.Close()

	client := &Client{
		conn:   conn,
		roomID: roomID,
		userID: c.(shared.Context).UserID,
	}

	h.hub.register <- client
	defer func() {
		h.hub.unregister <- client
	}()

	for {
		_, m, err := conn.ReadMessage()
		if err != nil {
			break
		}
		userID := c.(shared.Context).UserID
		msg := string(m)

		params := &chatsvc.CreateMessageParams{
			UserID:  userID,
			RoomID:  roomID,
			Message: msg,
		}
		if err := h.svc.CreateMessage(ctx, params); err != nil {
			continue
		}

		h.hub.broadcast <- &model.ChatMessage{
			SentBy: userID,
			RoomID: roomID,
			Body:   msg,
		}
	}

	return nil
}

func (h *Hub) run() {
	for {
		select {
		case c := <-h.register:
			if h.rooms[c.roomID] == nil {
				h.rooms[c.roomID] = make(map[*Client]bool)
			}
			h.rooms[c.roomID][c] = true
		case c := <-h.unregister:
			roomID := c.roomID
			if m, ok := h.rooms[roomID]; ok {
				if _, ok := m[c]; ok {
					delete(m, c)
					c.conn.Close()
				}
				if len(m) == 0 {
					delete(h.rooms, roomID)
				}
			}
		case message := <-h.broadcast:
			roomID := message.RoomID
			if clients, ok := h.rooms[roomID]; ok {
				for c := range clients {
					err := c.conn.WriteJSON(message)
					if err != nil {
						c.conn.Close()
						delete(clients, c)
					}
				}
			}
		}
	}
}
