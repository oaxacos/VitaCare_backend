package dto

type TokenRefreshResponse struct {
	AccessToken string `json:"access_token"`
	User        User   `json:"user"`
}

type TokenRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}
