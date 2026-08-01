package controllers

import (
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	loginWindow   = 15 * time.Minute
	loginBlock    = 15 * time.Minute
	loginMaxFails = 5
)

type loginAttempt struct {
	failures     int
	windowStart  time.Time
	blockedUntil time.Time
}

type loginAttemptStore struct {
	mu       sync.Mutex
	byIP     map[string]loginAttempt
	byUser   map[string]loginAttempt
	lastGCAt time.Time
}

var authAttempts = loginAttemptStore{
	byIP:   make(map[string]loginAttempt),
	byUser: make(map[string]loginAttempt),
}

func (store *loginAttemptStore) allowIP(ip string) bool {
	store.mu.Lock()
	defer store.mu.Unlock()
	now := time.Now()
	store.cleanup(now)
	return !isBlocked(store.byIP[ip], now)
}

// failure records an invalid login and reports whether that username was
// already throttled before this attempt. Username throttling is evaluated
// only after checking the password so an attacker cannot lock out a valid
// user merely by knowing the account name.
func (store *loginAttemptStore) failure(ip, username string) bool {
	store.mu.Lock()
	defer store.mu.Unlock()
	now := time.Now()
	store.cleanup(now)
	key := normalizeLoginUser(username)
	wasUserBlocked := isBlocked(store.byUser[key], now)
	store.byIP[ip] = addFailure(store.byIP[ip], now)
	store.byUser[key] = addFailure(store.byUser[key], now)
	return wasUserBlocked
}

func (store *loginAttemptStore) success(ip, username string) {
	store.mu.Lock()
	defer store.mu.Unlock()
	delete(store.byIP, ip)
	delete(store.byUser, normalizeLoginUser(username))
}

func (store *loginAttemptStore) cleanup(now time.Time) {
	if !store.lastGCAt.IsZero() && now.Sub(store.lastGCAt) < loginWindow {
		return
	}
	for key, attempt := range store.byIP {
		if now.After(attempt.blockedUntil) && now.Sub(attempt.windowStart) > loginWindow {
			delete(store.byIP, key)
		}
	}
	for key, attempt := range store.byUser {
		if now.After(attempt.blockedUntil) && now.Sub(attempt.windowStart) > loginWindow {
			delete(store.byUser, key)
		}
	}
	store.lastGCAt = now
}

func addFailure(attempt loginAttempt, now time.Time) loginAttempt {
	if attempt.windowStart.IsZero() || now.Sub(attempt.windowStart) > loginWindow {
		attempt = loginAttempt{windowStart: now}
	}
	attempt.failures++
	if attempt.failures >= loginMaxFails {
		attempt.blockedUntil = now.Add(loginBlock)
	}
	return attempt
}

func isBlocked(attempt loginAttempt, now time.Time) bool {
	return !attempt.blockedUntil.IsZero() && now.Before(attempt.blockedUntil)
}

func normalizeLoginUser(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func requestClientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err != nil {
		host = strings.TrimSpace(r.RemoteAddr)
	}
	remoteIP := net.ParseIP(host)
	if trustedProxy(remoteIP) {
		if forwarded := strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0]); net.ParseIP(forwarded) != nil {
			return forwarded
		}
	}
	return host
}

func trustedProxy(ip net.IP) bool {
	if ip == nil {
		return false
	}
	if ip.IsLoopback() {
		return true
	}
	for _, value := range strings.Split(os.Getenv("TRUSTED_PROXY_CIDRS"), ",") {
		_, network, err := net.ParseCIDR(strings.TrimSpace(value))
		if err == nil && network.Contains(ip) {
			return true
		}
	}
	return false
}
