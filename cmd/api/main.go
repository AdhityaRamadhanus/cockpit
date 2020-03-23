package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/AdhityaRamadhanus/cockpit/pkg/config"
	server "github.com/AdhityaRamadhanus/cockpit/pkg/server/api"
	"github.com/AdhityaRamadhanus/cockpit/pkg/server/api/handlers"
	"github.com/AdhityaRamadhanus/cockpit/pkg/storage/redis"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func buildConfig() *config.Configuration {
	envPath := ".env"
	if err := godotenv.Load(envPath); err != nil {
		logrus.Fatalf("godotenv.Load(%q) err = %v", envPath, err)
	}

	yamlPath := "config.yaml"
	envPrefix := ""
	c, err := config.Build(yamlPath, envPrefix)
	if err != nil {
		logrus.Fatalf("config.Build(%q, %q) err = %v", yamlPath, envPrefix, err)
	}

	return c
}

func setupLogger(cfg config.LogConfig) {
	switch cfg.Format {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	switch cfg.Level {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func main() {
	cfg := buildConfig()
	setupLogger(cfg.Log)

	logrus.Debug("Connecting to redis at", cfg.Redis)
	redisClient, err := redis.CreateClient(cfg.Redis, 0)
	if err != nil {
		logrus.Fatalf("redis.CreateClient(cfg, 0) err = %v", err)
	}
	defer redisClient.Close()

	redisRateClient, err := redis.CreateClient(cfg.Redis, 1)
	if err != nil {
		logrus.Fatalf("redis.CreateClient(cfg, 1) err = %v", err)
	}
	defer redisRateClient.Close()

	keyValueSvc := redis.NewKeyValueService(redisClient)
	healthHandler := handlers.HealthzHandler{}
	metricHandler := handlers.MetricHandler{}
	nationalHandler := handlers.NationalHandler{
		KeyValueService: keyValueSvc,
	}
	provincialHandler := handlers.ProvincialHandler{
		KeyValueService: keyValueSvc,
	}

	server := server.NewServer(cfg.API, metricHandler, healthHandler, nationalHandler, provincialHandler)
	srv := server.CreateHTTPServer()

	// Handle SIGINT, SIGTERN, SIGHUP signal from OS
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		<-termChan
		logrus.Warn("Receiving signal, Shutting down server")
		srv.Close()
	}()

	logrus.WithField("Address", server.Address).Info("COVID API Server is running")
	logrus.Fatal(srv.ListenAndServe())
}
