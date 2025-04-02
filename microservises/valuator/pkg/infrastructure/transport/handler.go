package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"server/pkg/app/event"
	"server/pkg/app/query"
	"server/pkg/app/service"
)

type Handler interface {
	Index(w http.ResponseWriter, r *http.Request)
	About(w http.ResponseWriter, r *http.Request)
	SummaryCreate(w http.ResponseWriter, r *http.Request)
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
	textService service.TextService,
	textQueryService query.TextQueryService,
) Handler {
	return &handler{
		ctx:              ctx,
		writer:           writer,
		textService:      textService,
		textQueryService: textQueryService,
	}
}

type handler struct {
	ctx              context.Context
	writer           event.Writer
	writerExchange   event.Writer
	textService      service.TextService
	textQueryService query.TextQueryService
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

func (h *handler) SummaryCreate(w http.ResponseWriter, r *http.Request) {
	textValue := r.FormValue("text")

	// TODO передавать от сюда айдишники, и со страницы отправлять запросы на получение данных, убрать rcp и диспатчить событие из rankcalcualtor
	// TODO опрашивать раз в секунду редис
	body, err := json.Marshal(map[string]any{
		"text": textValue,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = h.writer.Write(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	textID := h.textService.GetTextID(textValue)
	redirectURL := fmt.Sprintf("/summary?textID=%s", textID)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (h *handler) Summary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	textID := r.URL.Query().Get("textID")
	if textID == "" {
		http.Error(w, "textID parameter is required", http.StatusBadRequest)
		return
	}

	text, err := h.textQueryService.GetTextByID(h.ctx, textID)
	if err != nil {
		log.Panic(err)
	}

	data := SummaryData{
		Text:       text.Value,
		Rank:       text.Rank,
		Similarity: text.Similarity,
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

	err = h.writer.WriteExchange(event.SimilarityCalculated{
		TextID:     textID,
		Similarity: text.Similarity,
	})
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
