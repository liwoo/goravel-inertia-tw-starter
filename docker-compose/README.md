# Docker Compose Setup

This folder contains Docker Compose configurations for the Books Database application.

## Files

| File | Purpose |
|------|---------|
| `docker-compose.yml` | Main compose file - builds app from Dockerfile, includes PostgreSQL and Redis |
| `docker-compose.dev.yml` | Development mode - hot reload with Air (Go) and Vite (frontend) |
| `docker-compose.test.yml` | Testing - ephemeral databases for integration and E2E tests |

## Quick Start

### Production-like local environment

Builds the application from the Dockerfile and starts all services:

```bash
cd docker-compose
docker compose up -d
```

Access the application at: http://localhost:3000

### Development with hot reload

For active development with automatic reloading:

```bash
cd docker-compose
docker compose -f docker-compose.dev.yml up -d
```

- Backend (Go with Air): http://localhost:3000
- Frontend (Vite): http://localhost:5173

### Running tests

```bash
cd docker-compose

# Start test databases
docker compose -f docker-compose.test.yml up -d postgres-test redis-test

# Run Go tests
docker compose -f docker-compose.test.yml run --rm test

# Run frontend tests
docker compose -f docker-compose.test.yml run --rm frontend-test

# Cleanup
docker compose -f docker-compose.test.yml down -v
```

## Optional Services

The main `docker-compose.yml` includes optional services via profiles:

```bash
# Include MinIO object storage
docker compose --profile storage up -d

# Include Mailhog for email testing
docker compose --profile mail up -d

# Include pgAdmin and Redis Commander
docker compose --profile tools up -d

# All optional services
docker compose --profile storage --profile mail --profile tools up -d
```

## Service Ports

| Service | Port | Description |
|---------|------|-------------|
| app | 3000 | Main application |
| postgres | 5432 | PostgreSQL database |
| redis | 6379 | Redis cache |
| minio | 9000, 9001 | MinIO API and Console (profile: storage) |
| mailhog | 1025, 8025 | SMTP and Web UI (profile: mail) |
| pgadmin | 5050 | Database management (profile: tools) |
| redis-commander | 8081 | Redis management (profile: tools) |

## Environment Variables

The compose files set sensible defaults for local development. Key variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_DATABASE` | books_db | PostgreSQL database name |
| `DB_USERNAME` | books | PostgreSQL user |
| `DB_PASSWORD` | books_password | PostgreSQL password |
| `REDIS_HOST` | redis | Redis hostname |
| `APP_ENV` | local | Application environment |

## Volumes

Data is persisted in Docker volumes:

- `postgres-data` - PostgreSQL data
- `redis-data` - Redis data
- `app-storage` - Application storage (uploads, logs)

To reset all data:

```bash
docker compose down -v
```

## Troubleshooting

### View logs

```bash
docker compose logs -f app      # Application logs
docker compose logs -f postgres # Database logs
docker compose logs -f redis    # Redis logs
```

### Rebuild application

```bash
docker compose up -d --build
```

### Connect to PostgreSQL

```bash
docker compose exec postgres psql -U books -d books_db
```

### Connect to Redis

```bash
docker compose exec redis redis-cli
```
