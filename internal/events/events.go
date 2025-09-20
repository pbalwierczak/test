package events

import (
	"time"

	"github.com/google/uuid"
)

// BaseEvent represents the common structure for all events
type BaseEvent struct {
	EventType string    `json:"eventType"`
	EventID   string    `json:"eventId"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

// TripStartedEvent represents a trip start event
type TripStartedEvent struct {
	BaseEvent
	Data TripStartedData `json:"data"`
}

// TripStartedData contains the data for a trip start event
type TripStartedData struct {
	TripID         string  `json:"tripId"`
	ScooterID      string  `json:"scooterId"`
	UserID         string  `json:"userId"`
	StartLatitude  float64 `json:"startLatitude"`
	StartLongitude float64 `json:"startLongitude"`
	StartTime      string  `json:"startTime"`
}

// TripEndedEvent represents a trip end event
type TripEndedEvent struct {
	BaseEvent
	Data TripEndedData `json:"data"`
}

// TripEndedData contains the data for a trip end event
type TripEndedData struct {
	TripID          string  `json:"tripId"`
	ScooterID       string  `json:"scooterId"`
	UserID          string  `json:"userId"`
	EndLatitude     float64 `json:"endLatitude"`
	EndLongitude    float64 `json:"endLongitude"`
	EndTime         string  `json:"endTime"`
	DurationSeconds int     `json:"durationSeconds"`
}

// LocationUpdatedEvent represents a location update event
type LocationUpdatedEvent struct {
	BaseEvent
	Data LocationUpdatedData `json:"data"`
}

// LocationUpdatedData contains the data for a location update event
type LocationUpdatedData struct {
	ScooterID string  `json:"scooterId"`
	TripID    string  `json:"tripId"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Heading   float64 `json:"heading"`
	Speed     float64 `json:"speed"`
}

// NewTripStartedEvent creates a new trip started event
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

// NewTripEndedEvent creates a new trip ended event
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

// NewLocationUpdatedEvent creates a new location updated event
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
