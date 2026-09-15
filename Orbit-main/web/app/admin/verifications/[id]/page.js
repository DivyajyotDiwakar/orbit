import Link from "next/link";
import { notFound } from "next/navigation";
import SiteHeader from "@/components/site-header";
import VerificationActions from "@/components/verification-actions";
import { ApiError, formatDate } from "@/lib/api";
import { backendFetch } from "@/lib/api-server";
import { requireUser } from "@/lib/auth";
export default async function VerificationPage({ params }) {
  const { id } = await params;
  let data;
  try {
    data = await backendFetch(`/admin/verifications/${id}`);
  } catch (err) {
    if (err instanceof ApiError && err.status === 404) notFound();
    throw err;
  }
  const user = await requireUser("admin");
  const verification = data.verification;
  const fields = [
    ["Company", verification.companyName],
    ["GSTIN", verification.gstin],
    ["Account holder", verification.accountHolderName],
    ["Account number", verification.accountNumber],
    ["IFSC code", verification.ifscCode],
    ["Bank", verification.bankName],
    [
      "Submitted",
      formatDate(verification.createdAt || verification.submittedAt),
    ],
  ];
  return (
    <div>
      <SiteHeader user={user} />
      <main className="mx-auto max-w-3xl px-5 py-12 sm:px-8">
        <Link href="/admin" className="text-sm text-muted-foreground">
          ← Verifications
        </Link>
        <h1 className="mt-7 text-3xl font-semibold tracking-[-0.04em]">
          Review seller
        </h1>
        <div className="mt-8 grid gap-8 rounded-xl border border-border p-6 md:grid-cols-[1fr_220px]">
          <dl className="space-y-5">
            {fields.map(([label, value]) => (
              <div key={label}>
                <dt className="text-sm text-muted-foreground">{label}</dt>
                <dd className="mt-1 text-sm font-medium">{value || "—"}</dd>
              </div>
            ))}
          </dl>
          <VerificationActions id={id} />
        </div>
      </main>
    </div>
  );
}
