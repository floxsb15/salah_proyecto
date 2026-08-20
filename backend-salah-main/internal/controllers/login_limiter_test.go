package controllers

import (
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestClientIPIgnoresForwardedHeaderFromUntrustedPeer(t *testing.T) {
	t.Setenv("TRUSTED_PROXY_CIDRS", "")
	request := httptest.NewRequest("POST", "/api/v1/login", nil)
	request.RemoteAddr = "203.0.113.10:1234"
	request.Header.Set("X-Forwarded-For", "198.51.100.5")

	if got := requestClientIP(request); got != "203.0.113.10" {
		t.Fatalf("requestClientIP() = %q", got)
	}
}

func TestClientIPAcceptsForwardedHeaderFromConfiguredProxy(t *testing.T) {
	t.Setenv("TRUSTED_PROXY_CIDRS", "172.30.10.0/24")
	request := httptest.NewRequest("POST", "/api/v1/login", nil)
	request.RemoteAddr = "172.30.10.8:1234"
	request.Header.Set("X-Forwarded-For", "198.51.100.5")

	if got := requestClientIP(request); got != "198.51.100.5" {
		t.Fatalf("requestClientIP() = %q", got)
	}
}

func TestUsernameThrottleDoesNotPreventPasswordVerificationFromAnotherIP(t *testing.T) {
	store := loginAttemptStore{
		byIP:   make(map[string]loginAttempt),
		byUser: make(map[string]loginAttempt),
	}

	for attempt := 0; attempt < loginMaxFails; attempt++ {
		store.failure("198.51.100."+strconv.Itoa(attempt+1), "admin@example.com")
	}

	if !store.allowIP("203.0.113.10") {
		t.Fatal("username failures blocked a different IP before password verification")
	}
	if !store.failure("203.0.113.10", "admin@example.com") {
		t.Fatal("invalid attempt for throttled username was not identified")
	}
}

func TestParseLoginCredentialsAcceptsFormBody(t *testing.T) {
	request := httptest.NewRequest("POST", "/api/v1/login", strings.NewReader("usuario=admin&contra=secret"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	credentials, err := parseLoginCredentials(request)
	if err != nil {
		t.Fatalf("parseLoginCredentials() error = %v", err)
	}
	if credentials.Usuario != "admin" || credentials.Contra != "secret" {
		t.Fatalf("parseLoginCredentials() = %#v", credentials)
	}
}

func TestParseLoginCredentialsRejectsMalformedJSON(t *testing.T) {
	request := httptest.NewRequest("POST", "/api/v1/login", strings.NewReader(`{"usuario":`))
	request.Header.Set("Content-Type", "application/json")

	if _, err := parseLoginCredentials(request); err == nil {
		t.Fatal("parseLoginCredentials() error = nil")
	}
}
