"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import StatusMessage from "@/components/status-message";
import { clientApi } from "@/lib/api-client";

export default function SellerEventForm() {
  const router = useRouter();
  const [pending, setPending] = useState(false);
  const [error, setError] = useState("");
  async function submit(event) {
    event.preventDefault();
    const formElement = event.currentTarget;
    setPending(true);
    setError("");
    const form = new FormData(formElement);
    try {
      const result = await clientApi("/seller/events/create", {
        method: "POST",
        body: form,
      });
      if (!result?.eventId) {
        throw new Error(
          "The event was created, but Orbit did not receive its identifier.",
        );
      }
      router.push(`/seller/events/${result.eventId}`);
      router.refresh();
    } catch (err) {
      setError(err.message);
      setPending(false);
    }
  }
  return (
    <form
      onSubmit={submit}
      className="space-y-5 rounded-xl border border-border p-6"
    >
      <Field label="Event name" name="title" required />
      <Field label="Description" name="description" type="textarea" required />
      <Field
        label="Scheduled at"
        name="scheduledAt"
        type="datetime-local"
        required
      />
      <Field
        label="Banner image"
        name="imageBanner"
        type="file"
        accept="image/*"
        required
      />
      <StatusMessage message={error} />
      <button
        disabled={pending}
        className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground disabled:opacity-50"
      >
        {pending ? "Creating…" : "Create event"}
      </button>
    </form>
  );
}
export function Field({ label, name, type = "text", ...props }) {
  return (
    <label className="block space-y-2 text-sm font-medium">
      {label}
      {type === "textarea" ? (
        <textarea
          name={name}
          className="min-h-28 w-full rounded-lg border border-input bg-background px-3 py-2 text-sm font-normal outline-none focus:ring-2 focus:ring-ring"
          {...props}
        />
      ) : (
        <input
          name={name}
          type={type}
          className="h-10 w-full rounded-lg border border-input bg-background px-3 text-sm font-normal outline-none focus:ring-2 focus:ring-ring file:mr-3 file:border-0 file:bg-transparent file:text-sm file:font-medium"
          {...props}
        />
      )}
    </label>
  );
}
