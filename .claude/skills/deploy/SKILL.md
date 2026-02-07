---
name: deploy
description: Deploy the application to staging or production. Validates, builds, pushes Docker image, deploys via Helm, and runs smoke tests.
argument-hint: <staging|production> [--skip-build] [--dry-run]
disable-model-invocation: true
allowed-tools: Read, Bash, Grep, Glob, Edit
---

# Deploy Application

Deploy this application to the **$ARGUMENTS** environment.

Parse the arguments:
- First argument: environment (`staging` or `production`). Default to `staging` if not specified.
- `--skip-build`: Skip Docker build/push, deploy with existing image
- `--dry-run`: Render Helm templates and show diff without actually deploying

## Pre-flight Checks

Before deploying, verify ALL of these:

1. **Git status is clean** - no uncommitted changes (warn if dirty, ask to continue)
2. **On correct branch** - `develop` for staging, `main` for production
3. **Read `package.json`** to get the current app version
4. **Sync Chart.yaml appVersion** with the version from package.json
5. **Helm lint passes** with the target values file:
   ```bash
   helm lint helm/goravel-blog -f helm/goravel-blog/values.<env>.yaml
   ```

## Build Phase (skip if --skip-build)

Read `DOCKER_REGISTRY` and `DOCKER_IMAGE_NAME` from `.env` (fallbacks: `docker.io` and `books-database`).
Construct the full image name: `$DOCKER_REGISTRY/<dockerhub-username>/$DOCKER_IMAGE_NAME`

1. **Build Docker image:**
   ```bash
   docker build -t $DOCKER_REGISTRY/<username>/$DOCKER_IMAGE_NAME:<tag> \
     --build-arg BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ') \
     --build-arg VCS_REF=$(git rev-parse --short HEAD) \
     .
   ```
   - For staging: tag = `develop`
   - For production: tag = version from package.json

2. **Scan image with Trivy** (if available):
   ```bash
   docker run --rm aquasec/trivy image --severity CRITICAL,HIGH $IMAGE:<tag>
   ```
   Report findings but don't block unless CRITICAL vulnerabilities found.

3. **Push image:**
   ```bash
   docker push $IMAGE:<tag>
   ```

## Deploy Phase

### Dry Run Mode (--dry-run)
Render templates and show what would change:
```bash
helm template goravel-blog helm/goravel-blog -f helm/goravel-blog/values.<env>.yaml
helm diff upgrade goravel-blog helm/goravel-blog -f helm/goravel-blog/values.<env>.yaml -n <namespace>
```
Stop here for dry-run.

### Actual Deployment
1. **Determine namespace**: `books-staging` or `books-production`
2. **Run Helm upgrade:**
   ```bash
   helm upgrade --install goravel-blog helm/goravel-blog \
     -f helm/goravel-blog/values.<env>.yaml \
     -n <namespace> --create-namespace \
     --set image.tag=<tag> \
     --atomic --timeout 10m \
     --wait
   ```
   For production, ask for explicit confirmation before executing.

3. **Verify rollout:**
   ```bash
   kubectl rollout status deployment/goravel-blog -n <namespace> --timeout=300s
   ```

## Post-Deploy Verification

1. **Check pod status:**
   ```bash
   kubectl get pods -n <namespace> -l app.kubernetes.io/name=goravel-blog
   ```

2. **Get ingress URL:**
   ```bash
   kubectl get ingress -n <namespace>
   ```

3. **Smoke test** the health endpoint (if accessible):
   ```bash
   curl -sf https://<host>/health
   ```

## Rollback

If deployment fails, Helm `--atomic` auto-rolls back. If post-deploy checks fail, offer to run:
```bash
helm rollback goravel-blog -n <namespace>
```

## Summary

Output a deployment summary table:
| Item | Value |
|------|-------|
| Environment | staging/production |
| Version | x.y.z |
| Image | $DOCKER_REGISTRY/<username>/$DOCKER_IMAGE_NAME:tag |
| Namespace | books-xxx |
| Status | SUCCESS/FAILED |
| URL | https://... |
