/**
 * Card Actions Component
 * Provides actions for card management (report, change PIN, settings)
 */

"use client";

import { useState } from "react";
import { Button } from "@/shared/components/ui/Button";
import { Card } from "@/shared/components/ui/Card";
import { reportCard, changeCardPIN, updateCardLimits } from "../actions";
import type { Card as CardType } from "../types";

interface CardActionsProps {
  card: CardType;
}

export function CardActions({ card }: CardActionsProps) {
  const [activeAction, setActiveAction] = useState<
    "report" | "pin" | "limits" | null
  >(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  // Report Card Form State
  const [reportType, setReportType] = useState<"lost" | "stolen">("lost");
  const [reportDescription, setReportDescription] = useState("");

  // Change PIN Form State
  const [currentPIN, setCurrentPIN] = useState("");
  const [newPIN, setNewPIN] = useState("");
  const [confirmPIN, setConfirmPIN] = useState("");

  // Limits Form State
  const [dailyLimit, setDailyLimit] = useState(card.dailyLimit);
  const [monthlyLimit, setMonthlyLimit] = useState(card.monthlyLimit);

  const resetForm = () => {
    setActiveAction(null);
    setLoading(false);
    setError(null);
    setSuccess(false);
    setReportDescription("");
    setCurrentPIN("");
    setNewPIN("");
    setConfirmPIN("");
  };

  const handleReportCard = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    const result = await reportCard({
      cardId: card.id,
      reportType,
      description: reportDescription,
      requestReplacement: true,
    });

    if (result.success) {
      setSuccess(true);
      setTimeout(resetForm, 2000);
    } else {
      setError(result.error);
    }

    setLoading(false);
  };

  const handleChangePIN = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    const result = await changeCardPIN({
      cardId: card.id,
      currentPIN,
      newPIN,
      confirmPIN,
    });

    if (result.success) {
      setSuccess(true);
      setTimeout(resetForm, 2000);
    } else {
      setError(result.error);
    }

    setLoading(false);
  };

  const handleUpdateLimits = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    const result = await updateCardLimits({
      cardId: card.id,
      dailyLimit,
      monthlyLimit,
    });

    if (result.success) {
      setSuccess(true);
      setTimeout(resetForm, 2000);
    } else {
      setError(result.error);
    }

    setLoading(false);
  };

  return (
    <div className="space-y-4">
      {/* Action Buttons */}
      {!activeAction && (
        <div className="grid grid-cols-3 gap-2">
          <Button
            variant="outline"
            onClick={() => setActiveAction("report")}
            className="text-sm"
          >
            Reportar Perda
          </Button>
          <Button
            variant="outline"
            onClick={() => setActiveAction("pin")}
            className="text-sm"
          >
            Alterar Senha
          </Button>
          <Button
            variant="outline"
            onClick={() => setActiveAction("limits")}
            className="text-sm"
          >
            Ajustar Limites
          </Button>
        </div>
      )}

      {/* Report Card Form */}
      {activeAction === "report" && (
        <Card className="p-4">
          <h4 className="font-semibold mb-3">Reportar Cartão</h4>
          <form onSubmit={handleReportCard} className="space-y-3">
            <div>
              <label className="block text-sm font-medium mb-1">Motivo</label>
              <select
                value={reportType}
                onChange={(e) =>
                  setReportType(e.target.value as "lost" | "stolen")
                }
                className="w-full px-3 py-2 border rounded-lg text-sm"
              >
                <option value="lost">Perda</option>
                <option value="stolen">Roubo</option>
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium mb-1">
                Descrição
              </label>
              <textarea
                value={reportDescription}
                onChange={(e) => setReportDescription(e.target.value)}
                placeholder="Descreva o ocorrido"
                rows={3}
                required
                minLength={10}
                className="w-full px-3 py-2 border rounded-lg text-sm"
              />
            </div>

            {error && (
              <div className="p-2 bg-red-50 border border-red-200 rounded text-red-700 text-xs">
                {error}
              </div>
            )}

            {success && (
              <div className="p-2 bg-green-50 border border-green-200 rounded text-green-700 text-xs">
                Cartão reportado com sucesso!
              </div>
            )}

            <div className="flex gap-2">
              <Button type="submit" disabled={loading} size="sm" className="flex-1">
                {loading ? "Processando..." : "Confirmar"}
              </Button>
              <Button
                type="button"
                variant="outline"
                onClick={resetForm}
                size="sm"
              >
                Cancelar
              </Button>
            </div>
          </form>
        </Card>
      )}

      {/* Change PIN Form */}
      {activeAction === "pin" && (
        <Card className="p-4">
          <h4 className="font-semibold mb-3">Alterar Senha do Cartão</h4>
          <form onSubmit={handleChangePIN} className="space-y-3">
            <div>
              <label className="block text-sm font-medium mb-1">
                Senha Atual
              </label>
              <input
                type="password"
                value={currentPIN}
                onChange={(e) => setCurrentPIN(e.target.value)}
                placeholder="****"
                maxLength={4}
                pattern="\d{4}"
                required
                className="w-full px-3 py-2 border rounded-lg text-sm font-mono"
              />
            </div>

            <div>
              <label className="block text-sm font-medium mb-1">
                Nova Senha
              </label>
              <input
                type="password"
                value={newPIN}
                onChange={(e) => setNewPIN(e.target.value)}
                placeholder="****"
                maxLength={4}
                pattern="\d{4}"
                required
                className="w-full px-3 py-2 border rounded-lg text-sm font-mono"
              />
            </div>

            <div>
              <label className="block text-sm font-medium mb-1">
                Confirmar Nova Senha
              </label>
              <input
                type="password"
                value={confirmPIN}
                onChange={(e) => setConfirmPIN(e.target.value)}
                placeholder="****"
                maxLength={4}
                pattern="\d{4}"
                required
                className="w-full px-3 py-2 border rounded-lg text-sm font-mono"
              />
            </div>

            {error && (
              <div className="p-2 bg-red-50 border border-red-200 rounded text-red-700 text-xs">
                {error}
              </div>
            )}

            {success && (
              <div className="p-2 bg-green-50 border border-green-200 rounded text-green-700 text-xs">
                Senha alterada com sucesso!
              </div>
            )}

            <div className="flex gap-2">
              <Button type="submit" disabled={loading} size="sm" className="flex-1">
                {loading ? "Processando..." : "Alterar Senha"}
              </Button>
              <Button
                type="button"
                variant="outline"
                onClick={resetForm}
                size="sm"
              >
                Cancelar
              </Button>
            </div>
          </form>
        </Card>
      )}

      {/* Update Limits Form */}
      {activeAction === "limits" && (
        <Card className="p-4">
          <h4 className="font-semibold mb-3">Ajustar Limites</h4>
          <form onSubmit={handleUpdateLimits} className="space-y-3">
            <div>
              <label className="block text-sm font-medium mb-1">
                Limite Diário (R$)
              </label>
              <input
                type="number"
                step="0.01"
                min="0"
                max="50000"
                value={dailyLimit}
                onChange={(e) => setDailyLimit(parseFloat(e.target.value))}
                className="w-full px-3 py-2 border rounded-lg text-sm"
              />
            </div>

            <div>
              <label className="block text-sm font-medium mb-1">
                Limite Mensal (R$)
              </label>
              <input
                type="number"
                step="0.01"
                min="0"
                max="500000"
                value={monthlyLimit}
                onChange={(e) => setMonthlyLimit(parseFloat(e.target.value))}
                className="w-full px-3 py-2 border rounded-lg text-sm"
              />
            </div>

            {error && (
              <div className="p-2 bg-red-50 border border-red-200 rounded text-red-700 text-xs">
                {error}
              </div>
            )}

            {success && (
              <div className="p-2 bg-green-50 border border-green-200 rounded text-green-700 text-xs">
                Limites atualizados com sucesso!
              </div>
            )}

            <div className="flex gap-2">
              <Button type="submit" disabled={loading} size="sm" className="flex-1">
                {loading ? "Salvando..." : "Salvar Limites"}
              </Button>
              <Button
                type="button"
                variant="outline"
                onClick={resetForm}
                size="sm"
              >
                Cancelar
              </Button>
            </div>
          </form>
        </Card>
      )}
    </div>
  );
}
