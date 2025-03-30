package dto

type TokenRefreshResponse struct {
	AccessToken string `json:"access_token"`
}

type TokenRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}
