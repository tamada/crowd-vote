package main

import (
	"bytes"
	"crowd-vote/db"
	"crowd-vote/service"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHandlers(t *testing.T) {
	dbPath := "test_handlers.db"
	defer os.Remove(dbPath)

	store := db.NewSQLiteStore(dbPath)
	store.Init()
	defer store.Close()

	svc = service.NewCrowdVoteService(store)

	// Test POST /locations/{group}
	loc := db.Location{ID: "l1", Name: "Loc 1"}
	body, _ := json.Marshal(loc)
	req := httptest.NewRequest("POST", "/locations/g1", bytes.NewBuffer(body))
	req.SetPathValue("group", "g1")
	w := httptest.NewRecorder()
	createLocation(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected 201, got %d", w.Code)
	}

	// Test GET /locations/{group}
	req = httptest.NewRequest("GET", "/locations/g1", nil)
	req.SetPathValue("group", "g1")
	w = httptest.NewRecorder()
	listLocations(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	// Test GET /locations/{group}/{id}
	req = httptest.NewRequest("GET", "/locations/g1/l1", nil)
	req.SetPathValue("group", "g1")
	req.SetPathValue("id", "l1")
	w = httptest.NewRecorder()
	getLocation(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	// Test POST /votes/{group}/{id}
	vote := db.Vote{Level: 2}
	body, _ = json.Marshal(vote)
	req = httptest.NewRequest("POST", "/votes/g1/l1", bytes.NewBuffer(body))
	req.SetPathValue("group", "g1")
	req.SetPathValue("id", "l1")
	w = httptest.NewRecorder()
	createVote(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected 201, got %d", w.Code)
	}

	// Test GET /votes/{group}/{id}
	req = httptest.NewRequest("GET", "/votes/g1/l1?before=5&unit=hour", nil)
	req.SetPathValue("group", "g1")
	req.SetPathValue("id", "l1")
	w = httptest.NewRecorder()
	getCrowdRate(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	// Test GET /votes/{group}
	req = httptest.NewRequest("GET", "/votes/g1", nil)
	req.SetPathValue("group", "g1")
	w = httptest.NewRecorder()
	listCrowdRates(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	// Test GET /histories/{group}/{id}
	req = httptest.NewRequest("GET", "/histories/g1/l1", nil)
	req.SetPathValue("group", "g1")
	req.SetPathValue("id", "l1")
	w = httptest.NewRecorder()
	getHistory(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}
