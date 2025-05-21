package transport

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"html/template"
	"log"
	"net/http"
	"server/pkg/app/service"
	jwtinfra "server/pkg/infrastructure/jwt"
	"time"
)

type Handler interface {
	LoginPage(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	RefreshToken(w http.ResponseWriter, r *http.Request)
}

func NewHandler(
	ctx context.Context,
	jwtKey string,
	userService *service.UserService,
) Handler {
	return &handler{
		ctx:         ctx,
		jwtKey:      jwtKey,
		userService: userService,
	}
}

type handler struct {
	ctx         context.Context
	jwtKey      string
	userService *service.UserService
}

func (h *handler) LoginPage(w http.ResponseWriter, _ *http.Request) {
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

func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("reg-login")
	password := r.FormValue("reg-password")

	err := h.userService.CreateUser(h.ctx, login, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.login(w, login)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("login-password")

	log.Println("login", login)
	log.Println("password", password)

	isAuth, err := h.userService.Authenticate(h.ctx, login, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isAuth {
		http.Error(w, errors.New("password not matched").Error(), http.StatusUnauthorized)
		return
	}

	h.login(w, login)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func (h *handler) Logout(w http.ResponseWriter, r *http.Request) {
	deleteCookie(w, "access_token")
	deleteCookie(w, "refresh_token")
}

func (h *handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshTokenCookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "Refresh token not found", http.StatusUnauthorized)
		return
	}

	claims := &jwtinfra.Claims{}
	token, err := jwt.ParseWithClaims(refreshTokenCookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.jwtKey), nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	if time.Now().Unix() > claims.ExpiresAt {
		http.Error(w, "Refresh token expired", http.StatusUnauthorized)
		return
	}

	accessToken, accessExpirationTime, err := h.createToken(claims.Login, 5*time.Hour)
	if err != nil {
		http.Error(w, "Could not create access token", http.StatusInternalServerError)
		return
	}

	refreshToken, refreshExpirationTime, err := h.createToken(claims.Login, 30*24*time.Hour)
	if err != nil {
		http.Error(w, "Could not create refresh token", http.StatusInternalServerError)
		return
	}

	setCookie(w, "access_token", accessToken, accessExpirationTime)
	setCookie(w, "refresh_token", refreshToken, refreshExpirationTime)
}

func (h *handler) login(w http.ResponseWriter, login string) {
	accessToken, accessExpirationTime, err := h.createToken(login, 5*time.Hour)
	if err != nil {
		http.Error(w, errors.New("could not create access token").Error(), http.StatusUnauthorized)
		return
	}

	refreshToken, refreshExpirationTime, err := h.createToken(login, 30*24*time.Hour)
	if err != nil {
		http.Error(w, errors.New("could not create refresh token").Error(), http.StatusUnauthorized)
		return
	}

	setCookie(w, "access_token", accessToken, accessExpirationTime)
	setCookie(w, "refresh_token", refreshToken, refreshExpirationTime)
}

func (h *handler) createToken(login string, expirationTimeDur time.Duration) (string, time.Time, error) {
	expirationTime := time.Now().Add(expirationTimeDur)
	claims := &jwtinfra.Claims{
		Login: login,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.jwtKey))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}

func setCookie(w http.ResponseWriter, name, value string, expirationTime time.Time) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  expirationTime,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
}

func deleteCookie(w http.ResponseWriter, name string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
}
