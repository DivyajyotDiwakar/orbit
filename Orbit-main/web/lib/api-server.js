import { headers } from "next/headers";
import { ApiError, getErrorMessage } from "@/lib/api";

const backendUrl =
  process.env.BACKEND_URL || process.env.NEXT_PUBLIC_BACKEND_URL;

async function responseBody(response) {
  const text = await response.text();
  if (!text) return null;
  try {
    return JSON.parse(text);
  } catch {
    return null;
  }
}

export async function backendFetch(path, options = {}) {
  if (!backendUrl) throw new ApiError(500, "Backend URL is not configured.");
  const headerStore = await headers();
  const response = await fetch(`${backendUrl}${path}`, {
    ...options,
    headers: { cookie: headerStore.get("cookie") || "", ...options.headers },
    cache: "no-store",
  });
  const body = await responseBody(response);
  if (!response.ok)
    throw new ApiError(
      response.status,
      body?.error || body?.message || getErrorMessage(response.status),
    );
  return body;
}
export async function getCurrentUser() {
  return (await backendFetch("/auth/check")).user;
}
