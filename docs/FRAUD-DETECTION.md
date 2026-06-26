# Fraud Detection

Three layers: device (mobile), server (Go), biometric (Rust). See [FRAUD-DETECTION.md](FRAUD-DETECTION.md) fraud matrix for responses.

## Layer 1 — Device (Kotlin)

Checks run before camera opens. Results sent as `DeviceIntegrityReport`:

| Check | Flag |
|-------|------|
| Mock GPS (`Location.isMock` / `isFromMockProvider`) | MOCK_GPS |
| Developer options enabled | logged in report |
| Root binaries present | DEVICE_ROOTED |
| VPN active transport | VPN_DETECTED |

Critical device flags may block UI; warnings are recorded but may allow continue per policy.

## Layer 2 — Server (Go)

| Detection | Rule |
|-----------|------|
| Clock manipulation | `\|device_time - server_time\| > 5min` → FraudFlag; > 30min → CRITICAL reject |
| Impossible speed | > 600 km/h between consecutive punches → CRITICAL |
| IP anomaly | Same IP, different employees, < 30s → SUSPICIOUS |
| Duplicate punch | Valid punch within 60s → DUPLICATE_PUNCH |
| Out of geofence | Not inside any assigned zone → REJECTED |

## Layer 3 — Biometric (Rust)

| Detection | Rule |
|-----------|------|
| Liveness failed | ensemble score < 0.80 |
| Face not recognized | max cosine similarity < 0.60 |
| Adaptive threshold | After 5+ punches: `threshold = max(0.65, mean - 2*std)` |

## Fraud response matrix

| Fraud | Severity | Auto action | Notification |
|-------|----------|-------------|----------------|
| Mock GPS | HIGH | SUSPICIOUS | Manager |
| Clock manipulation > 5min | MEDIUM | SUSPICIOUS | HR |
| Clock manipulation > 30min | CRITICAL | REJECTED | HR + leadership |
| Liveness failed | HIGH | Block attempt | Manager |
| Face not recognized (3×) | HIGH | Device lock 30min | Manager + security |
| Impossible speed | CRITICAL | REJECTED + account lock | HR + legal |
| Device rooted | MEDIUM | SUSPICIOUS | IT |
| VPN active | LOW | Flag only | None |
| Duplicate IP < 30s | MEDIUM | SUSPICIOUS | HR |
| Out of geofence | HIGH | REJECTED | Manager |
| GPS low accuracy | LOW | Flag on VALID | None |

Full logic references: [BUSINESS-RULES.md](BUSINESS-RULES.md) BR-010–BR-022.
