package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error){
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPassword(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error){
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject: fmt.Sprintf("%v", userID),
	})

	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error){
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.UUID{}, err
	}
	if !token.Valid {
		return uuid.UUID{}, fmt.Errorf("invalid token")
	}
	userID, err := uuid.Parse(claims.Subject)

	if err != nil {
		return uuid.UUID{}, err
	}

	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	headerString := headers.Get("Authorization")

	if headerString == "" || !strings.Contains( headerString, "Bearer") {
		return "", fmt.Errorf("no token provided")
	}
	tokenString := strings.Split(headerString, " ")
	
	return strings.TrimSpace(tokenString[1]), nil
}

func MakeRefreshToken() string {
	key := make([]byte, 32)
	rand.Read(key)
	hexStr := hex.EncodeToString(key)
	return hexStr
}

func GetApiKey(headers http.Header) (string, error){
	headerString := headers.Get("Authorization")

	if headerString == "" || !strings.Contains( headerString, "ApiKey") {
		return "", fmt.Errorf("no key provided")
	}
	apiString := strings.Split(headerString, " ")
	
	return strings.TrimSpace(apiString[1]), nil
}

