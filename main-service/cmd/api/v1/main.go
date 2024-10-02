package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"log/slog"
	"main-service/config"
	"main-service/internal/api/handlers"
	"main-service/internal/api/middleware"
	"main-service/internal/kafka"
	"main-service/internal/repository/mongodb"
	repository2 "main-service/internal/repository/mongodb"
	repository "main-service/internal/repository/postgre"
	"main-service/internal/usecase"

	handlers_gen "main-service/internal/api/handlers/gen"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	err := godotenv.Load()
	if err != nil {
		logger.Error("failed to load .env file", slog.String("msg", err.Error()))
	}
	cfg, err := config.Read()

	if err != nil {
		log.Println("failed to read config:", err.Error())
		return
	}
	dbPool, err := repository.Connect(cfg)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer func() {
		if dbPool != nil {
			dbPool.Close()
		}
	}()

	producer, err := kafka.New(cfg)
	if err != nil {
		slog.Error("kafka connect", err)
		return
	}

	defer func() {
		if producer != nil {
			err := producer.Close()
			if err != nil {
				slog.Error("closing producer", err)
				return
			}
		}
	}()

	mongo, err := mongodb.Connect(context.Background(), cfg)
	if err != nil {
		slog.Error("error connecting to MongoDB", err)
		return
	}

	defer func() {
		if mongo != nil {
			if err := mongodb.Disconnect(context.Background(), mongo.Client()); err != nil {
				slog.Error("error disconnecting from MongoDB", err)
			}
		}
	}()

	storageUser := repository.NewStorageUser(dbPool)
	storageContent := repository.NewStorageContent(dbPool)
	storageMongo := repository2.NewStorageMongo(mongo.Client(), *cfg)

	userService := usecase.NewUserService(&storageUser)
	textContentService := usecase.NewTextService(storageMongo)
	contentService := usecase.NewContentService(&storageContent, storageMongo, producer, cfg.Kafka.KafkaTopic)

	handlerUser := handlers.NewUserHandler(userService)
	handlerContent := handlers.NewContentHandler(textContentService, contentService)

	apiHandler := handlers.NewAPIHandler(handlerUser, handlerContent)

	r := chi.NewRouter()

	fs := http.FileServer(http.Dir("../../../swagger-ui/dist"))
	r.Handle("/swagger-ui/*", http.StripPrefix("/swagger-ui/", fs))

	//TODO:
	userMiddleware := middleware.NewUserMiddleware(userService)
	r.Use(userMiddleware.Authenticate) //  для всех маршрутов

	r.Group(func(r chi.Router) {
		r.Route("/api/v1", func(r chi.Router) {
			r.Post("/user/login", handlers.LoginUserWrapper(apiHandler))

			r.Group(func(r chi.Router) {
				r.Use(userMiddleware.Authenticate)
				r.Mount("/", handlers_gen.HandlerFromMuxWithBaseURL(apiHandler, r, "/api/v1"))
			})
		})
	})

	server := &http.Server{
		Addr:    net.JoinHostPort(cfg.Host, cfg.Port),
		Handler: r,
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Starting server on port %v...\n", cfg.Port)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-stop

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Server gracefully stopped")
}
