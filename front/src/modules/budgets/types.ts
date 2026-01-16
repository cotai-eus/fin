/**
 * Budgets Module - Type Definitions
 * Domain: Budget planning and spending tracking by category
 */

export enum BudgetCategory {
  FOOD = "food",
  TRANSPORT = "transport",
  ENTERTAINMENT = "entertainment",
  SHOPPING = "shopping",
  BILLS = "bills",
  HEALTH = "health",
  EDUCATION = "education",
  OTHER = "other",
}

export enum BudgetPeriod {
  WEEKLY = "weekly",
  MONTHLY = "monthly",
  YEARLY = "yearly",
}

export enum AlertThreshold {
  FIFTY_PERCENT = 50,
  SEVENTY_FIVE_PERCENT = 75,
  NINETY_PERCENT = 90,
}

export interface Budget {
  id: string;
  userId: string;
  category: BudgetCategory;
  period: BudgetPeriod;
  limit: number;
  currency: string;
  currentSpent: number;
  alertThreshold: AlertThreshold;
  alertsEnabled: boolean;
  startDate: string;
  endDate: string;
  createdAt: string;
  updatedAt: string;
}

export interface BudgetSummary {
  totalBudget: number;
  totalSpent: number;
  remainingBudget: number;
  percentageUsed: number;
  budgets: Budget[];
}

export interface CategorySpending {
  category: BudgetCategory;
  spent: number;
  budget?: number;
  percentageOfTotal: number;
  transactionCount: number;
}

export interface SpendingTrend {
  period: string; // ISO date string
  amount: number;
  category?: BudgetCategory;
}

export interface BudgetAlert {
  id: string;
  budgetId: string;
  userId: string;
  category: BudgetCategory;
  threshold: number;
  currentPercentage: number;
  message: string;
  isRead: boolean;
  createdAt: string;
}
