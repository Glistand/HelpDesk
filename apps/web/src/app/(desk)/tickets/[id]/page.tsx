import Link from "next/link";
import { notFound } from "next/navigation";
import { PriorityLabel, SlaBadge, StatusBadge } from "@/components/badges";
import { Timeline } from "@/components/timeline";
import { formatDateTime, formatRelative } from "@/lib/format";
import { getTicket } from "@/lib/mock-data";

export default async function TicketPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;
  const ticket = getTicket(id);
  if (!ticket) notFound();

  return (
    <>
      <header className="shrink-0 border-b border-border px-5 py-4 md:px-6">
        <Link
          href="/inbox"
          className="text-[13px] text-tertiary hover:text-accent md:hidden"
        >
          ← Inbox
        </Link>
        <div className="mt-2 flex flex-wrap items-start justify-between gap-4">
          <div className="min-w-0 space-y-2">
            <div className="flex flex-wrap items-center gap-2">
              <span className="font-mono text-[12px] text-tertiary">
                {ticket.id}
              </span>
              <StatusBadge status={ticket.status} />
              <SlaBadge state={ticket.sla.state} />
            </div>
            <h1 className="font-ui-semibold text-lg leading-snug tracking-tight text-primary md:text-xl">
              {ticket.title}
            </h1>
            <p className="text-[13px] text-tertiary">
              {ticket.requester} · {ticket.category} · {formatRelative(ticket.updatedAt)}
            </p>
          </div>
          <div className="flex gap-2">
            <button
              type="button"
              className="rounded-md px-3 py-1.5 text-[13px] text-tertiary ring-surface transition-colors hover:bg-elevated hover:text-secondary"
            >
              Комментарий
            </button>
            <button
              type="button"
              className="rounded-md bg-accent px-3 py-1.5 text-[13px] font-ui-medium text-white hover:opacity-90"
            >
              Статус
            </button>
          </div>
        </div>
      </header>

      <div className="flex flex-1 overflow-hidden">
        <section className="flex-1 overflow-y-auto px-5 py-5 md:px-6">
          <Block title="Описание">
            <p className="max-w-2xl text-[15px] leading-relaxed text-secondary">
              {ticket.description}
            </p>
          </Block>

          <Block title="Activity">
            <Timeline events={ticket.timeline} />
          </Block>
        </section>

        <aside className="hidden w-64 shrink-0 overflow-y-auto border-l border-border bg-view/50 px-4 py-5 lg:block">
          <Meta label="Приоритет">
            <PriorityLabel priority={ticket.priority} />
          </Meta>
          <Meta label="Исполнитель">
            <span className="text-[13px] text-secondary">
              {ticket.assignee?.name ?? "Не назначен"}
            </span>
          </Meta>
          <Meta label="Первый ответ">
            <span className="font-mono text-[12px] text-secondary">
              {formatDateTime(ticket.sla.firstResponseDue)}
            </span>
          </Meta>
          <Meta label="Resolve">
            <span className="font-mono text-[12px] text-secondary">
              {formatDateTime(ticket.sla.resolveDue)}
            </span>
          </Meta>
          <Meta label="Создан">
            <span className="font-mono text-[12px] text-tertiary">
              {formatDateTime(ticket.createdAt)}
            </span>
          </Meta>
        </aside>
      </div>
    </>
  );
}

function Block({
  title,
  children,
}: {
  title: string;
  children: React.ReactNode;
}) {
  return (
    <div className="mb-8">
      <h2 className="mb-3 text-[11px] font-ui-medium uppercase tracking-widest text-tertiary">
        {title}
      </h2>
      {children}
    </div>
  );
}

function Meta({
  label,
  children,
}: {
  label: string;
  children: React.ReactNode;
}) {
  return (
    <div className="mb-4 border-b border-border pb-4 last:border-0">
      <p className="mb-1 text-[10px] uppercase tracking-widest text-tertiary">
        {label}
      </p>
      {children}
    </div>
  );
}
