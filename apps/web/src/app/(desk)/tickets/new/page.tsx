"use client";

import Link from "next/link";
import { useState } from "react";

export default function NewTicketPage() {
  const [done, setDone] = useState(false);

  return (
    <div className="flex flex-1 flex-col overflow-y-auto">
      <header className="border-b border-border px-5 py-4 md:px-6">
        <Link
          href="/inbox"
          className="text-[13px] text-tertiary hover:text-accent"
        >
          ← Inbox
        </Link>
        <h1 className="mt-2 font-ui-semibold text-xl text-primary">
          Создать тикет
        </h1>
        <p className="mt-1 text-[13px] text-tertiary">
          Preview формы · без сохранения
        </p>
      </header>

      <div className="mx-auto w-full max-w-lg px-5 py-8 md:px-6">
        {done ? (
          <div className="rounded-lg bg-status-green/10 p-6 text-center ring-surface">
            <p className="font-ui-semibold text-primary">Тикет создан (mock)</p>
            <p className="mt-2 text-[13px] text-secondary">
              Событие{" "}
              <code className="font-mono text-[12px] text-accent">
                helpdesk.ticket.created
              </code>
            </p>
            <Link
              href="/tickets/t-1042"
              className="mt-4 inline-block text-[13px] text-accent hover:underline"
            >
              Открыть t-1042 →
            </Link>
          </div>
        ) : (
          <form
            className="space-y-4"
            onSubmit={(e) => {
              e.preventDefault();
              setDone(true);
            }}
          >
            <Field label="Тема" required>
              <input
                required
                placeholder="Кратко опишите проблему"
                className="field-input"
              />
            </Field>
            <Field label="Категория">
              <select className="field-input">
                <option>Сеть / VPN</option>
                <option>Почта</option>
                <option>Доступы</option>
                <option>Оборудование</option>
              </select>
            </Field>
            <Field label="Приоритет">
              <select className="field-input">
                <option>Низкий</option>
                <option>Обычный</option>
                <option>Высокий</option>
                <option>Срочный</option>
              </select>
            </Field>
            <Field label="Описание" required>
              <textarea
                required
                rows={5}
                placeholder="Шаги воспроизведения, что уже пробовали"
                className="field-input resize-y"
              />
            </Field>
            <div className="flex gap-2 pt-2">
              <button
                type="submit"
                className="rounded-md bg-accent px-4 py-2 text-[13px] font-ui-medium text-white hover:opacity-90"
              >
                Создать
              </button>
              <Link
                href="/inbox"
                className="rounded-md px-4 py-2 text-[13px] text-tertiary ring-surface hover:bg-elevated"
              >
                Отмена
              </Link>
            </div>
          </form>
        )}
      </div>
    </div>
  );
}

function Field({
  label,
  required,
  children,
}: {
  label: string;
  required?: boolean;
  children: React.ReactNode;
}) {
  return (
    <label className="block space-y-1.5">
      <span className="text-[13px] font-ui-medium text-secondary">
        {label}
        {required && <span className="text-status-red"> *</span>}
      </span>
      {children}
    </label>
  );
}
