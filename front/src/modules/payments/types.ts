/**
 * src/modules/payments/types.ts
 *
 * Type definitions for payments domain
 * These types complement Zod schemas for full type safety
 */

import { z } from "zod";
import {
  transferSchema,
  transactionSchema,
  transactionListSchema,
} from "./validators";

export type Transfer = z.infer<typeof transferSchema>;
export type Transaction = z.infer<typeof transactionSchema>;
export type TransactionList = z.infer<typeof transactionListSchema>;

/**
 * Transaction status enum
 * Used for filtering and badge colors
 */
export enum TransactionStatus {
  Pending = "pending",
  Completed = "completed",
  Failed = "failed",
  Cancelled = "cancelled",
}

/**
 * Filter criteria for querying transactions
 */
export interface TransactionFilters {
  userId: string;
  page?: number;
  limit?: number;
  status?: TransactionStatus;
  startDate?: Date;
  endDate?: Date;
  minAmount?: number;
  maxAmount?: number;
}
