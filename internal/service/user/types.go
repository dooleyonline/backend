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
	ID         string
	Email      string
	FirstName  string
	LastName   string
	LikedItems []int64
}

type Seller struct {
	ID        string
	FirstName string
	LastName  string
}
