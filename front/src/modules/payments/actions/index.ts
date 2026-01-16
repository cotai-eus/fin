/**
 * src/modules/payments/actions/index.ts
 *
 * Server Actions for payments module.
 * All mutations (create, update, delete) are handled here with:
 * - Session verification (Zero-Trust)
 * - Input validation (Zod)
 * - Authorization checks
 * - Cache revalidation
 * - Error handling with Sentry
 */

"use server";

import { revalidatePath, revalidateTag } from "next/cache";
import { requireOrySession } from "@/core/ory/session";
import {
  transferSchema,
  exportTransactionSchema,
} from "../validators";
import * as Sentry from "@sentry/nextjs";

/**
 * Execute a transfer between two users.
 * 
 * Flow:
 * 1. Verify Ory session (Zero-Trust)
 * 2. Validate input with Zod
 * 3. Authorize: ensure user can only transfer from their own account
 * 4. Call backend API with session token
 * 5. Revalidate cache
 * 
 * @param formData - Form data from client
 * @returns Success response or error details
 */
export async function executeTransfer(formData: unknown) {
  try {
    // 1. Verify session (Zero-Trust principle)
    const session = await requireOrySession();
    const userId = session.identity?.id;

    if (!userId) {
      throw new Error("Invalid session: No user ID");
    }

    // 2. Validate input
    const validated = transferSchema.safeParse(formData);
    if (!validated.success) {
      return {
        error: "Validation failed",
        details: validated.error.flatten(),
      };
    }

    // 3. Authorization: User can only transfer from their own account
    if (validated.data.fromUserId !== userId) {
      Sentry.captureMessage(
        `Unauthorized transfer attempt by user ${userId}`,
        "warning"
      );

      throw new Error(
        "Unauthorized: Cannot transfer from another user's account"
      );
    }

    // 4. Call backend API
    const backendUrl = process.env.BACKEND_API_URL;
    if (!backendUrl) {
      throw new Error("Backend API URL not configured");
    }

    const response = await fetch(`${backendUrl}/api/transfers`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-Kratos-Authenticated-Identity-Id": userId,
        "X-Request-ID": crypto.randomUUID(),
      },
      body: JSON.stringify(validated.data),
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(
        errorData.message ||
        `Backend error: ${response.status} ${response.statusText}`
      );
    }

    const result = await response.json();

    // 5. Revalidate cache
    revalidatePath("/dashboard/payments", "page");
    revalidateTag("transactions", "max");

    return {
      success: true,
      data: result,
    };
  } catch (error) {
    // Log to Sentry for monitoring
    Sentry.captureException(error, {
      tags: { module: "payments", action: "executeTransfer" },
      extra: { formData },
    });

    // Return generic error to client (don't leak server details)
    console.error("[Payments] Transfer failed:", error);

    return {
      error:
        "Transfer could not be completed. Please check your information and try again.",
    };
  }
}

/**
 * Cancel a pending transfer.
 * Only the initiator can cancel their own transfer.
 * 
 * @param transferId - ID of transfer to cancel
 * @returns Success response or error
 */
export async function cancelTransfer(transferId: string) {
  try {
    const session = await requireOrySession();
    const userId = session.identity?.id;

    if (!userId) {
      throw new Error("Invalid session");
    }

    const backendUrl = process.env.BACKEND_API_URL;
    const response = await fetch(
      `${backendUrl}/api/transfers/${transferId}/cancel`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "X-Kratos-Authenticated-Identity-Id": userId,
          "X-Request-ID": crypto.randomUUID(),
        },
      }
    );

    if (!response.ok) {
      throw new Error("Could not cancel transfer");
    }

    revalidatePath("/dashboard/payments", "page");
    revalidateTag("transactions", "max");

    return { success: true };
  } catch (error) {
    Sentry.captureException(error, {
      tags: { action: "cancelTransfer" },
    });

    return { error: "Failed to cancel transfer" };
  }
}

/**
 * Export transaction as PDF or CSV.
 * User can only export their own transactions.
 * 
 * @param formData - Contains transactionId and format
 * @returns Download URL or error
 */
export async function exportTransaction(formData: unknown) {
  try {
    const session = await requireOrySession();
    const userId = session.identity?.id;

    if (!userId) {
      throw new Error("Invalid session");
    }

    const validated = exportTransactionSchema.safeParse(formData);
    if (!validated.success) {
      return { error: "Invalid request", details: validated.error.flatten() };
    }

    const backendUrl = process.env.BACKEND_API_URL;
    const response = await fetch(
      `${backendUrl}/api/transactions/${validated.data.transactionId}/export`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "X-Kratos-Authenticated-Identity-Id": userId,
          "X-Request-ID": crypto.randomUUID(),
        },
        body: JSON.stringify({
          format: validated.data.format,
        }),
      }
    );

    if (!response.ok) {
      throw new Error("Export failed");
    }

    const { downloadUrl } = await response.json();

    return { success: true, downloadUrl };
  } catch (error) {
    Sentry.captureException(error, {
      tags: { action: "exportTransaction" },
    });

    return { error: "Could not export transaction" };
  }
}

/**
 * Fetch user's transactions with pagination and filtering.
 * Server-side data fetching for RSCs.
 * Uses Next.js cache directives for ISR.
 */
export async function fetchUserTransactions(
  userId: string,
  page: number = 1,
  limit: number = 20
) {
  const backendUrl = process.env.BACKEND_API_URL;

  if (!backendUrl) {
    throw new Error("Backend API URL not configured");
  }

  try {
    const response = await fetch(
      `${backendUrl}/api/transactions?userId=${userId}&page=${page}&limit=${limit}`,
      {
        next: {
          revalidate: 60, // ISR: revalidate every 60s
          tags: ["transactions"], // Can revalidate manually
        },
        headers: {
          "X-User-ID": userId,
        },
      }
    );

    if (!response.ok) {
      throw new Error(
        `Failed to fetch transactions: ${response.statusText}`
      );
    }

    const data = await response.json();
    return data;
  } catch (error) {
    Sentry.captureException(error, {
      tags: { action: "fetchUserTransactions" },
      extra: { userId },
    });

    throw error;
  }
}
