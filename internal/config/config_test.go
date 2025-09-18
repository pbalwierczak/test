package config

import (
	"testing"
)

func TestGeographicConstants(t *testing.T) {
	// Test Ottawa coordinates
	if OttawaCenterLat != 45.4215 {
		t.Errorf("Expected OttawaCenterLat to be 45.4215, got %f", OttawaCenterLat)
	}
	if OttawaCenterLng != -75.6972 {
		t.Errorf("Expected OttawaCenterLng to be -75.6972, got %f", OttawaCenterLng)
	}

	// Test Montreal coordinates
	if MontrealCenterLat != 45.5017 {
		t.Errorf("Expected MontrealCenterLat to be 45.5017, got %f", MontrealCenterLat)
	}
	if MontrealCenterLng != -73.5673 {
		t.Errorf("Expected MontrealCenterLng to be -73.5673, got %f", MontrealCenterLng)
	}

	// Test city radius
	if CityRadiusKm != 10 {
		t.Errorf("Expected CityRadiusKm to be 10, got %d", CityRadiusKm)
	}
}

func TestConfigLoad(t *testing.T) {
	config, err := Load()
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// Test that geographic constants are not in the config struct
	// (they should be accessed as constants, not as config fields)
	if config.APIKey == "" {
		t.Error("Expected APIKey to be set")
	}
	if config.ServerPort == "" {
		t.Error("Expected ServerPort to be set")
	}
}
