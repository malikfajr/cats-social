package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/malikfajr/cats-social/config"
	"github.com/malikfajr/cats-social/exception"
	"github.com/malikfajr/cats-social/helper"
	"github.com/malikfajr/cats-social/httpmux"
	"github.com/malikfajr/cats-social/models"
)

func main() {
	config.InitEnv()
	httpmux.InitValidator()

	err := models.InitDb(config.GetDbAddress())
	helper.PanicIfError(err)

	router := initializeRoutes()
	wrapper := use(router, loggingMiddleware, exception.RecoverWrap)

	server := &http.Server{
		Addr:    ":8080",
		Handler: wrapper,
	}

	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

// wrapper global middleware
func use(r *http.ServeMux, middlewares ...func(next http.Handler) http.Handler) http.Handler {
	var s http.Handler
	s = r

	for _, mw := range middlewares {
		s = mw(s)
	}

	return s
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s", r.URL.String())
		next.ServeHTTP(w, r)
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")

		if authorization == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("invalid token"))
			return
		}

		tokenString := authorization[7:]

		token, err := jwt.ParseWithClaims(tokenString, &config.CustomJWTClaim{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Env.JWT_SECRET), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("invalid token"))
			return
		}
		if claims, ok := token.Claims.(*config.CustomJWTClaim); ok {
			r.Header.Set("email", claims.Email)
			r.Header.Set("name", claims.Name)
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("invalid token"))
			return
		}
	})
}

func initializeRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		// The "/" pattern matches everything, so we need to check
		// that we're at the root here.
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		fmt.Fprint(w, "hello world")
	})

	mux.HandleFunc("POST /v1/user/register", httpmux.RegisterHandler)
	mux.HandleFunc("POST /v1/user/login", httpmux.LoginHandler)

	SaveCat := http.HandlerFunc(httpmux.SaveCat)
	mux.Handle("POST /v1/cat", authMiddleware(SaveCat))

	GetCat := http.HandlerFunc(httpmux.GetCat)
	mux.Handle("GET /v1/cat", authMiddleware(GetCat))

	NewMatch := http.HandlerFunc(httpmux.CreateMatch)
	mux.Handle("POST /v1/cat/match", authMiddleware(NewMatch))

	GetMatch := http.HandlerFunc(httpmux.GetMyMatch)
	mux.Handle("GET /v1/cat/match", authMiddleware(GetMatch))

	ApproveMatch := http.HandlerFunc(httpmux.ApproveMatch)
	mux.Handle("POST /v1/cat/match/approve", authMiddleware(ApproveMatch))

	RejectMatch := http.HandlerFunc(httpmux.RejectMatch)
	mux.Handle("POST /v1/cat/match/reject", authMiddleware(RejectMatch))

	DeleteMatch := http.HandlerFunc(httpmux.DeleteMatch)
	mux.Handle("DELETE /v1/cat/match/{id}", authMiddleware(DeleteMatch))

	UpdateCat := http.HandlerFunc(httpmux.UpdateCat)
	mux.Handle("PUT /v1/cat/{id}", authMiddleware(UpdateCat))

	DeleteCat := http.HandlerFunc(httpmux.DestroyCat)
	mux.Handle("DELETE /v1/cat/{id}", authMiddleware(DeleteCat))

	return mux
}
