---
name: infra-debug
description: Debug infrastructure and deployment issues. Inspect pods, check logs, view events, diagnose ingress/networking, and troubleshoot Kubernetes or Docker problems.
argument-hint: [staging|production] [pods|logs|events|ingress|describe|all]
disable-model-invocation: true
allowed-tools: Bash, Read, Grep
---

# Infrastructure Debugging

Debug the deployment for **$ARGUMENTS**.

Parse arguments:
- First argument: environment (`staging` or `production`). Default to `staging`.
- Second argument: what to inspect (`pods`, `logs`, `events`, `ingress`, `describe`, `all`). Default to `all`.

Namespace mapping:
- staging → `books-staging`
- production → `books-production`

## Debug Commands

### `pods` - Pod Status
```bash
kubectl get pods -n <namespace> -l app.kubernetes.io/name=goravel-blog -o wide
kubectl top pods -n <namespace> -l app.kubernetes.io/name=goravel-blog 2>/dev/null
```
Flag any pod that is NOT `Running` or has restarts > 0.

### `logs` - Application Logs
For each pod:
```bash
kubectl logs -n <namespace> <pod-name> --tail=100
```
If a pod has restarted, also get previous logs:
```bash
kubectl logs -n <namespace> <pod-name> --previous --tail=50
```
Look for error patterns: `panic`, `fatal`, `error`, `failed to connect`, `timeout`.

### `events` - Cluster Events
```bash
kubectl get events -n <namespace> --sort-by='.lastTimestamp' | tail -30
```
Flag any Warning events and explain likely causes.

### `ingress` - Ingress & Networking
```bash
kubectl get ingress -n <namespace>
kubectl describe ingress -n <namespace> goravel-blog
kubectl get svc -n <namespace>
kubectl get endpoints -n <namespace>
```
Check:
- Ingress has an address assigned
- Backend endpoints exist and are healthy
- TLS certificate is valid (if applicable)

### `describe` - Full Resource Description
```bash
kubectl describe deployment goravel-blog -n <namespace>
kubectl describe pod -n <namespace> -l app.kubernetes.io/name=goravel-blog
```
Look for:
- Image pull errors
- Resource limit issues (OOMKilled)
- Scheduling failures (insufficient resources, node selectors)
- Init container failures (migration errors)

### `all` - Full Diagnostic
Run all of the above in sequence, then provide a summary:

| Check | Status | Details |
|-------|--------|---------|
| Pods Running | OK/FAIL | X/Y pods ready |
| Pod Restarts | OK/WARN | Total restart count |
| Recent Errors | OK/WARN | Error count in logs |
| Ingress | OK/FAIL | URL and status |
| Events | OK/WARN | Warning count |
| Resources | OK/WARN | CPU/Memory utilization |

## Common Issues & Fixes

Based on findings, suggest specific fixes:
- **CrashLoopBackOff**: Check logs for panic/connection errors. Likely DB or Redis unreachable.
- **ImagePullBackOff**: Image tag doesn't exist. Check if CI built it.
- **OOMKilled**: Increase memory limits in values file.
- **Init container failed**: Migration error. Check init container logs.
- **No endpoints**: Service selector doesn't match pod labels.
- **503 errors**: No healthy pods behind the ingress.

## Helm Status
Also show:
```bash
helm status goravel-blog -n <namespace>
helm history goravel-blog -n <namespace> --max 5
```
