package userhandler

import (
	"encoding/json"
	user "go-authentication/src/internal/core"
	userservice "go-authentication/src/internal/usecase"
	"go-authentication/src/pkg"
	"net/http"
	"time"
)

type UserHandler struct {
	userUsecase userservice.UserServiceImpl
}

func NewUserHandler(usecase userservice.UserServiceImpl) UserHandler {
	return UserHandler{
		userUsecase: usecase,
	}
}

func (u *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user user.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	returedUser, err := u.userUsecase.CreateUser(user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	user = returedUser

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// login
func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {

	var user user.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	returedUser, err := u.userUsecase.GetUser(user.Username)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid username"})
		return
	}

	if err := u.userUsecase.MatchPassword(returedUser, user.Password); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid credentials"})
		return
	}

	tokenString, err := pkg.GenerateJWT(returedUser.Username)
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
}

func (u *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {

	userId, ok := r.Context().Value("user").(string)

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	returedUser, err := u.userUsecase.GetUser(userId)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("x-user", returedUser.Username)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"Uid":returedUser.Uid, "Username": returedUser.Username})
}
