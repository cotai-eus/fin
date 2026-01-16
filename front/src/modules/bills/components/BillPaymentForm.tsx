/**
 * Bill Payment Form Component
 * Manual barcode input or camera scan
 */

"use client";

import { useState } from "react";
import { Button } from "@/shared/components/ui/Button";
import { Card } from "@/shared/components/ui/Card";
import { Badge } from "@/shared/components/ui/Badge";
import { validateBarcode, payBill } from "../actions";
import { formatBarcode } from "../validators";
import type { BarcodeData, Bill } from "../types";
import { BarcodeScanner } from "./BarcodeScanner";

export function BillPaymentForm() {
  const [barcodeInput, setBarcodeInput] = useState("");
  const [showScanner, setShowScanner] = useState(false);
  const [validating, setValidating] = useState(false);
  const [paying, setPaying] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [barcodeData, setBarcodeData] = useState<BarcodeData | null>(null);
  const [paidBill, setPaidBill] = useState<Bill | null>(null);

  const handleValidateBarcode = async (barcode: string) => {
    setValidating(true);
    setError(null);
    setBarcodeData(null);

    const result = await validateBarcode(barcode);

    if (result.success) {
      setBarcodeData(result.data);
    } else {
      setError(result.error);
    }

    setValidating(false);
  };

  const handleScan = (barcode: string) => {
    setBarcodeInput(barcode);
    setShowScanner(false);
    handleValidateBarcode(barcode);
  };

  const handlePayBill = async () => {
    if (!barcodeData) return;

    setPaying(true);
    setError(null);

    const result = await payBill({
      barcode: barcodeData.barcode,
      amount: barcodeData.amount,
    });

    if (result.success) {
      setPaidBill(result.data);
      setBarcodeData(null);
      setBarcodeInput("");
    } else {
      setError(result.error);
    }

    setPaying(false);
  };

  // Show success screen
  if (paidBill) {
    return (
      <Card className="p-6">
        <div className="text-center space-y-4">
          <div className="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto">
            <svg
              className="w-8 h-8 text-green-600"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M5 13l4 4L19 7"
              />
            </svg>
          </div>

          <div>
            <h3 className="text-lg font-semibold mb-2">Pagamento Realizado!</h3>
            <p className="text-gray-600 mb-4">
              Seu boleto foi pago com sucesso
            </p>
          </div>

          <div className="p-4 bg-gray-50 rounded-lg text-left space-y-2 text-sm">
            <div className="flex justify-between">
              <span className="text-gray-600">Valor:</span>
              <span className="font-semibold">
                R$ {paidBill.finalAmount.toFixed(2)}
              </span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-600">Beneficiário:</span>
              <span className="font-medium">{paidBill.recipientName}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-600">Data:</span>
              <span>
                {new Date(paidBill.createdAt).toLocaleDateString("pt-BR")}
              </span>
            </div>
          </div>

          <Button variant="primary" onClick={() => setPaidBill(null)} className="w-full">
            Pagar Outro Boleto
          </Button>
        </div>
      </Card>
    );
  }

  // Show scanner
  if (showScanner) {
    return (
      <BarcodeScanner
        onScan={handleScan}
        onClose={() => setShowScanner(false)}
      />
    );
  }

  return (
    <div className="space-y-6">
      <Card className="p-6">
        <h3 className="text-lg font-semibold mb-4">Pagamento de Boleto</h3>

        <div className="space-y-4">
          {/* Barcode Input */}
          <div>
            <label className="block text-sm font-medium mb-2">
              Código de Barras
            </label>
            <input
              type="text"
              value={barcodeInput}
              onChange={(e) => setBarcodeInput(e.target.value)}
              placeholder="Digite ou escaneie o código de barras"
              className="w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500"
              maxLength={60}
            />
            <p className="text-xs text-gray-500 mt-1">
              Digite os números ou use a câmera para escanear
            </p>
          </div>

          {/* Action Buttons */}
          <div className="grid grid-cols-2 gap-2">
            <Button
              variant="outline"
              onClick={() => setShowScanner(true)}
              className="w-full"
            >
              📷 Escanear
            </Button>
            <Button
              onClick={() => handleValidateBarcode(barcodeInput)}
              disabled={validating || barcodeInput.length < 44}
              className="w-full"
            >
              {validating ? "Validando..." : "Validar"}
            </Button>
          </div>

          {/* Error Message */}
          {error && (
            <div className="p-3 bg-red-50 border border-red-200 rounded-lg text-red-700 text-sm">
              {error}
            </div>
          )}
        </div>
      </Card>

      {/* Bill Details (after validation) */}
      {barcodeData && (
        <Card className="p-6">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold">Detalhes do Boleto</h3>
            <Badge variant="pending">Pendente</Badge>
          </div>

          <div className="space-y-3 mb-4">
            <div>
              <p className="text-sm text-gray-600">Código de barras</p>
              <p className="font-mono text-sm">
                {formatBarcode(barcodeData.barcode)}
              </p>
            </div>

            {barcodeData.recipientName && (
              <div>
                <p className="text-sm text-gray-600">Beneficiário</p>
                <p className="font-medium">{barcodeData.recipientName}</p>
              </div>
            )}

            {barcodeData.amount && (
              <div>
                <p className="text-sm text-gray-600">Valor</p>
                <p className="text-2xl font-bold">
                  R$ {barcodeData.amount.toFixed(2)}
                </p>
              </div>
            )}

            {barcodeData.dueDate && (
              <div>
                <p className="text-sm text-gray-600">Vencimento</p>
                <p className="font-medium">
                  {new Date(barcodeData.dueDate).toLocaleDateString("pt-BR")}
                </p>
              </div>
            )}

            {barcodeData.type && (
              <div>
                <p className="text-sm text-gray-600">Tipo</p>
                <Badge variant="cancelled">{barcodeData.type}</Badge>
              </div>
            )}
          </div>

          <Button
            onClick={handlePayBill}
            disabled={paying}
            className="w-full"
          >
            {paying ? "Processando..." : "Pagar Boleto"}
          </Button>
        </Card>
      )}
    </div>
  );
}
