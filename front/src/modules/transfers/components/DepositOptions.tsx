/**
 * Deposit Options Component
 * Shows different deposit methods (PIX, Boleto, Bank Transfer)
 */

"use client";

import { useState } from "react";
import Image from "next/image";
import { createDeposit } from "../actions";
import { DepositMethod } from "../types";
import type { Deposit } from "../types";
import { Button } from "@/shared/components/ui/Button";
import { Card } from "@/shared/components/ui/Card";
import { Badge } from "@/shared/components/ui/Badge";

export function DepositOptions() {
  const [selectedMethod, setSelectedMethod] = useState<DepositMethod>(
    DepositMethod.PIX
  );
  const [amount, setAmount] = useState<number>(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [deposit, setDeposit] = useState<Deposit | null>(null);

  const handleCreateDeposit = async () => {
    if (amount <= 0) {
      setError("Valor deve ser maior que zero");
      return;
    }

    setLoading(true);
    setError(null);

    const result = await createDeposit({
      method: selectedMethod,
      amount,
    });

    if (result.success) {
      setDeposit(result.data);
    } else {
      setError(result.error);
    }

    setLoading(false);
  };

  return (
    <div className="space-y-6">
      <Card className="p-6">
        <h3 className="text-lg font-semibold mb-4">Adicionar Dinheiro</h3>

        {/* Method Selection */}
        <div className="mb-4">
          <label className="block text-sm font-medium mb-2">
            Escolha o método
          </label>
          <div className="grid grid-cols-3 gap-2">
            <button
              type="button"
              onClick={() => setSelectedMethod(DepositMethod.PIX)}
              className={`p-3 border rounded-lg text-center transition ${
                selectedMethod === DepositMethod.PIX
                  ? "border-blue-500 bg-blue-50 text-blue-700"
                  : "border-gray-200 hover:border-gray-300"
              }`}
            >
              <div className="font-semibold">PIX</div>
              <div className="text-xs text-gray-500">Instantâneo</div>
            </button>

            <button
              type="button"
              onClick={() => setSelectedMethod(DepositMethod.BOLETO)}
              className={`p-3 border rounded-lg text-center transition ${
                selectedMethod === DepositMethod.BOLETO
                  ? "border-blue-500 bg-blue-50 text-blue-700"
                  : "border-gray-200 hover:border-gray-300"
              }`}
            >
              <div className="font-semibold">Boleto</div>
              <div className="text-xs text-gray-500">1-2 dias úteis</div>
            </button>

            <button
              type="button"
              onClick={() => setSelectedMethod(DepositMethod.BANK_TRANSFER)}
              className={`p-3 border rounded-lg text-center transition ${
                selectedMethod === DepositMethod.BANK_TRANSFER
                  ? "border-blue-500 bg-blue-50 text-blue-700"
                  : "border-gray-200 hover:border-gray-300"
              }`}
            >
              <div className="font-semibold">TED</div>
              <div className="text-xs text-gray-500">1 dia útil</div>
            </button>
          </div>
        </div>

        {/* Amount Input */}
        <div className="mb-4">
          <label className="block text-sm font-medium mb-2">Valor (R$)</label>
          <input
            type="number"
            step="0.01"
            min="0.01"
            max="1000000"
            value={amount || ""}
            onChange={(e) => setAmount(parseFloat(e.target.value) || 0)}
            placeholder="0,00"
            className="w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500"
          />
        </div>

        {/* Error Message */}
        {error && (
          <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-lg text-red-700 text-sm">
            {error}
          </div>
        )}

        {/* Generate Button */}
        {!deposit && (
          <Button
            onClick={handleCreateDeposit}
            disabled={loading}
            className="w-full"
          >
            {loading ? "Gerando..." : "Gerar Código"}
          </Button>
        )}
      </Card>

      {/* Deposit Details */}
      {deposit && (
        <Card className="p-6">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold">Instruções de Depósito</h3>
            <Badge variant={deposit.status === "pending" ? "pending" : "completed"}>
              {deposit.status}
            </Badge>
          </div>

          {/* PIX QR Code */}
          {deposit.method === DepositMethod.PIX && deposit.pixQRCode && (
            <div className="text-center mb-4">
              <div className="mx-auto mb-2 w-48 h-48 relative">
                <Image
                  src={deposit.pixQRCode}
                  alt="QR Code PIX"
                  width={192}
                  height={192}
                  className="w-full h-full"
                />
              </div>
              <p className="text-sm text-gray-600 mb-2">
                Escaneie o QR Code acima com o app do seu banco
              </p>
              {deposit.pixKey && (
                <div className="p-3 bg-gray-50 rounded-lg">
                  <p className="text-xs text-gray-500 mb-1">Ou copie a chave:</p>
                  <code className="text-sm break-all">{deposit.pixKey}</code>
                </div>
              )}
            </div>
          )}

          {/* Boleto */}
          {deposit.method === DepositMethod.BOLETO && (
            <div className="space-y-3">
              {deposit.boletoCode && (
                <div className="p-3 bg-gray-50 rounded-lg">
                  <p className="text-xs text-gray-500 mb-1">Código do boleto:</p>
                  <code className="text-sm break-all">{deposit.boletoCode}</code>
                </div>
              )}
              {deposit.boletoUrl && (
                <Button
                  variant="outline"
                  onClick={() => window.open(deposit.boletoUrl, "_blank")}
                  className="w-full"
                >
                  Visualizar Boleto
                </Button>
              )}
            </div>
          )}

          {/* Bank Transfer Instructions */}
          {deposit.method === DepositMethod.BANK_TRANSFER && (
            <div className="space-y-2 text-sm">
              <p className="font-medium">Dados para transferência:</p>
              <div className="p-3 bg-gray-50 rounded-lg space-y-1">
                <p>
                  <span className="text-gray-600">Banco:</span> Banco LauraTech
                  (XXX)
                </p>
                <p>
                  <span className="text-gray-600">Agência:</span> 0001
                </p>
                <p>
                  <span className="text-gray-600">Conta:</span> {deposit.id}
                </p>
                <p>
                  <span className="text-gray-600">Valor:</span> R${" "}
                  {deposit.amount.toFixed(2)}
                </p>
              </div>
            </div>
          )}

          {deposit.expiresAt && (
            <p className="text-xs text-gray-500 mt-4">
              Expira em: {new Date(deposit.expiresAt).toLocaleString("pt-BR")}
            </p>
          )}

          <Button
            variant="outline"
            onClick={() => setDeposit(null)}
            className="w-full mt-4"
          >
            Novo Depósito
          </Button>
        </Card>
      )}
    </div>
  );
}
