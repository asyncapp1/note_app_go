package types

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserClaims struct {
	ID   uint   `json:"user_id"`
	Type string `json:"token_type"`
}
