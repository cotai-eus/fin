/**
 * Bills Module - Server Actions
 * Handles bill payments
 */

"use server";

import { revalidatePath } from "next/cache";
import { requireOrySession } from "@/core/ory/session";
import {
  payBillSchema,
  cancelBillPaymentSchema,
  type PayBillInput,
  type CancelBillPaymentInput,
  parseBrazilianBarcode,
} from "../validators";
import type { Bill, BarcodeData } from "../types";

const BACKEND_URL = process.env.BACKEND_API_URL || "http://localhost:8080";

type ActionResult<T> =
  | { success: true; data: T }
  | { success: false; error: string };

/**
 * Validate Barcode
 * Checks barcode and retrieves bill information
 */
export async function validateBarcode(
  barcode: string
): Promise<ActionResult<BarcodeData>> {
  try {
    const session = await requireOrySession();
    const userId = session.identity?.id;
    if (!userId) {
      return { success: false, error: "Unauthorized" };
    }

    // Client-side validation
    const parsed = parseBrazilianBarcode(barcode);
    if (!parsed.isValid) {
      return { success: false, error: "Invalid barcode format" };
    }

    // Call backend to get bill details
    const response = await fetch(`${BACKEND_URL}/api/bills/validate`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-User-ID": userId,
        "X-Request-ID": crypto.randomUUID(),
      },
      body: JSON.stringify({ barcode: barcode.replace(/\D/g, "") }),
    });

    if (!response.ok) {
      const error = await response
        .json()
        .catch(() => ({ message: "Barcode validation failed" }));
      return { success: false, error: error.message || "Barcode validation failed" };
    }

    const barcodeData: BarcodeData = await response.json();
    return { success: true, data: barcodeData };
  } catch (error) {
    console.error("Validate barcode error:", error);
    return { success: false, error: "An unexpected error occurred" };
  }
}

/**
 * Pay Bill
 * Process bill payment
 */
export async function payBill(
  input: unknown
): Promise<ActionResult<Bill>> {
  try {
    const session = await requireOrySession();
    const userId = session.identity?.id;
    if (!userId) {
      return { success: false, error: "Unauthorized" };
    }

    const validated = payBillSchema.safeParse(input);
    if (!validated.success) {
      return {
        success: false,
        error: validated.error.issues[0]?.message || "Invalid input",
      };
    }

    // Normalize barcode (remove non-numeric)
    const normalizedBarcode = validated.data.barcode.replace(/\D/g, "");

    const response = await fetch(`${BACKEND_URL}/api/bills/pay`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-User-ID": userId,
        "X-Request-ID": crypto.randomUUID(),
      },
      body: JSON.stringify({
        ...validated.data,
        barcode: normalizedBarcode,
      }),
    });

    if (!response.ok) {
      const error = await response
        .json()
        .catch(() => ({ message: "Bill payment failed" }));
      return { success: false, error: error.message || "Bill payment failed" };
    }

    const bill: Bill = await response.json();
    revalidatePath("/dashboard/payments");
    revalidatePath("/dashboard");

    return { success: true, data: bill };
  } catch (error) {
    console.error("Pay bill error:", error);
    return { success: false, error: "An unexpected error occurred" };
  }
}

/**
 * Cancel Bill Payment
 * Only pending or scheduled bills can be cancelled
 */
export async function cancelBillPayment(
  input: unknown
): Promise<ActionResult<Bill>> {
  try {
    const session = await requireOrySession();
    const userId = session.identity?.id;
    if (!userId) {
      return { success: false, error: "Unauthorized" };
    }

    const validated = cancelBillPaymentSchema.safeParse(input);
    if (!validated.success) {
      return {
        success: false,
        error: validated.error.issues[0]?.message || "Invalid input",
      };
    }

    const response = await fetch(
      `${BACKEND_URL}/api/bills/${validated.data.billId}/cancel`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "X-User-ID": userId,
          "X-Request-ID": crypto.randomUUID(),
        },
        body: JSON.stringify({ reason: validated.data.reason }),
      }
    );

    if (!response.ok) {
      const error = await response
        .json()
        .catch(() => ({ message: "Cancellation failed" }));
      return { success: false, error: error.message || "Cancellation failed" };
    }

    const bill: Bill = await response.json();
    revalidatePath("/dashboard/payments");

    return { success: true, data: bill };
  } catch (error) {
    console.error("Cancel bill payment error:", error);
    return { success: false, error: "An unexpected error occurred" };
  }
}

/**
 * Fetch User's Bill Payments
 */
export async function fetchUserBills(
  page: number = 1,
  limit: number = 20
): Promise<ActionResult<{ data: Bill[]; total: number; page: number }>> {
  try {
    const session = await requireOrySession();
    const userId = session.identity?.id;
    if (!userId) {
      return { success: false, error: "Unauthorized" };
    }

    const response = await fetch(
      `${BACKEND_URL}/api/bills?page=${page}&limit=${limit}`,
      {
        method: "GET",
        headers: {
          "X-User-ID": userId,
        },
      }
    );

    if (!response.ok) {
      return { success: false, error: "Failed to fetch bills" };
    }

    const data = await response.json();
    return { success: true, data };
  } catch (error) {
    console.error("Fetch bills error:", error);
    return { success: false, error: "An unexpected error occurred" };
  }
}
