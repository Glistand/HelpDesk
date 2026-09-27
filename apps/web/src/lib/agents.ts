import type { Agent, AuthUser } from "./types";

/** @deprecated Prefer session user from getMe(); kept for ticket mapper fallbacks. */
export const currentAgent: Agent = {
  id: "a-1",
  name: "Алексей К.",
  initials: "АК",
  role: "L1 Support",
};

export const agents: Agent[] = [currentAgent];

const byId = Object.fromEntries(agents.map((a) => [a.id, a]));

export function agentById(id: string | undefined | null): Agent | undefined {
  if (!id) return undefined;
  return byId[id];
}

export function initialsFromName(name: string): string {
  const parts = name.trim().split(/\s+/).filter(Boolean);
  if (parts.length === 0) return "?";
  if (parts.length === 1) return parts[0].slice(0, 2).toUpperCase();
  return (parts[0][0] + parts[1][0]).toUpperCase();
}

export function userToAgent(u: AuthUser): Agent {
  return {
    id: u.id,
    name: u.name,
    initials: initialsFromName(u.name),
    role: u.role,
  };
}

export function nameByUserId(
  id: string | undefined | null,
  users: AuthUser[],
): string {
  if (!id) return "Не назначен";
  const u = users.find((x) => x.id === id);
  return u?.name || id;
}
