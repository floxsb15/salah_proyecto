package security

import (
	"backend-restaurant-delitto/internal/db"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	SessionCookieName = "session_token"
	CSRFCookieName    = "csrf_token"
	CSRFHeaderName    = "X-CSRF-Token"
	SessionDuration   = 8 * time.Hour
)

type contextKey string

const principalContextKey contextKey = "authenticated-principal"

type Principal struct {
	ID                 uint
	Role               string
	SessionVersion     uint
	CSRFToken          string
	ViaCookie          bool
	MustChangePassword bool
}

type tokenClaims struct {
	UserID         uint
	ExpiresAt      int64
	SessionVersion uint
	CSRFToken      string
}

var errInvalidToken = errors.New("token no valido")

func ValidateConfiguration() error {
	secret := strings.TrimSpace(os.Getenv("AUTH_TOKEN_SECRET"))
	if len(secret) < 32 {
		return errors.New("AUTH_TOKEN_SECRET debe contener al menos 32 caracteres aleatorios")
	}
	if isProduction() && !cookieSecure() {
		return errors.New("COOKIE_SECURE debe ser true en produccion")
	}
	return nil
}

func IssueToken(userID, sessionVersion uint) (string, string, error) {
	if userID == 0 {
		return "", "", errInvalidToken
	}
	if sessionVersion == 0 {
		sessionVersion = 1
	}
	csrfToken, err := randomToken(24)
	if err != nil {
		return "", "", err
	}
	claims := tokenClaims{
		UserID:         userID,
		ExpiresAt:      time.Now().Add(SessionDuration).Unix(),
		SessionVersion: sessionVersion,
		CSRFToken:      csrfToken,
	}
	payload := serializeClaims(claims)
	signature, err := sign(payload)
	if err != nil {
		return "", "", err
	}
	return payload + ":" + signature, csrfToken, nil
}

func SetSessionCookies(w http.ResponseWriter, token, csrfToken string) {
	maxAge := int(SessionDuration.Seconds())
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   cookieSecure(),
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     CSRFCookieName,
		Value:    csrfToken,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: false,
		Secure:   cookieSecure(),
		SameSite: http.SameSiteStrictMode,
	})
}

func ClearSessionCookies(w http.ResponseWriter) {
	for _, name := range []string{SessionCookieName, CSRFCookieName} {
		http.SetCookie(w, &http.Cookie{
			Name:     name,
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			Expires:  time.Unix(1, 0),
			HttpOnly: name == SessionCookieName,
			Secure:   cookieSecure(),
			SameSite: http.SameSiteStrictMode,
		})
	}
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, viaCookie := requestToken(r)
		claims, err := parseAndVerify(token)
		if err != nil {
			unauthorized(w)
			return
		}

		var account struct {
			ID                 uint
			Role               string
			SessionVersion     uint
			MustChangePassword bool
		}
		err = db.GDB.Table("usuarios AS u").
			Select("u.id, r.nombre AS role, u.session_version, u.must_change_password").
			Joins("JOIN roles AS r ON r.id = u.id_rol").
			Where("u.id = ? AND u.estado = true", claims.UserID).
			Scan(&account).Error
		if err != nil || account.ID == 0 || account.SessionVersion != claims.SessionVersion {
			unauthorized(w)
			return
		}

		principal := Principal{
			ID:                 account.ID,
			Role:               normalizeRole(account.Role),
			SessionVersion:     account.SessionVersion,
			CSRFToken:          claims.CSRFToken,
			ViaCookie:          viaCookie,
			MustChangePassword: account.MustChangePassword,
		}
		if viaCookie && requiresCSRF(r.Method) && !validCSRF(r, principal.CSRFToken) {
			http.Error(w, "Solicitud no autorizada", http.StatusForbidden)
			return
		}
		if principal.MustChangePassword && r.URL.Path != "/api/v1/me/password" && r.URL.Path != "/api/v1/logout" {
			http.Error(w, "Cambio de contrasena requerido", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), principalContextKey, principal)))
	})
}

func RequireRoles(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[normalizeRole(role)] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			principal, ok := PrincipalFromContext(r.Context())
			if !ok {
				unauthorized(w)
				return
			}
			if _, ok := allowed[principal.Role]; !ok {
				http.Error(w, "Acceso denegado", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func PrincipalFromContext(ctx context.Context) (Principal, bool) {
	principal, ok := ctx.Value(principalContextKey).(Principal)
	return principal, ok
}

func CurrentUserID(r *http.Request) (uint, bool) {
	principal, ok := PrincipalFromContext(r.Context())
	return principal.ID, ok
}

func CurrentUserHasRole(r *http.Request, roles ...string) bool {
	principal, ok := PrincipalFromContext(r.Context())
	if !ok {
		return false
	}
	for _, role := range roles {
		if principal.Role == normalizeRole(role) {
			return true
		}
	}
	return false
}

func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func LimitBody(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Body != nil {
				r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func serializeClaims(claims tokenClaims) string {
	return fmt.Sprintf("%d:%d:%d:%s", claims.UserID, claims.ExpiresAt, claims.SessionVersion, claims.CSRFToken)
}

func parseAndVerify(token string) (tokenClaims, error) {
	parts := strings.Split(strings.TrimSpace(token), ":")
	if len(parts) != 5 {
		return tokenClaims{}, errInvalidToken
	}
	payload := strings.Join(parts[:4], ":")
	expected, err := sign(payload)
	if err != nil || !hmac.Equal([]byte(expected), []byte(parts[4])) {
		return tokenClaims{}, errInvalidToken
	}
	userID, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil || userID == 0 {
		return tokenClaims{}, errInvalidToken
	}
	expiresAt, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || time.Now().Unix() >= expiresAt {
		return tokenClaims{}, errInvalidToken
	}
	version, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil || version == 0 || strings.TrimSpace(parts[3]) == "" {
		return tokenClaims{}, errInvalidToken
	}
	return tokenClaims{
		UserID:         uint(userID),
		ExpiresAt:      expiresAt,
		SessionVersion: uint(version),
		CSRFToken:      parts[3],
	}, nil
}

func sign(payload string) (string, error) {
	secret := strings.TrimSpace(os.Getenv("AUTH_TOKEN_SECRET"))
	if len(secret) < 32 {
		return "", errors.New("configuracion de autenticacion invalida")
	}
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil)), nil
}

func requestToken(r *http.Request) (string, bool) {
	if cookie, err := r.Cookie(SessionCookieName); err == nil && strings.TrimSpace(cookie.Value) != "" {
		return cookie.Value, true
	}
	header := strings.TrimSpace(r.Header.Get("Authorization"))
	parts := strings.Fields(header)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1], false
	}
	return "", false
}

func validCSRF(r *http.Request, expected string) bool {
	header := strings.TrimSpace(r.Header.Get(CSRFHeaderName))
	cookie, err := r.Cookie(CSRFCookieName)
	if err != nil || header == "" || cookie.Value == "" {
		return false
	}
	return hmac.Equal([]byte(header), []byte(expected)) && hmac.Equal([]byte(cookie.Value), []byte(expected))
}

func requiresCSRF(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return false
	default:
		return true
	}
}

func randomToken(size int) (string, error) {
	buffer := make([]byte, size)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buffer), nil
}

func normalizeRole(role string) string {
	return strings.ToLower(strings.TrimSpace(role))
}

func unauthorized(w http.ResponseWriter) {
	http.Error(w, "Autenticacion requerida", http.StatusUnauthorized)
}

func cookieSecure() bool {
	value, _ := strconv.ParseBool(strings.TrimSpace(os.Getenv("COOKIE_SECURE")))
	return value
}

func isProduction() bool {
	return strings.EqualFold(strings.TrimSpace(os.Getenv("APP_ENV")), "production")
}
