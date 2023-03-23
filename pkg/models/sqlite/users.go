package sqlite

import (
	"Alikhan/forum/pkg/models"
	"database/sql"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (name, email, hashed_password, created)
	VALUES (?,?,?,DATETIME('now'))`
	_, err = m.DB.Exec(stmt, name, email, string(hasedPassword))
	// fmt.Println(err)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.email") {
			return models.ErrDuplicateEmail
		}
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.name") {
			return models.ErrDuplicateName
		}
		// fmt.Println(err)
		return err
	}
	return nil
}

func (m *UserModel) UserNameByID(id int) (string, error) {
	var name string
	stmt := `
		SELECT name
		FROM users
		WHERE id = ?;
	`
	if err := m.DB.QueryRow(stmt, id).Scan(&name); err != nil {
		return "", err
	}
	return name, nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hasedPassword []byte
	stmt := `SELECT id, hashed_password FROM users WHERE email = ? AND active = TRUE`
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hasedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hasedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return id, nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
