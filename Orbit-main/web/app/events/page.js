import EventCard from "@/components/event-card";
import SiteFooter from "@/components/site-footer";
import SiteHeader from "@/components/site-header";
import { backendFetch } from "@/lib/api-server";
import { requireUser } from "@/lib/auth";

export default async function EventsPage() {
  const [user, data] = await Promise.all([requireUser(), backendFetch("/buyer/events")]);
  const events = data?.events || [];
  return <div className="flex min-h-screen flex-col"><SiteHeader user={user} /><main className="page-shell flex-1 py-12 sm:py-16"><p className="eyebrow">Discover</p><div className="mt-3 max-w-2xl"><h1 className="text-4xl font-semibold tracking-[-0.055em] sm:text-5xl">Events with a defined moment.</h1><p className="mt-4 text-base leading-7 text-muted-foreground">Every product is available only within its live event. Inventory shown here comes from the active sale state.</p></div>{events.length ? <div className="mt-10 grid gap-5 md:grid-cols-2 lg:grid-cols-3">{events.map((event) => <EventCard key={event.id} event={event} />)}</div> : <div className="surface mt-10 px-6 py-16 text-center"><p className="font-medium">Nothing is live at the moment.</p><p className="mt-2 text-sm text-muted-foreground">Check back when the next seller event opens.</p></div>}</main><SiteFooter /></div>;
}
