package models

type User struct {
	Id        string `json:"id"`
	FirtsName string `json:"name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Birthday  string `json:"birthday"`
	ImageUrl  string `json:"image_url"`
	Card_num  string `json:"card_num"`
	Phone     string `json:"phone"`
	Role      string `json:"role"`
}

type Users struct {
	Users []*User `json:"users"`
}

type RegisterResponseModel struct {
	UserID       string
	AccessToken  string
	RefreshToken string
}
