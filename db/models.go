package db

import "time"

// Location represents a physical place defined by its group and ID.
type Location struct {
	Group       string `json:"group" bson:"group"`
	ID          string `json:"id" bson:"id"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
}

// Vote represents a user's report of crowdedness at a specific location.
type Vote struct {
	Group         string    `json:"group" bson:"group"`
	LocationID    string    `json:"location_id" bson:"location_id"`
	ClientAddress string    `json:"client_address" bson:"client_address"`
	Level         int       `json:"level" bson:"level"`
	Timestamp     time.Time `json:"timestamp" bson:"timestamp"`
}

// CrowdRate represents the aggregated crowdedness status of a location
// within a specified time range.
type CrowdRate struct {
	Group      string    `json:"group" bson:"group"`
	LocationID string    `json:"location_id" bson:"location_id"`
	Votes      float64   `json:"votes" bson:"votes"`
	Total      int       `json:"total" bson:"total"`
	Since      time.Time `json:"since" bson:"since"`
}
