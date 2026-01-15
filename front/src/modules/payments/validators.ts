/**
 * src/modules/payments/validators.ts
 *
 * Zod schemas for payments domain.
 * Validates transfer requests, transaction queries, and API responses.
 */

import { z } from "zod";

/**
 * Schema for transfer/payment execution
 * Used in Server Actions to validate user input before sending to backend
 */
export const transferSchema = z.object({
  fromUserId: z
    .string()
    .uuid("Invalid sender ID")
    .describe("UUID of the user initiating the transfer"),

  toUserId: z
    .string()
    .uuid("Invalid recipient ID")
    .describe("UUID of the recipient user"),

  recipientCPF: z
    .string()
    .regex(/^\d{3}\.\d{3}\.\d{3}-\d{2}$/, "Invalid CPF format")
    .optional()
    .describe("Optional: recipient CPF (for validation)"),

  amount: z
    .number()
    .positive("Amount must be greater than zero")
    .max(1_000_000, "Amount exceeds maximum transfer limit (R$ 1,000,000)")
    .multipleOf(0.01, "Amount must have maximum 2 decimal places")
    .describe("Transfer amount in BRL"),

  description: z
    .string()
    .max(500, "Description must be 500 characters or less")
    .optional()
    .describe("Optional description or reference for the transfer"),

  metadata: z
    .record(z.string(), z.string())
    .optional()
    .describe("Optional metadata for tracking"),

  scheduleDate: z
    .date()
    .optional()
    .refine(
      (date) => !date || date >= new Date(),
      "Scheduled date must be in the future"
    )
    .describe("Optional: schedule transfer for future date"),
});

export type Transfer = z.infer<typeof transferSchema>;

/**
 * Schema for transaction list query parameters
 * Used to validate search/filter parameters in Server Components
 */
export const transactionQuerySchema = z.object({
  userId: z.string().uuid("Invalid user ID"),
  page: z.number().int().positive().default(1),
  limit: z.number().int().min(1).max(100).default(20),
  status: z
    .enum(["pending", "completed", "failed", "cancelled"])
    .optional(),
  startDate: z.date().optional(),
  endDate: z.date().optional(),
  minAmount: z.number().nonnegative().optional(),
  maxAmount: z.number().nonnegative().optional(),
});

export type TransactionQuery = z.infer<typeof transactionQuerySchema>;

/**
 * Schema for transaction response from backend
 * Ensures backend API response conforms to expected shape
 */
export const transactionSchema = z.object({
  id: z.string().uuid("Invalid transaction ID"),
  fromUserId: z.string().uuid(),
  toUserId: z.string().uuid(),
  amount: z.number().positive(),
  status: z.enum(["pending", "completed", "failed", "cancelled"]),
  description: z.string().optional(),
  createdAt: z.string().datetime("Invalid date"),
  completedAt: z.string().datetime().optional(),
  metadata: z.record(z.string(), z.string()).optional(),
});

export type Transaction = z.infer<typeof transactionSchema>;

/**
 * Schema for paginated transactions response
 */
export const transactionListSchema = z.object({
  data: z.array(transactionSchema),
  page: z.number().int().positive(),
  limit: z.number().int().positive(),
  total: z.number().int().nonnegative(),
  totalPages: z.number().int().nonnegative(),
});

export type TransactionList = z.infer<typeof transactionListSchema>;

/**
 * Schema for transaction export request
 */
export const exportTransactionSchema = z.object({
  transactionId: z.string().uuid("Invalid transaction ID"),
  format: z.enum(["pdf", "csv"]).default("pdf"),
});

export type ExportTransaction = z.infer<typeof exportTransactionSchema>;
