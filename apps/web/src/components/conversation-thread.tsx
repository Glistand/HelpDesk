"use client";

import { useRouter } from "next/navigation";
import { useEffect, useState, useTransition } from "react";
import type { ChatMessage, ConversationSummary } from "@/lib/types";
import { formatDateTime, formatRelative } from "@/lib/format";

const ROLE_LABEL: Record<string, string> = {
  visitor: "Посетитель",
  bot: "Бот",
  agent: "Агент",
  system: "Система",
};

export function ConversationThread({
  conversation,
  messages,
}: {
  conversation: ConversationSummary;
  messages: ChatMessage[];
}) {
  const router = useRouter();
  const [body, setBody] = useState("");
  const [error, setError] = useState("");
  const [pending, startTransition] = useTransition();

  useEffect(() => {
    if (conversation.status === "resolved") return;
    const t = window.setInterval(() => router.refresh(), 2000);
    return () => window.clearInterval(t);
  }, [conversation.status, router]);

  function send() {
    const text = body.trim();
    if (!text) return;
    setError("");
    startTransition(async () => {
      const res = await fetch(`/api/hd/conversations/${conversation.id}/messages`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ body: text }),
      });
      if (!res.ok) {
        const data = await res.json().catch(() => ({}));
        setError(data.error || "Не удалось отправить");
        return;
      }
      setBody("");
      router.refresh();
    });
  }

  function resolve() {
    setError("");
    startTransition(async () => {
      const res = await fetch(`/api/hd/conversations/${conversation.id}/resolve`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: "{}",
      });
      if (!res.ok) {
        const data = await res.json().catch(() => ({}));
        setError(data.error || "Не удалось закрыть");
        return;
      }
      router.refresh();
    });
  }

  return (
    <>
      <header className="shrink-0 border-b border-border px-5 py-4 md:px-6">
        <div className="flex flex-wrap items-start justify-between gap-4">
          <div className="min-w-0 space-y-1">
            <p className="font-mono text-[12px] text-tertiary">{conversation.id}</p>
            <h1 className="font-ui-semibold text-lg tracking-tight text-primary">
              {conversation.preview || "Беседа"}
            </h1>
            <p className="text-[13px] text-tertiary">
              {conversation.status} · visitor {conversation.visitorId.slice(0, 8)}…
              {" · "}
              {formatRelative(conversation.updatedAt)}
            </p>
          </div>
          {conversation.status !== "resolved" && (
            <button
              type="button"
              disabled={pending}
              onClick={resolve}
              className="rounded-md border border-border px-3 py-1.5 text-[13px] text-secondary hover:bg-elevated disabled:opacity-50"
            >
              Закрыть
            </button>
          )}
        </div>
      </header>

      <div className="flex flex-1 flex-col overflow-hidden">
        <section className="flex-1 space-y-3 overflow-y-auto px-5 py-5 md:px-6">
          {messages.length === 0 ? (
            <p className="text-[13px] text-tertiary">Пока нет сообщений</p>
          ) : (
            messages.map((m) => (
              <div
                key={m.id}
                className={`max-w-2xl rounded-lg px-3 py-2 ${
                  m.role === "agent"
                    ? "ml-auto bg-accent/15"
                    : m.role === "visitor"
                      ? "bg-elevated"
                      : "bg-view/80"
                }`}
              >
                <div className="mb-1 flex items-center justify-between gap-3 text-[11px] text-tertiary">
                  <span>{ROLE_LABEL[m.role] || m.role}</span>
                  <span className="tabular-nums">{formatDateTime(m.createdAt)}</span>
                </div>
                <p className="whitespace-pre-wrap text-[14px] leading-relaxed text-secondary">
                  {m.body}
                </p>
              </div>
            ))
          )}
        </section>

        <footer className="shrink-0 border-t border-border px-5 py-3 md:px-6">
          {error && <p className="mb-2 text-[12px] text-red-500">{error}</p>}
          <div className="flex gap-2">
            <input
              value={body}
              onChange={(e) => setBody(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === "Enter" && !e.shiftKey) {
                  e.preventDefault();
                  send();
                }
              }}
              disabled={pending || conversation.status === "resolved"}
              placeholder="Ответ агента…"
              className="min-w-0 flex-1 rounded-md border border-border bg-base px-3 py-2 text-[13px] text-primary outline-none focus:border-accent"
            />
            <button
              type="button"
              disabled={pending || !body.trim() || conversation.status === "resolved"}
              onClick={send}
              className="rounded-md bg-accent px-4 py-2 text-[13px] font-ui-medium text-white disabled:opacity-50"
            >
              Отправить
            </button>
          </div>
        </footer>
      </div>
    </>
  );
}
