/**
 * src/app/(dashboard)/payments/page.tsx
 *
 * Payment History Dashboard Page
 * 
 * This is a Server Component (RSC) that demonstrates best practices:
 * - Session verification via Ory middleware + getOrySession()
 * - Server-side data fetching with Next.js cache directives
 * - Progressive rendering with Suspense
 * - Streaming for better perceived performance
 * - Proper error handling and loading states
 */

import { Suspense } from "react";
import { redirect } from "next/navigation";
import { getOrySession } from "@/core/ory/session";
import { fetchUserTransactions } from "@/modules/payments/actions";
import { TransactionsList } from "@/modules/payments/components/TransactionsList";
import { TransactionsSkeleton } from "./components/TransactionsSkeleton";
import { PageHeader } from "@/shared/components/PageHeader";
import { ErrorBoundary } from "@/shared/components/ErrorBoundary";

interface PaymentsPageProps {
  searchParams: {
    page?: string;
    status?: "pending" | "completed" | "failed";
  };
}

/**
 * Main payments page (RSC)
 * Executed entirely on the server; no client-side JavaScript needed for initial render
 */
export default async function PaymentsPage({
  searchParams,
}: PaymentsPageProps) {
  // 1. Verify session (middleware guarantees, but we verify again for defense-in-depth)
  const session = await getOrySession();

  if (!session) {
    // Redirect to login if not authenticated
    // Middleware should normally prevent reaching here, but we redirect defensively
    redirect("/auth/login");
  }

  const userId = session.identity?.id;
  if (!userId) {
    throw new Error("Invalid session: No user ID");
  }

  const page = parseInt(searchParams.page || "1", 10);

  return (
    <div className="space-y-6 pb-10">
      {/* Page Header */}
      <PageHeader
        title="Payment History"
        subtitle={`Account: ${session.identity?.traits?.email || "User"}`}
        description="View and manage your recent transactions"
      />

      {/* Error Boundary + Suspense for progressive rendering */}
      <ErrorBoundary>
        <Suspense fallback={<TransactionsSkeleton />}>
          <TransactionsContent userId={userId} page={page} />
        </Suspense>
      </ErrorBoundary>
    </div>
  );
}

/**
 * Separated component to leverage Suspense boundaries
 * This allows the page to stream data as it becomes available
 */
async function TransactionsContent({
  userId,
  page,
}: {
  userId: string;
  page: number;
}) {
  let transactions;

  try {
    // Fetch data (this is the slow operation)
    transactions = await fetchUserTransactions(userId, page);
  } catch (error) {
    console.error("[Payments] Failed to load transactions:", error);

    return (
      <div className="rounded-lg border border-red-200 bg-red-50 p-6 text-red-800">
        <h3 className="font-semibold">Unable to Load Transactions</h3>
        <p className="mt-2 text-sm">
          {error instanceof Error
            ? error.message
            : "An unexpected error occurred. Please try again."}
        </p>
        <button
          onClick={() => window.location.reload()}
          className="mt-4 rounded bg-red-600 px-4 py-2 text-white hover:bg-red-700"
        >
          Retry
        </button>
      </div>
    );
  }

  return (
    <TransactionsList
      userId={userId}
      initialData={transactions}
      currentPage={page}
    />
  );
}

/**
 * Metadata for SEO and page title
 */
export const metadata = {
  title: "Payment History | LauraTech",
  description: "View your transaction history and payment details",
};

/**
 * ISR Configuration: Revalidate every 60 seconds
 * This means the page is cached for 60s, then regenerated on next request
 * Balance between fresh data and server performance
 */
export const revalidate = 60;

/**
 * Only allow authenticated users to access this page
 * (middleware will enforce this, but this documents the requirement)
 */
export const dynamic = "force-dynamic"; // Don't cache authenticated pages
