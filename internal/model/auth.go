package model

import "time"

const (
	AUTHORIZED = "authorized"
	CANCELLED  = "cancelled"
)

type Auth struct {
	ShipmentID     uint64    `json:"shipment_id"`
	TrackingNumber string    `json:"tracking_number"`
	DateCreated    time.Time `json:"create_at"`
	SenderID       uint64    `json:"sender_id"`
	DestinationID  uint64    `json:"destination_id"`
	FacilityID     string    `json:"facility_id"`
	Status         string    `json:"status"`
}
