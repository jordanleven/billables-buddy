package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	EnvironmentDevelopment                    = "development"
	EnvironmentConfigurationForecastAccountID = "FORECAST_ACCOUNT_ID"
	EnvironmentConfigurationHarvestToken      = "HARVEST_ACCOUNT_TOKEN"
	EnvironmentConfigurationHarvestAccountID  = "HARVEST_ACCOUNT_ID"
)

var (
	flagEnvironment     string = ""
	forecastAccountId   string
	harvestAccountToken string
	harvestAccountId    string
)

func getEnvironmentVariable(v string) string {
	return os.Getenv(v)
}

func getAppEnvironment() string {
	env := os.Getenv("APP_ENV")

	switch {
	case env != "":
		return env
	case flagEnvironment != "":
		return flagEnvironment
	default:
		return EnvironmentDevelopment
	}
}

func loadVariablesFromEnvFile() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading environmental variables:", err)
	}
	forecastAccountId = getEnvironmentVariable(EnvironmentConfigurationForecastAccountID)
	harvestAccountToken = getEnvironmentVariable(EnvironmentConfigurationHarvestToken)
	harvestAccountId = getEnvironmentVariable(EnvironmentConfigurationHarvestAccountID)
}

func init() {
	appEnv := getAppEnvironment()

	if appEnv == EnvironmentDevelopment {
		loadVariablesFromEnvFile()
	}
}
