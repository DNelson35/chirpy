package test

import (
	"testing"
	"time"
	"github.com/google/uuid"
	"github.com/DNelson35/chirpy/internal/auth"
)

func TestJWTFunctions(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "mysecretkey"
	expiresIn := time.Hour

	t.Run("MakeJWT", func(t *testing.T) {
		token, err := auth.MakeJWT(userID, tokenSecret, expiresIn)

		if err != nil {
			t.Errorf("MakeJWT returned an error: %v", err)
		}

		if token == "" {
			t.Errorf("MakeJWT returned an empty token")
		}
	})

	t.Run("ValidateJWT", func(t *testing.T) {
		token, err := auth.MakeJWT(userID, tokenSecret, expiresIn)
		if err != nil {
			t.Errorf("MakeJWT returned an error: %v", err)
		}

		parsedUserID, err := auth.ValidateJWT(token, tokenSecret)

		if err != nil {
			t.Errorf("ValidateJWT returned an error: %v", err)
		}

		if parsedUserID != userID {
			t.Errorf("ValidateJWT returned wrong userID, got: %v, want: %v", parsedUserID, userID)
		}
	})

	t.Run("InvalidToken", func(t *testing.T) {
		validToken, err := auth.MakeJWT(userID, tokenSecret, expiresIn)
		if err != nil {
			t.Errorf("MakeJWT returned an error: %v", err)
		}

		invalidToken := validToken + "invalid"
		_, err = auth.ValidateJWT(invalidToken, tokenSecret)
		
		if err == nil {
			t.Errorf("ValidateJWT should have returned an error for an invalid token")
		}
	})
}