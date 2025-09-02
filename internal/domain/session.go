package domain

import "time"

type Session struct {
	UserID     string `redis:"userID"`
	Name       string `redis:"name"`
	Role       string `redis:"role"`
	ImageURL   string `redis:"image_url"`
	Email      string `redis:"email"`
	Provider   string `redis:"provider"`
	ProviderID string `redis:"provider_id"`
	Phone      string `redis:"phone"`
	Verified   string `redis:"verified"`

	// Device info
	Device    string `redis:"device"`     // e.g. "iPhone 15 Pro", "MacBook", "Chrome on Windows"
	UserAgent string `redis:"user_agent"` // raw UA string
	IpAddress string `redis:"ip_address"` // IP login
}

// ToRedisMap - Convert Session to map for Redis
func (s *Session) ToRedisMap() map[string]interface{} {
	return map[string]interface{}{
		"userID":      s.UserID,
		"name":        s.Name,
		"role":        s.Role,
		"image_url":   s.ImageURL,
		"email":       s.Email,
		"provider":    s.Provider,
		"provider_id": s.ProviderID,
		"phone":       s.Phone,
		"verified":    s.Verified,
		"device":      s.Device,
		"user_agent":  s.UserAgent,
		"ip_address":  s.IpAddress,
	}
}

type SessionWithExpiry struct {
	Session    Session   `json:"session"`
	ExpireTime time.Time `json:"expire_time"`
	Token      string    `json:"token,omitempty"`
}
