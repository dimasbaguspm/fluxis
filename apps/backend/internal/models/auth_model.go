package models

type AuthLoginInputModel struct {
	Username string `json:"username" minLength:"1" doc:"Your username"`
	Password string `json:"password" minLength:"1" doc:"Your password"`
}

type AuthLoginOutputModel struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Username     string `json:"username"`
}

type AuthRefreshInputModel struct {
	RefreshToken string `json:"refreshToken"`
}

type AuthRefreshOutputModel struct {
	AccessToken string `json:"accessToken"`
}
