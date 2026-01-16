/**
 * Bills Payment Page
 * Pay utility bills using barcode scanner
 */

import { Suspense } from "react";
import { redirect } from "next/navigation";
import { getOrySession } from "@/core/ory/session";
import { PageHeader } from "@/shared/components/PageHeader";
import { ErrorBoundary } from "@/shared/components/ErrorBoundary";
import { Skeleton } from "@/shared/components/ui/Skeleton";
import { BillPaymentForm } from "@/modules/bills/components/BillPaymentForm";

export const metadata = {
  title: "Pagamento de Contas | LauraTech",
  description: "Pague suas contas com código de barras",
};

export const dynamic = "force-dynamic";

export default async function BillsPage() {
  const session = await getOrySession();
  if (!session) {
    redirect("/auth/login");
  }

  return (
    <div className="space-y-6">
      <PageHeader
        title="Pagamento de Contas"
        description="Pague boletos usando código de barras ou câmera"
      />

      <div className="max-w-2xl mx-auto">
        <ErrorBoundary>
          <Suspense fallback={<Skeleton className="h-96" />}>
            <BillPaymentForm />
          </Suspense>
        </ErrorBoundary>
      </div>
    </div>
  );
}
