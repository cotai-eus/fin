/**
 * Cards Page
 * Main page for card management
 */

import { Suspense } from "react";
import { redirect } from "next/navigation";
import { getOrySession } from "@/core/ory/session";
import { PageHeader } from "@/shared/components/PageHeader";
import { ErrorBoundary } from "@/shared/components/ErrorBoundary";
import { Skeleton } from "@/shared/components/ui/Skeleton";
import { fetchUserCards } from "@/modules/cards/actions";
import { CardsList } from "@/modules/cards/components/CardsList";
import { Button } from "@/shared/components/ui/Button";

export const metadata = {
  title: "Cartões | LauraTech",
  description: "Gerencie seus cartões virtuais e físicos",
};

export const dynamic = "force-dynamic";

async function CardsContent() {
  const result = await fetchUserCards();

  if (!result.success) {
    throw new Error(result.error);
  }

  return <CardsList cards={result.data} />;
}

function CardsSkeleton() {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
      <Skeleton className="h-64" />
      <Skeleton className="h-64" />
    </div>
  );
}

export default async function CardsPage() {
  const session = await getOrySession();
  if (!session) {
    redirect("/auth/login");
  }

  return (
    <div className="space-y-6">
      <PageHeader
        title="Meus Cartões"
        description="Gerencie seus cartões virtuais e físicos"
        action={
          <Button variant="primary">
            Criar Cartão Virtual
          </Button>
        }
      />

      <ErrorBoundary>
        <Suspense fallback={<CardsSkeleton />}>
          <CardsContent />
        </Suspense>
      </ErrorBoundary>
    </div>
  );
}
