import Link from "next/link";

export default function TicketsPage() {
  return (
    <div className="flex flex-1 flex-col items-center justify-center px-8 text-center">
      <div className="max-w-sm">
        <p className="text-[11px] font-ui-medium uppercase tracking-widest text-tertiary">
          Support Desk
        </p>
        <h1 className="mt-2 font-ui-semibold text-xl tracking-tight text-primary">
          Выберите тикет
        </h1>
        <p className="mt-2 text-[13px] leading-relaxed text-secondary">
          Используйте список слева для поиска, фильтрации и открытия тикетов.
        </p>
        <Link
          href="/tickets/new"
          className="mt-6 inline-flex rounded-md bg-accent px-4 py-2 text-[13px] font-ui-medium text-white transition-opacity hover:opacity-90"
        >
          Создать тикет
        </Link>
      </div>
    </div>
  );
}
