/**
 * src/modules/payments/components/FilterBar.tsx
 *
 * Filter Bar Molecule Component
 *
 * Allows users to filter transactions by status
 * Example of a pure client component with local state management
 */

"use client";

import { Button } from "@/shared/components/ui/Button";

interface FilterBarProps {
  value: "all" | "completed" | "pending" | "failed";
  onChange: (value: "all" | "completed" | "pending" | "failed") => void;
}

const FILTERS = [
  { id: "all", label: "All", color: "gray" },
  { id: "completed", label: "Completed", color: "green" },
  { id: "pending", label: "Pending", color: "yellow" },
  { id: "failed", label: "Failed", color: "red" },
] as const;

export function FilterBar({ value, onChange }: FilterBarProps) {
  return (
    <div className="flex flex-wrap gap-2">
      {FILTERS.map((filter) => (
        <Button
          key={filter.id}
          onClick={() =>
            onChange(filter.id as "all" | "completed" | "pending" | "failed")
          }
          variant={value === filter.id ? "primary" : "outline"}
          size="sm"
        >
          {filter.label}
        </Button>
      ))}
    </div>
  );
}
