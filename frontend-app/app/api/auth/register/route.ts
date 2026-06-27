import { apiProxyAuth } from "@/lib/api/proxy/api-proxy-auth";
import { NextRequest } from "next/server";

export async function POST(request: NextRequest) {
  return apiProxyAuth(request, "/auth/register");
}
