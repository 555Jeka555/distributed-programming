package transport

import (
	"context"
	"errors"
	"html/template"
	"log"
	"net/http"
	"server/pkg/app/query"
	"server/pkg/app/service"
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
	ctx              context.Context
	valuatorService  service.ValuatorService
	textQueryService query.TextQueryService
}

func NewHandler(
	ctx context.Context,
	valuatorService service.ValuatorService,
	textQueryService query.TextQueryService,
) Handler {
	return &handler{
		ctx:              ctx,
		valuatorService:  valuatorService,
		textQueryService: textQueryService,
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

	textValue := r.FormValue("text")

	textID, rankID, err := a.valuatorService.AddText(a.ctx, textValue)
	similarity := 0
	if errors.Is(err, service.ErrKeyAlreadyExists) {
		similarity = 1
	}
	if err != nil && !errors.Is(err, service.ErrKeyAlreadyExists) {
		log.Panic(err)
	}

	text, err := a.textQueryService.GetTextByID(a.ctx, string(textID), string(rankID))
	if err != nil {
		log.Panic(err)
	}

	data := SummaryData{
		Text:       text.Value,
		Rank:       text.Rank,
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
