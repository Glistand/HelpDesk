export type TicketStatus =
  | "new"
  | "open"
  | "pending"
  | "resolved"
  | "closed";

export type TicketPriority = "low" | "normal" | "high" | "urgent";
export type TicketSource = "manager" | "bot";

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
  source: TicketSource;
  conversationId?: string;
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
  source: TicketSource;
  conversationId?: string;
}

export interface AuthUser {
  id: string;
  email: string;
  name: string;
  role: string;
}

export type ConversationStatus =
  | "bot"
  | "waiting_agent"
  | "open"
  | "resolved"
  | "unspecified";

export type MessageRole = "visitor" | "bot" | "agent" | "system" | "unspecified";

export interface ConversationSummary {
  id: string;
  siteKey: string;
  visitorId: string;
  status: ConversationStatus;
  assigneeId: string;
  preview: string;
  createdAt: string;
  updatedAt: string;
}

export interface ChatMessage {
  id: string;
  conversationId: string;
  role: MessageRole;
  body: string;
  authorId: string;
  createdAt: string;
}

export interface ConversationDetail {
  conversation: ConversationSummary;
  messages: ChatMessage[];
}
