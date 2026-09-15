import SiteHeader from "@/components/site-header";
import SiteFooter from "@/components/site-footer";
import { requireUser } from "@/lib/auth";
import { backendFetch } from "@/lib/api-server";

export default async function OrdersPage() {
  const user = await requireUser("buyer");
  // TODO: Handle API errors gracefully
  const { orders } = await backendFetch("/buyer/orders");

  return (
    <div className="flex min-h-screen flex-col">
      <SiteHeader user={user} />
      <main className="mx-auto w-full max-w-3xl flex-1 px-5 py-16 sm:px-8">
        <p className="text-sm text-muted-foreground">Orders</p>
        <h1 className="mt-2 text-3xl font-semibold tracking-[-0.04em]">
          Your order history
        </h1>
        <div className="mt-10 space-y-4">
          {orders && orders.length > 0 ? (
            orders.map((order) => (
              <div
                key={order.orderId}
                className="rounded-xl border border-border bg-card p-5 transition-colors hover:bg-muted/40"
              >
                <div className="flex items-start justify-between">
                  <div>
                    <h2 className="font-medium leading-tight">
                      {order.productTitle}
                    </h2>
                    <p className="text-sm text-muted-foreground">
                      for event: {order.eventName}
                    </p>
                    <p className="mt-2 text-sm text-muted-foreground">
                      Booked on{" "}
                      {new Date(order.bookedAt).toLocaleDateString("en-US", {
                        year: "numeric",
                        month: "long",
                        day: "numeric",
                      })}
                    </p>
                  </div>
                  <div className="flex-shrink-0 text-right">
                    <p className="font-semibold">
                      ${order.price.toFixed(2)}
                    </p>
                    <p className="text-sm capitalize text-muted-foreground">
                      {order.status.toLowerCase()}
                    </p>
                  </div>
                </div>
              </div>
            ))
          ) : (
            <div className="rounded-xl border border-dashed border-border p-8 text-center text-sm text-muted-foreground">
              You have not placed any orders yet.
            </div>
          )}
        </div>
      </main>
      <SiteFooter />
    </div>
  );
}
