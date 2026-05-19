// Package service handles business logic and orchestration
// between the API layer and the storage layer.
package service

import (
	"crowd-vote/db"
	"time"
)

// CrowdVoteService coordinates operations related to locations and votes.
type CrowdVoteService struct {
	store db.Store
}

// NewCrowdVoteService creates a new instance of CrowdVoteService.
func NewCrowdVoteService(store db.Store) *CrowdVoteService {
	return &CrowdVoteService{store: store}
}

// ListLocations retrieves all locations for a given group.
func (s *CrowdVoteService) ListLocations(group string) ([]db.Location, error) {
	return s.store.ListLocations(group)
}

// CreateLocation saves a new location.
func (s *CrowdVoteService) CreateLocation(l db.Location) error {
	return s.store.CreateLocation(l)
}

// GetLocation retrieves a specific location details.
func (s *CrowdVoteService) GetLocation(group, id string) (db.Location, error) {
	return s.store.GetLocation(group, id)
}

// SubmitVote processes a new user vote with the current timestamp.
func (s *CrowdVoteService) SubmitVote(v db.Vote) error {
	v.Timestamp = time.Now()
	return s.store.CreateVote(v)
}

// GetCrowdRate computes the aggregated status for a location.
func (s *CrowdVoteService) GetCrowdRate(group, id string, before int, unit string) (db.CrowdRate, error) {
	since := CalculateSince(before, unit)
	return s.store.CalculateCrowdRate(group, id, since)
}

// ListCrowdRates returns a list of crowd rates for all locations in a group.
func (s *CrowdVoteService) ListCrowdRates(group string, before int, unit string) ([]db.CrowdRate, error) {
	since := CalculateSince(before, unit)
	ids, err := s.store.ListLocationIDs(group)
	if err != nil {
		return nil, err
	}

	rates := []db.CrowdRate{}
	for _, id := range ids {
		rate, _ := s.store.CalculateCrowdRate(group, id, since)
		rates = append(rates, rate)
	}
	return rates, nil
}

// CalculateSince computes a start time for aggregation based on a relative time range.
func CalculateSince(before int, unit string) time.Time {
	var duration time.Duration
	switch unit {
	case "second":
		duration = time.Duration(before) * time.Second
	case "hour":
		duration = time.Duration(before) * time.Hour
	default: // minute
		duration = time.Duration(before) * time.Minute
	}
	return time.Now().Add(-duration)
}
