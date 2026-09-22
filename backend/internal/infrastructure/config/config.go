package config

import "os"

type Config struct {
	Port          string
	MongoURI      string
	MongoDatabase string
}

func Load() Config {
	return Config{
		Port:          getEnv("PORT", "8080"),
		MongoURI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDatabase: getEnv("MONGO_DATABASE", "controle_chamados"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}

	return fallback
}
