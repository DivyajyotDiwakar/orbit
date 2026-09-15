"use client";
import { useState } from "react";
import { clientApi } from "@/lib/api-client";
import { Field } from "@/components/seller-event-form";
import StatusMessage from "@/components/status-message";
export default function VerificationForm() {
  const [pending, setPending] = useState(false);
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");
  async function submit(event) {
    event.preventDefault();
    setPending(true);
    setError("");
    const form = new FormData(event.currentTarget);
    const data = Object.fromEntries(form.entries());
    const body = {
      companyName: data.companyName,
      gstin: data.gstin,
      accountHolderName: data.accountHolderName,
      accountNumber: data.accountNumber,
      ifscCode: data.ifscCode,
      bankName: data.bankName,
    };
    try {
      const res = await clientApi("/seller/verify", {
        method: "POST",
        headers: { "content-type": "application/json" },
        body: JSON.stringify(body),
      });
      setMessage(res.message || "Verification request submitted.");
    } catch (err) {
      setError(err.message);
    } finally {
      setPending(false);
    }
  }
  return (
    <form
      onSubmit={submit}
      className="space-y-5 rounded-xl border border-border p-6"
    >
      <Field label="Company name" name="companyName" required />
      <Field label="GSTIN" name="gstin" required />
      <Field label="Account holder name" name="accountHolderName" required />
      <Field label="Account number" name="accountNumber" required />
      <div className="grid gap-4 sm:grid-cols-2">
        <Field label="IFSC code" name="ifscCode" required />
        <Field label="Bank name" name="bankName" required />
      </div>
      <StatusMessage message={message} tone="success" />
      <StatusMessage message={error} />
      <button
        disabled={pending}
        className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground disabled:opacity-50"
      >
        {pending ? "Submitting…" : "Submit verification"}
      </button>
    </form>
  );
}
