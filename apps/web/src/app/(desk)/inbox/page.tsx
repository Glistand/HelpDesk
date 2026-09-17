import Link from "next/link";
import { listTickets } from "@/lib/api";

export default async function InboxPage() {
  let count = 0;
  try {
    const tickets = await listTickets();
    count = tickets.filter((t) => !["resolved", "closed"].includes(t.status)).length;
  } catch {
    count = 0;
  }

  return (
    <div className="flex flex-1 flex-col items-center justify-center px-8 text-center">
      <div className="max-w-sm">
        <p className="text-[11px] font-ui-medium uppercase tracking-widest text-tertiary">
          Helpdesk Event Hub
        </p>
        <h2 className="mt-2 font-ui-semibold text-xl tracking-tight text-primary">
          Выберите тикет
        </h2>
        <p className="mt-2 text-[13px] leading-relaxed text-secondary">
          Список слева · детали справа. Открытых:{" "}
          <span className="tabular-nums text-primary">{count}</span>
        </p>
        <div className="mt-6 flex flex-col gap-2 sm:flex-row sm:justify-center">
          <Link
            href="/tickets/new"
            className="rounded-md bg-accent px-4 py-2 text-[13px] font-ui-medium text-white transition-opacity hover:opacity-90"
          >
            Создать тикет
          </Link>
        </div>
      </div>
      <p className="mt-12 text-[11px] text-tertiary">
        Live API · gateway :8080 · Meilisearch search
      </p>
    </div>
  );
}
