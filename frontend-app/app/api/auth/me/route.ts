import { NextRequest, NextResponse } from "next/server";

const API_BASE = process.env.NEXT_PUBLIC_HOST_API + "/api";

async function callBackend(path: string, token?: string) {
  const res = await fetch(`${API_BASE}${path}`, {
    headers: {
      "Content-Type": "application/json",
      ...(token && { Authorization: `Bearer ${token}` }),
    },
  });
  const data = await res.json();
  return { res, data };
}

export async function GET(request: NextRequest) {
  const accessToken = request.cookies.get("access_token")?.value;

  let { res, data } = await callBackend("/auth/me", accessToken);

  if (res.status === 401 && accessToken) {
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

        const retryRes = await callBackend(
          "/auth/me",
          refreshData.data.access_token,
        );
        res = retryRes.res;
        data = retryRes.data;

        const response = NextResponse.json(data, { status: res.status });
        response.cookies.set("access_token", refreshData.data.access_token, {
          httpOnly: true,
          sameSite: "lax",
          path: "/",
        });

        const setCookie = refreshRes.headers.get("set-cookie");
        if (setCookie) {
          const match = setCookie.match(/refresh_token=([^;]+)/);
          if (match) {
            response.cookies.set("refresh_token", match[1], {
              httpOnly: true,
              sameSite: "lax",
              path: "/",
            });
          }
        }

        return response;
      }

      const errData = await refreshRes.json();
      const cleared = NextResponse.json(errData, {
        status: refreshRes.status,
      });
      cleared.cookies.set("access_token", "", {
        maxAge: 0,
        httpOnly: true,
        sameSite: "lax",
        path: "/",
      });
      cleared.cookies.set("refresh_token", "", {
        maxAge: 0,
        httpOnly: true,
        sameSite: "lax",
        path: "/",
      });
      return cleared;
    }
  }

  return NextResponse.json(data, { status: res.status });
}
