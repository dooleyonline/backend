package chatsvc

import (
	"context"
	"slices"

	"github.com/bwmarrin/snowflake"
	"github.com/dooleyonline/backend/internal/config"
	"github.com/dooleyonline/backend/internal/db"
	chatmessage "github.com/dooleyonline/backend/internal/db/chat/message"
	chatparticipant "github.com/dooleyonline/backend/internal/db/chat/participant"
	chatroom "github.com/dooleyonline/backend/internal/db/chat/room"
	"github.com/dooleyonline/backend/internal/model"
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

	dbparams := chatmessage.CreateParams{
		ID:     s.GenerateID(),
		RoomID: p.RoomID,
		SentBy: p.UserID,
		Body:   p.Message,
	}
	if err := s.db.Chat.Message.Create(ctx, dbparams); err != nil {
		return err
	}
	return nil
}

type EditMessageParams struct {
	ID   int64
	Body string
}

func (s *Service) EditMessage(ctx context.Context, p *EditMessageParams) error {
	dbparams := chatmessage.EditMessageParams{
		ID:   p.ID,
		Body: p.Body,
	}
	if err := s.db.Chat.Message.EditMessage(ctx, dbparams); err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteMessage(ctx context.Context, messageID int64) error {
	return s.db.Chat.Message.Delete(ctx, messageID)
}

func (s *Service) CreateRoom(ctx context.Context, users []string) (string, error) {
	id, err := s.db.Chat.Room.CreateRoom(ctx, users)
	if err != nil {
		return "", err
	}

	for _, user := range users {
		if err := s.db.Chat.Participant.Create(ctx, chatparticipant.CreateParams{
			RoomID: id,
			UserID: user,
		}); err != nil {
			return "", err
		}
	}
	return id, nil
}

func (s *Service) DeleteRoom(ctx context.Context, roomID string) error {
	return s.db.Chat.Room.DeleteRoom(ctx, roomID)
}

func (s *Service) GetLatestMessages(ctx context.Context, roomID string) ([]model.ChatMessage, error) {
	return s.db.Chat.Message.GetMany(ctx, chatmessage.GetManyParams{
		RoomID: roomID,
		Limit:  20,
	})
}

type ParticipantParams struct {
	RoomID string
	UserID string
}

func (s *Service) AddParticipant(ctx context.Context, p *ParticipantParams) error {
	if err := s.db.Chat.Room.AddParticipant(ctx, chatroom.AddParticipantParams{
		UserID: p.UserID,
		RoomID: p.RoomID,
	}); err != nil {
		return err
	}

	if err := s.db.Chat.Participant.Create(ctx, chatparticipant.CreateParams{
		RoomID: p.RoomID,
		UserID: p.UserID,
	}); err != nil {
		return err
	}

	return nil
}

func (s *Service) RemoveParticipant(ctx context.Context, p *ParticipantParams) error {
	if err := s.db.Chat.Room.RemoveParticipant(ctx, chatroom.RemoveParticipantParams{
		UserID: p.UserID,
		RoomID: p.RoomID,
	}); err != nil {
		return err
	}

	if err := s.db.Chat.Participant.Delete(ctx, chatparticipant.DeleteParams{
		RoomID: p.RoomID,
		UserID: p.UserID,
	}); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetRooms(ctx context.Context, userID string) ([]model.ChatParticipant, error) {
	return s.db.Chat.Participant.GetByUserID(ctx, userID)
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
