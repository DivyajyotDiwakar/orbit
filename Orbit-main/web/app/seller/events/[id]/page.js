import Link from "next/link";
import { notFound } from "next/navigation";
import SiteHeader from "@/components/site-header";
import SellerEventActions from "@/components/seller-event-actions";
import ProductForm from "@/components/product-form";
import { ApiError, formatCurrency, formatDate } from "@/lib/api";
import { backendFetch } from "@/lib/api-server";
import { requireUser } from "@/lib/auth";
export default async function SellerEventPage({ params }) {
  const { id } = await params;
  let eventData;
  try {
    eventData = await backendFetch(`/seller/events/get/${id}`);
  } catch (err) {
    if (err instanceof ApiError && err.status === 404) notFound();
    throw err;
  }
  const [user, productsData] = await Promise.all([
    requireUser("seller"),
    backendFetch(`/seller/events/${id}/getProducts`),
  ]);
  const event = eventData.event;
  const products = productsData?.products || [];
  return (
    <div>
      <SiteHeader user={user} />
      <main className="mx-auto max-w-6xl px-5 py-12 sm:px-8">
        <Link
          href="/seller"
          className="text-sm text-muted-foreground hover:text-foreground"
        >
          ← Events
        </Link>
        <div className="mt-7 flex flex-col justify-between gap-6 border-b border-border pb-8 md:flex-row">
          <div>
            <p className="text-sm text-muted-foreground">
              {event.isLive ? "Live sale" : "Scheduled event"}
            </p>
            <h1 className="mt-2 text-3xl font-semibold tracking-[-0.04em]">
              {event.eventName}
            </h1>
            <p className="mt-3 max-w-2xl text-muted-foreground">
              {event.description}
            </p>
            <p className="mt-4 text-sm text-muted-foreground">
              {formatDate(event.scheduledAt)}
            </p>
          </div>
          <SellerEventActions eventId={id} isLive={event.isLive} />
        </div>
        <div className="mt-10 grid gap-10 lg:grid-cols-[1fr_360px]">
          <section>
            <div className="flex items-center justify-between">
              <h2 className="text-xl font-semibold tracking-tight">Products</h2>
              <Link
                href={`/seller/events/${id}/orders`}
                className="text-sm font-medium underline underline-offset-4"
              >
                Sales & orders
              </Link>
            </div>
            <div className="mt-5 space-y-3">
              {products.length ? (
                products.map((product) => (
                  <div
                    key={product.id}
                    className="flex items-center justify-between rounded-xl border border-border p-4"
                  >
                    <div>
                      <h3 className="font-medium">{product.title}</h3>
                      <p className="mt-1 text-sm text-muted-foreground">
                        {product.frequency} units ·{" "}
                        {formatCurrency(product.price)}
                      </p>
                    </div>
                  </div>
                ))
              ) : (
                <p className="rounded-xl border border-dashed border-border p-6 text-sm text-muted-foreground">
                  No products have been added.
                </p>
              )}
            </div>
          </section>
          <ProductForm eventId={id} />
        </div>
      </main>
    </div>
  );
}
