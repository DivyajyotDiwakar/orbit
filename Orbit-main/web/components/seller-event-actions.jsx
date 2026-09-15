"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { clientApi } from "@/lib/api-client";
import StatusMessage from "@/components/status-message";
export default function SellerEventActions({ eventId, isLive }) {
  const router = useRouter();
  const [pending, setPending] = useState("");
  const [error, setError] = useState("");
  async function action(name, endpoint) {
    setPending(name);
    setError("");
    try {
      await clientApi(endpoint, { method: "POST" });
      router.refresh();
    } catch (err) {
      setError(err.message);
    } finally {
      setPending("");
    }
  }
  async function remove() {
    if (!confirm("Delete this event and its products?")) return;
    setPending("delete");
    try {
      await clientApi(`/seller/events/delete/${eventId}`, { method: "DELETE" });
      router.push("/seller");
      router.refresh();
    } catch (err) {
      setError(err.message);
      setPending("");
    }
  }
  return (
    <div className="space-y-3">
      <div className="flex flex-wrap gap-2">
        {!isLive && (
          <button
            onClick={() => action("live", `/seller/events/${eventId}/Live`)}
            disabled={!!pending}
            className="rounded-lg bg-primary px-3 py-2 text-sm font-medium text-primary-foreground"
          >
            {pending === "live" ? "Starting…" : "Go live"}
          </button>
        )}
        {isLive && (
          <>
            <button
              onClick={() => action("pause", `/seller/events/${eventId}/Pause`)}
              disabled={!!pending}
              className="rounded-lg border border-border px-3 py-2 text-sm font-medium"
            >
              {pending === "pause" ? "Pausing…" : "Pause"}
            </button>
            <button
              onClick={() =>
                action("resume", `/seller/events/${eventId}/Resume`)
              }
              disabled={!!pending}
              className="rounded-lg border border-border px-3 py-2 text-sm font-medium"
            >
              {pending === "resume" ? "Resuming…" : "Resume"}
            </button>
            <button
              onClick={() => action("end", `/seller/events/${eventId}/End`)}
              disabled={!!pending}
              className="rounded-lg border border-border px-3 py-2 text-sm font-medium"
            >
              {pending === "end" ? "Ending…" : "End sale"}
            </button>
          </>
        )}
        <button
          onClick={remove}
          disabled={!!pending}
          className="rounded-lg px-3 py-2 text-sm font-medium text-destructive"
        >
          Delete
        </button>
      </div>
      <StatusMessage message={error} />
    </div>
  );
}
