package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort           string
	DBPath             string
	WorldWidth         float64
	WorldHeight        float64
	InitialPop         int
	FoodCount          int
	FoodEnergy         float64
	MoveCost           float64
	SpeedFactor        float64
	InputSize          int
	HiddenSize         int
	OutputSize         int
	EatRadius          float64
	MutationRate       float64
	MutationStrength   float64
	ReproduceThreshold   float64
	AsexualThresholdMult float64
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		HTTPPort:           getEnv("HTTP_PORT", "8080"),
		DBPath:             getEnv("DB_PATH", "./database.db"),
		WorldWidth:         getEnvAsFloat("WORLD_WIDTH", 800.0),
		WorldHeight:        getEnvAsFloat("WORLD_HEIGHT", 600.0),
		InitialPop:         getEnvAsInt("INITIAL_POP", 20),
		FoodCount:          getEnvAsInt("FOOD_COUNT", 50),
		FoodEnergy:         getEnvAsFloat("FOOD_ENERGY", 70.0),
		MoveCost:           getEnvAsFloat("MOVE_COST", 0.05),
		SpeedFactor:        getEnvAsFloat("SPEED_FACTOR", 1.5),
		InputSize:          getEnvAsInt("INPUT_SIZE", 10),
		HiddenSize:         getEnvAsInt("HIDDEN_SIZE", 4),
		OutputSize:         getEnvAsInt("OUTPUT_SIZE", 2),
		EatRadius:          getEnvAsFloat("EAT_RADIUS", 10.0),
		MutationRate:       getEnvAsFloat("MUTATION_RATE", 0.1),
		MutationStrength:   getEnvAsFloat("MUTATION_STRENGTH", 0.2),
		ReproduceThreshold:   getEnvAsFloat("REPRODUCE_THRESHOLD", 150.0),
		AsexualThresholdMult: getEnvAsFloat("ASEXUAL_THRESHOLD_MULT", 1.5),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

func getEnvAsFloat(key string, defaultVal float64) float64 {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
		return value
	}
	return defaultVal
}
