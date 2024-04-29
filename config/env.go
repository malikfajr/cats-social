package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	db_name     string
	db_port     int
	db_host     string
	db_username string
	db_password string
	db_params   string
	JWT_SECRET  string
	BCRYPT_SALT int
}

var Env Config

func InitEnv() {
	Env.db_name = getEnv("DB_NAME", "cats-social").(string)
	Env.db_port = getEnv("DB_PORT", 5432).(int)
	Env.db_host = getEnv("DB_HOST", "localhost").(string)
	Env.db_username = getEnv("DB_USERNAME", "postgres").(string)
	Env.db_password = getEnv("DB_PASSWORD", "secret").(string)
	Env.db_params = getEnv("DB_PARAMS", "sslmode=disable").(string)
	Env.JWT_SECRET = getEnv("JWT_SECRET", "not-define").(string)
	Env.BCRYPT_SALT = getEnv("BCRYPT_SALT", 8).(int)
}

func getEnv(key string, defaultValue interface{}) interface{} {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	switch defaultValue.(type) {
	case string:
		return value
	case int:
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return defaultValue
		}
		return intValue
	default:
		return defaultValue
	}
}

func GetDbAddress() string {
	address := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s", Env.db_username, Env.db_password, Env.db_host, Env.db_port, Env.db_name, Env.db_params)
	return address
}
