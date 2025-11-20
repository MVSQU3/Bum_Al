package repo

import (
	"database/sql"
	"fmt"
	"xxx/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserController(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) SignIn(input *models.Input) (*models.Input, error) {
	query := "SELECT email, password FROM users WHERE email=$1"
	err := r.db.QueryRow(query, input.Email).Scan(&input.Email, input.Password)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("User not found")
	}
	if err != nil {
		return nil, err
	}
	return input, nil
}

func (r *UserRepository) SignUp(user *models.User) (*models.User, error) {
	CheckQuery := "SELECT email FROM users WHERE email = $1"
	var existingEmail string

	err := r.db.QueryRow(CheckQuery, user.Email).Scan(&existingEmail)
	if err == nil {
		return nil, fmt.Errorf("email already exists")
	}
	if err != sql.ErrNoRows {
		// Erreur autre que "non trouvé"
		return nil, err
	}

	InsertQuery := "INSERT INTO users (fullname, email, password) VALUES (1$, 2$, 3$) RETURNING id, fullname, email"
	InsertErr := r.db.QueryRow(InsertQuery, user.FullName, user.Email, user.Password).Scan(&user.ID, &user.FullName, &user.Email)
	if InsertErr != nil {
		return nil, InsertErr
	}

	return user, nil
}
