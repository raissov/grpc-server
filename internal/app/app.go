package app

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"server/config"
	delivery "server/internal/delivery/http"
	"server/internal/repository"
	"server/internal/service"
	"server/pb"
	"syscall"
	"time"
)

func Run() error {
	var logger *zap.Logger
	var err error
	appMode := os.Getenv("APP_MODE")
	if appMode == "prod" {
		logger, err = zap.NewProduction()
		if err != nil {
			log.Println("Error while creating logger: ", err)
			return err
		}
	} else if appMode == "dev" {
		logger, err = zap.NewDevelopment()
		if err != nil {
			log.Println("Error while creating logger: ", err)
			return err
		}
	} else {
		log.Println("Error while creating logger: logger mode undefined")
		return fmt.Errorf("error while defining logger: app_mode is invalid %s", appMode)
	}
	//nolint
	defer logger.Sync()
	sugar := logger.Sugar()
	cfg, err := config.New()
	if err != nil {
		sugar.Errorf("errof while defining config: %v", err)
		return err
	}

	db, err := repository.NewPostgresPool(cfg.DB)
	if err != nil {
		sugar.Errorf("error while creating database: %v", err)
		return err
	}

	defer db.Close()

	repos := repository.NewRepository(db, time.Duration(cfg.DB.TimeOut)*time.Second, sugar)
	if err != nil {
		sugar.Errorf("error while creating repository: %v", err)
		return err
	}
	services := service.NewService(repos, cfg, sugar)
	handlers := delivery.NewHandler(services, sugar, cfg)

	srv := http.Server{
		Addr:    ":" + cfg.App.AppPort,
		Handler: handlers.InitRoutes(),
	}

	errChan := make(chan error, 1)

	go func(errChan chan<- error) {
		sugar.Infof("Starting client on port: %s\n", cfg.App.AppPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sugar.Error(err.Error())
			errChan <- err
		}
	}(errChan)
	log.Printf(cfg.GRPCConfig.Port)

	go func() {
		listener, err := net.Listen("tcp", "0.0.0.0"+":"+cfg.GRPCConfig.Port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)

		}
		s := grpc.NewServer()
		pb.RegisterMessageServiceServer(s, pb.UnimplementedMessageServiceServer{
			Logger:  sugar,
			Service: *services,
		})

		if err := s.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-quit:
		sugar.Error("Killing signal was received")
	case err := <-errChan:
		sugar.Errorf("HTTP client run error: %s", err)
	}

	sugar.Info("Shutting down client...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(cfg.App.AppShutdownTime))
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		sugar.Infof("Server forced to shutdown: %s", err)
	}
	return nil
}
