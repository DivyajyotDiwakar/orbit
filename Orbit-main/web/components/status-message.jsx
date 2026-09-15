"use client";

export default function StatusMessage({ message, tone = "error" }) {
  if (!message) return null;
  return (
    <p
      role="alert"
      className={
        tone === "success"
          ? "text-sm text-foreground"
          : "text-sm text-destructive"
      }
    >
      {message}
    </p>
  );
}
