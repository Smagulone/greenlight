package main

import (
	"fmt"
	"github.com/dexciuq/greenlight/internal/data"
	"net/http"
)

func (app *application) createDirectorHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name    string   `json:"name"`
		Surname string   `json:"surname"`
		Awards  []string `json:"awards"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
	}

	direction := &data.Director{
		Name:    input.Name,
		Surname: input.Surname,
		Awards:  input.Awards,
	}

	err = app.models.Directors.Insert(direction)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/directors/%d", direction.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"director": direction}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listDirectorsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name    string
		Surname string
		Awards  []string
		data.Filters
	}

	qs := r.URL.Query()
	input.Name = app.readString(qs, "name", "")
	input.Surname = app.readString(qs, "surname", "")
	input.Awards = app.readCSV(qs, "awards", []string{})
	input.Filters.Page = app.readInt(qs, "page", 1)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "name", "surname", "awards", "-id", "-name", "-surname", "-awards"}

	directors, metadata, err := app.models.Directors.GetAll(input.Name, input.Surname, input.Awards, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"directors": directors, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
