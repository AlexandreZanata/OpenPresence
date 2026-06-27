package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/app"
	infbiometric "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/infrastructure/biometric"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/interfaces/httpapi"
)

func main() {
	addr := envOr("ATTENDANCE_HTTP_ADDR", ":8080")
	dsn := envOr("DATABASE_URL", "postgres://attendance_app:attendance_app@localhost:5432/openpresence?sslmode=disable")
	bioAddr := envOr("BIOMETRIC_GRPC_ADDR", "127.0.0.1:9090")

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	bioClient, err := infbiometric.NewClient(bioAddr)
	if err != nil {
		log.Fatalf("biometric grpc: %v", err)
	}
	defer bioClient.Close()

	stack := app.NewPunchStack(app.PunchStackConfig{
		DB:        db,
		Biometric: app.BiometricGRPCAdapter{Client: bioClient},
	})
	handler := &httpapi.PunchHandler{Submit: stack.Handler}
	mux := httpapi.NewMux(handler)
	corsOrigins := httpapi.ParseAllowedOrigins(envOr("CORS_ALLOWED_ORIGINS", ""))
	srv := &http.Server{Addr: addr, Handler: httpapi.WithCORS(corsOrigins, mux)}

	go func() {
		log.Printf("attendance HTTP listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
