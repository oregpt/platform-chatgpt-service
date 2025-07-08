#!/bin/bash
# Deploy ChatGPT Service to Google Cloud Run

# Exit on any error
set -e

# Configuration
PROJECT_ID=$(gcloud config get-value project)
SERVICE_NAME="chatgpt-service"
REGION="us-central1"
IMAGE_NAME="gcr.io/$PROJECT_ID/$SERVICE_NAME:latest"

# Check if OPENAI_API_KEY is set
if [ -z "$OPENAI_API_KEY" ]; then
  echo "Error: OPENAI_API_KEY environment variable is not set"
  echo "Please set it with: export OPENAI_API_KEY=your_api_key"
  exit 1
fi

echo "Building $SERVICE_NAME Docker image..."
docker build -t $IMAGE_NAME .

echo "Pushing image to Google Container Registry..."
docker push $IMAGE_NAME

# Check if the secret exists, create if it doesn't
if ! gcloud secrets describe openai-api-key &>/dev/null; then
  echo "Creating secret for OPENAI_API_KEY..."
  echo -n "$OPENAI_API_KEY" | gcloud secrets create openai-api-key --data-file=-
else
  echo "Updating secret for OPENAI_API_KEY..."
  echo -n "$OPENAI_API_KEY" | gcloud secrets versions add openai-api-key --data-file=-
fi

echo "Deploying $SERVICE_NAME to Cloud Run..."
gcloud run deploy $SERVICE_NAME \
  --image=$IMAGE_NAME \
  --region=$REGION \
  --platform=managed \
  --allow-unauthenticated \
  --memory=512Mi \
  --cpu=1 \
  --min-instances=0 \
  --max-instances=10 \
  --set-env-vars=LOG_LEVEL=info,THREAD_TTL=60 \
  --set-secrets=OPENAI_API_KEY=openai-api-key:latest

echo "Deployment completed successfully!"
echo "Service URL: $(gcloud run services describe $SERVICE_NAME --region=$REGION --format='value(status.url)')"
