package config

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Configuration provide package level configuration for vendtron by reading from config.yaml first then overwrite if any with env (Default read from .env)
type Configuration struct {
	Env           string              `yaml:"env"`
	API           ApiConfig           `yaml:"api"`
	Redis         RedisConfig         `yaml:"redis"`
	Log           LogConfig           `yaml:"log"`
	Tableau       TableauConfig       `yaml:"tableau"`
	URL           URLConfig           `yaml:"url"`
	FeatureToggle FeatureToggleConfig `yaml:"feature_toggle"`
}

type ApiConfig struct {
	Port int    `yaml:"port" envconfig:"API_PORT"`
	Host string `yaml:"host" envconfig:"API_HOST"`
}

type URLConfig struct {
	JabarCoreData  string `yaml:"jabar_core_data"`
	JatimDraxiData string `yaml:"jatim_draxi_data"`
	JatengWeb      string `yaml:"jateng_web"`
}

type FeatureToggleConfig struct {
	Indonesia bool `yaml:"indonesia"`
	Jakarta   bool `yaml:"jakarta"`
	Jabar     bool `yaml:"jabar"`
	Jateng    bool `yaml:"jateng"`
	Jatim     bool `yaml:"jatim"`
}

type TableauConfig struct {
	Workbooks  []string `yaml:"workbooks"`
	Dashboards []string `yaml:"dashboards"`
}

type RedisConfig struct {
	Port     int    `yaml:"port" envconfig:"REDIS_PORT"`
	Host     string `yaml:"host" envconfig:"REDIS_HOST"`
	Password string `yaml:"password" envconfig:"REDIS_PASSWORD"`
	DB       int    `yaml:"db" envconfig:"REDIS_DB"`
}

type LogConfig struct {
	Level  string `yaml:"level" envconfig:"LOG_LEVEL"`
	Format string `yaml:"format" envconfig:"LOG_FORMAT"`
}

func Build(yamlPath, envPrefix string) (*Configuration, error) {
	var cfg Configuration
	f, err := os.Open(yamlPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	if err := envconfig.Process(envPrefix, &cfg.API); err != nil {
		return nil, errors.Wrap(err, "envconfig.Process(envPrefix, &cfg.API) err")
	}

	if err := envconfig.Process(envPrefix, &cfg.URL); err != nil {
		return nil, errors.Wrap(err, "envconfig.Process(envPrefix, &cfg.URL) err")
	}

	if err := envconfig.Process(envPrefix, &cfg.FeatureToggle); err != nil {
		return nil, errors.Wrap(err, "envconfig.Process(envPrefix, &cfg.FeatureToggle) err")
	}

	if err := envconfig.Process(envPrefix, &cfg.Tableau); err != nil {
		return nil, errors.Wrap(err, "envconfig.Process(envPrefix, &cfg.Tableau) err")
	}

	if err := envconfig.Process(envPrefix, &cfg.Redis); err != nil {
		return nil, errors.Wrap(err, "envconfig.Process(envPrefix, &cfg.Redis) err")
	}

	if err := envconfig.Process(envPrefix, &cfg.Log); err != nil {
		return nil, errors.Wrap(err, "envconfig.Process(envPrefix, &cfg.Log) err")
	}

	return &cfg, nil
}
