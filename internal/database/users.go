package database

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

func (db *DB) CreateUser(email string) (User, error) {
	appDB, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	id := len(appDB.Users) + 1

	newUser := User{
		Id:    id,
		Email: email,
	}
	appDB.Users[id] = newUser

	err = db.writeDB(appDB)
	if err != nil {
		return User{}, err
	}

	return newUser, nil
}
