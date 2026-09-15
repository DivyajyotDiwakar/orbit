import Image from "next/image";
import Link from "next/link";
import { ArrowUpRight, CalendarDays } from "lucide-react";
import { formatDate } from "@/lib/api";

export default function EventCard({ event }) {
  return (
    <Link href={`/events/${event.id}`} className="surface surface-hover group block overflow-hidden">
      <div className="relative aspect-[16/10] overflow-hidden bg-muted">
        {event.imageBanner ? <Image src={event.imageBanner} alt={event.eventName || "Event banner"} fill className="object-cover transition duration-500 group-hover:scale-[1.03]" sizes="(min-width: 1024px) 33vw, (min-width: 768px) 50vw, 100vw" /> : <div className="absolute inset-0 bg-[radial-gradient(circle_at_25%_20%,color-mix(in_oklch,var(--primary)_24%,transparent),transparent_33%),linear-gradient(135deg,var(--muted),color-mix(in_oklch,var(--accent)_65%,transparent))]" />}
        <span className="absolute left-4 top-4 rounded-full border border-background/20 bg-background/80 px-2.5 py-1 text-[0.68rem] font-semibold uppercase tracking-[0.13em] text-foreground backdrop-blur">Live now</span>
      </div>
      <div className="p-5"><div className="flex items-center gap-2 text-xs text-muted-foreground"><CalendarDays size={14} /> {formatDate(event.scheduledAt, { dateStyle: "medium", timeStyle: "short" })}</div><div className="mt-3 flex items-start justify-between gap-4"><h3 className="text-lg font-semibold tracking-[-0.025em]">{event.eventName}</h3><ArrowUpRight size={18} className="mt-0.5 shrink-0 text-muted-foreground transition group-hover:text-primary" /></div>{event.description && <p className="mt-2 line-clamp-2 text-sm leading-6 text-muted-foreground">{event.description}</p>}</div>
    </Link>
  );
}
