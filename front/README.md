# LauraTech Frontend: Next.js 16.1 Fintech Architecture

**Version:** 1.0  
**Last Updated:** 2026-01-14  
**Status:** MVP Foundation  

---

## 📋 Table of Contents

- [Overview](#overview)
- [Quick Start](#quick-start)
- [Architecture & Design](#architecture--design)
- [Project Structure](#project-structure)
- [Core Patterns](#core-patterns)
- [Authentication (Ory)](#authentication-ory)
- [Data Fetching Strategy](#data-fetching-strategy)
- [Security](#security)
- [Testing](#testing)
- [Development](#development)
- [Deployment](#deployment)
- [Troubleshooting](#troubleshooting)

---

## Overview

**LauraTech** is a fintech platform built on **Next.js 16.1** with React Server Components (RSCs) as the primary rendering paradigm. This README documents the architecture, conventions, and patterns used throughout the codebase.

### Tech Stack

```
Runtime:        Bun 1.0+
Framework:      Next.js 16.1 (App Router, RSC)
Language:       TypeScript 5 (strict mode)
Styling:        Tailwind CSS 4 + PostCSS
Auth:           Ory Kratos + Ory Elements
Validation:     Zod (schema-first validation)
Testing:        Vitest + Testing Library + Playwright
Observability:  Sentry + OpenTelemetry
```

### Core Principles

1. **Server-First:** React Server Components are default; `"use client"` only at interactive boundaries.
2. **Zero-Trust Security:** Every request verified; no implicit trust of sessions.
3. **Type Safety:** TypeScript strict mode; Zod schemas for all API boundaries.
4. **Atomic Design:** Components organized in Atoms → Molecules → Organisms → Containers.
5. **Domain-Driven Design:** Business logic isolated in `/src/modules/{domain}/`.

---

## Quick Start

### Prerequisites

- **Bun** 1.0+ installed ([bun.sh](https://bun.sh))
- **Node.js** 20+ (for tooling compatibility)
- **Docker** (for local Ory Kratos + APISIX)

### Setup

1. **Clone and install dependencies:**

```bash
cd fin/front
bun install
```

2. **Create environment file:**

```bash
cp .env.example .env.local
```

3. **Update `.env.local` with local development values:**

```bash
# Ory
ORY_SDK_URL=http://localhost:4433
ORY_API_KEY=

# Backend API
BACKEND_API_URL=http://localhost:9000

# Observability (optional for MVP)
SENTRY_DSN=
OTEL_EXPORTER_OTLP_ENDPOINT=

# Feature Flags
FEATURE_EXPORT_TRANSACTIONS=true
```

4. **Start Docker services** (Ory Kratos + APISIX):

```bash
cd ../docker
docker-compose up -d
```

5. **Start development server:**

```bash
bun run dev
```

Access at `http://localhost:3000`.

### Verification

- Navigate to `/auth/login` — should see Ory login form.
- Navigate to `/dashboard` — should redirect to login (protected route).
- (After login) Access dashboard and view transactions.

---

## Architecture & Design

### High-Level Data Flow

See [docs/architecture.md](../docs/architecture.md) for comprehensive architecture documentation.

### Architectural Decision Records (ADRs)

All major architecture decisions are documented in [docs/adr/](../docs/adr/):

| ADR | Topic |
|-----|-------|
| **ADR-001** | [Next.js 16.1 Adoption & RSC Strategy](../docs/adr/ADR-001-nextjs-rsc-adoption.md) |
| **ADR-002** | [Zero-Trust Identity with Ory](../docs/adr/ADR-002-ory-zero-trust.md) |
| **ADR-003** | [Atomic Design + Server Components](../docs/adr/ADR-003-atomic-server-components.md) |
| **ADR-004** | [Tri-Layer Testing Strategy](../docs/adr/ADR-004-tri-layer-testing.md) |

---

## Project Structure

### Directory Organization

```
src/
├── app/                           # Next.js App Router (routes + pages)
│   ├── (auth)/                   # Route group: public auth pages
│   │   ├── login/
│   │   │   └── page.tsx          # RSC: Login page with Ory Elements
│   │   ├── registration/
│   │   ├── recovery/
│   │   ├── verification/
│   │   └── layout.tsx            # Auth layout (no sidebar)
│   │
│   ├── (dashboard)/              # Route group: protected pages
│   │   ├── layout.tsx            # Dashboard layout with sidebar
│   │   ├── page.tsx              # Dashboard home
│   │   ├── payments/
│   │   │   └── page.tsx          # Transaction history (RSC)
│   │   ├── accounts/
│   │   └── profile/
│   │
│   ├── api/                      # API routes (webhooks, external integrations)
│   │   └── webhooks/
│   │       └── payment/route.ts
│   │
│   ├── error.tsx                 # Error boundary
│   ├── not-found.tsx
│   ├── layout.tsx                # Root layout + providers
│   └── globals.css               # Global styles
│
├── modules/                       # Domain-specific business logic
│   ├── auth/
│   │   ├── components/           # Auth-specific UI
│   │   ├── actions/              # Server Actions
│   │   ├── hooks/                # Auth hooks
│   │   ├── types.ts
│   │   └── validators.ts         # Zod schemas
│   ├── payments/
│   │   ├── components/
│   │   ├── actions/
│   │   ├── hooks/
│   │   ├── types.ts
│   │   └── validators.ts
│   └── dashboard/
│
├── shared/                        # Reusable across domains
│   ├── components/
│   │   ├── ui/                   # UI atoms (Button, Input, Card, etc.)
│   │   ├── layouts/              # Layout containers
│   │   └── feedback/             # Toast, Skeleton, EmptyState
│   ├── hooks/                    # Generic hooks (useAsync, useDebounce)
│   ├── utils/                    # Utilities (formatters, parsers)
│   └── types/                    # Shared domain types
│
├── core/                         # Infrastructure & configuration
│   ├── api/                      # HTTP client for backend
│   │   ├── client.ts            # Typed fetch wrapper
│   │   └── endpoints.ts         # API route definitions
│   │
│   ├── ory/                      # Ory integration
│   │   ├── client.ts            # Ory SDK instance
│   │   ├── middleware.ts        # Middleware utilities
│   │   ├── session.ts           # Session verification
│   │   └── hooks.ts             # Client-side Ory hooks
│   │
│   ├── validators/               # Global validation schemas
│   │   ├── email.ts
│   │   ├── currency.ts
│   │   └── index.ts
│   │
│   ├── config/
│   │   ├── constants.ts         # App constants
│   │   ├── env.ts              # Environment variables (typed)
│   │   └── csp.ts              # Content Security Policy
│   │
│   └── telemetry/
│       ├── sentry.ts           # Sentry SDK setup
│       └── otel.ts             # OpenTelemetry instrumentation
│
├── lib/                         # Legacy utility exports
│   └── [deprecated - migrate to modules/]
│
├── test/                        # Test utilities
│   ├── setup.ts                # Vitest setup
│   ├── mocks.ts                # MSW handlers
│   └── fixtures/               # Test data
│
├── ory.config.ts               # Ory SDK configuration
└── middleware.ts               # Next.js middleware (Ory session validation)

public/                          # Static assets (images, fonts)
docs/                           # Documentation
├── architecture.md             # Comprehensive architecture
└── adr/                        # Architectural Decision Records
```

### Module Structure (Per Domain)

Each domain module follows a consistent pattern:

```
src/modules/{domain}/
├── components/                 # UI components (atoms, molecules, organisms)
│   ├── {Feature}Card.tsx      # RSC or "use client" as needed
│   └── {Feature}Form.tsx
│
├── actions/                    # Server Actions (mutations)
│   └── index.ts               # "use server" - authenticated operations
│
├── hooks/                      # Client-side hooks
│   └── use{Feature}.ts        # "use client" - state management
│
├── types.ts                    # Domain-specific types
└── validators.ts              # Zod schemas for domain entities
```

---

## Core Patterns

### 1. Server Component (RSC) — Read-Only Pages

```typescript
// ✅ Default pattern for pages
// src/app/(dashboard)/payments/page.tsx

import { Suspense } from "react";
import { getOrySession } from "@/core/ory/session";
import { fetchTransactions } from "@/modules/payments/actions";
import { TransactionsList } from "./components/TransactionsList";

export default async function PaymentsPage({
  searchParams,
}: {
  searchParams: { page?: string };
}) {
  // Server-side: Verify session
  const session = await getOrySession();

  // Server-side: Fetch data
  const transactions = await fetchTransactions({
    userId: session.identity.id,
    page: parseInt(searchParams.page || "1"),
  });

  return (
    <div className="space-y-6">
      <h1>Payment History</h1>

      {/* Streaming: Progressive rendering with Suspense */}
      <Suspense fallback={<TransactionsSkeleton />}>
        <TransactionsList data={transactions} userId={session.identity.id} />
      </Suspense>
    </div>
  );
}

// Enable ISR (Incremental Static Regeneration)
export const revalidate = 60; // Revalidate every 60s
```

### 2. Client Component — Interactive Subtree

```typescript
// ✅ Client component only at leaf level
// src/modules/payments/components/FilterBar.tsx

"use client";

import { useState } from "react";

interface FilterBarProps {
  onFilter: (status: string) => Promise<void>;
}

export function FilterBar({ onFilter }: FilterBarProps) {
  const [status, setStatus] = useState("all");
  const [isLoading, setIsLoading] = useState(false);

  const handleFilter = async (newStatus: string) => {
    setIsLoading(true);
    setStatus(newStatus);
    await onFilter(newStatus);
    setIsLoading(false);
  };

  return (
    <div className="flex gap-2">
      <select
        value={status}
        onChange={(e) => handleFilter(e.target.value)}
        disabled={isLoading}
      >
        <option value="all">All</option>
        <option value="completed">Completed</option>
        <option value="pending">Pending</option>
        <option value="failed">Failed</option>
      </select>
    </div>
  );
}
```

### 3. Server Action — Protected Mutation

```typescript
// ✅ Server Actions for mutations
// src/modules/payments/actions/index.ts

"use server";

import { revalidatePath } from "next/cache";
import { requireOrySession } from "@/core/ory/session";
import { transferSchema } from "../validators";
import * as Sentry from "@sentry/nextjs";

export async function executeTransfer(formData: unknown) {
  try {
    // 1. Verify session (Zero-Trust)
    const session = await requireOrySession();

    // 2. Validate input
    const validated = transferSchema.safeParse(formData);
    if (!validated.success) {
      return { error: validated.error.flatten() };
    }

    // 3. Check authorization
    if (validated.data.fromUserId !== session.identity.id) {
      throw new Error("Unauthorized transfer");
    }

    // 4. Call backend API
    const response = await fetch(`${process.env.BACKEND_API_URL}/transfers`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${session.session_token}`,
      },
      body: JSON.stringify(validated.data),
    });

    if (!response.ok) {
      throw new Error("Transfer failed");
    }

    // 5. Revalidate cache
    revalidatePath("/dashboard/payments");

    return { success: true };
  } catch (error) {
    Sentry.captureException(error);
    return { error: "Transfer failed. Please try again." };
  }
}
```

### 4. Zod Validation — API Boundary Protection

```typescript
// ✅ Zod schema for type-safe validation
// src/modules/payments/validators.ts

import { z } from "zod";

export const transferSchema = z.object({
  fromUserId: z.string().uuid("Invalid user ID"),
  toUserId: z.string().uuid("Invalid recipient ID"),
  amount: z
    .number()
    .positive("Amount must be positive")
    .max(1_000_000, "Exceeds maximum transfer limit")
    .multipleOf(0.01, "Amount must have max 2 decimal places"),
  description: z
    .string()
    .max(500, "Description too long")
    .optional(),
  metadata: z.record(z.string()).optional(),
});

export type Transfer = z.infer<typeof transferSchema>;
```

### 5. Atomic Component — Reusable UI

```typescript
// ✅ Atom: Reusable across all modules
// src/shared/components/ui/Button.tsx

"use client";

interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: "primary" | "secondary" | "danger";
  size?: "sm" | "md" | "lg";
  isLoading?: boolean;
}

export function Button({
  variant = "primary",
  size = "md",
  isLoading = false,
  children,
  ...props
}: ButtonProps) {
  const baseStyles =
    "font-semibold rounded-lg transition-colors disabled:opacity-50";
  const variants = {
    primary: "bg-blue-600 text-white hover:bg-blue-700",
    secondary: "bg-gray-200 text-gray-900 hover:bg-gray-300",
    danger: "bg-red-600 text-white hover:bg-red-700",
  };
  const sizes = {
    sm: "px-2 py-1 text-xs",
    md: "px-4 py-2 text-sm",
    lg: "px-6 py-3 text-base",
  };

  return (
    <button
      {...props}
      disabled={props.disabled || isLoading}
      className={`${baseStyles} ${variants[variant]} ${sizes[size]}`}
    >
      {isLoading ? "Loading..." : children}
    </button>
  );
}
```

---

## Authentication (Ory)

### Login Flow (OIDC PKCE)

1. User clicks "Login" → Redirects to `/auth/login`.
2. Server renders Ory Elements (UserAuthForm) → Shows login UI.
3. User enters credentials → Ory Kratos validates.
4. Kratos returns session token + cookie.
5. Middleware verifies cookie on protected routes.

### Protected Routes

Protected routes use middleware to verify Ory session:

```typescript
// middleware.ts
const PUBLIC_ROUTES = ["/auth/login", "/auth/registration"];

export async function middleware(request: NextRequest) {
  if (PUBLIC_ROUTES.includes(request.nextUrl.pathname)) {
    return NextResponse.next();
  }

  // Verify session for protected routes
  const session = await verifyOrySession(request);
  if (!session) {
    return NextResponse.redirect(new URL("/auth/login", request.url));
  }

  return NextResponse.next();
}
```

### Access Session in Server Components

```typescript
import { getOrySession } from "@/core/ory/session";

export default async function ProtectedPage() {
  const session = await getOrySession();

  return (
    <div>
      Welcome, {session.identity.traits.email}!
      <p>User ID: {session.identity.id}</p>
    </div>
  );
}
```

### Using Ory in Client Components (Rare)

For client-side identity checks (e.g., avatar, logout button):

```typescript
"use client";

import { useOrySession } from "@ory/nextjs";

export function UserMenu() {
  const { data: session, isLoading } = useOrySession();

  if (isLoading) return <div>Loading...</div>;

  return (
    <div>
      <p>User: {session?.identity.traits.email}</p>
      <a href="/auth/logout">Logout</a>
    </div>
  );
}
```

---

## Data Fetching Strategy

### Hybrid Approach

| Scenario | Method | Caching |
|----------|--------|---------|
| Initial page load (read-only) | RSC + `fetch` | Next.js ISR |
| Form submission | Server Action | Manual revalidation |
| Complex filters | TanStack Query | In-memory |
| Real-time updates | WebSocket (future) | N/A |

### Server-Side Data Fetching (Default)

```typescript
// Use Next.js cache directives
async function fetchData() {
  const response = await fetch("https://api.example.com/data", {
    next: {
      revalidate: 60, // Revalidate every 60s
      tags: ["data"], // Can revalidate manually via revalidateTag()
    },
  });
  return response.json();
}
```

### Revalidation Strategies

```typescript
// On-demand revalidation after mutation
"use server";

import { revalidatePath, revalidateTag } from "next/cache";

export async function createTransaction() {
  // ... create transaction
  revalidatePath("/dashboard/payments"); // Revalidate route
  revalidateTag("transactions"); // Revalidate by tag
}
```

### Client-Side State (TanStack Query)

For complex client-side state (filtering, pagination with real-time):

```typescript
"use client";

import { useQuery } from "@tanstack/react-query";

export function useTransactions() {
  return useQuery({
    queryKey: ["transactions"],
    queryFn: async () => {
      const response = await fetch("/api/transactions");
      return response.json();
    },
    staleTime: 30000, // 30s
  });
}
```

---

## Security

### 1. Environment Variables

Never expose secrets in client-side code:

```typescript
// ✅ Safe: Server-only code
const backendKey = process.env.BACKEND_API_KEY; // Only in SSR

// ❌ Dangerous: Exposed to client
const clientSecret = process.env.NEXT_PUBLIC_SECRET; // Visible to browser
```

### 2. Input Validation (Zod)

All API boundaries protected by Zod:

```typescript
const schema = z.object({
  email: z.string().email(),
  amount: z.number().positive(),
});

const result = schema.safeParse(input);
if (!result.success) {
  return { error: result.error.flatten() };
}
```

### 3. Content Security Policy (CSP)

CSP headers prevent XSS attacks:

```typescript
// next.config.ts
export default {
  headers: async () => [
    {
      source: "/(.*)",
      headers: [
        {
          key: "Content-Security-Policy",
          value: "default-src 'self'; script-src 'self' 'unsafe-inline'",
        },
      ],
    },
  ],
};
```

### 4. HTTPS & Secure Cookies

```typescript
// Middleware sets secure flags on cookies
Set-Cookie: ory_session=xxx; Secure; HttpOnly; SameSite=Strict;
```

---

## Testing

### Unit Tests (Vitest)

```bash
bun test:unit
```

Test validators, utilities, pure functions:

```typescript
// src/modules/payments/validators.test.ts
import { describe, it, expect } from "vitest";
import { transferSchema } from "./validators";

describe("transferSchema", () => {
  it("validates valid transfer", () => {
    const result = transferSchema.safeParse({
      fromUserId: "550e8400-e29b-41d4-a716-446655440000",
      toUserId: "550e8400-e29b-41d4-a716-446655440001",
      amount: 100.50,
    });
    expect(result.success).toBe(true);
  });
});
```

### Integration Tests (Testing Library + MSW)

```bash
bun test:integration
```

Test components with mocked API responses.

### E2E Tests (Playwright)

```bash
bun test:e2e
```

Test complete user flows (login → transfer → confirmation).

---

## Development

### Commands

```bash
# Development
bun run dev                 # Start dev server

# Testing
bun run test              # Run all tests
bun run test:unit         # Run unit tests
bun run test:unit:watch   # Watch mode
bun run test:e2e          # Run E2E tests
bun run test:e2e:ui       # E2E UI mode

# Building
bun run build             # Build for production
bun run start             # Start production server

# Linting & Formatting
bun run lint              # Run ESLint
bun run format            # Format code with Prettier

# Observability
bun run analyze:bundle    # Analyze bundle size
```

### Code Style

- **Language:** TypeScript (strict mode)
- **Linting:** ESLint (Next.js rules)
- **Formatting:** Prettier
- **CSS:** Tailwind CSS utility-first

### Directory Navigation Tips

```
# View domain modules
ls src/modules/

# Find auth-related code
find src -name "*auth*"

# Find components for a feature
find src -path "*payments/components*"
```

---

## Deployment

### Build Process

```bash
bun run build
```

Outputs optimized Next.js build to `.next/`.

### Docker Deployment

```bash
# Build image
docker build -t laura-tech-frontend .

# Run container
docker run -p 3000:3000 laura-tech-frontend
```

### Environment Variables (Production)

Set in deployment platform:

```
ORY_SDK_URL=https://ory.lauratech.com
BACKEND_API_URL=https://api.lauratech.com
SENTRY_DSN=https://...@sentry.io/...
```

### CI/CD Pipeline

GitHub Actions workflow runs on every push:

1. **Lint:** ESLint checks
2. **Type Check:** TypeScript compilation
3. **Unit Tests:** Vitest
4. **Build:** Next.js build
5. **E2E Tests:** Playwright (critical flows only)

---

## Troubleshooting

### Session Not Persisting After Login

**Problem:** Login succeeds, but you're redirected to login on next request.

**Solution:** Ensure Ory cookies are being set:

```bash
# In browser DevTools → Application → Cookies
# Look for: ory_kratos_session, ory_kratos_session_csrf
```

If missing, check Ory Kratos is running:

```bash
curl http://localhost:4434/health/ready
```

### "use client" Hydration Mismatch

**Problem:** Error: "Hydration failed because initial UI does not match server-rendered HTML."

**Solution:** Ensure client components don't access server-only data (no async/promises):

```typescript
// ❌ Wrong: Async in client component
"use client";
export function Component() {
  const data = await fetch(...); // Can't use await in client component
}

// ✅ Correct: Parent RSC fetches, passes as prop
export async function Page() {
  const data = await fetch(...);
  return <ClientComponent data={data} />;
}
```

### Zod Validation Errors

**Problem:** "Expected string, received undefined."

**Solution:** Ensure all required fields are present in validation:

```typescript
const schema = z.object({
  email: z.string().email(), // Required
  phone: z.string().optional(), // Optional
});
```

---

## Contributing

1. Create feature branch: `git checkout -b feature/your-feature`
2. Follow code patterns (see [Core Patterns](#core-patterns))
3. Write tests for new features
4. Run `bun run lint` and `bun run test` before push
5. Submit PR with description of changes

---

## Resources

- **Documentation:** [docs/architecture.md](../docs/architecture.md)
- **ADRs:** [docs/adr/](../docs/adr/)
- **Next.js:** [nextjs.org](https://nextjs.org)
- **Ory:** [ory.sh](https://www.ory.sh)
- **Tailwind:** [tailwindcss.com](https://tailwindcss.com)
- **Zod:** [zod.dev](https://zod.dev)

---

**Made by:** Senior Principal Software Architect  
**For:** LauraTech MVP  
**Last Updated:** 2026-01-14
