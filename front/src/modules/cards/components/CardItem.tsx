/**
 * Card Item Component
 * Displays a single card with controls
 */

"use client";

import { useState } from "react";
import { Card } from "@/shared/components/ui/Card";
import { Badge } from "@/shared/components/ui/Badge";
import { Button } from "@/shared/components/ui/Button";
import type { Card as CardType, CardDetails } from "../types";
import { CardStatus, CardType as CardTypeEnum } from "../types";
import { getCardDetails, toggleCardStatus } from "../actions";

interface CardItemProps {
  card: CardType;
}

export function CardItem({ card }: CardItemProps) {
  const [showDetails, setShowDetails] = useState(false);
  const [details, setDetails] = useState<CardDetails | null>(null);
  const [loading, setLoading] = useState(false);

  const handleShowDetails = async () => {
    if (showDetails) {
      setShowDetails(false);
      setDetails(null);
      return;
    }

    setLoading(true);
    const result = await getCardDetails(card.id);
    if (result.success) {
      setDetails(result.data);
      setShowDetails(true);
    }
    setLoading(false);
  };

  const handleToggleBlock = async () => {
    const action = card.status === CardStatus.Active ? "block" : "unblock";
    setLoading(true);
    await toggleCardStatus({
      cardId: card.id,
      action,
      reason: action === "block" ? "User requested" : undefined,
    });
    setLoading(false);
  };

  const getStatusVariant = (status: CardStatus) => {
    switch (status) {
      case CardStatus.Active:
        return "completed" as const;
      case CardStatus.Blocked:
        return "pending" as const;
      case CardStatus.Cancelled:
      case CardStatus.Lost:
      case CardStatus.Stolen:
        return "failed" as const;
      default:
        return "cancelled" as const;
    }
  };

  const limitPercentage = (card.currentMonthlySpent / card.monthlyLimit) * 100;

  return (
    <Card className="p-6">
      <div className="flex items-start justify-between mb-4">
        <div className="flex-1">
          <div className="flex items-center gap-2 mb-1">
            <h3 className="font-semibold">
              {card.type === CardTypeEnum.PHYSICAL ? "Cartão Físico" : "Cartão Virtual"}
            </h3>
            <Badge variant={getStatusVariant(card.status)}>
              {card.status}
            </Badge>
          </div>
          <p className="text-2xl font-mono tracking-wider">
            **** **** **** {card.lastFourDigits}
          </p>
          <p className="text-sm text-gray-600 mt-1">{card.holderName}</p>
        </div>

        <div className="text-right">
          <p className="text-sm text-gray-600">Validade</p>
          <p className="font-mono">
            {String(card.expiryMonth).padStart(2, "0")}/{card.expiryYear}
          </p>
        </div>
      </div>

      {/* Spending Progress */}
      <div className="mb-4">
        <div className="flex justify-between text-sm mb-1">
          <span className="text-gray-600">Limite mensal</span>
          <span className="font-medium">
            R$ {card.currentMonthlySpent.toFixed(2)} / R${" "}
            {card.monthlyLimit.toFixed(2)}
          </span>
        </div>
        <div className="w-full bg-gray-200 rounded-full h-2">
          <div
            className={`h-2 rounded-full transition-all ${
              limitPercentage > 90
                ? "bg-red-500"
                : limitPercentage > 70
                ? "bg-yellow-500"
                : "bg-green-500"
            }`}
            style={{ width: `${Math.min(limitPercentage, 100)}%` }}
          />
        </div>
      </div>

      {/* Card Details (only when revealed) */}
      {showDetails && details && (
        <div className="mb-4 p-4 bg-gray-50 rounded-lg border">
          <p className="text-sm text-gray-600 mb-2">Número do cartão</p>
          <p className="font-mono text-lg mb-3">{details.cardNumber}</p>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <p className="text-sm text-gray-600">CVV</p>
              <p className="font-mono text-lg">{details.cvv}</p>
            </div>
            <div>
              <p className="text-sm text-gray-600">Validade</p>
              <p className="font-mono text-lg">
                {String(details.expiryMonth).padStart(2, "0")}/
                {details.expiryYear}
              </p>
            </div>
          </div>
        </div>
      )}

      {/* Features */}
      <div className="flex gap-2 mb-4 flex-wrap">
        {card.isContactless && (
          <Badge variant="pending" className="text-xs">
            Contactless
          </Badge>
        )}
        {card.isInternational && (
          <Badge variant="pending" className="text-xs">
            Internacional
          </Badge>
        )}
        <Badge variant="cancelled" className="text-xs">
          {card.brand.toUpperCase()}
        </Badge>
      </div>

      {/* Actions */}
      <div className="grid grid-cols-2 gap-2">
        <Button
          variant="outline"
          size="sm"
          onClick={handleShowDetails}
          disabled={loading || card.status !== CardStatus.Active}
        >
          {loading ? "Carregando..." : showDetails ? "Ocultar Dados" : "Ver Dados"}
        </Button>

        <Button
          variant={card.status === CardStatus.Active ? "destructive" : "primary"}
          size="sm"
          onClick={handleToggleBlock}
          disabled={
            loading ||
            card.status === CardStatus.Cancelled ||
            card.status === CardStatus.Lost ||
            card.status === CardStatus.Stolen
          }
        >
          {card.status === CardStatus.Active ? "Bloquear" : "Desbloquear"}
        </Button>
      </div>
    </Card>
  );
}
