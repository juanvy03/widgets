package model

import (
	"time"
)

type WidgetEvent struct {
	EventType string      `json:"event_type"`
	Event     interface{} `json:"event"`
	Timestamp time.Time   `json:"timestamp"`
}
