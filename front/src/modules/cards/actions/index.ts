/**
 * Cards Module - Server Actions
 * Handles card management operations
 */

"use server";

import { revalidatePath } from "next/cache";
import { requireOrySession } from "@/core/ory/session";
import {
  createVirtualCardSchema,
  updateCardLimitsSchema,
  toggleCardStatusSchema,
  reportCardSchema,
  changeCardPINSchema,
  updateSecuritySettingsSchema,
  cancelCardSchema,
  type CreateVirtualCardInput,
  type UpdateCardLimitsInput,
  type ToggleCardStatusInput,
  type ReportCardInput,
  type ChangeCardPINInput,
  type UpdateSecuritySettingsInput,
  type CancelCardInput,
} from "../validators";
import type { Card, CardDetails, CardTransaction, CardSecuritySettings } from "../types";

const BACKEND_URL = process.env.BACKEND_API_URL || "http://localhost:8080";

type ActionResult<T> =
  | { success: true; data: T }
  | { success: false; error: string };

/**
 * Fetch User's Cards
 */
export async function fetchUserCards(): Promise<ActionResult<Card[]>> {
  try {
    const session = await requireOrySession();
    const userId = session.identity?.id;
    if (!userId) {
      return { success: false, error: "Unauthorized" };
    }

    const response = await fetch(`${BACKEND_URL}/api/cards`, {
      method: "GET",
      headers: {
        "X-User-ID": userId,
      },
    });

    if (!response.ok) {
      return { success: false, error: "Failed to fetch cards" };
    }

    const cards: Card[] = await response.json();
    return { success: true, data: cards };
  } catch (error) {
    console.error("Fetch cards error:", error);
    return { success: false, error: "An unexpected error occurred" };
  }
}

/**
 * Get Card Details (sensitive data - number, CVV)
 * Requires additional verification for security
 */
export async function getCardDetails(
  cardId: string
): Promise<ActionResult<CardDetails>> {
  try {
    const session = await requireOrySession();
    const userId = session.identity?.id;
    if (!userId) {
      return { success: false, error: "Unauthorized" };
    }

    const response = await fetch(`${BACKEND_URL}/api/cards/${cardId}/details`, {
      method: "GET",
      headers: {
        "X-User-ID": userId,
        "X-Request-ID": crypto.randomUUID(),
      },
    });

    if (!response.ok) {
      return { success: false, error: "Failed to fetch card details" };
    }

    const details: CardDetails = await response.json();
    return { success: true, data: details };
  } catch (error) {
    console.error("Get card details error:", error);
    return { success: false, error: "An unexpected error occurred" };
  }
}

/**
 * Create Virtual Card
 */
export async function createVirtualCard(
  input: unknown
): Promise<ActionResult<Card>> {
  try {
    const session = await requireOrySession();
    const userId = session.identity?.id;
    if (!userId) {
      return { success: false, error: "Unauthorized" };
    }

    const validated = createVirtualCardSchema.safeParse(input);
    if (!validated.success) {
      return {
        success: false,
        error: validated.error.issues[0]?.message || "Invalid input",
      };
    }

    const response = await fetch(`${BACKEND_URL}/api/cards/virtual`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-User-ID": userId,
        "X-Request-ID": crypto.randomUUID(),
      },
      body: JSON.stringify(validated.data),
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ message: "Card creation failed" }));
      return { success: false, error: error.message || "Card creation failed" };
    }

    const card: Card = await response.json();
    revalidatePath("/dashboard/cards");

    return { success: true, data: card };
  } catch (error) {
    console.error("Create virtual card error:", error);
    return { success: false, error: "An unexpected error occurred" };
  }
}

/**
 * Update Card Limits
 */
export async function updateCardLimits(
  input: unknown
): Promise<ActionResult<Card>> {
  try {
    const session = await requireOrySession();
    const userId = session.identity?.id;
    if (!userId) {
      return { success: false, error: "Unauthorized" };
    }

    const validated = updateCardLimitsSchema.safeParse(input);
    if (!validated.success) {
      return {
        success: false,
        error: validated.error.issues[0]?.message || "Invalid input",
      };
    }

    const { cardId, ...limits } = validated.data;

    const response = await fetch(`${BACKEND_URL}/api/cards/${cardId}/limits`, {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
        "X-User-ID": userId,
        "X-Request-ID": crypto.randomUUID(),
      },
      body: JSON.stringify(limits),
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ message: "Update failed" }));
      return { success: false, error: error.message || "Update failed" };
    }

    const card: Card = await response.json();
    revalidatePath("/dashboard/cards");

    return { success: true, data: card };
  } catch (error) {
    console.error("Update card limits error:", error);
    return { success: false, error: "An unexpected error occurred" };
  }
}

/**
 * Block/Unblock Card
 */
export async function toggleCardStatus(
  input: unknown
): Promise<ActionResult<Card>> {
  try {
    const session = await requireOrySession();
    const userId = session.identity?.id;
    if (!userId) {
      return { success: false, error: "Unauthorized" };
    }

    const validated = toggleCardStatusSchema.safeParse(input);
    if (!validated.success) {
      return {
        success: false,
        error: validated.error.issues[0]?.message || "Invalid input",
      };
    }

    const { cardId, action, reason } = validated.data;

    const response = await fetch(
      `${BACKEND_URL}/api/cards/${cardId}/${action}`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "X-User-ID": userId,
          "X-Request-ID": crypto.randomUUID(),
        },
        body: JSON.stringify({ reason }),
      }
    );

    if (!response.ok) {
      const error = await response.json().catch(() => ({ message: "Operation failed" }));
      return { success: false, error: error.message || "Operation failed" };
    }

    const card: Card = await response.json();
    revalidatePath("/dashboard/cards");

    return { success: true, data: card };
  } catch (error) {
    console.error("Toggle card status error:", error);
    return { success: false, error: "An unexpected error occurred" };
  }
}

/**
 * Report Card Lost/Stolen
 */
export async function reportCard(
  input: unknown
): Promise<ActionResult<Card>> {
  try {
    const session = await requireOrySession();
    const userId = session.identity?.id;
    if (!userId) {
      return { success: false, error: "Unauthorized" };
    }

    const validated = reportCardSchema.safeParse(input);
    if (!validated.success) {
      return {
        success: false,
        error: validated.error.issues[0]?.message || "Invalid input",
      };
    }

    const { cardId, ...reportData } = validated.data;

    const response = await fetch(`${BACKEND_URL}/api/cards/${cardId}/report`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-User-ID": userId,
        "X-Request-ID": crypto.randomUUID(),
      },
      body: JSON.stringify(reportData),
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ message: "Report failed" }));
      return { success: false, error: error.message || "Report failed" };
    }

    const card: Card = await response.json();
    revalidatePath("/dashboard/cards");

    return { success: true, data: card };
  } catch (error) {
    console.error("Report card error:", error);
    return { success: false, error: "An unexpected error occurred" };
  }
}

/**
 * Change Card PIN
 */
export async function changeCardPIN(
  input: unknown
): Promise<ActionResult<{ success: boolean }>> {
  try {
    const session = await requireOrySession();
    const userId = session.identity?.id;
    if (!userId) {
      return { success: false, error: "Unauthorized" };
    }

    const validated = changeCardPINSchema.safeParse(input);
    if (!validated.success) {
      return {
        success: false,
        error: validated.error.issues[0]?.message || "Invalid input",
      };
    }

    const { cardId, currentPIN, newPIN } = validated.data;

    const response = await fetch(`${BACKEND_URL}/api/cards/${cardId}/pin`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-User-ID": userId,
        "X-Request-ID": crypto.randomUUID(),
      },
      body: JSON.stringify({ currentPIN, newPIN }),
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ message: "PIN change failed" }));
      return { success: false, error: error.message || "PIN change failed" };
    }

    return { success: true, data: { success: true } };
  } catch (error) {
    console.error("Change card PIN error:", error);
    return { success: false, error: "An unexpected error occurred" };
  }
}

/**
 * Update Card Security Settings
 */
export async function updateSecuritySettings(
  input: unknown
): Promise<ActionResult<CardSecuritySettings>> {
  try {
    const session = await requireOrySession();
    const userId = session.identity?.id;
    if (!userId) {
      return { success: false, error: "Unauthorized" };
    }

    const validated = updateSecuritySettingsSchema.safeParse(input);
    if (!validated.success) {
      return {
        success: false,
        error: validated.error.issues[0]?.message || "Invalid input",
      };
    }

    const { cardId, ...settings } = validated.data;

    const response = await fetch(`${BACKEND_URL}/api/cards/${cardId}/security`, {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
        "X-User-ID": userId,
        "X-Request-ID": crypto.randomUUID(),
      },
      body: JSON.stringify(settings),
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ message: "Update failed" }));
      return { success: false, error: error.message || "Update failed" };
    }

    const securitySettings: CardSecuritySettings = await response.json();
    revalidatePath("/dashboard/cards");

    return { success: true, data: securitySettings };
  } catch (error) {
    console.error("Update security settings error:", error);
    return { success: false, error: "An unexpected error occurred" };
  }
}

/**
 * Cancel Card
 */
export async function cancelCard(
  input: unknown
): Promise<ActionResult<Card>> {
  try {
    const session = await requireOrySession();
    const userId = session.identity?.id;
    if (!userId) {
      return { success: false, error: "Unauthorized" };
    }

    const validated = cancelCardSchema.safeParse(input);
    if (!validated.success) {
      return {
        success: false,
        error: validated.error.issues[0]?.message || "Invalid input",
      };
    }

    const { cardId, reason } = validated.data;

    const response = await fetch(`${BACKEND_URL}/api/cards/${cardId}/cancel`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-User-ID": userId,
        "X-Request-ID": crypto.randomUUID(),
      },
      body: JSON.stringify({ reason }),
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ message: "Cancellation failed" }));
      return { success: false, error: error.message || "Cancellation failed" };
    }

    const card: Card = await response.json();
    revalidatePath("/dashboard/cards");

    return { success: true, data: card };
  } catch (error) {
    console.error("Cancel card error:", error);
    return { success: false, error: "An unexpected error occurred" };
  }
}

/**
 * Fetch Card Transactions
 */
export async function fetchCardTransactions(
  cardId: string,
  page: number = 1,
  limit: number = 20
): Promise<ActionResult<{ data: CardTransaction[]; total: number }>> {
  try {
    const session = await requireOrySession();
    const userId = session.identity?.id;
    if (!userId) {
      return { success: false, error: "Unauthorized" };
    }

    const response = await fetch(
      `${BACKEND_URL}/api/cards/${cardId}/transactions?page=${page}&limit=${limit}`,
      {
        method: "GET",
        headers: {
          "X-User-ID": userId,
        },
      }
    );

    if (!response.ok) {
      return { success: false, error: "Failed to fetch transactions" };
    }

    const data = await response.json();
    return { success: true, data };
  } catch (error) {
    console.error("Fetch card transactions error:", error);
    return { success: false, error: "An unexpected error occurred" };
  }
}
