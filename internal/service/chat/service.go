package chatsvc

import (
	"context"
	"errors"
	"slices"

	"github.com/bwmarrin/snowflake"
	"github.com/dooleyonline/backend/internal/config"
	"github.com/dooleyonline/backend/internal/db"
	chatmessage "github.com/dooleyonline/backend/internal/db/chat/message"
	chatparticipant "github.com/dooleyonline/backend/internal/db/chat/participant"
	chatroom "github.com/dooleyonline/backend/internal/db/chat/room"
	"github.com/dooleyonline/backend/internal/model"
	"github.com/jackc/pgx/v5"
)

type Service struct {
	cfg *config.Config
	db  *db.DB

	node *snowflake.Node
}

func New(cfg *config.Config, db *db.DB) *Service {
	node, _ := snowflake.NewNode(1)
	return &Service{cfg, db, node}
}

func (s *Service) GenerateID() int64 {
	id := s.node.Generate()
	return id.Int64()
}

type CreateMessageParams struct {
	UserID  string
	RoomID  string
	Message string
}

func (s *Service) CreateMessage(ctx context.Context, p *CreateMessageParams) error {
	messageID := s.GenerateID()
	if err := s.db.Chat.Message.Create(ctx, chatmessage.CreateParams{
		ID:     messageID,
		RoomID: p.RoomID,
		SentBy: p.UserID,
		Body:   p.Message,
	}); err != nil {
		return err
	}

	return nil
}

type EditMessageParams struct {
	UserID    string
	MessageID int64
	Body      string
}

func (s *Service) EditMessage(ctx context.Context, p *EditMessageParams) error {
	msg, err := s.GetMessageByID(ctx, p.MessageID)
	if err != nil {
		return err
	}
	if msg.SentBy != p.UserID {
		return errors.New("user is not the sender of this message")
	}

	if err := s.db.Chat.Message.EditMessage(ctx, chatmessage.EditMessageParams{
		ID:   p.MessageID,
		Body: p.Body,
	}); err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteMessage(ctx context.Context, messageID int64) error {
	return s.db.Chat.Message.Delete(ctx, messageID)
}

func (s *Service) GetMessageByID(ctx context.Context, messageID int64) (model.ChatMessage, error) {
	return s.db.Chat.Message.GetByID(ctx, messageID)
}

func (s *Service) CreateRoom(ctx context.Context, users []string) (string, error) {
	tx, err := s.db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	room := chatroom.New(tx)
	participant := chatparticipant.New(tx)

	id, err := room.CreateRoom(ctx, users)
	if err != nil {
		return "", err
	}

	for _, user := range users {
		if err := participant.Create(ctx, chatparticipant.CreateParams{
			RoomID: id,
			UserID: user,
		}); err != nil {
			return "", err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return "", err
	}
	return id, nil
}

func (s *Service) DeleteRoom(ctx context.Context, roomID string) error {
	return s.db.Chat.Room.DeleteRoom(ctx, roomID)
}

func (s *Service) GetLatestMessages(ctx context.Context, roomID string, page int32) ([]model.ChatMessage, error) {
	return s.db.Chat.Message.GetMany(ctx, chatmessage.GetManyParams{
		RoomID: roomID,
		Page:   page,
	})
}

type ParticipantParams struct {
	RoomID string
	UserID string
}

func (s *Service) AddParticipant(ctx context.Context, p *ParticipantParams) error {
	tx, err := s.db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	room := chatroom.New(tx)
	participant := chatparticipant.New(tx)

	if err := room.AddParticipant(ctx, chatroom.AddParticipantParams{
		UserID: p.UserID,
		RoomID: p.RoomID,
	}); err != nil {
		return err
	}

	if err := participant.Create(ctx, chatparticipant.CreateParams{
		RoomID: p.RoomID,
		UserID: p.UserID,
	}); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *Service) RemoveParticipant(ctx context.Context, p *ParticipantParams) error {
	tx, err := s.db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	room := chatroom.New(tx)
	participant := chatparticipant.New(tx)

	if err := room.RemoveParticipant(ctx, chatroom.RemoveParticipantParams{
		UserID: p.UserID,
		RoomID: p.RoomID,
	}); err != nil {
		return err
	}

	if err := participant.Delete(ctx, chatparticipant.DeleteParams{
		RoomID: p.RoomID,
		UserID: p.UserID,
	}); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// type GetRoomsResult struct {
// 	RoomID            string `json:"room_id"`
// 	UserID            string `json:"user_id"`
// 	LastReadMessageID *int64 `json:"last_read_message_id"`
// 	ReadAll           bool   `json:"read_all"`
// }

type GetRoomsResult struct {
	RoomID      string             `json:"room_id"`
	Title       string             `json:"title"`
	LastMessage *model.ChatMessage `json:"last_message"`
	ReadAll     bool               `json:"read_all"`
}

func (s *Service) GetRooms(ctx context.Context, userID string) ([]GetRoomsResult, error) {
	participants, err := s.db.Chat.Participant.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	res := make([]GetRoomsResult, 0, len(participants))
	for _, p := range participants {
		var latestMessage *model.ChatMessage

		msg, err := s.db.Chat.Message.GetLatestMessage(ctx, p.RoomID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				latestMessage = nil
			} else {
				return nil, err
			}
		} else {
			latestMessage = &msg
		}

		room, err := s.db.Chat.Room.GetRoomByID(ctx, p.RoomID)
		if err != nil {
			return nil, err
		}

		// if last message of the room is nil, then read all
		// if last read message is not nil and last message of the room is equal to last read message, then read all
		readAll := latestMessage == nil || (p.LastReadMessageID != nil && latestMessage.ID == *p.LastReadMessageID)
		res = append(res, GetRoomsResult{
			RoomID:      p.RoomID,
			Title:       room.Title,
			LastMessage: latestMessage,
			ReadAll:     readAll,
		})
	}
	return res, nil
}

func (s *Service) GetParticipants(ctx context.Context, roomID string) ([]model.ChatParticipant, error) {
	return s.db.Chat.Participant.GetByRoomID(ctx, roomID)
}

func (s *Service) IsParticipant(ctx context.Context, roomID string, userID string) (bool, error) {
	participants, err := s.GetParticipants(ctx, roomID)
	if err != nil {
		return false, err
	}

	if !slices.ContainsFunc(
		participants,
		func(p model.ChatParticipant) bool { return p.UserID == userID },
	) {
		return false, nil
	}
	return true, nil
}

func (s *Service) UpdateLastReadMessageID(ctx context.Context, roomID string, userID string) error {
	messageID, err := s.db.Chat.Message.GetLatestMessage(ctx, roomID)
	if err != nil {
		return err
	}
	return s.db.Chat.Participant.UpdateLastReadMessageID(ctx, chatparticipant.UpdateLastReadMessageIDParams{
		RoomID:            roomID,
		UserID:            userID,
		LastReadMessageID: &messageID.ID,
	})
}
