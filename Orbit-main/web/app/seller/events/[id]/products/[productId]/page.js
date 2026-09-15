import Link from "next/link";
import { notFound } from "next/navigation";
import SiteHeader from "@/components/site-header";
import ProductEditor from "@/components/product-editor";
import { ApiError } from "@/lib/api";
import { backendFetch } from "@/lib/api-server";
import { requireUser } from "@/lib/auth";
export default async function ProductPage({ params }) {
  const { id, productId } = await params;
  let data;
  try {
    data = await backendFetch(`/seller/events/${id}/getProduct/${productId}`);
  } catch (err) {
    if (err instanceof ApiError && err.status === 404) notFound();
    throw err;
  }
  const user = await requireUser("seller");
  return (
    <div>
      <SiteHeader user={user} />
      <main className="mx-auto max-w-2xl px-5 py-12 sm:px-8">
        <Link
          href={`/seller/events/${id}`}
          className="text-sm text-muted-foreground"
        >
          ← Event
        </Link>
        <h1 className="mt-7 text-3xl font-semibold tracking-[-0.04em]">
          Edit product
        </h1>
        <div className="mt-8">
          <ProductEditor eventId={id} product={data.product} />
        </div>
      </main>
    </div>
  );
}
