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
	environment         string = EnvironmentDevelopment
	forecastAccountId   string
	harvestAccountToken string
	harvestAccountId    string
)

func getEnvironmentVariable(v string) string {
	return os.Getenv(v)
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
	if environment == EnvironmentDevelopment {
		loadVariablesFromEnvFile()
	}
}
