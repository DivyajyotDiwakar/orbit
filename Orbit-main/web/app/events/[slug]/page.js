import Image from "next/image";
import { notFound } from "next/navigation";
import SiteHeader from "@/components/site-header";
import SiteFooter from "@/components/site-footer";
import PurchaseButton from "@/components/purchase-button";
import { ApiError, formatCurrency } from "@/lib/api";
import { backendFetch } from "@/lib/api-server";
import { requireUser } from "@/lib/auth";

export default async function EventPage({ params }) {
  const { slug } = await params;
  let data;
  try {
    data = await backendFetch(`/buyer/event/${slug}`);
  } catch (error) {
    if (error instanceof ApiError && error.status === 404) notFound();
    throw error;
  }
  const user = await requireUser();
  const products = data?.products || [];
  return (
    <div className="flex min-h-screen flex-col">
      <SiteHeader user={user} />
      <main className="mx-auto w-full max-w-6xl flex-1 px-5 py-12 sm:px-8">
        <div className="mb-10">
          <p className="text-sm text-muted-foreground">Event</p>
          <h1 className="mt-2 text-3xl font-semibold tracking-[-0.04em]">
            Available products
          </h1>
          <p className="mt-3 max-w-2xl text-muted-foreground">
            Products currently offered in this live event.
          </p>
        </div>
        {products.length ? (
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {products.map((product) => (
              <article
                key={product.productId}
                className="overflow-hidden rounded-xl border border-border bg-card"
              >
                <div className="relative aspect-square bg-muted">
                  {product.image ? (
                    <Image
                      src={product.image}
                      alt={product.title}
                      fill
                      className="object-cover"
                      sizes="(min-width:1024px) 33vw, (min-width:768px) 50vw, 100vw"
                    />
                  ) : (
                    <div className="h-full" />
                  )}
                </div>
                <div className="space-y-4 p-5">
                  <div>
                    <h2 className="font-medium tracking-[-0.02em]">
                      {product.title}
                    </h2>
                    <p className="mt-2 text-sm leading-6 text-muted-foreground">
                      {product.description}
                    </p>
                  </div>
                  <div className="flex items-center justify-between border-t border-border pt-4">
                    <span className="font-medium">
                      {formatCurrency(product.price, product.currency)}
                    </span>
                    <span className="text-sm text-muted-foreground">
                      {product.availableStock} left
                    </span>
                  </div>
                  <PurchaseButton
                    eventId={slug}
                    productId={product.productId}
                    disabled={!product.availableStock}
                  />
                </div>
              </article>
            ))}
          </div>
        ) : (
          <div className="rounded-xl border border-dashed border-border px-6 py-14 text-center text-sm text-muted-foreground">
            No products are available for this event.
          </div>
        )}
      </main>
      <SiteFooter />
    </div>
  );
}
