package security

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	loginAttempts = make(map[string]int)
	lastAttempt   = make(map[string]time.Time)
	mu            sync.Mutex
)

// ==================== PASSWORD HASHING ====================

// HashPassword securely hashes a password using bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword compares a plaintext password with its bcrypt hash.
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ==================== LOGIN THROTTLING ====================

// RecordLoginFailure increases failed login count for a username.
func RecordLoginFailure(username string) {
	mu.Lock()
	defer mu.Unlock()
	loginAttempts[username]++
	lastAttempt[username] = time.Now()
}

// ResetLoginFailures clears failed attempt count for a username.
func ResetLoginFailures(username string) {
	mu.Lock()
	defer mu.Unlock()
	delete(loginAttempts, username)
	delete(lastAttempt, username)
}

// CheckLoginLock returns true if user is locked out due to multiple failed logins.
func CheckLoginLock(username string) bool {
	mu.Lock()
	defer mu.Unlock()

	count := loginAttempts[username]
	if count >= 5 {
		if time.Since(lastAttempt[username]) < 5*time.Minute {
			return true
		}
		// Reset lockout after cooldown
		delete(loginAttempts, username)
		delete(lastAttempt, username)
	}
	return false
}

// ==================== SECURE COOKIE & SESSION ====================

// SecureSession configures a secure session with best practices (HTTPS-ready).
func SecureSession(store cookie.Store) gin.HandlerFunc {
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   600, // 10-minute session timeout
		HttpOnly: true,
		Secure:   true, // Use HTTPS in production
		SameSite: http.SameSiteStrictMode,
	})
	return sessions.Sessions("mysession", store)
}

// ==================== SAFE LOGIN HELPER ====================

// SecureLogin safely verifies password & throttles brute-force attacks.
func SecureLogin(c *gin.Context, username, password, storedHash string) error {
	if CheckLoginLock(username) {
		return errors.New("too many failed attempts — please try again in 5 minutes")
	}

	if !CheckPassword(password, storedHash) {
		RecordLoginFailure(username)
		return errors.New("invalid username or password")
	}

	ResetLoginFailures(username)
	return nil
}
