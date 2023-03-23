package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

type SessionModel struct {
	DB *sql.DB
}

// func (m *SessionModel) IsThereSession(user_id int, session string) error {
// 	var trueUserId int
// 	var trueSession string
// 	stmt := `SELECT user_id, session FROM sessions WHERE user_id = ?`
// 	row := m.DB.QueryRow(stmt, user_id)
// 	err := row.Scan(&trueUserId, &trueSession)
// 	if err != nil {
// 		err = m.InsertSession(user_id, session)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	} else if trueSession != session {
// 		stmt = `UPDATE sessions SET session = ? WHERE user_id = ?`
// 		_, err = m.DB.Exec(stmt, session, user_id)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

func (m *SessionModel) GetSession(id int) (string, error) {
	var trueSession string
	stmt := `SELECT session FROM sessions WHERE user_id = ? AND expire >= DATETIME('now')`
	row := m.DB.QueryRow(stmt, id)
	err := row.Scan(&trueSession)
	if err != nil {
		return "", err
	}
	return trueSession, nil
}

func (m *SessionModel) DeleteSession(id int) error {
	query := `DELETE FROM sessions WHERE user_id = ?;`

	if _, err := m.DB.Exec(query, id); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (m *SessionModel) GenerateSession() (string, error) {
	token, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return token.String(), nil
}

func (m *SessionModel) UserByToken(token string) (int, error) {
	var id int
	stmt := `
		SELECT user_id
		FROM sessions
		WHERE session = ?;
	`
	if err := m.DB.QueryRow(stmt, token).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (m *SessionModel) InsertSession(id int, session string, expire time.Time) error {
	stmt := `INSERT INTO sessions (user_id, session, expire)
	VALUES (?,?,?)`
	_, err := m.DB.Exec(stmt, id, session, expire)
	if err != nil {
		return err
	}
	return nil
}
