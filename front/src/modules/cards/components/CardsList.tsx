/**
 * Cards List Component
 * Displays all user cards
 */

"use client";

import { CardItem } from "./CardItem";
import type { Card } from "../types";

interface CardsListProps {
  cards: Card[];
}

export function CardsList({ cards }: CardsListProps) {
  if (cards.length === 0) {
    return (
      <div className="text-center py-12 text-gray-500">
        <p className="text-lg mb-2">Nenhum cartão encontrado</p>
        <p className="text-sm">Crie um cartão virtual para começar</p>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
      {cards.map((card) => (
        <CardItem key={card.id} card={card} />
      ))}
    </div>
  );
}
