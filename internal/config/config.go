package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort             string
	DBPath               string
	WorldWidth           float64
	WorldHeight          float64
	InitialPop           int
	FoodCount            int
	FoodEnergy           float64
	MoveCost             float64
	SpeedFactor          float64
	InputSize  int
	OutputSize int
	EatRadius            float64
	MutationRate         float64
	MutationStrength     float64
	ReproduceThreshold   float64
	AsexualThresholdMult float64
	MaxAge               float64

	// Ecosystem Control
	FoodSpawnChance    float64
	CrowdingDistance   float64
	CrowdingMultiplier float64

	SpeciationThreshold     float64
	MatingDistanceThreshold float64

	// Bio-improvements
	CarrionEnergyMult   float64 // Multiplier for dead creature's mass → carrion energy
	CarrionLifespan     int     // Ticks before carrion fully decays
	MaturityAgeFraction float64 // Fraction of MaxAge before reproduction is possible
	InbreedingThreshold float64 // Min genetic distance for healthy offspring
	InbreedingPenalty   float64 // Energy reduction fraction for inbred offspring

	// Advanced bio
	BrainCostPerNeuron float64 // Energy cost per hidden neuron per tick
	PheromoneDeposit   float64 // Amount of pheromone deposited per tick
	PheromoneDecay     float64 // Decay factor per tick (0.98 = 2% decay)
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		HTTPPort:             getEnv("HTTP_PORT", "8080"),
		DBPath:               getEnv("DB_PATH", "./database.db"),
		WorldWidth:           getEnvAsFloat("WORLD_WIDTH", 800.0),
		WorldHeight:          getEnvAsFloat("WORLD_HEIGHT", 600.0),
		InitialPop:           getEnvAsInt("INITIAL_POP", 20),
		FoodCount:            getEnvAsInt("FOOD_COUNT", 50),
		FoodEnergy:           getEnvAsFloat("FOOD_ENERGY", 70.0),
		MoveCost:             getEnvAsFloat("MOVE_COST", 0.05),
		SpeedFactor:          getEnvAsFloat("SPEED_FACTOR", 1.5),
		InputSize:  getEnvAsInt("INPUT_SIZE", 11),
		OutputSize: getEnvAsInt("OUTPUT_SIZE", 2),
		EatRadius:            getEnvAsFloat("EAT_RADIUS", 10.0),
		MutationRate:         getEnvAsFloat("MUTATION_RATE", 0.1),
		MutationStrength:     getEnvAsFloat("MUTATION_STRENGTH", 0.2),
		ReproduceThreshold:   getEnvAsFloat("REPRODUCE_THRESHOLD", 150.0),
		AsexualThresholdMult: getEnvAsFloat("ASEXUAL_THRESHOLD_MULT", 1.5),
		MaxAge:               getEnvAsFloat("MAX_AGE", 10000.0),

		FoodSpawnChance:    getEnvAsFloat("FOOD_SPAWN_CHANCE", 0.05), // ~3 food/sec at 60fps
		CrowdingDistance:   getEnvAsFloat("CROWDING_DISTANCE", 50.0),
		CrowdingMultiplier: getEnvAsFloat("CROWDING_MULTIPLIER", 0.1), // +10% BMR per neighbor

		SpeciationThreshold:     getEnvAsFloat("SPECIATION_THRESHOLD", 1.0),
		MatingDistanceThreshold: getEnvAsFloat("MATING_DISTANCE_THRESHOLD", 0.5),

		CarrionEnergyMult:   getEnvAsFloat("CARRION_ENERGY_MULT", 30.0),
		CarrionLifespan:     getEnvAsInt("CARRION_LIFESPAN", 600),
		MaturityAgeFraction: getEnvAsFloat("MATURITY_AGE_FRACTION", 0.05),
		InbreedingThreshold: getEnvAsFloat("INBREEDING_THRESHOLD", 0.15),
		InbreedingPenalty:   getEnvAsFloat("INBREEDING_PENALTY", 0.2),

		BrainCostPerNeuron: getEnvAsFloat("BRAIN_COST_PER_NEURON", 0.005),
		PheromoneDeposit:   getEnvAsFloat("PHEROMONE_DEPOSIT", 0.1),
		PheromoneDecay:     getEnvAsFloat("PHEROMONE_DECAY", 0.98),
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
