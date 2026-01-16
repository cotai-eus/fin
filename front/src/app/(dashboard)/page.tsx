/**
 * Enhanced Dashboard Page
 * Shows financial overview with charts and analytics
 */

import { Suspense } from "react";
import { redirect } from "next/navigation";
import { getOrySession } from "@/core/ory/session";
import { PageHeader } from "@/shared/components/PageHeader";
import { ErrorBoundary } from "@/shared/components/ErrorBoundary";
import { Skeleton } from "@/shared/components/ui/Skeleton";
import { Card } from "@/shared/components/ui/Card";
import { getBudgetSummary, getCategorySpending, getSpendingTrends } from "@/modules/budgets/actions";
import { BudgetWidget } from "@/modules/budgets/components/BudgetWidget";
import { SpendingChart } from "@/modules/budgets/components/SpendingChart";
import { CategoryBreakdown } from "@/modules/budgets/components/CategoryBreakdown";

export const metadata = {
  title: "Dashboard | LauraTech",
  description: "Visão geral das suas finanças",
};

export const dynamic = "force-dynamic";

async function DashboardContent() {
  // Fetch all dashboard data
  const [budgetSummaryResult, categorySpendingResult, spendingTrendsResult] =
    await Promise.all([
      getBudgetSummary(),
      getCategorySpending(),
      getSpendingTrends(),
    ]);

  // Handle errors
  if (!budgetSummaryResult.success) {
    throw new Error(budgetSummaryResult.error);
  }

  const budgetSummary = budgetSummaryResult.data;
  const categorySpending = categorySpendingResult.success
    ? categorySpendingResult.data
    : [];
  const spendingTrends = spendingTrendsResult.success
    ? spendingTrendsResult.data
    : [];

  return (
    <div className="space-y-6">
      {/* Financial Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card className="p-6">
          <h3 className="text-sm text-gray-600 mb-1">Orçamento Total</h3>
          <p className="text-3xl font-bold">
            R$ {budgetSummary.totalBudget.toFixed(2)}
          </p>
        </Card>

        <Card className="p-6">
          <h3 className="text-sm text-gray-600 mb-1">Gasto Total</h3>
          <p className="text-3xl font-bold text-blue-600">
            R$ {budgetSummary.totalSpent.toFixed(2)}
          </p>
          <p className="text-xs text-gray-500 mt-1">
            {budgetSummary.percentageUsed.toFixed(1)}% utilizado
          </p>
        </Card>

        <Card className="p-6">
          <h3 className="text-sm text-gray-600 mb-1">Saldo Disponível</h3>
          <p className="text-3xl font-bold text-green-600">
            R$ {budgetSummary.remainingBudget.toFixed(2)}
          </p>
        </Card>
      </div>

      {/* Charts Section */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {spendingTrends.length > 0 && (
          <SpendingChart data={spendingTrends} />
        )}
        {categorySpending.length > 0 && (
          <CategoryBreakdown data={categorySpending} />
        )}
      </div>

      {/* Budget Widgets */}
      {budgetSummary.budgets.length > 0 && (
        <div>
          <h2 className="text-xl font-semibold mb-4">Orçamentos por Categoria</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {budgetSummary.budgets.map((budget) => (
              <BudgetWidget key={budget.id} budget={budget} />
            ))}
          </div>
        </div>
      )}

      {/* Empty State */}
      {budgetSummary.budgets.length === 0 && (
        <Card className="p-12 text-center">
          <h3 className="text-lg font-semibold mb-2">Nenhum orçamento criado</h3>
          <p className="text-gray-600 mb-4">
            Crie orçamentos por categoria para acompanhar seus gastos
          </p>
        </Card>
      )}
    </div>
  );
}

function DashboardSkeleton() {
  return (
    <div className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Skeleton className="h-32" />
        <Skeleton className="h-32" />
        <Skeleton className="h-32" />
      </div>
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <Skeleton className="h-80" />
        <Skeleton className="h-80" />
      </div>
    </div>
  );
}

export default async function DashboardPage() {
  const session = await getOrySession();
  if (!session) {
    redirect("/auth/login");
  }

  return (
    <div className="space-y-6">
      <PageHeader
        title="Dashboard"
        description="Visão geral das suas finanças"
      />

      <ErrorBoundary>
        <Suspense fallback={<DashboardSkeleton />}>
          <DashboardContent />
        </Suspense>
      </ErrorBoundary>
    </div>
  );
}
