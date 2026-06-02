# MyAnimeAPI — Frontend

Vite + React 18 + TypeScript + Tailwind SPA that consumes the Go API.

## Stack

- **Vite 5** — build + dev server
- **React 18** with React Router 6
- **TanStack Query 5** for server state
- **axios** with JWT interceptor for HTTP
- **React Hook Form** + **Zod** for forms
- **Tailwind 3** for styling
- **Playwright** for end-to-end tests

## Scripts

```bash
npm run dev        # vite dev server on :5173, proxies /v1 to API
npm run build      # type-check + production build to dist/
npm run preview    # serve dist/ on :4173 (used by CI for Playwright)
npm run lint       # eslint
npm run typecheck  # tsc --noEmit
npm run test:e2e   # playwright tests against PLAYWRIGHT_BASE_URL
```

## Environment

`VITE_API_URL` (default `http://localhost:8080`) — base URL the dev
proxy and the axios client point at. In dev the proxy keeps requests
same-origin so cookies/CORS behave; in production the built bundle reads
from `import.meta.env.VITE_API_URL` directly.

## Layout

```
src/
├── components/    # presentational + small interactive bits
├── lib/
│   ├── api/       # one file per resource (auth, animes, favorites, ...)
│   └── auth/      # AuthProvider context + helpers
├── pages/         # one file per route
└── routes/        # router + ProtectedRoute
tests/e2e/         # Playwright specs
```
