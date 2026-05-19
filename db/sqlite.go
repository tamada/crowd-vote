package db

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStore struct {
	dbPath string
	db     *sql.DB
}

func NewSQLiteStore(dbPath string) *SQLiteStore {
	return &SQLiteStore{dbPath: dbPath}
}

func (s *SQLiteStore) Init() error {
	var err error
	s.db, err = sql.Open("sqlite3", s.dbPath)
	if err != nil {
		return err
	}

	schema := `
	CREATE TABLE IF NOT EXISTS locations (
		group_name TEXT,
		id TEXT,
		name TEXT NOT NULL,
		description TEXT,
		PRIMARY KEY (group_name, id)
	);
	CREATE TABLE IF NOT EXISTS votes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		group_name TEXT,
		location_id TEXT,
		client_address TEXT,
		level INTEGER,
		timestamp DATETIME,
		FOREIGN KEY(group_name, location_id) REFERENCES locations(group_name, id)
	);`
	_, err = s.db.Exec(schema)
	return err
}

func (s *SQLiteStore) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *SQLiteStore) ListLocations(group string) ([]Location, error) {
	rows, err := s.db.Query("SELECT group_name, id, name, description FROM locations WHERE group_name = ?", group)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	locations := []Location{}
	for rows.Next() {
		var l Location
		if err := rows.Scan(&l.Group, &l.ID, &l.Name, &l.Description); err != nil {
			return nil, err
		}
		locations = append(locations, l)
	}
	return locations, nil
}

func (s *SQLiteStore) CreateLocation(l Location) error {
	_, err := s.db.Exec("INSERT INTO locations (group_name, id, name, description) VALUES (?, ?, ?, ?)",
		l.Group, l.ID, l.Name, l.Description)
	return err
}

func (s *SQLiteStore) GetLocation(group, id string) (Location, error) {
	var l Location
	err := s.db.QueryRow("SELECT group_name, id, name, description FROM locations WHERE group_name = ? AND id = ?", group, id).
		Scan(&l.Group, &l.ID, &l.Name, &l.Description)
	return l, err
}

func (s *SQLiteStore) CreateVote(v Vote) error {
	_, err := s.db.Exec("INSERT INTO votes (group_name, location_id, client_address, level, timestamp) VALUES (?, ?, ?, ?, ?)",
		v.Group, v.LocationID, v.ClientAddress, v.Level, v.Timestamp)
	return err
}

func (s *SQLiteStore) CalculateCrowdRate(group, id string, since time.Time) (CrowdRate, error) {
	var votes sql.NullFloat64
	var total int
	err := s.db.QueryRow("SELECT AVG(level), COUNT(*) FROM votes WHERE group_name = ? AND location_id = ? AND timestamp > ?", group, id, since).
		Scan(&votes, &total)

	rate := CrowdRate{
		Group:      group,
		LocationID: id,
		Total:      total,
		Since:      since,
	}
	if votes.Valid {
		rate.Votes = votes.Float64
	}
	return rate, err
}

func (s *SQLiteStore) ListLocationIDs(group string) ([]string, error) {
	rows, err := s.db.Query("SELECT id FROM locations WHERE group_name = ?", group)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			continue
		}
		ids = append(ids, id)
	}
	return ids, nil
}
