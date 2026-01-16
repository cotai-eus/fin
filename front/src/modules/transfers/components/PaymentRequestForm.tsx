/**
 * Payment Request Component
 * Generates payment link/QR code for requesting money
 */

"use client";

import { useState } from "react";
import Image from "next/image";
import { createPaymentRequest } from "../actions";
import type { PaymentRequest } from "../types";
import { Button } from "@/shared/components/ui/Button";
import { Card } from "@/shared/components/ui/Card";
import { Badge } from "@/shared/components/ui/Badge";

export function PaymentRequestForm() {
  const [amount, setAmount] = useState<number>(0);
  const [description, setDescription] = useState("");
  const [expiresInHours, setExpiresInHours] = useState<number>(24);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [paymentRequest, setPaymentRequest] = useState<PaymentRequest | null>(
    null
  );

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    const result = await createPaymentRequest({
      amount,
      description,
      expiresInHours,
    });

    if (result.success) {
      setPaymentRequest(result.data);
    } else {
      setError(result.error);
    }

    setLoading(false);
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    alert("Link copiado!");
  };

  if (paymentRequest) {
    return (
      <Card className="p-6">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-lg font-semibold">Solicitação de Pagamento</h3>
          <Badge variant={paymentRequest.status === "pending" ? "pending" : "completed"}>
            {paymentRequest.status}
          </Badge>
        </div>

        <div className="space-y-4">
          <div>
            <p className="text-sm text-gray-600">Valor solicitado</p>
            <p className="text-2xl font-bold">
              R$ {paymentRequest.amount.toFixed(2)}
            </p>
          </div>

          <div>
            <p className="text-sm text-gray-600">Descrição</p>
            <p className="font-medium">{paymentRequest.description}</p>
          </div>

          {paymentRequest.qrCode && (
            <div className="text-center">
              <div className="mx-auto w-48 h-48 relative">
                <Image
                  src={paymentRequest.qrCode}
                  alt="QR Code"
                  width={192}
                  height={192}
                  className="w-full h-full"
                />
              </div>
              <p className="text-xs text-gray-500 mt-2">
                Escaneie para pagar
              </p>
            </div>
          )}

          <div className="p-3 bg-gray-50 rounded-lg">
            <p className="text-xs text-gray-500 mb-1">Link de pagamento:</p>
            <div className="flex gap-2">
              <code className="text-sm flex-1 break-all">
                {paymentRequest.paymentLink}
              </code>
              <Button
                variant="outline"
                size="sm"
                onClick={() => copyToClipboard(paymentRequest.paymentLink)}
              >
                Copiar
              </Button>
            </div>
          </div>

          {paymentRequest.expiresAt && (
            <p className="text-xs text-gray-500">
              Expira em:{" "}
              {new Date(paymentRequest.expiresAt).toLocaleString("pt-BR")}
            </p>
          )}

          <Button
            variant="outline"
            onClick={() => setPaymentRequest(null)}
            className="w-full"
          >
            Nova Solicitação
          </Button>
        </div>
      </Card>
    );
  }

  return (
    <Card className="p-6">
      <h3 className="text-lg font-semibold mb-4">Solicitar Dinheiro</h3>

      <form onSubmit={handleCreate} className="space-y-4">
        <div>
          <label className="block text-sm font-medium mb-2">Valor (R$)</label>
          <input
            type="number"
            step="0.01"
            min="0.01"
            max="1000000"
            value={amount || ""}
            onChange={(e) => setAmount(parseFloat(e.target.value) || 0)}
            placeholder="0,00"
            required
            className="w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label className="block text-sm font-medium mb-2">Descrição</label>
          <textarea
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="Ex: Pagamento de almoço compartilhado"
            maxLength={200}
            rows={3}
            required
            className="w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label className="block text-sm font-medium mb-2">
            Expira em (horas)
          </label>
          <select
            value={expiresInHours}
            onChange={(e) => setExpiresInHours(parseInt(e.target.value))}
            className="w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500"
          >
            <option value={1}>1 hora</option>
            <option value={6}>6 horas</option>
            <option value={24}>24 horas</option>
            <option value={72}>3 dias</option>
            <option value={168}>7 dias</option>
          </select>
        </div>

        {error && (
          <div className="p-3 bg-red-50 border border-red-200 rounded-lg text-red-700 text-sm">
            {error}
          </div>
        )}

        <Button type="submit" disabled={loading} className="w-full">
          {loading ? "Criando..." : "Gerar Link de Pagamento"}
        </Button>
      </form>
    </Card>
  );
}
