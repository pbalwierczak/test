package repository

// Repository defines the interface for all repository operations
type Repository interface {
	Scooter() ScooterRepository
	Trip() TripRepository
	User() UserRepository
	LocationUpdate() LocationUpdateRepository
}
