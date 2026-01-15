# ADR-004: Tri-Layer Testing Strategy (Vitest, Testing Library, Playwright)

**Status:** Accepted  
**Context:** Quality Assurance & Testing  
**Date:** 2026-01-14  
**Ratification:** Senior Principal Software Architect  

---

## Problem Statement

Fintech applications require rigorous testing to prevent bugs in critical payment flows. However, testing in Next.js 16.1 with RSCs, Server Actions, and mixed Client/Server components introduces complexity:

1. **Unit Testing:** How to test server-side async functions, Zod validators, and utilities?
2. **Component Testing:** How to test RSCs that fetch data server-side?
3. **Integration Testing:** How to test API contracts without spinning up entire backend?
4. **E2E Testing:** How to test complete user flows including payment processing, auth, error recovery?

**Question:** What is the optimal testing pyramid for LauraTech that balances **coverage**, **speed**, **confidence**, and **maintainability**?

---

## Context

### Current State

- No test infrastructure configured (no Jest/Vitest, no Testing Library, no Playwright).
- MVP timeline is tight (4-7 weeks); expensive, slow tests are problematic.
- Critical payment flows must have high confidence before production release.
- Team has React testing experience but limited Server Component testing experience.

### Test Pyramid Principle

```
        /\
       /  \
      / E2E \        5-10% of tests
     /      \        (Slow, expensive, high confidence)
    /────────\
   /          \
  /Integration  \    20-30% of tests
 /              \    (Medium speed, medium maintenance)
/────────────────\
/                  \
/    Unit Tests      \  60-70% of tests
/                    \ (Fast, cheap, low maintenance)
──────────────────────
```

**For LauraTech:** Inverted pyramid emphasizes E2E for payment flows (high risk) + unit tests for logic.

### Regulatory Requirements

- PCI-DSS: Payment handling must be thoroughly tested.
- LGPD: Data handling flows must be audited.
- **Implication:** E2E tests for payment flows are mandatory.

---

## Decision

**Implement a tri-layer testing strategy:**

### Layer 1: Unit Tests (Vitest)

**What to test:**
- Utility functions (formatters, validators, parsers).
- Zod schemas.
- Server Actions (business logic, no DB calls).
- Domain hooks (not RSCs).

**Tools:**
- **Vitest:** Fast unit test runner (built on Vite).
- **No mocking required** for pure functions.

**Speed:** ~50-100ms per test.

### Layer 2: Integration Tests (Vitest + Testing Library + MSW)

**What to test:**
- Components rendering with various data states.
- Client components with user interactions.
- Server Actions that interact with mocked API responses.
- Form validation flows.

**Tools:**
- **Vitest:** Test runner.
- **React Testing Library:** Component testing (user-centric, not implementation-centric).
- **MSW (Mock Service Worker):** Mock API responses.

**Speed:** ~200-500ms per test suite.

### Layer 3: E2E Tests (Playwright)

**What to test:**
- Complete user journeys (login → transfer → confirmation).
- Payment flow (happy path + error scenarios).
- Session management (login → logout → re-login).
- Multi-step forms.
- Critical business flows.

**Tools:**
- **Playwright:** Browser automation.
- **Real backend** (staging environment) or **containerized services** (local).

**Speed:** ~2-5 seconds per test.

---

## Justification

### 1. Vitest for Unit Tests

**Why Vitest over Jest?**

| Aspect | Jest | Vitest |
|--------|------|--------|
| Speed | 🟡 Slower | ✅ 10-100x faster (uses Vite) |
| Config | 🟡 Babel + Jest config | ✅ Uses Vite config (already in Next.js) |
| ESM Support | 🟡 Limited | ✅ Native |
| Watch Mode | 🟡 Works | ✅ Faster feedback loop |
| Next.js Integration | ✅ Better known | ⚠️ Growing support (Next.js 16 supports) |

**Conclusion:** Vitest matches Next.js dev velocity. Faster feedback = faster iteration = faster MVP delivery.

### 2. Testing Library for Component Tests

**Why Testing Library over Enzyme?**

Testing Library encourages **user-centric testing** (how users interact) vs. implementation details:

```typescript
// ❌ Implementation-centric (Enzyme)
expect(component.find(Button).props().onClick).toBeDefined();

// ✅ User-centric (Testing Library)
const button = screen.getByRole("button", { name: /submit/i });
await userEvent.click(button);
expect(screen.getByText("Success")).toBeInTheDocument();
```

This prevents brittle tests that break when refactoring.

### 3. MSW for API Mocking

**Why MSW over Sinon/nock?**

```typescript
// MSW: Intercepts requests at network level
import { http, HttpResponse } from "msw";
import { setupServer } from "msw/node";

const server = setupServer(
  http.post("/api/transfers", () => {
    return HttpResponse.json({ success: true });
  })
);

beforeAll(() => server.listen());
afterEach(() => server.resetHandlers());
afterAll(() => server.close());

// Test automatically uses mocked API
test("transfer form submits correctly", async () => {
  render(<TransferForm />);
  // User fills form and submits
  // API call is intercepted by MSW, returns mocked response
});
```

**Benefits:**
- Works for both fetch and axios.
- No changes to application code needed.
- Can switch to real API in E2E tests.

### 4. Playwright for E2E Tests

**Why Playwright over Cypress?**

| Aspect | Cypress | Playwright |
|--------|---------|-----------|
| Speed | 🟡 Slower | ✅ Faster |
| Multi-browser | ⚠️ Limited | ✅ Chrome, Firefox, Safari |
| API | 🟡 Command queue | ✅ Modern async/await |
| Debugging | ✅ Built-in | ✅ Strong tooling |
| Payment Testing | ⚠️ Hard to mock payments | ✅ Can handle Stripe, etc. |

**Conclusion:** Playwright's multi-browser support + modern API align with LauraTech's requirements.

### 5. Tri-Layer Pyramid for Fintech

**For payments, high E2E coverage is justified:**

- **Unit tests:** Fast feedback on business logic.
- **Integration tests:** Verify components work together.
- **E2E tests:** Verify end-to-end payment flows (high risk = high investment).

```
           ┌────────┐
           │  E2E   │  (Payment flows, auth flows)
           │ Playwright
           └────────┘
         ┌──────────────┐
         │ Integration  │  (Components, forms)
         │ Testing Lib  │
         └──────────────┘
    ┌────────────────────────┐
    │      Unit Tests        │  (Logic, validators)
    │       Vitest           │
    └────────────────────────┘
```

---

## Alternatives Considered

### Alternative 1: Jest + Enzyme + Cypress

| Aspect | Evaluation |
|--------|------------|
| Familiarity | ✅ Most teams know Jest |
| Speed | ❌ Jest slower than Vitest |
| Integration | ❌ Enzyme couples to implementation |
| E2E | ⚠️ Cypress has browser limitations |
| **Decision** | ❌ **Rejected** — Slower feedback loop |

### Alternative 2: No Unit Tests (E2E Only)

| Aspect | Evaluation |
|--------|------------|
| Speed | ✅ No slow tests |
| Coverage | ❌ E2E doesn't cover edge cases (off-by-one, formatter bugs, etc.) |
| Stability | ❌ Flaky E2E tests (slow, environment-dependent) |
| Cost | ❌ Longer feedback loop = slower development |
| **Decision** | ❌ **Rejected** — E2E alone insufficient |

### Alternative 3: All Tests in E2E

| Aspect | Evaluation |
|--------|------------|
| Confidence | ✅ Tests real scenarios |
| Speed | ❌ Very slow (2-5s per test × 1000s tests = hours) |
| Stability | ❌ Flaky (environment dependent) |
| CI/CD | ❌ Expensive (cloud-based Playwright infrastructure) |
| **Decision** | ❌ **Rejected** — Breaks development velocity |

---

## Implementation Details

### 1. Vitest Setup

```typescript
// vitest.config.ts
import react from "@vitejs/plugin-react";
import path from "path";
import { defineConfig } from "vitest/config";

export default defineConfig({
  plugins: [react()],
  test: {
    environment: "jsdom",
    globals: true,
    setupFiles: ["./src/test/setup.ts"],
    coverage: {
      provider: "v8",
      reporter: ["text", "json", "html"],
      exclude: ["node_modules/", "src/test/"],
    },
  },
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
});
```

### 2. Unit Test Example (Zod Validator)

```typescript
// src/modules/payments/validators.test.ts
import { describe, it, expect } from "vitest";
import { transferSchema } from "./validators";

describe("transferSchema", () => {
  it("accepts valid transfer", () => {
    const data = {
      recipientId: "550e8400-e29b-41d4-a716-446655440000",
      amount: 100.50,
      description: "Payment for invoice",
    };

    const result = transferSchema.safeParse(data);
    expect(result.success).toBe(true);
    expect(result.data).toEqual(data);
  });

  it("rejects negative amount", () => {
    const data = {
      recipientId: "550e8400-e29b-41d4-a716-446655440000",
      amount: -50,
    };

    const result = transferSchema.safeParse(data);
    expect(result.success).toBe(false);
    expect(result.error?.flatten().fieldErrors.amount).toBeDefined();
  });

  it("rejects invalid UUID", () => {
    const data = {
      recipientId: "not-a-uuid",
      amount: 100,
    };

    const result = transferSchema.safeParse(data);
    expect(result.success).toBe(false);
  });

  it("rejects amount exceeding limit", () => {
    const data = {
      recipientId: "550e8400-e29b-41d4-a716-446655440000",
      amount: 2_000_000, // Exceeds max
    };

    const result = transferSchema.safeParse(data);
    expect(result.success).toBe(false);
  });
});
```

### 3. Integration Test Example (Component)

```typescript
// src/modules/payments/components/TransactionCard.test.tsx
import { render, screen } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import { TransactionCard } from "./TransactionCard";

describe("TransactionCard", () => {
  const mockTransaction = {
    id: "tx_123",
    amount: 150.75,
    recipientName: "Jane Doe",
    status: "completed",
    createdAt: new Date("2026-01-14"),
  };

  it("displays transaction details", () => {
    render(<TransactionCard transaction={mockTransaction} />);

    expect(screen.getByText("Jane Doe")).toBeInTheDocument();
    expect(screen.getByText(/R\$ 150,75/)).toBeInTheDocument();
  });

  it("shows correct status badge", () => {
    render(<TransactionCard transaction={mockTransaction} />);

    expect(screen.getByText("completed")).toHaveClass("bg-green");
  });

  it("formats date in Brazilian locale", () => {
    render(<TransactionCard transaction={mockTransaction} />);

    expect(screen.getByText("14/01/2026")).toBeInTheDocument();
  });
});
```

### 4. Integration Test Example (Form + API)

```typescript
// src/modules/payments/components/TransferForm.test.tsx
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, it, expect, beforeAll, afterEach, afterAll } from "vitest";
import { http, HttpResponse } from "msw";
import { setupServer } from "msw/node";
import { TransferForm } from "./TransferForm";

const server = setupServer(
  http.post("/api/transfers", async ({ request }) => {
    const body = await request.json();

    if (body.amount > 1000) {
      return HttpResponse.json(
        { error: "Amount exceeds limit" },
        { status: 400 }
      );
    }

    return HttpResponse.json({ success: true, transferId: "txn_123" });
  })
);

beforeAll(() => server.listen());
afterEach(() => server.resetHandlers());
afterAll(() => server.close());

describe("TransferForm", () => {
  it("submits transfer successfully", async () => {
    const user = userEvent.setup();
    render(<TransferForm />);

    // Fill form
    await user.type(screen.getByLabelText(/recipient/i), "user@example.com");
    await user.type(screen.getByLabelText(/amount/i), "100.50");

    // Submit
    await user.click(screen.getByRole("button", { name: /transfer/i }));

    // Verify success message
    await waitFor(() => {
      expect(screen.getByText(/transfer successful/i)).toBeInTheDocument();
    });
  });

  it("handles API error", async () => {
    const user = userEvent.setup();
    render(<TransferForm />);

    // Enter amount exceeding limit
    await user.type(screen.getByLabelText(/recipient/i), "user@example.com");
    await user.type(screen.getByLabelText(/amount/i), "1500"); // > 1000 limit

    await user.click(screen.getByRole("button", { name: /transfer/i }));

    await waitFor(() => {
      expect(screen.getByText(/exceeds limit/i)).toBeInTheDocument();
    });
  });
});
```

### 5. E2E Test Example (Payment Flow)

```typescript
// e2e/payment-flow.spec.ts
import { test, expect } from "@playwright/test";

test.describe("Payment Flow", () => {
  test.beforeEach(async ({ page }) => {
    // Login first
    await page.goto("/auth/login");
    await page.fill('input[name="email"]', "test@example.com");
    await page.fill('input[name="password"]', "TestPassword123!");
    await page.click('button:has-text("Sign In")');

    // Wait for redirect to dashboard
    await page.waitForURL("/dashboard");
  });

  test("Happy path: User transfers money", async ({ page }) => {
    // Navigate to payments
    await page.click('a:has-text("Payments")');
    await page.waitForURL("/dashboard/payments");

    // Verify transactions are displayed
    const transactionTable = page.locator("table");
    await expect(transactionTable).toBeVisible();

    // Click transfer button
    await page.click('button:has-text("New Transfer")');
    await page.waitForURL("/dashboard/payments/transfer");

    // Fill transfer form
    await page.fill('input[name="recipientEmail"]', "recipient@example.com");
    await page.fill('input[name="amount"]', "100.50");
    await page.fill(
      'textarea[name="description"]',
      "Payment for services"
    );

    // Submit
    await page.click('button:has-text("Confirm Transfer")');

    // Verify confirmation
    await expect(page.locator("text=Transfer successful")).toBeVisible();

    // Verify transaction appears in history
    await page.waitForURL("/dashboard/payments");
    await expect(
      page.locator("text=recipient@example.com")
    ).toBeVisible();
  });

  test("Error handling: Transfer with insufficient balance", async ({ page }) => {
    await page.click('a:has-text("Payments")');
    await page.click('button:has-text("New Transfer")');

    // Transfer large amount
    await page.fill('input[name="recipientEmail"]', "recipient@example.com");
    await page.fill('input[name="amount"]', "100000"); // Exceeds balance

    await page.click('button:has-text("Confirm Transfer")');

    // Verify error message
    await expect(
      page.locator("text=Insufficient balance")
    ).toBeVisible();
  });

  test("Session timeout: User is logged out after 60 minutes", async ({
    page,
  }) => {
    // Simulate session expiration
    // (In real test: manipulate cookies/local storage)
    await page.context().clearCookies();

    // Attempt to access protected route
    await page.goto("/dashboard/payments");

    // Should redirect to login
    await page.waitForURL("/auth/login");
    expect(page.url()).toContain("/auth/login");
  });
});
```

### 6. Playwright Config

```typescript
// playwright.config.ts
import { defineConfig, devices } from "@playwright/test";

export default defineConfig({
  testDir: "./e2e",
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: "html",
  use: {
    baseURL: "http://localhost:3000",
    trace: "on-first-retry",
  },

  projects: [
    {
      name: "chromium",
      use: { ...devices["Desktop Chrome"] },
    },
    {
      name: "firefox",
      use: { ...devices["Desktop Firefox"] },
    },
    {
      name: "webkit",
      use: { ...devices["Desktop Safari"] },
    },
  ],

  webServer: {
    command: "bun run dev",
    url: "http://localhost:3000",
    reuseExistingServer: !process.env.CI,
  },
});
```

---

## CI/CD Integration

### GitHub Actions Workflow

```yaml
# .github/workflows/test.yml
name: Tests

on: [push, pull_request]

jobs:
  unit:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: oven-sh/setup-bun@v1
      - run: bun install --frozen-lockfile
      - run: bun test:unit
      - run: bun test:unit -- --coverage

  integration:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: oven-sh/setup-bun@v1
      - run: bun install --frozen-lockfile
      - run: bun test:integration

  e2e:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: oven-sh/setup-bun@v1
      - uses: actions/setup-node@v4
        with:
          node-version: "20"
      - run: bun install --frozen-lockfile
      - run: bunx playwright install --with-deps
      - run: bun run build
      - run: bun test:e2e
      - uses: actions/upload-artifact@v4
        if: always()
        with:
          name: playwright-report
          path: playwright-report/
```

### Package.json Scripts

```json
{
  "scripts": {
    "test": "bun test:unit && bun test:integration",
    "test:unit": "vitest run src/**/*.test.ts",
    "test:unit:watch": "vitest watch src/**/*.test.ts",
    "test:integration": "vitest run src/**/*.integration.test.tsx",
    "test:e2e": "playwright test",
    "test:e2e:ui": "playwright test --ui",
    "test:coverage": "vitest run --coverage"
  }
}
```

---

## Consequences

### ✅ Positive Consequences

| Benefit | Impact |
|---------|--------|
| **Fast Feedback** | Unit tests run in < 1s (catch issues early) |
| **High Confidence** | E2E tests verify critical payment flows |
| **Maintainable** | Tests focus on user behavior, not implementation |
| **Scalable** | Test pyramid scales with codebase growth |
| **Cost-Effective** | Unit + integration tests cheap; E2E only for critical paths |

### ⚠️ Negative Consequences & Mitigations

| Challenge | Mitigation |
|-----------|-----------|
| **Test Writing Overhead** | Learning curve for Vitest + Testing Library + Playwright | Pair programming, test examples, templates |
| **MSW Learning** | New mocking library for API calls | Good documentation; start with simple examples |
| **Flaky E2E Tests** | Network delays, timing issues | Explicit waits, retry logic, stable staging env |
| **CI/CD Time** | Tests slow down pipeline | Parallelize jobs; run E2E only for main branch |

---

## Metrics & Success Criteria

- [ ] **Unit Test Coverage:** ≥ 80% for utility functions and validators.
- [ ] **Integration Test Coverage:** All components tested with user scenarios.
- [ ] **E2E Coverage:** All critical payment flows covered (happy path + error scenarios).
- [ ] **Test Speed:** Unit tests < 1s, integration < 5s, E2E < 10s per test.
- [ ] **CI/CD Time:** Full test suite runs in < 10 minutes.
- [ ] **Flakiness:** < 1% of E2E tests are flaky (fail randomly).
- [ ] **Team Adoption:** All PRs include test coverage; 0 coverage regressions.

---

## Testing Maturity Phases

| Phase | Duration | Focus | Coverage |
|-------|----------|-------|----------|
| **Phase 0 (MVP)** | Week 1-4 | Critical payment flows, auth | E2E only for critical |
| **Phase 1** | Week 5+ | Unit + integration tests | 80%+ overall coverage |
| **Phase 2** | Later | Performance tests, load testing | Sustained SLA validation |

---

## Related Decisions

- **ADR-001:** Server Components require server-side test infrastructure (Vitest + test databases).
- **ADR-003:** Componentization affects testing strategy (atoms vs. molecules vs. organisms).

---

## References

- [Vitest Documentation](https://vitest.dev/)
- [React Testing Library](https://testing-library.com/react)
- [MSW (Mock Service Worker)](https://mswjs.io/)
- [Playwright Documentation](https://playwright.dev/)
- [Testing Trophy (Kent C. Dodds)](https://kentcdodds.com/blog/the-testing-trophy-and-testing-javascript)

---

**Author:** Senior Principal Software Architect  
**Created:** 2026-01-14  
**Approved:** Pending  
