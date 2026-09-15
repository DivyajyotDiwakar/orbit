"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { clientApi } from "@/lib/api-client";
import StatusMessage from "@/components/status-message";
export default function VerificationActions({ id }) {
  const router = useRouter();
  const [pending, setPending] = useState("");
  const [error, setError] = useState("");
  async function approve() {
    setPending("approve");
    setError("");
    try {
      await clientApi(`/admin/verifications/${id}/approve`, { method: "POST" });
      router.push("/admin");
      router.refresh();
    } catch (err) {
      setError(err.message);
      setPending("");
    }
  }
  async function reject(event) {
    event.preventDefault();
    const reason = new FormData(event.currentTarget).get("reason");
    setPending("reject");
    setError("");
    try {
      await clientApi(`/admin/verifications/${id}/reject`, {
        method: "POST",
        headers: { "content-type": "application/json" },
        body: JSON.stringify({ reason }),
      });
      router.push("/admin");
      router.refresh();
    } catch (err) {
      setError(err.message);
      setPending("");
    }
  }
  return (
    <div className="space-y-5">
      <button
        disabled={!!pending}
        onClick={approve}
        className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground disabled:opacity-50"
      >
        {pending === "approve" ? "Approving…" : "Approve seller"}
      </button>
      <form onSubmit={reject} className="space-y-3 border-t border-border pt-5">
        <label className="block text-sm font-medium">
          Rejection reason
          <textarea
            name="reason"
            minLength="5"
            maxLength="500"
            required
            className="mt-2 min-h-24 w-full rounded-lg border border-input bg-background p-3 text-sm font-normal"
          />
        </label>
        <button
          disabled={!!pending}
          className="rounded-lg border border-border px-4 py-2 text-sm font-medium disabled:opacity-50"
        >
          {pending === "reject" ? "Rejecting…" : "Reject request"}
        </button>
      </form>
      <StatusMessage message={error} />
    </div>
  );
}
