package routers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestEveryNonLoginRouteRequiresAuthentication(t *testing.T) {
	router := mux.NewRouter()
	InitEndPoints(router)

	err := router.Walk(func(route *mux.Route, _ *mux.Router, _ []*mux.Route) error {
		path, pathErr := route.GetPathTemplate()
		if pathErr != nil || path == "/api/v1/login" || path == "/healthz" {
			return pathErr
		}
		methods, methodsErr := route.GetMethods()
		if methodsErr != nil {
			// The path-prefix route that owns the subrouter has no method.
			return nil
		}
		path = strings.NewReplacer("{id}", "1").Replace(path)
		for _, method := range methods {
			request := httptest.NewRequest(method, path, nil)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			if response.Code != http.StatusUnauthorized {
				t.Errorf("%s %s without a session returned %d; want %d", method, path, response.Code, http.StatusUnauthorized)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walking routes: %v", err)
	}
}

func TestHealthEndpointIsPublic(t *testing.T) {
	router := mux.NewRouter()
	InitEndPoints(router)

	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("GET /healthz returned %d; want %d", response.Code, http.StatusOK)
	}
	if strings.TrimSpace(response.Body.String()) != "ok" {
		t.Fatalf("GET /healthz returned %q; want %q", response.Body.String(), "ok")
	}
}
