import { NextResponse } from "next/server";

const publicRoutes = ["/signin", "/signup"];

export async function proxy(request) {
  const { pathname } = request.nextUrl;

  if (
    pathname.startsWith("/_next") ||
    pathname.startsWith("/favicon.ico") ||
    pathname.startsWith("/api")
  ) {
    return NextResponse.next();
  }

  const isPublic = publicRoutes.includes(pathname);

  const cookie = request.headers.get("cookie") || "";

  let authenticated = false;
  try {
    const res = await fetch(`${request.nextUrl.origin}/api/auth/check`, { headers: { cookie }, cache: "no-store" });
    authenticated = res.ok;
  } catch {
    authenticated = false;
  }

  if (!authenticated && !isPublic) {
    return NextResponse.redirect(new URL("/signin", request.url));
  }

  if (authenticated && isPublic) {
    return NextResponse.redirect(new URL("/", request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/:path*"],
};
