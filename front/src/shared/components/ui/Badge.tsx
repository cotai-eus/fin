/**
 * src/shared/components/ui/Badge.tsx
 *
 * Badge Atom Component
 *
 * Status indicator badge with color-coded variants
 */

"use client";

interface BadgeProps {
  status: "pending" | "completed" | "failed" | "cancelled";
  label?: string;
}

export function Badge({ status, label }: BadgeProps) {
  const variants = {
    pending: {
      bg: "bg-yellow-100",
      text: "text-yellow-800",
      border: "border-yellow-300",
      dot: "bg-yellow-500",
    },
    completed: {
      bg: "bg-green-100",
      text: "text-green-800",
      border: "border-green-300",
      dot: "bg-green-500",
    },
    failed: {
      bg: "bg-red-100",
      text: "text-red-800",
      border: "border-red-300",
      dot: "bg-red-500",
    },
    cancelled: {
      bg: "bg-gray-100",
      text: "text-gray-800",
      border: "border-gray-300",
      dot: "bg-gray-500",
    },
  };

  const variant = variants[status];
  const displayLabel = label || status.charAt(0).toUpperCase() + status.slice(1);

  return (
    <div
      className={`inline-flex items-center gap-2 px-3 py-1 rounded-full border ${variant.bg} ${variant.text} ${variant.border} text-xs font-semibold`}
    >
      <div className={`h-2 w-2 rounded-full ${variant.dot}`} />
      {displayLabel}
    </div>
  );
}
