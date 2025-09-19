package kafka

import (
	"context"
	"errors"
	"testing"

	"scootin-aboot/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTripEndedHandler_Handle(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		setupMocks  func(*MockTripService)
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid trip ended event",
			data: []byte(`{"eventType":"trip.ended","eventId":"test-id","timestamp":"2023-01-01T00:00:00Z","version":"1.0","data":{"tripId":"trip-123","scooterId":"550e8400-e29b-41d4-a716-446655440001","userId":"user-123","endLatitude":45.4216,"endLongitude":-75.6973,"endTime":"2023-01-01T00:30:00Z","durationSeconds":1800}}`),
			setupMocks: func(tripService *MockTripService) {
				tripService.On("EndTrip", mock.Anything, mock.Anything, 45.4216, -75.6973).Return(&models.Trip{}, nil)
			},
			expectError: false,
		},
		{
			name: "invalid JSON",
			data: []byte(`invalid json`),
			setupMocks: func(tripService *MockTripService) {
			},
			expectError: true,
			errorMsg:    "failed to unmarshal",
		},
		{
			name: "invalid scooter ID",
			data: []byte(`{"eventType":"trip.ended","eventId":"test-id","timestamp":"2023-01-01T00:00:00Z","version":"1.0","data":{"tripId":"trip-123","scooterId":"invalid-uuid","userId":"user-123","endLatitude":45.4216,"endLongitude":-75.6973,"endTime":"2023-01-01T00:30:00Z","durationSeconds":1800}}`),
			setupMocks: func(tripService *MockTripService) {
			},
			expectError: true,
			errorMsg:    "invalid scooter ID",
		},
		{
			name: "trip service error",
			data: []byte(`{"eventType":"trip.ended","eventId":"test-id","timestamp":"2023-01-01T00:00:00Z","version":"1.0","data":{"tripId":"trip-123","scooterId":"550e8400-e29b-41d4-a716-446655440001","userId":"user-123","endLatitude":45.4216,"endLongitude":-75.6973,"endTime":"2023-01-01T00:30:00Z","durationSeconds":1800}}`),
			setupMocks: func(tripService *MockTripService) {
				tripService.On("EndTrip", mock.Anything, mock.Anything, 45.4216, -75.6973).Return((*models.Trip)(nil), errors.New("service error"))
			},
			expectError: true,
			errorMsg:    "failed to end trip",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tripService := &MockTripService{}
			scooterService := &MockScooterService{}
			tt.setupMocks(tripService)

			handler := NewTripEndedHandler(HandlerDependencies{
				TripService:    tripService,
				ScooterService: scooterService,
			})

			err := handler.Handle(context.Background(), tt.data)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}

			tripService.AssertExpectations(t)
		})
	}
}
