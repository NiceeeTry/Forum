package main

import (
	"Alikhan/forum/pkg/forms"
	"Alikhan/forum/pkg/models"
	"text/template"
	"time"
)

type templateData struct {
	// Post            *models.Post
	Posts           []*models.Post
	CurrentYear     int
	Form            *forms.Form
	IsAuthenticated bool
	UserName        string
	// Comments        []*models.Comment
	Merge  MergePostAndComments
	ErrorM ErrorMes
}

type ErrorMes struct {
	ErrorMessage string
	ErrorCode    int
}

type MergePostAndComments struct {
	Comments       []*models.Comment
	Post           *models.Post
	LikesNumber    int
	DislikesNumber int
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}
