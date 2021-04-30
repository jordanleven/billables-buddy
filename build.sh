#!/bin/bash

# Variable names mapped to our .env file
ENV_VARIABLE_NAME_FORECAST_ACCOUNT="FORECAST_ACCOUNT_ID"
ENV_VARIABLE_NAME_HARVEST_ACCOUNT="HARVEST_ACCOUNT_ID"
ENV_VARIABLE_NAME_HARVEST_TOKEN="HARVEST_ACCOUNT_TOKEN"
PATH_BUILT_ARTIFACT=billablesbuddy.1m.goc
PATH_BITBAR_PLUGINS=~/Library/Application\ Support/xbar/plugins/

loadEnvVariables() {
  if [ -f .env ]
  then
    export $(cat .env | sed 's/#.*//g' | xargs)
  fi
}

buildBinary() {
  go build -o ./$PATH_BUILT_ARTIFACT -ldflags "
    -X main.flagEnvironment=production
    -X main.forecastAccountId=${!ENV_VARIABLE_NAME_FORECAST_ACCOUNT}
    -X main.harvestAccountId=${!ENV_VARIABLE_NAME_HARVEST_ACCOUNT}
    -X main.harvestAccountToken=${!ENV_VARIABLE_NAME_HARVEST_TOKEN}
  "
}

copyToPlugins() {
  cp "./$PATH_BUILT_ARTIFACT" "$PATH_BITBAR_PLUGINS"
}

# Load from our env file
loadEnvVariables

# Build the binary
buildBinary

# Copy to the plugins directory
copyToPlugins
