import Link from "next/link";
import { ArrowRight, CalendarDays, ShieldCheck, Zap } from "lucide-react";
import EventCard from "@/components/event-card";
import SiteFooter from "@/components/site-footer";
import SiteHeader from "@/components/site-header";
import { backendFetch } from "@/lib/api-server";
import { requireUser } from "@/lib/auth";

export default async function HomePage() {
  const [user, data] = await Promise.all([requireUser(), backendFetch("/buyer/events")]);
  const events = data?.events || [];
  const name = user.firstName || user.emailId?.split("@")[0] || "there";

  return (
    <div className="flex min-h-screen flex-col overflow-hidden">
      <SiteHeader user={user} />
      <main className="flex-1">
        <section className="relative border-b border-border/80">
          <div className="dot-grid pointer-events-none absolute inset-0 opacity-45 [mask-image:linear-gradient(to_bottom,black,transparent)]" />
          <div className="page-shell relative grid gap-12 py-18 sm:py-24 lg:grid-cols-[1.1fr_0.9fr] lg:items-end lg:py-28">
            <div>
              <p className="eyebrow">Live commerce platform</p>
              <h1 className="mt-5 max-w-3xl text-4xl font-semibold tracking-[-0.065em] text-balance sm:text-6xl lg:text-7xl">Good to see you, {name}.</h1>
              <p className="mt-6 max-w-xl text-base leading-7 text-muted-foreground sm:text-lg">Discover active events, see live inventory, and reserve products in a focused checkout flow.</p>
              <div className="mt-9 flex flex-wrap gap-3">
                <Link href="/events" className="inline-flex h-11 items-center gap-2 rounded-xl bg-primary px-4 text-sm font-medium text-primary-foreground shadow-sm transition hover:opacity-90">Explore events <ArrowRight size={16} /></Link>
                {user.role === "seller" && <Link href="/seller" className="inline-flex h-11 items-center rounded-xl border border-border bg-card px-4 text-sm font-medium transition hover:bg-muted">Open workspace</Link>}
              </div>
            </div>
            <aside className="surface relative overflow-hidden p-6 sm:p-7">
              <div className="absolute right-0 top-0 size-42 rounded-full bg-primary/10 blur-3xl" />
              <p className="eyebrow">How Orbit works</p>
              <ol className="relative mt-6 space-y-5">
                {[[CalendarDays, "Events are scoped", "Each product belongs to a specific live event."], [Zap, "Reservations are immediate", "Live stock is decided through the booking engine."], [ShieldCheck, "Confirmation is durable", "Completed sales are synchronized into orders."]].map(([Icon, title, copy], index) => <li key={title} className="flex gap-4"><span className="grid size-9 shrink-0 place-items-center rounded-xl bg-muted text-primary"><Icon size={17} /></span><div><p className="text-sm font-medium"><span className="mr-2 text-muted-foreground">0{index + 1}</span>{title}</p><p className="mt-1 text-sm leading-5 text-muted-foreground">{copy}</p></div></li>)}
              </ol>
            </aside>
          </div>
        </section>

        <section className="page-shell py-14 sm:py-18">
          <div className="flex flex-wrap items-end justify-between gap-5">
            <div><p className="eyebrow">Available now</p><h2 className="mt-3 text-2xl font-semibold tracking-[-0.045em] sm:text-3xl">Live events</h2></div>
            <Link href="/events" className="text-sm font-medium text-muted-foreground transition hover:text-foreground">View all events</Link>
          </div>
          {events.length ? <div className="mt-8 grid gap-5 md:grid-cols-2 lg:grid-cols-3">{events.slice(0, 6).map((event) => <EventCard key={event.id} event={event} />)}</div> : <div className="surface mt-8 px-6 py-14 text-center"><p className="text-base font-medium">No live events right now</p><p className="mt-2 text-sm text-muted-foreground">New events will appear here when a seller opens a sale.</p></div>}
        </section>
      </main>
      <SiteFooter />
    </div>
  );
}
