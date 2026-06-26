const API_BASE = process.env.NEXT_PUBLIC_HOST_API + "/api"

export class ApiError extends Error {
  status: number

  constructor(message: string, status: number) {
    super(message)
    this.status = status
    this.name = "ApiError"
  }
}

export async function apiPost<T>(
  path: string,
  body: Record<string, unknown>
): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  })

  const data = await res.json()

  if (!res.ok) {
    throw new ApiError(data.message ?? "Something went wrong", res.status)
  }

  return data as T
}
