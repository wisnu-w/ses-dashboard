#!/bin/bash

# Script to generate .env from config.yaml
# This allows Docker Compose to use values from config.yaml

CONFIG_FILE="ses-dashboard-monitoring/config/config.yaml"
ENV_FILE=".env"

echo "# Auto-generated from config.yaml" > $ENV_FILE
echo "# Edit config.yaml and run ./generate-env.sh to update" >> $ENV_FILE
echo "" >> $ENV_FILE

# Extract values from YAML and convert to environment variables
echo "# Database Configuration" >> $ENV_FILE
DB_HOST=$(grep -A 10 "database:" $CONFIG_FILE | grep "host:" | awk '{print $2}')
DB_PORT=$(grep -A 10 "database:" $CONFIG_FILE | grep "port:" | awk '{print $2}')
DB_NAME=$(grep -A 10 "database:" $CONFIG_FILE | grep "name:" | awk '{print $2}')
DB_USER=$(grep -A 10 "database:" $CONFIG_FILE | grep "user:" | awk '{print $2}')
DB_PASSWORD=$(grep -A 10 "database:" $CONFIG_FILE | grep "password:" | awk '{print $2}')
DB_SSLMODE=$(grep -A 10 "database:" $CONFIG_FILE | grep "sslmode:" | awk '{print $2}')

echo "DB_HOST=$DB_HOST" >> $ENV_FILE
echo "DB_PORT=$DB_PORT" >> $ENV_FILE
echo "DB_NAME=$DB_NAME" >> $ENV_FILE
echo "DB_USER=$DB_USER" >> $ENV_FILE
echo "DB_PASSWORD=$DB_PASSWORD" >> $ENV_FILE
echo "DB_SSLMODE=$DB_SSLMODE" >> $ENV_FILE
echo "" >> $ENV_FILE

echo "# Application Configuration" >> $ENV_FILE
APP_NAME=$(grep -A 10 "app:" $CONFIG_FILE | grep "name:" | awk '{print $2}')
APP_ENV=$(grep -A 10 "app:" $CONFIG_FILE | grep "env:" | awk '{print $2}')
APP_PORT=$(grep -A 10 "app:" $CONFIG_FILE | grep "port:" | awk '{print $2}')
JWT_SECRET=$(grep -A 10 "app:" $CONFIG_FILE | grep "jwt_secret:" | awk '{print $2}')
ENABLE_SWAGGER=$(grep -A 10 "app:" $CONFIG_FILE | grep "enable_swagger:" | awk '{print $2}')

echo "APP_NAME=$APP_NAME" >> $ENV_FILE
echo "APP_ENV=$APP_ENV" >> $ENV_FILE
echo "APP_PORT=$APP_PORT" >> $ENV_FILE
echo "JWT_SECRET=$JWT_SECRET" >> $ENV_FILE
echo "ENABLE_SWAGGER=$ENABLE_SWAGGER" >> $ENV_FILE
echo "" >> $ENV_FILE

echo "# AWS Configuration" >> $ENV_FILE
AWS_REGION=$(grep -A 10 "aws:" $CONFIG_FILE | grep "region:" | awk '{print $2}')
echo "AWS_REGION=$AWS_REGION" >> $ENV_FILE
echo "AWS_ACCESS_KEY=" >> $ENV_FILE
echo "AWS_SECRET_KEY=" >> $ENV_FILE
echo "" >> $ENV_FILE

echo "# Frontend Configuration" >> $ENV_FILE
echo "BACKEND_URL=http://backend:$APP_PORT" >> $ENV_FILE

echo "✅ Generated $ENV_FILE from $CONFIG_FILE"