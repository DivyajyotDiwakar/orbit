import Link from "next/link";
export default function NotFound() {
  return (
    <main className="mx-auto flex min-h-screen max-w-6xl items-center px-5 sm:px-8">
      <div>
        <p className="text-sm text-muted-foreground">404</p>
        <h1 className="mt-3 text-3xl font-semibold tracking-tight">
          This page does not exist.
        </h1>
        <Link
          href="/"
          className="mt-7 inline-block text-sm font-medium underline underline-offset-4"
        >
          Return to Orbit
        </Link>
      </div>
    </main>
  );
}
