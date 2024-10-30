package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type WidgetProperties struct {
	Name         string `json:"name"`
	SerialNumber string `json:"serial_number"`
	PortTypeP    bool   `json:"port_type_p"`
	PortTypeR    bool   `json:"port_type_r"`
	PortTypeQ    bool   `json:"port_type_q"`
	Timestamp    int64  `json:"timestamp"`
	IsActive     bool   `json:"is_active"`
	ValidUntil   any    `json:"valid_until"` // Very dangerous but IDK another idea
}

type WidgetMongoProperties struct {
	ID           primitive.ObjectID `bson:"_id"`
	Name         string             `json:"name"`
	SerialNumber string             `json:"serial_number"`
	PortTypeP    bool               `json:"port_type_p"`
	PortTypeR    bool               `json:"port_type_r"`
	PortTypeQ    bool               `json:"port_type_q"`
	Timestamp    int64              `json:"timestamp"`
	IsActive     bool               `json:"is_active"`
	ValidUntil   any                `json:"valid_until"` // Very dangerous but IDK another idea
}

type WidgetDeregistration struct {
	SerialNumber string `json:"serial_number"`
	IsActive     bool   `json:"is_active"`
	ValidUntil   any    `json:"valid_until"` // Very dangerous but IDK another idea
}
