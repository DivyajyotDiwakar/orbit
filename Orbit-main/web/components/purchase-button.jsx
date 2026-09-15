"use client";
import { useState } from "react";
import { clientApi } from "@/lib/api-client";
import StatusMessage from "@/components/status-message";

export default function PurchaseButton({ eventId, productId, disabled }) {
  const [pending, setPending] = useState(false);
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");
  async function purchase() {
    setPending(true);
    setMessage("");
    setError("");
    try {
      const data = await clientApi(
        `/buyer/event/${eventId}/purchase/${productId}`,
        { method: "POST" },
      );
      setMessage(
        data?.reservationId
          ? `Reservation ${data.reservationId} confirmed.`
          : "Your purchase was confirmed.",
      );
    } catch (err) {
      setError(err.message);
    } finally {
      setPending(false);
    }
  }
  return (
    <div className="space-y-3">
      <button
        type="button"
        disabled={disabled || pending}
        onClick={purchase}
        className="h-10 w-full rounded-lg bg-primary px-4 text-sm font-medium text-primary-foreground disabled:cursor-not-allowed disabled:opacity-50"
      >
        {pending ? "Processing…" : disabled ? "Sold out" : "Purchase"}
      </button>
      <StatusMessage message={message} tone="success" />
      <StatusMessage message={error} />
    </div>
  );
}
