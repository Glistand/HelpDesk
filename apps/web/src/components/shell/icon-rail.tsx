"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useTheme } from "@/components/theme-provider";

const items = [
  {
    href: "/inbox",
    label: "Inbox",
    match: (p: string) => p === "/inbox" || p.startsWith("/tickets/"),
    icon: (
      <path
        d="M4 6h16v2H4V6zm0 5h16v2H4v-2zm0 5h10v2H4v-2z"
        fill="currentColor"
      />
    ),
  },
  {
    href: "/tickets/new",
    label: "Создать",
    match: (p: string) => p === "/tickets/new",
    icon: (
      <path d="M11 5h2v14h-2V5zm-7 7h14v2H4v-2z" fill="currentColor" />
    ),
  },
  {
    href: "/design",
    label: "Design",
    match: (p: string) => p === "/design",
    icon: (
      <path
        d="M12 2l2.4 7.4H22l-6 4.6 2.3 7-6.3-4.6L5.7 21l2.3-7-6-4.6h7.6L12 2z"
        fill="currentColor"
      />
    ),
  },
];

export function IconRail() {
  const pathname = usePathname();
  const { theme, toggleTheme } = useTheme();

  return (
    <nav
      aria-label="Main"
      className="flex w-14 shrink-0 flex-col items-center border-r border-border bg-frame py-3"
    >
      <Link
        href="/inbox"
        className="mb-6 flex h-8 w-8 items-center justify-center rounded-md bg-accent font-ui-semibold text-sm text-white"
        title="Helpdesk"
      >
        H
      </Link>

      <div className="flex flex-1 flex-col gap-1">
        {items.map((item) => {
          const active = item.match(pathname);
          return (
            <Link
              key={item.href}
              href={item.href}
              title={item.label}
              className={`flex h-9 w-9 items-center justify-center rounded-md transition-colors ${
                active
                  ? "bg-elevated text-primary"
                  : "text-tertiary hover:bg-elevated/80 hover:text-secondary"
              }`}
            >
              <svg
                viewBox="0 0 24 24"
                className="h-[18px] w-[18px]"
                aria-hidden
              >
                {item.icon}
              </svg>
            </Link>
          );
        })}
      </div>

      <div className="mb-3 flex flex-col items-center gap-2">
        <button
          type="button"
          onClick={toggleTheme}
          title={theme === "dark" ? "Светлая тема" : "Тёмная тема"}
          aria-label={theme === "dark" ? "Включить светлую тему" : "Включить тёмную тему"}
          className="flex h-9 w-9 items-center justify-center rounded-md text-tertiary transition-colors hover:bg-elevated hover:text-secondary"
        >
          {theme === "dark" ? (
            <svg viewBox="0 0 24 24" className="h-[18px] w-[18px]" aria-hidden>
              <path
                fill="currentColor"
                d="M12 4a1 1 0 0 1 1 1v1a1 1 0 1 1-2 0V5a1 1 0 0 1 1-1zm0 12a4 4 0 1 0 0-8 4 4 0 0 0 0 8zm7-5a1 1 0 0 1 1 1 1 1 0 0 1-1 1h-1a1 1 0 1 1 0-2h1zM6 12a1 1 0 0 1-1 1H4a1 1 0 1 1 0-2h1a1 1 0 0 1 1 1zm10.95 5.536a1 1 0 0 1 0 1.414l-.707.707a1 1 0 1 1-1.414-1.414l.707-.707a1 1 0 0 1 1.414 0zM8.464 6.05a1 1 0 0 1 0 1.414l-.707.707A1 1 0 1 1 6.343 6.757l.707-.707a1 1 0 0 1 1.414 0zm9.192 0a1 1 0 0 1 1.414 0l.707.707a1 1 0 1 1-1.414 1.414l-.707-.707a1 1 0 0 1 0-1.414zM8.464 17.95a1 1 0 0 1-1.414 0l-.707-.707a1 1 0 1 1 1.414-1.414l.707.707a1 1 0 0 1 0 1.414zM12 17a1 1 0 0 1 1 1v1a1 1 0 1 1-2 0v-1a1 1 0 0 1 1-1z"
              />
            </svg>
          ) : (
            <svg viewBox="0 0 24 24" className="h-[18px] w-[18px]" aria-hidden>
              <path
                fill="currentColor"
                d="M12.1 22c-4.6 0-8.4-3.4-9-7.9-.1-.8.6-1.5 1.4-1.3 3.4.8 6.9-.9 8.4-3.9 1.5-3 1-6.6-1.1-9.1-.4-.5 0-1.3.6-1.4C18.5-.7 24 4.2 24 10.8 24 17 18.5 22 12.1 22z"
              />
            </svg>
          )}
        </button>
        <div
          className="flex h-8 w-8 items-center justify-center rounded-full bg-accent-muted text-[11px] font-ui-semibold text-accent"
          title="Алексей К."
        >
          АК
        </div>
      </div>
    </nav>
  );
}
