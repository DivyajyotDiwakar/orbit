import { NextResponse } from "next/server";

export async function GET(request) {
  const cookie = request.headers.get("cookie") || "";

  const res = await fetch(
    `${process.env.BACKEND_URL || process.env.NEXT_PUBLIC_BACKEND_URL}/auth/check`,
    {
      method: "GET",
      headers: {
        cookie,
      },
      credentials: "include",
    }
  );

  const body = await res.text();

  return new NextResponse(body, {
    status: res.status,
    headers: {
      "Content-Type": "application/json",
    },
  });
}
