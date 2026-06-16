package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"schedule/internal/logger"
	admin_repository "schedule/internal/repository/admin"
	schedule_repository "schedule/internal/repository/schedule"
	admin_service "schedule/internal/service/admin"
	schedule_service "schedule/internal/service/schedule"
	admin_handler "schedule/internal/transport/admin"
	schedule_handler "schedule/internal/transport/schedule"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	pid := os.Getpid()
	fmt.Println("PID: ", pid)
	logger, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	logger.Info("Logger initializated successfully")

	connectString := os.Getenv("POSTGRES_URL")

	appPort := os.Getenv("APP_PORT")

	db, err := sql.Open("postgres", connectString)
	if err != nil {
		logger.Error("Failed to open database connection", zap.Error(err))
		panic(err)
	}

	if err := db.Ping(); err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		panic(err)
	}
	logger.Info("Database connection established successfully")

	// Repo init
	scheduleRepo := schedule_repository.NewScheduleRepo(db, logger)
	adminRepo := admin_repository.NewAdminRepository(db, logger)
	logger.Info("repositories initializated successfully")

	// Service init
	scheduleService := schedule_service.NewScheduleService(scheduleRepo, logger)
	adminService := admin_service.NewUserService(adminRepo, logger)
	logger.Info("services initializated successfully")

	// Transport init
	scheduleHandler := schedule_handler.NewScheduleHandler(scheduleService, logger)
	adminHandler := admin_handler.NewUserHandler(adminService, logger)
	logger.Info("handlers initializated successfully")

	router := mux.NewRouter()

	// Schedule handlers
	router.Path("/schedule").Methods("GET").HandlerFunc(scheduleHandler.GetSchedule)
	router.Path("/schedule/week").Methods("GET").HandlerFunc(scheduleHandler.GetWeekSchedule)

	// Admin handlers
	router.Path("/teachers").Methods("GET").HandlerFunc(adminHandler.GetTeachers)
	router.Path("/teachers").Methods("POST").HandlerFunc(adminHandler.AddTeacher)
	router.Path("/classrooms").Methods("GET").HandlerFunc(adminHandler.GetClassrooms)
	router.Path("/classrooms").Methods("POST").HandlerFunc(adminHandler.AddClassroom)
	router.Path("/subjects").Methods("GET").HandlerFunc(adminHandler.GetSubjects)
	router.Path("/subjects").Methods("POST").HandlerFunc(adminHandler.AddSubject)
	router.Path("/groups").Methods("GET").HandlerFunc(adminHandler.GetGroups)
	router.Path("/groups").Methods("POST").HandlerFunc(adminHandler.AddGroup)
	router.Path("/schedule").Methods("POST").HandlerFunc(adminHandler.CreateSchedule)
	router.Path("/schedule").Methods("DELETE").HandlerFunc(adminHandler.DeleteSchedule)

	// Web
	assetsDir := filepath.Join("web", "assets")
	assetsHandler := http.StripPrefix("/assets/", http.FileServer(http.Dir(assetsDir)))
	router.PathPrefix("/assets/").Methods("GET").Handler(assetsHandler)
	router.Path("/").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("web", "index.html"))
	})

	// Server
	srv := &http.Server{
		Addr:    ":" + appPort,
		Handler: router,
	}

	go func() {
		logger.Info("Starting server", zap.String("port", appPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start server", zap.Error(err))
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Info("Stopping the server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Forced server stop")
		panic(err)
	}

	logger.Info("Server gracefully stopped.")
}
