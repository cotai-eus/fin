/**
 * Support Page
 * Central de ajuda com FAQ, chat ao vivo e tickets
 */

import { Suspense } from "react";
import { redirect } from "next/navigation";
import { getOrySession } from "@/core/ory/session";
import { PageHeader } from "@/shared/components/PageHeader";
import { ErrorBoundary } from "@/shared/components/ErrorBoundary";
import { Skeleton } from "@/shared/components/ui/Skeleton";
import { getFAQCategories, fetchUserTickets } from "@/modules/support/actions";
import { FAQSection } from "@/modules/support/components/FAQSection";
import { LiveChat } from "@/modules/support/components/LiveChat";
import { TicketHistory } from "@/modules/support/components/TicketHistory";

export const metadata = {
  title: "Suporte | LauraTech",
  description: "Central de ajuda e suporte ao cliente",
};

export const dynamic = "force-dynamic";

async function SupportContent() {
  const [faqResult, ticketsResult] = await Promise.all([
    getFAQCategories(),
    fetchUserTickets(),
  ]);

  if (!faqResult.success) {
    throw new Error(faqResult.error);
  }

  const faqCategories = faqResult.data;
  const tickets = ticketsResult.success ? ticketsResult.data : [];

  return (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
      {/* Left Column: FAQ + Tickets */}
      <div className="lg:col-span-2 space-y-6">
        {/* FAQ Section */}
        <div>
          <h2 className="text-xl font-semibold mb-4">
            Perguntas Frequentes
          </h2>
          <FAQSection categories={faqCategories} />
        </div>

        {/* Ticket History */}
        <div>
          <h2 className="text-xl font-semibold mb-4">Meus Tickets</h2>
          <TicketHistory tickets={tickets} />
        </div>
      </div>

      {/* Right Column: Live Chat */}
      <div className="lg:col-span-1">
        <LiveChat />
      </div>
    </div>
  );
}

function SupportSkeleton() {
  return (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <div className="lg:col-span-2 space-y-6">
        <Skeleton className="h-96" />
        <Skeleton className="h-64" />
      </div>
      <div className="lg:col-span-1">
        <Skeleton className="h-[600px]" />
      </div>
    </div>
  );
}

export default async function SupportPage() {
  const session = await getOrySession();
  if (!session) {
    redirect("/auth/login");
  }

  return (
    <div className="space-y-6">
      <PageHeader
        title="Central de Suporte"
        description="Encontre respostas ou fale conosco"
      />

      <ErrorBoundary>
        <Suspense fallback={<SupportSkeleton />}>
          <SupportContent />
        </Suspense>
      </ErrorBoundary>
    </div>
  );
}
