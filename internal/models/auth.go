package models

type RegisterBody struct {
	Login      string `json:"login"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	InviteCode string `json:"inviteCode"`
}

type LoginBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
