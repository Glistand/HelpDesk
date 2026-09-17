import type { Agent } from "./types";

export const currentAgent: Agent = {
  id: "a-1",
  name: "Алексей К.",
  initials: "АК",
  role: "L1 Support",
};

export const agents: Agent[] = [
  currentAgent,
  { id: "a-2", name: "Марина С.", initials: "МС", role: "L1 Support" },
  { id: "a-3", name: "Денис В.", initials: "ДВ", role: "L2 Escalation" },
];

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
