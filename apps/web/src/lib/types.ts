export type TicketStatus =
  | "new"
  | "open"
  | "pending"
  | "resolved"
  | "closed";

export type TicketPriority = "low" | "normal" | "high" | "urgent";

export type SlaState = "ok" | "warning" | "breached" | "cancelled";

export type TimelineKind =
  | "created"
  | "assigned"
  | "status"
  | "comment"
  | "sla_warn"
  | "sla_breach"
  | "escalated"
  | "notification"
  | "other";

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

export interface TicketSLA {
  firstResponseDue: string;
  resolveDue: string;
  state: SlaState;
  warnedAt?: string;
  breachedAt?: string;
  policy?: string;
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
  sla: TicketSLA;
  timeline: TimelineEvent[];
}

export interface TicketSummary {
  id: string;
  title: string;
  status: TicketStatus;
  priority: TicketPriority;
  category: string;
  requester: string;
  assigneeId: string;
  assignee?: Agent;
  createdAt: string;
  updatedAt: string;
  slaState: SlaState;
}

export interface AuthUser {
  id: string;
  email: string;
  name: string;
  role: string;
}
