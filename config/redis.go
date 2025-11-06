package config

import (
	"os"
	"strconv"
	"time"
)

type RedisConfig struct {
	Host         string
	Port         string
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
	MaxRetries   int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// LoadRedisConfig loads Redis configuration from environment variables
func LoadRedisConfig() *RedisConfig {
	host := getEnv("REDIS_HOST", "localhost")
	port := getEnv("REDIS_PORT", "6379")
	password := getEnv("REDIS_PASSWORD", "")
	db := getEnvAsInt("REDIS_DB", 0)
	poolSize := getEnvAsInt("REDIS_POOL_SIZE", 10)
	minIdleConns := getEnvAsInt("REDIS_MIN_IDLE_CONNS", 5)
	maxRetries := getEnvAsInt("REDIS_MAX_RETRIES", 3)

	dialTimeout := getEnvAsDuration("REDIS_DIAL_TIMEOUT", 5*time.Second)
	readTimeout := getEnvAsDuration("REDIS_READ_TIMEOUT", 3*time.Second)
	writeTimeout := getEnvAsDuration("REDIS_WRITE_TIMEOUT", 3*time.Second)
	idleTimeout := getEnvAsDuration("REDIS_IDLE_TIMEOUT", 5*time.Minute)

	return &RedisConfig{
		Host:         host,
		Port:         port,
		Password:     password,
		DB:           db,
		PoolSize:     poolSize,
		MinIdleConns: minIdleConns,
		MaxRetries:   maxRetries,
		DialTimeout:  dialTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}
}

// GetRedisAddr returns the Redis address
func (rc *RedisConfig) GetRedisAddr() string {
	return rc.Host + ":" + rc.Port
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as integer with a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsDuration gets an environment variable as duration with a default value
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
