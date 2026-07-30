package main

import (
	"backend-restaurant-delitto/internal/db"
	"backend-restaurant-delitto/internal/functions"
	"backend-restaurant-delitto/internal/models"
	"backend-restaurant-delitto/internal/routers"
	"backend-restaurant-delitto/internal/security"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	var err error

	// Cargar el archivo .env
	if err := godotenv.Load(".env"); err != nil && !os.IsNotExist(err) {
		log.Printf("No se pudo cargar .env: %v", err)
	}
	if err := security.ValidateConfiguration(); err != nil {
		log.Fatal("Configuracion de seguridad invalida: ", err)
	}

	err = db.Connection()
	if err != nil {
		log.Printf("Error al conectar a la base de datos: %v", err)
		return
	}

	if migrationsEnabled() {
		if err := db.GDB.AutoMigrate(
			/* migraciones */
			&models.Rol{},
			&models.CategoriaVehiculo{},
			&models.SegmentoVehiculo{},
			&models.MarcaVehiculo{},
			&models.AnioVehiculo{},
			&models.Usuario{},
			&models.Cliente{},
			&models.Vehiculo{},
			&models.VentaVehiculo{},
			&models.CuotaCredito{},
			&models.GastoVario{},
			&models.Movimiento{},
		); err != nil {
			log.Fatal("Error al migrar los modelos de la db:", err)
		}

		if err := functions.CreacionInicial(); err != nil {
			log.Fatal("Error al iniciar los datos predeterminados: ", err)
		}
		if err := security.MigratePlaintextPasswords(); err != nil {
			log.Fatal("Error al migrar contrasenas heredadas: ", err)
		}
	}

	r := mux.NewRouter()
	routers.InitEndPoints(r)

	host := envOrDefault("SERVER_HOST", "127.0.0.1")
	port := envOrDefault("SERVER_PORT", "5000")
	server := &http.Server{
		Addr:              host + ":" + port,
		Handler:           security.SecurityHeaders(r),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       35 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	shutdownContext, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	serverErrors := make(chan error, 1)
	go func() {
		fmt.Printf("Servidor corriendo en http://%s\n", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("Error al iniciar el servidor: ", err)
		}
	case <-shutdownContext.Done():
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Cierre forzado del servidor: %v", err)
		}
	}
}

func migrationsEnabled() bool {
	if value := strings.TrimSpace(os.Getenv("RUN_MIGRATIONS")); value != "" {
		enabled, err := strconv.ParseBool(value)
		if err != nil {
			log.Fatal("RUN_MIGRATIONS debe ser true o false")
		}
		return enabled
	}
	return !strings.EqualFold(strings.TrimSpace(os.Getenv("APP_ENV")), "production")
}

func envOrDefault(name, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(name)); value != "" {
		return value
	}
	return fallback
}
