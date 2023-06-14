package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alan1420/mobile-app-api/model/auth"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) AuthLoginUser(dataByte []byte) (string, int, error) {
	data := auth.UserLogin{}

	err := json.Unmarshal(dataByte, &data)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	if data.Email == "" || data.Password == "" {
		return "", http.StatusBadRequest, errors.New("email or password is empty")
	}

	// check in db
	query := `
		SELECT id, password FROM users
			WHERE email = ?
	;`

	row := r.db.QueryRow(query, data.Email)

	var password string
	var userId int

	err = row.Scan(&userId, &password)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", http.StatusUnauthorized, errors.New("wrong email or password")
		}
		return "", http.StatusInternalServerError, err
	}

	if !CheckPasswordHash(data.Password, password) {
		return "", http.StatusUnauthorized, errors.New("wrong email or password")
	}

	// Create a new random session token
	// we use the "github.com/google/uuid" library to generate UUIDs
	sessionId := uuid.NewString()

	// Save session in db
	query = `
		INSERT INTO sessions
			(id, user_id)
		VALUES
			(UUID_TO_BIN(?), ?)
	;`

	_, err = r.db.Exec(query, sessionId, userId)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	return sessionId, 200, nil
}

func (r *AuthRepository) AuthRegisterUser(dataByte []byte) (int, error) {
	data := auth.UserRegister{}

	err := json.Unmarshal(dataByte, &data)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if data.Email == "" || data.Password == "" || data.Name == "" {
		return http.StatusBadRequest, errors.New("email, password, or name is empty")
	}

	// count rows with the same email
	query := `
		SELECT COUNT(*) FROM users
			WHERE email = ?
	;`

	row := r.db.QueryRow(query, data.Email)

	var count int

	err = row.Scan(&count)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if count > 0 {
		return http.StatusBadRequest, errors.New("email already exists")
	}

	// hash password
	password, err := HashPassword(data.Password)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// insert into db
	query = `
		INSERT INTO users
			(email, password, name)
		VALUES
			(?, ?, ?)
	;`

	_, err = r.db.Exec(query, data.Email, password, data.Name)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return 200, nil
}

func (r *AuthRepository) verifySession(sessionId string) (*auth.UserData, int, error) {
	query := `
		SELECT user_id FROM sessions
			WHERE id = UUID_TO_BIN(?)
	;`

	row := r.db.QueryRow(query, sessionId)

	var userId int

	err := row.Scan(&userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, http.StatusUnauthorized, errors.New("session not found")
		}
		return nil, http.StatusInternalServerError, err
	}

	// find user data from user id
	query = `
		SELECT id, email, name FROM users
			WHERE id = ?
	;`

	row = r.db.QueryRow(query, userId)

	var userData auth.UserData

	err = row.Scan(
		&userData.ID,
		&userData.Email,
		&userData.Name,
	)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &userData, 200, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
