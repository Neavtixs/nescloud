"use client";

import { useRouter } from "next/navigation";

export default function HomePage() {
  const router = useRouter();

  async function handleLogout() {
    await fetch("/api/auth/logout", {
      method: "POST",
    });
    router.push("/login");
  }

  return (
    <div>
      <h1>ini home page</h1>
      <button onClick={handleLogout}>Logout</button>
    </div>
  );
}
