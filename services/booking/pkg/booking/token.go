package booking

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
)

const SECRET = "xxx"

type UserClaim struct {
	jwt.StandardClaims
	ID    string `json:"id"`
	Email string `json:"email"`
}

func GetTokenFromAuthorization(token string) string {
	return strings.Replace(token, "Bearer ", "", 1)
}

func DecodeUserFromToken(token string) (*UserClaim, error) {
	claim, err := jwt.ParseWithClaims(token, &UserClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET), nil
	})
	if err != nil {
		return nil, err
	}
	return claim.Claims.(*UserClaim), nil
}

func GetUserClaimFromRequest(r *http.Request) (*UserClaim, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("unauthorized")
	}
	authToken := GetTokenFromAuthorization(authHeader)
	userClaim, err := DecodeUserFromToken(authToken)
	if err != nil {
		return nil, errors.New("unauthorized")
	}
	return userClaim, nil
}
