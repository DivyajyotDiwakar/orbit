import Link from "next/link";
import { notFound } from "next/navigation";
import SiteHeader from "@/components/site-header";
import { ApiError } from "@/lib/api";
import { backendFetch } from "@/lib/api-server";
import { requireUser } from "@/lib/auth";

export default async function SellerEventOrdersPage({ params }) {
  const { id } = params;
  let eventData;
  try {
    eventData = await backendFetch(`/seller/events/get/${id}`);
  } catch (err) {
    if (err instanceof ApiError && err.status === 404) notFound();
    throw err;
  }

  const [user, ordersData] = await Promise.all([
    requireUser("seller"),
    backendFetch(`/seller/events/booking/fetch/${id}`),
  ]);

  const event = eventData.event;
  const orders = ordersData?.orders || [];

  return (
    <div>
      <SiteHeader user={user} />
      <main className="mx-auto max-w-6xl px-5 py-12 sm:px-8">
        <Link
          href={`/seller/events/${id}`}
          className="text-sm text-muted-foreground hover:text-foreground"
        >
          ← Back to event
        </Link>
        <div className="mt-7">
          <p className="text-sm text-muted-foreground">Orders for</p>
          <h1 className="mt-2 text-3xl font-semibold tracking-[-0.04em]">
            {event.eventName}
          </h1>
        </div>

        <div className="mt-10">
          <h2 className="text-xl font-semibold tracking-tight">
            All orders
          </h2>
          {orders.length > 0 ? (
            <div className="mt-5 overflow-x-auto rounded-lg border border-border">
              <table className="min-w-full divide-y divide-border">
                <thead className="bg-muted/40">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground">
                      Order ID
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground">
                      Customer
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground">
                      Product
                    </th>
                    <th className="px-6 py-3 text-right text-xs font-medium uppercase tracking-wider text-muted-foreground">
                      Amount
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground">
                      Status
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground">
                      Date
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border">
                  {orders.map((order) => (
                    <tr key={order.orderId}>
                      <td className="whitespace-nowrap px-6 py-4 text-sm font-medium text-foreground">
                        #{order.orderId.substring(0, 7)}...
                      </td>
                      <td className="whitespace-nowrap px-6 py-4 text-sm text-muted-foreground">
                        <div className="font-medium text-foreground">{order.customerName}</div>
                        <div>{order.customerEmail}</div>
                      </td>
                      <td className="whitespace-nowrap px-6 py-4 text-sm text-muted-foreground">
                        {order.productTitle}
                      </td>
                      <td className="whitespace-nowrap px-6 py-4 text-right text-sm font-medium text-foreground">
                        ${order.price.toFixed(2)}
                      </td>
                      <td className="whitespace-nowrap px-6 py-4 text-sm capitalize text-muted-foreground">
                        {order.status.toLowerCase()}
                      </td>
                      <td className="whitespace-nowrap px-6 py-4 text-sm text-muted-foreground">
                        {new Date(order.bookedAt).toLocaleDateString("en-US", {
                          year: "numeric",
                          month: "short",
                          day: "numeric",
                        })}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ) : (
            <div className="mt-6 rounded-xl border border-dashed border-border p-8 text-center text-sm text-muted-foreground">
              <p>No orders have been placed for this event yet.</p>
            </div>
          )}
        </div>
      </main>
    </div>
  );
}
