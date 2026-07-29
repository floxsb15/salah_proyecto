package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

const authTokenDuration = 24 * time.Hour

func generarAuthToken(userID uint) string {
	exp := time.Now().Add(authTokenDuration).Unix()
	payload := strconv.FormatUint(uint64(userID), 10) + ":" + strconv.FormatInt(exp, 10)
	signature := firmarAuthToken(payload)
	return payload + ":" + signature
}

func obtenerUsuarioDesdeAuthHeader(header string) (uint, error) {
	token := strings.TrimSpace(header)
	token = strings.TrimPrefix(token, "Bearer ")
	parts := strings.Split(token, ":")
	if len(parts) != 3 {
		return 0, errors.New("token no valido")
	}

	payload := parts[0] + ":" + parts[1]
	expected := firmarAuthToken(payload)
	if !hmac.Equal([]byte(expected), []byte(parts[2])) {
		return 0, errors.New("token no valido")
	}

	exp, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || time.Now().Unix() > exp {
		return 0, errors.New("token expirado")
	}

	id, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil || id == 0 {
		return 0, errors.New("token no valido")
	}
	return uint(id), nil
}

func firmarAuthToken(payload string) string {
	mac := hmac.New(sha256.New, []byte(authTokenSecret()))
	mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func authTokenSecret() string {
	if secret := strings.TrimSpace(os.Getenv("AUTH_TOKEN_SECRET")); secret != "" {
		return secret
	}
	return "salah-motors-local-secret"
}
