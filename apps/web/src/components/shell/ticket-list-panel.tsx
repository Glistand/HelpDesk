"use client";

import Link from "next/link";
import { usePathname, useRouter, useSearchParams } from "next/navigation";
import { useCallback, useEffect, useMemo, useState, useTransition } from "react";
import { PriorityDot, SlaDot } from "@/components/badges";
import { formatRelative } from "@/lib/format";
import { currentAgent } from "@/lib/agents";
import type { TicketSummary } from "@/lib/types";

type Props = {
  initialTickets: TicketSummary[];
};

export function TicketListPanel({ initialTickets }: Props) {
  const pathname = usePathname();
  const router = useRouter();
  const searchParams = useSearchParams();
  const filter = searchParams.get("filter");
  const qParam = searchParams.get("q") || "";
  const [query, setQuery] = useState(qParam);
  const [tickets, setTickets] = useState(initialTickets);
  const [pending, startTransition] = useTransition();

  useEffect(() => {
    setTickets(initialTickets);
  }, [initialTickets]);

  useEffect(() => {
    setQuery(qParam);
  }, [qParam]);

  const activeId = pathname.startsWith("/tickets/")
    ? pathname.split("/")[2]
    : null;

  const runSearch = useCallback(
    (value: string) => {
      startTransition(async () => {
        if (!value.trim()) {
          router.replace(filter ? `/inbox?filter=${filter}` : "/inbox");
          router.refresh();
          return;
        }
        const res = await fetch(`/api/hd/search?q=${encodeURIComponent(value)}&limit=30`);
        if (!res.ok) return;
        const data = await res.json();
        const hits: TicketSummary[] = (data.hits || []).map(
          (h: {
            id: string;
            title: string;
            status: string;
            priority: string;
            category: string;
            requester: string;
            assignee_id: string;
            updated_at: string;
          }) => ({
            id: h.id,
            title: h.title,
            status: h.status,
            priority: h.priority,
            category: h.category,
            requester: h.requester,
            assigneeId: h.assignee_id || "",
            updatedAt: h.updated_at,
            createdAt: h.updated_at,
            slaState: "ok" as const,
          }),
        );
        setTickets(hits);
        const params = new URLSearchParams();
        if (filter) params.set("filter", filter);
        params.set("q", value);
        router.replace(`/inbox?${params.toString()}`);
      });
    },
    [filter, router],
  );

  const list = useMemo(() => {
    if (qParam.trim()) return tickets;
    if (filter === "mine") {
      return tickets.filter((t) => t.assigneeId === currentAgent.id);
    }
    if (filter === "breached") {
      return tickets.filter((t) => t.slaState === "breached");
    }
    return tickets.filter((t) => t.status !== "resolved" && t.status !== "closed");
  }, [tickets, filter, qParam]);

  const stats = useMemo(() => {
    const open = tickets.filter((t) => !["resolved", "closed"].includes(t.status)).length;
    const mine = tickets.filter((t) => t.assigneeId === currentAgent.id).length;
    const breached = tickets.filter((t) => t.slaState === "breached").length;
    return { open, mine, breached };
  }, [tickets]);

  const filters = [
    { key: null as string | null, label: "Все", count: stats.open },
    { key: "mine", label: "Мои", count: stats.mine },
    { key: "breached", label: "Breach", count: stats.breached },
  ];

  return (
    <aside className="hidden w-[340px] shrink-0 flex-col border-r border-border bg-base md:flex">
      <div className="border-b border-border px-4 py-3">
        <div className="flex items-center justify-between">
          <h1 className="font-ui-semibold text-[15px] text-primary">Inbox</h1>
          <span className="text-xs text-tertiary tabular-nums">
            {pending ? "…" : list.length}
          </span>
        </div>
        <div className="mt-3 flex gap-1">
          {filters.map((f) => {
            const href = f.key ? `/inbox?filter=${f.key}` : "/inbox";
            const active = (f.key === null && !filter) || filter === f.key;
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
        <form
          className="mt-3"
          onSubmit={(e) => {
            e.preventDefault();
            runSearch(query);
          }}
        >
          <input
            type="search"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Поиск тикетов…"
            className="w-full rounded-md border border-border bg-input px-3 py-1.5 text-[13px] text-secondary placeholder:text-tertiary focus:border-accent/40 focus:outline-none"
          />
        </form>
      </div>

      <div className="flex-1 overflow-y-auto">
        {list.length === 0 ? (
          <p className="px-4 py-8 text-center text-[13px] text-tertiary">
            Нет тикетов
          </p>
        ) : (
          list.map((ticket) => {
            const active = activeId === ticket.id;
            return (
              <Link
                key={ticket.id}
                href={`/tickets/${ticket.id}`}
                className={`block border-b border-border px-4 py-3 transition-colors ${
                  active ? "bg-view ring-surface" : "hover:bg-elevated/50"
                }`}
              >
                <div className="flex items-start gap-2">
                  <PriorityDot priority={ticket.priority} />
                  <div className="min-w-0 flex-1">
                    <div className="flex items-center gap-2">
                      <span className="font-mono text-[11px] text-tertiary">
                        {ticket.id}
                      </span>
                      <SlaDot state={ticket.slaState} />
                    </div>
                    <p
                      className={`mt-0.5 truncate text-[13px] leading-snug ${
                        active ? "font-ui-medium text-primary" : "text-secondary"
                      }`}
                    >
                      {ticket.title}
                    </p>
                    <p className="mt-1 truncate text-[11px] text-tertiary">
                      {ticket.requester || "—"} · {formatRelative(ticket.updatedAt)}
                    </p>
                  </div>
                </div>
              </Link>
            );
          })
        )}
      </div>
    </aside>
  );
}
