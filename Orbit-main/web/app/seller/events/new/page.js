import Link from "next/link";
import SiteHeader from "@/components/site-header";
import SellerEventForm from "@/components/seller-event-form";
import { requireUser } from "@/lib/auth";
export default async function NewEventPage() {
  const user = await requireUser("seller");
  return (
    <div>
      <SiteHeader user={user} />
      <main className="mx-auto max-w-2xl px-5 py-12 sm:px-8">
        <Link
          href="/seller"
          className="text-sm text-muted-foreground hover:text-foreground"
        >
          ← Events
        </Link>
        <h1 className="mt-7 text-3xl font-semibold tracking-[-0.04em]">
          Create event
        </h1>
        <p className="mt-2 text-sm text-muted-foreground">
          Set the details and banner for a future sale.
        </p>
        <div className="mt-8">
          <SellerEventForm />
        </div>
      </main>
    </div>
  );
}
