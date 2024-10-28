package responses

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	Exp          int64  `json:"exp"`
}

func NewLoginResponse(token string, exp int64) *LoginResponse {
	return &LoginResponse{
		AccessToken:  token,
		Exp:          exp,
	}
}