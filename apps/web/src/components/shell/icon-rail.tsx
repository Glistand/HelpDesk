"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useEffect, useRef, useState } from "react";
import { useTheme } from "@/components/theme-provider";
import { initialsFromName } from "@/lib/agents";
import type { AuthUser } from "@/lib/types";

const items = [
  {
    href: "/inbox",
    label: "Inbox",
    match: (p: string) => p === "/inbox" || p.startsWith("/conversations/"),
    icon: (
      <path
        d="M4 6h16v2H4V6zm0 5h16v2H4v-2zm0 5h10v2H4v-2z"
        fill="currentColor"
      />
    ),
  },
  {
    href: "/tickets",
    label: "Тикеты",
    match: (p: string) =>
      p === "/tickets" || (p.startsWith("/tickets/") && p !== "/tickets/new"),
    icon: (
      <path
        d="M4 6h16v12H4V6zm2 3v2h12V9H6zm0 4v2h8v-2H6z"
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
    href: "/team",
    label: "Команда",
    adminOnly: true,
    match: (p: string) => p === "/team",
    icon: (
      <path
        d="M12 12a4 4 0 1 0-4-4 4 4 0 0 0 4 4zm0 2c-4 0-8 2-8 4v1h16v-1c0-2-4-4-8-4z"
        fill="currentColor"
      />
    ),
  },
];

const COLLAPSE_KEY = "hd-nav-expanded";

export function IconRail({ user }: { user: AuthUser | null }) {
  const pathname = usePathname();
  const { theme, toggleTheme } = useTheme();
  const [expanded, setExpanded] = useState(false);
  const [menuOpen, setMenuOpen] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    try {
      setExpanded(localStorage.getItem(COLLAPSE_KEY) === "1");
    } catch {
      /* ignore */
    }
  }, []);

  useEffect(() => {
    function onDoc(e: MouseEvent) {
      if (!menuRef.current?.contains(e.target as Node)) setMenuOpen(false);
    }
    document.addEventListener("mousedown", onDoc);
    return () => document.removeEventListener("mousedown", onDoc);
  }, []);

  function toggleExpanded() {
    setExpanded((v) => {
      const next = !v;
      try {
        localStorage.setItem(COLLAPSE_KEY, next ? "1" : "0");
      } catch {
        /* ignore */
      }
      return next;
    });
  }

  const initials = user ? initialsFromName(user.name) : "?";
  const isAdmin = user?.role === "admin";
  const navItems = items.filter((i) => !i.adminOnly || isAdmin);

  async function logout() {
    await fetch("/api/auth/logout", { method: "POST" });
    window.location.href = "/login";
  }

  return (
    <nav
      aria-label="Main"
      className={`flex shrink-0 flex-col border-r border-border bg-frame py-3 transition-[width] ${
        expanded ? "w-52" : "w-14"
      }`}
    >
      <div className={`mb-4 flex items-center ${expanded ? "px-3 gap-2" : "justify-center"}`}>
        <Link
          href="/inbox"
          className="flex h-8 w-8 shrink-0 items-center justify-center rounded-md bg-accent font-ui-semibold text-sm text-white"
          title="Support Desk"
        >
          S
        </Link>
        {expanded && (
          <span className="truncate font-ui-semibold text-sm text-primary">Support Desk</span>
        )}
      </div>

      <button
        type="button"
        onClick={toggleExpanded}
        title={expanded ? "Свернуть" : "Развернуть"}
        aria-label={expanded ? "Свернуть меню" : "Развернуть меню"}
        className={`mb-3 flex h-8 items-center rounded-md text-tertiary transition-colors hover:bg-elevated hover:text-secondary ${
          expanded ? "mx-2 gap-2 px-2" : "mx-auto w-9 justify-center"
        }`}
      >
        <svg viewBox="0 0 24 24" className="h-4 w-4 shrink-0" aria-hidden>
          <path
            fill="currentColor"
            d={
              expanded
                ? "M15.41 7.41 14 6l-6 6 6 6 1.41-1.41L10.83 12z"
                : "M10 6 8.59 7.41 13.17 12l-4.58 4.59L10 18l6-6z"
            }
          />
        </svg>
        {expanded && <span className="text-xs">Свернуть</span>}
      </button>

      <div className={`flex flex-1 flex-col gap-1 ${expanded ? "px-2" : "items-center"}`}>
        {navItems.map((item) => {
          const active = item.match(pathname);
          return (
            <Link
              key={item.href}
              href={item.href}
              title={item.label}
              className={`flex h-9 items-center rounded-md transition-colors ${
                expanded ? "gap-2.5 px-2.5" : "w-9 justify-center"
              } ${
                active
                  ? "bg-elevated text-primary"
                  : "text-tertiary hover:bg-elevated/80 hover:text-secondary"
              }`}
            >
              <svg
                viewBox="0 0 24 24"
                className="h-[18px] w-[18px] shrink-0"
                aria-hidden
              >
                {item.icon}
              </svg>
              {expanded && (
                <span className="truncate text-[13px] font-ui-medium">{item.label}</span>
              )}
            </Link>
          );
        })}
      </div>

      <div className={`mt-auto flex flex-col gap-2 ${expanded ? "px-2" : "items-center"}`}>
        <button
          type="button"
          onClick={toggleTheme}
          title={theme === "dark" ? "Светлая тема" : "Тёмная тема"}
          aria-label={theme === "dark" ? "Включить светлую тему" : "Включить тёмную тему"}
          className={`flex h-9 items-center rounded-md text-tertiary transition-colors hover:bg-elevated hover:text-secondary ${
            expanded ? "gap-2.5 px-2.5" : "w-9 justify-center"
          }`}
        >
          {theme === "dark" ? (
            <svg viewBox="0 0 24 24" className="h-[18px] w-[18px] shrink-0" aria-hidden>
              <path
                fill="currentColor"
                d="M12 4a1 1 0 0 1 1 1v1a1 1 0 1 1-2 0V5a1 1 0 0 1 1-1zm0 12a4 4 0 1 0 0-8 4 4 0 0 0 0 8zm7-5a1 1 0 0 1 1 1 1 1 0 0 1-1 1h-1a1 1 0 1 1 0-2h1zM6 12a1 1 0 0 1-1 1H4a1 1 0 1 1 0-2h1a1 1 0 0 1 1 1zm10.95 5.536a1 1 0 0 1 0 1.414l-.707.707a1 1 0 1 1-1.414-1.414l.707-.707a1 1 0 0 1 1.414 0zM8.464 6.05a1 1 0 0 1 0 1.414l-.707.707A1 1 0 1 1 6.343 6.757l.707-.707a1 1 0 0 1 1.414 0zm9.192 0a1 1 0 0 1 1.414 0l.707.707a1 1 0 1 1-1.414 1.414l-.707-.707a1 1 0 0 1 0-1.414zM8.464 17.95a1 1 0 0 1-1.414 0l-.707-.707a1 1 0 1 1 1.414-1.414l.707.707a1 1 0 0 1 0 1.414zM12 17a1 1 0 0 1 1 1v1a1 1 0 1 1-2 0v-1a1 1 0 0 1 1-1z"
              />
            </svg>
          ) : (
            <svg viewBox="0 0 24 24" className="h-[18px] w-[18px] shrink-0" aria-hidden>
              <path
                fill="currentColor"
                d="M12.1 22c-4.6 0-8.4-3.4-9-7.9-.1-.8.6-1.5 1.4-1.3 3.4.8 6.9-.9 8.4-3.9 1.5-3 1-6.6-1.1-9.1-.4-.5 0-1.3.6-1.4C18.5-.7 24 4.2 24 10.8 24 17 18.5 22 12.1 22z"
              />
            </svg>
          )}
          {expanded && (
            <span className="text-[13px]">{theme === "dark" ? "Светлая" : "Тёмная"}</span>
          )}
        </button>

        <div className="relative" ref={menuRef}>
          <button
            type="button"
            title={user?.name || "Профиль"}
            aria-haspopup="menu"
            aria-expanded={menuOpen}
            onClick={() => setMenuOpen((o) => !o)}
            className={`flex items-center rounded-md transition-colors hover:bg-elevated ${
              expanded ? "w-full gap-2.5 px-2.5 py-1.5" : "mx-auto h-8 w-8 justify-center"
            }`}
          >
            <span className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-accent-muted text-[11px] font-ui-semibold text-accent">
              {initials}
            </span>
            {expanded && (
              <span className="min-w-0 flex-1 text-left">
                <span className="block truncate text-[13px] font-ui-medium text-primary">
                  {user?.name || "—"}
                </span>
                <span className="block truncate text-[11px] text-tertiary">{user?.role}</span>
              </span>
            )}
          </button>

          {menuOpen && (
            <div
              role="menu"
              className="absolute bottom-full left-0 z-50 mb-2 w-52 overflow-hidden rounded-md border border-border bg-elevated shadow-lg"
            >
              <div className="border-b border-border px-3 py-2">
                <p className="truncate text-[13px] font-ui-medium text-primary">
                  {user?.name}
                </p>
                <p className="truncate text-[11px] text-tertiary">{user?.email}</p>
              </div>
              <Link
                href="/profile"
                role="menuitem"
                onClick={() => setMenuOpen(false)}
                className="block px-3 py-2 text-[13px] text-secondary hover:bg-base"
              >
                Настройки профиля
              </Link>
              {isAdmin && (
                <Link
                  href="/team"
                  role="menuitem"
                  onClick={() => setMenuOpen(false)}
                  className="block px-3 py-2 text-[13px] text-secondary hover:bg-base"
                >
                  Команда и аккаунты
                </Link>
              )}
              <button
                type="button"
                role="menuitem"
                onClick={logout}
                className="block w-full px-3 py-2 text-left text-[13px] text-red-500 hover:bg-base"
              >
                Выйти
              </button>
            </div>
          )}
        </div>
      </div>
    </nav>
  );
}
