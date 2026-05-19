package db

import (
	"os"
	"testing"
	"time"
)

func TestSQLiteStore(t *testing.T) {
	dbPath := "test_crowd_vote.db"
	defer os.Remove(dbPath)

	store := NewSQLiteStore(dbPath)
	if err := store.Init(); err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}
	defer store.Close()

	// Test Create and Get Location
	loc := Location{
		Group:       "test-group",
		ID:          "test-loc",
		Name:        "Test Location",
		Description: "A test description",
	}

	if err := store.CreateLocation(loc); err != nil {
		t.Errorf("Failed to create location: %v", err)
	}

	gotLoc, err := store.GetLocation("test-group", "test-loc")
	if err != nil {
		t.Errorf("Failed to get location: %v", err)
	}
	if gotLoc.Name != loc.Name {
		t.Errorf("Expected name %s, got %s", loc.Name, gotLoc.Name)
	}

	// Test List Locations
	locs, err := store.ListLocations("test-group")
	if err != nil {
		t.Errorf("Failed to list locations: %v", err)
	}
	if len(locs) != 1 {
		t.Errorf("Expected 1 location, got %d", len(locs))
	}

	// Test List Location IDs
	ids, err := store.ListLocationIDs("test-group")
	if err != nil {
		t.Errorf("Failed to list location IDs: %v", err)
	}
	if len(ids) != 1 || ids[0] != "test-loc" {
		t.Errorf("Expected [test-loc], got %v", ids)
	}

	// Test Create Vote and Calculate Crowd Rate
	vote := Vote{
		Group:         "test-group",
		LocationID:    "test-loc",
		ClientAddress: "127.0.0.1",
		Level:         2,
		Timestamp:     time.Now(),
	}

	if err := store.CreateVote(vote); err != nil {
		t.Errorf("Failed to create vote: %v", err)
	}

	since := time.Now().Add(-1 * time.Hour)
	rate, err := store.CalculateCrowdRate("test-group", "test-loc", since)
	if err != nil {
		t.Errorf("Failed to calculate crowd rate: %v", err)
	}
	if rate.Total != 1 {
		t.Errorf("Expected total 1, got %d", rate.Total)
	}
	if rate.Votes != 2.0 {
		t.Errorf("Expected votes 2.0, got %f", rate.Votes)
	}
}
