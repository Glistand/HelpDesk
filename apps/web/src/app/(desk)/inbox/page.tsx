import Link from "next/link";

export default function InboxPage() {
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
          Список слева · детали справа. Как в Linear и Amie inbox — master-detail
          layout для агента.
        </p>
        <div className="mt-6 flex flex-col gap-2 sm:flex-row sm:justify-center">
          <Link
            href="/tickets/t-1042"
            className="rounded-md bg-accent px-4 py-2 text-[13px] font-ui-medium text-white transition-opacity hover:opacity-90"
          >
            Открыть t-1042
          </Link>
          <Link
            href="/tickets/t-1035"
            className="rounded-md px-4 py-2 text-[13px] text-tertiary ring-surface transition-colors hover:bg-elevated hover:text-secondary"
          >
            SLA breach
          </Link>
        </div>
      </div>
      <p className="mt-12 text-[11px] text-tertiary">
        Mock-данные · preview UI · Sparkbites refs: linear.app, amie.so, attio.com
      </p>
    </div>
  );
}
