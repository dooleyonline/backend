package model

import (
	useruser "github.com/dooleyonline/backend/internal/db/user/user"
	userverify "github.com/dooleyonline/backend/internal/db/user/verify"
)

type User = useruser.UserUser

type Verify = userverify.UserVerify
