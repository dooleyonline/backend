package chathandler

import (
	"net/http"

	"github.com/dooleyonline/backend/internal/api/shared"
	chatmessage "github.com/dooleyonline/backend/internal/db/chat/message"
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
	return &Handler{svc, hub}
}

type Client struct {
	conn   *websocket.Conn
	send chan *Message
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
//	@Param		roomID		path string	false	"Room ID"
//	@Success	200			{array}	model.ChatMessage
//	@Router		/chat/:roomID/messages [get]
func (h *Handler) GetLatest(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)

	var roomID string
	if err := echo.PathParamsBinder(c).String(&roomID).BindError(); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	res, err := h.svc.GetLatestMessages(ctx, roomID)
	if err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.JSON(http.StatusOK, res)
}

// CreateRoom godoc
//
//	@Summary	Create a new room
//	@Tags		chat
//	@Produce	json
//	@Param		userIDs	body []string	true	"User IDs"
//	@Success	201			{object}	model.ChatRoom
//	@Router		/chat [post]
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
		return echo.ErrBadRequest
	}

	res, err := h.svc.CreateRoom(ctx, userIDs)
	if err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.JSON(http.StatusCreated, res)
}

var (
	upgrader = websocket.Upgrader{}
)

func (h *Handler) handleConnections(c echo.Context) error {
	var roomID string
	if err := echo.PathParamsBinder(c).String("id", &roomID).BindError(); err != nil {
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

	h.hub.register <- &Client{
		conn:   conn,
		userID: c.(shared.Context).UserID,
	}
	defer func() {
		h.hub.unregister <- &Client{
			conn:
		}
	}

	for {
		_, m, err := conn.ReadMessage()
		if err != nil {
			break
		}
		userID := c.(shared.Context).UserID
		msg := string(m)

		if err := h.svc.CreateMessage(ctx, userID, roomID, msg); err != nil {
			continue
		}

		h.hub.broadcast <- &model.ChatMessage{
			SentBy: userID,
			RoomID: roomID,
			Body:   msg,
		}
	}

	return c.JSON(http.StatusCreated, conn)
}

func (h *Handler) run() error {
	for {
		select {
		case c := <-h.hub.register:
			if h.hub.rooms[c.roomID] == nil {
				h.hub.rooms[c.roomID] = make(map[*Client]bool)
			}
		 	h.rooms[c.roomID][c] = true
		}
		case := <- h.hub.unregister:
			if m, ok := h.hub.rooms[c.roomID]; ok{
				if _, ok := m[c]; ok {
					delete(m, c)
					close(c.send)
				}
				if  
			}
	}
}
