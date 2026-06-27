import type { ApiResponse } from "@/lib/api/types/api-response";
import type { UserData } from "@/lib/atoms/auth-atoms";

export class ApiError extends Error {
  status: number;

  constructor(message: string, status: number) {
    super(message);
    this.status = status;
    this.name = "ApiError";
  }
}

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`/api${path}`, {
    headers: { "Content-Type": "application/json" },
    ...options,
  });

  const data = await res.json();

  if (!res.ok) {
    throw new ApiError(data.message ?? "Something went wrong", res.status);
  }

  return data as T;
}

export const api = {
  get: <T>(path: string) => request<T>(path),
  post: <T>(path: string, body: unknown) =>
    request<T>(path, { method: "POST", body: JSON.stringify(body) }),
  put: <T>(path: string, body: unknown) =>
    request<T>(path, { method: "PUT", body: JSON.stringify(body) }),
  patch: <T>(path: string, body: unknown) =>
    request<T>(path, { method: "PATCH", body: JSON.stringify(body) }),
  delete: <T>(path: string) => request<T>(path, { method: "DELETE" }),
};

type AuthData = { access_token: string };

export const authApi = {
  login: (body: { email: string; password: string }) =>
    api.post<ApiResponse<AuthData>>("/auth/login", body),
  register: (body: { name: string; email: string; password: string }) =>
    api.post<ApiResponse<AuthData>>("/auth/register", body),
  me: () => api.get<ApiResponse<UserData>>("/auth/me"),
  logout: () => api.post<ApiResponse<null>>("/auth/logout", {}),
  refresh: () => api.post<ApiResponse<AuthData>>("/auth/refresh", {}),
};
