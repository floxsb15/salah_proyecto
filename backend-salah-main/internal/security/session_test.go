package security

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testSecret = "0123456789abcdef0123456789abcdef"

func TestIssueTokenCanBeVerified(t *testing.T) {
	t.Setenv("AUTH_TOKEN_SECRET", testSecret)

	token, csrf, err := IssueToken(42, 3)
	if err != nil {
		t.Fatalf("IssueToken() error = %v", err)
	}
	claims, err := parseAndVerify(token)
	if err != nil {
		t.Fatalf("parseAndVerify() error = %v", err)
	}
	if claims.UserID != 42 || claims.SessionVersion != 3 || claims.CSRFToken != csrf {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestTamperedTokenIsRejected(t *testing.T) {
	t.Setenv("AUTH_TOKEN_SECRET", testSecret)

	token, _, err := IssueToken(1, 1)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := parseAndVerify(token + "alterado"); err == nil {
		t.Fatal("tampered token was accepted")
	}
}

func TestConfigurationRejectsWeakSecret(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("AUTH_TOKEN_SECRET", "corta")
	if err := ValidateConfiguration(); err == nil {
		t.Fatal("weak secret was accepted")
	}
}

func TestProductionRequiresSecureCookies(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("AUTH_TOKEN_SECRET", testSecret)
	t.Setenv("COOKIE_SECURE", "false")
	if err := ValidateConfiguration(); err == nil {
		t.Fatal("production accepted insecure cookies")
	}
}

func TestCSRFRequiresMatchingCookieAndHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/logout", nil)
	req.AddCookie(&http.Cookie{Name: CSRFCookieName, Value: "csrf-value"})
	req.Header.Set(CSRFHeaderName, "csrf-value")
	if !validCSRF(req, "csrf-value") {
		t.Fatal("valid CSRF proof was rejected")
	}
	req.Header.Set(CSRFHeaderName, "otro")
	if validCSRF(req, "csrf-value") {
		t.Fatal("mismatched CSRF proof was accepted")
	}
}

func TestRequireRolesDeniesWrongRole(t *testing.T) {
	nextCalled := false
	handler := RequireRoles("admin")(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		nextCalled = true
	}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	principal := Principal{ID: 5, Role: "vendedor"}
	req = req.WithContext(context.WithValue(req.Context(), principalContextKey, principal))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, req)
	if response.Code != http.StatusForbidden || nextCalled {
		t.Fatalf("status = %d, nextCalled = %v", response.Code, nextCalled)
	}
}
