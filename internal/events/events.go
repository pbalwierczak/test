package events

import (
	"time"

	"github.com/google/uuid"
)

type BaseEvent struct {
	EventType string    `json:"eventType"`
	EventID   string    `json:"eventId"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

type TripStartedEvent struct {
	BaseEvent
	Data TripStartedData `json:"data"`
}

type TripStartedData struct {
	TripID         string  `json:"tripId"`
	ScooterID      string  `json:"scooterId"`
	UserID         string  `json:"userId"`
	StartLatitude  float64 `json:"startLatitude"`
	StartLongitude float64 `json:"startLongitude"`
	StartTime      string  `json:"startTime"`
}

type TripEndedEvent struct {
	BaseEvent
	Data TripEndedData `json:"data"`
}

type TripEndedData struct {
	TripID          string  `json:"tripId"`
	ScooterID       string  `json:"scooterId"`
	UserID          string  `json:"userId"`
	EndLatitude     float64 `json:"endLatitude"`
	EndLongitude    float64 `json:"endLongitude"`
	EndTime         string  `json:"endTime"`
	DurationSeconds int     `json:"durationSeconds"`
}

type LocationUpdatedEvent struct {
	BaseEvent
	Data LocationUpdatedData `json:"data"`
}

type LocationUpdatedData struct {
	ScooterID string  `json:"scooterId"`
	TripID    string  `json:"tripId"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Heading   float64 `json:"heading"`
	Speed     float64 `json:"speed"`
}

func NewTripStartedEvent(tripID, scooterID, userID string, startLat, startLng float64) *TripStartedEvent {
	now := time.Now()
	return &TripStartedEvent{
		BaseEvent: BaseEvent{
			EventType: "trip.started",
			EventID:   uuid.New().String(),
			Timestamp: now,
			Version:   "1.0",
		},
		Data: TripStartedData{
			TripID:         tripID,
			ScooterID:      scooterID,
			UserID:         userID,
			StartLatitude:  startLat,
			StartLongitude: startLng,
			StartTime:      now.Format(time.RFC3339),
		},
	}
}

func NewTripEndedEvent(tripID, scooterID, userID string, endLat, endLng float64, startTime time.Time) *TripEndedEvent {
	now := time.Now()
	duration := int(now.Sub(startTime).Seconds())

	return &TripEndedEvent{
		BaseEvent: BaseEvent{
			EventType: "trip.ended",
			EventID:   uuid.New().String(),
			Timestamp: now,
			Version:   "1.0",
		},
		Data: TripEndedData{
			TripID:          tripID,
			ScooterID:       scooterID,
			UserID:          userID,
			EndLatitude:     endLat,
			EndLongitude:    endLng,
			EndTime:         now.Format(time.RFC3339),
			DurationSeconds: duration,
		},
	}
}

func NewLocationUpdatedEvent(scooterID, tripID string, lat, lng, heading, speed float64) *LocationUpdatedEvent {
	return &LocationUpdatedEvent{
		BaseEvent: BaseEvent{
			EventType: "location.updated",
			EventID:   uuid.New().String(),
			Timestamp: time.Now(),
			Version:   "1.0",
		},
		Data: LocationUpdatedData{
			ScooterID: scooterID,
			TripID:    tripID,
			Latitude:  lat,
			Longitude: lng,
			Heading:   heading,
			Speed:     speed,
		},
	}
}
