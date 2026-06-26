# OpenPresence Admin Panel

Web administrative UI built with [TanStack Start](https://tanstack.com/start), [Router](https://tanstack.com/router), and [Query](https://tanstack.com/query).

## Prerequisites

```bash
./scripts/dev-backend.sh start   # API at http://127.0.0.1:8088
```

## Setup

```bash
cd web/admin
cp .env.example .env.local
npm install
npm run dev
```

Open http://localhost:5174

## Scripts

| Command | Description |
|---------|-------------|
| `npm run dev` | Dev server on port 5174 |
| `npm run build` | Production build + typecheck |
| `npm run preview` | Preview production build |

## Environment

| Variable | Default | Description |
|----------|---------|-------------|
| `VITE_API_BASE_URL` | `http://127.0.0.1:8088` | Attendance / API gateway base URL |
| `VITE_AUTH_MOCK` | `true` | Dev mock login until auth service exists |

## Verification

```bash
./scripts/verify-admin-scaffold.sh
```
