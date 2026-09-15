import Link from "next/link";
import SiteHeader from "@/components/site-header";
import SiteFooter from "@/components/site-footer";
import { formatDate } from "@/lib/api";
import { backendFetch } from "@/lib/api-server";
import { requireUser } from "@/lib/auth";
export default async function AdminPage() {
  const [user, data] = await Promise.all([
    requireUser("admin"),
    backendFetch("/admin/verifications/pending"),
  ]);
  const verifications = data?.verifications || [];
  return (
    <div className="flex min-h-screen flex-col">
      <SiteHeader user={user} />
      <main className="mx-auto w-full max-w-5xl flex-1 px-5 py-12 sm:px-8">
        <p className="text-sm text-muted-foreground">Admin</p>
        <h1 className="mt-2 text-3xl font-semibold tracking-[-0.04em]">
          Seller verifications
        </h1>
        <div className="mt-9 overflow-x-auto rounded-xl border border-border">
          <table className="w-full min-w-150 text-left text-sm">
            <thead className="border-b border-border text-muted-foreground">
              <tr>
                <th className="px-5 py-3 font-medium">Company</th>
                <th className="px-5 py-3 font-medium">GSTIN</th>
                <th className="px-5 py-3 font-medium">Submitted</th>
                <th className="px-5 py-3" />
              </tr>
            </thead>
            <tbody>
              {verifications.length ? (
                verifications.map((item) => (
                  <tr
                    key={item.id}
                    className="border-b border-border last:border-0"
                  >
                    <td className="px-5 py-4 font-medium">
                      {item.companyName}
                    </td>
                    <td className="px-5 py-4 text-muted-foreground">
                      {item.gstin}
                    </td>
                    <td className="px-5 py-4 text-muted-foreground">
                      {formatDate(item.submittedAt || item.createdAt)}
                    </td>
                    <td className="px-5 py-4 text-right">
                      <Link
                        href={`/admin/verifications/${item.id}`}
                        className="font-medium underline underline-offset-4"
                      >
                        Review
                      </Link>
                    </td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td
                    colSpan="4"
                    className="px-5 py-10 text-center text-muted-foreground"
                  >
                    There are no pending verification requests.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </main>
      <SiteFooter />
    </div>
  );
}
