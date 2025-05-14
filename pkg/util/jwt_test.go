package util

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"testing"
)

// helper to create a signed token
func createTestToken(secret string, claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

func TestVerifyJWT(t *testing.T) {

	t.Run("valid token", func(t *testing.T) {
		tokenString := "eyJhbGciOiJIUzI1NiIsImtpZCI6Ik5aTGxqK2RLcG9sZjNuR20iLCJ0eXAiOiJKV1QifQ.eyJhYWwiOiJhYWwxIiwiYW1yIjpbeyJtZXRob2QiOiJwYXNzd29yZCIsInRpbWVzdGFtcCI6MTc0NzIwOTkxMH1dLCJhcHBfbWV0YWRhdGEiOnsicHJvdmlkZXIiOiJlbWFpbCIsInByb3ZpZGVycyI6WyJlbWFpbCJdfSwiYXVkIjoiYXV0aGVudGljYXRlZCIsImVtYWlsIjoibnRubTIwMDNAZ21haWwuY29tIiwiZXhwIjoxNzQ3MjEzNTEwLCJpYXQiOjE3NDcyMDk5MTAsImlzX2Fub255bW91cyI6ZmFsc2UsImlzcyI6Imh0dHBzOi8veGhuZGp3dHRyb2dubGp4ZGt1em0uc3VwYWJhc2UuY28vYXV0aC92MSIsInBob25lIjoiIiwicm9sZSI6ImF1dGhlbnRpY2F0ZWQiLCJzZXNzaW9uX2lkIjoiNzgzM2FiMzctMjg4NS00NjQ5LWE2MmUtZmY1N2Y1ODhkMzk5Iiwic3ViIjoiMzVhMWU3MGYtZGU5Ni00NWQ4LTg2MTctYTg1YTk4ZGZiYmQ3IiwidXNlcl9tZXRhZGF0YSI6eyJlbWFpbF92ZXJpZmllZCI6dHJ1ZX0sInVzZXJfcm9sZSI6ImFkbWluIn0.ZNu7a9WxMQpVvgfC1V1O80NHt1QEgh1FsgBQwmrl3KU"

		token, err := VerifyJWT(tokenString)
		assert.NoError(t, err)
		assert.NotNil(t, token)
		assert.True(t, token.Valid)
	})
}
