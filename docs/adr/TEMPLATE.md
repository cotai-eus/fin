# ADR {NNNN}: {Short Title}

**Status:** [Proposed | Accepted | Deprecated | Superseded]  
**Context:** {Category - e.g., Frontend Architecture, Data Management, Security, Observability}  
**Date:** {YYYY-MM-DD}  
**Supersedes:** [ADR-XXXX if applicable]  
**Ratification:** [Name/Team, e.g., "Senior Principal Architect"]  

---

## Problem Statement

Concise description of the challenge or question being addressed. Include business context and constraints that necessitate this decision.

Example:
> Next.js 15 offers client-side rendering (CSR), server-side rendering (SSR), and static generation (SSG). Each has tradeoffs in build time, runtime performance, and operational complexity. We need to choose a primary rendering paradigm that optimizes for development velocity, bundle size, and fintech-grade security for the LauraTech MVP.

---

## Context

Deeper explanation of the situation:

1. **Background:** What led to this decision? What problem were we trying to solve?
2. **Constraints:** Budget, timeline, team expertise, regulatory requirements, performance SLAs.
3. **Scope:** Which parts of the system does this apply to?

Example:
> - Team has 2 years of React experience; Node.js backend experience varies.
> - Must support HIPAA/PCI compliance requirements (sensitive financial data).
> - Target bundle size: < 100KB gzipped (Core Web Vitals optimization).
> - MVP delivery in 4-7 weeks.

---

## Decision

**The chosen approach.** Be specific and actionable.

Example:
> We will adopt **Next.js 16.1 with React Server Components (RSCs)** as the primary rendering paradigm. Server Components will be the default for all new pages unless explicit client-side interactivity is required. This decision includes:
>
> 1. All page components are RSCs by default.
> 2. `"use client"` boundary is added only for interactive subtrees (forms, charts, filters).
> 3. Server Actions (asynchronous mutations) are the standard for form submissions.
> 4. Client-side state is minimized; server-driven state is preferred.

---

## Justification (Why This Decision)

Explain the reasoning behind the choice. Address:

1. **Advantages:** What benefits does this bring?
2. **Alignment with Goals:** How does this support project objectives?
3. **Long-term Implications:** Is this future-proof?

Example:
> - **Bundle Size Reduction:** RSCs eliminate the need to ship client-side framework code (React Router, state management libraries) for read-only pages. Estimated 40-60% reduction in JavaScript.
> - **Security:** Sensitive operations (database queries, API calls) execute server-side only. No credential exposure in client bundles.
> - **Developer Experience:** Fewer mental models required; developers write synchronous server code and async client code separately.
> - **SEO:** RSCs are inherently server-rendered, improving Core Web Vitals and search ranking.

---

## Alternatives Considered

List other options evaluated and why they were rejected.

| Alternative | Pros | Cons | Decision |
|-------------|------|------|----------|
| CSR (React SPA) | Instant interactivity, familiar patterns | Large JS bundle, SEO challenges, fintech security risk | Rejected |
| SSR (Next.js 15) | Good SEO, moderate bundle | Higher server load, harder to scale, less modern | Rejected |
| RSC (Selected) | Optimal bundle, security-first, modern DX | Learning curve, less mature ecosystem | **Accepted** |
| Remix | Full-stack framework, excellent DX | Overkill for MVP, less widespread adoption | Rejected |
| Astro | Excellent static generation | Not suitable for interactive fintech UIs | Rejected |

---

## Implementation Details

Specific technical guidance for implementing this decision.

Example:
> 1. All page files in `src/app/` are `.tsx` RSCs by default.
> 2. Place `"use client"` at the top of a file only when you need React hooks (useState, useEffect) or browser APIs.
> 3. Use `async` components for server-side data fetching:
>    ```tsx
>    export default async function Page() {
>      const data = await fetch(...);
>      return <div>{data}</div>;
>    }
>    ```
> 4. Use Server Actions for mutations:
>    ```tsx
>    "use server";
>    export async function submitForm(formData) { ... }
>    ```
> 5. Use `Suspense` for progressive rendering:
>    ```tsx
>    <Suspense fallback={<Skeleton />}>
>      <SlowComponent />
>    </Suspense>
>    ```

---

## Consequences

What are the trade-offs and implications?

### Positive Consequences
- 🎯 Smaller bundle size → faster initial page load.
- 🔒 Enhanced security posture → sensitive data never touches client.
- ⚡ Reduced client-side re-renders → improved performance.
- 📚 Clearer separation of concerns.

### Negative Consequences
- 📚 Steeper learning curve for teams new to RSCs.
- 🔄 Limited real-time capabilities without WebSocket/polling.
- 🧪 Testing RSCs requires server-side test infrastructure.

### Mitigation Strategies
- Provide extensive documentation and code examples.
- Pair program with team members during adoption phase.
- Establish clear guidelines for RSC vs. Client Component boundaries.

---

## Metrics & Success Criteria

How will we measure if this decision was sound?

- [ ] **Bundle Size:** Core JavaScript < 100KB gzipped.
- [ ] **LCP (Largest Contentful Paint):** < 2.5 seconds (Core Web Vitals threshold).
- [ ] **Team Velocity:** No > 20% regression in sprint velocity during adoption phase.
- [ ] **Security Audits:** Zero credential exposures in client bundles.
- [ ] **User Satisfaction:** No increase in error reports related to rendering/hydration issues.

---

## Related Decisions

Links to dependent or related ADRs:

- [ADR-002: Ory Zero-Trust Identity Architecture](./ADR-002-ory-zero-trust.md) — Session management relies on RSC server-side verification.
- [ADR-003: Atomic Design + Server Components](./ADR-003-componentization.md) — Component boundaries must respect RSC/Client boundaries.
- [ADR-004: Tri-Layer Testing Strategy](./ADR-004-testing.md) — Testing approach adapted for RSC serverside code.

---

## Timeline & Phases

How will this be rolled out?

- **Phase 1 (Week 1-2):** Establish RSC patterns and documentation.
- **Phase 2 (Week 3-4):** Migrate existing pages to RSCs.
- **Phase 3 (Week 5+):** New features built on RSC-first paradigm.

---

## Document Metadata

| Field | Value |
|-------|-------|
| **Author** | Senior Principal Software Architect |
| **Created** | 2026-01-14 |
| **Last Updated** | 2026-01-14 |
| **Reviewed By** | [Pending] |
| **Approved By** | [Pending] |

---

## Questions & Answers

**Q: When should we use client components?**  
A: When you need React hooks (useState, useEffect, useContext), browser APIs (localStorage, geolocation), or event handlers that require interactivity.

**Q: Can RSCs call databases directly?**  
A: Yes, RSCs run on the server and can access databases, APIs, and secrets securely.

**Q: What about real-time updates?**  
A: For real-time, consider WebSockets or Server-Sent Events (SSE) within client components, or use TanStack Query for polling-based updates.

---

## References

- [Next.js Server Components Documentation](https://nextjs.org/docs/app/building-your-application/rendering/server-components)
- [React RFC: Server Components](https://github.com/reactjs/rfcs/blob/main/text/0188-server-components.md)
- [Web Vitals Optimization Guide](https://web.dev/vitals/)
