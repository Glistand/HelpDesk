import type { CSSProperties } from "react";

export type Theme = "dark" | "light";

const colorVariables = {
  WEB_ACCENT_COLOR: "--web-accent",
  WEB_ACCENT_MUTED_COLOR: "--web-accent-muted",
  WEB_DARK_BG_BASE: "--web-dark-bg-base",
  WEB_DARK_BG_FRAME: "--web-dark-bg-frame",
  WEB_DARK_BG_ELEVATED: "--web-dark-bg-elevated",
  WEB_DARK_TEXT_PRIMARY: "--web-dark-text-primary",
  WEB_DARK_TEXT_SECONDARY: "--web-dark-text-secondary",
  WEB_LIGHT_BG_BASE: "--web-light-bg-base",
  WEB_LIGHT_BG_FRAME: "--web-light-bg-frame",
  WEB_LIGHT_BG_ELEVATED: "--web-light-bg-elevated",
  WEB_LIGHT_TEXT_PRIMARY: "--web-light-text-primary",
  WEB_LIGHT_TEXT_SECONDARY: "--web-light-text-secondary",
} as const;

const HEX_COLOR = /^#[0-9a-f]{3,8}$/i;

function value(name: string): string | undefined {
  const raw = process.env[name]?.trim();
  return raw && HEX_COLOR.test(raw) ? raw : undefined;
}

export function branding() {
  const defaultTheme: Theme =
    process.env.WEB_DEFAULT_THEME === "light" ? "light" : "dark";
  const name = process.env.WEB_BRAND_NAME?.trim() || "Support Desk";
  const style: CSSProperties & Record<`--${string}`, string> = {};

  for (const [env, variable] of Object.entries(colorVariables)) {
    const color = value(env);
    if (color) style[variable] = color;
  }

  return { defaultTheme, name, style };
}
