package test

import (
	"testing"
	"net/http"
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



func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name         string
		headers      http.Header
		expected     string
		expectingErr bool
	}{
		{
			name:         "Valid Bearer token",
			headers:      http.Header{"Authorization": {"Bearer mytoken123"}},
			expected:     "mytoken123",
			expectingErr: false,
		},
		{
			name:         "No Authorization header",
			headers:      http.Header{},
			expected:     "",
			expectingErr: true,
		},
		{
			name:         "Authorization header missing Bearer",
			headers:      http.Header{"Authorization": {"mytoken123"}},
			expected:     "",
			expectingErr: true,
		},
		{
			name:         "Empty Authorization header",
			headers:      http.Header{"Authorization": {""}},
			expected:     "",
			expectingErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := auth.GetBearerToken(tt.headers)

			if tt.expectingErr && err == nil {
				t.Errorf("Expected error, but got nil")
			}
			if !tt.expectingErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if got != tt.expected {
				t.Errorf("Expected token %v, but got %v", tt.expected, got)
			}
		})
	}
}