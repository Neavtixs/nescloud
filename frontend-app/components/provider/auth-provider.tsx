"use client";

import { Suspense, use, useEffect } from "react";
import { useSetAtom } from "jotai";
import { userAtom, isAuthLoadingAtom } from "@/lib/atoms/auth-atoms";
import { authApi } from "@/lib/api/api-call";
import type { UserData } from "@/lib/atoms/auth-atoms";

let userPromise: Promise<UserData | null> | null = null;

function getUserPromise() {
  if (!userPromise) {
    userPromise = authApi
      .me()
      .then((res) => res.data)
      .catch(() => null);
  }
  return userPromise;
}

export function invalidateAuth() {
  userPromise = null;
}

function AuthGate({ children }: { children: React.ReactNode }) {
  const user = use(getUserPromise());
  const setUser = useSetAtom(userAtom);
  const setIsLoading = useSetAtom(isAuthLoadingAtom);

  useEffect(() => {
    setUser(user);
    setIsLoading(false);
  }, [user, setUser, setIsLoading]);

  return <>{children}</>;
}

function LoadingFallback() {
  return <div>Loading...fall</div>;
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
