# ADR-002: Zero-Trust Identity Architecture with Ory

**Status:** Accepted  
**Context:** Authentication & Security  
**Date:** 2026-01-14  
**Ratification:** Senior Principal Software Architect  

---

## Problem Statement

Fintech applications handle sensitive user data and financial transactions. A robust identity system must:

1. **Prevent unauthorized access** to protected routes and operations.
2. **Enforce session validation** on every request (Zero-Trust principle).
3. **Support multiple auth methods** (email/password, MFA, social) without fragmenting codebase.
4. **Isolate auth logic** from business logic to enable independent scaling and security audits.
5. **Maintain PCI/LGPD compliance** by never exposing credentials or sensitive identity data to the frontend.

**Question:** How do we implement a Zero-Trust identity architecture that leverages Ory Kratos (already configured) while integrating seamlessly with Next.js 16.1 Server Components?

---

## Context

### Requirements

- Ory Kratos is already containerized and available at development/staging.
- Must use OIDC PKCE flow (OAuth 2.0 with Proof Key for Public Clients).
- Frontend must not directly access user password or secrets.
- All protected routes must verify session integrity.
- Session data must be available in Next.js Server Components for secure operations.
- Compliance: LGPD (Brazil's data protection law), PCI-DSS equivalent.

### Current State

- Ory Elements React library is installed (`@ory/elements-react` v1.1.0).
- `@ory/nextjs` (v1.0.0-rc.0) is available but not fully integrated.
- Middleware is scaffolded but minimal.
- Auth routes exist (login, registration, recovery, verification) but have no implementation.

### Threat Model

```
Threat                              | Risk Level | Mitigation
------------------------------------+------------+----------------------------------
Session Hijacking (stolen cookie)   | Critical   | HTTPS only, secure flags, short TTL
Token Exposure (in URL/query)       | Critical   | PKCE flow, no tokens in logs
CSRF Attacks                        | High       | SameSite=Strict, CSRF tokens
Replay Attacks                      | High       | Nonce validation, timestamp checks
Privilege Escalation                | High       | Verify permissions on every action
Data Exposure (logs)                | Medium     | Sanitize PII from logs
```

---

## Decision

**Implement a Zero-Trust identity architecture with Ory Kratos as the exclusive identity provider, enforcing session validation at two layers:**

### Layer 1: Middleware (Request-Level Verification)

All incoming requests to protected routes must pass middleware that:
- Extracts session cookie from request.
- Validates Ory session token against Kratos API.
- Injects user identity context into request headers.
- Redirects unauthenticated users to Ory login UI.

### Layer 2: Server Component (Operation-Level Verification)

All sensitive Server Components and Server Actions must:
- Call `getSession()` utility to retrieve current session.
- Verify session is valid and user has required permissions.
- Throw errors if session is missing or invalid.
- Never trust implicit session (re-verify per operation).

### Architecture: Ory OIDC PKCE Flow

```
┌──────────────┐         ┌──────────────┐         ┌──────────────┐
│  Frontend    │         │  Next.js     │         │  Ory Kratos  │
│  (Browser)   │         │  Server      │         │  (Auth)      │
└──────────────┘         └──────────────┘         └──────────────┘
       │                       │                        │
       ├──────────────────────────────────────────────>│
       │ 1. User clicks "Login"                       │
       │ Redirect to /auth/login                      │
       │                                               │
       │ <────────────────────────────────────────────┤
       │ Redirect to Ory UI: https://ory/auth?...     │
       │                                               │
       ├──── 2. User enters credentials ─────────────>│
       │                                               │
       │ <───── 3. Auth success, return code+state ───┤
       │                                               │
       ├──────────────────────────────────────────────>│
       │ 4. Redirect to callback: /auth/callback?code │
       │                                               │
       │                    ┌─────────────────────────>│
       │                    │ 5. Exchange code for    │
       │                    │ session (backend only)   │
       │                    │ <────────────────────────┤
       │                    │ Return session JWT       │
       │                    │                          │
       │ <───────────────────────────────────────────>│
       │ 6. Set HTTP-only cookie                      │
       │ Redirect to /dashboard                       │
       │                                               │
       ├──────────────────────────────────────────────>│
       │ 7. Request /dashboard with cookie             │
       │                                               │
       │             ┌──────────────────────────────>│
       │             │ 8. Middleware verifies        │
       │             │ session with Kratos           │
       │             │ <──────────────────────────────┤
       │             │ OK, return session data        │
       │             │                                │
       │ <──────────────────────────────────────────>│
       │ 9. Render protected page + user context      │
       │                                               │
       └                 (session valid)               ┘
```

---

## Justification

### 1. Zero-Trust Principle

> "Never trust, always verify."

Every request and operation is independently verified, regardless of prior authentication. This is the industry standard for fintech applications (Stripe, Square, Wise all employ this model).

**Benefits:**
- Prevents lateral movement if one session is compromised.
- Each Server Action re-verifies user identity before executing.
- Permission checks happen at operation level, not route level.

### 2. Ory Kratos Specialization

Ory Kratos is purpose-built for identity management:
- ✅ Credential management (password hashing, validation).
- ✅ Session lifecycle (creation, refresh, revocation).
- ✅ Multi-factor authentication (built-in).
- ✅ Account recovery flows (email recovery, password reset).
- ✅ OIDC/OAuth2 compliance.

Using Ory vs. building in-house eliminates a massive security surface.

### 3. PKCE Flow Over Password Grant

**PKCE (Proof Key for Public Clients):**
- Browser never sends password to Next.js server.
- Password only sent to Ory Kratos (over HTTPS).
- Kratos issues code + state, frontend exchanges for token.
- **Result:** Frontend server never handles plaintext passwords.

**Comparison:**

| Approach | Security | Complexity | Compliance |
|----------|----------|-----------|-----------|
| Password Grant | ❌ (server stores password) | Low | ❌ Risky |
| PKCE | ✅ (no password on frontend) | Medium | ✅ PCI-DSS, LGPD friendly |

### 4. Middleware + Server Component Double Validation

Validating at two layers provides defense-in-depth:

```
Layer 1 (Middleware)             Layer 2 (Server Component)
┌──────────────────────┐        ┌──────────────────────────┐
│ Quick session check   │        │ Deep permission check    │
│ Reject invalid tokens │        │ Verify user_id matches  │
│ Rate limit attempts   │        │ Enforce row-level access │
│ Log suspicious access │        │ Validate action scope    │
└──────────────────────┘        └──────────────────────────┘
```

**Example:**
```typescript
// Middleware catches invalid tokens
// Server Component catches unauthorized operations
export async function executeTransfer(targetUserId: string) {
  const session = await getSession();
  
  // Middleware verified session exists
  // But we still verify the operation is authorized
  if (session.identity.id !== targetUserId) {
    throw new Error("Cannot transfer from another user's account");
  }
  // Proceed with transfer
}
```

### 5. HTTP-Only Cookies + Secure Flags

Sessions are stored in **HTTP-only cookies** which:
- Cannot be accessed by JavaScript (XSS protection).
- Cannot be sent to different domain (CSRF protection).
- Automatically included by browser in same-site requests.

```http
Set-Cookie: ory_session_id=xxxx; 
  HttpOnly;              ← JavaScript cannot access
  Secure;                ← HTTPS only
  SameSite=Strict;       ← No cross-site requests
  Path=/;
  Max-Age=3600           ← 1 hour session
```

### 6. Alignment with Next.js 16.1 Server Components

RSCs naturally integrate with Ory:

```typescript
// ✅ Server Component can call getSession() server-side
export default async function TransferPage() {
  const session = await getSession();  // Calls Ory
  const user = await db.users.findUnique({ id: session.identity.id });
  return <TransferForm user={user} />;
}

// ✅ Server Actions are protected by design
"use server";
export async function executeTransfer(data) {
  const session = await getSession();
  // Ory verifies session; if invalid, throws error
  // Client never knows why (security through obscurity)
  await db.transfers.create({ from_user_id: session.identity.id, ...data });
}
```

---

## Alternatives Considered

### Alternative 1: In-House Session Management

| Aspect | Evaluation |
|--------|------------|
| Control | ✅ Full control over implementation |
| Maintenance | ❌ Must handle password hashing, validation, rotation |
| Security | ❌ High risk of implementation errors |
| Compliance | ❌ Difficult to audit and certify |
| Time | ❌ 2-3 weeks to build correctly |
| **Decision** | ❌ **Rejected** — Too risky for fintech |

### Alternative 2: OAuth 2.0 Resource Owner Password Grant

| Aspect | Evaluation |
|--------|------------|
| Security | ❌ Password stored on frontend server (PCI risk) |
| UX | ✅ Simple flow |
| Compliance | ❌ Non-compliant with OAuth best practices |
| **Decision** | ❌ **Rejected** — Security risk outweighs simplicity |

### Alternative 3: JWT-Only (No Cookies)

| Aspect | Evaluation |
|--------|------------|
| Simplicity | ✅ Easier client-side state management |
| Security | ❌ JWTs in localStorage vulnerable to XSS |
| CSRF | ❌ Custom CSRF tokens required |
| Refresh | ❌ Refresh token rotation more complex |
| **Decision** | ❌ **Rejected** — HTTP-only cookies are best practice |

### Alternative 4: Ory Hydra (Full OAuth2/OIDC Provider)

| Aspect | Evaluation |
|--------|------------|
| Use Case | ✅ Useful if LauraTech becomes identity provider for partners |
| Complexity | ❌ Over-engineered for MVP (Kratos alone sufficient) |
| Overhead | ❌ Extra infrastructure, more moving parts |
| **Decision** | ⚠️ **Deferred** — Consider for Phase 2 if needed |

---

## Implementation Details

### 1. Middleware: Session Validation

```typescript
// middleware.ts
import { NextResponse, type NextRequest } from "next/server";
import { verifyOrySession } from "@/core/ory/session";

const PUBLIC_ROUTES = [
  "/auth/login",
  "/auth/registration",
  "/auth/recovery",
  "/auth/verification",
];

const PROTECTED_ROUTES = ["/(dashboard)/:path*"];

export async function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;

  // Allow public auth routes
  if (PUBLIC_ROUTES.some((route) => pathname.startsWith(route))) {
    return NextResponse.next();
  }

  // Verify session for protected routes
  try {
    const session = await verifyOrySession(request.cookies);

    if (!session?.active) {
      return NextResponse.redirect(new URL("/auth/login", request.url));
    }

    // Inject user context into response headers
    // Available to Server Components via headers()
    const response = NextResponse.next();
    response.headers.set("x-user-id", session.identity.id);
    response.headers.set("x-user-email", session.identity.traits.email);

    return response;
  } catch (error) {
    // Session validation failed (invalid token, Kratos unreachable, etc.)
    return NextResponse.redirect(new URL("/auth/login", request.url));
  }
}

export const config = {
  matcher: [
    "/((?!_next/static|_next/image|favicon.ico|public/).*)",
  ],
};
```

### 2. Session Utility: Server-Side Verification

```typescript
// src/core/ory/session.ts
import { FrontendApi } from "@ory/client";
import { cookies, headers } from "next/headers";

const oryClient = new FrontendApi({
  basePath: process.env.ORY_SDK_URL,
  baseOptions: {
    withCredentials: true,
  },
});

export async function getOrySession() {
  try {
    // Retrieve Ory session cookie
    const cookieStore = await cookies();
    const sessionCookie = cookieStore.get("ory_kratos_session");

    if (!sessionCookie) {
      throw new Error("No session cookie found");
    }

    // Verify with Ory Kratos
    const { data: session } = await oryClient.toSession(undefined, {
      headers: {
        Cookie: `ory_kratos_session=${sessionCookie.value}`,
      },
    });

    if (!session?.active) {
      throw new Error("Session inactive");
    }

    return session;
  } catch (error) {
    console.error("Session verification failed:", error);
    throw new Error("Unauthorized");
  }
}

// Utility for Server Actions
export async function requireOrySession() {
  const session = await getOrySession();
  if (!session) throw new Error("Unauthorized");
  return session;
}
```

### 3. Server Action: Protected Mutation

```typescript
// src/modules/payments/actions.ts
"use server";

import { revalidatePath } from "next/cache";
import { requireOrySession } from "@/core/ory/session";
import { transferSchema } from "./validators";
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
    // (Verify user isn't transferring from another account)
    if (validated.data.fromUserId !== session.identity.id) {
      Sentry.captureMessage(
        `Unauthorized transfer attempt by user ${session.identity.id}`,
        "warning"
      );
      throw new Error("Unauthorized transfer");
    }

    // 4. Call backend API
    // (Pass Ory token for backend to re-verify)
    const response = await fetch(
      `${process.env.BACKEND_API_URL}/transfers`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Authorization": `Bearer ${session.session_token}`,
          "X-Request-ID": crypto.randomUUID(),
          "X-User-ID": session.identity.id,
        },
        body: JSON.stringify(validated.data),
      }
    );

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message);
    }

    const result = await response.json();

    // 5. Revalidate cache
    revalidatePath("/dashboard/payments");

    return { success: true, data: result };
  } catch (error) {
    Sentry.captureException(error, {
      tags: { module: "payments", action: "executeTransfer" },
    });

    // Return generic error to prevent information leakage
    return { error: "Transfer failed. Please try again." };
  }
}
```

### 4. Protected Page: Server Component

```typescript
// src/app/(dashboard)/payments/page.tsx
import { Suspense } from "react";
import { getOrySession } from "@/core/ory/session";
import { fetchUserTransactions } from "@/modules/payments/actions";
import TransactionsList from "./components/TransactionsList";
import { PageHeader } from "@/shared/components/PageHeader";

export default async function PaymentsPage({
  searchParams,
}: {
  searchParams: { page?: string };
}) {
  // Verify session (middleware verified, but we re-verify for defense-in-depth)
  const session = await getOrySession();

  // Fetch user transactions (authorization happens server-side)
  const transactions = await fetchUserTransactions({
    userId: session.identity.id,
    page: parseInt(searchParams.page || "1"),
  });

  return (
    <div className="space-y-6">
      <PageHeader
        title="Payment History"
        subtitle={`Logged in as ${session.identity.traits.email}`}
      />

      <Suspense fallback={<TransactionsSkeleton />}>
        <TransactionsList
          userId={session.identity.id}
          initialData={transactions}
        />
      </Suspense>
    </div>
  );
}

export const revalidate = 60; // ISR: revalidate every 60s
```

### 5. Login Flow: Ory Elements Integration

```typescript
// src/app/(auth)/login/page.tsx
"use client";

import { useEffect } from "react";
import { UserAuthForm } from "@ory/elements-react";
import { useOryFlowConfig } from "@/core/ory/hooks";

export default function LoginPage() {
  const { flow, isLoading, error } = useOryFlowConfig("login");

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;

  return (
    <div className="flex justify-center items-center min-h-screen">
      {flow && <UserAuthForm flow={flow} />}
    </div>
  );
}
```

---

## Consequences

### ✅ Positive Consequences

| Benefit | Impact |
|---------|--------|
| **Security** | Zero-Trust principle prevents unauthorized operations |
| **Compliance** | PCI-DSS, LGPD-compliant identity handling |
| **Audit Trail** | Every operation is verified and can be logged |
| **Flexibility** | Easy to add MFA, social login, account recovery |
| **Maintainability** | Ory handles identity; backend/frontend focus on business logic |

### ⚠️ Negative Consequences & Mitigations

| Challenge | Mitigation |
|-----------|-----------|
| **Ory Dependency** | If Kratos is down, all auth fails. Mitigation: High availability setup, fallback cache for valid sessions |
| **API Call Overhead** | Every request validates session with Kratos. Mitigation: Cache session in Redis, short TTL |
| **Debugging Complexity** | Multiple services (frontend, Kratos, backend) involved. Mitigation: Structured logging, request IDs |

---

## Metrics & Success Criteria

- [ ] All protected routes require valid Ory session.
- [ ] Invalid sessions redirect to login within < 100ms.
- [ ] No passwords logged or exposed in error messages.
- [ ] Session revocation takes effect within 5 minutes.
- [ ] MFA support implemented (Phase 2).
- [ ] Zero unauthorized access incidents.

---

## Related Decisions

- **ADR-001:** Server Components naturally integrate with server-side session verification.
- **ADR-003:** Componentization respects authentication boundaries.
- **ADR-004:** Testing strategy accounts for server-side session verification.

---

## References

- [Ory Kratos Documentation](https://www.ory.sh/kratos/)
- [OAuth 2.0 PKCE (RFC 7636)](https://tools.ietf.org/html/rfc7636)
- [OWASP Session Management Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html)
- [Zero Trust Architecture (NIST SP 800-207)](https://nvlpubs.nist.gov/nistpubs/SpecialPublications/NIST.SP.800-207.pdf)

---

**Author:** Senior Principal Software Architect  
**Created:** 2026-01-14  
**Approved:** Pending  
