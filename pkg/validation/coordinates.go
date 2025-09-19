package validation

import "errors"

func ValidateCoordinates(lat, lng float64) error {
	if lat < -90 || lat > 90 {
		return errors.New("invalid latitude: must be between -90 and 90")
	}

	if lng < -180 || lng > 180 {
		return errors.New("invalid longitude: must be between -180 and 180")
	}

	return nil
}
