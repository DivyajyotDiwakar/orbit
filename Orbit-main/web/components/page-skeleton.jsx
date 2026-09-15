export default function PageSkeleton({ rows = 3 }) {
  return (
    <div className="mx-auto max-w-6xl animate-pulse space-y-5 px-5 py-12 sm:px-8">
      <div className="h-8 w-48 rounded bg-muted" />
      {Array.from({ length: rows }).map((_, index) => (
        <div
          className="h-28 rounded-xl border border-border bg-muted/50"
          key={index}
        />
      ))}
    </div>
  );
}
