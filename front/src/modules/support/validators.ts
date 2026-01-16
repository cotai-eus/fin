/**
 * Support Module - Zod Validation Schemas
 */

import { z } from "zod";
import { TicketCategory, TicketPriority } from "./types";

/**
 * Create Support Ticket Schema
 */
export const createTicketSchema = z.object({
  category: z.nativeEnum(TicketCategory),
  subject: z
    .string()
    .min(5, "Subject must be at least 5 characters")
    .max(200, "Subject too long"),
  description: z
    .string()
    .min(20, "Please provide more details (at least 20 characters)")
    .max(2000, "Description too long"),
  priority: z.nativeEnum(TicketPriority).default(TicketPriority.Medium),
});

export type CreateTicketInput = z.infer<typeof createTicketSchema>;

/**
 * Add Ticket Message Schema
 */
export const addTicketMessageSchema = z.object({
  ticketId: z.string().uuid("Invalid ticket ID"),
  message: z
    .string()
    .min(1, "Message cannot be empty")
    .max(1000, "Message too long"),
});

export type AddTicketMessageInput = z.infer<typeof addTicketMessageSchema>;

/**
 * Send Chat Message Schema
 */
export const sendChatMessageSchema = z.object({
  message: z
    .string()
    .min(1, "Message cannot be empty")
    .max(500, "Message too long"),
});

export type SendChatMessageInput = z.infer<typeof sendChatMessageSchema>;

/**
 * Category Labels
 */
export const CATEGORY_LABELS: Record<TicketCategory, string> = {
  [TicketCategory.ACCOUNT]: "Conta",
  [TicketCategory.CARD]: "Cartões",
  [TicketCategory.TRANSFER]: "Transferências",
  [TicketCategory.BILL]: "Pagamentos",
  [TicketCategory.TECHNICAL]: "Técnico",
  [TicketCategory.OTHER]: "Outros",
};

/**
 * Priority Labels
 */
export const PRIORITY_LABELS: Record<TicketPriority, string> = {
  [TicketPriority.Low]: "Baixa",
  [TicketPriority.Medium]: "Média",
  [TicketPriority.High]: "Alta",
  [TicketPriority.Urgent]: "Urgente",
};

/**
 * Status Labels
 */
export const STATUS_LABELS = {
  open: "Aberto",
  in_progress: "Em Andamento",
  resolved: "Resolvido",
  closed: "Fechado",
};
