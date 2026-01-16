/**
 * Support Module - Type Definitions
 * Domain: Customer support (FAQ, live chat, tickets)
 */

export enum TicketStatus {
  Open = "open",
  InProgress = "in_progress",
  Resolved = "resolved",
  Closed = "closed",
}

export enum TicketPriority {
  Low = "low",
  Medium = "medium",
  High = "high",
  Urgent = "urgent",
}

export enum TicketCategory {
  ACCOUNT = "account",
  CARD = "card",
  TRANSFER = "transfer",
  BILL = "bill",
  TECHNICAL = "technical",
  OTHER = "other",
}

export interface SupportTicket {
  id: string;
  userId: string;
  category: TicketCategory;
  subject: string;
  description: string;
  status: TicketStatus;
  priority: TicketPriority;
  assignedTo?: string;
  createdAt: string;
  updatedAt: string;
  resolvedAt?: string;
}

export interface TicketMessage {
  id: string;
  ticketId: string;
  userId: string;
  userName: string;
  message: string;
  isStaff: boolean;
  createdAt: string;
}

export interface ChatMessage {
  id: string;
  userId: string;
  userName: string;
  message: string;
  isStaff: boolean;
  timestamp: string;
}

export interface FAQItem {
  id: string;
  category: string;
  question: string;
  answer: string;
  order: number;
}

export interface FAQCategory {
  name: string;
  icon: string;
  items: FAQItem[];
}
