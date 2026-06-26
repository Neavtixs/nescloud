import { cookies } from "next/headers";
import { NextRequest, NextResponse } from "next/server";

const API_BASE = process.env.NEXT_PUBLIC_HOST_API + "/api";

export async function POST(request: NextRequest) {
  const body = await request.text();

  const resp = await fetch(`${API_BASE}/auth/register`, {
    method: "POST",
    body,
    headers: {
      "Content-Type": request.headers.get("Content-Type") ?? "application/json",
    },
  });
  const setCookie = resp.headers.get("set-cookie");

  const accessToken = setCookie
    ? getCookieValue(setCookie, "refresh_token")
    : null;

  console.log(accessToken);

  const data = await resp.json();

  const cookieStore = await cookies();

  cookieStore.set({
    name: "access_token",
    value: data.access_token,
    httpOnly: true,
    secure: process.env.NODE_ENV === "production",
    sameSite: "lax",
    path: "/",
  });

  if (accessToken) {
    cookieStore.set({
      name: "refresh_token",
      value: accessToken,
      httpOnly: true,
      secure: process.env.NODE_ENV === "production",
      sameSite: "lax",
      path: "/",
    });
  }
  return NextResponse.json({
    user: data.user,
  });
}

function getCookieValue(cookieHeader: string, name: string) {
  const match = cookieHeader.match(new RegExp(`${name}=([^;]+)`));
  return match?.[1];
}
