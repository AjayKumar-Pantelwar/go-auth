package util

import (
	"go-authentication/src/db"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("jhdsdjaiuuhgh")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, err
		}
		return nil, err
	}
	if !token.Valid {
		return nil, err
	}
	return claims, nil
}

func generateSession(userId int) (string, error) {
	tokenId := uuid.New()
	expiresAt := time.Now().Add(24 * time.Hour)
	issuedAt := time.Now()

	hashToken, err := bcrypt.GenerateFromPassword([]byte(tokenId.String()), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	session := db.Session{
		Id:        tokenId,
		Uid:       userId,
		TokenHash: string(hashToken),
		ExpiresAt: expiresAt,
		IssuedAt:  issuedAt,
	}

}
