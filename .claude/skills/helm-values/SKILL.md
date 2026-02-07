---
name: helm-values
description: Update Helm chart values, bump versions, validate templates, or compare staging vs production config. Use when modifying deployment configuration.
argument-hint: <bump-version|validate|diff|update> [args...]
disable-model-invocation: true
allowed-tools: Read, Edit, Write, Bash, Grep, Glob
---

# Helm Chart Management

Action: **$ARGUMENTS**

Parse the arguments:
- `bump-version [major|minor|patch]` - Bump the app version
- `validate [staging|production|all]` - Lint and render templates
- `diff` - Compare staging vs production values
- `update <key> <value> [staging|production]` - Update a Helm value

The chart is at `helm/goravel-blog/`.

## Actions

### `bump-version`

1. Read current version from `package.json`
2. Increment based on argument (default: `patch`):
   - major: 1.0.0 → 2.0.0
   - minor: 1.0.0 → 1.1.0
   - patch: 1.0.0 → 1.0.1
3. Update `package.json` version field
4. Sync `helm/goravel-blog/Chart.yaml` appVersion to match
5. Show the old → new version

### `validate`

For the target environment (default: `all`):

1. **Lint the chart:**
   ```bash
   helm lint helm/goravel-blog -f helm/goravel-blog/values.<env>.yaml
   ```

2. **Render templates** and check for issues:
   ```bash
   helm template goravel-blog helm/goravel-blog -f helm/goravel-blog/values.<env>.yaml > /tmp/rendered-<env>.yaml
   ```

3. **Validate rendered YAML** (check for empty values, missing secrets, resource limits):
   - Deployment has resource limits set
   - Ingress has a host defined
   - All required env vars are present in ConfigMap
   - Secrets reference valid base64 values

4. **Report findings** as a checklist:
   - [ ] Chart lint passes
   - [ ] Templates render without errors
   - [ ] Resource limits defined
   - [ ] Ingress configured
   - [ ] Health checks defined
   - [ ] Security context set

### `diff`

Compare staging and production configurations side-by-side:

1. Read `values.staging.yaml` and `values.production.yaml`
2. Show a comparison table of key differences:

| Setting | Staging | Production |
|---------|---------|------------|
| replicas | 1 | 3 |
| resources.limits.cpu | 500m | 1000m |
| resources.limits.memory | 512Mi | 1Gi |
| autoscaling.enabled | false | true |
| debug | true | false |
| ... | ... | ... |

3. Flag any security concerns (debug=true in prod, weak passwords, etc.)

### `update`

Update a specific value in a values file:

1. Parse: `update <dotted.key.path> <value> [environment]`
   - Environment defaults to both staging and production
2. Read the target values file(s)
3. Find and update the key
4. Re-validate the chart after the change
5. Show what changed

Example: `/helm-values update replicas 3 production`

## Safety Checks

For all operations:
- Never modify `values.yaml` (base defaults) without explicit confirmation
- Warn if production values have debug=true
- Warn if any passwords are hardcoded (not referencing secrets)
- Warn if resource limits are missing
