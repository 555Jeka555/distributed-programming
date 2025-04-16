package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"html/template"
	"log"
	"net/http"
	"server/pkg/app/event"
	"server/pkg/app/query"
	"server/pkg/app/service"
	"server/pkg/infrastructure/redis/repo"
	"time"
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
	Text            string
	Rank            float64
	Similarity      int
	CentrifugoURL   string
	CentrifugoToken string
	Channel         string
	ProcessingID    string
}

type IndexData struct {
	Countries map[string]string
}

func NewHandler(
	ctx context.Context,
	writer event.Writer,
	textService service.TextService,
	textQueryService query.TextQueryService,
	regions map[string]string,
) Handler {
	return &handler{
		ctx:              ctx,
		writer:           writer,
		textService:      textService,
		textQueryService: textQueryService,
		regions:          regions,
	}
}

type handler struct {
	ctx              context.Context
	writer           event.Writer
	writerExchange   event.Writer
	textService      service.TextService
	textQueryService query.TextQueryService
	regions          map[string]string
}

func (h *handler) Index(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	data := IndexData{
		Countries: map[string]string{
			"RU": "Россия",
			"FR": "Франция",
			"DE": "Германия",
			"AE": "ОАЭ",
			"IN": "Индия",
		},
	}

	tmplParsed, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmplParsed.Execute(w, data)
	if err != nil {
		log.Panic(err)
	}
}

func (h *handler) SummaryCreate(w http.ResponseWriter, r *http.Request) {
	textValue := r.FormValue("text")
	county := r.FormValue("country")

	body, err := json.Marshal(map[string]any{
		"text":   textValue,
		"region": h.regions[county],
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
	if err != nil && !errors.Is(err, repo.NotFoundRegion) {
		log.Panic(err)
	}

	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}

	channel := "personal#" + textID
	data := SummaryData{
		Text:            text.Value,
		Rank:            text.Rank,
		Similarity:      text.Similarity,
		CentrifugoToken: generateCentrifugoToken(ip, channel),
		CentrifugoURL:   "ws://localhost:8000/connection/websocket",
		Channel:         channel,
		ProcessingID:    textID,
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

func generateCentrifugoToken(identifier string, channel string) string {
	claims := jwt.MapClaims{
		"sub":      identifier,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"channels": []string{channel},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("my_secret"))
	if err != nil {
		log.Printf("Ошибка генерации токена: %v", err)
		return ""
	}

	return signedToken
}
