package validation

import "errors"

// ValidateCoordinates validates latitude and longitude values
// Returns an error if coordinates are invalid
func ValidateCoordinates(lat, lng float64) error {
	// Validate latitude (-90 to 90)
	if lat < -90 || lat > 90 {
		return errors.New("invalid latitude: must be between -90 and 90")
	}

	// Validate longitude (-180 to 180)
	if lng < -180 || lng > 180 {
		return errors.New("invalid longitude: must be between -180 and 180")
	}

	return nil
}
