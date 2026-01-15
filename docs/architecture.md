# LauraTech: Comprehensive Technical Architecture

**Document Version:** 1.0  
**Last Updated:** 2026-01-14  
**Status:** Active  
**Owner:** Senior Principal Software Architect  

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [System Architecture Overview](#system-architecture-overview)
3. [Frontend Architecture (Next.js 16.1)](#frontend-architecture-nextjs-161)
4. [Authentication & Identity (Ory Integration)](#authentication--identity-ory-integration)
5. [Data Fetching & State Management](#data-fetching--state-management)
6. [Security Architecture](#security-architecture)
7. [Observability & Monitoring](#observability--monitoring)
8. [Deployment & Infrastructure](#deployment--infrastructure)
9. [Decision Records Reference](#decision-records-reference)

---

## Executive Summary

LauraTech is a fintech platform built on **Next.js 16.1** with **Server Components** as the primary rendering paradigm. The architecture emphasizes:

- **Zero-Trust Security:** All routes protected by Ory Elements middleware; no unauthenticated access to sensitive data.
- **Server-Driven Architecture:** Leverage React Server Components (RSCs) to minimize client-side JavaScript, reduce bundle size, and enable secure server-side operations.
- **Domain-Driven Design:** Frontend organized into isolated business domains (`auth`, `payments`, `dashboard`) enabling independent scaling and feature teams.
- **Fintech-Grade Observability:** Sentry for error tracking, OpenTelemetry for performance metrics, structured logging for audit trails.
- **Type Safety:** Zod schemas for all API boundaries; TypeScript strict mode throughout.

---

## System Architecture Overview

### High-Level Data Flow Diagram

```
┌─────────────────────────────────────────────────────────────────────┐
│                          CLIENT BROWSER                             │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │  React Client Components (Interactive Forms, Charts, Real-   │  │
│  │  time Updates)                                               │  │
│  └────────────────┬─────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────┘
                     │
         ┌───────────┴──────────────┬────────────────┐
         │                          │                │
    Server Action                Fetch             WebSocket
    + Zod Validation            (GET/POST)         (Optional)
         │                          │                │
┌────────▼──────────────────────────▼────────────────▼────────────────┐
│                      NEXT.JS SERVER LAYER                           │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │ Middleware: Ory Token Validation + Session Verification       │ │
│  └────────────────────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │ RSC (React Server Components) - Protected Routes              │ │
│  │  • Page Components                                             │ │
│  │  • Data Fetching (with Next.js cache invalidation)            │ │
│  │  • Server Actions (authenticated mutations)                   │ │
│  └────────────────────────────────────────────────────────────────┘ │
└────────────────────────────┬──────────────────────────────────────────┘
                     │
         ┌───────────┴──────────────┬────────────────┐
         │                          │                │
      Ory API               Backend Gateway         Cache
      (OIDC/PKCE)          (APISIX)                (Redis)
         │                          │                │
┌────────▼──────────────────────────▼────────────────▼────────────────┐
│                    MICROSERVICES BACKEND                            │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │ Ory Kratos: Identity Management (Login, MFA, Recovery)        │ │
│  └────────────────────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │ Core Service (Go/Gin): Ledger, Account Management, Profiles   │ │
│  └────────────────────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │ Payments Service: Transaction Orchestration, Settlement        │ │
│  └────────────────────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │ Apache APISIX Gateway: Rate Limiting, Auth, Request Logging    │ │
│  └────────────────────────────────────────────────────────────────┘ │
└────────────────────────────┬──────────────────────────────────────────┘
                     │
         ┌───────────┴──────────────┬────────────────┐
         │                          │                │
      PostgreSQL              Redis              S3/Blob
      (Ledger)               (Cache)            (Documents)
```

---

## Frontend Architecture (Next.js 16.1)

### Core Principles

1. **React Server Components as Default:** All new pages/routes are RSCs unless explicit interactivity required.
2. **Minimize Client JavaScript:** Use `"use client"` sparingly; push logic to server.
3. **Streaming & Progressive Enhancement:** Leverage `suspense` + Server Components for streaming HTML responses.
4. **Type-Safe Props:** All component interfaces fully typed; no `any` types.

### Page Lifecycle: Protected Route Example

```typescript
// src/app/(dashboard)/payments/page.tsx

// 1. RSC - Server-side only
// Executed on server; no client-side JS footprint
export default async function PaymentsPage({
  searchParams,
}: {
  searchParams: { page?: string };
}) {
  // 2. Verify session via Ory (middleware guarantees this)
  const session = await getOrySession();
  
  // 3. Fetch data with Next.js cache directives
  const transactions = await fetchUserTransactions({
    userId: session.identity.id,
    page: parseInt(searchParams.page || "1"),
  });

  return (
    <div className="space-y-6">
      {/* Static content */}
      <Header title="Payment History" />
      
      {/* Interactive component - requires client hydration */}
      <Suspense fallback={<TransactionsSkeleton />}>
        <TransactionsList initialData={transactions} />
      </Suspense>
      
      {/* Server Action - mutate state securely */}
      <ExportButton userId={session.identity.id} />
    </div>
  );
}

// 4. Revalidation strategy
export const revalidateTransactions = () => {
  revalidatePath("/dashboard/payments");
  revalidateTag("user-transactions");
};
```

### Client Component: Interactivity

```typescript
// src/shared/components/TransactionsList.tsx
"use client";

import { useTransition } from "react";
import { reexportTransactions } from "@/modules/payments/actions";

export function TransactionsList({ initialData }) {
  const [transactions, setTransactions] = useState(initialData);
  const [isPending, startTransition] = useTransition();

  const handleFilter = (filter: PaymentFilter) => {
    startTransition(async () => {
      // Server Action - securely filtered on backend
      const filtered = await fetchFilteredTransactions(filter);
      setTransactions(filtered);
    });
  };

  return (
    <div>
      <FilterUI onFilter={handleFilter} />
      <Table data={transactions} isLoading={isPending} />
    </div>
  );
}
```

### Folder Structure Rationale

```
src/
├── app/                           # Next.js App Router (routes + layouts)
│   ├── (auth)/                   # Route group for public auth pages
│   │   ├── login/page.tsx
│   │   ├── registration/page.tsx
│   │   └── layout.tsx
│   ├── (dashboard)/              # Route group for protected dashboard
│   │   ├── layout.tsx            # Shared dashboard layout (sidebar, nav)
│   │   ├── page.tsx              # Dashboard home
│   │   ├── payments/
│   │   │   └── page.tsx          # Transaction history (RSC)
│   │   ├── accounts/
│   │   │   └── page.tsx          # Account management
│   │   └── profile/
│   │       └── page.tsx          # User settings
│   ├── api/                      # API routes (webhooks, external integrations)
│   │   └── webhooks/
│   │       └── payment/route.ts
│   ├── error.tsx                 # Error boundary
│   ├── not-found.tsx
│   ├── globals.css
│   └── layout.tsx                # Root layout + providers
│
├── modules/                       # Domain-specific business logic
│   ├── auth/
│   │   ├── components/           # Auth-specific UI (LoginForm, etc)
│   │   ├── actions/              # Server Actions for auth
│   │   ├── hooks/                # Auth hooks (useSession, useUser)
│   │   ├── types.ts              # Auth domain types
│   │   └── validators.ts         # Zod schemas for auth
│   ├── payments/
│   │   ├── components/
│   │   ├── actions/              # Server Actions (executeTransfer, etc)
│   │   ├── hooks/
│   │   ├── types.ts
│   │   └── validators.ts
│   └── dashboard/
│       ├── components/
│       ├── types.ts
│       └── validators.ts
│
├── shared/                       # Reusable across modules
│   ├── components/              # UI atoms (Button, Input, Card, etc)
│   │   ├── ui/
│   │   │   ├── Button.tsx
│   │   │   ├── Input.tsx
│   │   │   ├── Card.tsx
│   │   │   ├── Modal.tsx
│   │   │   └── Table.tsx
│   │   ├── layouts/
│   │   │   ├── DashboardLayout.tsx
│   │   │   └── AuthLayout.tsx
│   │   └── feedback/
│   │       ├── Toast.tsx
│   │       ├── Skeleton.tsx
│   │       └── EmptyState.tsx
│   ├── hooks/                   # Generic hooks (useAsync, useDebounce, etc)
│   │   ├── useAsync.ts
│   │   └── useDebounce.ts
│   ├── utils/                   # Utilities (formatters, parsers, etc)
│   │   ├── currency.ts          # Format BRL amounts
│   │   ├── date.ts
│   │   └── string.ts
│   └── types/                   # Shared domain types
│       └── index.ts
│
├── core/                        # Infrastructure & configuration
│   ├── api/                     # HTTP client for backend
│   │   ├── client.ts           # Typed fetch wrapper
│   │   └── endpoints.ts         # API route definitions
│   ├── ory/                     # Ory integration
│   │   ├── client.ts           # Ory SDK instance
│   │   ├── middleware.ts       # Middleware for token validation
│   │   └── session.ts          # Session verification utilities
│   ├── validators/              # Global validation schemas
│   │   ├── email.ts
│   │   ├── currency.ts
│   │   └── index.ts
│   ├── config/
│   │   ├── constants.ts         # App constants
│   │   ├── env.ts              # Environment variables (typed)
│   │   └── csp.ts              # Content Security Policy
│   └── telemetry/
│       ├── sentry.ts           # Sentry SDK setup
│       └── otel.ts             # OpenTelemetry instrumentation
│
└── lib/                         # Legacy utility exports (migrate to modules/)
    └── [deprecated]
```

---

## Authentication & Identity (Ory Integration)

### Ory Flow: OIDC PKCE

1. **User initiates login** → Browser redirects to Ory UI (hosted on separate domain)
2. **Ory handles credential validation** → Creates session JWT + refresh token
3. **Redirect to callback** → `/auth/callback?code=XXX&state=YYY`
4. **Next.js exchanges code** → Server Action calls Ory API → Returns session
5. **Session stored** → Secure HTTP-only cookie set by Ory
6. **Subsequent requests** → Middleware validates token, injects user context

### Middleware: Token Validation

```typescript
// middleware.ts
import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";
import { verifyOrySession } from "@/core/ory/session";

const PUBLIC_ROUTES = ["/auth/login", "/auth/registration"];

export async function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;

  // Allow public routes
  if (PUBLIC_ROUTES.includes(pathname)) {
    return NextResponse.next();
  }

  // Verify Ory session for protected routes
  try {
    const session = await verifyOrySession(request);
    
    if (!session) {
      // Redirect to login
      return NextResponse.redirect(new URL("/auth/login", request.url));
    }

    // Attach session to request headers for RSCs
    const response = NextResponse.next();
    response.headers.set("x-user-id", session.identity.id);
    response.headers.set("x-user-email", session.identity.traits.email);
    
    return response;
  } catch (error) {
    return NextResponse.redirect(new URL("/auth/login", request.url));
  }
}

export const config = {
  matcher: ["/((?!_next/static|_next/image|favicon.ico).*)"],
};
```

### Server Component: Session Access

```typescript
// src/core/ory/session.ts

export async function getOrySession() {
  const ory = new FrontendApi(
    new Configuration({
      basePath: process.env.ORY_SDK_URL,
      baseOptions: {
        withCredentials: true,
        headers: {
          "Content-Type": "application/json",
        },
      },
    })
  );

  try {
    const { data } = await ory.toSession();
    return data;
  } catch (error) {
    throw new Error("Session verification failed");
  }
}

// Usage in RSC
export default async function ProtectedPage() {
  const session = await getOrySession();
  
  return <div>Welcome, {session.identity.traits.email}</div>;
}
```

### Protection: Server Actions

```typescript
// src/modules/payments/actions.ts
"use server";

import { verifyOrySession } from "@/core/ory/session";
import { executePaymentSchema } from "./validators";

export async function executePayment(formData: unknown) {
  // 1. Authenticate request
  const session = await verifyOrySession();
  if (!session) throw new Error("Unauthorized");

  // 2. Validate input with Zod
  const parsed = executePaymentSchema.safeParse(formData);
  if (!parsed.success) {
    return { error: "Invalid input", details: parsed.error.flatten() };
  }

  // 3. Call backend API (Ory token included automatically)
  const response = await fetch(
    `${process.env.BACKEND_API_URL}/payments/execute`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${session.session_token}`,
      },
      body: JSON.stringify(parsed.data),
    }
  );

  if (!response.ok) {
    throw new Error("Payment failed");
  }

  // 4. Revalidate cache
  revalidatePath("/dashboard/payments");
  
  return { success: true };
}
```

---

## Data Fetching & State Management

### Strategy: Hybrid Approach

| Use Case | Method | Caching | Real-time |
|----------|--------|---------|-----------|
| Initial page load (read-only) | RSC + `fetch` | Next.js ISR/revalidate | No |
| Form submission (mutation) | Server Action | Manual revalidation | No |
| Complex client filters | TanStack Query | In-memory | Yes (optional polling) |
| Real-time updates (future) | WebSocket | N/A | Yes |

### Server-Side Data Fetching (Default)

```typescript
// src/modules/payments/actions.ts

export async function fetchUserTransactions(userId: string, page: number) {
  const response = await fetch(
    `${process.env.BACKEND_API_URL}/transactions?userId=${userId}&page=${page}`,
    {
      // Next.js cache directives
      next: { 
        revalidate: 60, // ISR: revalidate every 60 seconds
        tags: ["transactions"], // Invalidate via revalidateTag()
      },
      headers: {
        Authorization: `Bearer ${getSessionToken()}`,
      },
    }
  );

  if (!response.ok) throw new Error("Failed to fetch transactions");
  
  return response.json();
}
```

### Client-Side Mutation (Server Action)

```typescript
// src/modules/payments/actions.ts
"use server";

export async function createTransfer(payload: TransferPayload) {
  const session = await getOrySession();
  
  // Zod validation
  const validated = transferSchema.parse(payload);

  // API call
  const response = await fetch(
    `${process.env.BACKEND_API_URL}/transfers`,
    {
      method: "POST",
      headers: {
        "Authorization": `Bearer ${session.session_token}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify(validated),
    }
  );

  // Revalidate affected caches
  revalidatePath("/dashboard/payments");
  revalidateTag("user-transactions");

  return response.json();
}
```

### Complex Client State (TanStack Query)

```typescript
// src/modules/payments/hooks/useTransactions.ts
"use client";

import { useQuery } from "@tanstack/react-query";
import { fetchFilteredTransactions } from "../actions";

export function useTransactions(userId: string, filters: TransactionFilter) {
  return useQuery({
    queryKey: ["transactions", userId, filters],
    queryFn: async () => {
      return fetchFilteredTransactions(userId, filters);
    },
    staleTime: 30000, // 30 seconds
  });
}
```

---

## Security Architecture

### Defense-in-Depth Layers

#### 1. **Content Security Policy (CSP)**

```typescript
// src/core/config/csp.ts

export const CSP_HEADER = `
  default-src 'self';
  script-src 'self' 'nonce-{NONCE}';
  style-src 'self' 'unsafe-inline' fonts.googleapis.com;
  font-src 'self' fonts.gstatic.com;
  img-src 'self' data: https:;
  connect-src 'self' ${process.env.ORY_SDK_URL} ${process.env.BACKEND_API_URL};
  frame-ancestors 'none';
  base-uri 'self';
  form-action 'self';
  upgrade-insecure-requests;
`;
```

#### 2. **Input Validation (Zod)**

```typescript
// src/modules/payments/validators.ts

import { z } from "zod";

export const transferSchema = z.object({
  recipientId: z.string().uuid("Invalid recipient ID"),
  amount: z
    .number()
    .positive("Amount must be positive")
    .max(1000000, "Amount exceeds limit"),
  description: z.string().max(500).optional(),
  metadata: z.record(z.string()).optional(),
});

export type Transfer = z.infer<typeof transferSchema>;
```

#### 3. **CORS & CSRF Protection**

```typescript
// next.config.ts

export default {
  headers: async () => [
    {
      source: "/api/:path*",
      headers: [
        {
          key: "Access-Control-Allow-Origin",
          value: process.env.FRONTEND_URL,
        },
        {
          key: "Access-Control-Allow-Methods",
          value: "GET, POST, PUT, DELETE",
        },
        {
          key: "X-Content-Type-Options",
          value: "nosniff",
        },
        {
          key: "X-Frame-Options",
          value: "DENY",
        },
      ],
    },
  ],
};
```

#### 4. **Rate Limiting (APISIX Gateway)**

All API calls routed through Apache APISIX which enforces:
- Per-IP rate limiting: 1000 req/min
- Per-user rate limiting: 100 req/min (sensitive operations)
- Token bucket algorithm for burst protection

#### 5. **Secrets Management**

```typescript
// src/core/config/env.ts

import { z } from "zod";

const envSchema = z.object({
  ORY_SDK_URL: z.string().url(),
  BACKEND_API_URL: z.string().url(),
  SENTRY_DSN: z.string().url().optional(),
  OTEL_EXPORTER_OTLP_ENDPOINT: z.string().url().optional(),
});

export const env = envSchema.parse(process.env);

// Never expose to client:
// - API keys
// - Database credentials
// - Signing keys
```

---

## Observability & Monitoring

### Three Pillars

#### 1. **Error Tracking (Sentry)**

```typescript
// src/core/telemetry/sentry.ts

import * as Sentry from "@sentry/nextjs";

export function initSentry() {
  Sentry.init({
    dsn: process.env.SENTRY_DSN,
    environment: process.env.NODE_ENV,
    tracesSampleRate: process.env.NODE_ENV === "production" ? 0.1 : 1,
    beforeSend(event) {
      // Don't send auth errors (handled separately)
      if (event.tags?.auth === "true") {
        return null;
      }
      return event;
    },
  });
}

// Capture exceptions in Server Actions
export async function executePayment(data: unknown) {
  try {
    // ...payment logic
  } catch (error) {
    Sentry.captureException(error, {
      tags: { module: "payments" },
      extra: { userId: session.identity.id },
    });
    throw error;
  }
}
```

#### 2. **Performance Metrics (OpenTelemetry)**

```typescript
// src/core/telemetry/otel.ts

import { getNodeAutoInstrumentations } from "@opentelemetry/auto-instrumentations-node";
import { NodeSDK } from "@opentelemetry/sdk-node";

export const sdk = new NodeSDK({
  traceExporter: new OTLPTraceExporter({
    url: process.env.OTEL_EXPORTER_OTLP_ENDPOINT,
  }),
  instrumentations: [getNodeAutoInstrumentations()],
});

sdk.start();

// Instrument slow operations
export const tracer = opentelemetry.trace.getTracer("laura-tech");

export function traceServerAction<T>(
  name: string,
  fn: () => Promise<T>
): Promise<T> {
  return tracer.startActiveSpan(name, (span) => {
    return fn()
      .then((result) => {
        span.setStatus({ code: SpanStatusCode.OK });
        return result;
      })
      .catch((error) => {
        span.setStatus({
          code: SpanStatusCode.ERROR,
          message: error.message,
        });
        throw error;
      })
      .finally(() => span.end());
  });
}
```

#### 3. **Structured Logging**

```typescript
// src/core/telemetry/logger.ts

export interface LogContext {
  userId?: string;
  requestId: string;
  module: string;
  action: string;
}

export function logPaymentSuccess(context: LogContext, amount: number) {
  console.log(
    JSON.stringify({
      level: "info",
      timestamp: new Date().toISOString(),
      event: "payment_executed",
      amount,
      ...context,
    })
  );
}

export function logPaymentFailure(
  context: LogContext,
  error: Error,
  reason: string
) {
  console.error(
    JSON.stringify({
      level: "error",
      timestamp: new Date().toISOString(),
      event: "payment_failed",
      reason,
      error: {
        message: error.message,
        stack: error.stack,
      },
      ...context,
    })
  );
}
```

---

## Deployment & Infrastructure

### Build Pipeline

```dockerfile
# Dockerfile (Multi-stage, Bun-based)

FROM oven/bun:latest as builder
WORKDIR /app
COPY package.json bun.lockb ./
RUN bun install --frozen-lockfile
COPY . .
RUN bun run build

FROM oven/bun:slim
WORKDIR /app
COPY --from=builder /app/.next .next
COPY --from=builder /app/public public
COPY package.json bun.lockb ./
RUN bun install --frozen-lockfile --production

EXPOSE 3000
CMD ["bun", "run", "start"]
```

### Docker Compose (Local Development)

```yaml
# docker-compose.yml
version: "3.8"
services:
  frontend:
    build:
      context: ./front
    ports:
      - "3000:3000"
    environment:
      ORY_SDK_URL: http://kratos:4433
      BACKEND_API_URL: http://apisix:9000
      NODE_ENV: development
    depends_on:
      - kratos
      - apisix

  kratos:
    image: oryd/kratos:latest
    command: serve -c /etc/kratos/kratos.yml
    volumes:
      - ./docker/kratos:/etc/kratos
    ports:
      - "4433:4433"
      - "4434:4434"

  apisix:
    image: apache/apisix:latest
    volumes:
      - ./docker/apisix:/etc/apisix
    ports:
      - "9000:9000"
    depends_on:
      - backend

  backend:
    build:
      context: ./back
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgresql://user:pass@postgres:5432/laura_tech
```

### Environment Variables (`.env.local`)

```bash
# Ory Configuration
ORY_SDK_URL=https://ory.lauratech.com
ORY_API_KEY=<secret-api-key>

# Backend API
BACKEND_API_URL=https://api.lauratech.com

# Observability
SENTRY_DSN=https://<sentry-key>@sentry.io/<project-id>
OTEL_EXPORTER_OTLP_ENDPOINT=https://otel-collector.lauratech.com

# Feature Flags
FEATURE_EXPORT_TRANSACTIONS=true
```

---

## Decision Records Reference

The following Architectural Decision Records (ADRs) document critical choices:

- **ADR-001:** Next.js 16.1 Adoption & RSC Strategy
- **ADR-002:** Ory Zero-Trust Identity Architecture
- **ADR-003:** Atomic Design + Server Components Componentization
- **ADR-004:** Tri-Layer Testing Strategy (Vitest, Testing Library, Playwright)

See `docs/adr/` directory for full details.

---

## Glossary

- **RSC:** React Server Component (server-only component)
- **CSP:** Content Security Policy (HTTP header for XSS prevention)
- **OIDC:** OpenID Connect (identity protocol)
- **PKCE:** Proof Key for Public Clients (secure OAuth 2.0 flow)
- **ISR:** Incremental Static Regeneration (Next.js caching strategy)
- **OTEL:** OpenTelemetry (observability framework)

---

**Document created by:** Senior Principal Software Architect  
**Last review:** 2026-01-14  
**Next review:** 2026-02-14
  - OpenTelemetry (Futuro) para tracing distribuído.
- **CI/CD:** GitHub Actions (Build, Test, Lint).

### 5. Pagamentos (Microserviço)
- Serviço isolado para lidar com gateways externos (Stripe, Adyen, Banco Central).
- Desenho desacoplado para permitir troca fácil de provedor.
