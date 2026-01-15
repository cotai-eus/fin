# ADR-003: Atomic Design + Server Components Componentization

**Status:** Accepted  
**Context:** Frontend Architecture & Component Design  
**Date:** 2026-01-14  
**Ratification:** Senior Principal Software Architect  

---

## Problem Statement

As LauraTech scales from MVP to Phase 2, component reusability becomes critical. Previous approaches (Container/Presentational, Compound Components) don't account for React Server Components' unique constraints. We need a componentization strategy that:

1. **Prevents needless client-side hydration** by pushing non-interactive content to Server Components.
2. **Enables component reuse** across multiple domains without coupling business logic.
3. **Maintains clear boundaries** between server and client rendering (prevents hydration mismatches).
4. **Scales to multiple feature teams** without coordination overhead.

**Question:** How do we adapt industry-proven Atomic Design principles to Next.js 16.1's hybrid Server/Client component model?

---

## Context

### Design System Requirements

LauraTech's fintech UI has three categories of components:

1. **UI Atoms** (Button, Input, Card, Badge) — 100% reusable, zero business logic.
2. **Business Patterns** (TransactionsList, TransferForm) — Domain-specific, may have client/server boundary.
3. **Page Layouts** (DashboardLayout, AuthLayout) — Structural, combine atoms + patterns.

### Existing Components

- None formalized; scaffolded pages have no componentization.
- Ory Elements provides pre-built auth forms (UserAuthForm).
- Tailwind CSS + Headless UI provide unstyled foundation.

### Team Constraints

- 3-4 frontend engineers; minimal design system experience.
- No component library (Storybook) yet (Phase 2 consideration).
- Need quick iteration for MVP.

---

## Decision

**Adopt Atomic Design adapted for Server Components (ASC — "Atomic Server Components"), organized in a three-tier hierarchy:**

### Tier 1: Atoms (Presentation Layer)

**Characteristics:**
- Zero business logic.
- Reusable across all domains.
- Live in `src/shared/components/ui/`.
- **Can be RSC or Client Component** depending on interactivity need.

**Examples:** `Button`, `Input`, `Card`, `Modal`, `Badge`, `Table`.

```typescript
// src/shared/components/ui/Button.tsx
// Typically a Client Component (has onClick, onFocus, etc.)
"use client";

interface ButtonProps {
  children: React.ReactNode;
  onClick?: () => void;
  variant?: "primary" | "secondary" | "danger";
  size?: "sm" | "md" | "lg";
  disabled?: boolean;
}

export function Button({
  children,
  onClick,
  variant = "primary",
  size = "md",
  disabled = false,
}: ButtonProps) {
  const baseStyles = "font-semibold rounded-lg transition";
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
      onClick={onClick}
      disabled={disabled}
      className={`${baseStyles} ${variants[variant]} ${sizes[size]} ${disabled ? "opacity-50 cursor-not-allowed" : ""}`}
    >
      {children}
    </button>
  );
}
```

### Tier 2: Molecules (Feature Pattern Layer)

**Characteristics:**
- Combine atoms into reusable patterns.
- **Always Client Component** if they have state or events; **RSC if read-only with server data**.
- Domain-agnostic (can be used in payments, accounts, profile).
- Live in `src/shared/components/` or domain-specific `src/modules/{domain}/components/`.

**Examples:** `SearchBar`, `PaginationControls`, `DataTable`, `FormField`, `TransactionCard`.

```typescript
// src/shared/components/DataTable.tsx
// RSC version (read-only, receives data from parent RSC)

interface DataTableProps<T> {
  columns: Array<{
    header: string;
    accessorKey: string;
    cell?: (value: T) => React.ReactNode;
  }>;
  data: T[];
  striped?: boolean;
}

export function DataTable<T extends Record<string, any>>({
  columns,
  data,
  striped = true,
}: DataTableProps<T>) {
  return (
    <div className="overflow-x-auto">
      <table className="w-full border-collapse">
        <thead>
          <tr className="bg-gray-100">
            {columns.map((col, i) => (
              <th
                key={i}
                className="border px-4 py-2 text-left font-semibold"
              >
                {col.header}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {data.map((row, i) => (
            <tr
              key={i}
              className={striped && i % 2 === 0 ? "bg-gray-50" : ""}
            >
              {columns.map((col, j) => (
                <td key={j} className="border px-4 py-2">
                  {col.cell ? col.cell(row) : String(row[col.accessorKey])}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
```

```typescript
// src/shared/components/SearchBar.tsx
// Client Component (has useState, onChange)
"use client";

import { useState } from "react";
import { Input } from "./ui/Input";
import { Button } from "./ui/Button";

interface SearchBarProps {
  onSearch: (query: string) => void;
  placeholder?: string;
}

export function SearchBar({ onSearch, placeholder = "Search..." }: SearchBarProps) {
  const [query, setQuery] = useState("");

  const handleSearch = () => {
    onSearch(query);
  };

  return (
    <div className="flex gap-2">
      <Input
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        placeholder={placeholder}
      />
      <Button onClick={handleSearch}>Search</Button>
    </div>
  );
}
```

### Tier 3: Organisms (Domain-Specific Layer)

**Characteristics:**
- Combine molecules + atoms into domain-specific features.
- **Always RSC** unless they have significant client-side interactivity.
- Live in domain directories: `src/modules/{domain}/components/`.
- Encapsulate business logic (permissions, data fetching, transformations).

**Examples:** `TransactionsList`, `TransferForm`, `AccountsGrid`, `KYCForm`.

```typescript
// src/modules/payments/components/TransactionsList.tsx
// RSC Organism: Fetches data server-side, passes to Client Component for interactivity

import { getOrySession } from "@/core/ory/session";
import { fetchUserTransactions } from "../actions";
import { DataTable } from "@/shared/components/DataTable";
import { TransactionActions } from "./TransactionActions";

interface TransactionsListProps {
  userId: string;
  page?: number;
}

export async function TransactionsList({
  userId,
  page = 1,
}: TransactionsListProps) {
  // Server-side data fetching (no client bundle impact)
  const transactions = await fetchUserTransactions(userId, page);

  // Define table structure
  const columns = [
    {
      header: "Date",
      accessorKey: "createdAt",
      cell: (tx) => new Date(tx.createdAt).toLocaleDateString("pt-BR"),
    },
    {
      header: "Amount",
      accessorKey: "amount",
      cell: (tx) => `R$ ${tx.amount.toFixed(2)}`,
    },
    {
      header: "Status",
      accessorKey: "status",
      cell: (tx) => <StatusBadge status={tx.status} />,
    },
    {
      header: "Actions",
      accessorKey: "id",
      cell: (tx) => <TransactionActions transactionId={tx.id} />,
    },
  ];

  return (
    <div>
      <DataTable columns={columns} data={transactions} striped />
      <Pagination currentPage={page} totalPages={Math.ceil(transactions.length / 20)} />
    </div>
  );
}
```

```typescript
// src/modules/payments/components/TransactionActions.tsx
// Client Component Organism: Interactive actions for each transaction

"use client";

import { useState } from "react";
import { exportTransaction } from "../actions";
import { Button } from "@/shared/components/ui/Button";
import { Modal } from "@/shared/components/ui/Modal";

interface TransactionActionsProps {
  transactionId: string;
}

export function TransactionActions({ transactionId }: TransactionActionsProps) {
  const [showModal, setShowModal] = useState(false);

  const handleExport = async () => {
    const result = await exportTransaction(transactionId);
    if (result.success) {
      setShowModal(false);
      // Show toast notification
    }
  };

  return (
    <>
      <Button size="sm" onClick={() => setShowModal(true)}>
        Export
      </Button>

      {showModal && (
        <Modal
          title="Export Transaction"
          onClose={() => setShowModal(false)}
          onConfirm={handleExport}
        >
          <p>Export transaction {transactionId} as PDF?</p>
        </Modal>
      )}
    </>
  );
}
```

### Tier 4: Containers (Layout Layer)

**Characteristics:**
- Combine organisms into full-page layouts.
- Often RSCs that orchestrate data fetching + layout structure.
- Live in `src/shared/components/layouts/` or domain-specific folders.

**Examples:** `DashboardLayout`, `AuthLayout`, `SettingsLayout`.

```typescript
// src/shared/components/layouts/DashboardLayout.tsx
// RSC: Fetches session, renders sidebar + content

import { getOrySession } from "@/core/ory/session";
import { Sidebar } from "./Sidebar";
import { TopNav } from "./TopNav";

interface DashboardLayoutProps {
  children: React.ReactNode;
}

export async function DashboardLayout({ children }: DashboardLayoutProps) {
  const session = await getOrySession();

  return (
    <div className="flex h-screen">
      <Sidebar user={session.identity} />
      <div className="flex-1 flex flex-col">
        <TopNav user={session.identity} />
        <main className="flex-1 overflow-y-auto p-6">{children}</main>
      </div>
    </div>
  );
}
```

---

## Componentization Flowchart

```
┌─────────────────────────────────────────┐
│          Page (RSC by default)          │  src/app/(dashboard)/payments/page.tsx
│  - Fetches session, permissions, data   │
│  - Orchestrates page layout              │
│  - No interactivity                      │
└──────────────────┬──────────────────────┘
                   │
        ┌──────────┴──────────┐
        ▼                     ▼
┌──────────────┐      ┌──────────────┐
│  Container   │      │  Container   │
│  (RSC)       │      │  (Client)    │
│  Layout      │      │  Sidebar     │
└──────┬───────┘      └──────┬───────┘
       │                     │
   ┌───┴────┬────────────┬───┴────┬─────────┐
   ▼        ▼            ▼        ▼         ▼
 Molecule Molecule  Molecule  Molecule  Organism
 (RSC)    (Client)  (RSC)     (Client)  (Client)
DataTable Search   CardList Filters   Actions
   │        │         │         │         │
   └────────┴─────────┴─────────┴─────────┘
            ▼
       ┌────────────┐
       │   Atoms    │
       │ (Shared)   │
       │Button, etc │
       └────────────┘
```

---

## Justification

### 1. Server Components as Default

Pushing components to RSC by default:
- **Eliminates unnecessary client-side code.**
- **Improves security** (business logic server-only).
- **Reduces hydration mismatch** (less client rendering).

### 2. Atomic Design Separation of Concerns

Clear tier separation enables:
- **Team Scalability:** Multiple teams work on different organisms without collision.
- **Code Reuse:** Atoms + molecules reused across domains.
- **Testability:** Each tier has clear testing strategy (Tier 1: Unit, Tier 2-3: Integration, Tier 4: E2E).

### 3. Server/Client Boundary Clarity

`"use client"` is explicit at each level:
- Atoms: May be client (interactive) or server (presentational).
- Molecules: Client only if they have state; otherwise RSC.
- Organisms: RSC by default, client subtree only for interactivity.
- Containers: RSC (layout, data orchestration).

This prevents accidental client-side expansion.

### 4. Scalability to Feature Teams

With clear boundaries, teams can:

```
Team A (Payments)        Team B (Accounts)       Team C (Dashboard)
├── src/modules/         ├── src/modules/         ├── src/modules/
│   payments/            │   accounts/            │   dashboard/
│   ├── components/      │   ├── components/      │   ├── components/
│   ├── actions/         │   ├── actions/         │   ├── actions/
│   └── validators.ts    │   └── validators.ts    │   └── validators.ts
│                        │                        │
└── Shared components ◄──┴────────────────────────┘
    src/shared/components/
    (Atoms + Molecules shared across all teams)
```

---

## Alternatives Considered

### Alternative 1: Container/Presentational (Legacy Pattern)

| Aspect | Evaluation |
|--------|------------|
| Familiar | ✅ Most React devs know it |
| Server Components | ❌ Doesn't account for RSCs |
| Hydration | ❌ Increases hydration mismatches |
| Bundle | ❌ No bundle benefit |
| **Decision** | ❌ **Rejected** — Outdated for RSC era |

### Alternative 2: Feature-Based Folder Structure (No Atomic Hierarchy)

```
src/
├── features/
│   ├── payments/
│   │   ├── components/
│   │   │   ├── Button.tsx        ❌ Ties UI to domain
│   │   │   ├── TransactionsList
│   │   │   └── PaymentForm
│   │   └── actions/
│   └── accounts/
│       ├── components/
│       │   ├── Button.tsx        ❌ Duplicate Button component
│       │   └── AccountCard
│       └── actions/
```

| Aspect | Evaluation |
|--------|------------|
| Isolation | ✅ Features isolated |
| Reusability | ❌ No sharing of UI atoms |
| Duplication | ❌ Components redefined per domain |
| Scalability | ❌ Harder for shared UI libraries |
| **Decision** | ❌ **Rejected** — Limits reusability |

### Alternative 3: Strictly Separated Client Library + Server Layer

| Aspect | Evaluation |
|--------|------------|
| Clarity | ✅ Clear boundary (client vs. server) |
| Flexibility | ❌ Rigid; hard to mix RSCs with client |
| Pragmatism | ❌ Overcomplicates simple molecules |
| **Decision** | ❌ **Rejected** — Too dogmatic |

**Recommendation:** ✅ **Atomic Server Components (ASC)** provides best balance of clarity, reusability, and pragmatism.

---

## Implementation Guidelines

### Rule 1: RSC-First Mentality

All components start as RSC unless interactivity is required:

```typescript
// ✅ Start here: RSC
export async function TransactionsList() {
  const data = await fetch(...);
  return <Table data={data} />;
}

// If you need state/events, add "use client" to **client boundary only**:
// ❌ DON'T do this (moves entire component to client)
"use client";
export function TransactionsList() { ... }

// ✅ DO this instead (only interactive subset is client)
function TransactionsList() { // Still RSC
  return (
    <>
      <DataTable data={data} />
      <FilterControls onFilter={...} />  {/* This is the "use client" boundary */}
    </>
  );
}
```

### Rule 2: `"use client"` at Leaf Level

`"use client"` directives appear only at the smallest interactive subtree:

```typescript
// src/modules/payments/components/TransactionsList.tsx
// RSC: Doesn't need "use client"
export function TransactionsList() {
  return (
    <>
      <DataTable data={transactions} />
      <FilterBar onFilter={handleFilter} />
      {/* ☝️ FilterBar is Client Component */}
    </>
  );
}

// src/modules/payments/components/FilterBar.tsx
"use client";  // ← Only here, at the leaf
export function FilterBar({ onFilter }) {
  const [filters, setFilters] = useState();
  return (
    <div>
      <input onChange={(e) => setFilters(e.target.value)} />
      <button onClick={() => onFilter(filters)}>Filter</button>
    </div>
  );
}
```

### Rule 3: Props Are the Bridge

Client Components accept data via props (no direct database access):

```typescript
// ✅ Correct: Data flows from RSC parent via props
export async function PaymentPage() {
  const transactions = await db.query(...);  // Server-side
  return <TransactionsList data={transactions} />;  {/* Pass as prop */}
}

// ❌ Wrong: Client Component accessing database
"use client";
export function TransactionsList() {
  useEffect(() => {
    const data = await fetch(...);  // Exposed to client network
  }, []);
}
```

### Rule 4: Naming Conventions

```
src/
├── app/                                    # Next.js routes
├── shared/
│   ├── components/
│   │   ├── ui/                            # Atoms (Button, Input, Card, etc.)
│   │   │   └── *.tsx
│   │   ├── layouts/                       # Container-level layouts
│   │   │   └── DashboardLayout.tsx
│   │   └── [Molecules if domain-agnostic]
│   ├── hooks/
│   │   └── use*.ts                        # Generic hooks (useDebounce, etc.)
│   └── utils/
│       └── formatters/                    # Shared utilities
│
├── modules/
│   ├── payments/                          # Domain: Payments
│   │   ├── components/                    # Organisms + specialized molecules
│   │   │   ├── TransactionsList.tsx       # RSC Organism
│   │   │   ├── TransactionActions.tsx     # Client Organism
│   │   │   └── TransactionCard.tsx        # Molecule (reusable within domain)
│   │   ├── actions/
│   │   │   └── index.ts                   # Server Actions
│   │   ├── hooks/
│   │   │   └── useTransactions.ts         # Domain-specific hook
│   │   └── validators.ts
│   ├── accounts/
│   └── dashboard/
│
└── core/                                  # Infrastructure
    ├── api/
    ├── ory/
    └── config/
```

---

## Consequences

### ✅ Positive Consequences

| Benefit | Impact |
|---------|--------|
| **Bundle Size** | Atoms + read-only molecules stay server-side |
| **Reusability** | Shared atoms across all domains |
| **Clarity** | Clear RSC/Client boundaries |
| **Scalability** | Multiple teams work in parallel |
| **Testability** | Each tier has specific test strategy |

### ⚠️ Negative Consequences & Mitigations

| Challenge | Mitigation |
|-----------|-----------|
| **Learning Curve** | RSC patterns unfamiliar to CSR-trained devs | Pair programming, documentation |
| **Naming Complexity** | 4 tiers (atoms, molecules, organisms, containers) | Clear naming convention, linter rules |
| **Testing Atoms** | Unit testing atoms in isolation | Vitest + React Testing Library |

---

## Testing Strategy (Per Tier)

| Tier | Test Type | Tool | Example |
|------|-----------|------|---------|
| **Atoms** | Unit | Vitest + Testing Library | Button click, input change |
| **Molecules** | Unit + Integration | Vitest + Testing Library | SearchBar searches correctly |
| **Organisms** | Integration | Testing Library + MSW | TransactionsList fetches and renders |
| **Containers** | E2E | Playwright | Full page flow (login → view transactions) |

---

## Metrics & Success Criteria

- [ ] All atoms tested with Vitest (100% coverage target).
- [ ] No `"use client"` above leaf level (lint rule enforced).
- [ ] Zero hydration mismatches in production.
- [ ] New features developed 20% faster (using atomic components).
- [ ] Component reuse rate > 80% (atoms reused across ≥ 80% of pages).

---

## Related Decisions

- **ADR-001:** Server Components are the rendering paradigm; componentization must respect this.
- **ADR-004:** Testing strategy aligned to tier types (unit for atoms, integration for molecules, E2E for pages).

---

## References

- [Atomic Design Methodology](https://atomicdesign.bradfrost.com/)
- [Next.js Server Components](https://nextjs.org/docs/app/building-your-application/rendering/server-components)
- [React Composition Patterns](https://reactpatterns.com/)

---

**Author:** Senior Principal Software Architect  
**Created:** 2026-01-14  
**Approved:** Pending  
