package sqlite

import (
	"database/sql"
	"fmt"
)

type DisLikesModel struct {
	Db *sql.DB
}

func (m *DisLikesModel) InsertToPostDis(author_id, post_id int) (int, error) {
	stmt := `INSERT INTO dislikes (author_id, post_id)
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

func (m *DisLikesModel) IsThereDisLikePost(post_id, author_id int) (bool, error) {
	stmt := `SELECT id FROM dislikes WHERE author_id = ? AND post_id = ? AND comment_id IS NULL`
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

func (m *DisLikesModel) DeleteDislikePost(post_id, author_id int) error {
	query := `DELETE FROM dislikes WHERE author_id = ? AND post_id = ? AND comment_id IS NULL;`

	if _, err := m.Db.Exec(query, author_id, post_id); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (m *DisLikesModel) CountOfDislikesPosts(post_id int) (int, error) {
	var num int
	stmt := `SELECT COUNT(author_id) FROM dislikes WHERE post_id = ? AND comment_id IS NULL`
	row := m.Db.QueryRow(stmt, post_id)
	row.Scan(&num)
	// fmt.Println(num)
	return num, nil
}

func (m *DisLikesModel) InsertToComment(author_id, post_id, comment_id int) (int, error) {
	stmt := `INSERT INTO dislikes (author_id, post_id, comment_id)
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

func (m *DisLikesModel) IsThereDisLikeComment(post_id, author_id, comment_id int) (bool, error) {
	stmt := `SELECT id FROM dislikes WHERE author_id = ? AND post_id = ? AND comment_id = ?`
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

func (m *DisLikesModel) DeleteDisLikeComment(post_id, author_id, comment_id int) error {
	query := `DELETE FROM dislikes WHERE author_id = ? AND post_id = ? AND comment_id = ?;`

	if _, err := m.Db.Exec(query, author_id, post_id, comment_id); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (m *DisLikesModel) CountOfDisLikesComment(post_id, comment_id int) (int, error) {
	var num int
	stmt := `SELECT COUNT(author_id) FROM dislikes WHERE post_id = ? AND comment_id = ? `
	row := m.Db.QueryRow(stmt, post_id, comment_id)
	row.Scan(&num)
	// fmt.Println(num)
	return num, nil
}
