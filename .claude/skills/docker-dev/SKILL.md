---
name: docker-dev
description: Manage the Docker development environment. Start, stop, restart services, view logs, troubleshoot containers, or run with optional tools (pgadmin, redis-commander, mailhog, minio).
argument-hint: <start|stop|restart|status|logs|clean> [--tools] [--mail] [--storage]
disable-model-invocation: true
allowed-tools: Bash, Read, Grep
---

# Docker Development Environment

Manage the local Docker development environment. Action: **$ARGUMENTS**

Parse the arguments:
- First argument: action (`start`, `stop`, `restart`, `status`, `logs`, `clean`, `rebuild`)
- `--tools`: Include pgadmin and redis-commander
- `--mail`: Include mailhog for email testing
- `--storage`: Include minio for object storage
- `--all`: Include all optional services
- If no action given, default to `status`

The compose files are in `docker-compose/`.

## Actions

### `start`
Start the development environment with hot reload:
```bash
cd /Users/liwu/Projects/go/blog
docker compose -f docker-compose/docker-compose.dev.yml up -d
```

Add profiles based on flags:
- `--tools`: `docker compose -f docker-compose/docker-compose.yml --profile tools up -d pgadmin redis-commander`
- `--mail`: `docker compose -f docker-compose/docker-compose.yml --profile mail up -d mailhog`
- `--storage`: `docker compose -f docker-compose/docker-compose.yml --profile storage up -d minio`

After starting, show:
- Container status table
- Port mappings (backend: 3000, frontend: 5173, postgres: 5432, redis: 6379)
- Optional service URLs if started (pgadmin: 5050, mailhog: 8025, minio: 9001, redis-commander: 8081)

### `stop`
```bash
docker compose -f docker-compose/docker-compose.dev.yml down
docker compose -f docker-compose/docker-compose.yml down
```

### `restart`
Stop then start (preserving volumes).

### `status`
Show detailed status:
```bash
docker compose -f docker-compose/docker-compose.dev.yml ps -a
docker compose -f docker-compose/docker-compose.yml ps -a
```
Also check:
- Database connectivity: `docker exec <postgres-container> pg_isready`
- Redis connectivity: `docker exec <redis-container> redis-cli ping`
- Volume usage: `docker system df -v | head -20`

### `logs`
Follow logs for all or a specific service:
```bash
docker compose -f docker-compose/docker-compose.dev.yml logs -f --tail=100
```
If a service name is passed after `logs`, filter to that service.

### `clean`
Remove all containers, volumes, and orphaned networks:
```bash
docker compose -f docker-compose/docker-compose.dev.yml down -v --remove-orphans
docker compose -f docker-compose/docker-compose.yml down -v --remove-orphans
docker volume prune -f
```
Warn the user that this deletes database data.

### `rebuild`
Force rebuild of dev containers (useful after Dockerfile or dependency changes):
```bash
docker compose -f docker-compose/docker-compose.dev.yml up -d --build --force-recreate
```

## Troubleshooting

If any container is unhealthy or not running, automatically:
1. Show the container logs (last 50 lines)
2. Check if ports are already in use (`lsof -i :<port>`)
3. Suggest fixes based on common issues:
   - Port conflict → suggest stopping the conflicting process
   - OOM → suggest increasing Docker memory
   - DB connection refused → check if postgres is healthy first
