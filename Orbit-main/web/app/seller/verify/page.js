import Link from "next/link";
import SiteHeader from "@/components/site-header";
import VerificationForm from "@/components/verification-form";
import { requireUser } from "@/lib/auth";
export default async function VerifySellerPage() {
  const user = await requireUser("seller");
  return (
    <div>
      <SiteHeader user={user} />
      <main className="mx-auto max-w-2xl px-5 py-12 sm:px-8">
        <Link href="/seller" className="text-sm text-muted-foreground">
          ← Seller workspace
        </Link>
        <h1 className="mt-7 text-3xl font-semibold tracking-[-0.04em]">
          Verify your seller account
        </h1>
        <p className="mt-2 text-sm leading-6 text-muted-foreground">
          Provide the business and bank details required for review.
        </p>
        <div className="mt-8">
          <VerificationForm />
        </div>
      </main>
    </div>
  );
}
