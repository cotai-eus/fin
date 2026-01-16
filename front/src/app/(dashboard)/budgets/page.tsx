/**
 * Budgets Page
 * Manage budgets and view spending analysis
 */

import { Suspense } from "react";
import { redirect } from "next/navigation";
import { getOrySession } from "@/core/ory/session";
import { PageHeader } from "@/shared/components/PageHeader";
import { ErrorBoundary } from "@/shared/components/ErrorBoundary";
import { Skeleton } from "@/shared/components/ui/Skeleton";
import { Button } from "@/shared/components/ui/Button";
import { fetchUserBudgets } from "@/modules/budgets/actions";
import { BudgetWidget } from "@/modules/budgets/components/BudgetWidget";
import { Card } from "@/shared/components/ui/Card";

export const metadata = {
  title: "Orçamentos | LauraTech",
  description: "Gerencie seus orçamentos por categoria",
};

export const dynamic = "force-dynamic";

async function BudgetsContent() {
  const result = await fetchUserBudgets();

  if (!result.success) {
    throw new Error(result.error);
  }

  const budgets = result.data;

  return (
    <div className="space-y-6">
      {budgets.length === 0 ? (
        <Card className="p-12 text-center">
          <h3 className="text-lg font-semibold mb-2">
            Nenhum orçamento criado
          </h3>
          <p className="text-gray-600 mb-4">
            Crie orçamentos para controlar seus gastos por categoria
          </p>
          <Button variant="primary">Criar Primeiro Orçamento</Button>
        </Card>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {budgets.map((budget) => (
            <BudgetWidget key={budget.id} budget={budget} />
          ))}
        </div>
      )}
    </div>
  );
}

function BudgetsSkeleton() {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <Skeleton className="h-48" />
      <Skeleton className="h-48" />
      <Skeleton className="h-48" />
    </div>
  );
}

export default async function BudgetsPage() {
  const session = await getOrySession();
  if (!session) {
    redirect("/auth/login");
  }

  return (
    <div className="space-y-6">
      <PageHeader
        title="Orçamentos"
        description="Controle seus gastos por categoria"
        action={
          <Button variant="primary">Criar Orçamento</Button>
        }
      />

      <ErrorBoundary>
        <Suspense fallback={<BudgetsSkeleton />}>
          <BudgetsContent />
        </Suspense>
      </ErrorBoundary>
    </div>
  );
}
