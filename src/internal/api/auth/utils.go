package auth

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// issues a new session cookie for the authenticated user id
func NewUserSession(userId *int, w http.ResponseWriter, r *http.Request) {
	// Only one active session per user
	/* Before issuing a new session for this user, we want to remove all old ones */
	st := authSessions.GetSessions()
	for token, v := range st {
		if *v.Data == *userId {
			// Invalidate any existing sessions for this user
			authSessions.Remove(token)
			break
		}
	}

	authSessions.PutSession(w, r, userId)
}

// Password hashing
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
