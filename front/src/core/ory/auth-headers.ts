/**
 * src/core/ory/auth-headers.ts
 *
 * Authentication headers utilities for API requests
 */

import { Session } from "@ory/client";

/**
 * Get authentication headers for backend API requests
 * Ory uses cookie-based authentication, so we just need to ensure
 * the session is valid and let the cookies be sent automatically
 */
export function getAuthHeaders(session: Session): Record<string, string> {
  // Validate session is active
  if (!session.active || !session.identity?.id) {
    throw new Error("Invalid or inactive session");
  }

  // Return headers for backend API requests
  // The session cookie (ory_kratos_session) will be sent automatically by the browser
  return {
    "Content-Type": "application/json",
    "X-User-Id": session.identity.id,
  };
}

/**
 * Get user ID from session
 */
export function getUserId(session: Session): string {
  if (!session.identity?.id) {
    throw new Error("No user identity found in session");
  }
  return session.identity.id;
}

/**
 * Get user email from session
 */
export function getUserEmail(session: Session): string {
  const email = session.identity?.traits?.email;
  if (!email) {
    throw new Error("No email found in session");
  }
  return email as string;
}
