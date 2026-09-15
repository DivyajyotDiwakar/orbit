"use client";

import Link from "next/link";
import { ArrowRight, LockKeyhole, Mail } from "lucide-react";
import { useRouter } from "next/navigation";
import { useState } from "react";

export default function SignInPage() {
  const router = useRouter();
  const [pending, setPending] = useState(false);
  const [error, setError] = useState("");

  async function handleSubmit(event) {
    event.preventDefault();
    setPending(true);
    setError("");
    const form = new FormData(event.currentTarget);
    try {
      const response = await fetch("/api/auth/signin", { method: "POST", headers: { "Content-Type": "application/json" }, credentials: "include", body: JSON.stringify({ emailId: form.get("email"), password: form.get("password") }) });
      if (!response.ok) {
        const data = await response.json().catch(() => ({}));
        throw new Error(data.error || "Unable to sign in with those credentials.");
      }
      router.replace("/");
      router.refresh();
    } catch (err) {
      setError(err.message);
      setPending(false);
    }
  }

  return <main className="relative grid min-h-screen place-items-center overflow-hidden px-5 py-10"><div className="dot-grid absolute inset-0 opacity-45" /><div className="absolute left-1/2 top-0 h-90 w-140 -translate-x-1/2 rounded-full bg-primary/10 blur-3xl" /><section className="surface relative w-full max-w-md p-7 sm:p-8"><Link href="/signin" className="inline-flex items-center gap-2 text-base font-semibold tracking-[-0.04em]"><span className="grid size-8 place-items-center rounded-xl bg-primary text-sm font-bold text-primary-foreground">O</span> Orbit</Link><p className="eyebrow mt-10">Welcome back</p><h1 className="mt-3 text-3xl font-semibold tracking-[-0.05em]">Sign in to your account.</h1><p className="mt-3 text-sm leading-6 text-muted-foreground">Use the account you created for buying, selling, or reviewing sellers.</p><form onSubmit={handleSubmit} className="mt-8 space-y-5"><label className="block"><span className="mb-2 block text-sm font-medium">Email address</span><span className="flex h-11 items-center gap-3 rounded-xl border border-input bg-background px-3 focus-within:ring-2 focus-within:ring-ring"><Mail size={16} className="text-muted-foreground" /><input name="email" type="email" autoComplete="email" placeholder="you@example.com" className="w-full bg-transparent text-sm outline-none" required /></span></label><label className="block"><span className="mb-2 block text-sm font-medium">Password</span><span className="flex h-11 items-center gap-3 rounded-xl border border-input bg-background px-3 focus-within:ring-2 focus-within:ring-ring"><LockKeyhole size={16} className="text-muted-foreground" /><input name="password" type="password" autoComplete="current-password" placeholder="Your password" className="w-full bg-transparent text-sm outline-none" required /></span></label>{error && <p role="alert" className="rounded-lg bg-destructive/10 px-3 py-2.5 text-sm text-destructive">{error}</p>}<button disabled={pending} className="inline-flex h-11 w-full items-center justify-center gap-2 rounded-xl bg-primary px-4 text-sm font-medium text-primary-foreground transition hover:opacity-90 disabled:cursor-not-allowed disabled:opacity-60">{pending ? "Signing in..." : <>Sign in <ArrowRight size={16} /></>}</button></form><p className="mt-7 text-center text-sm text-muted-foreground">New to Orbit? <Link href="/signup" className="font-medium text-foreground underline underline-offset-4">Create an account</Link></p></section></main>;
}
