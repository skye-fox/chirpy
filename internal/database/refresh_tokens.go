package database

import (
	"time"
)

type RefreshToken struct {
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (db *DB) SaveRefreshToken(token string, id int) error {
	appDB, err := db.loadDB()
	if err != nil {
		return err
	}

	expiration := 60 * 24 * time.Hour
	refreshToken := RefreshToken{
		UserID:    id,
		Token:     token,
		ExpiresAt: time.Now().Add(expiration),
	}
	appDB.RefreshTokens[token] = refreshToken

	err = db.writeDB(appDB)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetUserByRefreshToken(token string) (User, error) {
	appDB, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	refreshToken, ok := appDB.RefreshTokens[token]
	if !ok {
		return User{}, ErrNotExist
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return User{}, ErrNotExist
	}

	user, err := db.GetUserByID(refreshToken.UserID)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) RevokeToken(token string) error {
	appDB, err := db.loadDB()
	if err != nil {
		return err
	}

	delete(appDB.RefreshTokens, token)

	err = db.writeDB(appDB)
	if err != nil {
		return err
	}

	return nil
}
