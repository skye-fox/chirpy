package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type DB struct {
	path string
	mu   *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

func NewDB(path string) (*DB, error) {
	newDB := DB{
		path: path,
		mu:   &sync.RWMutex{},
	}
	err := newDB.ensureDB()
	return &newDB, err
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	chirpDB, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(chirpDB.Chirps) + 1

	newChirp := Chirp{
		Id:   id,
		Body: body,
	}
	chirpDB.Chirps[newChirp.Id] = newChirp

	err = db.writeDB(chirpDB)
	if err != nil {
		return Chirp{}, err
	}

	return newChirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	chirpDB, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(chirpDB.Chirps))
	for _, chirp := range chirpDB.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) GetChirpByID(id int) (Chirp, error) {
	chirpDB, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := chirpDB.Chirps[id]
	if !ok {
		return Chirp{}, errors.New("Invalid chirp ID")
	}

	return chirp, nil
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.writeDB(DBStructure{
			Chirps: map[int]Chirp{},
		})
	}
	return err
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	newDBStructure := DBStructure{}
	file, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return newDBStructure, err
	}

	err = json.Unmarshal(file, &newDBStructure)
	if err != nil {
		return newDBStructure, err
	}
	return newDBStructure, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, data, 0600)
	if err != nil {
		return err
	}
	return nil
}
