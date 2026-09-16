import { PriorityLabel, SlaBadge, StatusBadge } from "@/components/badges";

const refs = [
  { name: "linear.app", note: "Dark canvas, Inter 510/590, indigo accent, ring shadows" },
  { name: "amie.so", note: "Icon rail + inbox panel widths, app-like density" },
  { name: "attio.com", note: "Achromatic precision, border-driven states" },
  { name: "boardui.com", note: "Dashboard component patterns (Sparkbites bookmark)" },
];

const darkSwatches = [
  { name: "Base", hex: "#08090a" },
  { name: "Elevated", hex: "#0f1011" },
  { name: "Card", hex: "#151617" },
  { name: "Primary text", hex: "#f7f8f8" },
  { name: "Secondary", hex: "#d0d6e0" },
  { name: "Tertiary", hex: "#8a8f98" },
  { name: "Accent", hex: "#5e6ad2" },
  { name: "Border", hex: "rgba(255,255,255,0.08)" },
];

const lightSwatches = [
  { name: "Base", hex: "#ffffff" },
  { name: "Frame", hex: "#fafafa" },
  { name: "Elevated", hex: "#f4f4f5" },
  { name: "Primary text", hex: "#171717" },
  { name: "Secondary", hex: "#52525b" },
  { name: "Tertiary", hex: "#a1a1aa" },
  { name: "Accent", hex: "#5e6ad2" },
  { name: "Border", hex: "#e5e5e5" },
];

export default function DesignPage() {
  return (
    <div className="flex-1 overflow-y-auto">
      <header className="border-b border-border px-6 py-5">
        <p className="text-[11px] font-ui-medium uppercase tracking-widest text-tertiary">
          Sparkbites · frontend-design skill
        </p>
        <h1 className="mt-1 font-ui-semibold text-xl text-primary">
          Design system preview
        </h1>
        <p className="mt-2 max-w-xl text-[13px] leading-relaxed text-secondary">
          Dark — Linear. Light — Attio. Layout как у Amie. Переключатель темы —
          иконка солнца/луны внизу левого rail.
        </p>
      </header>

      <div className="space-y-8 px-6 py-8">
        <section>
          <SectionTitle>References (Sparkbites)</SectionTitle>
          <ul className="space-y-2">
            {refs.map((r) => (
              <li
                key={r.name}
                className="rounded-lg bg-card px-4 py-3 text-[13px] ring-surface"
              >
                <span className="font-mono text-accent">{r.name}</span>
                <span className="text-tertiary"> — {r.note}</span>
              </li>
            ))}
          </ul>
        </section>

        <section>
          <SectionTitle>Layout</SectionTitle>
          <pre className="overflow-x-auto rounded-lg bg-card p-4 font-mono text-[12px] leading-relaxed text-secondary ring-surface">
{`┌──┬──────────────┬─────────────────────────┐
│56│    340px     │        flex-1           │
│px│  ticket list │   detail / empty        │
│  │  (Amie ref)  │   (Linear surfaces)     │
└──┴──────────────┴─────────────────────────┘`}
          </pre>
        </section>

        <section>
          <SectionTitle>Palette · dark</SectionTitle>
          <SwatchGrid swatches={darkSwatches} />
        </section>

        <section>
          <SectionTitle>Palette · light</SectionTitle>
          <SwatchGrid swatches={lightSwatches} />
        </section>

        <section>
          <SectionTitle>Typography</SectionTitle>
          <div className="space-y-3 rounded-lg bg-card p-5 ring-surface">
            <p className="font-ui-semibold text-xl tracking-tight text-primary">
              Inter 590 — заголовки
            </p>
            <p className="text-[15px] text-secondary">
              Inter 400 — body · 15px для описаний
            </p>
            <p className="font-ui-medium text-[13px] text-secondary">
              Inter 510 — строки списка, ссылки
            </p>
            <p className="font-mono text-[12px] text-tertiary">
              mono — t-1042 · timestamps
            </p>
          </div>
        </section>

        <section>
          <SectionTitle>Components</SectionTitle>
          <div className="flex flex-wrap gap-2 rounded-lg bg-card p-4 ring-surface">
            <StatusBadge status="open" />
            <StatusBadge status="pending" />
            <SlaBadge state="warning" />
            <SlaBadge state="breached" />
          </div>
          <div className="mt-3 flex flex-wrap gap-4 rounded-lg bg-card p-4 ring-surface">
            <PriorityLabel priority="normal" />
            <PriorityLabel priority="high" />
            <PriorityLabel priority="urgent" />
          </div>
          <div className="mt-3 flex flex-wrap gap-2">
            <button
              type="button"
              className="rounded-md bg-accent px-4 py-2 text-[13px] font-ui-medium text-white"
            >
              Primary
            </button>
            <button
              type="button"
              className="rounded-md px-4 py-2 text-[13px] text-tertiary ring-surface hover:bg-elevated"
            >
              Ghost
            </button>
          </div>
        </section>
      </div>
    </div>
  );
}

function SectionTitle({ children }: { children: React.ReactNode }) {
  return (
    <h2 className="mb-3 text-[11px] font-ui-medium uppercase tracking-widest text-tertiary">
      {children}
    </h2>
  );
}

function SwatchGrid({
  swatches,
}: {
  swatches: { name: string; hex: string }[];
}) {
  return (
    <div className="grid grid-cols-2 gap-2 sm:grid-cols-4">
      {swatches.map((s) => (
        <div key={s.name} className="overflow-hidden rounded-lg ring-surface">
          <div
            className="h-12 border-b border-border"
            style={{ background: s.hex }}
          />
          <div className="bg-card px-3 py-2">
            <p className="text-[13px] text-primary">{s.name}</p>
            <p className="font-mono text-[11px] text-tertiary">{s.hex}</p>
          </div>
        </div>
      ))}
    </div>
  );
}
