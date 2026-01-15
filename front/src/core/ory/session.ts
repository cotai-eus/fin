/**
 * src/core/ory/session.ts
 *
 * Session verification utilities for Ory Kratos integration.
 * Handles server-side session retrieval and validation for RSCs and Server Actions.
 */

import { FrontendApi, Configuration } from "@ory/client";
import { cookies } from "next/headers";

/**
 * Initialize Ory SDK with server configuration
 */
function getOryClient() {
  const sdkUrl = process.env.ORY_SDK_URL || "http://kratos:4433";

  return new FrontendApi(
    new Configuration({
      basePath: sdkUrl,
      baseOptions: {
        withCredentials: true,
        headers: {
          "Content-Type": "application/json",
        },
      },
    })
  );
}

/**
 * Get Ory session from cookies.
 * Called from Server Components and Server Actions.
 *
 * @returns Ory session object or null if not authenticated
 * @throws Error if session verification fails
 */
export async function getOrySession() {
  try {
    const ory = getOryClient();
    const cookieStore = await cookies();
    const sessionCookie = cookieStore.get("ory_kratos_session");

    if (!sessionCookie?.value) {
      return null;
    }

    const { data: session } = await ory.toSession(undefined, {
      headers: {
        Cookie: `ory_kratos_session=${sessionCookie.value}`,
      },
    });

    if (!session?.active) {
      return null;
    }

    return session;
  } catch (error) {
    console.error("[Ory] Session verification failed:", error);
    return null;
  }
}

/**
 * Require Ory session (throws if not authenticated).
 * Use in Server Actions that must be authenticated.
 *
 * @returns Ory session object
 * @throws Error if user is not authenticated
 */
export async function requireOrySession() {
  const session = await getOrySession();

  if (!session) {
    throw new Error("Unauthorized: No active session");
  }

  return session;
}

/**
 * Get user ID from session.
 * Convenience utility for extracting user ID.
 */
export async function getUserId(): Promise<string | null> {
  const session = await getOrySession();
  return session?.identity?.id ?? null;
}

/**
 * Get user email from session.
 * Convenience utility for extracting user email.
 */
export async function getUserEmail(): Promise<string | null> {
  const session = await getOrySession();
  return session?.identity?.traits?.email ?? null;
}

/**
 * Verify session and extract userId for use in RSCs.
 * Called from route handlers to ensure user is authenticated.
 *
 * @returns userId
 * @throws Error if user is not authenticated
 */
export async function getAuthenticatedUserId(): Promise<string> {
  const session = await requireOrySession();
  const userId = session.identity?.id;

  if (!userId) {
    throw new Error("Invalid session: No user ID");
  }

  return userId;
}
