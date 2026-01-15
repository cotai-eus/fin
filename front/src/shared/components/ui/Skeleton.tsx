/**
 * src/shared/components/ui/Skeleton.tsx
 *
 * Skeleton Atom Component
 *
 * Loading placeholder skeleton, commonly used with Suspense
 */

"use client";

interface SkeletonProps {
  count?: number;
  height?: string;
  className?: string;
}

export function Skeleton({
  count = 1,
  height = "h-4",
  className = "",
}: SkeletonProps) {
  return (
    <div className="space-y-2">
      {[...Array(count)].map((_, i) => (
        <div
          key={i}
          className={`${height} bg-gray-200 rounded animate-pulse ${className}`}
        />
      ))}
    </div>
  );
}
