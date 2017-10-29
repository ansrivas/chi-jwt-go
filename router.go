package main

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/pkcs12"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
)

var tokenAuth *jwtauth.JwtAuth
var (
	// For simplicity these files are in the same folder as the app binary.
	// You shouldn't do this in production.
	privKeyPath, privKeyPass, pubKeyPath string
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func initKeys() {
	privKeyPath = filepath.Join(keyPath, "jwtsig-test-prv-ks.p12")
	privKeyPass = "test123"
	pubKeyPath = filepath.Join(keyPath, "jwtsig-test-pub-ks.pem")

	signBytes, err := ioutil.ReadFile(privKeyPath)
	fatal(err)

	sign, _, err := pkcs12.Decode(signBytes, privKeyPass)
	fatal(err)
	//If successful assign it to signkey
	signKey = sign.(*rsa.PrivateKey)

	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	fatal(err)
	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	fatal(err)

	tokenAuth = jwtauth.New("RS256", signKey, verifyKey)
}

// UserCredentials ...
type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Response ...
type Response struct {
	Data string `json:"data"`
}

// Token ...
type Token struct {
	Token string `json:"token"`
}

// NewRouter creates a new router with some basic end points
func NewRouter() *chi.Mux {
	initKeys()
	router := chi.NewRouter()

	router.Use(
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
				next.ServeHTTP(w, r)
			})
		},
		middleware.Recoverer,
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.StripSlashes,
	)

	// Non-Protected Endpoint(s)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to a public end point"))
	})

	router.Post("/login", LoginHandler)

	// Protected Endpoints
	router.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Get("/resource", ProtectedHandler)
	})

	return router
}

// ProtectedHandler ...
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	response := Response{fmt.Sprintf("protected area. hi %v", claims["user_id"])}
	JSONResponse(response, w)
}

// LoginHandler ...
func LoginHandler(w http.ResponseWriter, r *http.Request) {

	var user UserCredentials

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Error in request")
		return
	}

	if (strings.ToLower(user.Username) != "someone") || (user.Password != "p@ssword") {
		w.WriteHeader(http.StatusForbidden)
		fmt.Println("Error logging in")
		fmt.Fprint(w, "Invalid credentials")
		return
	}

	_, tokenString, err := tokenAuth.Encode(jwtauth.Claims{"user_id": 123,
		"exp": time.Now().Add(time.Hour * time.Duration(1)).Unix(),
		"iat": time.Now().Unix(),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error while signing the token")
		fatal(err)
	}

	response := Token{tokenString}
	JSONResponse(response, w)

}

// JSONResponse ...
func JSONResponse(response interface{}, w http.ResponseWriter) {

	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
