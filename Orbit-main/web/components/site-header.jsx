"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { LogOut, Menu, X } from "lucide-react";
import { useState } from "react";
import { ThemeSwitcher } from "@/components/theme-switcher";

const linksForRole = {
  buyer: [["Discover", "/events"], ["Orders", "/orders"]],
  seller: [["Workspace", "/seller"], ["Verification", "/seller/verify"]],
  admin: [["Verifications", "/admin"]],
};

export default function SiteHeader({ user }) {
  const pathname = usePathname();
  const router = useRouter();
  const [open, setOpen] = useState(false);
  const links = linksForRole[user?.role] || [["Discover", "/events"]];

  async function signOut() {
    await fetch("/api/auth/signout", { method: "POST" });
    router.replace("/signin");
    router.refresh();
  }

  const isActive = (href) => href === "/events" ? pathname === "/" || pathname.startsWith("/events") : pathname.startsWith(href);

  return (
    <header className="sticky top-0 z-40 border-b border-border/80 bg-background/80 backdrop-blur-xl">
      <div className="page-shell flex h-16 items-center justify-between gap-5">
        <Link href="/" className="inline-flex items-center gap-2.5" aria-label="Orbit home">
          <span className="grid size-8 place-items-center rounded-xl bg-primary text-sm font-bold tracking-[-0.08em] text-primary-foreground shadow-sm">O</span>
          <span className="text-base font-semibold tracking-[-0.045em]">Orbit</span>
        </Link>
        <nav className="hidden items-center gap-1 md:flex" aria-label="Primary navigation">
          {links.map(([label, href]) => <Link key={href} href={href} className={`rounded-lg px-3 py-2 text-sm transition ${isActive(href) ? "bg-muted font-medium text-foreground" : "text-muted-foreground hover:bg-muted hover:text-foreground"}`}>{label}</Link>)}
        </nav>
        <div className="hidden items-center gap-2 md:flex">
          <ThemeSwitcher />
          <button onClick={signOut} className="inline-flex h-9 items-center gap-2 rounded-lg px-3 text-sm text-muted-foreground transition hover:bg-muted hover:text-foreground"><LogOut size={15} /> Sign out</button>
        </div>
        <div className="flex items-center gap-1 md:hidden">
          <ThemeSwitcher />
          <button onClick={() => setOpen((value) => !value)} className="grid size-9 place-items-center rounded-lg hover:bg-muted" aria-expanded={open} aria-label="Toggle navigation">{open ? <X size={18} /> : <Menu size={18} />}</button>
        </div>
      </div>
      {open && <div className="border-t border-border bg-background px-5 py-3 md:hidden"><nav className="mx-auto flex max-w-7xl flex-col gap-1" aria-label="Mobile navigation">{links.map(([label, href]) => <Link key={href} href={href} onClick={() => setOpen(false)} className="rounded-lg px-3 py-2.5 text-sm hover:bg-muted">{label}</Link>)}<button onClick={signOut} className="mt-1 inline-flex items-center gap-2 rounded-lg px-3 py-2.5 text-left text-sm text-muted-foreground hover:bg-muted"><LogOut size={15} /> Sign out</button></nav></div>}
    </header>
  );
}
