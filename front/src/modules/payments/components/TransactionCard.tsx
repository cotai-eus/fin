/**
 * src/modules/payments/components/TransactionCard.tsx
 *
 * Transaction Card Molecule Component
 *
 * Displays a single transaction with status badge, amounts, and action buttons
 * This is a client component (uses "use client" for interactivity)
 */

"use client";

import { useState } from "react";
import { formatCurrency, formatDate } from "@/shared/utils/formatters";
import { Badge } from "@/shared/components/ui/Badge";

import { TransactionActions } from "./TransactionActions";
import { type Transaction } from "../types";

interface TransactionCardProps {
  transaction: Transaction;
  userId: string;
}

/**
 * TransactionCard renders a single transaction row/card
 * 
 * Features:
 * - Status badge (color-coded)
 * - Amount with currency formatting
 * - Date with relative time (e.g., "2 hours ago")
 * - Action buttons (export, repeat, cancel)
 * - Expandable details
 */
export function TransactionCard({
  transaction,
  userId,
}: TransactionCardProps) {
  const [isExpanded, setIsExpanded] = useState(false);
  const isOutgoing = transaction.fromUserId === userId;

  return (
    <div className="px-6 py-4 hover:bg-gray-50 transition-colors">
      {/* Main transaction row */}
      <div
        className="flex items-center justify-between cursor-pointer"
        onClick={() => setIsExpanded(!isExpanded)}
      >
        {/* Left: Icon + Description */}
        <div className="flex items-center gap-4 flex-1">
          <div
            className={`flex h-10 w-10 items-center justify-center rounded-full ${isOutgoing ? "bg-red-100" : "bg-green-100"
              }`}
          >
            <span
              className={`text-lg ${isOutgoing ? "text-red-600" : "text-green-600"}`}
            >
              {isOutgoing ? "→" : "←"}
            </span>
          </div>

          <div className="min-w-0 flex-1">
            <p className="font-semibold text-gray-900">
              {isOutgoing ? "Transfer Out" : "Transfer In"}
            </p>
            <p className="text-sm text-gray-600">
              {transaction.description || "No description"}
            </p>
          </div>
        </div>

        {/* Right: Amount + Status */}
        <div className="flex items-center gap-4">
          <div className="text-right">
            <p
              className={`font-semibold text-lg ${isOutgoing ? "text-red-600" : "text-green-600"
                }`}
            >
              {isOutgoing ? "-" : "+"}
              {formatCurrency(transaction.amount)}
            </p>
            <p className="text-xs text-gray-500">
              {formatDate(transaction.createdAt)}
            </p>
          </div>

          <Badge status={transaction.status} />

          <button className="text-gray-400 hover:text-gray-600">
            {isExpanded ? "▼" : "▶"}
          </button>
        </div>
      </div>

      {/* Expanded details */}
      {isExpanded && (
        <div className="mt-4 border-t border-gray-200 pt-4">
          <div className="grid grid-cols-2 gap-4 text-sm md:grid-cols-4">
            <div>
              <p className="text-xs font-semibold text-gray-500 uppercase">
                Transaction ID
              </p>
              <p className="mt-1 font-mono text-gray-900">{transaction.id}</p>
            </div>

            <div>
              <p className="text-xs font-semibold text-gray-500 uppercase">
                {isOutgoing ? "Recipient" : "Sender"}
              </p>
              <p className="mt-1 text-gray-900">
                {isOutgoing ? transaction.toUserId : transaction.fromUserId}
              </p>
            </div>

            <div>
              <p className="text-xs font-semibold text-gray-500 uppercase">
                Completed
              </p>
              <p className="mt-1 text-gray-900">
                {transaction.completedAt
                  ? formatDate(transaction.completedAt)
                  : "Pending"}
              </p>
            </div>


          </div>

          {/* Action buttons */}
          <div className="mt-4 flex gap-2">
            <TransactionActions transaction={transaction} userId={userId} />
          </div>
        </div>
      )}
    </div>
  );
}
