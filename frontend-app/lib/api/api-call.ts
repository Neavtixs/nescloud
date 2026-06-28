import { getDefaultStore } from "jotai";
import type { UserData } from "@/lib/atoms/auth-atoms";
import { accessTokenAtom } from "@/lib/atoms/auth-atoms";
import { ApiResponse } from "./api-response";

const store = getDefaultStore();

export class ApiError extends Error {
  status: number;

  constructor(message: string, status: number) {
    super(message);
    this.status = status;
    this.name = "ApiError";
  }
}

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const token = store.get(accessTokenAtom);

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(options?.headers as Record<string, string>),
  };

  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  const res = await fetch(`${process.env.NEXT_PUBLIC_HOST_API}/api${path}`, {
    headers,
    credentials: "include",
    ...options,
  });

  const data = await res.json();
  console.log(data + res.status);
  if (res.status === 401) {
    console.log("401" + data);
    const refreshRes = await fetch(
      `${process.env.NEXT_PUBLIC_HOST_API}/api/auth/refresh`,
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
      },
    );

    if (refreshRes.ok) {
      const refreshData = await refreshRes.json();
      if (refreshData.data?.access_token) {
        store.set(accessTokenAtom, refreshData.data.access_token);
      }
      return request<T>(path, options);
    }

    if (typeof window !== "undefined") {
      window.location.href = "/login";
    }
    throw new ApiError("Session expired", 401);
  }
  console.log("lewwat ");

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
  logout: () => {
    store.set(accessTokenAtom, null);
    return api.post<ApiResponse<null>>("/auth/logout", {});
  },
  refresh: () => api.post<ApiResponse<AuthData>>("/auth/refresh", {}),
};
