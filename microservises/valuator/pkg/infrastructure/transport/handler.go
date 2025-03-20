package transport

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"server/pkg/app/event"
	"server/pkg/app/query"
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

func NewHandler(
	ctx context.Context,
	writer event.Writer,
	textQueryService query.TextQueryService,
) Handler {
	return &handler{
		ctx:              ctx,
		writer:           writer,
		textQueryService: textQueryService,
	}
}

type handler struct {
	ctx              context.Context
	writer           event.Writer
	textQueryService query.TextQueryService
}

type eventBody struct {
	TextID     string `json:"text_id"`
	RankID     string `json:"rank_id"`
	Similarity int    `json:"similarity"`
}

func (h *handler) Index(w http.ResponseWriter, _ *http.Request) {
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

func (h *handler) Summary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	textValue := r.FormValue("text")

	body, err := json.Marshal(map[string]any{
		"text": textValue,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	respBody, err := h.writer.Write(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	var evt eventBody

	err = json.Unmarshal(respBody, &evt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	text, err := h.textQueryService.GetTextByID(h.ctx, evt.TextID, evt.RankID)
	if err != nil {
		log.Panic(err)
	}

	data := SummaryData{
		Text:       text.Value,
		Rank:       text.Rank,
		Similarity: evt.Similarity,
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

func (h *handler) About(w http.ResponseWriter, _ *http.Request) {
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
