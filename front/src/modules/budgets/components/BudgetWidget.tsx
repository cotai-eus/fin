/**
 * Budget Widget Component
 * Shows budget progress with visual indicators
 */

"use client";

import { Card } from "@/shared/components/ui/Card";
import { Badge } from "@/shared/components/ui/Badge";
import type { Budget } from "../types";
import { CATEGORY_LABELS, CATEGORY_ICONS, getCategoryBudgetStatus } from "../validators";

interface BudgetWidgetProps {
  budget: Budget;
  onEdit?: (budget: Budget) => void;
}

export function BudgetWidget({ budget, onEdit }: BudgetWidgetProps) {
  const { percentage, status, color } = getCategoryBudgetStatus(
    budget.currentSpent,
    budget.limit
  );

  const getStatusBadge = () => {
    switch (status) {
      case "safe":
        return <Badge variant="completed">Normal</Badge>;
      case "warning":
        return <Badge variant="pending">Atenção</Badge>;
      case "danger":
      case "exceeded":
        return <Badge variant="failed">Limite atingido</Badge>;
    }
  };

  const getProgressColor = () => {
    switch (color) {
      case "green":
        return "bg-green-500";
      case "yellow":
        return "bg-yellow-500";
      case "red":
        return "bg-red-500";
      default:
        return "bg-gray-500";
    }
  };

  return (
    <Card
      className="p-4 cursor-pointer hover:shadow-md transition"
      onClick={() => onEdit?.(budget)}
    >
      <div className="flex items-start justify-between mb-3">
        <div className="flex items-center gap-2">
          <span className="text-2xl">
            {CATEGORY_ICONS[budget.category]}
          </span>
          <div>
            <h4 className="font-semibold">
              {CATEGORY_LABELS[budget.category]}
            </h4>
            <p className="text-xs text-gray-500">
              {budget.period === "monthly" ? "Mensal" : "Semanal"}
            </p>
          </div>
        </div>
        {getStatusBadge()}
      </div>

      <div className="space-y-2">
        <div className="flex justify-between text-sm">
          <span className="text-gray-600">Gasto</span>
          <span className="font-semibold">
            R$ {budget.currentSpent.toFixed(2)} / R$ {budget.limit.toFixed(2)}
          </span>
        </div>

        <div className="w-full bg-gray-200 rounded-full h-2">
          <div
            className={`h-2 rounded-full transition-all ${getProgressColor()}`}
            style={{ width: `${Math.min(percentage, 100)}%` }}
          />
        </div>

        <div className="flex justify-between text-xs text-gray-500">
          <span>{percentage.toFixed(0)}% utilizado</span>
          <span>R$ {(budget.limit - budget.currentSpent).toFixed(2)} restante</span>
        </div>
      </div>

      {budget.alertsEnabled && percentage >= budget.alertThreshold && (
        <div className="mt-3 p-2 bg-yellow-50 border border-yellow-200 rounded text-xs text-yellow-800">
          ⚠️ Você atingiu {percentage.toFixed(0)}% do seu orçamento
        </div>
      )}
    </Card>
  );
}
