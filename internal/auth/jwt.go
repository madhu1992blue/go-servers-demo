package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	registeredClaims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		registeredClaims,
	)
	signed, err  := token.SignedString([]byte(tokenSecret))
	return signed, err

}


func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	var claims jwt.RegisteredClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims,func (token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err !=  nil {
		return uuid.UUID{}, err
	}
	subject, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, nil
	}
	return uuid.Parse(subject)
}
