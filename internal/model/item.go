package model

import (
	itemcategory "github.com/dooleyonline/backend/internal/db/item/category"
	itemitem "github.com/dooleyonline/backend/internal/db/item/item"
	useruser "github.com/dooleyonline/backend/internal/db/user/user"
)

type Item = itemitem.ItemItem

type Category = itemcategory.ItemCategory

type User = useruser.UserUser
