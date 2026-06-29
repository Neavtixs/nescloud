"use client";

import { Suspense, use, useEffect, useState } from "react";
import { useSetAtom } from "jotai";
import { userAtom, isAuthLoadingAtom } from "@/lib/atoms/auth-atoms";
import { authApi } from "@/lib/api/api-call";
import type { UserData } from "@/lib/atoms/auth-atoms";
import { WifiOff, RefreshCw } from "lucide-react";

type AuthState = {
  user: UserData | null;
  error: "connection" | null;
};

let authPromise: Promise<AuthState> | null = null;

function isNetworkError(err: unknown): boolean {
  if (err instanceof TypeError) return true;
  if (
    err &&
    typeof err === "object" &&
    "message" in err &&
    typeof (err as Record<string, unknown>).message === "string"
  ) {
    const msg = (err as Record<string, string>).message.toLowerCase();
    return msg.includes("fetch") || msg.includes("network");
  }
  return false;
}

function getAuthState(): Promise<AuthState> {
  if (!authPromise) {
    authPromise = authApi
      .me()
      .then(
        (res): AuthState => ({ user: res.data, error: null }),
      )
      .catch((err): AuthState => {
        if (isNetworkError(err)) {
          return { user: null, error: "connection" };
        }
        return { user: null, error: null };
      });
  }
  return authPromise;
}

export function invalidateAuth() {
  authPromise = null;
}

function ConnectionError({ onRetry }: { onRetry: () => void }) {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-gray-50 px-4 dark:bg-gray-950">
      <div className="text-center">
        <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-red-100 dark:bg-red-900/30">
          <WifiOff size={28} className="text-red-600 dark:text-red-400" />
        </div>
        <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100">
          Connection Error
        </h2>
        <p className="mt-2 text-sm text-gray-500 dark:text-gray-400">
          Unable to reach the server. Please check your internet connection and
          try again.
        </p>
        <button
          onClick={onRetry}
          className="mt-6 inline-flex items-center gap-2 rounded-lg bg-blue-600 px-5 py-2.5 text-sm font-medium text-white hover:bg-blue-700"
        >
          <RefreshCw size={16} />
          Try Again
        </button>
      </div>
    </div>
  );
}

function AuthGate({ children }: { children: React.ReactNode }) {
  const { user, error } = use(getAuthState());
  const setUser = useSetAtom(userAtom);
  const setIsLoading = useSetAtom(isAuthLoadingAtom);
  const [, forceUpdate] = useState(0);

  useEffect(() => {
    setUser(user);
    setIsLoading(false);
  }, [user, setUser, setIsLoading]);

  if (error === "connection") {
    return (
      <ConnectionError
        onRetry={() => {
          invalidateAuth();
          forceUpdate((n) => n + 1);
        }}
      />
    );
  }

  return <>{children}</>;
}

function LoadingFallback() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50 text-sm text-gray-400 dark:bg-gray-950 dark:text-gray-500">
      Loading...
    </div>
  );
}

export default function AuthProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <Suspense fallback={<LoadingFallback />}>
      <AuthGate>{children}</AuthGate>
    </Suspense>
  );
}
