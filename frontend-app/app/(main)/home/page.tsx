"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useAtom } from "jotai";
import { userAtom } from "@/lib/atoms/auth-atoms";
import { authApi } from "@/lib/api/api-call";

export default function HomePage() {
  const router = useRouter();
  const [user, setUser] = useAtom(userAtom);
  const [isReloading, setIsReloading] = useState(false);

  async function handleLogout() {
    await fetch("/api/auth/logout", {
      method: "POST",
    });
    router.push("/login");
  }

  async function handleReload() {
    setIsReloading(true);
    try {
      const data = await authApi.me();
      setUser(data.data);
    } finally {
      setIsReloading(false);
    }
  }

  return (
    <div>
      <h1>ini home page</h1>

      <div>
        <button onClick={handleReload} disabled={isReloading}>
          {isReloading ? "Memuat..." : "Reload"}
        </button>
      </div>
      <div>
        {user && (
          <div>
            <p>Nama: {user.name}</p>
            <p>Email: {user.email}</p>
            <p>
              Bergabung: {new Date(user.created_at).toLocaleDateString("id-ID")}
            </p>
          </div>
        )}

        <button onClick={handleLogout}>Logout</button>
      </div>
    </div>
  );
}
