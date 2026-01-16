/**
 * Transfers Module - Zod Validation Schemas
 * Validates transfer forms and API requests
 */

import { z } from "zod";
import { PIXKeyType, DepositMethod } from "./types";

/**
 * CPF validation (Brazilian tax ID)
 * Format: XXX.XXX.XXX-XX or XXXXXXXXXXX
 */
const cpfRegex = /^\d{3}\.?\d{3}\.?\d{3}-?\d{2}$/;

/**
 * CNPJ validation (Brazilian company tax ID)
 * Format: XX.XXX.XXX/XXXX-XX or XXXXXXXXXXXXXX
 */
const cnpjRegex = /^\d{2}\.?\d{3}\.?\d{3}\/?\d{4}-?\d{2}$/;

/**
 * Brazilian phone number
 * Format: (XX) XXXXX-XXXX or (XX) XXXX-XXXX
 */
const phoneRegex = /^\(?[1-9]{2}\)?\s?9?\d{4}-?\d{4}$/;

/**
 * PIX random key format (UUID-like)
 */
const pixRandomKeyRegex = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;

/**
 * Monetary amount validation
 * - Positive number
 * - Max 1 million BRL
 * - 2 decimal precision
 */
const amountSchema = z
  .number()
  .positive("Amount must be positive")
  .max(1_000_000, "Amount cannot exceed R$ 1,000,000")
  .multipleOf(0.01, "Amount must have at most 2 decimal places");

/**
 * PIX Transfer Schema
 */
export const pixTransferSchema = z.object({
  pixKey: z.string().min(1, "PIX key is required"),
  pixKeyType: z.nativeEnum(PIXKeyType),
  amount: amountSchema,
  description: z.string().max(500, "Description too long").optional(),
});

export type PIXTransferInput = z.infer<typeof pixTransferSchema>;

/**
 * TED Transfer Schema
 */
export const tedTransferSchema = z.object({
  recipientName: z.string().min(3, "Recipient name is required").max(100),
  recipientDocument: z
    .string()
    .regex(cpfRegex, "Invalid CPF format")
    .or(z.string().regex(cnpjRegex, "Invalid CNPJ format")),
  recipientBank: z.string().min(3, "Bank code is required"),
  recipientBranch: z.string().min(1, "Branch is required").max(10),
  recipientAccount: z.string().min(1, "Account is required").max(20),
  recipientAccountType: z.enum(["checking", "savings"]),
  amount: amountSchema,
  description: z.string().max(500).optional(),
  scheduledFor: z.string().datetime().optional(),
});

export type TEDTransferInput = z.infer<typeof tedTransferSchema>;

/**
 * P2P (Platform User to User) Transfer Schema
 */
export const p2pTransferSchema = z.object({
  recipientId: z.string().uuid("Invalid recipient ID"),
  amount: amountSchema,
  description: z.string().max(500).optional(),
});

export type P2PTransferInput = z.infer<typeof p2pTransferSchema>;

/**
 * Create Deposit Schema
 */
export const createDepositSchema = z.object({
  method: z.nativeEnum(DepositMethod),
  amount: amountSchema,
});

export type CreateDepositInput = z.infer<typeof createDepositSchema>;

/**
 * Create Payment Request Schema
 */
export const createPaymentRequestSchema = z.object({
  amount: amountSchema,
  description: z.string().min(3, "Description is required").max(200),
  expiresInHours: z.number().int().min(1).max(720).optional(), // Max 30 days
});

export type CreatePaymentRequestInput = z.infer<typeof createPaymentRequestSchema>;

/**
 * PIX Key Validator
 * Validates PIX key based on type
 */
export function validatePIXKey(key: string, type: PIXKeyType): boolean {
  switch (type) {
    case PIXKeyType.CPF:
      return cpfRegex.test(key);
    case PIXKeyType.CNPJ:
      return cnpjRegex.test(key);
    case PIXKeyType.EMAIL:
      return z.string().email().safeParse(key).success;
    case PIXKeyType.PHONE:
      return phoneRegex.test(key);
    case PIXKeyType.RANDOM:
      return pixRandomKeyRegex.test(key);
    default:
      return false;
  }
}

/**
 * Cancel Transfer Schema
 */
export const cancelTransferSchema = z.object({
  transferId: z.string().uuid("Invalid transfer ID"),
  reason: z.string().min(3, "Reason is required").max(500),
});

export type CancelTransferInput = z.infer<typeof cancelTransferSchema>;
