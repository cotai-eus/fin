/**
 * src/modules/payments/components/TransactionsList.tsx
 *
 * Transactions List Organism Component
 *
 * This demonstrates the Atomic Design pattern with Server Components:
 * - Parent (this component) = RSC receiving data as props
 * - Child components (TransactionCard, FilterBar) = Client components for interactivity
 * - This component is a "Smart Container" that fetches data and passes to children
 */

"use client";

import { useState } from "react";
import { TransactionCard } from "./TransactionCard";
import { FilterBar } from "./FilterBar";
import { Skeleton } from "@/shared/components/ui/Skeleton";
import { Card } from "@/shared/components/ui/Card";
import { type TransactionList as TransactionListType } from "../types";

interface TransactionsListProps {
  userId: string;
  initialData: TransactionListType;
  currentPage: number;
}

/**
 * TransactionsList renders filtered, paginated transaction history
 * 
 * Pattern breakdown:
 * 1. Receives pre-fetched data from RSC parent (initialData)
 * 2. Manages client-side filtering state (status, dateRange)
 * 3. Handles pagination by calling Server Action to fetch new page
 * 4. Renders TransactionCard for each transaction (molecule component)
 * 5. Renders FilterBar for filtering (molecule component)
 */
export function TransactionsList({
  userId,
  initialData,
  currentPage,
}: TransactionsListProps) {
  const [filterStatus, setFilterStatus] = useState<
    "all" | "completed" | "pending" | "failed"
  >("all");

  const [isLoadingMore, setIsLoadingMore] = useState(false);

  // Filter transactions client-side based on status
  const filteredTransactions =
    filterStatus === "all"
      ? initialData.data
      : initialData.data.filter((t) => t.status === filterStatus);

  const handleLoadMore = async () => {
    setIsLoadingMore(true);
    try {
      // In a real implementation, this would call a Server Action
      // to fetch the next page of data
      // await fetchUserTransactions(userId, currentPage + 1);
      
      // For now, just show that pagination exists
      window.location.href = `?page=${currentPage + 1}`;
    } catch (error) {
      console.error("Failed to load more transactions:", error);
    } finally {
      setIsLoadingMore(false);
    }
  };

  return (
    <div className="space-y-4">
      {/* Filter bar (molecule component) */}
      <FilterBar value={filterStatus} onChange={setFilterStatus} />

      {/* Transactions list */}
      <Card className="overflow-hidden">
        {filteredTransactions.length === 0 ? (
          <div className="flex flex-col items-center justify-center px-6 py-12 text-center">
            <svg
              className="h-12 w-12 text-gray-400"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
              />
            </svg>
            <h3 className="mt-4 font-semibold text-gray-900">
              No transactions found
            </h3>
            <p className="mt-2 text-sm text-gray-600">
              {filterStatus === "all"
                ? "Start making transfers to see them here"
                : `No ${filterStatus} transactions`}
            </p>
          </div>
        ) : (
          <div className="divide-y divide-gray-200">
            {filteredTransactions.map((transaction) => (
              <TransactionCard
                key={transaction.id}
                transaction={transaction}
                userId={userId}
              />
            ))}
          </div>
        )}
      </Card>

      {/* Pagination */}
      {initialData.totalPages > 1 && (
        <div className="flex items-center justify-between">
          <p className="text-sm text-gray-600">
            Showing page {currentPage} of {initialData.totalPages}
          </p>
          <button
            onClick={handleLoadMore}
            disabled={isLoadingMore || currentPage >= initialData.totalPages}
            className="rounded-lg bg-blue-600 px-4 py-2 text-sm font-semibold text-white hover:bg-blue-700 disabled:opacity-50"
          >
            {isLoadingMore ? "Loading..." : "Load More"}
          </button>
        </div>
      )}
    </div>
  );
}

/**
 * Loading skeleton while data is being fetched
 * Used by Suspense fallback
 */
export function TransactionsListSkeleton() {
  return (
    <div className="space-y-4">
      <div className="h-10 bg-gray-200 rounded" />
      <Card className="space-y-2">
        {[...Array(5)].map((_, i) => (
          <div key={i} className="border-b border-gray-200 p-4">
            <Skeleton count={3} />
          </div>
        ))}
      </Card>
    </div>
  );
}
