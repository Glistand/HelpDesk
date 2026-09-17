import { agentById } from "./agents";
import type {
  SlaState,
  TimelineEvent,
  TimelineKind,
  Ticket,
  TicketPriority,
  TicketStatus,
  TicketSummary,
} from "./types";

type ApiTicket = {
  id: string;
  title: string;
  description: string;
  status: string;
  priority: string;
  category: string;
  requester: string;
  assignee_id: string;
  created_at: string;
  updated_at: string;
};

type ApiSLA = {
  ticket_id: string;
  state: string;
  first_response_due: string;
  resolve_due: string;
  warned_at?: string;
  breached_at?: string;
  policy?: string;
} | null;

type ApiTimelineEvent = {
  id: string;
  ticket_id: string;
  event_id: string;
  event_type: string;
  title: string;
  detail: string;
  actor: string;
  occurred_at: string;
};

type ApiAssignment = {
  ticket_id: string;
  assignee_id: string;
  assignee_name: string;
  assigned_at: string;
} | null;

type ApiCard = {
  ticket: ApiTicket;
  assignment: ApiAssignment;
  sla: ApiSLA;
  timeline: ApiTimelineEvent[];
};

function asStatus(s: string): TicketStatus {
  switch (s) {
    case "new":
    case "open":
    case "pending":
    case "resolved":
    case "closed":
      return s;
    default:
      return "new";
  }
}

function asPriority(p: string): TicketPriority {
  switch (p) {
    case "low":
    case "normal":
    case "high":
    case "urgent":
      return p;
    default:
      return "normal";
  }
}

function asSlaState(s: string | undefined): SlaState {
  switch (s) {
    case "ok":
    case "warning":
    case "breached":
    case "cancelled":
      return s;
    default:
      return "ok";
  }
}

function timelineKind(eventType: string): TimelineKind {
  switch (eventType) {
    case "ticket.created":
      return "created";
    case "ticket.assigned":
      return "assigned";
    case "ticket.updated":
      return "status";
    case "sla.warned":
      return "sla_warn";
    case "sla.breached":
      return "sla_breach";
    case "ticket.escalated":
      return "escalated";
    case "notification.sent":
      return "notification";
    default:
      return "other";
  }
}

export function mapTicketSummary(
  t: ApiTicket,
  slaState: SlaState = "ok",
): TicketSummary {
  return {
    id: t.id,
    title: t.title,
    status: asStatus(t.status),
    priority: asPriority(t.priority),
    category: t.category,
    requester: t.requester,
    assigneeId: t.assignee_id || "",
    assignee: agentById(t.assignee_id),
    createdAt: t.created_at,
    updatedAt: t.updated_at,
    slaState,
  };
}

export function mapCard(card: ApiCard): Ticket {
  const assigneeId =
    card.assignment?.assignee_id || card.ticket.assignee_id || "";
  const assignee =
    agentById(assigneeId) ||
    (card.assignment
      ? {
          id: card.assignment.assignee_id,
          name: card.assignment.assignee_name,
          initials: card.assignment.assignee_name.slice(0, 2).toUpperCase(),
          role: "Agent",
        }
      : undefined);

  const timeline: TimelineEvent[] = (card.timeline || []).map((e) => ({
    id: e.id,
    kind: timelineKind(e.event_type),
    title: e.title,
    detail: e.detail || undefined,
    actor: e.actor || undefined,
    at: e.occurred_at,
  }));

  return {
    id: card.ticket.id,
    title: card.ticket.title,
    description: card.ticket.description,
    status: asStatus(card.ticket.status),
    priority: asPriority(card.ticket.priority),
    category: card.ticket.category,
    requester: card.ticket.requester,
    assignee,
    createdAt: card.ticket.created_at,
    updatedAt: card.ticket.updated_at,
    sla: {
      firstResponseDue: card.sla?.first_response_due || "",
      resolveDue: card.sla?.resolve_due || "",
      state: asSlaState(card.sla?.state),
      warnedAt: card.sla?.warned_at,
      breachedAt: card.sla?.breached_at,
      policy: card.sla?.policy,
    },
    timeline,
  };
}

export type { ApiTicket, ApiCard, ApiSLA };
