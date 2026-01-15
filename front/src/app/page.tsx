import { redirect } from "next/navigation";
import { getOrySession } from "@/core/ory/session";

export const dynamic = 'force-dynamic';

export default async function Home() {
  // Check if user is authenticated
  const session = await getOrySession();

  // If authenticated, redirect to dashboard
  if (session) {
    redirect("/dashboard/payments");
  }

  // Not authenticated: show landing page
  return (
    <div className="flex min-h-screen items-center justify-center bg-zinc-50 font-sans dark:bg-black">
      <main className="flex min-h-screen w-full max-w-3xl flex-col items-center justify-between py-32 px-16 bg-white dark:bg-black sm:items-start">
        <div className="flex flex-col items-center gap-6 text-center sm:items-start sm:text-left">
          <h1 className="max-w-lg text-4xl font-bold leading-12 tracking-tight text-black dark:text-zinc-50">
            Welcome to FinFlow
          </h1>
          <p className="max-w-md text-lg leading-8 text-zinc-600 dark:text-zinc-400">
            Manage your payments and transactions securely with Ory-powered authentication.
          </p>
          <p className="max-w-md text-sm text-zinc-500 dark:text-zinc-500">
            Sign in to view your payment history, track transactions, and manage your account.
          </p>
        </div>
        <div className="flex flex-col gap-4 text-base font-medium sm:flex-row">
          <a
            className="flex h-12 w-full items-center justify-center gap-2 rounded-full bg-black px-6 text-white transition-colors hover:bg-zinc-800 dark:bg-white dark:text-black dark:hover:bg-zinc-200 md:w-auto"
            href="/auth/login"
          >
            Sign In
          </a>
          <a
            className="flex h-12 w-full items-center justify-center rounded-full border border-solid border-black/[.08] px-6 transition-colors hover:border-transparent hover:bg-black/[.04] dark:border-white/[.145] dark:hover:bg-[#1a1a1a] md:w-auto"
            href="/auth/registration"
          >
            Create Account
          </a>
        </div>
      </main>
    </div>
  );
}
