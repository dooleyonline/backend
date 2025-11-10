package model

import (
	userliked "github.com/dooleyonline/backend/internal/db/user/liked"
	useruser "github.com/dooleyonline/backend/internal/db/user/user"
  userverify "github.com/dooleyonline/backend/internal/db/user/verify"
	userviewed "github.com/dooleyonline/backend/internal/db/user/viewed"
)

type User = useruser.UserUser
type Liked = userliked.UserLiked
type Viewed = userviewed.UserViewed
type Verification = userverify.UserVerification
