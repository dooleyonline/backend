package chatsvc

import (
	"context"

	"github.com/bwmarrin/snowflake"
	"github.com/dooleyonline/backend/internal/config"
	"github.com/dooleyonline/backend/internal/db"
	chatmessage "github.com/dooleyonline/backend/internal/db/chat/message"
	chatparticipant "github.com/dooleyonline/backend/internal/db/chat/participant"
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

func (s *Service) CreateMessage(ctx context.Context, user string, room string, message string) error {

	dbparams := chatmessage.CreateParams{
		ID:     s.GenerateID(),
		RoomID: room,
		SentBy: user,
		Body:   message,
	}
	if err := s.db.Chat.Message.Create(ctx, dbparams); err != nil {
		return err
	}
	return nil
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

func (s *Service) GetLatestMessages(ctx context.Context, roomID string) ([]model.ChatMessage, error) {
	return s.db.Chat.Message.GetMany(ctx, chatmessage.GetManyParams{
		RoomID: roomID,
		Limit:  20,
	})
}
