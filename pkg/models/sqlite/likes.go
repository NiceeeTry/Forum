package sqlite

import (
	"database/sql"
	"fmt"
)

type LikesModel struct {
	Db *sql.DB
}

func (m *LikesModel) InsertToPost(author_id, post_id int) (int, error) {
	stmt := `INSERT INTO likes (author_id, post_id)
VALUES(?, ?)`
	result, err := m.Db.Exec(stmt, author_id, post_id)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *LikesModel) CountOfLikesPosts(post_id int) (int, error) {
	var num int
	stmt := `SELECT COUNT(author_id) FROM likes WHERE post_id = ? AND comment_id IS NULL `
	row := m.Db.QueryRow(stmt, post_id)
	row.Scan(&num)
	// fmt.Println(num)
	return num, nil
}

func (m *LikesModel) IsThereLikePost(post_id, author_id int) (bool, error) {
	stmt := `SELECT id FROM likes WHERE author_id = ? AND post_id = ? AND comment_id IS NULL`
	row := m.Db.QueryRow(stmt, author_id, post_id)
	id := 0
	err := row.Scan(&id)
	if err != nil {
		return false, err
	}
	if id == 0 {
		return false, nil
	}
	return true, nil
}

func (m *LikesModel) DeleteLikePost(post_id, author_id int) error {
	query := `DELETE FROM likes WHERE author_id = ? AND post_id = ? AND comment_id IS NULL;`

	if _, err := m.Db.Exec(query, author_id, post_id); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (m *LikesModel) InsertToComment(author_id, post_id, comment_id int) (int, error) {
	stmt := `INSERT INTO likes (author_id, post_id, comment_id)
VALUES(?, ?, ?)`
	result, err := m.Db.Exec(stmt, author_id, post_id, comment_id)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *LikesModel) IsThereLikeComment(post_id, author_id, comment_id int) (bool, error) {
	stmt := `SELECT id FROM likes WHERE author_id = ? AND post_id = ? AND comment_id = ?`
	row := m.Db.QueryRow(stmt, author_id, post_id, comment_id)
	id := 0
	err := row.Scan(&id)
	if err != nil {
		return false, err
	}
	if id == 0 {
		return false, nil
	}
	return true, nil
}

func (m *LikesModel) DeleteLikeComment(post_id, author_id, comment_id int) error {
	query := `DELETE FROM likes WHERE author_id = ? AND post_id = ? AND comment_id = ?;`

	if _, err := m.Db.Exec(query, author_id, post_id, comment_id); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (m *LikesModel) CountOfLikesComment(post_id, comment_id int) (int, error) {
	var num int
	stmt := `SELECT COUNT(author_id) FROM likes WHERE post_id = ? AND comment_id = ? `
	row := m.Db.QueryRow(stmt, post_id, comment_id)
	row.Scan(&num)
	// fmt.Println(num)
	return num, nil
}
