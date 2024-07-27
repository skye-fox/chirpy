package database

import "errors"

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	appDB, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(appDB.Chirps) + 1

	newChirp := Chirp{
		Id:   id,
		Body: body,
	}
	appDB.Chirps[id] = newChirp

	err = db.writeDB(appDB)
	if err != nil {
		return Chirp{}, err
	}

	return newChirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	appDB, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(appDB.Chirps))
	for _, chirp := range appDB.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) GetChirpByID(id int) (Chirp, error) {
	appDB, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := appDB.Chirps[id]
	if !ok {
		return Chirp{}, errors.New("Invalid chirp ID")
	}

	return chirp, nil
}