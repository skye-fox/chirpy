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
	Users  map[int]User  `json:"users"`
}

func NewDB(path string) (*DB, error) {
	newDB := DB{
		path: path,
		mu:   &sync.RWMutex{},
	}
	err := newDB.ensureDB()
	return &newDB, err
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.writeDB(DBStructure{
			Chirps: map[int]Chirp{},
			Users:  map[int]User{},
		})
	}
	return err
}

func (db *DB) ResetDB() error {
	err := os.Remove(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	return db.ensureDB()
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
