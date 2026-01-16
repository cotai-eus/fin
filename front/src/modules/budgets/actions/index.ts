/**
 * Budgets Module - Server Actions
 * Handles budget management and spending analysis
 */

"use server";

import { revalidatePath } from "next/cache";
import { requireOrySession } from "@/core/ory/session";
import {
  createBudgetSchema,
  updateBudgetSchema,
  deleteBudgetSchema,
  getSpendingAnalysisSchema,
  type CreateBudgetInput,
  type UpdateBudgetInput,
  type DeleteBudgetInput,
  type GetSpendingAnalysisInput,
} from "../validators";
import type { Budget, BudgetSummary, CategorySpending, SpendingTrend } from "../types";

const BACKEND_URL = process.env.BACKEND_API_URL || "http://localhost:8080";

type ActionResult<T> =
  | { success: true; data: T }
  | { success: false; error: string };

/**
 * Create Budget
 */
export async function createBudget(
  input: unknown
): Promise<ActionResult<Budget>> {
  try {
    const session = await requireOrySession();
    const userId = session.identity?.id;
    if (!userId) {
      return { success: false, error: "Unauthorized" };
    }

    const validated = createBudgetSchema.safeParse(input);
    if (!validated.success) {
      return {
        success: false,
        error: validated.error.issues[0]?.message || "Invalid input",
      };
    }

    const response = await fetch(`${BACKEND_URL}/api/budgets`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-User-ID": userId,
        "X-Request-ID": crypto.randomUUID(),
      },
      body: JSON.stringify(validated.data),
    });

    if (!response.ok) {
      const error = await response
        .json()
        .catch(() => ({ message: "Budget creation failed" }));
      return { success: false, error: error.message || "Budget creation failed" };
    }

    const budget: Budget = await response.json();
    revalidatePath("/dashboard");

    return { success: true, data: budget };
  } catch (error) {
    console.error("Create budget error:", error);
    return { success: false, error: "An unexpected error occurred" };
  }
}

/**
 * Update Budget
 */
export async function updateBudget(
  input: unknown
): Promise<ActionResult<Budget>> {
  try {
    const session = await requireOrySession();
    const userId = session.identity?.id;
    if (!userId) {
      return { success: false, error: "Unauthorized" };
    }

    const validated = updateBudgetSchema.safeParse(input);
    if (!validated.success) {
      return {
        success: false,
        error: validated.error.issues[0]?.message || "Invalid input",
      };
    }

    const { budgetId, ...updates } = validated.data;

    const response = await fetch(`${BACKEND_URL}/api/budgets/${budgetId}`, {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
        "X-User-ID": userId,
        "X-Request-ID": crypto.randomUUID(),
      },
      body: JSON.stringify(updates),
    });

    if (!response.ok) {
      const error = await response
        .json()
        .catch(() => ({ message: "Budget update failed" }));
      return { success: false, error: error.message || "Budget update failed" };
    }

    const budget: Budget = await response.json();
    revalidatePath("/dashboard");

    return { success: true, data: budget };
  } catch (error) {
    console.error("Update budget error:", error);
    return { success: false, error: "An unexpected error occurred" };
  }
}

/**
 * Delete Budget
 */
export async function deleteBudget(
  input: unknown
): Promise<ActionResult<{ success: boolean }>> {
  try {
    const session = await requireOrySession();
    const userId = session.identity?.id;
    if (!userId) {
      return { success: false, error: "Unauthorized" };
    }

    const validated = deleteBudgetSchema.safeParse(input);
    if (!validated.success) {
      return {
        success: false,
        error: validated.error.issues[0]?.message || "Invalid input",
      };
    }

    const response = await fetch(
      `${BACKEND_URL}/api/budgets/${validated.data.budgetId}`,
      {
        method: "DELETE",
        headers: {
          "X-User-ID": userId,
          "X-Request-ID": crypto.randomUUID(),
        },
      }
    );

    if (!response.ok) {
      const error = await response
        .json()
        .catch(() => ({ message: "Budget deletion failed" }));
      return { success: false, error: error.message || "Budget deletion failed" };
    }

    revalidatePath("/dashboard");
    return { success: true, data: { success: true } };
  } catch (error) {
    console.error("Delete budget error:", error);
    return { success: false, error: "An unexpected error occurred" };
  }
}

/**
 * Fetch User's Budgets
 */
export async function fetchUserBudgets(): Promise<ActionResult<Budget[]>> {
  try {
    const session = await requireOrySession();
    const userId = session.identity?.id;
    if (!userId) {
      return { success: false, error: "Unauthorized" };
    }

    const response = await fetch(`${BACKEND_URL}/api/budgets`, {
      method: "GET",
      headers: {
        "X-User-ID": userId,
      },
    });

    if (!response.ok) {
      return { success: false, error: "Failed to fetch budgets" };
    }

    const budgets: Budget[] = await response.json();
    return { success: true, data: budgets };
  } catch (error) {
    console.error("Fetch budgets error:", error);
    return { success: false, error: "An unexpected error occurred" };
  }
}

/**
 * Get Budget Summary
 */
export async function getBudgetSummary(): Promise<ActionResult<BudgetSummary>> {
  try {
    const session = await requireOrySession();
    const userId = session.identity?.id;
    if (!userId) {
      return { success: false, error: "Unauthorized" };
    }

    const response = await fetch(`${BACKEND_URL}/api/budgets/summary`, {
      method: "GET",
      headers: {
        "X-User-ID": userId,
      },
    });

    if (!response.ok) {
      return { success: false, error: "Failed to fetch budget summary" };
    }

    const summary: BudgetSummary = await response.json();
    return { success: true, data: summary };
  } catch (error) {
    console.error("Get budget summary error:", error);
    return { success: false, error: "An unexpected error occurred" };
  }
}

/**
 * Get Category Spending Analysis
 */
export async function getCategorySpending(
  input?: unknown
): Promise<ActionResult<CategorySpending[]>> {
  try {
    const session = await requireOrySession();
    const userId = session.identity?.id;
    if (!userId) {
      return { success: false, error: "Unauthorized" };
    }

    let queryParams = "";
    if (input) {
      const validated = getSpendingAnalysisSchema.safeParse(input);
      if (validated.success) {
        const params = new URLSearchParams({
          startDate: validated.data.startDate,
          endDate: validated.data.endDate,
        });
        if (validated.data.category) {
          params.append("category", validated.data.category);
        }
        queryParams = `?${params.toString()}`;
      }
    }

    const response = await fetch(
      `${BACKEND_URL}/api/analytics/category-spending${queryParams}`,
      {
        method: "GET",
        headers: {
          "X-User-ID": userId,
        },
      }
    );

    if (!response.ok) {
      return { success: false, error: "Failed to fetch category spending" };
    }

    const spending: CategorySpending[] = await response.json();
    return { success: true, data: spending };
  } catch (error) {
    console.error("Get category spending error:", error);
    return { success: false, error: "An unexpected error occurred" };
  }
}

/**
 * Get Spending Trends
 */
export async function getSpendingTrends(
  input?: unknown
): Promise<ActionResult<SpendingTrend[]>> {
  try {
    const session = await requireOrySession();
    const userId = session.identity?.id;
    if (!userId) {
      return { success: false, error: "Unauthorized" };
    }

    let queryParams = "";
    if (input) {
      const validated = getSpendingAnalysisSchema.safeParse(input);
      if (validated.success) {
        const params = new URLSearchParams({
          startDate: validated.data.startDate,
          endDate: validated.data.endDate,
        });
        if (validated.data.category) {
          params.append("category", validated.data.category);
        }
        queryParams = `?${params.toString()}`;
      }
    }

    const response = await fetch(
      `${BACKEND_URL}/api/analytics/spending-trends${queryParams}`,
      {
        method: "GET",
        headers: {
          "X-User-ID": userId,
        },
      }
    );

    if (!response.ok) {
      return { success: false, error: "Failed to fetch spending trends" };
    }

    const trends: SpendingTrend[] = await response.json();
    return { success: true, data: trends };
  } catch (error) {
    console.error("Get spending trends error:", error);
    return { success: false, error: "An unexpected error occurred" };
  }
}
