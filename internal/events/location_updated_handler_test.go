package events

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLocationUpdatedHandler_Handle(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		setupMocks  func(*MockTripService)
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid location updated event",
			data: []byte(`{"eventType":"location.updated","eventId":"test-id","timestamp":"2023-01-01T00:00:00Z","version":"1.0","data":{"scooterId":"550e8400-e29b-41d4-a716-446655440001","tripId":"trip-123","latitude":45.4216,"longitude":-75.6973,"heading":90.0,"speed":15.5}}`),
			setupMocks: func(tripService *MockTripService) {
				tripService.On("UpdateLocation", mock.Anything, mock.Anything, 45.4216, -75.6973).Return(nil)
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
			data: []byte(`{"eventType":"location.updated","eventId":"test-id","timestamp":"2023-01-01T00:00:00Z","version":"1.0","data":{"scooterId":"invalid-uuid","tripId":"trip-123","latitude":45.4216,"longitude":-75.6973,"heading":90.0,"speed":15.5}}`),
			setupMocks: func(tripService *MockTripService) {
			},
			expectError: true,
			errorMsg:    "invalid scooter ID",
		},
		{
			name: "trip service error",
			data: []byte(`{"eventType":"location.updated","eventId":"test-id","timestamp":"2023-01-01T00:00:00Z","version":"1.0","data":{"scooterId":"550e8400-e29b-41d4-a716-446655440001","tripId":"trip-123","latitude":45.4216,"longitude":-75.6973,"heading":90.0,"speed":15.5}}`),
			setupMocks: func(tripService *MockTripService) {
				tripService.On("UpdateLocation", mock.Anything, mock.Anything, 45.4216, -75.6973).Return(errors.New("service error"))
			},
			expectError: true,
			errorMsg:    "failed to update location",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tripService := &MockTripService{}
			scooterService := &MockScooterService{}
			tt.setupMocks(tripService)

			handler := NewLocationUpdatedHandler(HandlerDependencies{
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
