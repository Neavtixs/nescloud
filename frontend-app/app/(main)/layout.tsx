"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAtom } from "jotai";
import { userAtom, isAuthLoadingAtom } from "@/lib/atoms/auth-atoms";

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
        const res = await fetch("/api/auth/me", {
          signal: controller.signal,
        });

        if (!res.ok) {
          router.replace("/login");
          return;
        }

        const { data } = await res.json();
        setUser(data);
      } catch (err) {
        if ((err as Error).name !== "AbortError") {
          router.replace("/login");
        }
      } finally {
        if (!controller.signal.aborted) {
          setIsLoading(false);
        }
      }
    }

    validate();

    return () => controller.abort();
  }, [router, setUser, setIsLoading]);

  if (isLoading) return <div>Loading...</div>;

  return <>{children}</>;
}
