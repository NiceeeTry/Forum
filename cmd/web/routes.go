package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/post", app.showPost)
	mux.HandleFunc("/post/create", app.requireAuthentication(app.createPost))
	mux.HandleFunc("/post/like", app.requireAuthentication(app.likes))
	mux.HandleFunc("/post/comment", app.requireAuthentication(app.commentLikes))
	mux.HandleFunc("/myposts", app.requireAuthentication(app.myPosts))
	mux.HandleFunc("/likedPosts", app.requireAuthentication(app.likedPosts))

	mux.HandleFunc("/user/signup", app.signup)
	mux.HandleFunc("/user/login", app.login)
	mux.HandleFunc("/user/logout", app.requireAuthentication(app.logout))

	fs := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fs))
	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
