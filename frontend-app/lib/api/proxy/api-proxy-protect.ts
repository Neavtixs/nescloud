import { NextRequest, NextResponse } from "next/server";

const API_BASE = process.env.NEXT_PUBLIC_HOST_API + "/api";

async function fetchBackend(path: string, token?: string, init?: RequestInit) {
  const res = await fetch(`${API_BASE}${path}`, {
    headers: {
      "Content-Type": "application/json",
      ...(token && { Authorization: `Bearer ${token}` }),
    },
    ...init,
  });
  const data = await res.json();
  return { res, data };
}

function setAccessToken(response: NextResponse, value: string) {
  response.cookies.set("access_token", value, {
    httpOnly: true,
    sameSite: "lax",
    path: "/",
  });
}

function setRefreshToken(response: NextResponse, value: string) {
  response.cookies.set("refresh_token", value, {
    maxAge: 604800,
    httpOnly: true,
    sameSite: "lax",
    path: "/",
  });
}

function clearCookie(response: NextResponse, name: string) {
  response.cookies.set(name, "", {
    maxAge: -1,
    httpOnly: true,
    sameSite: "lax",
    path: "/",
  });
}

export async function apiProxyProtect(
  request: NextRequest,
  path: string,
  options?: { method?: string; body?: string },
) {
  const accessToken = request.cookies.get("access_token")?.value;
  const { method, body } = options ?? {};

  let { res, data } = await fetchBackend(path, accessToken, { method, body });

  if (res.status === 401) {
    const refreshToken = request.cookies.get("refresh_token")?.value;

    if (refreshToken) {
      const refreshRes = await fetch(`${API_BASE}/auth/refresh`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Cookie: `refresh_token=${refreshToken}`,
        },
      });

      if (refreshRes.ok) {
        const refreshData = await refreshRes.json();
        const newToken = refreshData.data.access_token;

        const retry = await fetchBackend(path, newToken, { method, body });
        res = retry.res;
        data = retry.data;

        const response = NextResponse.json(data, { status: res.status });
        setAccessToken(response, newToken);

        const setCookieHeader = refreshRes.headers.get("set-cookie");
        if (setCookieHeader) {
          const match = setCookieHeader.match(/refresh_token=([^;]+)/);
          if (match) {
            setRefreshToken(response, match[1]);
          }
        }

        return response;
      }

      const errData = await refreshRes.json();
      const cleared = NextResponse.json(errData, {
        status: refreshRes.status,
      });
      clearCookie(cleared, "access_token");
      clearCookie(cleared, "refresh_token");
      return cleared;
    }
  }

  return NextResponse.json(data, { status: res.status });
}
