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

  /*
  useEffect(() => {

    async function validate() {
      try {
        const res = await fetch("/api/auth/me");
        if (!res.ok) throw new Error("unauthorized");
        const json = await res.json();
        setUser(json.data);
      } catch {
        router.push("/login");
      } finally {
        setIsLoading(false);
      }
    }
    return () => validate();
  }, []);
	*/

  useEffect(() => {
    const controller = new AbortController();

    async function validate() {
      try {
        const res = await fetch("/api/auth/me", {
          signal: controller.signal,
        });

        if (!res.ok) throw new Error();

        const json = await res.json();

        setUser(json.data);
      } catch (err) {
        if ((err as Error).name === "AbortError") return;

        router.replace("/login");
      } finally {
        setIsLoading(false);
      }
    }

    validate();

    return () => controller.abort();
  }, []);

  if (isLoading) return <div>Loading...</div>;

  return <>{children}</>;
}
