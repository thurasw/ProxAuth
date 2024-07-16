package sessions

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"sync"
	"time"
)

type SessionStore[T any] struct {
	name       string
	expiration time.Duration
	secure     bool
	domain     string
	store      map[string]Data[T]
	lock       sync.RWMutex // Lock to synchronize access to the session store
}
type Data[T any] struct {
	Data      *T
	ExpiresAt time.Time
}

func New[T any](name string, expiration time.Duration, secure bool, domain string) *SessionStore[T] {
	ss := &SessionStore[T]{
		name:       name,
		expiration: expiration,
		secure:     secure,
		domain:     domain,
		store:      make(map[string]Data[T]),
	}
	go ss.cleanExpired()
	return ss
}

// PutSession will store the session in the SessionStore.
func (st *SessionStore[T]) PutSession(w http.ResponseWriter, r *http.Request, sess *T) {
	// Generate new session data
	cookieValue := generateSessionToken()
	data := Data[T]{
		Data:      sess,
		ExpiresAt: time.Now().Add(st.expiration),
	}

	// Set cookie
	cookie := st.makeCookie()
	cookie.Value = cookieValue
	cookie.Expires = data.ExpiresAt

	http.SetCookie(w, cookie)
	w.Header().Add("Cache-Control", `no-cache="Set-Cookie"`)

	// Store session token
	st.lock.Lock()
	defer st.lock.Unlock()

	st.store[cookieValue] = data
}

// DeleteSession will delete the session from the SessionStore and set an empty cookie
func (st *SessionStore[T]) DeleteSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(st.name)
	if err != nil {
		return
	}
	st.Remove(cookie.Value)

	newCookie := st.makeCookie()
	cookie.Value = ""
	cookie.MaxAge = -1

	http.SetCookie(w, newCookie)
	w.Header().Add("Cache-Control", `no-cache="Set-Cookie"`)
}

// GetSessionFromRequest retrieves the session from the http.Request cookies.
// The function will return nil if the session does not exist within the http.Request cookies OR the session is not found.
func (st *SessionStore[T]) GetSession(r *http.Request) *T {
	cookie, err := r.Cookie(st.name)
	if err != nil {
		return nil
	}

	st.lock.RLock()
	defer st.lock.RUnlock()

	sess, ok := st.store[cookie.Value]
	if !ok || sess.ExpiresAt.Before(time.Now()) {
		return nil
	}
	return sess.Data
}

func (st *SessionStore[T]) Remove(token string) {
	st.lock.Lock()
	defer st.lock.Unlock()

	delete(st.store, token)
}

// This returns a shallow copy of the session store.
func (st *SessionStore[T]) GetSessions() map[string]Data[T] {
	st.lock.RLock()
	defer st.lock.RUnlock()

	desc := make(map[string]Data[T], len(st.store))
	for token, data := range st.store {
		desc[token] = data
	}
	return desc
}

func (st *SessionStore[T]) makeCookie() *http.Cookie {
	return &http.Cookie{
		Name:     st.name,
		Secure:   st.secure,
		Domain:   st.domain,
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Path:     "/",
	}
}

func generateSessionToken() string {
	b := make([]byte, 32) //32 bytes of entropy
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func (st *SessionStore[T]) cleanExpired() {
	c := time.Tick(30 * time.Minute)

	for range c {
		st.lock.Lock()
		now := time.Now()

		for token, data := range st.store {
			if data.ExpiresAt.Before(now) {
				delete(st.store, token)
			}
		}
		st.lock.Unlock()
	}
}
