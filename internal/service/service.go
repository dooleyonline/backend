package service

import (
	"github.com/dooleyonline/backend/internal/config"
	"github.com/dooleyonline/backend/internal/db"
	authsvc "github.com/dooleyonline/backend/internal/service/auth"
	categorysvc "github.com/dooleyonline/backend/internal/service/category"
	chatsvc "github.com/dooleyonline/backend/internal/service/chat"
	itemsvc "github.com/dooleyonline/backend/internal/service/item"
	usersvc "github.com/dooleyonline/backend/internal/service/user"
)

type Service struct {
	Auth     *authsvc.Service
	Category *categorysvc.Service
	Item     *itemsvc.Service
	User     *usersvc.Service
	Chat     *chatsvc.Service
}

func New(cfg *config.Config, db *db.DB) *Service {
	return &Service{
		Auth:     authsvc.New(cfg, db),
		Category: categorysvc.New(db),
		Item:     itemsvc.New(cfg, db),
		User:     usersvc.New(db),
		Chat:     chatsvc.New(cfg, db),
	}
}
