package domain

type HttpResponse struct {
	Success bool `json:"success"`
	Message any  `json:"message,omitempty"`
	Data    any  `json:"data,omitempty"`
}
