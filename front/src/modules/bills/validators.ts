/**
 * Bills Module - Zod Validation Schemas
 * Validates bill payment forms and barcode data
 */

import { z } from "zod";
import { BillType } from "./types";

/**
 * Brazilian Barcode Validation
 * Boleto bancário: 47 digits (with formatting) or 44 digits (numeric only)
 * Concessionária: 48 digits (with formatting) or 46-48 digits
 */
const barcodeRegex = /^[0-9]{44,48}$/;

/**
 * Formatted barcode (with spaces/dots)
 */
const formattedBarcodeRegex = /^[0-9\s.]+$/;

/**
 * Pay Bill Schema
 */
export const payBillSchema = z.object({
  barcode: z
    .string()
    .min(44, "Barcode must have at least 44 digits")
    .max(60, "Barcode is too long")
    .refine(
      (value) => {
        // Remove non-numeric characters
        const numeric = value.replace(/\D/g, "");
        return barcodeRegex.test(numeric);
      },
      { message: "Invalid barcode format" }
    ),
  amount: z
    .number()
    .positive("Amount must be positive")
    .max(1_000_000, "Amount cannot exceed R$ 1,000,000")
    .multipleOf(0.01)
    .optional(), // Some barcodes have embedded amounts
  scheduledFor: z.string().datetime().optional(),
});

export type PayBillInput = z.infer<typeof payBillSchema>;

/**
 * Scan Barcode Result Schema
 */
export const scanBarcodeSchema = z.object({
  barcode: z.string().min(44),
});

export type ScanBarcodeInput = z.infer<typeof scanBarcodeSchema>;

/**
 * Cancel Bill Payment Schema
 */
export const cancelBillPaymentSchema = z.object({
  billId: z.string().uuid("Invalid bill ID"),
  reason: z.string().min(3, "Reason is required").max(500),
});

export type CancelBillPaymentInput = z.infer<typeof cancelBillPaymentSchema>;

/**
 * Validate and Parse Barcode
 * Extracts information from barcode string
 */
export function parseBrazilianBarcode(barcode: string): {
  isValid: boolean;
  type: "boleto" | "concessionaria" | "unknown";
  amount?: number;
  dueDate?: string;
} {
  const numeric = barcode.replace(/\D/g, "");

  // Boleto bancário: 44 digits
  if (numeric.length === 44) {
    return {
      isValid: true,
      type: "boleto",
      // Amount extraction logic (simplified)
      // In real implementation, this would parse the barcode structure
    };
  }

  // Concessionária: 46-48 digits
  if (numeric.length >= 46 && numeric.length <= 48) {
    return {
      isValid: true,
      type: "concessionaria",
    };
  }

  return {
    isValid: false,
    type: "unknown",
  };
}

/**
 * Format barcode for display
 * Adds spaces for readability
 */
export function formatBarcode(barcode: string): string {
  const numeric = barcode.replace(/\D/g, "");

  if (numeric.length === 44) {
    // Boleto format: XXXXX.XXXXX XXXXX.XXXXXX XXXXX.XXXXXX X XXXXXXXXXXXXXX
    return numeric.replace(
      /^(\d{5})(\d{5})(\d{5})(\d{6})(\d{5})(\d{6})(\d{1})(\d{14})$/,
      "$1.$2 $3.$4 $5.$6 $7 $8"
    );
  }

  if (numeric.length >= 46 && numeric.length <= 48) {
    // Concessionária format (simplified)
    return numeric.match(/.{1,4}/g)?.join(" ") || numeric;
  }

  return numeric;
}
