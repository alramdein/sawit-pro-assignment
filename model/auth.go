package model

type JwtClaims struct {
	UserId    string `json:"user_id,omitempty"`
	ExpiresAt string `json:"expires_at,omitempty"`
}
