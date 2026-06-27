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

Mock login (when `VITE_AUTH_MOCK=true`):

| Registration | Password | Role |
|--------------|----------|------|
| `admin` | `admin` | ORG_ADMIN |
| `manager` | `manager` | MANAGER |
| `hr` | `hr` | HR_ANALYST |
| `auditor` | `auditor` | AUDITOR |

Credentials are listed on the login screen in dev mock mode. DB seed: `./scripts/seed-dev-users.sh`.

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
./scripts/verify-admin-guards.sh
./scripts/verify-admin.sh
```

## Login (`/login`)

TanStack Form with `registrationId` + `password`. Mock: `admin` / `admin`. Redirects to `/dashboard` (or `?redirect=`).

## Router guards (admin-04)

Protected routes live under pathless layout `_authenticated` (`src/routes/_authenticated.tsx`). `beforeLoad` checks `context.auth.isAuthenticated` and redirects guests to `/login?redirect=...`. `/` redirects to `/dashboard` or `/login`; `/login` redirects authenticated users away.

```bash
./scripts/verify-admin-guards.sh
```

## Admin shell (admin-05)

Post-login layout via `AdminShell`: sidebar (Dashboard + future modules), header with user/tenant, sign out. Dashboard polls `GET /health/live` with TanStack Query.

| Module | Purpose |
|--------|---------|
| `src/components/AdminShell.tsx` | Sidebar, header, content slot |
| `src/lib/api/health.ts` | Live health probe |

```bash
./scripts/verify-admin.sh
```

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

Browser requests require CORS on the attendance API. Local dev sets `CORS_ALLOWED_ORIGINS=http://localhost:5174,http://127.0.0.1:5174` via `./scripts/dev-backend.sh`.

## Verification

```bash
./scripts/verify-admin-scaffold.sh
```
