# Goravel Blog Helm Chart

This Helm chart deploys the Goravel Blog application on a Kubernetes cluster.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.0+
- PV provisioner support in the underlying infrastructure (for storage persistence)
- MySQL/PostgreSQL database (external or in-cluster)
- Redis (optional, for caching)

## Installation

### Add the repository (if hosted)

```bash
helm repo add goravel-blog https://your-helm-repo.com
helm repo update
```

### Install the chart

```bash
helm install goravel-blog ./helm/goravel-blog \
  --namespace blog \
  --create-namespace \
  --set app.database.password=your-db-password \
  --set app.jwt.secret=your-jwt-secret
```

### Install with custom values

```bash
helm install goravel-blog ./helm/goravel-blog \
  --namespace blog \
  --create-namespace \
  -f custom-values.yaml
```

## Configuration

The following table lists the configurable parameters and their default values.

### General Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `replicaCount` | Number of replicas | `2` |
| `image.repository` | Image repository | `goravel-blog` |
| `image.tag` | Image tag | `""` (uses appVersion) |
| `image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `nameOverride` | Override chart name | `""` |
| `fullnameOverride` | Override full name | `""` |

### Service Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `service.type` | Kubernetes service type | `ClusterIP` |
| `service.port` | Service port | `80` |
| `service.targetPort` | Container port | `3000` |

### Ingress Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `ingress.enabled` | Enable ingress | `true` |
| `ingress.className` | Ingress class name | `nginx` |
| `ingress.hosts[0].host` | Hostname | `blog.example.com` |
| `ingress.tls` | TLS configuration | `[]` |

### Application Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `app.name` | Application name | `Goravel Blog` |
| `app.env` | Environment | `production` |
| `app.debug` | Debug mode | `false` |
| `app.key` | Application key | Auto-generated |
| `app.url` | Application URL | `https://blog.example.com` |

### Database Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `app.database.connection` | Database type | `mysql` |
| `app.database.host` | Database host | `mysql.default.svc.cluster.local` |
| `app.database.port` | Database port | `3306` |
| `app.database.database` | Database name | `goravel_blog` |
| `app.database.username` | Database username | `goravel` |
| `app.database.password` | Database password | `""` (must be set) |

### Resources

| Parameter | Description | Default |
|-----------|-------------|---------|
| `resources.limits.cpu` | CPU limit | `500m` |
| `resources.limits.memory` | Memory limit | `512Mi` |
| `resources.requests.cpu` | CPU request | `250m` |
| `resources.requests.memory` | Memory request | `256Mi` |

### Autoscaling

| Parameter | Description | Default |
|-----------|-------------|---------|
| `autoscaling.enabled` | Enable HPA | `true` |
| `autoscaling.minReplicas` | Minimum replicas | `2` |
| `autoscaling.maxReplicas` | Maximum replicas | `10` |
| `autoscaling.targetCPUUtilizationPercentage` | Target CPU utilization | `80` |
| `autoscaling.targetMemoryUtilizationPercentage` | Target memory utilization | `80` |

## Persistence

The chart mounts a Persistent Volume for storage. You can configure it with:

```yaml
app:
  storage:
    persistence:
      enabled: true
      storageClass: "standard"
      accessMode: ReadWriteOnce
      size: 10Gi
```

## Database Migrations

After installation, run database migrations:

```bash
kubectl exec -n blog deployment/goravel-blog -- ./main artisan migrate
```

## Upgrading

```bash
helm upgrade goravel-blog ./helm/goravel-blog \
  --namespace blog \
  --reuse-values \
  --set image.tag=new-version
```

## Uninstalling

```bash
helm uninstall goravel-blog --namespace blog
```

This will remove all the Kubernetes components associated with the chart and delete the release.

## Troubleshooting

### Check pod status
```bash
kubectl get pods -n blog -l app.kubernetes.io/name=goravel-blog
```

### View logs
```bash
kubectl logs -n blog -l app.kubernetes.io/name=goravel-blog
```

### Describe deployment
```bash
kubectl describe deployment -n blog goravel-blog
```