package dto

type ProductInput struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type UserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type JWTInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Error struct {
	Message string `json:"message"`
}

type JWTOutput struct {
	AccessToken string `json:"access_token"`
}
