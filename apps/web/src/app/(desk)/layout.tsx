import { AppShell } from "@/components/shell/app-shell";
import { listTickets } from "@/lib/api";
import type { TicketSummary } from "@/lib/types";

export default async function DeskLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  let tickets: TicketSummary[] = [];
  try {
    tickets = await listTickets();
  } catch {
    tickets = [];
  }
  return <AppShell tickets={tickets}>{children}</AppShell>;
}
