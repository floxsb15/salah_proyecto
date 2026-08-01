package reports

import (
	"log"
	"net/http"
)

func respondInternalError(w http.ResponseWriter, operation string, err error) {
	log.Printf("%s: %v", operation, err)
	http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
}
