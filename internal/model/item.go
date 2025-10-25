package model

import (
	itemcategory "github.com/dooleyonline/backend/internal/db/item/category"
	itemitem "github.com/dooleyonline/backend/internal/db/item/item"
)

type Item = itemitem.ItemItem

type Category = itemcategory.ItemCategory
