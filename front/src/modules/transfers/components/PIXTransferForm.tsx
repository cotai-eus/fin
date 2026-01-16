/**
 * PIX Transfer Form Component
 * Client component for PIX transfer input and submission
 */

"use client";

import { useState } from "react";
import { executePIXTransfer } from "../actions";
import { PIXKeyType } from "../types";
import type { PIXTransferInput } from "../validators";
import { Button } from "@/shared/components/ui/Button";
import { Card } from "@/shared/components/ui/Card";

const PIX_KEY_TYPES = [
  { value: PIXKeyType.CPF, label: "CPF" },
  { value: PIXKeyType.CNPJ, label: "CNPJ" },
  { value: PIXKeyType.EMAIL, label: "E-mail" },
  { value: PIXKeyType.PHONE, label: "Telefone" },
  { value: PIXKeyType.RANDOM, label: "Chave Aleatória" },
];

export function PIXTransferForm() {
  const [formData, setFormData] = useState<PIXTransferInput>({
    pixKey: "",
    pixKeyType: PIXKeyType.CPF,
    amount: 0,
    description: "",
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setSuccess(false);

    const result = await executePIXTransfer(formData);

    if (result.success) {
      setSuccess(true);
      setFormData({
        pixKey: "",
        pixKeyType: PIXKeyType.CPF,
        amount: 0,
        description: "",
      });
    } else {
      setError(result.error);
    }

    setLoading(false);
  };

  return (
    <Card className="p-6">
      <h3 className="text-lg font-semibold mb-4">Transferência PIX</h3>

      <form onSubmit={handleSubmit} className="space-y-4">
        {/* PIX Key Type */}
        <div>
          <label className="block text-sm font-medium mb-2">
            Tipo de Chave PIX
          </label>
          <select
            value={formData.pixKeyType}
            onChange={(e) =>
              setFormData({ ...formData, pixKeyType: e.target.value as PIXKeyType })
            }
            className="w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500"
          >
            {PIX_KEY_TYPES.map((type) => (
              <option key={type.value} value={type.value}>
                {type.label}
              </option>
            ))}
          </select>
        </div>

        {/* PIX Key */}
        <div>
          <label className="block text-sm font-medium mb-2">Chave PIX</label>
          <input
            type="text"
            value={formData.pixKey}
            onChange={(e) => setFormData({ ...formData, pixKey: e.target.value })}
            placeholder="Digite a chave PIX"
            required
            className="w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500"
          />
        </div>

        {/* Amount */}
        <div>
          <label className="block text-sm font-medium mb-2">Valor (R$)</label>
          <input
            type="number"
            step="0.01"
            min="0.01"
            max="1000000"
            value={formData.amount || ""}
            onChange={(e) =>
              setFormData({ ...formData, amount: parseFloat(e.target.value) || 0 })
            }
            placeholder="0,00"
            required
            className="w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500"
          />
        </div>

        {/* Description */}
        <div>
          <label className="block text-sm font-medium mb-2">
            Descrição (opcional)
          </label>
          <textarea
            value={formData.description}
            onChange={(e) =>
              setFormData({ ...formData, description: e.target.value })
            }
            placeholder="Ex: Pagamento de almoço"
            maxLength={500}
            rows={3}
            className="w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500"
          />
        </div>

        {/* Error Message */}
        {error && (
          <div className="p-3 bg-red-50 border border-red-200 rounded-lg text-red-700 text-sm">
            {error}
          </div>
        )}

        {/* Success Message */}
        {success && (
          <div className="p-3 bg-green-50 border border-green-200 rounded-lg text-green-700 text-sm">
            Transferência PIX realizada com sucesso!
          </div>
        )}

        {/* Submit Button */}
        <Button type="submit" disabled={loading} className="w-full">
          {loading ? "Processando..." : "Transferir"}
        </Button>
      </form>
    </Card>
  );
}
