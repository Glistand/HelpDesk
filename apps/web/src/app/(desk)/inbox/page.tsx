import Link from "next/link";
import { listConversations } from "@/lib/api";

export default async function InboxPage() {
  let waiting = 0;
  let active = 0;
  try {
    const conversations = await listConversations();
    waiting = conversations.filter((c) => c.status === "waiting_agent").length;
    active = conversations.filter((c) => c.status !== "resolved").length;
  } catch {
    waiting = 0;
    active = 0;
  }

  return (
    <div className="flex flex-1 flex-col items-center justify-center px-8 text-center">
      <div className="max-w-sm">
        <p className="text-[11px] font-ui-medium uppercase tracking-widest text-tertiary">
          Support Desk
        </p>
        <h2 className="mt-2 font-ui-semibold text-xl tracking-tight text-primary">
          Выберите беседу
        </h2>
        <p className="mt-2 text-[13px] leading-relaxed text-secondary">
          Inbox слева · тред справа. Активных:{" "}
          <span className="tabular-nums text-primary">{active}</span>
          {" · "}handoff:{" "}
          <span className="tabular-nums text-primary">{waiting}</span>
        </p>
        <div className="mt-6 flex flex-col gap-2 sm:flex-row sm:justify-center">
          <a
            href="/widget-demo.html"
            target="_blank"
            rel="noreferrer"
            className="rounded-md bg-accent px-4 py-2 text-[13px] font-ui-medium text-white transition-opacity hover:opacity-90"
          >
            Открыть виджет-демо
          </a>
          <Link
            href="/tickets/new"
            className="rounded-md border border-border px-4 py-2 text-[13px] font-ui-medium text-secondary hover:bg-elevated"
          >
            Тикеты (legacy)
          </Link>
        </div>
      </div>
      <p className="mt-12 text-[11px] text-tertiary">
        Widget → OpenRouter bot → handoff → agent reply
      </p>
    </div>
  );
}
