package main

import (
	"Alikhan/forum/pkg/forms"
	"Alikhan/forum/pkg/models"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w, r)
		return
	}
	files := []string{
		"./ui/html/home.html",
		"./ui/html/base.html",
		"./ui/html/footer.html",
	}

	switch r.Method {
	case http.MethodGet:
		s, err := app.posts.Latest()
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		app.render(w, r, files, "home.html", &templateData{Posts: s})
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			app.clientError(w, r, http.StatusBadRequest)
			return
		}

		var category string
		if r.Form.Get("Auto") != "" {
			category = r.Form.Get("Auto")
		} else if r.Form.Get("Bike") != "" {
			category = r.Form.Get("Bike")
		} else if r.Form.Get("Ship") != "" {
			category = r.Form.Get("Ship")
		} else if r.Form.Get("Track") != "" {
			category = r.Form.Get("Track")
		} else if r.Form.Get("Plane") != "" {
			category = r.Form.Get("Plane")
		}
		var s []*models.Post
		if category != "" {
			s, err = app.posts.GetByCategory(category)
			if err != nil {
				app.serverError(w, r, err)
				return
			}
		}

		app.render(w, r, files, "home.html", &templateData{Posts: s})
	default:
		app.clientError(w, r, http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) showPost(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/show.html",
		"./ui/html/base.html",
		"./ui/html/footer.html",
	}

	switch r.Method {
	case http.MethodGet:
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			app.notFound(w, r)
			return
		}
		s, err := app.posts.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.notFound(w, r)
			} else {
				app.serverError(w, r, err)
			}
			return
		}

		comm, err := app.comments.GetCommentByPost(app.likeToPostOrComment, app.dislikeToPostOrComment, id)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		likesNum, err := app.likeToPostOrComment.CountOfLikesPosts(id)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		dislikesNum, err := app.dislikeToPostOrComment.CountOfDislikesPosts(id)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		postAndComments := MergePostAndComments{
			Post:           s,
			Comments:       comm,
			LikesNumber:    likesNum,
			DislikesNumber: dislikesNum,
		}
		// data := &templateData{Post: s}
		app.render(w, r, files, "show.html", &templateData{Merge: postAndComments})

	case http.MethodPost:
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			app.notFound(w, r)
			return
		}
		r.ParseForm()

		comment := r.Form.Get("comment")
		if strings.TrimSpace(comment) == "" {
			app.clientError(w, r, http.StatusBadRequest)
			return
		}
		c, _ := r.Cookie("session")
		userId, _ := app.sessions.UserByToken(c.Value)
		userName, _ := app.users.UserNameByID(userId)
		_, err = app.comments.InsertComment(r.Form.Get("comment"), userName, id)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/post?id=%d", id), http.StatusSeeOther)
	default:
		app.clientError(w, r, http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) createPost(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/create.html",
		"./ui/html/base.html",
		"./ui/html/footer.html",
	}
	switch r.Method {
	case http.MethodGet:
		app.render(w, r, files, "create.html", &templateData{
			Form: forms.New(nil),
		})
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			app.clientError(w, r, http.StatusBadRequest)
			return
		}
		form := forms.New(r.PostForm)
		form.Required("title", "content")
		form.MaxLength("title", 100)

		if !form.Valid() {
			app.render(w, r, files, "create.html", &templateData{Form: form})
			return
		}
		token, err := r.Cookie("session")
		user_id, err := app.sessions.UserByToken(token.Value)
		id, err := app.posts.Insert(form.Get("title"), form.Get("content"), user_id, form.Get("category"))
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/post?id=%d", id), http.StatusSeeOther)
	default:
		// r.Header.Set("Allow", http.MethodPost)
		app.clientError(w, r, http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/signup.html",
		"./ui/html/base.html",
		"./ui/html/footer.html",
	}
	switch r.Method {
	case http.MethodGet:
		app.render(w, r, files, "signup.html", &templateData{Form: forms.New(nil)})
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			app.clientError(w, r, http.StatusBadRequest)
			return
		}

		form := forms.New(r.PostForm)
		form.Required("name", "email", "password")
		form.MaxLength("name", 255)
		form.MaxLength("email", 255)
		form.MathcesPattern("email", forms.EmailRX)
		form.MinLength("password", 10)

		if !form.Valid() {
			app.render(w, r, files, "signup.html", &templateData{Form: form})
			return
		}
		err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
		if err != nil {
			if errors.Is(err, models.ErrDuplicateEmail) {
				form.Errors.Add("email", "Address is already in use")
				app.render(w, r, files, "signup.html", &templateData{Form: form})
			} else if errors.Is(err, models.ErrDuplicateName) {
				form.Errors.Add("name", "Username is already in use")
				app.render(w, r, files, "signup.html", &templateData{Form: form})
			} else {
				app.serverError(w, r, err)
			}
			return
		}
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	default:
		app.clientError(w, r, http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/login.html",
		"./ui/html/base.html",
		"./ui/html/footer.html",
	}
	switch r.Method {
	case http.MethodGet:
		app.render(w, r, files, "login.html", &templateData{Form: forms.New(nil)})
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			app.clientError(w, r, http.StatusBadRequest)
			return
		}

		form := forms.New(r.Form)
		id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
		if err != nil {
			if errors.Is(err, models.ErrInvalidCredentials) {
				form.Errors.Add("generic", "Email or Password is incorrect")
				app.render(w, r, files, "login.html", &templateData{Form: form})
			} else {
				app.serverError(w, r, err)
			}
			return
		}

		err = app.sessions.DeleteSession(id)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		token, err := app.sessions.GenerateSession()
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		c := &http.Cookie{
			Name:     "session",
			Value:    token,
			MaxAge:   int(time.Hour) * 2,
			Path:     "/",
			HttpOnly: true,
		}
		http.SetCookie(w, c)
		// fmt.Println(time.Now().Add(time.Hour * 2))
		err = app.sessions.InsertSession(id, token, time.Now().Add(time.Hour*2))
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		http.Redirect(w, r, "/post/create", http.StatusSeeOther)
	default:
		app.clientError(w, r, http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.clientError(w, r, http.StatusMethodNotAllowed)
		return
	}
	c, err := r.Cookie("session")
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	c = &http.Cookie{
		Name:     "session",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, c)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) likes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.clientError(w, r, http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()
	c, _ := r.Cookie("session")
	userId, _ := app.sessions.UserByToken(c.Value)
	info := strings.Split(r.PostForm.Get("like"), " ")
	Postid, err := strconv.Atoi(info[1])
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	if info[0] == "like" {
		isLike, err := app.likeToPostOrComment.IsThereLikePost(Postid, userId)
		if isLike {
			err := app.likeToPostOrComment.DeleteLikePost(Postid, userId)
			if err != nil {
				app.serverError(w, r, err)
				return
			}
			http.Redirect(w, r, fmt.Sprintf("/post?id=%d", Postid), http.StatusSeeOther)
			return
		}
		if err != nil && err.Error() != "sql: no rows in result set" {
			app.serverError(w, r, err)
			return
		}
		isDislike, err := app.dislikeToPostOrComment.IsThereDisLikePost(Postid, userId)
		if err != nil && err.Error() != "sql: no rows in result set" {
			app.serverError(w, r, err)
			return
		}
		if isDislike {
			err := app.dislikeToPostOrComment.DeleteDislikePost(Postid, userId)
			if err != nil {
				app.serverError(w, r, err)
				return
			}
		}
		_, err = app.likeToPostOrComment.InsertToPost(userId, Postid)

		if err != nil {
			app.serverError(w, r, err)
			return
		}

	} else if info[0] == "dislike" {
		isDislike, err := app.dislikeToPostOrComment.IsThereDisLikePost(Postid, userId)
		if isDislike {
			err := app.dislikeToPostOrComment.DeleteDislikePost(Postid, userId)
			if err != nil {
				app.serverError(w, r, err)
				return
			}
			http.Redirect(w, r, fmt.Sprintf("/post?id=%d", Postid), http.StatusSeeOther)
			return
		}
		if err != nil && err.Error() != "sql: no rows in result set" {
			app.serverError(w, r, err)
			return
		}
		isLike, err := app.likeToPostOrComment.IsThereLikePost(Postid, userId)
		if err != nil && err.Error() != "sql: no rows in result set" {
			app.serverError(w, r, err)
			return
		}
		// fmt.Println(isLike)
		if isLike {
			err := app.likeToPostOrComment.DeleteLikePost(Postid, userId)
			if err != nil {
				app.serverError(w, r, err)
				return
			}
		}
		_, err = app.dislikeToPostOrComment.InsertToPostDis(userId, Postid)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	}
	// fmt.Println(r.PostForm.Get("like"))
	http.Redirect(w, r, fmt.Sprintf("/post?id=%d", Postid), http.StatusSeeOther)
}

func (app *application) commentLikes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.clientError(w, r, http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()
	c, _ := r.Cookie("session")
	userId, _ := app.sessions.UserByToken(c.Value)
	info := strings.Split(r.PostForm.Get("comments"), " ")
	commentId, err := strconv.Atoi(info[1])
	PostId, err := strconv.Atoi(info[2])
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	if info[0] == "like" {
		isLike, err := app.likeToPostOrComment.IsThereLikeComment(PostId, userId, commentId)
		if isLike {
			err := app.likeToPostOrComment.DeleteLikeComment(PostId, userId, commentId)
			if err != nil {
				app.serverError(w, r, err)
				return
			}
			http.Redirect(w, r, fmt.Sprintf("/post?id=%d", PostId), http.StatusSeeOther)
			return
		}
		if err != nil && err.Error() != "sql: no rows in result set" {
			app.serverError(w, r, err)
			return
		}
		isDislike, err := app.dislikeToPostOrComment.IsThereDisLikeComment(PostId, userId, commentId)
		if err != nil && err.Error() != "sql: no rows in result set" {
			app.serverError(w, r, err)
			return
		}
		if isDislike {
			err := app.dislikeToPostOrComment.DeleteDisLikeComment(PostId, userId, commentId)
			if err != nil {
				app.serverError(w, r, err)
				return
			}
		}
		_, err = app.likeToPostOrComment.InsertToComment(userId, PostId, commentId)

		if err != nil {
			app.serverError(w, r, err)
			return
		}

	} else if info[0] == "dislike" {

		isDislike, err := app.dislikeToPostOrComment.IsThereDisLikeComment(PostId, userId, commentId)
		if isDislike {
			err := app.dislikeToPostOrComment.DeleteDisLikeComment(PostId, userId, commentId)
			if err != nil {
				app.serverError(w, r, err)
				return
			}
			http.Redirect(w, r, fmt.Sprintf("/post?id=%d", PostId), http.StatusSeeOther)
			return
		}
		if err != nil && err.Error() != "sql: no rows in result set" {
			app.serverError(w, r, err)
			return
		}
		isLike, err := app.likeToPostOrComment.IsThereLikeComment(PostId, userId, commentId)
		if err != nil && err.Error() != "sql: no rows in result set" {
			app.serverError(w, r, err)
			return
		}
		// fmt.Println(isLike)
		if isLike {
			err := app.likeToPostOrComment.DeleteLikeComment(PostId, userId, commentId)
			if err != nil {
				app.serverError(w, r, err)
				return
			}
		}

		_, err = app.dislikeToPostOrComment.InsertToComment(userId, PostId, commentId)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

	}
	http.Redirect(w, r, fmt.Sprintf("/post?id=%d", PostId), http.StatusSeeOther)
}

func (app *application) myPosts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/myposts" {
		app.notFound(w, r)
		return
	}
	files := []string{
		"./ui/html/myposts.html",
		"./ui/html/base.html",
		"./ui/html/footer.html",
	}

	switch r.Method {
	case http.MethodGet:
		c, _ := r.Cookie("session")
		userId, _ := app.sessions.UserByToken(c.Value)
		s, err := app.posts.GetMyPosts(userId)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		app.render(w, r, files, "myposts.html", &templateData{Posts: s})
	default:
		app.clientError(w, r, http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) likedPosts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/likedPosts" {
		app.notFound(w, r)
		return
	}
	files := []string{
		"./ui/html/liked.html",
		"./ui/html/base.html",
		"./ui/html/footer.html",
	}

	switch r.Method {
	case http.MethodGet:
		c, _ := r.Cookie("session")
		userId, _ := app.sessions.UserByToken(c.Value)
		s, err := app.posts.GetLikedPosts(userId)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		app.render(w, r, files, "liked.html", &templateData{Posts: s})
	default:
		app.clientError(w, r, http.StatusMethodNotAllowed)
		return
	}
}
