import type { AuthResponse, LoginPayload, RegisterPayload } from "@/lib/types/auth"

export class AuthError extends Error {
  status: number

  constructor(message: string, status: number) {
    super(message)
    this.status = status
    this.name = "AuthError"
  }
}

async function authFetch<T>(
  path: string,
  body: object
): Promise<T> {
  const res = await fetch(path, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  })

  const data = await res.json()

  if (!res.ok) {
    throw new AuthError(data.message ?? "Something went wrong", res.status)
  }

  return data as T
}

export async function login(payload: LoginPayload) {
  return authFetch<AuthResponse>("/api/auth/login", payload)
}

export async function register(payload: RegisterPayload) {
  return authFetch<AuthResponse>("/api/auth/register", payload)
}
