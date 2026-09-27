import { AppShell } from "@/components/shell/app-shell";
import { listConversations, listTickets } from "@/lib/api";
import type { ConversationSummary, TicketSummary } from "@/lib/types";

export default async function DeskLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const [conversations, tickets] = await Promise.all([
    listConversations().catch((): ConversationSummary[] => []),
    listTickets().catch((): TicketSummary[] => []),
  ]);

  return (
    <AppShell conversations={conversations} tickets={tickets}>
      {children}
    </AppShell>
  );
}
