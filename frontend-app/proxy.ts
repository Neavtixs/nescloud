import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

const protectedRoutes = ["/home"];
const guestRoutes = ["/login", "/register"];

export function proxy(request: NextRequest) {
  const { pathname } = request.nextUrl;
  const accessToken = request.cookies.get("refresh_token")?.value;

  if (protectedRoutes.some((route) => pathname.startsWith(route))) {
    if (!accessToken) {
      return NextResponse.redirect(new URL("/login", request.url));
    }
  }

  if (guestRoutes.some((route) => pathname.startsWith(route))) {
    if (accessToken) {
      return NextResponse.redirect(new URL("/home", request.url));
    }
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/home/:path*", "/login/:path*", "/register/:path*"],
};
