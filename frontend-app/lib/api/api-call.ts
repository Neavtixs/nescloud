import { getDefaultStore } from "jotai";
import type { UserData } from "@/lib/atoms/auth-atoms";
import { accessTokenAtom } from "@/lib/atoms/auth-atoms";
import { ApiResponse } from "./api-response";

const store = getDefaultStore();

let refreshLock: Promise<boolean> | null = null;

async function doRefresh(signal?: AbortSignal): Promise<boolean> {
  if (refreshLock) return refreshLock;

  refreshLock = (async () => {
    try {
      const res = await fetch(
        `${process.env.NEXT_PUBLIC_HOST_API}/api/auth/refresh`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          credentials: "include",
          signal,
        },
      );

      if (res.ok) {
        const data = await res.json();
        if (data.data?.access_token) {
          store.set(accessTokenAtom, data.data.access_token);
        }
        return true;
      }

      return false;
    } finally {
      refreshLock = null;
    }
  })();

  return refreshLock;
}

export class ApiError extends Error {
  status: number;

  constructor(message: string, status: number) {
    super(message);
    this.status = status;
    this.name = "ApiError";
  }
}

async function request<T>(
  path: string,
  options?: RequestInit & { needsAuth?: boolean },
  _retryCount = 0,
): Promise<T> {
  const { needsAuth = true, signal } = options ?? {};
  const token = store.get(accessTokenAtom);

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(options?.headers as Record<string, string>),
  };

  if (needsAuth && token) {
    headers.Authorization = `Bearer ${token}`;
  }

  const res = await fetch(`${process.env.NEXT_PUBLIC_HOST_API}/api${path}`, {
    headers,
    credentials: "include",
    signal,
    ...options,
  });

  const data = await res.json();
  if (needsAuth && res.status === 401 && _retryCount < 1) {
    const refreshed = await doRefresh(signal);

    if (refreshed) {
      return request<T>(path, options, _retryCount + 1);
    }

    if (typeof window !== "undefined") {
      window.location.href = "/login";
    }
    throw new ApiError("Session expired", 401);
  }

  if (!res.ok) {
    throw new ApiError(data.message ?? "Something went wrong", res.status);
  }

  return data as T;
}

export const api = {
  get: <T>(path: string, options?: RequestInit & { needsAuth?: boolean }) =>
    request<T>(path, { ...options, method: "GET" }),
  post: <T>(
    path: string,
    body: unknown,
    options?: RequestInit & { needsAuth?: boolean },
  ) =>
    request<T>(path, {
      method: "POST",
      body: JSON.stringify(body),
      ...options,
    }),
  put: <T>(
    path: string,
    body: unknown,
    options?: RequestInit & { needsAuth?: boolean },
  ) =>
    request<T>(path, { method: "PUT", body: JSON.stringify(body), ...options }),
  patch: <T>(
    path: string,
    body: unknown,
    options?: RequestInit & { needsAuth?: boolean },
  ) =>
    request<T>(path, {
      method: "PATCH",
      body: JSON.stringify(body),
      ...options,
    }),
  delete: <T>(path: string, options?: RequestInit & { needsAuth?: boolean }) =>
    request<T>(path, { ...options, method: "DELETE" }),
};

type AuthData = { access_token: string };

export const authApi = {
  login: (
    body: { email: string; password: string },
    options?: RequestInit & { needsAuth?: boolean },
  ) =>
    api.post<ApiResponse<AuthData>>("/auth/login", body, {
      ...options,
      needsAuth: false,
    }),
  register: (
    body: { name: string; email: string; password: string },
    options?: RequestInit & { needsAuth?: boolean },
  ) =>
    api.post<ApiResponse<AuthData>>("/auth/register", body, {
      ...options,
      needsAuth: false,
    }),
  me: (options?: RequestInit & { needsAuth?: boolean }) =>
    api.get<ApiResponse<UserData>>("/auth/me", options),
  logout: (options?: RequestInit & { needsAuth?: boolean }) => {
    store.set(accessTokenAtom, null);
    return api.post<ApiResponse<null>>("/auth/logout", {}, options);
  },
};
