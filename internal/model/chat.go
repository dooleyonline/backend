package model

import (
	chatmessage "github.com/dooleyonline/backend/internal/db/chat/message"
	chatparticipant "github.com/dooleyonline/backend/internal/db/chat/participant"
	chatroom "github.com/dooleyonline/backend/internal/db/chat/room"
)

type ChatMessage = chatmessage.ChatMessage

type ChatRoom = chatroom.ChatRoom

type ChatParticipant = chatparticipant.ChatParticipant
