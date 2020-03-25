package main

import (
	"github.com/AdhityaRamadhanus/cockpit/pkg/config"
	"github.com/AdhityaRamadhanus/cockpit/pkg/storage/redis"
	"github.com/AdhityaRamadhanus/cockpit/pkg/tasks"
	"github.com/jasonlvhit/gocron"
	"github.com/joho/godotenv"
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

func task() {
	logrus.Info("Run COVID Cron")
	cfg := buildConfig()
	setupLogger(cfg.Log)

	logrus.Debug("Connecting to redis at", cfg.Redis)
	redisClient, err := redis.CreateClient(cfg.Redis, 0)
	if err != nil {
		logrus.Fatalf("redis.CreateClient(cfg, 0) err = %v", err)
	}
	defer redisClient.Close()

	keyValueSvc := redis.NewKeyValueService(redisClient)

	taskJabar := tasks.TaskJabar{
		KeyValueService: keyValueSvc,
		URL:             cfg.URL,
		FeatureToggle:   cfg.FeatureToggle,
	}

	taskJatim := tasks.TaskJatim{
		KeyValueService: keyValueSvc,
		URL:             cfg.URL,
		FeatureToggle:   cfg.FeatureToggle,
	}

	taskJateng := tasks.TaskJateng{
		KeyValueService: keyValueSvc,
		URL:             cfg.URL,
		FeatureToggle:   cfg.FeatureToggle,
	}

	taskJakarta := tasks.TaskJakarta{
		KeyValueService: keyValueSvc,
		Tableau:         cfg.Tableau,
		FeatureToggle:   cfg.FeatureToggle,
	}

	taskJogja := tasks.TaskJogja{
		KeyValueService: keyValueSvc,
		Tableau:         cfg.Tableau,
		FeatureToggle:   cfg.FeatureToggle,
	}

	taskBanten := tasks.TaskBanten{
		KeyValueService: keyValueSvc,
		URL:             cfg.URL,
		FeatureToggle:   cfg.FeatureToggle,
	}

	taskIndonesia := tasks.TaskIndonesia{
		KeyValueService: keyValueSvc,
		Tableau:         cfg.Tableau,
		FeatureToggle:   cfg.FeatureToggle,
	}

	taskJabar.SaveJabarProvincialLevelData()
	taskJatim.SaveJatimData()
	taskJateng.SaveJatengProvincialLevelData()
	taskJakarta.SaveJakartaProvincialLevelData()
	taskIndonesia.SaveIndonesiaNationalLevelData()
	taskJogja.SaveJogjaProvincialLevelData()
	taskBanten.SaveBantenProvincialLevelData()
}

func main() {
	gocron.Every(10).Minute().From(gocron.NextTick()).Do(task)
	<-gocron.Start()
}
