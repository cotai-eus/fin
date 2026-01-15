# LauraTech Quick Reference Guide

## 🚀 Getting Started (5 minutes)

### 1. Install & Run
```bash
cd /home/user/fin/front

# Install dependencies (using Bun)
bun install

# Start development server
bun run dev

# Visit http://localhost:3000
```

### 2. Verify Setup
```bash
# Type checking
bun run typecheck

# Run tests
bun run test

# Check linting
bun run lint
```

### 3. Create `.env.local`
```env
# Ory Kratos
NEXT_PUBLIC_ORY_SDK_URL=http://localhost:4433

# Backend API
NEXT_PUBLIC_API_URL=http://localhost:8080

# Sentry (optional)
NEXT_PUBLIC_SENTRY_DSN=https://...
```

---

## 📖 Code Patterns

### Pattern 1: Creating a Server Component Page
```typescript
// src/app/payments/page.tsx
import { getOrySession } from "@/core/ory/session";
import { fetchData } from "@/modules/payments/actions";
import { Suspense } from "react";

export default async function Page() {
  // ✓ Verify session on server
  const session = await getOrySession();

  return (
    <div>
      <h1>Payments</h1>
      <Suspense fallback={<div>Loading...</div>}>
        <Content />
      </Suspense>
    </div>
  );
}

async function Content() {
  const data = await fetchData();
  return <div>{/* render data */}</div>;
}
```

### Pattern 2: Creating a Server Action
```typescript
// src/modules/payments/actions/index.ts
"use server";

import { getAuthenticatedUserId } from "@/core/ory/session";
import { transferSchema } from "../validators";
import * as Sentry from "@sentry/nextjs";

export async function executeTransfer(formData: FormData) {
  try {
    // 1. Verify session
    const userId = await getAuthenticatedUserId();

    // 2. Parse & validate input
    const data = transferSchema.safeParse({
      fromUserId: userId,
      toUserId: formData.get("toUserId"),
      amount: parseFloat(formData.get("amount") as string),
    });

    if (!data.success) {
      return { success: false, error: "Invalid input" };
    }

    // 3. Authorize (user can only transfer from own account)
    if (data.data.fromUserId !== userId) {
      return { success: false, error: "Unauthorized" };
    }

    // 4. Call backend API
    const response = await fetch(`${process.env.API_URL}/transfers`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${process.env.API_TOKEN}`,
      },
      body: JSON.stringify(data.data),
    });

    if (!response.ok) throw new Error(`API error: ${response.status}`);

    // 5. Revalidate cache
    revalidatePath("/payments");

    return { success: true, data: await response.json() };
  } catch (error) {
    // 6. Capture errors
    Sentry.captureException(error, {
      tags: { action: "executeTransfer" },
    });

    return { success: false, error: "Transfer failed" };
  }
}
```

### Pattern 3: Creating a Client Component
```typescript
// src/modules/payments/components/TransactionCard.tsx
"use client";

import { useState } from "react";
import { exportTransaction } from "../actions";

interface Props {
  transaction: Transaction;
  userId: string;
}

export function TransactionCard({ transaction, userId }: Props) {
  const [isLoading, setIsLoading] = useState(false);

  const handleExport = async () => {
    setIsLoading(true);
    try {
      const result = await exportTransaction({
        transactionId: transaction.id,
        format: "pdf",
      });

      if (result.success) {
        // Handle success
        window.open(result.downloadUrl, "_blank");
      } else {
        // Handle error
        console.error(result.error);
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div>
      <h3>{transaction.id}</h3>
      <button onClick={handleExport} disabled={isLoading}>
        {isLoading ? "Exporting..." : "Export"}
      </button>
    </div>
  );
}
```

### Pattern 4: Creating Validators
```typescript
// src/modules/payments/validators.ts
import { z } from "zod";

export const transferSchema = z.object({
  fromUserId: z.string().min(1, "From user required"),
  toUserId: z.string().min(1, "To user required"),
  amount: z
    .number()
    .positive("Amount must be positive")
    .multipleOf(0.01, "Must have 2 decimal places")
    .max(1_000_000, "Max 1M BRL per transfer")
    .describe("Transfer amount in BRL"),
  description: z.string().max(256, "Max 256 characters").optional(),
});

// Type extraction (no manual types needed!)
export type Transfer = z.infer<typeof transferSchema>;
```

### Pattern 5: Using Session in RSC
```typescript
// src/app/dashboard/page.tsx
import { requireOrySession, getUserId } from "@/core/ory/session";

export default async function DashboardPage() {
  // Get current user ID (throws if not authenticated)
  const userId = await getUserId();

  // Or get full session object
  const session = await requireOrySession();
  const email = session.identity?.traits?.email;

  return <div>Welcome, {email}!</div>;
}
```

---

## 🧪 Testing Examples

### Unit Test (Vitest)
```typescript
// src/modules/payments/validators.test.ts
import { describe, it, expect } from "vitest";
import { transferSchema } from "./validators";

describe("transferSchema", () => {
  it("validates valid transfer", () => {
    const result = transferSchema.safeParse({
      fromUserId: "user1",
      toUserId: "user2",
      amount: 1000.50,
    });
    expect(result.success).toBe(true);
  });

  it("rejects invalid amount", () => {
    const result = transferSchema.safeParse({
      fromUserId: "user1",
      toUserId: "user2",
      amount: 1000.999, // Invalid decimals
    });
    expect(result.success).toBe(false);
  });
});
```

### Integration Test (Testing Library + MSW)
```typescript
// src/modules/payments/components/TransactionCard.test.tsx
import { render, screen } from "@testing-library/react";
import { TransactionCard } from "./TransactionCard";

describe("TransactionCard", () => {
  it("renders transaction details", () => {
    const tx = {
      id: "tx-123",
      amount: 100.00,
      status: "completed",
      fromUserId: "user1",
      toUserId: "user2",
      createdAt: new Date(),
    };

    render(<TransactionCard transaction={tx} userId="user1" />);

    expect(screen.getByText("tx-123")).toBeInTheDocument();
    expect(screen.getByText("R$ 100,00")).toBeInTheDocument();
  });
});
```

### E2E Test (Playwright)
```typescript
// e2e/payments.spec.ts
import { test, expect } from "@playwright/test";

test("complete payment flow", async ({ page }) => {
  // Navigate to payments
  await page.goto("http://localhost:3000/payments");

  // Wait for page to load
  await expect(page.locator("h1")).toContainText("Payments");

  // Filter transactions
  await page.click("button:has-text('Completed')");

  // View transaction details
  await page.click("[data-testid='transaction-row']:first-child");

  // Export transaction
  await page.click("button:has-text('Export PDF')");

  // Verify download
  const downloadPath = await page.waitForEvent("popup");
  expect(downloadPath.url()).toContain("download");
});
```

---

## 📁 File Organization Reference

### Create a new feature module
```
src/modules/{feature}/
├── components/
│   ├── MyComponent.tsx        # Client component with "use client"
│   └── MyContainer.tsx        # RSC component
├── actions/
│   └── index.ts               # Server Actions (mutations)
├── hooks/
│   └── useMyHook.ts           # Custom hooks
├── types.ts                   # TypeScript types
└── validators.ts              # Zod schemas
```

### Component hierarchy
```
RSC Page (src/app/...)
  └─ RSC Container (src/modules/{domain}/components/)
      └─ Client Organism (src/modules/{domain}/components/)
          ├─ Client Molecule (src/shared/components/)
          │   └─ Atoms (src/shared/components/ui/)
          └─ Client Molecule
              └─ Atoms
```

---

## 🔒 Security Checklist for New Features

- [ ] All mutations use Server Actions
- [ ] Server Actions verify session with `requireOrySession()`
- [ ] All inputs validated with Zod schemas
- [ ] Authorization checks (user ID matching)
- [ ] Errors captured to Sentry
- [ ] Cache revalidated after mutations
- [ ] No secrets in client-side code
- [ ] HTTP-only cookies for sessions
- [ ] Rate limiting on sensitive endpoints (backend)

---

## 🐛 Common Issues & Solutions

### Issue: Hydration Mismatch
```
Error: Text content did not match. Server: "..." Client: "..."
```

**Solution:** Use `useId()` instead of `Math.random()` or generate IDs on the server
```typescript
"use client";
import { useId } from "react";

export function Component() {
  const id = useId(); // ✓ Stable across render
  return <div id={id}>{id}</div>;
}
```

### Issue: Can't import server module in client component
```
Error: Module X is not supported in the browser
```

**Solution:** Move to a separate server file and call it via Server Action
```typescript
// ✓ Correct
"use server";
export async function getData() {
  // import heavy module here
}

// Call from client component
const result = await getData();
```

### Issue: Zod validation error
```
Error: Expected number, received string
```

**Solution:** Parse to correct type before validating
```typescript
const schema = z.object({
  amount: z.number(),
});

// ✓ Correct
schema.parse({ amount: parseFloat(formData.get("amount")) });

// ✗ Wrong
schema.parse({ amount: formData.get("amount") }); // Still a string!
```

### Issue: Session not persisting
```
getOrySession() returns null
```

**Solution:** Check cookie settings and Ory configuration
```bash
# Verify cookie exists
devtools → Application → Cookies

# Check Ory logs
docker logs kratos

# Verify COOKIE_DOMAIN matches your domain
```

---

## 📊 Development Commands Reference

```bash
# Development
bun run dev              # Start dev server on port 3000
bun run dev --open       # Auto-open browser

# Testing
bun run test             # Run all tests once
bun run test:watch       # Watch mode
bun run test:coverage    # Coverage report

# Quality
bun run typecheck        # TypeScript type checking
bun run lint             # ESLint
bun run format           # Prettier formatting

# Building
bun run build            # Production build
bun run start            # Start production server
bun run analyze          # Analyze bundle size

# E2E Testing
bun run test:e2e         # Run Playwright tests
bun run test:e2e --ui    # Interactive UI mode
```

---

## 🎯 Architecture Decision Quick Reference

| Decision | Rationale | Location |
|----------|-----------|----------|
| **Next.js 16 + RSC** | 40-60% bundle reduction, security | [ADR-001](../docs/adr/ADR-001-nextjs-rsc-adoption.md) |
| **Ory Kratos OIDC** | Industry-standard auth, self-hosted | [ADR-002](../docs/adr/ADR-002-ory-zero-trust.md) |
| **Atomic Design** | Scalable component hierarchy | [ADR-003](../docs/adr/ADR-003-atomic-server-components.md) |
| **Tri-layer Testing** | Balance speed/confidence | [ADR-004](../docs/adr/ADR-004-tri-layer-testing.md) |
| **Zod Validation** | Runtime safety + type inference | validators.ts |
| **Server Actions** | Type-safe mutations, auth-aware | actions/index.ts |

---

## 📚 Documentation Links

- **Architecture:** [docs/architecture.md](../docs/architecture.md)
- **Development Guide:** [front/README.md](../front/README.md)
- **ADR Template:** [docs/adr/TEMPLATE.md](../docs/adr/TEMPLATE.md)
- **Example Code:** [src/app/(dashboard)/payments/](../src/app/\(dashboard\)/payments/)
- **Implementation Summary:** [IMPLEMENTATION_COMPLETE.md](../IMPLEMENTATION_COMPLETE.md)

---

## 💡 Pro Tips

1. **Use `revalidatePath()`** after mutations to refresh data:
   ```typescript
   revalidatePath("/dashboard"); // Revalidate page
   revalidateTag("transactions"); // Revalidate by tag
   ```

2. **Type your Server Actions** for IDE autocomplete:
   ```typescript
   async function action(formData: FormData): Promise<Result> {}
   ```

3. **Always validate responses** from external APIs:
   ```typescript
   const result = transactionSchema.parse(await api.call());
   ```

4. **Use error boundaries** for better UX:
   ```typescript
   <ErrorBoundary>
     <ComponentThatMayFail />
   </ErrorBoundary>
   ```

5. **Leverage Suspense** for streaming:
   ```typescript
   <Suspense fallback={<Loading />}>
     <SlowComponent />
   </Suspense>
   ```

---

**Last Updated:** 2024  
**For Questions:** See the full documentation in [IMPLEMENTATION_COMPLETE.md](../IMPLEMENTATION_COMPLETE.md)
