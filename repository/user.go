package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/alan1420/mobile-app-api/model/auth"
	"github.com/alan1420/mobile-app-api/model/user"
	. "github.com/alan1420/mobile-app-api/model/user"

	val "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type UserRepository struct {
	db       *sql.DB
	authRepo *AuthRepository
}

func NewUserRepository(
	db *sql.DB,
	authRepo *AuthRepository,
) *UserRepository {
	return &UserRepository{
		db:       db,
		authRepo: authRepo,
	}
}

func (r *UserRepository) GetUser(session string) (*auth.UserData, int, error) {
	return r.authRepo.verifySession(session)
}

func (r *UserRepository) GetTicket(userID int) (*UserTicket, error) {
	query := `
		SELECT 
			depart_date, depart_time, travel_time, seat, from_location, to_location, qrcode 
		FROM 
			tickets 
		WHERE 
			user_id = ?;
	`

	row := r.db.QueryRow(query, userID)
	ticket := UserTicket{}

	err := row.Scan(
		&ticket.TicketDate,
		&ticket.DepartureOn,
		&ticket.TravelTime,
		&ticket.Seat,
		&ticket.From,
		&ticket.To,
		&ticket.Qrcode,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &ticket, nil
}

func (r *UserRepository) CreateTicket(userID int, data []byte) (int, error) {
	valid := val.New()
	ticket := user.CreateUserTicket{}

	err := json.Unmarshal(data, &ticket)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if err := valid.Struct(ticket); err != nil {
		return http.StatusBadRequest, errors.New("Invalid request body")
	}

	// generate qrcode
	// for now use uuid
	qrcode := uuid.NewString()

	query := `
		INSERT INTO
			tickets
			(user_id, depart_date, depart_time, travel_time, seat, from_location, to_location, qrcode)
		VALUES
			(?, ?, ?, ?, ?, ?, ?, ?);
	`

	_, err = r.db.Exec(
		query,
		userID,
		ticket.TicketDate,
		ticket.DepartureOn,
		ticket.TravelTime,
		ticket.Seat,
		ticket.From,
		ticket.To,
		qrcode,
	)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
