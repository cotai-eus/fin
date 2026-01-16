/**
 * Transfers Module - Server Actions
 * Handles money transfers, deposits, and payment requests
 */

"use server";

import { revalidatePath } from "next/cache";
import { requireOrySession, getUserId } from "@/core/ory/session";
import {
  pixTransferSchema,
  tedTransferSchema,
  p2pTransferSchema,
  createDepositSchema,
  createPaymentRequestSchema,
  cancelTransferSchema,
  type PIXTransferInput,
  type TEDTransferInput,
  type P2PTransferInput,
  type CreateDepositInput,
  type CreatePaymentRequestInput,
  type CancelTransferInput,
} from "../validators";
import type { Transfer, Deposit, PaymentRequest } from "../types";

const BACKEND_URL = process.env.BACKEND_API_URL || "http://localhost:8080";

type ActionResult<T> = 
  | { success: true; data: T }
  | { success: false; error: string };

/**
 * Execute PIX Transfer
 * Zero-Trust Flow: Session → Validate → Authorize → Execute → Revalidate
 */
export async function executePIXTransfer(
  input: unknown
): Promise<ActionResult<Transfer>> {
  try {
    // 1. Verify session
    const session = await requireOrySession();
    const userId = session.identity?.id;
    if (!userId) {
      return { success: false, error: "Unauthorized: No user ID" };
    }

    // 2. Validate input
    const validated = pixTransferSchema.safeParse(input);
    if (!validated.success) {
      return { 
        success: false, 
        error: validated.error.issues[0]?.message || "Invalid input" 
      };
    }

    // 3. Call backend API (convert BRL to cents)
    const response = await fetch(`${BACKEND_URL}/api/transfers/pix`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-Kratos-Authenticated-Identity-Id": userId,
        "X-Request-ID": crypto.randomUUID(),
      },
      body: JSON.stringify({
        pix_key: validated.data.pixKey,
        pix_key_type: validated.data.pixKeyType,
        amount_cents: Math.round(validated.data.amount * 100),
        description: validated.data.description,
      }),
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ message: "Transfer failed" }));
      return { success: false, error: error.message || "Transfer failed" };
    }

    const rawTransfer = await response.json();
    
    // Convert backend response (cents) to frontend format (BRL)
    const transfer: Transfer = {
      id: rawTransfer.id,
      userId: rawTransfer.user_id,
      type: rawTransfer.type,
      status: rawTransfer.status,
      amount: rawTransfer.amount_cents / 100,
      currency: rawTransfer.currency,
      pixKey: rawTransfer.pix_key,
      pixKeyType: rawTransfer.pix_key_type,
      fee: rawTransfer.fee_cents / 100,
      completedAt: rawTransfer.completed_at,
      createdAt: rawTransfer.created_at,
      updatedAt: rawTransfer.updated_at,
    };

    // 4. Revalidate cache
    revalidatePath("/dashboard/payments");
    revalidatePath("/dashboard");

    return { success: true, data: transfer };
  } catch (error) {
    console.error("PIX transfer error:", error);
    return { success: false, error: "An unexpected error occurred" };
  }
}

/**
 * Execute TED Transfer
 */
export async function executeTEDTransfer(
  input: unknown
): Promise<ActionResult<Transfer>> {
  try {
    const session = await requireOrySession();
    const userId = session.identity?.id;
    if (!userId) {
      return { success: false, error: "Unauthorized" };
    }

    const validated = tedTransferSchema.safeParse(input);
    if (!validated.success) {
      return { 
        success: false, 
        error: validated.error.issues[0]?.message || "Invalid input" 
      };
    }

    const response = await fetch(`${BACKEND_URL}/api/transfers/ted`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-Kratos-Authenticated-Identity-Id": userId,
        "X-Request-ID": crypto.randomUUID(),
      },
      body: JSON.stringify({
        recipient_name: validated.data.recipientName,
        recipient_document: validated.data.recipientDocument,
        recipient_bank: validated.data.recipientBank,
        recipient_branch: validated.data.recipientBranch,
        recipient_account: validated.data.recipientAccount,
        recipient_account_type: validated.data.recipientAccountType,
        amount_cents: Math.round(validated.data.amount * 100),
        description: validated.data.description,
      }),
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ message: "Transfer failed" }));
      return { success: false, error: error.message || "Transfer failed" };
    }

    const rawTransfer = await response.json();
    const transfer: Transfer = {
      id: rawTransfer.id,
      userId: rawTransfer.user_id,
      type: rawTransfer.type,
      status: rawTransfer.status,
      amount: rawTransfer.amount_cents / 100,
      currency: rawTransfer.currency,
      recipientName: rawTransfer.recipient_name,
      recipientDocument: rawTransfer.recipient_document,
      recipientBank: rawTransfer.recipient_bank,
      recipientBranch: rawTransfer.recipient_branch,
      recipientAccount: rawTransfer.recipient_account,
      recipientAccountType: rawTransfer.recipient_account_type,
      fee: rawTransfer.fee_cents / 100,
      completedAt: rawTransfer.completed_at,
      createdAt: rawTransfer.created_at,
      updatedAt: rawTransfer.updated_at,
    };
    revalidatePath("/dashboard/payments");
    revalidatePath("/dashboard");

    return { success: true, data: transfer };
  } catch (error) {
    console.error("TED transfer error:", error);
    return { success: false, error: "An unexpected error occurred" };
  }
}

/**
 * Execute P2P Transfer (User to User on platform)
 */
export async function executeP2PTransfer(
  input: unknown
): Promise<ActionResult<Transfer>> {
  try {
    const session = await requireOrySession();
    const userId = session.identity?.id;
    if (!userId) {
      return { success: false, error: "Unauthorized" };
    }

    const validated = p2pTransferSchema.safeParse(input);
    if (!validated.success) {
      return { 
        success: false, 
        error: validated.error.issues[0]?.message || "Invalid input" 
      };
    }

    // Prevent self-transfer
    if (validated.data.recipientId === userId) {
      return { success: false, error: "Cannot transfer to yourself" };
    }

    const response = await fetch(`${BACKEND_URL}/api/transfers/p2p`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-Kratos-Authenticated-Identity-Id": userId,
        "X-Request-ID": crypto.randomUUID(),
      },
      body: JSON.stringify({
        recipient_user_id: validated.data.recipientId,
        amount_cents: Math.round(validated.data.amount * 100),
        description: validated.data.description,
      }),
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ message: "Transfer failed" }));
      return { success: false, error: error.message || "Transfer failed" };
    }

    const rawTransfer = await response.json();
    const transfer: Transfer = {
      id: rawTransfer.id,
      userId: rawTransfer.user_id,
      type: rawTransfer.type,
      status: rawTransfer.status,
      amount: rawTransfer.amount_cents / 100,
      currency: rawTransfer.currency,
      fee: rawTransfer.fee_cents / 100,
      completedAt: rawTransfer.completed_at,
      createdAt: rawTransfer.created_at,
      updatedAt: rawTransfer.updated_at,
    };
    revalidatePath("/dashboard/payments");
    revalidatePath("/dashboard");

    return { success: true, data: transfer };
  } catch (error) {
    console.error("P2P transfer error:", error);
    return { success: false, error: "An unexpected error occurred" };
  }
}

/**
 * Create Deposit
 * Generates PIX QR code, boleto, or bank transfer instructions
 */
export async function createDeposit(
  input: unknown
): Promise<ActionResult<Deposit>> {
  try {
    const session = await requireOrySession();
    const userId = session.identity?.id;
    if (!userId) {
      return { success: false, error: "Unauthorized" };
    }

    const validated = createDepositSchema.safeParse(input);
    if (!validated.success) {
      return { 
        success: false, 
        error: validated.error.issues[0]?.message || "Invalid input" 
      };
    }

    const response = await fetch(`${BACKEND_URL}/api/deposits`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-User-ID": userId,
        "X-Request-ID": crypto.randomUUID(),
      },
      body: JSON.stringify(validated.data),
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ message: "Deposit creation failed" }));
      return { success: false, error: error.message || "Deposit creation failed" };
    }

    const deposit: Deposit = await response.json();
    revalidatePath("/dashboard");

    return { success: true, data: deposit };
  } catch (error) {
    console.error("Create deposit error:", error);
    return { success: false, error: "An unexpected error occurred" };
  }
}

/**
 * Create Payment Request
 * Generates a payment link/QR code for requesting money from others
 */
export async function createPaymentRequest(
  input: unknown
): Promise<ActionResult<PaymentRequest>> {
  try {
    const session = await requireOrySession();
    const userId = session.identity?.id;
    if (!userId) {
      return { success: false, error: "Unauthorized" };
    }

    const validated = createPaymentRequestSchema.safeParse(input);
    if (!validated.success) {
      return { 
        success: false, 
        error: validated.error.issues[0]?.message || "Invalid input" 
      };
    }

    const response = await fetch(`${BACKEND_URL}/api/payment-requests`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-User-ID": userId,
        "X-Request-ID": crypto.randomUUID(),
      },
      body: JSON.stringify(validated.data),
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ message: "Payment request creation failed" }));
      return { success: false, error: error.message || "Payment request creation failed" };
    }

    const paymentRequest: PaymentRequest = await response.json();
    revalidatePath("/dashboard/payments");

    return { success: true, data: paymentRequest };
  } catch (error) {
    console.error("Create payment request error:", error);
    return { success: false, error: "An unexpected error occurred" };
  }
}

/**
 * Cancel Transfer
 * Only pending or scheduled transfers can be cancelled
 */
export async function cancelTransfer(
  input: unknown
): Promise<ActionResult<Transfer>> {
  try {
    const session = await requireOrySession();
    const userId = session.identity?.id;
    if (!userId) {
      return { success: false, error: "Unauthorized" };
    }

    const validated = cancelTransferSchema.safeParse(input);
    if (!validated.success) {
      return { 
        success: false, 
        error: validated.error.issues[0]?.message || "Invalid input" 
      };
    }

    const response = await fetch(
      `${BACKEND_URL}/api/transfers/${validated.data.transferId}/cancel`,
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
      const error = await response.json().catch(() => ({ message: "Cancel failed" }));
      return { success: false, error: error.message || "Cancel failed" };
    }

    const transfer: Transfer = await response.json();
    revalidatePath("/dashboard/payments");

    return { success: true, data: transfer };
  } catch (error) {
    console.error("Cancel transfer error:", error);
    return { success: false, error: "An unexpected error occurred" };
  }
}

/**
 * Fetch user's transfers with pagination
 */
export async function fetchUserTransfers(
  page: number = 1,
  limit: number = 20
): Promise<ActionResult<{ data: Transfer[]; pagination: { page: number; limit: number; total: number; total_pages: number; has_more: boolean } }>> {
  try {
    const session = await requireOrySession();
    const userId = session.identity?.id;
    if (!userId) {
      return { success: false, error: "Unauthorized" };
    }

    const response = await fetch(
      `${BACKEND_URL}/api/transfers?page=${page}&limit=${limit}`,
      {
        method: "GET",
        headers: {
          "X-Kratos-Authenticated-Identity-Id": userId,
        },
      }
    );

    if (!response.ok) {
      return { success: false, error: "Failed to fetch transfers" };
    }

    const rawData = await response.json();
    
    // Convert backend response (cents) to frontend format (BRL)
    const transfers: Transfer[] = rawData.data.map((t: any) => ({
      id: t.id,
      userId: t.user_id,
      type: t.type,
      status: t.status,
      amount: t.amount_cents / 100,
      currency: t.currency,
      pixKey: t.pix_key,
      pixKeyType: t.pix_key_type,
      recipientName: t.recipient_name,
      recipientDocument: t.recipient_document,
      recipientBank: t.recipient_bank,
      recipientBranch: t.recipient_branch,
      recipientAccount: t.recipient_account,
      recipientAccountType: t.recipient_account_type,
      fee: t.fee_cents / 100,
      completedAt: t.completed_at,
      createdAt: t.created_at,
      updatedAt: t.updated_at,
    }));

    return { 
      success: true, 
      data: { 
        data: transfers, 
        pagination: rawData.pagination 
      } 
    };
  } catch (error) {
    console.error("Fetch transfers error:", error);
    return { success: false, error: "An unexpected error occurred" };
  }
}
