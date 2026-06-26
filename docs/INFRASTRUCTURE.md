# Infrastructure

## Local development (Docker Compose)

Placeholder at `infra/docker-compose.yml` — services are commented until each microservice exists. Uncomment incrementally per [STACK.md](STACK.md).

**UC-001 E2E stack** (Postgres + attendance HTTP + biometric gRPC): `infra/docker-compose.e2e.yml`

```bash
docker compose -f infra/docker-compose.e2e.yml up -d --wait
./scripts/verify-uc001-e2e.sh --curl
docker compose -f infra/docker-compose.e2e.yml down
```

Target stack (when enabled):

```yaml
# infra/docker-compose.yml
services:
  postgres:
    image: timescale/timescaledb-ha:pg16
  redis:
    image: valkey/valkey:7.2
  nats:
    image: nats:2.10-alpine
    command: -js
  api-gateway:
    build: ../services/api-gateway
  attendance:
    build: ../services/attendance
  biometric-service:
    build: ../services/biometric
```

Verify repo layout before enabling containers:

```bash
./scripts/verify-scaffold.sh
```

Ports (dev): Postgres 5432, Redis 6379, NATS 4222, API 8080, Biometric gRPC 9090.

## Production

```
GitHub Actions CI → container images → Helm charts (infra/k8s/) → Kubernetes
```

Optional IaC: `infra/terraform/` for cloud resources.

## Observability

| Signal | Stack |
|--------|-------|
| Traces | OpenTelemetry → Tempo |
| Logs | Loki |
| Metrics | Prometheus → Grafana |

### Required punch span attributes

`employee_id`, `tenant_id`, `punch_type`, `liveness_score`, `fraud_flags_count`, `geofence_id` — **never** raw biometric data.

### Prometheus metrics

- `openpresence_punch_total{status, tenant, fraud_type}`
- `openpresence_biometric_latency_seconds`
- `openpresence_liveness_score_histogram`

## Reliability

- Circuit breaker: API Gateway → Biometric Service
- NATS JetStream retries: max 3, exponential backoff
- Health: `/health/live`, `/health/ready` on every service
- Graceful shutdown: drain connections before exit

## Performance targets

| Metric | Target |
|--------|--------|
| Punch endpoint P99 | < 500ms (incl. biometric) |
| Embedding cache (Redis) | TTL 5min per tenant |
| pgvector IVFFlat | `lists=100` up to ~1M embeddings |
