package transport

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"server/pkg/app"
)

type Handler interface {
	Index(w http.ResponseWriter, r *http.Request)
	About(w http.ResponseWriter, r *http.Request)
	Summary(w http.ResponseWriter, r *http.Request)
}

type Response struct {
	Rank       float64 `json:"rank"`
	Similarity int     `json:"similarity"`
}

type SummaryData struct {
	Text       string
	Rank       float64
	Similarity int
}

type handler struct {
	ctx             context.Context
	valuatorService app.ValuatorService
}

func NewHandler(
	ctx context.Context,
	valuatorService app.ValuatorService,
) *handler {
	return &handler{
		ctx:             ctx,
		valuatorService: valuatorService,
	}
}

func (a *handler) Index(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	tmplParsed, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmplParsed.Execute(w, nil)
	if err != nil {
		log.Panic(err)
	}
}

func (a *handler) Summary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	text := r.FormValue("text")

	rank := a.valuatorService.CalculateRank(text)
	similarity := a.valuatorService.AddText(a.ctx, text)

	data := SummaryData{
		Text:       text,
		Rank:       rank,
		Similarity: similarity,
	}

	tmplParsed, err := template.ParseFiles("./templates/summary.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmplParsed.Execute(w, data)
	if err != nil {
		log.Panic(err)
	}
}

func (a *handler) About(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	tmplParsed, err := template.ParseFiles("./templates/about.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmplParsed.Execute(w, nil)
	if err != nil {
		log.Panic(err)
	}
}
