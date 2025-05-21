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
	jwtinfra "server/pkg/infrastructure/jwt"
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
	jwtKey string,
	publisher event.Publisher,
	textService service.TextService,
	textQueryService query.TextQueryService,
	regions map[string]string,
) Handler {
	return &handler{
		ctx:              ctx,
		jwtKey:           jwtKey,
		publisher:        publisher,
		textService:      textService,
		textQueryService: textQueryService,
		regions:          regions,
	}
}

type handler struct {
	ctx              context.Context
	jwtKey           string
	publisher        event.Publisher
	textService      service.TextService
	textQueryService query.TextQueryService
	regions          map[string]string
}

func (h *handler) Index(w http.ResponseWriter, r *http.Request) {
	_, err := h.getLoginFromToken(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login-page", http.StatusMovedPermanently)
		return
	}

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
	login, err := h.getLoginFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	textValue := r.FormValue("text")
	county := r.FormValue("country")

	body, err := json.Marshal(map[string]any{
		"login":  login,
		"text":   textValue,
		"region": h.regions[county],
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = h.publisher.Publish(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	textID := h.textService.GetTextID(textValue)
	redirectURL := fmt.Sprintf("/valuator/summary?textID=%s", textID)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (h *handler) Summary(w http.ResponseWriter, r *http.Request) {
	login, err := h.getLoginFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

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

	log.Println("Text login", text.Login)
	log.Println("Token Login", login)
	if text.Login != "" && text.Login != login {
		http.Error(w, fmt.Sprintf("text is unavailable for this login %s", login), http.StatusForbidden)
		return
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

	err = h.publisher.PublishInExchange(event.SimilarityCalculated{
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

func (h *handler) getLoginFromToken(r *http.Request) (string, error) {
	tokenString, err := h.getTokenFromCookie(r, "access_token")
	if err != nil {
		return "", err
	}

	claims, err := h.parseAndValidateToken(tokenString)
	if err != nil {
		return "", fmt.Errorf("invalid token: %v", err)
	}

	log.Println("claims", claims)

	return claims.Login, nil
}

func (h *handler) parseAndValidateToken(tokenString string) (*jwtinfra.Claims, error) {
	claims := &jwtinfra.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(h.jwtKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if time.Now().Unix() > claims.ExpiresAt {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

func (h *handler) getTokenFromCookie(r *http.Request, cookieName string) (string, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return "", fmt.Errorf("cookie %s not found: %v", cookieName, err)
	}
	return cookie.Value, nil
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
