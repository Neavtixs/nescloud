import { cookies } from "next/headers";
import { NextRequest, NextResponse } from "next/server";

const API_BASE = process.env.NEXT_PUBLIC_HOST_API + "/api";

type ProxyOptions = {
  method?: string;
  forwardBody?: boolean;
  setCookies?: boolean;
  clearCookies?: boolean;
};

export async function apiProxyAuth(
  request: NextRequest,
  path: string,
  options?: ProxyOptions,
) {
  const {
    method = "POST",
    forwardBody = true,
    setCookies = true,
    clearCookies = false,
  } = options ?? {};

  const body = forwardBody ? await request.text() : undefined;

  let resp: Response;
  let data: any;

  try {
    resp = await fetch(`${API_BASE}${path}`, {
    method,
    headers: {
      "Content-Type": request.headers.get("Content-Type") ?? "application/json",
    },
    body,
    });

    data = await resp.json();
    console.log(data);
  } catch {
    return NextResponse.json(
      { message: "Service unavailable" },
      { status: 502 },
    );
  }

  if (setCookies && resp.ok) {
    const cookieStore = await cookies();

    if (clearCookies) {
      cookieStore.set({
        name: "access_token",
        value: "",
        maxAge: 0,
        httpOnly: true,
        secure: process.env.NODE_ENV === "production",
        sameSite: "lax",
        path: "/",
      });

      cookieStore.set({
        name: "refresh_token",
        value: "",
        maxAge: 0,
        httpOnly: true,
        secure: process.env.NODE_ENV === "production",
        sameSite: "lax",
        path: "/",
      });
    } else {
      const setCookie = resp.headers.get("set-cookie");
      const refreshToken = setCookie
        ? getCookieValue(setCookie, "refresh_token")
        : null;

      if (data.data?.access_token) {
        cookieStore.set({
          name: "access_token",
          value: data.data.access_token,
          httpOnly: true,
          secure: process.env.NODE_ENV === "production",
          sameSite: "lax",
          path: "/",
        });
      }

      if (refreshToken) {
        cookieStore.set({
          name: "refresh_token",
          value: refreshToken,
          httpOnly: true,
          secure: process.env.NODE_ENV === "production",
          sameSite: "lax",
          path: "/",
        });
      }
    }

    if (data.data?.access_token) {
      delete data.data.access_token;
    }
  }

  return NextResponse.json(data, { status: resp.status });
}

function getCookieValue(cookieHeader: string, name: string) {
  const match = cookieHeader.match(new RegExp(`${name}=([^;]+)`));
  return match?.[1];
}
