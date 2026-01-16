/**
 * Ticket History Component
 * Shows user's support tickets
 */

"use client";

import { Card } from "@/shared/components/ui/Card";
import { Badge } from "@/shared/components/ui/Badge";
import type { SupportTicket } from "../types";
import { TicketStatus } from "../types";
import { CATEGORY_LABELS, PRIORITY_LABELS, STATUS_LABELS } from "../validators";

interface TicketHistoryProps {
  tickets: SupportTicket[];
}

export function TicketHistory({ tickets }: TicketHistoryProps) {
  const getStatusVariant = (status: TicketStatus) => {
    switch (status) {
      case TicketStatus.Open:
        return "pending" as const;
      case TicketStatus.InProgress:
        return "pending" as const;
      case TicketStatus.Resolved:
        return "completed" as const;
      case TicketStatus.Closed:
        return "cancelled" as const;
      default:
        return "cancelled" as const;
    }
  };

  if (tickets.length === 0) {
    return (
      <Card className="p-12 text-center">
        <p className="text-gray-600">Nenhum ticket de suporte encontrado</p>
      </Card>
    );
  }

  return (
    <div className="space-y-3">
      {tickets.map((ticket) => (
        <Card key={ticket.id} className="p-4 hover:shadow-md transition cursor-pointer">
          <div className="flex items-start justify-between mb-2">
            <div className="flex-1">
              <div className="flex items-center gap-2 mb-1">
                <h4 className="font-semibold">{ticket.subject}</h4>
                <Badge variant={getStatusVariant(ticket.status)}>
                  {STATUS_LABELS[ticket.status]}
                </Badge>
              </div>
              <p className="text-sm text-gray-600 line-clamp-2">
                {ticket.description}
              </p>
            </div>
          </div>

          <div className="flex items-center gap-4 text-xs text-gray-500 mt-3">
            <span>
              Categoria: <strong>{CATEGORY_LABELS[ticket.category]}</strong>
            </span>
            <span>
              Prioridade: <strong>{PRIORITY_LABELS[ticket.priority]}</strong>
            </span>
            <span>
              Criado em:{" "}
              <strong>
                {new Date(ticket.createdAt).toLocaleDateString("pt-BR")}
              </strong>
            </span>
          </div>
        </Card>
      ))}
    </div>
  );
}
