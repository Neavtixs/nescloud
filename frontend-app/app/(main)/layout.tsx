"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAtom } from "jotai";
import { userAtom, isAuthLoadingAtom } from "@/lib/atoms/auth-atoms";
import { ApiError, authApi } from "@/lib/api/api-call";

export default function MainLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();
  const [, setUser] = useAtom(userAtom);
  const [isLoading, setIsLoading] = useAtom(isAuthLoadingAtom);

  useEffect(() => {
    const controller = new AbortController();

    async function validate() {
      try {
        const res = await authApi.me();
        setUser(res.data);
      } catch (err) {
        if (err instanceof ApiError) {
          if (err.status == 401) {
            router.replace("/login");
          }
        } else {
        }
      } finally {
        setIsLoading(false);
      }
    }

    validate();

    return () => controller.abort();
  }, [router, setUser, setIsLoading]);

  if (isLoading) return <div>Loading...</div>;

  return <>{children}</>;
}
