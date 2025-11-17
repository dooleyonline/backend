package chathandler

import (
	"errors"
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
//	@Router		/chat/{roomID}/ws [get]
func (h *Handler) HandleConnections(c echo.Context) error {
	var (
		req    = c.Request()
		ctx    = req.Context()
		userID = c.(shared.Context).UserID
	)

	var roomID string
	if err := echo.PathParamsBinder(c).String("roomID", &roomID).BindError(); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	if roomID == "" {
		return echo.ErrBadRequest
	}

	isParticipant, err := h.svc.IsParticipant(ctx, roomID, userID)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if !isParticipant {
		return echo.ErrForbidden.WithInternal(errors.New("user is not a participant"))
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	defer conn.Close()

	client := &Client{
		conn:   conn,
		roomID: roomID,
		userID: userID,
	}

	h.hub.register <- client
	defer func() {
		h.hub.unregister <- client
	}()

	var lastReadMsgID *int64
	defer func() {
		if lastReadMsgID == nil {
			return
		}
		if err := h.svc.UpdateLastReadMessageID(ctx, roomID, userID, *lastReadMsgID); err != nil {
			echo.ErrInternalServerError.WithInternal(err)
		}
	}()

	for {
		_, m, err := conn.ReadMessage()
		if err != nil {
			break
		}

		msg := string(m)

		params := &chatsvc.CreateMessageParams{
			UserID:  userID,
			RoomID:  roomID,
			Message: msg,
		}

		created, err := h.svc.CreateMessage(ctx, params)
		if err != nil {
			_ = conn.WriteJSON(map[string]string{"error": "failed to create message"})
			continue
		}

		id := created.ID
		lastReadMsgID = &id

		msgCopy := created
		h.hub.broadcast <- &msgCopy
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
