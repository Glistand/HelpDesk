"use client";

import { useRouter } from "next/navigation";
import { useState, useTransition } from "react";
import type { AuthUser } from "@/lib/types";

export function AssignActions({
  ticketId,
  currentAssigneeId,
  users,
}: {
  ticketId: string;
  currentAssigneeId?: string;
  users: AuthUser[];
}) {
  const router = useRouter();
  const [open, setOpen] = useState(false);
  const [pending, startTransition] = useTransition();
  const [error, setError] = useState("");

  const agents = users.filter((u) => u.role === "agent" || u.role === "admin");

  function assign(assigneeId: string) {
    setOpen(false);
    setError("");
    startTransition(async () => {
      const res = await fetch(`/api/hd/tickets/${ticketId}/assign`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ assignee_id: assigneeId }),
      });
      if (!res.ok) {
        const data = await res.json().catch(() => ({}));
        setError(data.error || "Не удалось назначить");
        return;
      }
      router.refresh();
    });
  }

  return (
    <div className="relative">
      <button
        type="button"
        disabled={pending || agents.length === 0}
        onClick={() => setOpen((v) => !v)}
        className="btn-accent px-3 py-1.5 text-[13px] disabled:opacity-60"
      >
        {pending ? "…" : "Назначить"}
      </button>
      {open && (
        <div className="absolute right-0 z-20 mt-1 min-w-[200px] rounded-xl border border-border bg-view py-1 shadow-elevated">
          {agents.map((u) => (
            <button
              key={u.id}
              type="button"
              onClick={() => assign(u.id)}
              className={`block w-full px-3 py-2 text-left text-[13px] hover:bg-elevated ${
                u.id === currentAssigneeId ? "text-accent" : "text-secondary"
              }`}
            >
              <span className="font-ui-medium text-primary">{u.name}</span>
              <span className="mt-0.5 block text-[11px] text-tertiary">
                {u.email} · {u.role}
              </span>
            </button>
          ))}
        </div>
      )}
      {error && <p className="mt-1 text-[11px] text-status-red">{error}</p>}
    </div>
  );
}
