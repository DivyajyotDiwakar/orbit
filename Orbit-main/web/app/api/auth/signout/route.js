import { NextResponse } from "next/server";

export async function POST(request) {
  const response = NextResponse.json({ ok: true });
  response.cookies.set("token", "", { httpOnly: true, sameSite: "lax", secure: request.nextUrl.protocol === "https:", path: "/", maxAge: 0 });
  return response;
}
