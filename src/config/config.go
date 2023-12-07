package config

import (
	validator "github.com/asaskevich/govalidator"
	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	log "github.com/sirupsen/logrus"
	"portfolio-cms-server/utils"
	"strings"
)

type configurations struct {
	AppEnv string `json:"app_env" koanf:"APP_ENV" valid:"required"`

	DBHosts    string `json:"db_hosts" koanf:"DB_HOSTS" valid:"required"`
	DBUsername string `json:"db_username" koanf:"DB_USERNAME" valid:"required"`
	DBPassword string `json:"db_password" koanf:"DB_PASSWORD" valid:"required"`
	DBPort     string `json:"db_port" koanf:"DB_PORT" valid:"required"`
	DBName     string `json:"db_name" koanf:"DB_NAME" valid:"required"`
}

var (
	parser = koanf.New(".")
	config configurations
)

// Init consumes the env file, validates it and parses it to a struct
func Init() (configurations, error) {
	err := parser.Load(file.Provider("config.env"), dotenv.Parser())
	if err != nil {
		err = loadFromOsEnv()
		if err != nil {
			return configurations{}, err
		}
	}

	err = parser.Unmarshal("", &config)
	if err != nil {
		return configurations{}, err
	}

	_, err = validator.ValidateStruct(config)
	if err != nil {
		utils.GetLogger().WithFields(log.Fields{"error": err.Error()}).Error("Error on config validation")
		return configurations{}, err
	}
	return config, nil
}

func loadFromOsEnv() (err error) {
	err = parser.Load(env.Provider("GO_CMS", "_", transformOsEnvToStructKeys), nil)
	if err != nil {
		return err
	}
	mapOsEnvToConfigurations()
	return
}

func mapOsEnvToConfigurations() {
	for key, value := range parser.All() {
		transformedKey := strings.Replace(key, ".", "_", -1)
		_ = parser.Set(transformedKey, value)
	}
}

func transformOsEnvToStructKeys(env string) string {
	return strings.Replace(env, "GO_CMS", "", -1)
}