package sqlite

import (
	"Alikhan/forum/pkg/models"
	"database/sql"
	"errors"
)

type PostModel struct {
	Db *sql.DB
}

func (m *PostModel) Insert(title, content string, user_id int, category string) (int, error) {
	stmt := `INSERT INTO posts (title, content, created, user_id, category)
VALUES(?, ?, DATETIME('now'),?,?)`
	result, err := m.Db.Exec(stmt, title, content, user_id, category)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *PostModel) Get(id int) (*models.Post, error) {
	stmt := `SELECT posts.id, title, content, posts.created, name FROM posts JOIN users ON posts.user_id = users.id WHERE posts.id = ?`
	row := m.Db.QueryRow(stmt, id)
	s := &models.Post{}
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *PostModel) GetByCategory(category string) ([]*models.Post, error) {
	stmt := `SELECT posts.id, title, content, posts.created, name FROM posts JOIN users ON posts.user_id = users.id WHERE category = ?`
	rows, err := m.Db.Query(stmt, category)
	if err != nil {
		// fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	posts := []*models.Post{}
	for rows.Next() {
		s := &models.Post{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Name)
		if err != nil {
			return nil, err
		}
		posts = append(posts, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (m *PostModel) Latest() ([]*models.Post, error) {
	stmt := `SELECT posts.id, title, content, posts.created, name FROM posts JOIN users ON posts.user_id = users.id ORDER BY posts.created DESC`
	rows, err := m.Db.Query(stmt)
	if err != nil {
		// fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	posts := []*models.Post{}
	for rows.Next() {
		s := &models.Post{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Name)
		if err != nil {
			return nil, err
		}
		posts = append(posts, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (m *PostModel) GetMyPosts(id int) ([]*models.Post, error) {
	stmt := `SELECT posts.id, title, content, posts.created, name FROM posts JOIN users ON posts.user_id = users.id WHERE user_id = ?`
	rows, err := m.Db.Query(stmt, id)
	if err != nil {
		// fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	posts := []*models.Post{}
	for rows.Next() {
		s := &models.Post{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Name)
		if err != nil {
			return nil, err
		}
		posts = append(posts, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (m *PostModel) GetLikedPosts(id int) ([]*models.Post, error) {
	stmt := `SELECT posts.id, title, content, posts.created, name FROM posts JOIN users ON posts.user_id = users.id WHERE posts.id IN (SELECT post_id FROM likes WHERE author_id = ? AND comment_id IS NULL)`
	rows, err := m.Db.Query(stmt, id)
	if err != nil {
		// fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	posts := []*models.Post{}
	for rows.Next() {
		s := &models.Post{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Name)
		if err != nil {
			return nil, err
		}
		posts = append(posts, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}
