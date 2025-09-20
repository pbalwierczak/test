package events

import (
	"context"
	"errors"
	"testing"

	"scootin-aboot/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTripStartedHandler_Handle(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		setupMocks  func(*MockTripService)
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid trip started event",
			data: []byte(`{"eventType":"trip.started","eventId":"test-id","timestamp":"2023-01-01T00:00:00Z","version":"1.0","data":{"tripId":"550e8400-e29b-41d4-a716-446655440000","scooterId":"550e8400-e29b-41d4-a716-446655440001","userId":"550e8400-e29b-41d4-a716-446655440002","startLatitude":45.4215,"startLongitude":-75.6972,"startTime":"2023-01-01T00:00:00Z"}}`),
			setupMocks: func(tripService *MockTripService) {
				tripService.On("StartTrip", mock.Anything, mock.Anything, mock.Anything, 45.4215, -75.6972).Return(&models.Trip{}, nil)
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
			data: []byte(`{"eventType":"trip.started","eventId":"test-id","timestamp":"2023-01-01T00:00:00Z","version":"1.0","data":{"tripId":"550e8400-e29b-41d4-a716-446655440000","scooterId":"invalid-uuid","userId":"550e8400-e29b-41d4-a716-446655440002","startLatitude":45.4215,"startLongitude":-75.6972,"startTime":"2023-01-01T00:00:00Z"}}`),
			setupMocks: func(tripService *MockTripService) {
			},
			expectError: true,
			errorMsg:    "invalid scooter ID",
		},
		{
			name: "invalid user ID",
			data: []byte(`{"eventType":"trip.started","eventId":"test-id","timestamp":"2023-01-01T00:00:00Z","version":"1.0","data":{"tripId":"550e8400-e29b-41d4-a716-446655440000","scooterId":"550e8400-e29b-41d4-a716-446655440001","userId":"invalid-uuid","startLatitude":45.4215,"startLongitude":-75.6972,"startTime":"2023-01-01T00:00:00Z"}}`),
			setupMocks: func(tripService *MockTripService) {
			},
			expectError: true,
			errorMsg:    "invalid user ID",
		},
		{
			name: "trip service error",
			data: []byte(`{"eventType":"trip.started","eventId":"test-id","timestamp":"2023-01-01T00:00:00Z","version":"1.0","data":{"tripId":"550e8400-e29b-41d4-a716-446655440000","scooterId":"550e8400-e29b-41d4-a716-446655440001","userId":"550e8400-e29b-41d4-a716-446655440002","startLatitude":45.4215,"startLongitude":-75.6972,"startTime":"2023-01-01T00:00:00Z"}}`),
			setupMocks: func(tripService *MockTripService) {
				tripService.On("StartTrip", mock.Anything, mock.Anything, mock.Anything, 45.4215, -75.6972).Return((*models.Trip)(nil), errors.New("service error"))
			},
			expectError: true,
			errorMsg:    "failed to start trip",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tripService := &MockTripService{}
			scooterService := &MockScooterService{}
			tt.setupMocks(tripService)

			handler := NewTripStartedHandler(HandlerDependencies{
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
