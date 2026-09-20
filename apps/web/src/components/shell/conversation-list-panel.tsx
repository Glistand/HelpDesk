"use client";

import Link from "next/link";
import { usePathname, useRouter, useSearchParams } from "next/navigation";
import { useMemo } from "react";
import { formatRelative } from "@/lib/format";
import type { ConversationStatus, ConversationSummary } from "@/lib/types";

type Props = {
  initialConversations: ConversationSummary[];
};

const STATUS_LABEL: Record<ConversationStatus, string> = {
  bot: "Бот",
  waiting_agent: "Ждёт агента",
  open: "Открыт",
  resolved: "Решён",
  unspecified: "—",
};

function statusTone(status: ConversationStatus): string {
  switch (status) {
    case "waiting_agent":
      return "bg-amber-500/15 text-amber-600 dark:text-amber-400";
    case "bot":
      return "bg-sky-500/15 text-sky-700 dark:text-sky-300";
    case "open":
      return "bg-emerald-500/15 text-emerald-700 dark:text-emerald-300";
    case "resolved":
      return "bg-zinc-500/15 text-tertiary";
    default:
      return "bg-elevated text-tertiary";
  }
}

export function ConversationListPanel({ initialConversations }: Props) {
  const pathname = usePathname();
  const router = useRouter();
  const searchParams = useSearchParams();
  const filter = searchParams.get("filter");

  const activeId = pathname.startsWith("/conversations/")
    ? pathname.split("/")[2]
    : null;

  const stats = useMemo(() => {
    const waiting = initialConversations.filter((c) => c.status === "waiting_agent").length;
    const bot = initialConversations.filter((c) => c.status === "bot").length;
    const open = initialConversations.filter((c) =>
      ["bot", "waiting_agent", "open"].includes(c.status),
    ).length;
    return { waiting, bot, open };
  }, [initialConversations]);

  const list = useMemo(() => {
    if (filter === "waiting") {
      return initialConversations.filter((c) => c.status === "waiting_agent");
    }
    if (filter === "bot") {
      return initialConversations.filter((c) => c.status === "bot");
    }
    if (filter === "resolved") {
      return initialConversations.filter((c) => c.status === "resolved");
    }
    return initialConversations.filter((c) => c.status !== "resolved");
  }, [initialConversations, filter]);

  const filters = [
    { key: null as string | null, label: "Активные", count: stats.open },
    { key: "waiting", label: "Handoff", count: stats.waiting },
    { key: "bot", label: "Бот", count: stats.bot },
    { key: "resolved", label: "Решённые", count: initialConversations.filter((c) => c.status === "resolved").length },
  ];

  return (
    <aside className="hidden w-[340px] shrink-0 flex-col border-r border-border bg-base md:flex">
      <div className="border-b border-border px-4 py-3">
        <div className="flex items-center justify-between">
          <h1 className="font-ui-semibold text-[15px] text-primary">Беседы</h1>
          <button
            type="button"
            onClick={() => router.refresh()}
            className="text-xs text-tertiary hover:text-accent"
          >
            Обновить
          </button>
        </div>
        <div className="mt-3 flex flex-wrap gap-1">
          {filters.map((f) => {
            const href = f.key ? `/inbox?filter=${f.key}` : "/inbox";
            const active = (f.key === null && !filter) || filter === f.key;
            return (
              <Link
                key={f.label}
                href={href}
                className={`rounded-md px-2 py-1 text-[11px] font-ui-medium transition-colors ${
                  active
                    ? "bg-elevated text-primary"
                    : "text-tertiary hover:bg-elevated/70 hover:text-secondary"
                }`}
              >
                {f.label}
                <span className="ml-1 tabular-nums opacity-70">{f.count}</span>
              </Link>
            );
          })}
        </div>
      </div>

      <div className="flex-1 overflow-y-auto">
        {list.length === 0 ? (
          <p className="px-4 py-8 text-center text-[13px] text-tertiary">
            Нет бесед. Откройте{" "}
            <a className="text-accent underline" href="/widget-demo.html" target="_blank" rel="noreferrer">
              виджет-демо
            </a>
            .
          </p>
        ) : (
          <ul className="divide-y divide-border">
            {list.map((c) => {
              const active = activeId === c.id;
              return (
                <li key={c.id}>
                  <Link
                    href={`/conversations/${c.id}`}
                    className={`block px-4 py-3 transition-colors ${
                      active ? "bg-elevated" : "hover:bg-elevated/60"
                    }`}
                  >
                    <div className="flex items-start justify-between gap-2">
                      <span
                        className={`rounded px-1.5 py-0.5 text-[10px] font-ui-medium ${statusTone(c.status)}`}
                      >
                        {STATUS_LABEL[c.status]}
                      </span>
                      <span className="shrink-0 text-[11px] text-tertiary tabular-nums">
                        {formatRelative(c.updatedAt)}
                      </span>
                    </div>
                    <p className="mt-1.5 line-clamp-2 text-[13px] leading-snug text-secondary">
                      {c.preview || "Новая беседа"}
                    </p>
                    <p className="mt-1 font-mono text-[10px] text-tertiary">
                      {c.visitorId.slice(0, 8)}…
                    </p>
                  </Link>
                </li>
              );
            })}
          </ul>
        )}
      </div>
    </aside>
  );
}
