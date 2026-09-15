import { NextResponse } from "next/server";

const BACKEND = process.env.BACKEND_URL || process.env.NEXT_PUBLIC_BACKEND_URL;

export async function GET(request) {
  try {
    const cookie = request.headers.get("cookie") || "";

    const res = await fetch(`${BACKEND}/buyer/events`, {
      method: "GET",
      headers: {
        cookie,
      },
      cache: "no-store",
    });

    const body = await res.text();

    return new NextResponse(body, {
      status: res.status,
      headers: {
        "Content-Type": "application/json",
      },
    });
  } catch (err) {
    console.error(err);

    return NextResponse.json(
      {
        error: "Failed to fetch events",
      },
      {
        status: 500,
      },
    );
  }
}
