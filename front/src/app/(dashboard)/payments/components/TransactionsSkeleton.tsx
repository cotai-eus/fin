/**
 * src/app/(dashboard)/payments/components/TransactionsSkeleton.tsx
 *
 * Loading skeleton for transactions list
 * Shown while data is being fetched via Suspense fallback
 */

import { Card } from "@/shared/components/ui/Card";
import { Skeleton } from "@/shared/components/ui/Skeleton";

export function TransactionsSkeleton() {
  return (
    <div className="space-y-4">
      {/* Filter bar skeleton */}
      <div className="flex gap-2">
        {[...Array(4)].map((_, i) => (
          <div key={i} className="h-10 w-24 bg-gray-200 rounded animate-pulse" />
        ))}
      </div>

      {/* Transactions list skeleton */}
      <Card>
        {[...Array(5)].map((_, i) => (
          <div
            key={i}
            className="flex items-center justify-between border-b border-gray-200 px-6 py-4"
          >
            <div className="flex items-center gap-4 flex-1">
              <div className="h-10 w-10 bg-gray-200 rounded-full animate-pulse" />
              <div className="space-y-2 flex-1">
                <Skeleton count={2} />
              </div>
            </div>
            <div className="text-right">
              <div className="h-6 w-24 bg-gray-200 rounded animate-pulse mb-2" />
              <div className="h-4 w-20 bg-gray-200 rounded animate-pulse" />
            </div>
          </div>
        ))}
      </Card>
    </div>
  );
}
