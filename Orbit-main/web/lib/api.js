export class ApiError extends Error {
  constructor(status, message) {
    super(message);
    this.status = status;
  }
}

export function getErrorMessage(
  status,
  fallback = "Something went wrong. Please try again.",
) {
  if (status === 401) return "Your session has ended. Please sign in again.";
  if (status === 403) return "You do not have permission to access this area.";
  if (status === 404) return "The requested item could not be found.";
  if (status === 409)
    return "This action cannot be completed because the item has changed.";
  if (status >= 500)
    return "The service is temporarily unavailable. Please try again.";
  return fallback;
}

async function responseBody(response) {
  const text = await response.text();
  if (!text) return null;
  try {
    return JSON.parse(text);
  } catch {
    return null;
  }
}

export function formatCurrency(value, currency = "INR") {
  const normalizedCurrency = /^[A-Z]{3}$/.test(String(currency).toUpperCase())
    ? String(currency).toUpperCase()
    : "INR";

  return new Intl.NumberFormat("en-IN", {
    style: "currency",
    currency: normalizedCurrency,
  }).format(Number(value) || 0);
}

export function formatDate(
  value,
  options = { dateStyle: "medium", timeStyle: "short" },
) {
  if (!value) return "—";
  return new Intl.DateTimeFormat("en-IN", options).format(new Date(value));
}
