package dto

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type RefreshRequest struct {
	RefreshToken string `json:"RefreshToken" binding:"required"`
	DeviceID     string `json:"DeviceID" binding:"required"`
}
