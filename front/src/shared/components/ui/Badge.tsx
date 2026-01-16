/**
 * src/shared/components/ui/Badge.tsx
 *
 * Badge Atom Component
 *
 * Status indicator badge with color-coded variants
 */

"use client";

interface BadgeProps {
  variant?: "pending" | "completed" | "failed" | "cancelled";
  status?: "pending" | "completed" | "failed" | "cancelled";
  label?: string;
  children?: React.ReactNode;
  className?: string;
}

export function Badge({ variant, status, label, children, className = "" }: BadgeProps) {
  const statusValue = variant || status || "cancelled";
  
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

  const variantStyles = variants[statusValue];
  const displayLabel = children || label || statusValue.charAt(0).toUpperCase() + statusValue.slice(1);

  return (
    <div
      className={`inline-flex items-center gap-2 px-3 py-1 rounded-full border ${variantStyles.bg} ${variantStyles.text} ${variantStyles.border} text-xs font-semibold ${className}`}
    >
      <div className={`w-2 h-2 rounded-full ${variantStyles.dot}`} />
      {displayLabel}
    </div>
  );
}
