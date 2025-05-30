# this file helps automatically deploy all the services to GCP

#!/bin/bash

set -e  # Stop on error


PROJECT_ID="capstone-457718"
REGION="us-west2"
REPO="quickstart-docker-repo"
SQL_INSTANCE="$PROJECT_ID:$REGION:capstone-db"
SERVICE_ACCOUNT="cloud-run-sa@$PROJECT_ID.iam.gserviceaccount.com"
REDIS_URL=redis://default:AXRzAAIjcDE3ZGI2YmViMjMxNGQ0MGE0OGQyYTU5MGZhZjY0MTFiOXAxMA@dashing-gibbon-29811.upstash.io:6379

# commands to deploy the fastapi service
echo " Deploying FASTAPI service..."

# creates fastapi Docker build and pushes to artifact registry
docker buildx build \
  --platform linux/amd64 \
  -t $REGION-docker.pkg.dev/$PROJECT_ID/$REPO/fastapi-app:latest \
  --push ./backend/fastapi

# deploys fastapi to Google Cloud Platform
gcloud run deploy fastapi-service \
  --image=$REGION-docker.pkg.dev/$PROJECT_ID/$REPO/fastapi-app:latest \
  --platform=managed \
  --region=$REGION \
  --add-cloudsql-instances=$SQL_INSTANCE \
  --service-account=$SERVICE_ACCOUNT \
  --set-secrets="DATABASE_URL=DATABASE_URL:latest,SESSION_KEY=SESSION_KEY:latest,REDIS_URL=REDIS_URL:latest" \
  --set-env-vars="CORS_ALLOWED_ORIGINS=https://frontend-app-182441842390.$REGION.run.app" \
  --allow-unauthenticated

# commands to deploy frontend
echo "Deploying FRONTEND service..."

# creates frontend Docker build and pushes to artifact registry
docker buildx build \
  --platform linux/amd64 \
  --no-cache \
  -t $REGION-docker.pkg.dev/$PROJECT_ID/$REPO/frontend-app:latest \
  --push ./frontend

# deploys frontend to Google Cloud Platform
gcloud run deploy frontend-app \
  --image=$REGION-docker.pkg.dev/$PROJECT_ID/$REPO/frontend-app:latest \
  --region=$REGION \
  --platform=managed \
  --allow-unauthenticated \
  --set-env-vars="VITE_FASTAPI_URL=https://fastapi-service-182441842390.$REGION.run.app,VITE_AUTH_API_URL=https://go-auth-service-182441842390.$REGION.run.app"

# commands to deploy auth
echo "Deploying AUTH service..."

# creates auth docker and pushes to artifact registry
docker buildx build \
  --platform linux/amd64 \
  -t $REGION-docker.pkg.dev/$PROJECT_ID/$REPO/go-auth:latest \
  --push ./backend/auth

# deploys auth service to Google Cloud Platform
gcloud run deploy go-auth-service \
  --image=$REGION-docker.pkg.dev/$PROJECT_ID/$REPO/go-auth:latest \
  --platform=managed \
  --region=$REGION \
  --add-cloudsql-instances=$SQL_INSTANCE \
  --service-account=$SERVICE_ACCOUNT \
  --set-secrets="DATABASE_URL=DATABASE_URL:latest,DISCUSSION_APP_SESSION_KEY=SESSION_KEY:latest" \
  --set-env-vars="CORS_ALLOWED_ORIGINS=https://frontend-app-182441842390.$REGION.run.app" \
  --allow-unauthenticated

echo "All services deployed"
