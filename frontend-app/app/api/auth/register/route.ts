import { NextRequest } from "next/server";
import { apiProxyAuth } from "@/lib/api-proxy-auth";

export async function POST(request: NextRequest) {
  return apiProxyAuth(request, "/auth/register");
}
