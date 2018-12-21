package authorize

import (
	"controller/config"
	"controller/events"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"model"
	"net/http"
	"regexp"
)

const UserContext = "user-context"

var (
	// secret must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	secret = []byte("1FE3B7IIB05CGBK2F6D17KJ61H36OLJJ")
	Store  = sessions.NewCookieStore(secret)
)

func CreateTokenEndpoint(w http.ResponseWriter, r *http.Request) {
	var user model.User
	_ = json.NewDecoder(r.Body).Decode(&user)

	search := regexp.MustCompile(`(?mi)^(` + user.Username + `):(.*)$`)

	// Check hashPwd authentication
	file, err := ioutil.ReadFile(config.Config.Users)
	if err != nil {
		json.NewEncoder(w).Encode(model.Exception{Msg: "Login error opening users file"})
		return
	}

	hashPwd := search.FindSubmatch(file)
	if hashPwd == nil || hashPwd[2] == nil {
		events.Msg <- "Authentication failed for: " + user.Username
		json.NewEncoder(w).Encode(model.Exception{Msg: "Login error invalid username/password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword(hashPwd[2], []byte(user.Password)); err != nil {
		events.Msg <- "Authentication failed for: " + user.Username
		json.NewEncoder(w).Encode(model.Exception{Msg: "Login error invalid username/password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"password": user.Password,
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		json.NewEncoder(w).Encode(model.Exception{Msg: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(model.JwtToken{Token: tokenString})
}

func ValidateMiddlewareToken(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("X-Token")
		if authorizationHeader != "undefined" && authorizationHeader != "" {
			token, err := jwt.Parse(authorizationHeader, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an err")
				}
				return secret, nil
			})
			if err != nil {
				json.NewEncoder(w).Encode(model.Exception{Msg: err.Error()})
				return
			}
			if token.Valid {
				context.Set(r, "decoded", token.Claims)
				next(w, r)
			} else {
				json.NewEncoder(w).Encode(model.Exception{Msg: "Invalid authorization token"})
			}
		} else {
			json.NewEncoder(w).Encode(model.Exception{Msg: "An authorization header is required"})
		}
	})
}

func ValidateMiddlewareCookie(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := Store.Get(r, UserContext)

		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			session.Values["message"] = "Unauthorized access"
			session.Save(r, w)

			if config.Config.Ssl {
				http.Redirect(w, r, "https://"+r.Host+"/login", http.StatusSeeOther)
			} else {
				http.Redirect(w, r, "http://"+r.Host+"/login", http.StatusSeeOther)
			}
			return
		}
		next(w, r)
	})
}
