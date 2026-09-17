"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { FormEvent, useState } from "react";

export default function NewTicketPage() {
  const router = useRouter();
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function onSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setLoading(true);
    setError("");
    const fd = new FormData(e.currentTarget);
    const body = {
      title: String(fd.get("title") || ""),
      description: String(fd.get("description") || ""),
      category: String(fd.get("category") || ""),
      priority: String(fd.get("priority") || "normal"),
    };
    try {
      const res = await fetch("/api/hd/tickets", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(body),
      });
      if (!res.ok) {
        const data = await res.json().catch(() => ({}));
        setError(data.error || "Не удалось создать тикет");
        return;
      }
      const ticket = await res.json();
      router.push(`/tickets/${ticket.id}`);
      router.refresh();
    } catch {
      setError("Нет связи с API");
    } finally {
      setLoading(false);
    }
  }

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
          Отправка в ticket-service через gateway
        </p>
      </header>

      <div className="mx-auto w-full max-w-lg px-5 py-8 md:px-6">
        <form className="space-y-4" onSubmit={onSubmit}>
          <Field label="Тема" required>
            <input
              name="title"
              required
              placeholder="Кратко опишите проблему"
              className="field-input"
            />
          </Field>
          <Field label="Категория">
            <select name="category" className="field-input" defaultValue="Сеть / VPN">
              <option>Сеть / VPN</option>
              <option>Почта</option>
              <option>Доступы</option>
              <option>Оборудование</option>
            </select>
          </Field>
          <Field label="Приоритет">
            <select name="priority" className="field-input" defaultValue="normal">
              <option value="low">Низкий</option>
              <option value="normal">Обычный</option>
              <option value="high">Высокий</option>
              <option value="urgent">Срочный</option>
            </select>
          </Field>
          <Field label="Описание" required>
            <textarea
              name="description"
              required
              rows={5}
              placeholder="Шаги воспроизведения, что уже пробовали"
              className="field-input resize-y"
            />
          </Field>
          {error && <p className="text-[13px] text-status-red">{error}</p>}
          <div className="flex gap-2 pt-2">
            <button
              type="submit"
              disabled={loading}
              className="rounded-md bg-accent px-4 py-2 text-[13px] font-ui-medium text-white hover:opacity-90 disabled:opacity-60"
            >
              {loading ? "Создание…" : "Создать"}
            </button>
            <Link
              href="/inbox"
              className="rounded-md px-4 py-2 text-[13px] text-tertiary ring-surface hover:bg-elevated"
            >
              Отмена
            </Link>
          </div>
        </form>
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
