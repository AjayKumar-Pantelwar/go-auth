package main

import (
	"encoding/json"
	"go-authentication/src/db"
	"go-authentication/src/middleware"
	util "go-authentication/src/utils"
	"log"
	"net/http"
	"time"
)

func main() {
	database, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if http.MethodPost != r.Method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var user db.User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		returedUser, err := db.CreateUser(database, user)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		user = returedUser

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	})

	// login
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if http.MethodPost != r.Method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var user db.User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		returedUser, err := db.GetUser(database, user.Username)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid username"})
			return
		}

		if err := db.MatchPassword(database, returedUser, user.Password); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid credentials"})
			return
		}

		tokenString, err := util.GenerateJWT(returedUser.Username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		cookie := http.Cookie{
			Name:     "at",
			Value:    tokenString,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
			Path:     "/",
		}

		http.SetCookie(w, &cookie)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("x-user", returedUser.Username)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "successful login"})
	})

	http.Handle("/profile", middleware.Authenticate((http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if http.MethodGet != r.Method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		userId, ok := r.Context().Value("user").(string)

		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		returedUser, err := db.GetUser(database, userId)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("x-user", returedUser.Username)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "successful login"})
	}))))

	http.ListenAndServe(":8080", nil)
}
