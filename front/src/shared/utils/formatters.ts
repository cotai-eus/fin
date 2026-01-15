/**
 * src/shared/utils/formatters.ts
 *
 * Formatting utilities for common data types
 * Used across all components for consistent display
 */

/**
 * Format a number as Brazilian Real currency
 * @example formatCurrency(1234.56) => "R$ 1.234,56"
 */
export function formatCurrency(amount: number): string {
  return new Intl.NumberFormat("pt-BR", {
    style: "currency",
    currency: "BRL",
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(amount);
}

/**
 * Format a date with relative time display
 * @example formatDate("2024-01-15T10:30:00Z") => "Jan 15, 2024 • 2 hours ago"
 */
export function formatDate(date: string | Date): string {
  const d = typeof date === "string" ? new Date(date) : date;
  const today = new Date();
  const diffMs = today.getTime() - d.getTime();
  const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

  let relativeTime = "";
  if (diffHours < 1) {
    relativeTime = "just now";
  } else if (diffHours < 24) {
    relativeTime = `${diffHours} ${diffHours === 1 ? "hour" : "hours"} ago`;
  } else if (diffDays < 7) {
    relativeTime = `${diffDays} ${diffDays === 1 ? "day" : "days"} ago`;
  } else {
    relativeTime = `${Math.floor(diffDays / 7)} weeks ago`;
  }

  const formatted = new Intl.DateTimeFormat("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
  }).format(d);

  return `${formatted} • ${relativeTime}`;
}

/**
 * Format a phone number to Brazilian format
 * @example formatPhoneNumber("11987654321") => "(11) 98765-4321"
 */
export function formatPhoneNumber(phone: string): string {
  const cleaned = phone.replace(/\D/g, "");
  if (cleaned.length !== 11) return phone;
  return `(${cleaned.slice(0, 2)}) ${cleaned.slice(2, 7)}-${cleaned.slice(7)}`;
}

/**
 * Format a CPF number to standard format
 * @example formatCPF("12345678901") => "123.456.789-01"
 */
export function formatCPF(cpf: string): string {
  const cleaned = cpf.replace(/\D/g, "");
  if (cleaned.length !== 11) return cpf;
  return `${cleaned.slice(0, 3)}.${cleaned.slice(3, 6)}.${cleaned.slice(6, 9)}-${cleaned.slice(9)}`;
}

/**
 * Truncate a string and add ellipsis
 * @example truncate("Long string text", 10) => "Long stri..."
 */
export function truncate(text: string, length: number): string {
  return text.length > length ? `${text.slice(0, length)}...` : text;
}

/**
 * Format a transaction ID for display (show first and last 4 chars)
 * @example formatTransactionId("abc123def456") => "abc1...f456"
 */
export function formatTransactionId(id: string): string {
  if (id.length <= 8) return id;
  return `${id.slice(0, 4)}...${id.slice(-4)}`;
}
