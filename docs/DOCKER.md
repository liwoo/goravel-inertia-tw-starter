# Docker Setup for Goravel Blog

This project includes Docker configurations for both development and production environments.

## Development Setup

For local development with hot-reloading:

```bash
# Using docker-compose (recommended)
docker-compose up

# Or build and run manually
docker build -f Dockerfile.dev -t goravel-blog:dev .
docker run -p 3000:3000 -p 5173:5173 -v $(pwd):/app goravel-blog:dev
```

This will:
- Start the Goravel app on http://localhost:3000
- Start Vite dev server on http://localhost:5173
- Enable hot-reloading for both Go (using Air) and frontend (using Vite)
- Set up MySQL and Redis services

## Production Build

For production deployment:

```bash
# Build the production image
docker build -t goravel-blog:latest .

# Run the production container
docker run -p 3000:3000 goravel-blog:latest
```

## Known Issues

### TypeScript Errors
The production build currently bypasses TypeScript checking due to type errors in the codebase. These need to be fixed:
- Import errors with `@inertiajs/react`
- Missing type definitions
- Component type mismatches

To fix these locally:
```bash
# Run TypeScript compiler to see all errors
npm run build

# Or just check types without building
npx tsc --noEmit
```

## Environment Variables

The following environment variables are used:
- `APP_ENV`: Application environment (local, production)
- `APP_PORT`: Port for the Goravel app (default: 3000)
- `DB_HOST`, `DB_PORT`, `DB_DATABASE`, `DB_USERNAME`, `DB_PASSWORD`: Database configuration
- `VITE_HOST`: Host for Vite dev server (set to 0.0.0.0 in Docker)

## CI/CD Integration

The production Dockerfile is used in the GitHub Actions workflow for building and deploying to Kubernetes.