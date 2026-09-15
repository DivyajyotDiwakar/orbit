"use client";
import { ApiError, getErrorMessage } from "@/lib/api";
async function responseBody(response) {
  const text = await response.text();
  if (!text) return null;
  try {
    return JSON.parse(text);
  } catch {
    return null;
  }
}
export async function clientApi(path, options = {}) {
  const response = await fetch(`/api/backend${path}`, {
    ...options,
    credentials: "include",
    headers: { ...options.headers },
  });
  const body = await responseBody(response);
  if (!response.ok)
    throw new ApiError(
      response.status,
      body?.error || body?.message || getErrorMessage(response.status),
    );
  return body;
}
