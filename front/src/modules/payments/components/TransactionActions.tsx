/**
 * src/modules/payments/components/TransactionActions.tsx
 *
 * Transaction Actions Component
 *
 * Provides action buttons for transactions (export, cancel, repeat)
 * Demonstrates calling Server Actions from client components
 */

"use client";

import { useState } from "react";
import { Button } from "@/shared/components/ui/Button";
import { exportTransaction, cancelTransfer } from "../actions";
import { type Transaction } from "../types";
import * as Sentry from "@sentry/nextjs";

interface TransactionActionsProps {
  transaction: Transaction;
  userId: string;
}

export function TransactionActions({
  transaction,
  userId,
}: TransactionActionsProps) {
  const [isLoading, setIsLoading] = useState<"export" | "cancel" | null>(null);
  const [error, setError] = useState<string | null>(null);

  const canCancel =
    transaction.status === "pending" &&
    transaction.fromUserId === userId;

  const handleExport = async (format: "pdf" | "csv") => {
    setIsLoading("export");
    setError(null);

    try {
      const result = await exportTransaction({
        transactionId: transaction.id,
        format,
      });

      if (!result.success) {
        throw new Error(result.error || "Export failed");
      }

      // Trigger download
      const link = document.createElement("a");
      link.href = result.downloadUrl || "";
      link.download = `transaction-${transaction.id}.${format}`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to export transaction";
      setError(message);
      Sentry.captureException(err, {
        tags: { action: "export_transaction" },
      });
    } finally {
      setIsLoading(null);
    }
  };

  const handleCancel = async () => {
    if (!canCancel) return;

    setIsLoading("cancel");
    setError(null);

    try {
      const result = await cancelTransfer(transaction.id);

      if (!result.success) {
        throw new Error(result.error || "Cancellation failed");
      }

      // Optionally refresh the page
      window.location.reload();
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to cancel transfer";
      setError(message);
      Sentry.captureException(err, {
        tags: { action: "cancel_transfer" },
      });
    } finally {
      setIsLoading(null);
    }
  };

  return (
    <div className="flex flex-col gap-2">
      {error && (
        <p className="text-sm text-red-600">
          {error}
        </p>
      )}

      <div className="flex gap-2 flex-wrap">
        <Button
          size="sm"
          variant="outline"
          onClick={() => handleExport("pdf")}
          disabled={isLoading === "export"}
        >
          {isLoading === "export" ? "Exporting..." : "Export PDF"}
        </Button>

        <Button
          size="sm"
          variant="outline"
          onClick={() => handleExport("csv")}
          disabled={isLoading === "export"}
        >
          Export CSV
        </Button>

        {canCancel && (
          <Button
            size="sm"
            variant="destructive"
            onClick={handleCancel}
            disabled={isLoading === "cancel"}
          >
            {isLoading === "cancel" ? "Canceling..." : "Cancel Transfer"}
          </Button>
        )}
      </div>
    </div>
  );
}
