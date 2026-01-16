/**
 * Support Module - Server Actions
 */

"use server";

import { revalidatePath } from "next/cache";
import { requireOrySession } from "@/core/ory/session";
import {
  createTicketSchema,
  addTicketMessageSchema,
  type CreateTicketInput,
  type AddTicketMessageInput,
} from "../validators";
import type { SupportTicket, TicketMessage, FAQCategory } from "../types";

const BACKEND_URL = process.env.BACKEND_API_URL || "http://localhost:8080";

type ActionResult<T> =
  | { success: true; data: T }
  | { success: false; error: string };

/**
 * Create Support Ticket
 */
export async function createSupportTicket(
  input: unknown
): Promise<ActionResult<SupportTicket>> {
  try {
    const session = await requireOrySession();
    const userId = session.identity?.id;
    if (!userId) {
      return { success: false, error: "Unauthorized" };
    }

    const validated = createTicketSchema.safeParse(input);
    if (!validated.success) {
      return {
        success: false,
        error: validated.error.issues[0]?.message || "Invalid input",
      };
    }

    const response = await fetch(`${BACKEND_URL}/api/support/tickets`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-User-ID": userId,
        "X-Request-ID": crypto.randomUUID(),
      },
      body: JSON.stringify(validated.data),
    });

    if (!response.ok) {
      const error = await response
        .json()
        .catch(() => ({ message: "Ticket creation failed" }));
      return { success: false, error: error.message || "Ticket creation failed" };
    }

    const ticket: SupportTicket = await response.json();
    revalidatePath("/dashboard/support");

    return { success: true, data: ticket };
  } catch (error) {
    console.error("Create ticket error:", error);
    return { success: false, error: "An unexpected error occurred" };
  }
}

/**
 * Fetch User's Support Tickets
 */
export async function fetchUserTickets(): Promise<ActionResult<SupportTicket[]>> {
  try {
    const session = await requireOrySession();
    const userId = session.identity?.id;
    if (!userId) {
      return { success: false, error: "Unauthorized" };
    }

    const response = await fetch(`${BACKEND_URL}/api/support/tickets`, {
      method: "GET",
      headers: {
        "X-User-ID": userId,
      },
    });

    if (!response.ok) {
      return { success: false, error: "Failed to fetch tickets" };
    }

    const tickets: SupportTicket[] = await response.json();
    return { success: true, data: tickets };
  } catch (error) {
    console.error("Fetch tickets error:", error);
    return { success: false, error: "An unexpected error occurred" };
  }
}

/**
 * Add Message to Ticket
 */
export async function addTicketMessage(
  input: unknown
): Promise<ActionResult<TicketMessage>> {
  try {
    const session = await requireOrySession();
    const userId = session.identity?.id;
    if (!userId) {
      return { success: false, error: "Unauthorized" };
    }

    const validated = addTicketMessageSchema.safeParse(input);
    if (!validated.success) {
      return {
        success: false,
        error: validated.error.issues[0]?.message || "Invalid input",
      };
    }

    const { ticketId, message } = validated.data;

    const response = await fetch(
      `${BACKEND_URL}/api/support/tickets/${ticketId}/messages`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "X-User-ID": userId,
          "X-Request-ID": crypto.randomUUID(),
        },
        body: JSON.stringify({ message }),
      }
    );

    if (!response.ok) {
      const error = await response
        .json()
        .catch(() => ({ message: "Failed to send message" }));
      return { success: false, error: error.message || "Failed to send message" };
    }

    const ticketMessage: TicketMessage = await response.json();
    revalidatePath("/dashboard/support");

    return { success: true, data: ticketMessage };
  } catch (error) {
    console.error("Add ticket message error:", error);
    return { success: false, error: "An unexpected error occurred" };
  }
}

/**
 * Fetch Ticket Messages
 */
export async function fetchTicketMessages(
  ticketId: string
): Promise<ActionResult<TicketMessage[]>> {
  try {
    const session = await requireOrySession();
    const userId = session.identity?.id;
    if (!userId) {
      return { success: false, error: "Unauthorized" };
    }

    const response = await fetch(
      `${BACKEND_URL}/api/support/tickets/${ticketId}/messages`,
      {
        method: "GET",
        headers: {
          "X-User-ID": userId,
        },
      }
    );

    if (!response.ok) {
      return { success: false, error: "Failed to fetch messages" };
    }

    const messages: TicketMessage[] = await response.json();
    return { success: true, data: messages };
  } catch (error) {
    console.error("Fetch ticket messages error:", error);
    return { success: false, error: "An unexpected error occurred" };
  }
}

/**
 * Get FAQ Items (Static - could be from CMS)
 */
export async function getFAQCategories(): Promise<ActionResult<FAQCategory[]>> {
  // In a real app, this would fetch from backend/CMS
  // For now, return static FAQ data
  const faqData: FAQCategory[] = [
    {
      name: "Conta e Segurança",
      icon: "🔒",
      items: [
        {
          id: "1",
          category: "account",
          question: "Como altero minha senha?",
          answer:
            "Acesse Configurações > Segurança > Alterar Senha. Você precisará informar sua senha atual e a nova senha duas vezes para confirmação.",
          order: 1,
        },
        {
          id: "2",
          category: "account",
          question: "Como ativo a autenticação de dois fatores?",
          answer:
            "Vá em Configurações > Segurança > Autenticação de Dois Fatores. Você pode escolher entre SMS ou aplicativo autenticador (recomendado).",
          order: 2,
        },
      ],
    },
    {
      name: "Transferências e Pagamentos",
      icon: "💸",
      items: [
        {
          id: "3",
          category: "transfer",
          question: "Qual o limite para transferências PIX?",
          answer:
            "O limite padrão para transferências PIX é de R$ 1.000 por transação e R$ 5.000 por dia. Você pode solicitar aumento de limite através do suporte.",
          order: 1,
        },
        {
          id: "4",
          category: "transfer",
          question: "Quanto tempo leva uma TED?",
          answer:
            "TEDs realizadas em dias úteis até às 17h são processadas no mesmo dia. Após esse horário ou em fins de semana/feriados, serão processadas no próximo dia útil.",
          order: 2,
        },
      ],
    },
    {
      name: "Cartões",
      icon: "💳",
      items: [
        {
          id: "5",
          category: "card",
          question: "Como bloqueio meu cartão?",
          answer:
            "Acesse Cartões, selecione o cartão desejado e clique em 'Bloquear'. O bloqueio é instantâneo e você pode desbloquear a qualquer momento.",
          order: 1,
        },
        {
          id: "6",
          category: "card",
          question: "Posso criar um cartão virtual?",
          answer:
            "Sim! Na área de Cartões, clique em 'Criar Cartão Virtual'. Você pode definir limites personalizados e usar para compras online com mais segurança.",
          order: 2,
        },
      ],
    },
  ];

  return { success: true, data: faqData };
}
