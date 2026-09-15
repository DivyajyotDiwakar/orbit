import { NextResponse } from "next/server";

export async function POST(request) {
  const body = await request.json();

  const res = await fetch(
    `${process.env.BACKEND_URL || process.env.NEXT_PUBLIC_BACKEND_URL}/auth/signin`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }
  );

  const data = await res.text();

  const response = new NextResponse(data, {
    status: res.status,
    headers: {
      "Content-Type": "application/json",
    },
  });

  const cookie = res.headers.get("set-cookie");

  if (cookie) {
    response.headers.set("set-cookie", cookie);
  }

  return response;
}
