# Product Overview

## Purpose

OpenPresence solves fraud and inefficiency in employee time tracking for on-site and field workers. The platform supports two primary customer profiles:

| Profile | Use cases |
|---------|-----------|
| **Private companies** | CLT employees, contractors, service providers by physical site or geographic area |
| **Public sector / municipalities** | Civil servants, cooperatives, outsourced staff across secretariats and locations |

## Product assumptions

1. Employees are **pre-registered** by an administrator with photo and biometric data.
2. Punch is allowed only inside a **pre-defined geographic area** (geofence: circle or polygon).
3. Every punch requires **facial recognition with active liveness detection**.
4. Fraud attempts are detected and recorded — never silently ignored; everything is auditable.
5. The mobile app works **100% offline** for recognition and local registration, syncing when connectivity is available.

## Core capabilities

- Biometric enrollment (multi-angle face capture)
- Clock in / out, break start / end with sequence validation
- Geofence validation (circle and polygon)
- Multi-layer fraud detection (device, server, biometric)
- Organizational hierarchy (companies and municipalities)
- Suspicious punch review workflow
- Offline punch queue with policy-based TTL
- Immutable audit trail (TimescaleDB)

## Non-goals (v1)

- High-quality 3D mask and deepfake detection beyond liveness ensemble
- Cloud-hosted biometric APIs (self-hosted only — LGPD Art. 11)
- Employee self-enrollment without administrator presence
