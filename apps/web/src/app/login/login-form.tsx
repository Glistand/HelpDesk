"use client";

import { FormEvent, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";

export default function LoginForm() {
  const router = useRouter();
  const search = useSearchParams();
  const [email, setEmail] = useState("agent@helpdesk.local");
  const [password, setPassword] = useState("password");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError("");
    try {
      const res = await fetch("/api/auth/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
      });
      if (!res.ok) {
        setError("Неверный email или пароль");
        return;
      }
      const next = search.get("next") || "/inbox";
      router.replace(next);
      router.refresh();
    } catch {
      setError("Нет связи с API");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-base px-4">
      <div className="w-full max-w-sm">
        <div className="mb-8 text-center">
          <div className="mx-auto mb-3 flex h-10 w-10 items-center justify-center rounded-md bg-accent font-ui-semibold text-white">
            H
          </div>
          <h1 className="font-ui-semibold text-xl text-primary">Helpdesk</h1>
          <p className="mt-1 text-[13px] text-tertiary">
            Войдите как агент · Event Hub
          </p>
        </div>
        <form
          onSubmit={onSubmit}
          className="space-y-3 rounded-lg border border-border bg-view p-5"
        >
          <label className="block space-y-1.5">
            <span className="text-[13px] font-ui-medium text-secondary">Email</span>
            <input
              type="email"
              required
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="field-input"
            />
          </label>
          <label className="block space-y-1.5">
            <span className="text-[13px] font-ui-medium text-secondary">Пароль</span>
            <input
              type="password"
              required
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="field-input"
            />
          </label>
          {error && (
            <p className="text-[13px] text-status-red">{error}</p>
          )}
          <button
            type="submit"
            disabled={loading}
            className="w-full rounded-md bg-accent px-4 py-2 text-[13px] font-ui-medium text-white hover:opacity-90 disabled:opacity-60"
          >
            {loading ? "Вход…" : "Войти"}
          </button>
          <p className="text-center text-[11px] text-tertiary">
            seed: agent@helpdesk.local / password
          </p>
        </form>
      </div>
    </div>
  );
}
