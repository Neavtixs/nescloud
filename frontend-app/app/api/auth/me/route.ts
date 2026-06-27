import { NextRequest } from "next/server"
import { apiProxyProtect } from "@/lib/api/proxy/api-proxy-protect"

export async function GET(request: NextRequest) {
  return apiProxyProtect(request, "/auth/me")
}
