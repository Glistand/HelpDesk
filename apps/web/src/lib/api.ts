import { cookies } from "next/headers";
import { mapCard, mapTicketSummary, type ApiCard, type ApiTicket } from "./mappers";
import type { AuthUser, SlaState, Ticket, TicketSummary } from "./types";

export const AUTH_COOKIE = "hd_token";

export function gatewayURL(): string {
  return (
    process.env.GATEWAY_URL ||
    process.env.NEXT_PUBLIC_GATEWAY_URL ||
    "http://localhost:8080"
  );
}

async function token(): Promise<string | undefined> {
  const jar = await cookies();
  return jar.get(AUTH_COOKIE)?.value;
}

export class ApiError extends Error {
  status: number;
  constructor(status: number, message: string) {
    super(message);
    this.status = status;
  }
}

export async function apiFetch<T>(
  path: string,
  init: RequestInit = {},
): Promise<T> {
  const t = await token();
  const headers = new Headers(init.headers);
  headers.set("Content-Type", "application/json");
  if (t) headers.set("Authorization", `Bearer ${t}`);

  const res = await fetch(`${gatewayURL()}${path}`, {
    ...init,
    headers,
    cache: "no-store",
  });
  if (!res.ok) {
    let msg = res.statusText;
    try {
      const body = await res.json();
      msg = body.error || msg;
    } catch {
      /* ignore */
    }
    throw new ApiError(res.status, msg);
  }
  if (res.status === 204) return undefined as T;
  return res.json() as Promise<T>;
}

export async function login(
  email: string,
  password: string,
): Promise<{ access_token: string; user: AuthUser }> {
  const res = await fetch(`${gatewayURL()}/auth/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
    cache: "no-store",
  });
  if (!res.ok) {
    throw new ApiError(res.status, "invalid email or password");
  }
  const data = await res.json();
  return {
    access_token: data.access_token,
    user: data.user as AuthUser,
  };
}

export async function listTickets(opts?: {
  assigneeId?: string;
}): Promise<TicketSummary[]> {
  const q = new URLSearchParams();
  if (opts?.assigneeId) q.set("assignee_id", opts.assigneeId);
  const qs = q.toString();
  const data = await apiFetch<{ tickets: ApiTicket[] }>(
    `/tickets${qs ? `?${qs}` : ""}`,
  );
  const tickets = data.tickets || [];

  // Enrich with SLA state (parallel, capped).
  const slaStates = await Promise.all(
    tickets.slice(0, 50).map(async (t) => {
      try {
        const sla = await apiFetch<{ state: string }>(`/tickets/${t.id}/sla`);
        return [t.id, (sla.state as SlaState) || "ok"] as const;
      } catch {
        return [t.id, "ok" as SlaState] as const;
      }
    }),
  );
  const slaMap = Object.fromEntries(slaStates);
  return tickets.map((t) => mapTicketSummary(t, slaMap[t.id] || "ok"));
}

export async function getTicketCard(id: string): Promise<Ticket> {
  const card = await apiFetch<ApiCard>(`/tickets/${id}/card`);
  return mapCard(card);
}

export async function createTicket(input: {
  title: string;
  description: string;
  priority: string;
  category: string;
  requester?: string;
}): Promise<ApiTicket> {
  return apiFetch<ApiTicket>("/tickets", {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export async function updateTicketStatus(
  id: string,
  status: string,
): Promise<ApiTicket> {
  return apiFetch<ApiTicket>(`/tickets/${id}/status`, {
    method: "PATCH",
    body: JSON.stringify({ status }),
  });
}

export async function searchTickets(query: string): Promise<TicketSummary[]> {
  const q = encodeURIComponent(query);
  const data = await apiFetch<{
    hits: Array<{
      id: string;
      title: string;
      description: string;
      status: string;
      priority: string;
      category: string;
      requester: string;
      assignee_id: string;
      updated_at: string;
    }>;
  }>(`/search?q=${q}&limit=30`);
  return (data.hits || []).map((h) =>
    mapTicketSummary(
      {
        id: h.id,
        title: h.title,
        description: h.description,
        status: h.status,
        priority: h.priority,
        category: h.category,
        requester: h.requester,
        assignee_id: h.assignee_id,
        created_at: h.updated_at,
        updated_at: h.updated_at,
      },
      "ok",
    ),
  );
}
