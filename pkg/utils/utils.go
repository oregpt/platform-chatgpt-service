package utils

import (
	"encoding/json"
	"fmt"
	"time"
)

// PrettyJSON converts a struct to a pretty-printed JSON string
func PrettyJSON(v interface{}) string {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error marshaling JSON: %v", err)
	}
	return string(bytes)
}

// FormatDuration formats a duration in a human-readable format
func FormatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%d µs", d.Microseconds())
	}
	if d < time.Second {
		return fmt.Sprintf("%d ms", d.Milliseconds())
	}
	if d < time.Minute {
		return fmt.Sprintf("%.2f s", float64(d)/float64(time.Second))
	}
	return fmt.Sprintf("%.2f min", float64(d)/float64(time.Minute))
}

// GenerateUUID generates a simple UUID-like string
func GenerateUUID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// TruncateString truncates a string to the specified length and adds "..." if truncated
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
