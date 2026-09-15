import Link from "next/link";
import SiteHeader from "@/components/site-header";
import SiteFooter from "@/components/site-footer";
import { formatDate } from "@/lib/api";
import { backendFetch } from "@/lib/api-server";
import { requireUser } from "@/lib/auth";
import { Button } from "@/components/ui/button";

export default async function SellerPage() {
  const [user, data, verificationData] = await Promise.all([
    requireUser("seller"),
    backendFetch("/seller/events/get"),
    backendFetch("/seller/isVerified"),
  ]);
  const events = data?.events || [];
  const isVerified = verificationData?.isApproved;

  return (
    <div className="flex min-h-screen flex-col">
      <SiteHeader user={user} />
      <main className="mx-auto w-full max-w-6xl flex-1 px-5 py-12 sm:px-8">
        <div className="flex items-end justify-between gap-4">
          <div>
            <p className="text-sm text-muted-foreground">Seller workspace</p>
            <h1 className="mt-2 text-3xl font-semibold tracking-[-0.04em]">
              Events
            </h1>
          </div>
          <Button asChild disabled={!isVerified}>
            <Link
                href="/seller/events/new"
            >
                Create event
            </Link>
          </Button>
        </div>
        {!isVerified && (
            <div className="mt-4 rounded-lg border border-yellow-300 bg-yellow-50 p-4 text-sm text-yellow-800 dark:bg-yellow-900/10 dark:text-yellow-500 dark:border-yellow-900/20">
                Your account is not verified. Please{" "}
                <Link href="/seller/verify" className="font-medium underline">
                    verify your account
                </Link>{" "}
                to create events.
            </div>
        )}
        <div className="mt-9 space-y-3">
          {events.length ? (
            events.map((event) => (
              <div
                key={event.id}
                className="flex items-center justify-between rounded-xl border border-border p-5"
              >
                <div>
                  <h2 className="font-medium">{event.eventName}</h2>
                  <p className="mt-1 text-sm text-muted-foreground">
                    {formatDate(event.scheduledAt)}
                  </p>
                </div>
                <div className="flex items-center gap-4">
                  <Link
                    href={`/seller/events/${event.id}/orders`}
                    className="text-sm font-medium text-primary underline-offset-4 hover:underline"
                  >
                    View Orders
                  </Link>
                  <Link
                    href={`/seller/events/${event.id}`}
                    className="text-sm font-medium text-muted-foreground underline-offset-4 hover:underline"
                  >
                    Manage Event
                  </Link>
                </div>
              </div>
            ))
          ) : (
            <div className="rounded-xl border border-dashed border-border p-8 text-sm text-muted-foreground">
              Create your first event to start selling.
            </div>
          )}
        </div>
        <Link
          href="/seller/verify"
          className="mt-8 inline-block text-sm font-medium underline underline-offset-4"
        >
          Seller verification
        </Link>
      </main>
      <SiteFooter />
    </div>
  );
}
