"use client";

import Link from "next/link";
import { usePathname, useSearchParams } from "next/navigation";
import { PriorityDot, SlaDot } from "@/components/badges";
import { formatRelative } from "@/lib/format";
import { inboxStats, tickets } from "@/lib/mock-data";
import type { Ticket } from "@/lib/types";

function filterTickets(filter: string | null): Ticket[] {
  if (filter === "mine") {
    return tickets.filter((t) => t.assignee?.name === "Алексей К.");
  }
  if (filter === "breached") {
    return tickets.filter((t) => t.sla.state === "breached");
  }
  return tickets.filter(
    (t) => t.status !== "resolved" && t.status !== "closed",
  );
}

export function TicketListPanel() {
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const filter = searchParams.get("filter");
  const list = filterTickets(filter);
  const activeId = pathname.startsWith("/tickets/")
    ? pathname.split("/")[2]
    : null;

  const filters = [
    { key: null, label: "Все", count: inboxStats.open },
    { key: "mine", label: "Мои", count: inboxStats.mine },
    { key: "breached", label: "Breach", count: inboxStats.breached },
  ];

  return (
    <aside className="hidden w-[340px] shrink-0 flex-col border-r border-border bg-base md:flex">
      <div className="border-b border-border px-4 py-3">
        <div className="flex items-center justify-between">
          <h1 className="font-ui-semibold text-[15px] text-primary">Inbox</h1>
          <span className="text-xs text-tertiary tabular-nums">{list.length}</span>
        </div>
        <div className="mt-3 flex gap-1">
          {filters.map((f) => {
            const href = f.key ? `/inbox?filter=${f.key}` : "/inbox";
            const active =
              (f.key === null && !filter) || filter === f.key;
            return (
              <Link
                key={f.label}
                href={href}
                className={`rounded-md px-2 py-1 text-xs transition-colors ${
                  active
                    ? "bg-elevated text-primary ring-surface"
                    : "text-tertiary hover:text-secondary"
                }`}
              >
                {f.label}
                {f.count > 0 && (
                  <span className="ml-1 tabular-nums opacity-70">{f.count}</span>
                )}
              </Link>
            );
          })}
        </div>
        <input
          type="search"
          placeholder="Поиск тикетов…"
          className="mt-3 w-full rounded-md border border-border bg-input px-3 py-1.5 text-[13px] text-secondary placeholder:text-tertiary focus:border-accent/40 focus:outline-none"
        />
      </div>

      <div className="flex-1 overflow-y-auto">
        {list.map((ticket) => {
          const active = activeId === ticket.id;
          return (
            <Link
              key={ticket.id}
              href={`/tickets/${ticket.id}`}
              className={`block border-b border-border px-4 py-3 transition-colors ${
                active
                  ? "bg-view ring-surface"
                  : "hover:bg-elevated/50"
              }`}
            >
              <div className="flex items-start gap-2">
                <PriorityDot priority={ticket.priority} />
                <div className="min-w-0 flex-1">
                  <div className="flex items-center gap-2">
                    <span className="font-mono text-[11px] text-tertiary">
                      {ticket.id}
                    </span>
                    <SlaDot state={ticket.sla.state} />
                  </div>
                  <p
                    className={`mt-0.5 truncate text-[13px] leading-snug ${
                      active ? "font-ui-medium text-primary" : "text-secondary"
                    }`}
                  >
                    {ticket.title}
                  </p>
                  <p className="mt-1 truncate text-[11px] text-tertiary">
                    {ticket.requester} · {formatRelative(ticket.updatedAt)}
                  </p>
                </div>
              </div>
            </Link>
          );
        })}
      </div>
    </aside>
  );
}
