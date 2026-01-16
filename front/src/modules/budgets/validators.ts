/**
 * Budgets Module - Zod Validation Schemas
 * Validates budget operations
 */

import { z } from "zod";
import { BudgetCategory, BudgetPeriod, AlertThreshold } from "./types";

/**
 * Create Budget Schema
 */
export const createBudgetSchema = z.object({
  category: z.nativeEnum(BudgetCategory),
  period: z.nativeEnum(BudgetPeriod),
  limit: z
    .number()
    .positive("Budget limit must be positive")
    .max(1_000_000, "Budget limit cannot exceed R$ 1,000,000"),
  alertThreshold: z.nativeEnum(AlertThreshold).default(AlertThreshold.NINETY_PERCENT),
  alertsEnabled: z.boolean().default(true),
  startDate: z.string().datetime().optional(),
});

export type CreateBudgetInput = z.infer<typeof createBudgetSchema>;

/**
 * Update Budget Schema
 */
export const updateBudgetSchema = z.object({
  budgetId: z.string().uuid("Invalid budget ID"),
  limit: z
    .number()
    .positive("Budget limit must be positive")
    .max(1_000_000, "Budget limit cannot exceed R$ 1,000,000")
    .optional(),
  alertThreshold: z.nativeEnum(AlertThreshold).optional(),
  alertsEnabled: z.boolean().optional(),
});

export type UpdateBudgetInput = z.infer<typeof updateBudgetSchema>;

/**
 * Delete Budget Schema
 */
export const deleteBudgetSchema = z.object({
  budgetId: z.string().uuid("Invalid budget ID"),
});

export type DeleteBudgetInput = z.infer<typeof deleteBudgetSchema>;

/**
 * Get Spending Analysis Schema
 */
export const getSpendingAnalysisSchema = z.object({
  startDate: z.string().datetime(),
  endDate: z.string().datetime(),
  category: z.nativeEnum(BudgetCategory).optional(),
});

export type GetSpendingAnalysisInput = z.infer<typeof getSpendingAnalysisSchema>;

/**
 * Get Category Budget Status
 * Calculates percentage and status
 */
export function getCategoryBudgetStatus(
  spent: number,
  limit: number
): {
  percentage: number;
  status: "safe" | "warning" | "danger" | "exceeded";
  color: string;
} {
  const percentage = (spent / limit) * 100;

  if (percentage >= 100) {
    return { percentage, status: "exceeded", color: "red" };
  } else if (percentage >= 90) {
    return { percentage, status: "danger", color: "red" };
  } else if (percentage >= 75) {
    return { percentage, status: "warning", color: "yellow" };
  } else {
    return { percentage, status: "safe", color: "green" };
  }
}

/**
 * Category Labels for Display
 */
export const CATEGORY_LABELS: Record<BudgetCategory, string> = {
  [BudgetCategory.FOOD]: "Alimentação",
  [BudgetCategory.TRANSPORT]: "Transporte",
  [BudgetCategory.ENTERTAINMENT]: "Lazer",
  [BudgetCategory.SHOPPING]: "Compras",
  [BudgetCategory.BILLS]: "Contas",
  [BudgetCategory.HEALTH]: "Saúde",
  [BudgetCategory.EDUCATION]: "Educação",
  [BudgetCategory.OTHER]: "Outros",
};

/**
 * Category Icons (emoji)
 */
export const CATEGORY_ICONS: Record<BudgetCategory, string> = {
  [BudgetCategory.FOOD]: "🍔",
  [BudgetCategory.TRANSPORT]: "🚗",
  [BudgetCategory.ENTERTAINMENT]: "🎬",
  [BudgetCategory.SHOPPING]: "🛍️",
  [BudgetCategory.BILLS]: "💡",
  [BudgetCategory.HEALTH]: "⚕️",
  [BudgetCategory.EDUCATION]: "📚",
  [BudgetCategory.OTHER]: "📦",
};
