import Link from "next/link";
import { SignupForm } from "@/components/signup-form";

export default function SignupPage() {
  return <main className="relative grid min-h-screen place-items-center overflow-hidden px-5 py-10"><div className="dot-grid absolute inset-0 opacity-45" /><div className="absolute right-1/4 top-1/4 size-80 rounded-full bg-primary/10 blur-3xl" /><section className="surface relative w-full max-w-xl p-7 sm:p-8"><Link href="/signup" className="inline-flex items-center gap-2 text-base font-semibold tracking-[-0.04em]"><span className="grid size-8 place-items-center rounded-xl bg-primary text-sm font-bold text-primary-foreground">O</span> Orbit</Link><div className="mt-9"><p className="eyebrow">Create your account</p><h1 className="mt-3 text-3xl font-semibold tracking-[-0.05em]">Join the next live event.</h1><p className="mt-3 text-sm leading-6 text-muted-foreground">Choose the role that matches how you want to use Orbit. You can create a buyer or seller account.</p></div><SignupForm className="mt-8 border-0 bg-transparent p-0 shadow-none" /></section></main>;
}
