// Package db provides data access abstractions and implementations
// for the Crowd Vote application. It defines the storage interface
// and provides backends like SQLite and MongoDB.
package db

import "time"

// Store defines the interface for database operations, allowing for
// pluggable backends such as SQLite or MongoDB.
type Store interface {
	// Init initializes the database connection and required schemas/indexes.
	Init() error
	// Close gracefully closes the database connection.
	Close() error

	// ListLocations returns all locations registered within a group.
	ListLocations(group string) ([]Location, error)
	// CreateLocation persists a new location.
	CreateLocation(l Location) error
	// GetLocation retrieves a single location by its group and ID.
	GetLocation(group, id string) (Location, error)

	// CreateVote persists a new user vote report.
	CreateVote(v Vote) error
	// CalculateCrowdRate aggregates vote data to compute a crowdedness status.
	CalculateCrowdRate(group, id string, since time.Time) (CrowdRate, error)
	// ListLocationIDs returns all location IDs associated with a specific group.
	ListLocationIDs(group string) ([]string, error)
}
