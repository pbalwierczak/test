package validation

import (
	"testing"
)

func TestValidateCoordinates(t *testing.T) {
	tests := []struct {
		name        string
		lat         float64
		lng         float64
		expectError bool
		errorMsg    string
	}{
		// Valid coordinates
		{
			name:        "valid coordinates - center of map",
			lat:         0.0,
			lng:         0.0,
			expectError: false,
		},
		{
			name:        "valid coordinates - positive values",
			lat:         40.7128,
			lng:         -74.0060,
			expectError: false,
		},
		{
			name:        "valid coordinates - negative values",
			lat:         -33.9249,
			lng:         18.4241,
			expectError: false,
		},
		{
			name:        "valid coordinates - boundary values",
			lat:         90.0,
			lng:         180.0,
			expectError: false,
		},
		{
			name:        "valid coordinates - negative boundary values",
			lat:         -90.0,
			lng:         -180.0,
			expectError: false,
		},
		{
			name:        "valid coordinates - small values",
			lat:         0.0001,
			lng:         0.0001,
			expectError: false,
		},
		{
			name:        "valid coordinates - large values within bounds",
			lat:         89.9999,
			lng:         179.9999,
			expectError: false,
		},

		// Invalid latitude
		{
			name:        "invalid latitude - too high",
			lat:         90.1,
			lng:         0.0,
			expectError: true,
			errorMsg:    "invalid latitude: must be between -90 and 90",
		},
		{
			name:        "invalid latitude - too low",
			lat:         -90.1,
			lng:         0.0,
			expectError: true,
			errorMsg:    "invalid latitude: must be between -90 and 90",
		},
		{
			name:        "invalid latitude - way too high",
			lat:         100.0,
			lng:         0.0,
			expectError: true,
			errorMsg:    "invalid latitude: must be between -90 and 90",
		},
		{
			name:        "invalid latitude - way too low",
			lat:         -100.0,
			lng:         0.0,
			expectError: true,
			errorMsg:    "invalid latitude: must be between -90 and 90",
		},

		// Invalid longitude
		{
			name:        "invalid longitude - too high",
			lat:         0.0,
			lng:         180.1,
			expectError: true,
			errorMsg:    "invalid longitude: must be between -180 and 180",
		},
		{
			name:        "invalid longitude - too low",
			lat:         0.0,
			lng:         -180.1,
			expectError: true,
			errorMsg:    "invalid longitude: must be between -180 and 180",
		},
		{
			name:        "invalid longitude - way too high",
			lat:         0.0,
			lng:         200.0,
			expectError: true,
			errorMsg:    "invalid longitude: must be between -180 and 180",
		},
		{
			name:        "invalid longitude - way too low",
			lat:         0.0,
			lng:         -200.0,
			expectError: true,
			errorMsg:    "invalid longitude: must be between -180 and 180",
		},

		// Both invalid
		{
			name:        "both coordinates invalid - latitude error should be returned first",
			lat:         91.0,
			lng:         181.0,
			expectError: true,
			errorMsg:    "invalid latitude: must be between -90 and 90",
		},
		{
			name:        "both coordinates invalid - negative values",
			lat:         -91.0,
			lng:         -181.0,
			expectError: true,
			errorMsg:    "invalid latitude: must be between -90 and 90",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCoordinates(tt.lat, tt.lng)

			if tt.expectError {
				if err == nil {
					t.Errorf("ValidateCoordinates() expected error but got none")
					return
				}
				if err.Error() != tt.errorMsg {
					t.Errorf("ValidateCoordinates() error message = %v, want %v", err.Error(), tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateCoordinates() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestValidateLatitude(t *testing.T) {
	tests := []struct {
		name        string
		lat         float64
		expectError bool
		errorMsg    string
	}{
		// Valid latitudes
		{
			name:        "valid latitude - zero",
			lat:         0.0,
			expectError: false,
		},
		{
			name:        "valid latitude - positive value",
			lat:         40.7128,
			expectError: false,
		},
		{
			name:        "valid latitude - negative value",
			lat:         -33.9249,
			expectError: false,
		},
		{
			name:        "valid latitude - maximum positive",
			lat:         90.0,
			expectError: false,
		},
		{
			name:        "valid latitude - maximum negative",
			lat:         -90.0,
			expectError: false,
		},
		{
			name:        "valid latitude - small positive value",
			lat:         0.0001,
			expectError: false,
		},
		{
			name:        "valid latitude - small negative value",
			lat:         -0.0001,
			expectError: false,
		},
		{
			name:        "valid latitude - close to maximum positive",
			lat:         89.9999,
			expectError: false,
		},
		{
			name:        "valid latitude - close to maximum negative",
			lat:         -89.9999,
			expectError: false,
		},

		// Invalid latitudes
		{
			name:        "invalid latitude - too high",
			lat:         90.1,
			expectError: true,
			errorMsg:    "invalid latitude: must be between -90 and 90",
		},
		{
			name:        "invalid latitude - too low",
			lat:         -90.1,
			expectError: true,
			errorMsg:    "invalid latitude: must be between -90 and 90",
		},
		{
			name:        "invalid latitude - way too high",
			lat:         100.0,
			expectError: true,
			errorMsg:    "invalid latitude: must be between -90 and 90",
		},
		{
			name:        "invalid latitude - way too low",
			lat:         -100.0,
			expectError: true,
			errorMsg:    "invalid latitude: must be between -90 and 90",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLatitude(tt.lat)

			if tt.expectError {
				if err == nil {
					t.Errorf("ValidateLatitude() expected error but got none")
					return
				}
				if err.Error() != tt.errorMsg {
					t.Errorf("ValidateLatitude() error message = %v, want %v", err.Error(), tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateLatitude() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestValidateLongitude(t *testing.T) {
	tests := []struct {
		name        string
		lng         float64
		expectError bool
		errorMsg    string
	}{
		// Valid longitudes
		{
			name:        "valid longitude - zero",
			lng:         0.0,
			expectError: false,
		},
		{
			name:        "valid longitude - positive value",
			lng:         120.0,
			expectError: false,
		},
		{
			name:        "valid longitude - negative value",
			lng:         -120.0,
			expectError: false,
		},
		{
			name:        "valid longitude - maximum positive",
			lng:         180.0,
			expectError: false,
		},
		{
			name:        "valid longitude - maximum negative",
			lng:         -180.0,
			expectError: false,
		},
		{
			name:        "valid longitude - small positive value",
			lng:         0.0001,
			expectError: false,
		},
		{
			name:        "valid longitude - small negative value",
			lng:         -0.0001,
			expectError: false,
		},
		{
			name:        "valid longitude - close to maximum positive",
			lng:         179.9999,
			expectError: false,
		},
		{
			name:        "valid longitude - close to maximum negative",
			lng:         -179.9999,
			expectError: false,
		},

		// Invalid longitudes
		{
			name:        "invalid longitude - too high",
			lng:         180.1,
			expectError: true,
			errorMsg:    "invalid longitude: must be between -180 and 180",
		},
		{
			name:        "invalid longitude - too low",
			lng:         -180.1,
			expectError: true,
			errorMsg:    "invalid longitude: must be between -180 and 180",
		},
		{
			name:        "invalid longitude - way too high",
			lng:         200.0,
			expectError: true,
			errorMsg:    "invalid longitude: must be between -180 and 180",
		},
		{
			name:        "invalid longitude - way too low",
			lng:         -200.0,
			expectError: true,
			errorMsg:    "invalid longitude: must be between -180 and 180",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLongitude(tt.lng)

			if tt.expectError {
				if err == nil {
					t.Errorf("ValidateLongitude() expected error but got none")
					return
				}
				if err.Error() != tt.errorMsg {
					t.Errorf("ValidateLongitude() error message = %v, want %v", err.Error(), tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateLongitude() unexpected error = %v", err)
				}
			}
		})
	}
}

// Benchmark tests to ensure performance is acceptable
func BenchmarkValidateCoordinates(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidateCoordinates(40.7128, -74.0060)
	}
}

func BenchmarkValidateLatitude(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidateLatitude(40.7128)
	}
}

func BenchmarkValidateLongitude(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidateLongitude(-74.0060)
	}
}
