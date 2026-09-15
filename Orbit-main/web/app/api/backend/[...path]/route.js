import { NextResponse } from "next/server";

const backendUrl = process.env.BACKEND_URL || process.env.NEXT_PUBLIC_BACKEND_URL;

function publicError(status, payload) {
  if (status >= 500) return "Orbit could not complete that request. Please try again shortly.";
  if (status === 401) return "Your session has ended. Please sign in again.";
  if (status === 403) return "You do not have permission to perform this action.";
  if (status === 404) return "The requested item could not be found.";
  if (status === 409) return payload?.error || "This action is not available in the event's current state.";
  return payload?.error || payload?.message || "Please check the submitted details and try again.";
}

async function forward(request, { params }) {
  if (!backendUrl)
    return NextResponse.json(
      { error: "Backend URL is not configured." },
      { status: 500 },
    );

  const { path } = await params;
  const url = new URL(`${backendUrl}/${path.join("/")}`);
  url.search = request.nextUrl.search;
  const headers = new Headers();
  const contentType = request.headers.get("content-type");
  const cookie = request.headers.get("cookie");
  if (contentType) headers.set("content-type", contentType);
  if (cookie) headers.set("cookie", cookie);

  try {
    const response = await fetch(url, {
      method: request.method,
      headers,
      body: ["GET", "HEAD"].includes(request.method)
        ? undefined
        : await request.arrayBuffer(),
      cache: "no-store",
    });
    if (!response.ok) {
      let payload = null;
      try {
        payload = await response.json();
      } catch {
        // The backend can return a non-JSON error. Never relay it to the browser.
      }
      return NextResponse.json(
        { error: publicError(response.status, payload) },
        { status: response.status },
      );
    }

    const result = new NextResponse(response.body, { status: response.status });
    result.headers.set("content-type", response.headers.get("content-type") || "application/json");
    return result;
  } catch {
    return NextResponse.json(
      { error: "Unable to reach the Orbit service." },
      { status: 500 },
    );
  }
}

export const GET = forward;
export const POST = forward;
export const PUT = forward;
export const DELETE = forward;
