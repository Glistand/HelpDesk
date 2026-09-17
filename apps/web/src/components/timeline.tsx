import { formatDateTime } from "@/lib/format";
import type { TimelineEvent, TimelineKind } from "@/lib/types";

const kindMeta: Record<TimelineKind, { label: string; color: string }> = {
  created: { label: "Created", color: "bg-accent" },
  assigned: { label: "Assigned", color: "bg-status-violet" },
  status: { label: "Status", color: "bg-tertiary" },
  comment: { label: "Comment", color: "bg-status-green" },
  sla_warn: { label: "SLA", color: "bg-status-amber" },
  sla_breach: { label: "Breach", color: "bg-status-red" },
  escalated: { label: "Escalated", color: "bg-status-red" },
  notification: { label: "Notify", color: "bg-accent/70" },
  other: { label: "Event", color: "bg-tertiary" },
};

export function Timeline({ events }: { events: TimelineEvent[] }) {
  return (
    <div className="space-y-0">
      {events.map((event, index) => {
        const meta = kindMeta[event.kind];
        const isLast = index === events.length - 1;

        return (
          <div key={event.id} className="relative flex gap-3 pb-5">
            {!isLast && (
              <span
                aria-hidden
                className="absolute left-[5px] top-3 h-[calc(100%-8px)] w-px bg-border"
              />
            )}
            <span
              aria-hidden
              className={`relative z-10 mt-1 h-2.5 w-2.5 shrink-0 rounded-full ${meta.color}`}
            />
            <div className="min-w-0 flex-1 pt-0.5">
              <div className="flex flex-wrap items-baseline gap-x-2">
                <span className="text-[13px] font-ui-medium text-primary">
                  {event.title}
                </span>
                <span className="text-[10px] uppercase tracking-wider text-tertiary">
                  {meta.label}
                </span>
              </div>
              {event.detail && (
                <p className="mt-0.5 text-[13px] leading-relaxed text-secondary">
                  {event.detail}
                </p>
              )}
              <p className="mt-1 text-[11px] text-tertiary">
                {event.actor && `${event.actor} · `}
                {formatDateTime(event.at)}
              </p>
            </div>
          </div>
        );
      })}
    </div>
  );
}
