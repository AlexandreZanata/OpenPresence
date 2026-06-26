# OpenPresence Admin Panel

Web administrative UI built with [TanStack Start](https://tanstack.com/start), [Router](https://tanstack.com/router), and [Query](https://tanstack.com/query).

## Prerequisites

```bash
./scripts/dev-backend.sh start   # API at http://127.0.0.1:8088
```

## Setup

```bash
cd web/admin
npm install          # required — node_modules is gitignored
cp .env.example .env.local
npm run dev
```

Open http://localhost:5174

Mock login (when `VITE_AUTH_MOCK=true`): registration `admin` / password `admin`.

## Auth (admin-02)

| Module | Purpose |
|--------|---------|
| `src/lib/auth/AuthProvider.tsx` | React context + `useAuth()` |
| `src/lib/auth/storage.ts` | Session in `localStorage` (dev only) |
| `src/lib/auth/dev-mock.ts` | Mock login until auth service |
| `src/lib/api/client.ts` | `fetch` + Bearer + API error envelope |
| `src/client.tsx` | Injects `auth` into `RouterProvider` context |

```bash
npm run auth-smoke
./scripts/verify-admin-auth.sh
./scripts/verify-admin-login.sh
```

## Login (`/login`)

TanStack Form with `registrationId` + `password`. Mock: `admin` / `admin`. Redirects to `/dashboard` (or `?redirect=`).

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
