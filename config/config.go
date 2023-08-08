package config

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

type EnvConfig struct {
	Mode                  string
	SchedulerPort         string
	SchedulerAllowOrigins []string
	PublicKeySrc          []byte
	DBHost                string
	DBUser                string
	DBPassWord            string
	DBName                string
}

type JsonConfig struct{}

type ServerConfig struct {
	EnvConf  EnvConfig
	JsonConf JsonConfig
}

func getEnvVar(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", errors.New("environment variable not found: " + key)
	}
	return value, nil
}

func QueryServerConfig() (*ServerConfig, error) {
	mode, err := getEnvVar("MODE")
	if err != nil {
		return nil, err
	}

	schedulerPort, err := getEnvVar("SCHEDULER_PORT")
	if err != nil {
		return nil, err
	}

	schedulerAllowOriginsStr, err := getEnvVar("SCHEDULER_ALLOW_ORIGINS")
	if err != nil {
		return nil, err
	}
	schedulerAllowOrigins := strings.Split(schedulerAllowOriginsStr, ",")

	publicKeyFile, err := getEnvVar("PUBLIC_KEY_FILE")
	if err != nil {
		return nil, err
	}

	publicKeySrc, err := os.ReadFile(publicKeyFile)
	if err != nil {
		return nil, err
	}

	dbHost, err := getEnvVar("DB_HOST")
	if err != nil {
		return nil, err
	}

	dbUser, err := getEnvVar("DB_USER")
	if err != nil {
		return nil, err
	}

	dbPassword, err := getEnvVar("DB_PASSWORD")
	if err != nil {
		return nil, err
	}

	dbName, err := getEnvVar("DB_NAME")
	if err != nil {
		return nil, err
	}

	jsonPath, err := getEnvVar("CONFIG_JSON_FILE")
	if err != nil {
		return nil, err
	}

	// open the JSON file
	file, err := os.Open(jsonPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// create a JSON decoder and decode the file contents into jsonConf
	var jsonConf JsonConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&jsonConf); err != nil {
		return nil, err
	}

	return &ServerConfig{
		EnvConf: EnvConfig{
			Mode:                  mode,
			SchedulerPort:         schedulerPort,
			SchedulerAllowOrigins: schedulerAllowOrigins,
			PublicKeySrc:          publicKeySrc,
			DBHost:                dbHost,
			DBUser:                dbUser,
			DBPassWord:            dbPassword,
			DBName:                dbName,
		},
		JsonConf: jsonConf,
	}, nil
}
