# ChatGPT Integration Service for OreGPT Agent Platform

This service provides integration with OpenAI's ChatGPT API for the OreGPT Agent Platform.

## Features

- Standardized input/output schema for consistent integration
- Stateful conversation management using Chat Completion API
- Thread management per session
- Error handling and retry logic
- Containerized for Google Cloud Run deployment

## Technical Details

- Built in Go
- Uses OpenAI Chat Completion API for conversations
- Maps sessionId to conversation threads
- Maintains conversation history in memory
- Compatible with OpenAI SDK v1.14.2

## Environment Variables

- `OPENAI_API_KEY`: Your OpenAI API key
- `PORT`: Port for the service (default: 8080)
- `LOG_LEVEL`: Logging level (default: info)
- `THREAD_TTL`: Time-to-live for cached conversation threads in minutes (default: 60)

## API Endpoints

- `POST /chat`: Main endpoint for chat interactions

## Development

```bash
# Build
go build -o chatgpt-service

# Run locally
./chatgpt-service

# Build Docker image
docker build -t oregpt/chatgpt-service:latest .
```

## Deployment

```bash
# Deploy to Google Cloud Run
gcloud run deploy chatgpt-service --image oregpt/chatgpt-service:latest --platform managed
```
