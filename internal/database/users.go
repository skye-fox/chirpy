package database

import (
	"log"
)

type User struct {
	Id             int    `json:"id"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
}

func (db *DB) CreateUser(email, hashedPassword string) (User, error) {
	appDB, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	id := len(appDB.Users) + 1

	newUser := User{
		Id:             id,
		Email:          email,
		HashedPassword: hashedPassword,
	}
	appDB.Users[id] = newUser

	err = db.writeDB(appDB)
	if err != nil {
		return User{}, err
	}

	return newUser, nil
}

func (db *DB) UpdateUsers(id int, email, hashedPassword string) (User, error) {
	appDB, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := appDB.Users[id]
	if !ok {
		return User{}, ErrNotExist
	}

	user.HashedPassword = hashedPassword
	user.Email = email
	appDB.Users[id] = user

	err = db.writeDB(appDB)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) CheckDuplicateEmail(email string) bool {
	appDB, err := db.loadDB()
	if err != nil {
		log.Printf("Error loading db: %s", err)
	}

	for _, user := range appDB.Users {
		if user.Email == email {
			return true
		}
	}

	return false
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	appDB, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range appDB.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, ErrNotExist
}

func (db *DB) GetUserByID(id int) (User, error) {
	appDB, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := appDB.Users[id]
	if !ok {
		return User{}, ErrNotExist
	}

	return user, nil
}
