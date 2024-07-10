package userhandler

import (
	"encoding/json"
	"fmt"
	"go-authentication/src/internal/core/user"
	userservice "go-authentication/src/internal/usecase"
	"go-authentication/src/pkg"
	"net/http"
	"time"
)

type UserHandler struct {
	userService userservice.UserService
}

func NewUserHandler(usecase userservice.UserService) UserHandler {
	return UserHandler{
		userService: usecase,
	}
}

func (u *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user user.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	returedUser, err := u.userService.RegisterUser(user)

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

func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {

	var requestUser user.User

	if err := json.NewDecoder(r.Body).Decode(&requestUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	returnedUser, err := u.userService.GetUser(requestUser.Username)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid username"})
		return
	}

	if err := u.userService.MatchPassword(returnedUser, requestUser.Password); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid credentials"})
		return
	}

	tokenString, expireTime, err := pkg.GenerateJWT(returnedUser.Uid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	session, err := pkg.GenerateSession(returnedUser.Uid)
	fmt.Println("session generated", session.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = u.userService.CreateSession(session)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	
	atCookie := http.Cookie{
		Name:     "at",
		Value:    tokenString,
		Expires:  expireTime,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	}

	http.SetCookie(w, &atCookie)

	sessCookie := http.Cookie{
		Name:     "sess",
		Value:    session.Id.String(),
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	}
	fmt.Println("cookie generated", sessCookie.Value)

	http.SetCookie(w, &sessCookie)


	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("x-user", returnedUser.Username)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "successful login"})
}

func (u *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user").(int)

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{"error":"user not found in context"})
		return
	}

	returedUser, err := u.userService.GetUserById(userId)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("x-user", returedUser.Username)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(returedUser)
}

func (u *UserHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sess")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	returnedSession, err := u.userService.GetSession(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	err = u.userService.MatchSessionToken(cookie.Value, returnedSession.TokenHash)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	tokenString, expireTime, err := pkg.GenerateJWT(returnedSession.Uid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	atCookie := http.Cookie{
		Name:     "at",
		Value:    tokenString,
		Expires:  expireTime,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	}

	http.SetCookie(w, &atCookie)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "cookie refreshed succesfully"})
}

func (u *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user").(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{"error":"user not found in context"})
		return
	}

	err := u.userService.DeleteSession(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	atCookie := http.Cookie{
		Name:     "at",
		Value:    "",
		Expires:  time.Now(),
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	}
	http.SetCookie(w, &atCookie)

	sessCookie := http.Cookie{
		Name:     "sess",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	}
	http.SetCookie(w, &sessCookie)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "successful logout"})
}
