package main

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gitlab.com/fqazi/snippet-box/internal/models"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, req *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(req)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) snippetView(w http.ResponseWriter, req *http.Request) {
	params := httprouter.ParamsFromContext(req.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(req)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.tmpl.html", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Display form for creating a new snippet ..."))
}

func (app *application) snippetCreatePost(w http.ResponseWriter, req *http.Request) {
	title := "o snail"
	content := "O snail\bClimb mount Fuji,\nBut slowly, slowly,!\n\n- Kobayahi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
	}

	http.Redirect(w, req, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
