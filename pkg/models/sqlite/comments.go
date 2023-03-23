package sqlite

import (
	"Alikhan/forum/pkg/models"
	"database/sql"
)

type CommentModel struct {
	DB *sql.DB
}

func (m *CommentModel) InsertComment(comment, author string, post_id int) (int, error) {
	stmt := `INSERT INTO comments (comment, author, post_id, created)
VALUES(?, ?, ?, DATETIME('now'))`
	result, err := m.DB.Exec(stmt, comment, author, post_id)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *CommentModel) GetCommentByPost(likeMod *LikesModel, dislike *DisLikesModel, post_id int) ([]*models.Comment, error) {
	stmt := `SELECT comment, author, created, id FROM comments WHERE post_id = ?`
	rows, err := m.DB.Query(stmt, post_id)
	if err != nil {
		// fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	comments := []*models.Comment{}
	for rows.Next() {
		s := &models.Comment{}
		err = rows.Scan(&s.CommentText, &s.Author, &s.Created, &s.ID)
		if err != nil {
			return nil, err
		}
		s.PostId = post_id

		likesNumComment, err := likeMod.CountOfLikesComment(post_id, s.ID)
		if err != nil {
			return nil, err
		}
		s.LikesComment = likesNumComment

		dislikesNumComment, err := dislike.CountOfDisLikesComment(post_id, s.ID)
		if err != nil {
			return nil, err
		}
		s.DislikesComment = dislikesNumComment

		comments = append(comments, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}
