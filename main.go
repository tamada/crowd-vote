package main

import (
	"crowd-vote/db"
	"crowd-vote/service"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
)

var svc *service.CrowdVoteService

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./crowd_vote.db"
	}

	// Setup Store
	store := db.NewSQLiteStore(dbPath)
	if err := store.Init(); err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	// Setup Service (Business Logic)
	svc = service.NewCrowdVoteService(store)

	// Setup Routes (REST)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /locations/{group}", listLocations)
	mux.HandleFunc("POST /locations/{group}", createLocation)
	mux.HandleFunc("GET /locations/{group}/{id}", getLocation)
	mux.HandleFunc("POST /votes/{group}/{id}", createVote)
	mux.HandleFunc("GET /votes/{group}/{id}", getCrowdRate)
	mux.HandleFunc("GET /votes/{group}", listCrowdRates)
	mux.HandleFunc("GET /histories/{group}/{id}", getHistory)

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// REST Handlers

func listLocations(w http.ResponseWriter, r *http.Request) {
	group := r.PathValue("group")
	locations, err := svc.ListLocations(group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(locations)
}

func createLocation(w http.ResponseWriter, r *http.Request) {
	group := r.PathValue("group")
	var l db.Location
	if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	l.Group = group
	if err := svc.CreateLocation(l); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(l)
}

func getLocation(w http.ResponseWriter, r *http.Request) {
	group := r.PathValue("group")
	id := r.PathValue("id")
	l, err := svc.GetLocation(group, id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(l)
}

func createVote(w http.ResponseWriter, r *http.Request) {
	group := r.PathValue("group")
	id := r.PathValue("id")
	var v db.Vote
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	v.Group = group
	v.LocationID = id
	if v.ClientAddress == "" {
		v.ClientAddress = r.RemoteAddr
	}

	if err := svc.SubmitVote(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(v)
}

func getCrowdRate(w http.ResponseWriter, r *http.Request) {
	group := r.PathValue("group")
	id := r.PathValue("id")
	before, unit := parseQueryParams(r)

	rate, err := svc.GetCrowdRate(group, id, before, unit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(rate)
}

func listCrowdRates(w http.ResponseWriter, r *http.Request) {
	group := r.PathValue("group")
	before, unit := parseQueryParams(r)

	rates, err := svc.ListCrowdRates(group, before, unit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(rates)
}

func getHistory(w http.ResponseWriter, r *http.Request) {
	group := r.PathValue("group")
	id := r.PathValue("id")
	before, unit := parseQueryParams(r)

	rate, err := svc.GetCrowdRate(group, id, before, unit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode([]db.CrowdRate{rate})
}

// REST Helpers

func parseQueryParams(r *http.Request) (int, string) {
	beforeStr := r.URL.Query().Get("before")
	unit := r.URL.Query().Get("unit")

	before, err := strconv.Atoi(beforeStr)
	if err != nil {
		before = 10
		unit = "minute"
	}
	if unit == "" {
		unit = "minute"
	}
	return before, unit
}
