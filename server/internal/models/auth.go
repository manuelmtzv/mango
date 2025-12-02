package models

type LoginResponse struct {
	Username    string `json:"username"`
	AccessToken string `json:"accessToken"`
}
