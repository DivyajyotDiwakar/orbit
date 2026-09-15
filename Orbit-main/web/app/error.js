"use client";
export default function ErrorPage({ reset }) {
  return (
    <main className="mx-auto flex min-h-screen max-w-6xl items-center px-5 sm:px-8">
      <div>
        <p className="text-sm text-muted-foreground">Orbit</p>
        <h1 className="mt-3 text-3xl font-semibold tracking-tight">
          We couldn’t load this page.
        </h1>
        <p className="mt-3 text-muted-foreground">
          Please try again in a moment.
        </p>
        <button
          onClick={reset}
          className="mt-7 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground"
        >
          Try again
        </button>
      </div>
    </main>
  );
}
