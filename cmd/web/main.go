package main

import (
	"Alikhan/forum/pkg/models/sqlite"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	posts    *sqlite.PostModel
	// session  []http.Cookie
	users                  *sqlite.UserModel
	sessions               *sqlite.SessionModel
	comments               *sqlite.CommentModel
	likeToPostOrComment    *sqlite.LikesModel
	dislikeToPostOrComment *sqlite.DisLikesModel
}

func main() {
	addr := flag.String("addr", ":8080", "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB()
	if err != nil {
		errorLog.Fatal(err)
	}
	stmt, _ := db.Prepare(`CREATE TABLE IF NOT EXISTS posts (id INTEGER PRIMARY KEY, title TEXT NOT NULL,
	content TEXT NOT NULL,
	created DATETIME NOT NULL,
	user_id INTEGER,
	category TEXT);`)
	stmt.Exec()
	stmt, _ = db.Prepare(`CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		hashed_password TEXT NOT NULL,
		created DATETIME NOT NULL,
		active BOOLEAN NOT NULL DEFAULT TRUE,
		CONSTRAINT users_uc_email UNIQUE (email),
		CONSTRAINT users_uc_name UNIQUE (name));`)
	stmt.Exec()
	stmt, _ = db.Prepare(`CREATE TABLE IF NOT EXISTS sessions (id INTEGER PRIMARY KEY,
		user_id INTEGER NOT NULL,
		session TEXT NOT NULL,
		expire DATETIME);`)
	stmt.Exec()
	stmt, _ = db.Prepare(`CREATE TABLE IF NOT EXISTS comments (id INTEGER PRIMARY KEY,
		comment TEXT NOT NULL,
		author TEXT NOT NULL,
		post_id INTEGER,
		created DATETIME NOT NULL);`)
	stmt.Exec()
	stmt, _ = db.Prepare(`CREATE TABLE IF NOT EXISTS likes (id INTEGER PRIMARY KEY,
		author_id INTEGER NOT NULL,
		post_id INTEGER NOT NULL,
		comment_id INTEGER);`)
	stmt.Exec()
	stmt, _ = db.Prepare(`CREATE TABLE IF NOT EXISTS dislikes (id INTEGER PRIMARY KEY,
		author_id INTEGER NOT NULL,
		post_id INTEGER NOT NULL,
		comment_id INTEGER);`)
	stmt.Exec()
	defer db.Close()

	// session := &http.Cookie{
	// 	Name: "session",
	// 	Value: string(uuid.NewV4()),
	// }
	// var session []http.Cookie
	app := &application{
		errorLog:               errorLog,
		infoLog:                infoLog,
		posts:                  &sqlite.PostModel{Db: db},
		users:                  &sqlite.UserModel{DB: db},
		sessions:               &sqlite.SessionModel{DB: db},
		comments:               &sqlite.CommentModel{DB: db},
		likeToPostOrComment:    &sqlite.LikesModel{Db: db},
		dislikeToPostOrComment: &sqlite.DisLikesModel{Db: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	infoLog.Printf("Starting server on http://localhost%s/", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
