"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { clientApi } from "@/lib/api-client";
import { Field } from "@/components/seller-event-form";
import StatusMessage from "@/components/status-message";
export default function ProductEditor({ eventId, product }) {
  const router = useRouter();
  const [pending, setPending] = useState(false);
  const [error, setError] = useState("");
    async function save(event) {
        event.preventDefault();
        const formElement = event.currentTarget;
        setPending(true);
        setError("");
        try {
            await clientApi(`/seller/events/${eventId}/updateProduct/${product.id}`, {
                method: "PUT",
                body: new FormData(formElement),
      });
      router.push(`/seller/events/${eventId}`);
      router.refresh();
    } catch (err) {
      setError(err.message);
      setPending(false);
    }
  }
  async function remove() {
    if (!confirm("Delete this product?")) return;
    setPending(true);
    try {
      await clientApi(`/seller/events/${eventId}/deleteProduct/${product.id}`, {
        method: "DELETE",
      });
      router.push(`/seller/events/${eventId}`);
      router.refresh();
    } catch (err) {
      setError(err.message);
      setPending(false);
    }
  }
  return (
    <form
      onSubmit={save}
      className="space-y-5 rounded-xl border border-border p-6"
    >
      <Field label="Title" name="title" defaultValue={product.title} required />
      <Field
        label="Description"
        name="description"
        type="textarea"
        defaultValue={product.description}
        required
      />
      <div className="grid grid-cols-2 gap-4">
        <Field
          label="Price"
          name="price"
          type="number"
          step="0.01"
          min="0.01"
          defaultValue={product.price}
          required
        />
        <Field
          label="Inventory"
          name="frequency"
          type="number"
          min="1"
          defaultValue={product.frequency}
          required
        />
      </div>
      <Field label="Replace image" name="image" type="file" accept="image/*" />
      <StatusMessage message={error} />
      <div className="flex gap-3">
        <button
          disabled={pending}
          className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground"
        >
          Save changes
        </button>
        <button
          type="button"
          disabled={pending}
          onClick={remove}
          className="rounded-lg px-4 py-2 text-sm font-medium text-destructive"
        >
          Delete product
        </button>
      </div>
    </form>
  );
}
