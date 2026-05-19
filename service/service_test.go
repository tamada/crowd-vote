package service

import (
	"crowd-vote/db"
	"testing"
	"time"
)

type mockStore struct {
	locations []db.Location
	votes     []db.Vote
}

func (m *mockStore) Init() error  { return nil }
func (m *mockStore) Close() error { return nil }

func (m *mockStore) ListLocations(group string) ([]db.Location, error) {
	return m.locations, nil
}

func (m *mockStore) CreateLocation(l db.Location) error {
	m.locations = append(m.locations, l)
	return nil
}

func (m *mockStore) GetLocation(group, id string) (db.Location, error) {
	for _, l := range m.locations {
		if l.Group == group && l.ID == id {
			return l, nil
		}
	}
	return db.Location{}, nil
}

func (m *mockStore) CreateVote(v db.Vote) error {
	m.votes = append(m.votes, v)
	return nil
}

func (m *mockStore) CalculateCrowdRate(group, id string, since time.Time) (db.CrowdRate, error) {
	total := 0
	sum := 0.0
	for _, v := range m.votes {
		if v.Group == group && v.LocationID == id && v.Timestamp.After(since) {
			total++
			sum += float64(v.Level)
		}
	}
	avg := 0.0
	if total > 0 {
		avg = sum / float64(total)
	}
	return db.CrowdRate{Group: group, LocationID: id, Votes: avg, Total: total, Since: since}, nil
}

func (m *mockStore) ListLocationIDs(group string) ([]string, error) {
	ids := []string{}
	for _, l := range m.locations {
		if l.Group == group {
			ids = append(ids, l.ID)
		}
	}
	return ids, nil
}

func TestService(t *testing.T) {
	store := &mockStore{}
	svc := NewCrowdVoteService(store)

	// Test Locations
	loc := db.Location{Group: "g1", ID: "l1", Name: "Loc 1"}
	svc.CreateLocation(loc)
	
	locs, _ := svc.ListLocations("g1")
	if len(locs) != 1 {
		t.Errorf("Expected 1 location, got %d", len(locs))
	}

	// Test Votes
	svc.SubmitVote(db.Vote{Group: "g1", LocationID: "l1", Level: 3})
	
	rate, _ := svc.GetCrowdRate("g1", "l1", 10, "minute")
	if rate.Total != 1 || rate.Votes != 3.0 {
		t.Errorf("Expected total 1 and votes 3.0, got total %d and votes %f", rate.Total, rate.Votes)
	}

	// Test GetLocation
	gotLoc, _ := svc.GetLocation("g1", "l1")
	if gotLoc.ID != "l1" {
		t.Errorf("Expected location l1, got %s", gotLoc.ID)
	}

	// Test List Crowd Rates
	rates, _ := svc.ListCrowdRates("g1", 10, "minute")
	if len(rates) != 1 {
		t.Errorf("Expected 1 crowd rate, got %d", len(rates))
	}

	// Test CalculateSince branches
	s1 := CalculateSince(1, "second")
	if time.Since(s1) < 1*time.Second {
		t.Log("Second branch hit")
	}
	s2 := CalculateSince(1, "hour")
	if time.Since(s2) > 59*time.Minute {
		t.Log("Hour branch hit")
	}
}
