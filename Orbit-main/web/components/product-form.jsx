"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { clientApi } from "@/lib/api-client";
import { Field } from "@/components/seller-event-form";
import StatusMessage from "@/components/status-message";
export default function ProductForm({ eventId }) {
  const router = useRouter();
  const [pending, setPending] = useState(false);
  const [error, setError] = useState("");
  async function submit(event) {
    event.preventDefault();
    const formElement = event.currentTarget;
    setPending(true);
    setError("");
    const form = new FormData(formElement);
    const image = form.get("image");
    if (!(image instanceof File) || image.size === 0) {
      setError("Please choose a product image.");
      setPending(false);
      return;
    }
    try {
      await clientApi(`/seller/events/${eventId}/registerProducts`, {
        method: "POST",
        body: form,
      });
      formElement.reset();
      router.refresh();
    } catch (err) {
      setError(err.message);
    } finally {
      setPending(false);
    }
  }
  return (
    <form
      onSubmit={submit}
      className="space-y-4 rounded-xl border border-border p-5"
    >
      <h2 className="font-medium">Add a product</h2>
      <Field label="Title" name="title" required />
      <Field label="Description" name="description" type="textarea" required />
      <div className="grid grid-cols-2 gap-4">
        <Field
          label="Price"
          name="price"
          type="number"
          min="0.01"
          step="0.01"
          required
        />
        <Field
          label="Inventory"
          name="frequency"
          type="number"
          min="1"
          required
        />
      </div>
      <Field
        label="Product image"
        name="image"
        type="file"
        accept="image/*"
        required
      />
      <StatusMessage message={error} />
      <button
        disabled={pending}
        className="rounded-lg bg-primary px-3 py-2 text-sm font-medium text-primary-foreground disabled:opacity-50"
      >
        {pending ? "Adding…" : "Add product"}
      </button>
    </form>
  );
}
