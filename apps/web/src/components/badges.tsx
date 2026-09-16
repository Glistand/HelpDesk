import type { SlaState, TicketPriority, TicketStatus } from "@/lib/types";

const statusLabels: Record<TicketStatus, string> = {
  new: "Новый",
  open: "Открыт",
  pending: "Ожидание",
  resolved: "Решён",
  closed: "Закрыт",
};

const priorityColors: Record<TicketPriority, string> = {
  low: "bg-tertiary",
  normal: "bg-secondary",
  high: "bg-status-amber",
  urgent: "bg-status-red",
};

const slaColors: Record<SlaState, string> = {
  ok: "bg-status-green",
  warning: "bg-status-amber",
  breached: "bg-status-red",
};

export function PriorityDot({ priority }: { priority: TicketPriority }) {
  return (
    <span
      className={`mt-1.5 h-2 w-2 shrink-0 rounded-full ${priorityColors[priority]}`}
      title={priority}
    />
  );
}

export function SlaDot({ state }: { state: SlaState }) {
  if (state === "ok") return null;
  return (
    <span
      className={`h-1.5 w-1.5 rounded-full ${slaColors[state]} ${
        state === "warning" ? "animate-pulse" : ""
      }`}
      title={state === "breached" ? "SLA нарушен" : "SLA скоро"}
    />
  );
}

export function StatusBadge({ status }: { status: TicketStatus }) {
  return (
    <span className="inline-flex items-center rounded-sm bg-elevated px-1.5 py-0.5 text-[11px] font-ui-medium text-secondary ring-surface">
      {statusLabels[status]}
    </span>
  );
}

export function SlaBadge({ state }: { state: SlaState }) {
  const labels: Record<SlaState, string> = {
    ok: "SLA OK",
    warning: "SLA · скоро",
    breached: "SLA · breach",
  };
  return (
    <span
      className={`inline-flex items-center gap-1.5 rounded-sm px-1.5 py-0.5 text-[11px] font-ui-medium ring-surface ${
        state === "breached"
          ? "bg-status-red/10 text-status-red"
          : state === "warning"
            ? "bg-status-amber/10 text-status-amber"
            : "bg-status-green/10 text-status-green"
      }`}
    >
      <span className={`h-1.5 w-1.5 rounded-full ${slaColors[state]}`} />
      {labels[state]}
    </span>
  );
}

export function PriorityLabel({ priority }: { priority: TicketPriority }) {
  const labels: Record<TicketPriority, string> = {
    low: "Низкий",
    normal: "Обычный",
    high: "Высокий",
    urgent: "Срочный",
  };
  return (
    <span className="text-[13px] text-secondary">{labels[priority]}</span>
  );
}
