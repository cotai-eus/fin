/**
 * Cards Module - Zod Validation Schemas
 * Validates card operations and settings
 */

import { z } from "zod";
import { CardType } from "./types";

/**
 * Create Virtual Card Schema
 */
export const createVirtualCardSchema = z.object({
  holderName: z
    .string()
    .min(3, "Name must be at least 3 characters")
    .max(100, "Name too long"),
  dailyLimit: z
    .number()
    .positive("Daily limit must be positive")
    .max(50_000, "Daily limit cannot exceed R$ 50,000"),
  monthlyLimit: z
    .number()
    .positive("Monthly limit must be positive")
    .max(500_000, "Monthly limit cannot exceed R$ 500,000"),
  isInternational: z.boolean().default(false),
});

export type CreateVirtualCardInput = z.infer<typeof createVirtualCardSchema>;

/**
 * Update Card Limits Schema
 */
export const updateCardLimitsSchema = z.object({
  cardId: z.string().uuid("Invalid card ID"),
  dailyLimit: z
    .number()
    .positive("Daily limit must be positive")
    .max(50_000, "Daily limit cannot exceed R$ 50,000")
    .optional(),
  monthlyLimit: z
    .number()
    .positive("Monthly limit must be positive")
    .max(500_000, "Monthly limit cannot exceed R$ 500,000")
    .optional(),
  perTransactionLimit: z
    .number()
    .positive("Per transaction limit must be positive")
    .max(50_000, "Per transaction limit cannot exceed R$ 50,000")
    .optional(),
});

export type UpdateCardLimitsInput = z.infer<typeof updateCardLimitsSchema>;

/**
 * Block/Unblock Card Schema
 */
export const toggleCardStatusSchema = z.object({
  cardId: z.string().uuid("Invalid card ID"),
  action: z.enum(["block", "unblock"]),
  reason: z.string().min(3, "Reason is required").max(500).optional(),
});

export type ToggleCardStatusInput = z.infer<typeof toggleCardStatusSchema>;

/**
 * Report Card Lost/Stolen Schema
 */
export const reportCardSchema = z.object({
  cardId: z.string().uuid("Invalid card ID"),
  reportType: z.enum(["lost", "stolen"]),
  description: z.string().min(10, "Please provide more details").max(1000),
  requestReplacement: z.boolean().default(true),
});

export type ReportCardInput = z.infer<typeof reportCardSchema>;

/**
 * Change Card PIN Schema
 */
export const changeCardPINSchema = z.object({
  cardId: z.string().uuid("Invalid card ID"),
  currentPIN: z
    .string()
    .length(4, "PIN must be 4 digits")
    .regex(/^\d{4}$/, "PIN must contain only numbers"),
  newPIN: z
    .string()
    .length(4, "PIN must be 4 digits")
    .regex(/^\d{4}$/, "PIN must contain only numbers"),
  confirmPIN: z.string(),
}).refine((data) => data.newPIN === data.confirmPIN, {
  message: "PINs do not match",
  path: ["confirmPIN"],
});

export type ChangeCardPINInput = z.infer<typeof changeCardPINSchema>;

/**
 * Update Card Security Settings Schema
 */
export const updateSecuritySettingsSchema = z.object({
  cardId: z.string().uuid("Invalid card ID"),
  blockInternational: z.boolean().optional(),
  blockOnline: z.boolean().optional(),
  blockContactless: z.boolean().optional(),
  blockAtm: z.boolean().optional(),
  requireNotification: z.boolean().optional(),
});

export type UpdateSecuritySettingsInput = z.infer<
  typeof updateSecuritySettingsSchema
>;

/**
 * Cancel Card Schema
 */
export const cancelCardSchema = z.object({
  cardId: z.string().uuid("Invalid card ID"),
  reason: z.string().min(3, "Reason is required").max(500),
});

export type CancelCardInput = z.infer<typeof cancelCardSchema>;
