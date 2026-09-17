"use client";

import { useRouter } from "next/navigation";
import { useState, useTransition } from "react";
import type { TicketStatus } from "@/lib/types";

const statuses: { value: TicketStatus; label: string }[] = [
  { value: "open", label: "Открыт" },
  { value: "pending", label: "Ожидание" },
  { value: "resolved", label: "Решён" },
  { value: "closed", label: "Закрыт" },
];

export function StatusActions({
  ticketId,
  current,
}: {
  ticketId: string;
  current: TicketStatus;
}) {
  const router = useRouter();
  const [open, setOpen] = useState(false);
  const [pending, startTransition] = useTransition();
  const [error, setError] = useState("");

  function setStatus(status: TicketStatus) {
    setOpen(false);
    setError("");
    startTransition(async () => {
      const res = await fetch(`/api/hd/tickets/${ticketId}/status`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ status }),
      });
      if (!res.ok) {
        setError("Не удалось обновить статус");
        return;
      }
      router.refresh();
    });
  }

  return (
    <div className="relative">
      <button
        type="button"
        disabled={pending}
        onClick={() => setOpen((v) => !v)}
        className="rounded-md bg-accent px-3 py-1.5 text-[13px] font-ui-medium text-white hover:opacity-90 disabled:opacity-60"
      >
        {pending ? "…" : "Статус"}
      </button>
      {open && (
        <div className="absolute right-0 z-20 mt-1 min-w-[140px] rounded-md border border-border bg-view py-1 shadow-lg">
          {statuses.map((s) => (
            <button
              key={s.value}
              type="button"
              onClick={() => setStatus(s.value)}
              className={`block w-full px-3 py-1.5 text-left text-[13px] hover:bg-elevated ${
                s.value === current ? "text-accent" : "text-secondary"
              }`}
            >
              {s.label}
            </button>
          ))}
        </div>
      )}
      {error && <p className="mt-1 text-[11px] text-status-red">{error}</p>}
    </div>
  );
}
