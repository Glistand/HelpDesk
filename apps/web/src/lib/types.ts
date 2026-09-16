export type TicketStatus =
  | "new"
  | "open"
  | "pending"
  | "resolved"
  | "closed";

export type TicketPriority = "low" | "normal" | "high" | "urgent";

export type SlaState = "ok" | "warning" | "breached";

export type TimelineKind =
  | "created"
  | "assigned"
  | "status"
  | "comment"
  | "sla_warn"
  | "sla_breach"
  | "escalated"
  | "notification";

export interface Agent {
  id: string;
  name: string;
  initials: string;
  role: string;
}

export interface TimelineEvent {
  id: string;
  kind: TimelineKind;
  title: string;
  detail?: string;
  actor?: string;
  at: string;
}

export interface Ticket {
  id: string;
  title: string;
  description: string;
  status: TicketStatus;
  priority: TicketPriority;
  category: string;
  requester: string;
  assignee?: Agent;
  createdAt: string;
  updatedAt: string;
  sla: {
    firstResponseDue: string;
    resolveDue: string;
    state: SlaState;
  };
  timeline: TimelineEvent[];
}
