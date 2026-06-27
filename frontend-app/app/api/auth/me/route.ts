import { NextRequest } from "next/server";
import { apiProxyAuth } from "@/lib/api-proxy-auth";

export async function GET(request: NextRequest) {
  return apiProxyAuth(request, "/auth/me", { method: "GET", forwardBody: false });
}
