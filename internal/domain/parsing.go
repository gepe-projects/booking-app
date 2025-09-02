package domain

import "time"

func NilStringHandler(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func BoolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func TimeToString(t time.Time) string {
	return t.Format(time.RFC3339)
}
