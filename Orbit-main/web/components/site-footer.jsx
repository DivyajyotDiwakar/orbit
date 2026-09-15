export default function SiteFooter() {
  return <footer className="border-t border-border/80"><div className="page-shell flex flex-col gap-2 py-7 text-sm text-muted-foreground sm:flex-row sm:items-center sm:justify-between"><p><span className="font-medium text-foreground">Orbit</span> <span className="mx-1">/</span> live commerce, deliberately designed.</p><p>© {new Date().getFullYear()}</p></div></footer>;
}
