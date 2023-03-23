package main

import (
	"bytes"
	"fmt"
	"net/http"
	"text/template"
	"time"
)

var files = []string{
	"./ui/html/error.html",
	"./ui/html/base.html",
	"./ui/html/footer.html",
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	// trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	trace := fmt.Sprintf("%s", err.Error())
	app.errorLog.Output(2, trace)

	app.render(w, r, files, "error.html", &templateData{ErrorM: ErrorMes{http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError}})
	// http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, r *http.Request, status int) {
	app.render(w, r, files, "error.html", &templateData{ErrorM: ErrorMes{http.StatusText(status), status}})
	// http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request) {
	app.clientError(w, r, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name []string, nameFunc string, td *templateData) {
	ts, err := template.New(nameFunc).Funcs(functions).ParseFiles(name...)
	if err != nil {
		app.serverError(w, r, fmt.Errorf("The template %s does not exist", name))
		return
	}
	if nameFunc == "error.html" {
		w.WriteHeader(td.ErrorM.ErrorCode)
		ts.Execute(w, app.addDefaultData(td, r))
		return
	}
	buf := new(bytes.Buffer)
	err = ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	buf.WriteTo(w)
}

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CurrentYear = time.Now().Year()
	td.IsAuthenticated = app.isAuthenticated(r)
	if td.IsAuthenticated {
		c, _ := r.Cookie("session")
		userId, _ := app.sessions.UserByToken(c.Value)
		userName, _ := app.users.UserNameByID(userId)
		td.UserName = userName
	}

	return td
}

func (app *application) isAuthenticated(r *http.Request) bool {
	c, err := r.Cookie("session")
	// fmt.Println(c.Value)
	if err != nil {
		// log.Println(err.Error())
		return false
	}
	id, err := app.sessions.UserByToken(c.Value)
	if err != nil {
		// log.Println(err.Error())
		return false
	}
	if id == 0 {
		return false
	}
	// fmt.Println(token.Value)
	return true
}

// func (app *application) LikesAndDislikes(userId int, r *http.Request, w http.ResponseWriter) error {
// 	info := strings.Split(r.PostForm.Get("like"), " ")
// 	Postid, err := strconv.Atoi(info[1])
// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}
// 	if info[0] == "like" {
// 		isLike, err := app.likeToPostOrComment.IsThereLikePost(Postid, userId)
// 		if isLike {
// 			return err
// 		}
// 		if err != nil && err.Error() != "sql: no rows in result set" {
// 			fmt.Println(err)
// 			return err
// 		}
// 		isDislike, err := app.dislikeToPostOrComment.IsThereDisLikePost(Postid, userId)
// 		if err != nil && err.Error() != "sql: no rows in result set" {
// 			fmt.Println(err)
// 			fmt.Print("1st errr")
// 			return err
// 		}
// 		if isDislike {
// 			err := app.dislikeToPostOrComment.DeleteDislikePost(Postid, userId)
// 			if err != nil {
// 				fmt.Println(err)
// 				return err
// 			}
// 		}
// 		_, err = app.likeToPostOrComment.InsertToPost(userId, Postid)

// 		if err != nil {
// 			fmt.Println(err)
// 			app.serverError(w, err)
// 			return err
// 		}

// 	} else if info[0] == "dislike" {
// 		isDislike, err := app.dislikeToPostOrComment.IsThereDisLikePost(Postid, userId)
// 		if isDislike {
// 			return err
// 		}
// 		if err != nil && err.Error() != "sql: no rows in result set" {
// 			fmt.Println(err)
// 			fmt.Print("1st errr")
// 			return err
// 		}
// 		isLike, err := app.likeToPostOrComment.IsThereLikePost(Postid, userId)
// 		if err != nil && err.Error() != "sql: no rows in result set" {
// 			fmt.Println(err)
// 			return err
// 		}
// 		fmt.Println(isLike)
// 		if isLike {
// 			err := app.likeToPostOrComment.DeleteLikePost(Postid, userId)
// 			if err != nil {
// 				fmt.Println(err)
// 				return err
// 			}
// 		}
// 		_, err = app.dislikeToPostOrComment.InsertToPostDis(userId, Postid)
// 		if err != nil {
// 			app.serverError(w, err)
// 			return err
// 		}
// 	}
// 	return nil
// 	// fmt.Println(r.PostForm.Get("like"))
// }
