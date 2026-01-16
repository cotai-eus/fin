/**
 * Transfers Page
 * Hub for all money transfer operations
 */

import { Suspense } from "react";
import { redirect } from "next/navigation";
import { getOrySession } from "@/core/ory/session";
import { PageHeader } from "@/shared/components/PageHeader";
import { ErrorBoundary } from "@/shared/components/ErrorBoundary";
import { Skeleton } from "@/shared/components/ui/Skeleton";
import { Card } from "@/shared/components/ui/Card";
import { PIXTransferForm } from "@/modules/transfers/components/PIXTransferForm";
import { DepositOptions } from "@/modules/transfers/components/DepositOptions";
import { PaymentRequestForm } from "@/modules/transfers/components/PaymentRequestForm";

export const metadata = {
  title: "Transferências | LauraTech",
  description: "Realize transferências PIX, TED e P2P",
};

export const dynamic = "force-dynamic";

export default async function TransfersPage() {
  const session = await getOrySession();
  if (!session) {
    redirect("/auth/login");
  }

  return (
    <div className="space-y-6">
      <PageHeader
        title="Transferências"
        description="Envie e receba dinheiro de forma rápida e segura"
      />

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* PIX Transfer */}
        <ErrorBoundary>
          <PIXTransferForm />
        </ErrorBoundary>

        {/* Deposit Money */}
        <ErrorBoundary>
          <DepositOptions />
        </ErrorBoundary>

        {/* Request Payment */}
        <ErrorBoundary>
          <PaymentRequestForm />
        </ErrorBoundary>

        {/* Quick Actions */}
        <Card className="p-6">
          <h3 className="text-lg font-semibold mb-4">Outras Opções</h3>
          <div className="space-y-2">
            <button className="w-full p-3 text-left border rounded-lg hover:bg-gray-50 transition">
              <div className="font-medium">Transferência TED</div>
              <div className="text-sm text-gray-600">
                Para contas em outros bancos (1 dia útil)
              </div>
            </button>
            <button className="w-full p-3 text-left border rounded-lg hover:bg-gray-50 transition">
              <div className="font-medium">Transferência P2P</div>
              <div className="text-sm text-gray-600">
                Para outros usuários da plataforma
              </div>
            </button>
            <button className="w-full p-3 text-left border rounded-lg hover:bg-gray-50 transition">
              <div className="font-medium">Agendar Transferência</div>
              <div className="text-sm text-gray-600">
                Programe transferências futuras
              </div>
            </button>
          </div>
        </Card>
      </div>
    </div>
  );
}
