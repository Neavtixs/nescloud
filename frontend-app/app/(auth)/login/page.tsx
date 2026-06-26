"use client";

import { useActionState, useEffect } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";

type State = {
  error?: string;
  success?: boolean;
  data?: { access_token: string; expires_in: number };
} | null;

async function loginAction(_prev: State, formData: FormData) {
  const email = formData.get("email") as string;
  const password = formData.get("password") as string;

  const res = await fetch("/api/auth/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  });

  const data = await res.json();

  if (!res.ok) {
    return { error: data.message ?? "Login failed" };
  }

  return { success: true, data: data.data };
}

export default function LoginPage() {
  const [state, formAction, pending] = useActionState(loginAction, null);
  const router = useRouter();

  useEffect(() => {
    if (state?.success) {
      //    router.push("/")
    }
  }, [state, router]);

  return (
    <form action={formAction} className="flex flex-col gap-6">
      <div className="flex flex-col gap-2 text-center">
        <h1 className="text-2xl font-semibold tracking-tight">Welcome back</h1>
        <p className="text-sm text-zinc-500 dark:text-zinc-400">
          Sign in to your account
        </p>
      </div>

      <div className="flex flex-col gap-4">
        <div className="flex flex-col gap-1.5">
          <label htmlFor="email" className="text-sm font-medium">
            Email
          </label>
          <input
            id="email"
            name="email"
            type="email"
            required
            placeholder="lawrient@mail.com"
            className="rounded-lg border border-zinc-300 px-3 py-2 text-sm outline-none focus:border-zinc-500 dark:border-zinc-700 dark:bg-zinc-900 dark:focus:border-zinc-400"
          />
        </div>

        <div className="flex flex-col gap-1.5">
          <label htmlFor="password" className="text-sm font-medium">
            Password
          </label>
          <input
            id="password"
            name="password"
            type="password"
            required
            placeholder="password"
            className="rounded-lg border border-zinc-300 px-3 py-2 text-sm outline-none focus:border-zinc-500 dark:border-zinc-700 dark:bg-zinc-900 dark:focus:border-zinc-400"
          />
        </div>
      </div>

      {state?.error && (
        <p className="text-sm text-red-600 dark:text-red-400">{state.error}</p>
      )}

      <button
        type="submit"
        disabled={pending}
        className="rounded-lg bg-zinc-900 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-zinc-800 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-200"
      >
        {pending ? "Signing in..." : "Sign in"}
      </button>

      <p className="text-center text-sm text-zinc-500 dark:text-zinc-400">
        Don&apos;t have an account?{" "}
        <Link
          href="/register"
          className="font-medium text-zinc-900 underline underline-offset-2 dark:text-zinc-100"
        >
          Create one
        </Link>
      </p>
    </form>
  );
}
