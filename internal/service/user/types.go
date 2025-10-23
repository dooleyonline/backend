package usersvc

type CreateInput struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
}

type UserSummary struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type Me struct {
	ID         string  `json:"id"`
	Email      string  `json:"email"`
	FirstName  string  `json:"first_name"`
	LastName   string  `json:"last_name"`
	LikedItems []int64 `json:"liked_items"`
}

type Seller struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
